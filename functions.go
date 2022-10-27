package templateManager

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/grokify/html-strip-tags-go" // => strip
)

/*
Returns a function map for use with the Go template standard library
*/
func getDefaultFunctions() map[string]any {
	return map[string]any{
		"add":				add,
		"capfirst":			capfirst,
		"collection":		collection, 
		"contains":			contains,
		"cut":				cut,
		"date":				date,
		"datetime":			datetime,
		"default":			defaultVal,
		"divide":			divide,
		"divisibleby":		divisibleBy,
		"dl":				dl,
		"equal":			equal,
		"first":			first,
		"firstof":			firstOf,
		"formattime":		formattime,
		"htmldecode":		htmlDecode,
		"htmlencode":		htmlEncode,
		"join":				join,
		"jsondecode":		jsonDecode,
		"jsonencode":		jsonEncode,
		"key":				keyFn,
		"kind":				kind,
		"last":				last,
		"length":			length,
		"localtime":		localtime,
		"lower":			lower,
		"ltrim":			ltrim,
		"mktime":			mktime,
		"multiply":			multiply,
		"nl2br":			nl2br,
		"notequal":			notequal,
		"now":				now, 
		"ol":				ol,
		"ordinal":			ordinal,
		"paragraph":		paragraph,
		"pluralise":		pluralise,
		"prefix":			prefix,
		"random":			random,
		"regexp":			regexpFindAll,
		"regexpreplace":	regexpReplaceAll,
		"replace":			replaceAll,
		"round":			round,
		"rtrim":			rtrim,
		"split":			split,
		"striptags":		stripTags,
		"subtract": 		subtract,
		"suffix":			suffix,
		"time":				timeFn,
		"timesince":		timeSince,
		"timeuntil":		timeUntil,
		"title":			title,
		"trim":				trim,
		"truncate":			truncate,
		"truncatewords":	truncatewords,
		"type":				typeFn, 
		"ul":				ul,
		"upper":			upper,
		"urldecode":		urlDecode,
		"urlencode":		urlEncode,
		"wordcount":		wordcount,
		"wrap":				wrap,
		"year":				year,
		"yesno":			yesno,
	}
}

/*
 func add[T any](value T, to T) T
Adds a value to the existing item.
For numeric items this is a simple addition. For other types this is appended / merged as appropriate.
*/
func add(value reflect.Value, to reflect.Value) reflect.Value {
	sig := "add(value any, to any)"

	value	= reflectHelperUnpackInterface(value)
	to		= reflectHelperUnpackInterface(to)

	if !value.IsValid() {
		logError(sig + " value added cannot be an untyped nil value")
		return to
	}

	if !to.IsValid() {
		logError(sig + " value added to cannot be an untyped nil value")
		return to
	}

	// It's a simple type, do it recursively
	if reflectHelperIsNumeric(value) || value.Kind() == reflect.String {
		switch to.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
			reflect.Uint64:
				addVal, _ := reflectHelperConvertToFloat64(value)
				toVal, _ := reflectHelperConvertToFloat64(to)

				return reflect.ValueOf(int64(roundFloat(addVal + toVal, 0))).Convert(to.Type())
			case reflect.Float32, reflect.Float64:
				addVal, _ := reflectHelperConvertToFloat64(value)
				toVal, _ := reflectHelperConvertToFloat64(to)

				return reflect.ValueOf(addVal + toVal).Convert(to.Type())
			case reflect.String:
				addVal, _ := reflectHelperConvertToString(value)

				return reflect.ValueOf(to.String() + addVal)
		}

		return recursiveHelper(to, reflect.ValueOf(add), value)
	}

	// It's a more complex type, no recursion and stricter checks
	if err := reflectHelperLooseTypeCompatibility(value, to); err != nil {
		logError(sig + fmt.Sprintf(" the value and addition must have the same approximate types. trying to add %s to %s", value.Type(), to.Type()))
		return to
	}

	switch to.Kind() {
		case reflect.Slice:
			slice, _ := reflectHelperCreateEmptySlice(value)
			slice = reflect.AppendSlice(slice, to)
			slice = reflect.AppendSlice(slice, value)
			return slice
		case reflect.Array:
			slice, _ := reflectHelperCreateEmptySlice(value)
			for i := 0; i < to.Len(); i++ {
				slice = reflect.Append(slice, to.Index(i))
			}
			for i := 0; i < value.Len(); i++ {
				slice = reflect.Append(slice, value.Index(i))
			}
			arr, _ := reflectHelperConvertSliceToArray(slice)
			return arr
		case reflect.Map:
			tmp := reflect.MakeMap(to.Type())

			iter := to.MapRange()
			for iter.Next() {
				tmp.SetMapIndex(iter.Key(), iter.Value())
			}

			iter = value.MapRange()
			for iter.Next() {
				if val := tmp.MapIndex(iter.Key()); val.IsValid() {
					tmp.SetMapIndex(iter.Key(), add(iter.Value(), val))
				} else {
					tmp.SetMapIndex(iter.Key(), iter.Value())
				}
			}
			return tmp
	}

	return to
}

/*
 func capfirst[T any](value T) T
Capitalises the first letter of strings.
Does not alter any other letters.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func capfirst(value reflect.Value) reflect.Value {
	switch value.Kind() {
		case reflect.String:
			runes := []rune(value.String())
			for i, r := range runes {
				if unicode.IsLetter(r) {
					if ! unicode.IsTitle(r) {
						runes[i] = []rune(strings.ToUpper(string(r)))[0]
					}
					break
				}
			}
			return reflect.ValueOf(string(runes))
	}

	return recursiveHelper(value, reflect.ValueOf(capfirst))
}

/*
 func collection(key string, value any) map[string]any
Allows several variables to be packaged together into a map for passing to templates. 
*/
func collection(pairs ...any) map[string]any {
	sig := "collection(pairs ...any)"
	length := len(pairs)
	if length == 0 || length % 2 != 0 {
		logError(sig + " can only accept pairs of arguments (string / any)")
		return map[string]any{}
	}

	collection := make(map[string]any, length / 2)

	for i := 0; i < length; i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			logError(sig + " first member of a pair must be a string")
			return map[string]any{}
		}
		collection[key] = pairs[i + 1]
	}

	return collection
}

/*
 func contains(find any, within any) bool
Returns a boolean value to determine whether the `find` value is contained in the `within` value.
The `find` value can act on strings, slices, arrays and maps.
*/
func contains(find reflect.Value, within reflect.Value) bool {
	sig		:= "contains(find any, within string|slice|map)"
	find	= reflectHelperUnpackInterface(find)
	within	= reflectHelperUnpackInterface(within)

	if !find.IsValid() {
		logWarning(sig + " is trying to search for an untyped nil value")
		return false
	}

	if !within.IsValid() {
		logWarning(sig + " is trying to search within an untyped nil value")
		return false
	}

	switch within.Kind() {
		case reflect.String:
			if find.Kind() == reflect.String {
				return strings.Contains(within.String(), find.String())
			}
			logError(sig + fmt.Sprintf(" can't search within a string using a %s", find.Type()))
			return false
		case reflect.Array, reflect.Slice:
			if reflectHelperGetSliceType(within) == find.Type().String() {
				for i := 0; i < within.Len(); i++ {
					if reflect.DeepEqual(within.Index(i).Interface(), find.Interface()) {
						return true
					}
				}
			} else {
				logError(sig + fmt.Sprintf(" can't search within a slice type %s using a %s", reflectHelperGetSliceType(within), find.Type()))
			}
			return false
		case reflect.Map:
			if reflectHelperGetMapType(within) == find.Type().String() {
				iter := within.MapRange()
				for iter.Next() {
					if reflect.DeepEqual(iter.Value().Interface(), find.Interface()) {
						return true
					}
				}
			} else {
				logError(sig + fmt.Sprintf(" can't search within a map type %s using a %s", reflectHelperGetMapType(within), find.Type()))
			}
			return false
		case reflect.Struct:
			for i := 0; i < within.NumField(); i++ {
				field, err := reflectHelperGetStructValue(within, reflect.ValueOf(i))
				if err == nil {
					if field.Type().String() == find.Type().String() {
						if reflect.DeepEqual(field.Interface(), find.Interface()) {
							return true
						}
					}
				}
			}
			return false
		case reflect.Invalid:
			logWarning(sig + " invalid value passed")
			return false
		default:
			logWarning(sig + fmt.Sprintf(" can't search within an item of type %s", within.Type()))
			return false
	}
}

/*
 func cut[T any](remove string, from T) T
Will `remove` a string value that is contained in the `from` value.
If `from` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func cut(remove reflect.Value, from reflect.Value) reflect.Value {
	return replaceAll(remove, reflect.ValueOf(""), from)
}

/*
Returns a simple date string (by default: "d/m/Y").
Supports Go, Python and PHP formatting standards.
It can accept various parameter combinations:
 date()                                          // Current date and default output format
 date(time time.Time)                            // Passed in time and default output format
 date(format string)                             // Current date and custom output format
 date(format string, time time.Time)             // Time returned in the format specified
 date(format string, time string)                // Time in `time.RFC3339` format parsed into the format specified
                                                     // date "15:04" "2019-04-23T11:30:05Z"
 date(format string, layout string, time string) // Time with a custom layout rule specifying an output format
                                                     // date "15:04" "Jan 2, 2006 at 3:04pm (MST)" "Feb 3, 2013 at 7:54pm (PST)"
                                                     // date "H:i" "Y-m-d H:i:s (T)" "2013-02-03 19:54:00 (PST)"
*/
func date(params ...any) string {
	format := dateDefaultDateFormat

	if len(params) == 0 {
		return timeFn(format)
	} else if len(params) == 1 {
		switch val := params[0].(type) {
			case time.Time:
				return timeFn(format, val)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				return timeFn(format, params[0])
		}
	}
	
	return timeFn(params...)
}

/*
Returns a simple datetime string (by default: "d/m/Y H:i").
Supports Go, Python and PHP formatting standards.
It can accept various parameter combinations:
 datetime()                                          // Current date and default output format
 datetime(time time.Time)                            // Passed in time and default output format
 datetime(format string)                             // Current date and custom output format
 datetime(format string, time time.Time)             // Time returned in the format specified
 datetime(format string, time string)                // Time in `time.RFC3339` format parsed into the format specified
                                                         // datetime "02/01 15:04" "2019-04-23T11:30:05Z"
 datetime(format string, layout string, time string) // Time with a custom layout rule specifying an output format
                                                         // datetime "02/01 15:04" "1 2, 2006 at 3:04pm" "2 3, 2013 at 7:54pm"
                                                         // datetime "m/d H:i" "Y-m-d H:i:s (T)" "2013-02-03 19:54:00 (PST)"
*/
func datetime(params ...any) string {
	format := dateDefaultDatetimeFormat

	if len(params) == 0 {
		return timeFn(format)
	} else if len(params) == 1 {
		switch val := params[0].(type) {
			case time.Time:
				return timeFn(format, val)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				return timeFn(format, params[0])
		}
	}
	
	return timeFn(params...)
}

/*
 func defaultVal(def any, test any) any
Will return the second `test` value if it is not empty, else return the `def` value
*/
func defaultVal(def reflect.Value, test reflect.Value) reflect.Value {
	sig		:= "default(def any, value any)"
	def		= reflectHelperUnpackInterface(def)
	test	= reflectHelperUnpackInterface(test)

	if !def.IsValid() {
		logError(sig + " default cannot be an untyped nil value")
		return reflect.Value{}
	}

	if !test.IsValid() {
		return def
	}

	switch test.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			if test.Len() > 0 {
				return test
			}
		case reflect.Bool:
			if test.Bool() {
				return test
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64:
			if integer, err := reflectHelperConvertToInt(test); err == nil {
				if integer != 0 {
					return test
				}
			}
		case reflect.Struct:
			if !reflectHelperIsEmptyStruct(test) {
				return test
			}
		case reflect.Invalid:
			logWarning(sig + " invalid value passed")
	}

	return def
}

/*
 func divide[T any](divisor int|float, value T) T
Divides the `value` by the `divisor`.
If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain.
*/
func divide(divisor reflect.Value, value reflect.Value) reflect.Value {
	sig		:= "divide(divisor int, value any)"
	divisor	= reflectHelperUnpackInterface(divisor)
	value	= reflectHelperUnpackInterface(value)

	if !reflectHelperIsNumeric(divisor) {
		logError(sig + fmt.Sprintf(" divisor must be numeric, not %s", value.Type()))
		return value
	}

	div, _ := reflectHelperConvertToFloat64(divisor)
	if div == 0.0 {
		logError(sig + " divisor must not be zero")
		return value
	}

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64:
			val, _ := reflectHelperConvertToFloat64(value)
			op := val / div
			return reflect.ValueOf(int64(roundFloat(op, 0))).Convert(value.Type())
		case reflect.Float32, reflect.Float64:
			val, _ := reflectHelperConvertToFloat64(value)
			op := val / div
			return reflect.ValueOf(op).Convert(value.Type())
		case reflect.String, reflect.Bool:
			logWarning(sig + fmt.Sprintf(" trying to divide a %s", value.Type()))
			return value
	}

	return recursiveHelper(value, reflect.ValueOf(divide), divisor)
}

/*
 func divisibleby[T any](divisor int, value T) bool
Determines if the `value` is divisible by the `divisor`
*/
func divisibleBy(divisor reflect.Value, value reflect.Value) bool {
	sig		:= "divisibleby(divisor int, value any)"
	value	= reflectHelperUnpackInterface(value)

	if !reflectHelperIsNumeric(divisor) {
		logError(sig + fmt.Sprintf(" divisor must be numeric, not %s", value.Type()))
		return false
	}

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64:
			val, _ := reflectHelperConvertToFloat64(value)
			div, _ := reflectHelperConvertToFloat64(divisor)

			if div == 0.0 {
				logWarning(sig + " divisor must not be zero")
				return false
			}
			result := val / div

			return equalFloats(result, roundFloat(result, 0))
		default:
			logWarning(sig + " attempting division of non numeric type")
	}

	return false
}

/*
 func dl(value any) string
Converts slices, arrays or maps into an HTML definition list.
For maps this will use the keys as the dt elements.
*/
func dl(value reflect.Value) string {
	return listHelper(value, "dl")
}

/*
 func equal(values ...any) bool
Determines whether any values are equal.
*/
func equal(values ...reflect.Value) bool {
	sig := "equal(values ...any)"

	if len(values) < 2 {
		logWarning(sig + fmt.Sprintf(" at least two values required, %d provided", len(values)))
		return false
	}

	for i, value := range values {
		values[i] = reflectHelperUnpackInterface(value)

		if !value.IsValid() {
			logWarning(sig + " values cannot be untyped nil values")
			return false
		}

		if i > 0 {
			err := reflectHelperVeryLooseTypeCompatibility(value, values[i - 1])
			if err != nil {
				return false
			}

			switch value.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
					val1, _ := reflectHelperConvertToFloat64(value)
					val2, _ := reflectHelperConvertToFloat64(values[i - 1])

					if !equalFloats(val1, val2) {
						return false
					}
				case reflect.String:
					if value.String() != values[i - 1].String() {
						return false
					}
				case reflect.Bool:
					if value.Bool() != values[i - 1].Bool() {
						return false
					}
				case reflect.Array, reflect.Slice:
					if !reflect.DeepEqual(value.Slice(0, value.Len() - 1).Interface(), values[i - 1].Slice(0, values[i - 1].Len() - 1).Interface()) {
						return false
					}
				case reflect.Map, reflect.Struct:
					if !reflect.DeepEqual(value.Interface(), values[i - 1].Interface()) {
						return false
					}
				default:
					return false
			}
		}
	}

	return true
}

/*
 func first(value string|slice|array) any
Gets the first value from slices / arrays or the first word from strings.
*/
func first(value reflect.Value) reflect.Value {
	sig		:= "first(value string|slice)"
	value	= reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		logError(sig + " value cannot be an untyped nil value")
		return reflect.Value{}
	}

	switch value.Kind() {
		case reflect.String:
			if value.Len() > 0 {
				str := strings.Split(strings.TrimLeft(value.String(), " \n\r\t"), " ")[0]
				return reflect.ValueOf(str)
			}
		case reflect.Array, reflect.Slice:
			if value.Len() > 0 {
				return value.Index(0)
			}
		case reflect.Struct:
			value, err := reflectHelperGetStructValue(value, reflect.ValueOf(0))
			if err != nil {
				logError(sig + " " + err.Error())
				return reflect.Value{}
			}
			return value
		case reflect.Invalid:
			logError(sig + " invalid value passed")
		default:
			logError(sig + fmt.Sprintf(" can't handle items of type %s", value.Type()))
	}

	return reflect.Value{}
}

/*
 func firstOf(values ...any) any
Accepts any number of values and returns the first one of them that exists and is not empty.
*/
func firstOf(values ...reflect.Value) reflect.Value {
	if len(values) < 1 {
		logWarning("firstof(values ...any) being called without any parameters")
		return reflect.Value{}
	}

	for _, value := range values {
		value = reflectHelperUnpackInterface(value)

		if !value.IsValid() {
			continue
		}

		switch value.Kind() {
			case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
				if value.Len() > 0 {
					return value
				}
			case reflect.Struct:
				empty := reflect.New(value.Type()).Elem().Interface()
				if !reflect.DeepEqual(value.Interface(), empty) {
					return value
				}
			case reflect.Bool:
				if value.Bool() {
					return value
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
			reflect.Uint64, reflect.Float32, reflect.Float64:
				if integer, err := reflectHelperConvertToInt(value); err == nil {
					if integer != 0 {
						return value
					}
				}
		}
	}

	return reflect.Value{}
}

/*
 func formattime(format string, t time.Time) string

Formats a time.Time object for display.
*/
func formattime(format string, t time.Time) string {
	return t.Format(dateFormatHelper(format))
}

/*
 func htmlDecode[T any](value T) T
Converts HTML character-entity equivalents back into their literal, usable forms.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func htmlDecode(value reflect.Value) reflect.Value {
	switch value.Kind() {
		case reflect.String:
			find	:= []string{ "&lt;", "&gt;", "&#34;", "&#x22;", "&quot;", "&#39;", "&#x27;", "&amp;" }
			replace	:= []string{ "<",    ">",    `"`,     `"`,      `"`,      "'",     "'",      "&" }
			replacer, err := replaceHelper(find, replace)
			if err != nil {
				logError(err.Error())
				return reflect.Value{}
			}
			return reflect.ValueOf(replacer.Replace(value.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(htmlDecode))
}

/*
 func htmlEncode[T any](value T) T
Converts literal HTML special characters into safe, character-entity equivalents.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func htmlEncode(value reflect.Value) reflect.Value {
	switch value.Kind() {
		case reflect.String:
			find	:= []string{ "<",    ">",    `"`,     `"`,      `"`,      "'",     "'",      "&" }
			replace	:= []string{ "&lt;", "&gt;", "&#34;", "&#x22;", "&quot;", "&#39;", "&#x27;", "&amp;" }
			replacer, err := replaceHelper(find, replace)
			if err != nil {
				logError(err.Error())
				return reflect.Value{}
			}
			return reflect.ValueOf(replacer.Replace(value.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(htmlEncode))
}

/*
 func join(separator string, values any) string
Joins slice or map `values` together as a string spaced by the `separator`.
*/
func join(separator string, values reflect.Value) string {
	sig		:= "join(separator string, values any)"
	values	= reflectHelperUnpackInterface(values)
	str		:= ""

	if !values.IsValid() {
		logError(sig + " is trying to search within an untyped nil value")
		return ""
	}

	switch values.Kind() {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool:
			return fmt.Sprintf("%v", values)
		case reflect.Array, reflect.Slice:
			for i := 0; i < values.Len(); i++ {
				if i > 0 { str += separator }
				str += join(separator, values.Index(i))
			}
		case reflect.Map:
			keys, err := reflectHelperMapSort(values)
			if err == nil {
				for i := 0; i < keys.Len(); i++ {
					if i > 0 { str += separator }
					str += join(separator, values.MapIndex(keys.Index(i)))
				}
			} else {
				iter := values.MapRange()
				i := 0
				for iter.Next() {
					if i > 0 { str += separator }
					str += join(separator, iter.Value())
					i++
				}
			}
		case reflect.Struct:
			for i := 0; i < values.NumField(); i++ {
				if i > 0 { str += separator }
				str += join(separator, values.Field(i))
			}
		case reflect.Invalid:
			logError(sig + " invalid value passed")
			return ""
		default:
			logError(sig + fmt.Sprintf(" can't join items of type %s", values.Type()))
			return ""
	}

	return str
}

/*
 func jsonDecode(value any) string
Decodes any JSON string value to a map.
*/
func jsonDecode(value string) any {
	var result any
	err := json.Unmarshal([]byte(value), &result)
     
    if err != nil {
        logError(err.Error())
    }
     
    return result
}

/*
 func jsonEncode(value any) string
Encodes any value to a JSON string.
*/
func jsonEncode(value any) string {
	result, err := json.Marshal(value)
     
    if err != nil {
        logError(err.Error())
		return ""
    }
     
    return string(result)
}

/*
 func keyFn(input ...any) any
Accepts any number of nested keys and returns the result of indexing its final argument by them.
For strings this returns a byte value.
The indexed item must be a string, map, slice, or array.
*/
func keyFn(input ...reflect.Value) reflect.Value {
	sig := "key(indexes ...any, value any)"
	if len(input) < 1 {
		logError(sig + " requires at least one argument")
		return reflect.Value{}
	}

	if len(input) == 1 {
		logWarning(sig + " requires at least two arguments to run as intended")
		return input[0]
	}

	value	:= reflectHelperUnpackInterface(input[len(input) - 1])
	indexes	:= input[:len(input) - 1]
	var nilPointer bool
	var err error

	if !value.IsValid() {
		logError(sig + " is trying to access an untyped nil value")
		return reflect.Value{}
	}

	for _, index := range indexes {
		index = reflectHelperUnpackInterface(index)

		if value, nilPointer = reflectHelperCheckNilPointers(value); nilPointer {
			logError(sig + " is trying to access an index of a nil pointer")
			return reflect.Value{}
		}

		switch value.Kind() {
			case reflect.Array, reflect.Slice, reflect.String:
				origKind := value.Kind()
				value, err = reflectHelperGetSliceValue(value, index)
				if err != nil {
					logError(sig + " " + err.Error())
					return value
				}
				if origKind == reflect.String && value.Kind() == reflect.Uint8 {
					value = reflect.ValueOf(string(uint8(value.Uint())))
				}
			case reflect.Map:
				value, err = reflectHelperGetMapValue(value, index)
				if err != nil {
					logError(sig + " " + err.Error())
					return value
				}
			case reflect.Struct:
				value, err = reflectHelperGetStructValue(value, index)
				if err != nil {
					logError(sig + " " + err.Error())
					return value
				}
			case reflect.Invalid:
				logError(sig + " invalid value passed")
				return reflect.Value{}
			default:
				logError(sig + " " + fmt.Sprintf("can't index item of type %s", value.Type()))
				return reflect.Value{}
		}
	}

	return value
}

/*
 func kind[T any](value T) string
Returns a string representation of the reflection Kind
*/
func kind(value reflect.Value) string {
	return value.Kind().String()
}

/*
 func last(value string|slice|array) any
Gets the last value from slices / arrays or the last word from strings.
*/
func last(value reflect.Value) reflect.Value {
	sig		:= "last(value string|slice)"
	value	= reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		logError(sig + " value cannot be an untyped nil value")
		return reflect.Value{}
	}

	switch value.Kind() {
		case reflect.String:
			if value.Len() > 0 {
				sl := strings.Split(strings.TrimRight(value.String(), " \n\r\t"), " ")
				str := sl[len(sl) - 1]
				return reflect.ValueOf(str)
			}
		case reflect.Array, reflect.Slice:
			if value.Len() > 0 {
				return value.Index(value.Len() - 1)
			}
		case reflect.Struct:
			value, err := reflectHelperGetStructValue(value, reflect.ValueOf(value.NumField() - 1))
			if err != nil {
				logError(sig + " " + err.Error())
				return reflect.Value{}
			}
			return value
		case reflect.Invalid:
			logError(sig + " invalid value passed")
		default:
			logError(sig + fmt.Sprintf(" can't handle items of type %s", value.Type()))
	}

	return reflect.Value{}
}

/*
 func length(value any) int
Gets the length of any type without panics.
*/
func length(value reflect.Value) int {
	sig		:= "length(value any)"
	value	= reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		return 0
	}

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return len(fmt.Sprintf("%v", value))
		case reflect.Float32, reflect.Float64:
			return len(fmt.Sprintf("%#v", value))
		case reflect.Bool:
			return 1
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			return value.Len()
		case reflect.Struct:
			return value.NumField()
		default:
			logError(sig + fmt.Sprintf(" can't handle items of type %s", value.Type()))
	}

	return 0
}

/*
 func localtime(location string|time.Location, t time.Time) time.Time
Localises a time.Time object to display local times / dates.
*/
func localtime(location any, t time.Time) time.Time {
	var tz *time.Location
	switch v := location.(type) {
		case time.Location:
			tz = &v
		case string:
			tmp, err := time.LoadLocation(v)
			if err != nil {
				return t
			}
			tz = tmp
	}

	return t.In(tz)
}

/*
 func lower[T any](value T) T
Converts string text to lower case.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func lower(value reflect.Value) reflect.Value {
	switch value.Kind() {
		case reflect.String: return reflect.ValueOf(strings.ToLower(value.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(lower))
}

/*
 func ltrim[T any](remove string, value T) T
Removes the passed characters from the left end of string values.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func ltrim(remove reflect.Value, value reflect.Value) reflect.Value {
	sig := "ltrim(remove string, value any)"
	if remove.Kind() != reflect.String {
		logError(sig + " remove can only be a string")
		return reflect.Value{}
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.TrimLeft(value.String(), remove.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(ltrim), remove)
}

/*
The mktime() function creates new `time.Time` struct from simple time strings.
Returns the current time if an invalid input is given. 
Supports Go, Python and PHP formatting standards.
It can accept various parameter combinations:
 mktime()                           // Current time
 mktime(time string)                // Parse from a `time.RFC3339` formatted string
                                        // mktime "2019-04-23T11:30:05Z"
 mktime(layout string, time string) // Parse from a custom formatted string using the given layout
                                        // mktime "2006-01-02T15:04:05Z07:00" "2019-04-23T11:30:05Z"
                                        // mktime "Y-m-d\\TH:i:sZ" "2019-04-23T11:30:05Z"
                                        // mktime "MYSQL" "2019-04-23 11:30:05"
*/
func mktime(params ...string) time.Time {
	sig	:= "mktime(params ...string)"
	t	:= time.Now()

	if len(params) == 1 {
		tmp, err := time.Parse(time.RFC3339, params[0])
		if err != nil {
			logError(sig + " Invalid RFC3339 (\"" + time.RFC3339 + "\") passed: mktime(\"" + params[0] + "\")")
			return t.In(dateLocalTimezone)
		}
		t = tmp
	} else if len(params) == 2 {
		tmp, err := time.Parse(dateFormatHelper(params[0]), params[1])
		if err != nil {
			logError(sig + "Invalid date / format passed: mktime(\"" + params[0] + "\", \"" + params[1] + "\")")
			return t.In(dateLocalTimezone)
		}
		t = tmp
	}

	return t.In(dateLocalTimezone)
}

/*
 func multiply[T any](multiplier int|float, value T) T
Multiplies the `value` by the `multiplier`.
If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain.
*/
func multiply(multiplier reflect.Value, value reflect.Value) reflect.Value {
	sig			:= "multiply(multiplier int, value any)"
	multiplier	= reflectHelperUnpackInterface(multiplier)
	value		= reflectHelperUnpackInterface(value)

	if !reflectHelperIsNumeric(multiplier) {
		logError(sig + fmt.Sprintf(" multiplier must be numeric, not %s", value.Type()))
		return value
	}

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64:
			val, _ := reflectHelperConvertToFloat64(value)
			mul, _ := reflectHelperConvertToFloat64(multiplier)
			op := val * mul
			return reflect.ValueOf(int64(roundFloat(op, 0))).Convert(value.Type())
		case reflect.Float32, reflect.Float64:
			val, _ := reflectHelperConvertToFloat64(value)
			mul, _ := reflectHelperConvertToFloat64(multiplier)
			op := val * mul
			return reflect.ValueOf(op).Convert(value.Type())
		case reflect.String, reflect.Bool:
			logWarning(sig + fmt.Sprintf(" trying to multiply a %s", value.Type()))
			return value
	}

	return recursiveHelper(value, reflect.ValueOf(multiply), multiplier)
}

/*
 func nl2br[T any](value T) T
Replaces all instances of "\n" (new line) with instances of "<br>" within `value`.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func nl2br(value reflect.Value) reflect.Value {
	return replaceAll(reflect.ValueOf("\n"), reflect.ValueOf("<br>"), value)
}

/*
 func notequal(values ...any) bool
Determines whether any values are not equal.
*/
func notequal(values ...reflect.Value) bool {
	return !equal(values...)
}

/*
 func now() time.Time
Returns the current `time.Time` value
*/
func now() time.Time {
	return time.Now().In(dateLocalTimezone)
}

/*
 func ol(value any) string
Converts slices, arrays or maps into an HTML ordered list.
*/
func ol(value reflect.Value) string {
	return listHelper(value, "ol")
}

/*
 func ordinal[T int|float64|string](value T) string
Suffixes a number with the correct English ordinal
*/
func ordinal(value reflect.Value) string {
	sig		:= "ordinal(value int)"
	value	= reflectHelperUnpackInterface(value)

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64:
			suffix := "th"
			integer, _ := reflectHelperConvertToInt(value)
			switch integer % 10 {
				case 1:
					if integer % 100 != 11 {
						suffix = "st"
					}
				case 2:
					if integer % 100 != 12 {
						suffix = "nd"
					}
				case 3:
					if integer % 100 != 13 {
						suffix = "rd"
					}
			}
			return strconv.Itoa(integer) + suffix
		default:
			logError(sig + " attempting an ordinal conversion on a non numeric type")
	}

	return ""
}

/*
 func paragraph[T any](value T) T
Replaces all instances of "\n+" (multiple new lines) with paragraphs and instances of "\n" (new line) with instances of "<br>" within `value`
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func paragraph(value reflect.Value) reflect.Value {
	value = replaceAll(reflect.ValueOf("\r\n"), reflect.ValueOf("\n"), value)
	value = replaceAll(reflect.ValueOf("\r"), reflect.ValueOf("\n"), value)
	value = regexpReplaceAll(reflect.ValueOf("(\\s*\\n\\s*){2,}"), reflect.ValueOf("</p><p>"), value)
	value = replaceAll(reflect.ValueOf("\n"), reflect.ValueOf("<br>"), value)
	value = wrap(reflect.ValueOf("<p>"), reflect.ValueOf("</p>"), value)

	return value
}

/*
Allows pluralisation of word endings.
Allows basic customisation of the possible plural forms.
 // Returns empty string for `count` == 1 and "s" for `count` != 1
 pluralise(count int)

 // Returns empty string for `count` == 1 and `plural` for `count` != 1
 pluralise(plural string, count int)

 // Returns `singular` for `count` == 1 and `plural` for `count` != 1
 pluralise(singular string, plural string, count int)
*/
func pluralise(values ...any) string {
	if len(values) < 1 {
		logError("pluralise(): called without argument")
		return ""
	}

	num := 1
	suffixSingular := ""
	suffixPlural := "s"

	if len(values) == 1 {
		switch v := values[0].(type) {
			case int: num = v
			default: 
				logError("pluralise(int): single value should be an integer")
				return ""
		}
	} else if len(values) == 2 {
		switch v := values[0].(type) {
			case string: suffixPlural = v
			default: 
				logError("pluralise(string, int): first value should be a string")
				return ""
		}

		switch v := values[1].(type) {
			case int: num = v
			default: 
				logError("pluralise(string, int): final value should be an integer")
				return ""
		}
	} else if len(values) == 3 {
		switch v := values[0].(type) {
			case string: suffixSingular = v
			default: 
				logError("pluralise(string, string int): first value should be a string")
				return ""
		}

		switch v := values[1].(type) {
			case string: suffixPlural = v
			default: 
				logError("pluralise(string, string, int): second value should be a string")
				return ""
		}

		switch v := values[2].(type) {
			case int: num = v
			default: 
				logError("pluralise(string, string, int): final value should be an integer")
				return ""
		}
	}

	if num == 1 {
		return suffixSingular
	}

	return suffixPlural
}

/*
 func prefixFn[T any](prefix string, value T) T
Prefixes all strings within `value` with `prefix`
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func prefix(prefixValue reflect.Value, value reflect.Value) reflect.Value {
	return wrap(prefixValue, reflect.ValueOf(""), value)
}

/*
Generates random numbers
 random()                 // Returns a random number between 0 and 10000
 random(limit int)        // Returns a random number between 0 and `limit`
 random(min int, max int) // Returns a random number between `min` and `max`
*/
func random(values ...int) int {
	rand.Seed(time.Now().UnixNano())

	if len(values) < 1 {
		return rand.Intn(10000)
	} else if len(values) == 1 {
		return rand.Intn(values[0])
	}

	min := values[0]
	max := values[1]
	if min > max {
		min = values[1]
		max = values[0]
	} else if min == max {
		return min
	}

	return rand.Intn(max - min) + min
}

/*
 func regexpFindAll(find string, value string) [][]string
Finds all instances of `find` regexp within `value`.
It ONLY acts on strings
*/
func regexpFindAll(find reflect.Value, value reflect.Value) reflect.Value {
	sig		:= "regexp(find string, value string)"
	value	= reflectHelperUnpackInterface(value)

	if find.Kind() != reflect.String {
		logError(sig + " can only find string values")
		return reflect.ValueOf([][]string{})
	}

	findRegexp, err := regexp.Compile(find.String())
	if err != nil {
		logError(sig + " invalid regexp: " + find.String())
		return reflect.ValueOf([][]string{})
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(findRegexp.FindAllStringSubmatch(value.String(), -1))
	}

	return reflect.ValueOf([][]string{})
}

/*
 func regexpReplaceAll[T any](find string, replace string, value T) T
Replaces all instances of `find` regexp with instances of `replace` within `value`
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func regexpReplaceAll(find reflect.Value, replace reflect.Value, value reflect.Value) reflect.Value {
	sig		:= "regexpreplace(find string, replace string, value any)"
	value	= reflectHelperUnpackInterface(value)

	if find.Kind() != reflect.String {
		logError(sig + " can only find string values")
		return value
	}

	if replace.Kind() != reflect.String {
		logError(sig + " can only replace string values")
		return value
	}

	findRegexp, err := regexp.Compile(find.String())
	if err != nil {
		logError(sig + " invalid regexp: " + find.String())
		return value
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(findRegexp.ReplaceAllString(value.String(), replace.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(regexpReplaceAll), find, replace)
}

/*
 func replaceAll[T any](find string, replace string, value T) T
Replaces all instances of `find` with instances of `replace` within `value`
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func replaceAll(find reflect.Value, replace reflect.Value, value reflect.Value) reflect.Value {
	sig		:= "replace(find string, replace string, value any)"
	value	= reflectHelperUnpackInterface(value)

	if find.Kind() != reflect.String {
		logError(sig + " can only find string values")
		return value
	}

	if replace.Kind() != reflect.String {
		logError(sig + " can only replace string values")
		return value
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.ReplaceAll(value.String(), find.String(), replace.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(replaceAll), find, replace)
}

/*
 func round[T any](precision int, value T) T
Rounds any floats to the required precision.
If `value` is a slice, array or map it will apply this conversion to any float elements that they contain.
*/
func round(precision reflect.Value, value reflect.Value) reflect.Value {
	sig := "round(precision int, value any)"
	if precision.Kind() != reflect.Int {
		logError(sig + " precision can only be an integer")
		return value
	}

	switch value.Kind() {
		case reflect.Float32, reflect.Float64:
			val, _ := reflectHelperConvertToFloat64(value)
			return reflect.ValueOf(roundFloat(val, uint(precision.Int()))).Convert(value.Type())
	}

	return recursiveHelper(value, reflect.ValueOf(round), precision)
}

/*
 func rtrim[T any](remove string, value T) T
Removes the passed characters from the right end of string values.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func rtrim(remove reflect.Value, value reflect.Value) reflect.Value {
	sig := "rtrim(remove string, value any)"
	if remove.Kind() != reflect.String {
		logError(sig + " remove can only be a string")
		return reflect.Value{}
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.TrimRight(value.String(), remove.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(rtrim), remove)
}

/*
 func split(separator string, value string) []string
Splits strings on the `separator` value and returns a slice of the pieces.
*/
func split(separator reflect.Value, value reflect.Value) reflect.Value {
	sig := "split(separator string, value any)"
	if separator.Kind() != reflect.String {
		logError(sig + " separator can only be a string")
		return reflect.Value{}
	}

	switch value.Kind() {
		case reflect.String:
			tmp := []string{}
			for _, val := range strings.Split(value.String(), separator.String()) {
				if val != "" {
					tmp = append(tmp, val)
				}
			}
			return reflect.ValueOf(tmp)
	}

	//return recursiveHelper(value, reflect.ValueOf(split), separator)
	return reflect.Value{}
}

/*
 func stripTags[T any](value T) T
Strips HTML tags from strings.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func stripTags(value reflect.Value) reflect.Value {
	switch value.Kind() {
		case reflect.String: return reflect.ValueOf(strip.StripTags(value.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(stripTags))
}

/*
 func subtract[T any](value T, to T) T
Removes a value from the existing item.
For numeric items this is a simple subtraction. For other types this is removed as appropriate.
*/
func subtract(value reflect.Value, from reflect.Value) reflect.Value {
	sig := "subtract(value any, from any)"

	value	= reflectHelperUnpackInterface(value)
	from	= reflectHelperUnpackInterface(from)

	if !value.IsValid() {
		logError(sig + " value subtracted cannot be an untyped nil value")
		return from
	}

	if !from.IsValid() {
		logError(sig + " value subtracted from cannot be an untyped nil value")
		return from
	}

	// It's a simple type, do it recursively
	if reflectHelperIsNumeric(value) || value.Kind() == reflect.String {
		switch from.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
			reflect.Uint64:
				subVal, _ := reflectHelperConvertToFloat64(value)
				fromVal, _ := reflectHelperConvertToFloat64(from)

				return reflect.ValueOf(int64(roundFloat(fromVal - subVal, 0))).Convert(from.Type())
			case reflect.Float32, reflect.Float64:
				subVal, _ := reflectHelperConvertToFloat64(value)
				fromVal, _ := reflectHelperConvertToFloat64(from)

				return reflect.ValueOf(fromVal - subVal).Convert(from.Type())
			case reflect.String:
				subVal, _ := reflectHelperConvertToString(value)

				return cut(reflect.ValueOf(subVal), from)
		}

		return recursiveHelper(from, reflect.ValueOf(subtract), value)
	}

	if err := reflectHelperLooseTypeCompatibility(value, from); err != nil {
		logError(sig + fmt.Sprintf(" the value and subtraction must have the same types. trying to remove %s from %s", value.Type(), from.Type()))
		return from
	}

	switch from.Kind() {
		case reflect.Slice, reflect.Array:
			slice, _ := reflectHelperCreateEmptySlice(value)
			var found bool
			for i := 0; i < from.Len(); i++ {
				found = false
				for j := 0; j < value.Len(); j++ {
					if reflect.DeepEqual(from.Index(i).Interface(), value.Index(j).Interface()) {
						found = true
						break
					}
				}
				if !found {
					slice = reflect.Append(slice, from.Index(i))
				}
			}
			if from.Kind() == reflect.Array {
				slice, _ = reflectHelperConvertSliceToArray(slice)
			}
			return slice
		case reflect.Map:
			tmp := reflect.MakeMap(from.Type())
			iter := from.MapRange()
			for iter.Next() {
				if val := value.MapIndex(iter.Key()); val.IsValid() {
					subtracted := subtract(val, iter.Value())
					if !reflectHelperIsEmpty(subtracted) {
						tmp.SetMapIndex(iter.Key(), subtracted)
					}
				} else {
					tmp.SetMapIndex(iter.Key(), iter.Value())
				}
			}
			return tmp
	}

	return from
}

/*
 func suffixFn[T any](suffix string, value T) T
Prefixes all strings within `value` with `prefix`
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func suffix(suffixValue reflect.Value, value reflect.Value) reflect.Value {
	return wrap(reflect.ValueOf(""), suffixValue, value)
}

/*
Returns a simple time string (by default: "HH:MM")
Supports Go, Python and PHP formatting standards.
It can accept various argument combinations:
 time()                                          // Current time and default output format
 time(time time.Time)                            // Passed in time and default output format
 time(format string)                             // Current time and custom output format
 time(format string, time time.Time)             // Time returned in the format specified
 time(format string, time string)                // Time in `time.RFC3339` format parsed into the format specified
                                                     // time "15:04" "2019-04-23T11:30:05Z"
 time(format string, layout string, time string) // Time with a custom layout rule specifying an output format
                                                     // time "15:04" "Jan 2, 2006 at 3:04pm (MST)" "Feb 3, 2013 at 7:54pm (PST)"
                                                     // time "H:i" "Y-m-d H:i:s (T)" "2013-02-03 19:54:00 (PST)"
*/
func timeFn(params ...any) string {
	sig		:= "time(params ...time.Time|string)"
	t		:= time.Now()
	f		:= dateFormatHelper(dateDefaultTimeFormat)

	if len(params) == 1 {
		switch val := params[0].(type) {
			case time.Time: t = val
			case string: f = dateFormatHelper(val)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				num, _ := interfaceHelperConvertToInt64(val)
				t = time.Unix(num, 0)
		}
	} else if len(params) == 2 {
		f = dateFormatHelper(params[0].(string))
		switch val := params[1].(type)  {
			case time.Time:
				t = val
			case string:
				tmp, err := time.Parse(time.RFC3339, val)
				t = tmp
				if err != nil {
					logError(sig + " Invalid RFC3339 date passed: time(\"" + f + "\", \"" + val + "\")")
					return ""
				}
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				num, _ := interfaceHelperConvertToInt64(val)
				t = time.Unix(num, 0)
		}
	} else if len(params) > 2 {
		f = dateFormatHelper(params[0].(string))
		l := dateFormatHelper(params[1].(string))
		tmp, err := time.Parse(l, params[2].(string))
		if err != nil {
			logError(sig + " Invalid date / format passed to time in template: time(\"" + f + "\", \"" + params[1].(string) + "\", \"" + params[2].(string) + "\")\n" + err.Error())
			return ""
		}
		if strings.Contains(l, "MST") {
			location, err := time.LoadLocation(tmp.Location().String())
			if err != nil {
				logError(err.Error())
				return ""
			}
			tmp, _ = time.ParseInLocation(l, params[2].(string), location)
		}
		t = tmp
	}

	return t.In(dateLocalTimezone).Format(f)
}

/*
 func timeSince(t time.Time) map[string]int
Calculates the approximate duration since the `time.Time` value.
The map of integers contains the keys: `years`, `weeks`, `days`, `hours`, `minutes`, `seconds`
*/
func timeSince(t time.Time) map[string]int {
	return formatDuration(time.Since(t))
}

/*
 func timeUntil(t time.Time) map[string]int
Calculates the approximate duration until the `time.Time` value.
The map of integers contains the keys: `years`, `weeks`, `days`, `hours`, `minutes`, `seconds`
*/
func timeUntil(t time.Time) map[string]int {
	return formatDuration(time.Until(t))
}

/*
 func title[T any](value T) T
Converts string text to title case.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func title(value reflect.Value) reflect.Value {
	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.Title(strings.ToLower(value.String())))
	}

	return recursiveHelper(value, reflect.ValueOf(title))
}

/*
 func trim[T any](remove string, value T) T
Removes the passed characters from the ends of string values.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func trim(remove reflect.Value, value reflect.Value) reflect.Value {
	sig := "trim(remove string, value any)"
	if remove.Kind() != reflect.String {
		logError(sig + " remove can only be a string")
		return value
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.Trim(value.String(), remove.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(trim), remove)
}

/*
 func truncate[T any](length int, value T) T
Truncates strings to a certain number of characters. It is multi-byte safe.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func truncate(length reflect.Value, value reflect.Value) reflect.Value {
	sig := "truncate(length int, value any)"
	if !reflectHelperIsNumeric(length) {
		logError(sig + " length can only be a number")
		return value
	}

	intLength, _ := reflectHelperConvertToInt(length)

	switch value.Kind() {
		case reflect.String:
			if intLength <= 0 {
				return reflect.ValueOf("")
			}

			runes := []rune(value.String())
			
			if stringLength := len(runes); intLength >= stringLength {
				return value
			}

			if !strings.Contains(value.String(), "<") && !strings.Contains(value.String(), "&") {
				return reflect.ValueOf(string(runes[:intLength]))
			}

			// We have HTML
			singleTags := map[string]int8{"area":1, "base":1, "br":1, "col":1, "embed":1, "hr":1, "img":1, "input":1, "keygen":1, "link":1, "meta":1, "param":1, "source":1, "track":1, "wbr":1}
			history := []string{}
			count := 0
			all := 0
			tag := ""
			entityLength := 0
			skip := false
			var attribute rune = 0
			var waitFor rune = 0
			for i, r := range runes {

				if waitFor == 0 && r == '<' {
					waitFor = '>'
					tag += string(r)
				} else if waitFor == 0 && r == '&' {
					waitFor = ';'
				} else if waitFor == r {
					if attribute == 0 {
						waitFor = 0
						skip = true
					}
				}

				// Incorrect usage of > in attributes
				if waitFor == '>' && (r == '"' || r == '\'') {
					var prev rune = 0
					for j := (i - 1); j >= 0; j-- {
						if runes[j] != ' ' {
							prev = runes[j]
							break
						}
					}

					if attribute == 0 && prev == '=' {
						attribute = r
					} else if attribute == r && prev != '\\' {
						attribute = 0
					}
				}

				// Incorrect usage of raw ampersand or non-terminated entities within HTML (far from perfect)
				if waitFor == ';' {
					entityLength++
					if r == ' ' || r == '&' {
						count += entityLength - 1
						waitFor = 0
						entityLength = 0

						if r == '&' {
							waitFor = ';'
							entityLength++
						}
					}
				}

				if len(tag) > 0 {
					tag += string(r)

					if waitFor == 0 {
						tag = strings.Split(strings.Trim(tag, "<>"), " ")[0]
						if tag[0:1] == "/" && len(history) > 0 && tag[1:] == history[len(history) - 1] {
							history = history[:len(history) - 1]
						} else {
							if _, ok := singleTags[tag]; !ok {
								history = append(history, tag)
							}
						}
						tag = ""
					}
				}
				
				all++
				if waitFor == 0 && !skip {
					count++

					if count >= intLength {
						break
					}
				}

				skip = false
			}

			output := string(runes[:all])

			if historyLength := len(history); historyLength > 0 {
				for i := historyLength - 1; i >= 0; i-- {
					output += "</" + history[i] + ">"
				}
			}

			return reflect.ValueOf(output)
	}

	return recursiveHelper(value, reflect.ValueOf(truncate), length)
}

/*
 func truncatewords[T any](length int, value T) T
Truncates strings to a certain number of words. It is multi-byte safe.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func truncatewords(length reflect.Value, value reflect.Value) reflect.Value {
	sig := "truncatewords(length int, value any)"
	if !reflectHelperIsNumeric(length) {
		logError(sig + " length can only be a number")
		return value
	}

	intLength, _ := reflectHelperConvertToInt(length)

	switch value.Kind() {
		case reflect.String:
			if intLength <= 0 {
				return reflect.ValueOf("") 
			}

			runes := []rune(value.String())

			singleTags := map[string]int8{"area":1, "base":1, "br":1, "col":1, "embed":1, "hr":1, "img":1, "input":1, "keygen":1, "link":1, "meta":1, "param":1, "source":1, "track":1, "wbr":1}
			wordBreaks := map[rune]int8{' ':1, '\n':1, '\r':1, '\t':1, '.':1, '?':1, '!':1, ',':1}
			history := []string{}
			count := 0
			all := 0
			tag := ""
			entityLength := 0
			skip := false
			wordbreak := false
			var attribute rune = 0
			var waitFor rune = 0
			for i, r := range runes {

				if waitFor == 0 && r == '<' {
					waitFor = '>'
					tag += string(r)
				} else if waitFor == 0 && r == '&' {
					waitFor = ';'
				} else if waitFor == r {
					if attribute == 0 {
						waitFor = 0
						skip = true
					}
				}

				// Incorrect usage of > in attributes
				if waitFor == '>' && (r == '"' || r == '\'') {
					var prev rune = 0
					for j := (i - 1); j >= 0; j-- {
						if runes[j] != ' ' {
							prev = runes[j]
							break
						}
					}

					if attribute == 0 && prev == '=' {
						attribute = r
					} else if attribute == r && prev != '\\' {
						attribute = 0
					}
				}

				// Incorrect usage of raw ampersand or non-terminated entities within HTML (far from perfect)
				if waitFor == ';' {
					entityLength++
					if r == ' ' || r == '&' {
						count += entityLength - 1
						waitFor = 0
						entityLength = 0

						if r == '&' {
							waitFor = ';'
							entityLength++
						}
					}
				}

				if len(tag) > 0 {
					tag += string(r)

					if waitFor == 0 {
						tag = strings.Split(strings.Trim(tag, "<>"), " ")[0]
						if tag[0:1] == "/" && len(history) > 0 && tag[1:] == history[len(history) - 1] {
							history = history[:len(history) - 1]
						} else {
							if _, ok := singleTags[tag]; !ok {
								history = append(history, tag)
							}
						}
						tag = ""
					}
				}
				
				all++
				if waitFor == 0 && !skip {
					
					if _, ok := wordBreaks[r]; ok {
						if !wordbreak {
							count++
							wordbreak = true
						}
					} else {
						wordbreak = false
					}

					if count >= intLength {
						break
					}
				}

				skip = false
			}

			output := strings.TrimRight(string(runes[:all]), " \n\r\t")

			if historyLength := len(history); historyLength > 0 {
				for i := historyLength - 1; i >= 0; i-- {
					output += "</" + history[i] + ">"
				}
			}

			return reflect.ValueOf(output)
	}

	return recursiveHelper(value, reflect.ValueOf(truncatewords), length)
}

/*
 func type[T any](value T) string
Returns a string representation of the reflection Type
*/
func typeFn(value reflect.Value) string {
	return value.Type().String()
}

/*
 func ul(value any) string
Converts slices, arrays or maps into an HTML unordered list.
*/
func ul(value reflect.Value) string {
	return listHelper(value, "ul")
}

/*
 func upper[T any](value T) T
Converts string text to upper case.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func upper(value reflect.Value) reflect.Value {
	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.ToUpper(value.String()))
	}

	return recursiveHelper(value, reflect.ValueOf(upper))
}

/*
 func urlDecode[T any](url T) T
Converts URL character-entity equivalents back into their literal, URL-unsafe forms.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func urlDecode(url reflect.Value) reflect.Value {
	switch url.Kind() {
		case reflect.String:
			find	:= []string{ "%21", "%2A", "%27", "%28", "%29", "%3B", "%3A", "%40", "%26", "%3D", "%2B", "%24", "%2C", "%2F", "%3F", "%25", "%23", "%5B", "%5D" }
			replace	:= []string{ "!",   "*",   "'",   "(",   ")",   ";",   ":",   "@",   "&",   "=",   "+",   "$",   ",",   "/",   "?",   "%",   "#",   "[",   "]" }
			replacer, err := replaceHelper(find, replace)
			if err != nil {
				logError(err.Error())
				return reflect.Value{}
			}
			return reflect.ValueOf(replacer.Replace(url.String()))
	}
	
	return recursiveHelper(url, reflect.ValueOf(urlDecode))
}

/*
 func urlEncode[T any](url T) T
Converts URL-unsafe characters into character-entity equivalents to allow the string to be used as part of a URL.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func urlEncode(url reflect.Value) reflect.Value {
	switch url.Kind() {
		case reflect.String:
			find	:= []string{ "!",   "*",   "'",   "(",   ")",   ";",   ":",   "@",   "&",   "=",   "+",   "$",   ",",   "/",   "?",   "%",   "#",   "[",   "]" }
			replace	:= []string{ "%21", "%2A", "%27", "%28", "%29", "%3B", "%3A", "%40", "%26", "%3D", "%2B", "%24", "%2C", "%2F", "%3F", "%25", "%23", "%5B", "%5D" }
			replacer, err := replaceHelper(find, replace)
			if err != nil {
				logError(err.Error())
				return reflect.Value{}
			}
			return reflect.ValueOf(replacer.Replace(url.String()))
	}

	return recursiveHelper(url, reflect.ValueOf(urlEncode))
}

/*
 func wordcount(value string) int
Counts the number of words (excluding HTML, numbers and special characters) in a string.
*/
func wordcount(value reflect.Value) int {
	switch value.Kind() {
		case reflect.String:
			tmp := strip.StripTags(value.String())
			tmp = urlDecode(reflect.ValueOf(tmp)).String()
			strip := map[string]string{
				"!": " ", "*": " ", "'": " ", "(": " ", ")": " ", ";": " ", ":": " ", "@": " ", "&": " ", "=": " ", "+": " ", "$": " ", ",": " ", "/": " ", "?": " ", 
				"%": " ", "#": " ", "[": " ", "]": " ", "`": " ", "¬": " ", `"`: " ", "£": " ", "^": " ", "-": " ", "_": " ", "{": " ", "}": " ", ".": " ", "~": " ", 
				"\\": " ", "<": " ", ">": " ", "|": " ", "0": " ", "1": " ", "2": " ", "3": " ", "4": " ", "5": " ", "6": " ", "7": " ", "8": " ", "9": " ", 
			}
			replacer, err := replaceHelper(strip)
			if err != nil {
				logError(err.Error())
				return 0
			}
			tmp = replacer.Replace(tmp)
			words := strings.Fields(tmp)
			return len(words)
	}

	logWarning("warning: wordcount(string) being called on a none string variable")
	return 0
}

/*
 func wrap[T any](prefix string, suffix string, value T) T
Wraps all strings within `value` with a prefix and suffix
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func wrap(prefixValue reflect.Value, suffixValue reflect.Value, value reflect.Value) reflect.Value {
	sig		:= "wrap(prefix string, suffix string, value any)"
	value	= reflectHelperUnpackInterface(value)

	if prefixValue.Kind() != reflect.String {
		logError(sig + " can only prefix string values")
		return value
	}

	if suffixValue.Kind() != reflect.String {
		logError(sig + " can only suffix string values")
		return value
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(prefixValue.String() + value.String() + suffixValue.String())
	}

	return recursiveHelper(value, reflect.ValueOf(wrap), prefixValue, suffixValue)
}

/*
 func year(times nil|time.Time) int
Returns an integer year from a `time.Time` input, or the current year if no time is provided.
*/
func year(times ...time.Time) int {
	t := time.Now().In(dateLocalTimezone)
	if len(times) > 0 {
		t = times[0]
	}
	year, _, _ := t.Date()

	return year
}

/*
Returns "Yes" for true values, "No" for false values and "Maybe" for empty values (`maybe` defaults to "No" unless maybe is specifically defined)
Return string options may be customised.
If numeric arguments are used, it treats numeric zero as "No", positive numbers as "Yes" and negative numbers as "Maybe"
If string, slice, array or map arguments are used, it treats empty as "Maybe", and populated as "Yes"
 yesno(test any)                                      // Uses the default "Yes" / "No" returns
 yesno(yes string, test any)                          // Customises the string used for "Yes"
 yesno(yes string, no string, test any)               // Customises the strings used for "Yes" and "No"
 yesno(yes string, no string, maybe string, test any) // Customises the strings used for "Yes", "No" and "Maybe" (enables `maybe`)
*/
func yesno(values ...reflect.Value) string {
	sig		:= "yesno(values ...any)"
	test	:= reflect.Value{}

	yes		:= reflect.ValueOf("Yes")
	no		:= reflect.ValueOf("No")
	maybe	:= reflect.ValueOf("No")

	if len(values) < 1 {
		logError(sig + " requires at least one argument")
		return ""
	} else if len(values) == 1 {
		test = values[0]
	} else if len(values) == 2 {
		yes		= values[0]
		test	= values[1]
	} else if len(values) == 3 {
		yes		= values[0]
		no		= values[1]
		maybe	= values[1]
		test	= values[2]
	} else if len(values) == 4 {
		yes		= values[0]
		no		= values[1]
		maybe	= values[2]
		test	= values[3]
	}

	test	= reflectHelperUnpackInterface(test)
	yes		= reflectHelperUnpackInterface(yes)
	no		= reflectHelperUnpackInterface(no)
	maybe	= reflectHelperUnpackInterface(maybe)

	if !no.IsValid() || no.Kind() != reflect.String {
		logError(sig + " value for `No` must be a string")
		return "No"
	}

	if !yes.IsValid() || yes.Kind() != reflect.String {
		logError(sig + " value for `Yes` must be a string")
		return no.String()
	}

	if !maybe.IsValid() || maybe.Kind() != reflect.String {
		logError(sig + " value for `Maybe` must be a string")
		return no.String()
	}

	if !test.IsValid() {
		return no.String()
	}

	switch test.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			if test.Len() > 0 {
				return yes.String()
			}
			return no.String()
		case reflect.Bool:
			if test.Bool() {
				return yes.String()
			}
			return no.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64:
			if integer, err := reflectHelperConvertToInt(test); err == nil {
				if integer > 0 {
					return yes.String()
				} else if integer != 0 {
					return maybe.String()
				}
				return no.String()
			}
			return maybe.String()
		case reflect.Struct:
			if !reflectHelperIsEmptyStruct(test) {
				return yes.String()
			}
			return no.String()
		default:
			return no.String()
	}
}