package templateManager

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/grokify/html-strip-tags-go" // => strip
	"github.com/google/uuid"
)

/*
Returns a function map for use with the Go template standard library
*/
func getDefaultFunctions() map[string]any {
	return map[string]any{
		"add":				add,
		"bool":				toBool,
		"capfirst":			capfirst,
		"collection":		collection, 
		"concat":			concat,
		"contains":			contains,
		"cut":				cut,
		"date":				date,
		"datetime":			datetime,
		"default":			defaultVal,
		"divide":			divide,
		"divideceil":		divideCeil,
		"dividefloor":		divideFloor,
		"divisibleby":		divisibleBy,
		"dl":				dl,
		"endswith":			endswith,
		"equal":			equal,
		"first":			first,
		"firstof":			firstOf,
		"float":			toFloat,
		"formattime":		formattime,
		"gto":				greaterThan,
		"gte":				greaterThanEqual,
		"htmldecode":		htmlDecode,
		"htmlencode":		htmlEncode,
		"int":				toInt,
		"iterable":			iterable,
		"join":				join,
		"jsondecode":		jsonDecode,
		"jsonencode":		jsonEncode,
		"key":				keyFn,
		"keys":				keys,
		"kind":				kind,
		"last":				last,
		"length":			length,
		"list":				list,
		"lto":				lessThan,
		"lte":				lessThanEqual,
		"localtime":		localtime,
		"lower":			lower,
		"ltrim":			ltrim,
		"md5":				md5Fn,
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
		"query":			query, 
		"random":			random,
		"regexp":			regexpFindAll,
		"regexpreplace":	regexpReplaceAll,
		"replace":			replaceAll,
		"round":			round,
		"rtrim":			rtrim,
		"sha1":				sha1Fn,
		"sha256":			sha256Fn,
		"sha512":			sha512Fn,
		"split":			split,
		"startswith":		startswith,
		"string":			toString,
		"striptags":		stripTags,
		"substr":			substr,
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
		"uuid":				uuid.NewString,
		"ul":				ul,
		"upper":			upper,
		"urldecode":		urlDecode,
		"urlencode":		urlEncode,
		"values":			values,
		"wordcount":		wordcount,
		"wrap":				wrap,
		"year":				year,
		"yesno":			yesno,
	}
}

/*
Returns a function map for use with the Go template standard library that will replace many of their functions 
with more consistent, fault tolerant and chainable alternatives
*/
func getOverloadFunctions() map[string]any {
	return map[string]any{
		"eq":				equal,
		"gt":				greaterThan,
		"ge":				greaterThanEqual,
		"len":				length,
		"index":			keyFn,
		"lt":				lessThan,
		"le":				lessThanEqual, 
		"ne":				notequal,
		"html":				htmlEncode,
		"urlquery":			urlEncode,
	}
}

/*
 func add[T any](value T, to T) (T, error)
Adds a value to the existing item.
For numeric items this is a simple addition. For other types this is appended / merged as appropriate.
*/
func add(value reflect.Value, to reflect.Value) (reflect.Value, error) {
	sig := "add(value any, to any)"

	value	= reflectHelperUnpackInterface(value)
	to		= reflectHelperUnpackInterface(to)

	if !value.IsValid() {
		err := logError(sig + " `value` added cannot be an untyped nil value")
		return to, err
	}

	if !to.IsValid() {
		err := logError(sig + " value being added `to` cannot be an untyped nil value")
		return to, err
	}

	// It's a simple type, do it recursively
	if reflectHelperIsNumeric(value) || value.Kind() == reflect.String {
		switch to.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
			reflect.Uint64:
				addVal, _ := reflectHelperConvertToFloat64(value)
				toVal, _ := reflectHelperConvertToFloat64(to)

				return reflect.ValueOf(int64(roundFloat(addVal + toVal, 0))).Convert(to.Type()), nil
			case reflect.Float32, reflect.Float64:
				addVal, _ := reflectHelperConvertToFloat64(value)
				toVal, _ := reflectHelperConvertToFloat64(to)

				return reflect.ValueOf(addVal + toVal).Convert(to.Type()), nil
			case reflect.String:
				addVal, _ := reflectHelperConvertToString(value)

				return reflect.ValueOf(to.String() + addVal), nil
		}

		return recursiveHelper(to, reflect.ValueOf(add), value)
	}

	// It's a more complex type, no recursion and stricter checks
	if err := reflectHelperLooseTypeCompatibility(value, to); err != nil {
		err = logError(sig + " the `value` and `to` parameters must have the same approximate types; trying to add %s to %s", value.Type(), to.Type())
		return to, err
	}

	switch to.Kind() {
		case reflect.Slice:
			slice, _ := reflectHelperCreateEmptySlice(value)
			slice = reflect.AppendSlice(slice, to)
			slice = reflect.AppendSlice(slice, value)
			return slice, nil
		case reflect.Array:
			slice, _ := reflectHelperCreateEmptySlice(value)
			for i := 0; i < to.Len(); i++ {
				slice = reflect.Append(slice, to.Index(i))
			}
			for i := 0; i < value.Len(); i++ {
				slice = reflect.Append(slice, value.Index(i))
			}
			arr, _ := reflectHelperConvertSliceToArray(slice)
			return arr, nil
		case reflect.Map:
			tmp := reflect.MakeMap(to.Type())

			iter := to.MapRange()
			for iter.Next() {
				tmp.SetMapIndex(iter.Key(), iter.Value())
			}

			iter = value.MapRange()
			for iter.Next() {
				if val := tmp.MapIndex(iter.Key()); val.IsValid() {
					recurse, _ := add(iter.Value(), val)
					tmp.SetMapIndex(iter.Key(), recurse)
				} else {
					tmp.SetMapIndex(iter.Key(), iter.Value())
				}
			}
			return tmp, nil
	}

	return to, nil
}

/*
 func capfirst[T any](value T) (T, error)
Capitalises the first letter of strings. Does not alter any other letters.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func capfirst(value reflect.Value) (reflect.Value, error) {
	sig := "capfirst(value string)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logWarning(sig + " cannot accept an untyped nil `value`")
		return value, err
	}

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
			return reflect.ValueOf(string(runes)), nil
	}

	return recursiveHelper(value, reflect.ValueOf(capfirst))
}

/*
 func collection(pairs ...any) (map[string]any, error)
Allows several variables to be packaged together into a map for passing to templates. 
*/
func collection(pairs ...any) (map[string]any, error) {
	sig := "collection(pairs ...any)"

	length := len(pairs)
	if length == 0 || length % 2 != 0 {
		err := logError(sig + " can only accept pairs of arguments (string / any)")
		return map[string]any{}, err
	}

	collection := make(map[string]any, length / 2)

	for i := 0; i < length; i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			err := logError(sig + " first member of a pair must be a string")
			return map[string]any{}, err
		}
		collection[key] = pairs[i + 1]
	}

	return collection, nil
}

/*
 func concat(values ...any) (string, error)
Concatenates any number of string-able values together in the order that they were declared.
*/
func concat(values ...reflect.Value) (reflect.Value, error) {
	sig := "concat(values ...any)"

	if len(values) < 1 {
		err := logError(sig + " requires at least 1 parameter")
		return reflect.ValueOf(""), err
	}

	str := ""
	for _, value := range values {
		value = reflectHelperUnpackInterface(value)
		val, err := reflectHelperConvertAnythingToString(value)
		if err == nil {
			str += val
		} else {
			logWarning(sig + " attempting to append an invalid value: %v (%s) - it was ignored", value, value.Type())
		}
	}

	return reflect.ValueOf(str), nil
}

/*
 func contains(find any, within any) (bool, error)
Returns a boolean value to determine whether the `find` value is contained in the `within` value.
The `find` value can act on strings, slices, arrays and maps.
*/
func contains(find reflect.Value, within reflect.Value) (bool, error) {
	sig := "contains(find any, within string|slice|map)"

	find	= reflectHelperUnpackInterface(find)
	within	= reflectHelperUnpackInterface(within)

	if !find.IsValid() {
		err := logWarning(sig + " is trying to search for an untyped nil value")
		return false, err
	}

	if !within.IsValid() {
		err := logWarning(sig + " is trying to search within an untyped nil value")
		return false, err
	}

	switch within.Kind() {
		case reflect.String:
			val, err := reflectHelperConvertToString(find)
			if err == nil {
				return strings.Contains(within.String(), val), nil
			}
			err = logError(sig + " can't search within a string using a %s", find.Type())
			return false, err
		case reflect.Array, reflect.Slice:
			var err error = nil
			if reflectHelperGetSliceType(within) == find.Type().String() {
				for i := 0; i < within.Len(); i++ {
					if reflect.DeepEqual(within.Index(i).Interface(), find.Interface()) {
						return true, err
					}
				}
			} else {
				err = logError(sig + " can't search within a slice type %s using a %s", reflectHelperGetSliceType(within), find.Type())
			}
			return false, err
		case reflect.Map:
			var err error = nil
			if reflectHelperGetMapType(within) == find.Type().String() {
				iter := within.MapRange()
				for iter.Next() {
					if reflect.DeepEqual(iter.Value().Interface(), find.Interface()) {
						return true, err
					}
				}
			} else {
				err = logError(sig + " can't search within a map type %s using a %s", reflectHelperGetMapType(within), find.Type())
			}
			return false, err
		case reflect.Struct:
			for i := 0; i < within.NumField(); i++ {
				field, err := reflectHelperGetStructValue(within, reflect.ValueOf(i))
				if err == nil {
					if field.Type().String() == find.Type().String() {
						if reflect.DeepEqual(field.Interface(), find.Interface()) {
							return true, nil
						}
					}
				}
			}
			return false, nil	
	}

	err := logWarning(sig + " can't search within an item of type %s", within.Type())
	return false, err
}

/*
 func cut[T any](remove string, from T) (T, error)
Will `remove` a string value that is contained in the `from` value.
If `from` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func cut(remove reflect.Value, from reflect.Value) (reflect.Value, error) {
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
func date(params ...any) (string, error) {
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
func datetime(params ...any) (string, error) {
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
 func defaultVal(def any, test any) (any, error)
Will return the second `test` value if it is not empty, else return the `def` value
*/
func defaultVal(def reflect.Value, test reflect.Value) (reflect.Value, error) {
	sig := "default(def any, value any)"

	def		= reflectHelperUnpackInterface(def)
	test	= reflectHelperUnpackInterface(test)

	if !def.IsValid() {
		err := logError(sig + " cannot set an untyped nil value as the default")
		return reflect.Value{}, err
	}

	switch test.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			if test.Len() > 0 {
				return test, nil
			}
		case reflect.Bool:
			if test.Bool() {
				return test, nil
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64:
			if integer, err := reflectHelperConvertToInt(test); err == nil {
				if integer != 0 {
					return test, nil
				}
			}
		case reflect.Struct:
			if !reflectHelperIsEmptyStruct(test) {
				return test, nil
			}
	}

	return def, nil
}

/*
 func divide[T any](divisor int|float, value T) (T, error)
Divides the `value` by the `divisor` and rounds if a float to integer conversion is required.
If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain.
*/
func divide(divisor reflect.Value, value reflect.Value) (reflect.Value, error) {
	return divideHelper(reflect.ValueOf("round"), divisor, value)
}

/*
 func divideceil[T any](divisor int|float, value T) (T, error)
Divides the `value` by the `divisor` and rounds up if a float to integer conversion is required.
If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain.
*/
func divideCeil(divisor reflect.Value, value reflect.Value) (reflect.Value, error) {
	return divideHelper(reflect.ValueOf("ceil"), divisor, value)
}

/*
 func dividefloor[T any](divisor int|float, value T) (T, error)
Divides the `value` by the `divisor` and rounds down if a float to integer conversion is required.
If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain.
*/
func divideFloor(divisor reflect.Value, value reflect.Value) (reflect.Value, error) {
	return divideHelper(reflect.ValueOf("floor"), divisor, value)
}

/*
 func divisibleby[T any](divisor int, value T) (bool, error)
Determines if the `value` is divisible by the `divisor`
*/
func divisibleBy(divisor reflect.Value, value reflect.Value) (bool, error) {
	sig := "divisibleby(divisor int, value any)"

	value = reflectHelperUnpackInterface(value)

	if !divisor.IsValid() {
		err := logError(sig + " divisor cannot be an untyped nil value")
		return false, err
	}

	if !reflectHelperIsNumeric(divisor) {
		err := logError(sig + " divisor must be numeric, not %s", value.Type())
		return false, err
	}

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64:
			val, _ := reflectHelperConvertToFloat64(value)
			div, _ := reflectHelperConvertToFloat64(divisor)

			if div == 0.0 {
				err := logWarning(sig + " divisor must not be zero")
				return false, err
			}
			result := val / div

			return equalFloats(result, roundFloat(result, 0)), nil
	}

	err := logWarning(sig + " attempting division of non numeric type: %s", value.Type())
	return false, err
}

/*
 func dl(value any) (string, error)
Converts slices, arrays or maps into an HTML definition list.
For maps this will use the keys as the dt elements.
*/
func dl(value reflect.Value) (string, error) {
	return listHelper(value, "dl")
}

/*
 func endswith(find any, value any) (bool, error)
Determines if a string ends with a certain value.
*/
func endswith(find reflect.Value, value reflect.Value) (bool, error) {
	sig := "endswith(find any, value any)"

	find	= reflectHelperUnpackInterface(find)
	value	= reflectHelperUnpackInterface(value)

	if !find.IsValid() || find.Kind() != reflect.String {
		err := logError(sig + " can only be used to find strings")
		return false, err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return false, err
	}

	switch value.Kind() {
		case reflect.String:
			return strings.HasSuffix(value.String(), find.String()), nil		
	}

	err := logError(sig + " can't handle items of type %s", value.Type())
	return false, err
}

/*
 func equal(values ...any) (bool, error)
Determines whether any values are equal.
*/
func equal(values ...reflect.Value) (bool, error) {
	sig := "equal(values ...any)"

	if len(values) < 2 {
		err := logError(sig + " at least two values required, %d provided", len(values))
		return false, err
	}

	for i, value := range values {
		values[i] = reflectHelperUnpackInterface(value)
		value = values[i]

		if !value.IsValid() {
			err := logWarning(sig + " cannot compare untyped nil values")
			return false, err
		}

		if i > 0 {
			err := reflectHelperVeryLooseTypeCompatibility(value, values[i - 1])
			if err != nil {
				return false, nil
			}

			switch value.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
					val1, _ := reflectHelperConvertToFloat64(value)
					val2, _ := reflectHelperConvertToFloat64(values[i - 1])

					if !equalFloats(val1, val2) {
						return false, nil
					}
				case reflect.String:
					if value.String() != values[i - 1].String() {
						return false, nil
					}
				case reflect.Bool:
					if value.Bool() != values[i - 1].Bool() {
						return false, nil
					}
				case reflect.Array, reflect.Slice:
					if !reflect.DeepEqual(value.Slice(0, value.Len() - 1).Interface(), values[i - 1].Slice(0, values[i - 1].Len() - 1).Interface()) {
						return false, nil
					}
				case reflect.Map, reflect.Struct:
					if !reflect.DeepEqual(value.Interface(), values[i - 1].Interface()) {
						return false, nil
					}
				default:
					return false, nil
			}
		}
	}

	return true, nil
}

/*
 func first(value string|slice|array) (any, error)
Gets the first value from slices / arrays / maps / structs or the first word from strings.
*/
func first(value reflect.Value) (reflect.Value, error) {
	sig := "first(value string|slice)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " value cannot be an untyped nil value")
		return reflect.Value{}, err
	}

	switch value.Kind() {
		case reflect.String:
			if value.Len() > 0 {
				str := strings.Split(strings.TrimLeft(value.String(), " \n\r\t"), " ")[0]
				return reflect.ValueOf(str), nil
			}
		case reflect.Array, reflect.Slice:
			if value.Len() > 0 {
				return value.Index(0), nil
			}
		case reflect.Map:
			if value.Len() > 0 {
				keys, err := reflectHelperMapSort(value)
				if err == nil {
					return value.MapIndex(keys.Index(0)), nil
				} else {
					iter := value.MapRange()
					for iter.Next() {
						return iter.Value(), nil
					}
				}
			}
		case reflect.Struct:
			value, err := reflectHelperGetStructValue(value, reflect.ValueOf(0))
			if err != nil {
				err := logError(sig + " " + err.Error())
				return reflect.Value{}, err
			}
			return value, nil
	}

	err := logError(sig + fmt.Sprintf(" can't handle items of type %s", value.Type()))
	return reflect.Value{}, err
}

/*
 func firstOf(values ...any) (any, error)
Accepts any number of values and returns the first one of them that exists and is not empty.
*/
func firstOf(values ...reflect.Value) (reflect.Value, error) {
	sig := "firstof(values ...any)"

	if len(values) < 1 {
		err := logError(sig + " being called without any parameters")
		return reflect.Value{}, err
	}

	for _, value := range values {
		value = reflectHelperUnpackInterface(value)

		if !value.IsValid() {
			continue
		}

		switch value.Kind() {
			case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
				if value.Len() > 0 {
					return value, nil
				}
			case reflect.Struct:
				empty := reflect.New(value.Type()).Elem().Interface()
				if !reflect.DeepEqual(value.Interface(), empty) {
					return value, nil
				}
			case reflect.Bool:
				if value.Bool() {
					return value, nil
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
			reflect.Uint64, reflect.Float32, reflect.Float64:
				if integer, err := reflectHelperConvertToInt(value); err == nil {
					if integer != 0 {
						return value, nil
					}
				}
		}
	}

	return reflect.Value{}, nil
}

/*
 func formattime(format string, t time.Time) (string, error)

Formats a time.Time object for display.
*/
func formattime(format string, t time.Time) (string, error) {
	return t.Format(dateFormatHelper(format)), nil
}

/*
 func greaterThan(value1 any, value2 any) (bool, error)
Determines if `value2` is greater than `value1`
*/
func greaterThan(value1 reflect.Value, value2 reflect.Value) (bool, error) {
	sig := "gto(value any, value any)"
	
	value1 = reflectHelperUnpackInterface(value1)
	if !value1.IsValid() {
		err := logError(sig + " values cannot be untyped nil values")
		return false, err
	}

	value2 = reflectHelperUnpackInterface(value2)
	if !value2.IsValid() {
		err := logError(sig + " values cannot be untyped nil values")
		return false, err
	}

	err := reflectHelperVeryLooseTypeCompatibility(value1, value2)
	if err != nil {
		err := logError(sig + " values of dramatically different types cannot be compared")
		return false, err
	}

	switch value1.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			val1, _ := reflectHelperConvertToFloat64(value1)
			val2, _ := reflectHelperConvertToFloat64(value2)

			if val2 > val1 {
				return true, nil
			}
	}

	err = logError(sig + " values cannot be type %s", value1.Type())
	return false, err
}

/*
 func greaterThanEqual(value1 any, value2 any) (bool, error)
Determines if `value2` is greater than or equal to `value1`
*/
func greaterThanEqual(value1 reflect.Value, value2 reflect.Value) (bool, error) {
	sig := "gte(value any, value any)"
	
	value1 = reflectHelperUnpackInterface(value1)
	if !value1.IsValid() {
		err := logError(sig + " values cannot be untyped nils")
		return false, err
	}

	value2 = reflectHelperUnpackInterface(value2)
	if !value2.IsValid() {
		err := logError(sig + " values cannot be untyped nils")
		return false, err
	}

	err := reflectHelperVeryLooseTypeCompatibility(value1, value2)
	if err != nil {
		err := logError(sig + " values of dramatically different types cannot be compared")
		return false, err
	}

	switch value1.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			val1, _ := reflectHelperConvertToFloat64(value1)
			val2, _ := reflectHelperConvertToFloat64(value2)

			if val2 >= val1 || equalFloats(val1, val2) {
				return true, err
			}
	}

	err = logError(sig + " values cannot be type %s", value1.Type())
	return false, err
}

/*
 func htmlDecode[T any](value T) (T, error)
Converts HTML character-entity equivalents back into their literal, usable forms.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func htmlDecode(value reflect.Value) (reflect.Value, error) {
	sig := "htmldecode(value any)"
	
	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " values cannot be untyped nils")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			find	:= []string{ "&lt;", "&gt;", "&#34;", "&#x22;", "&quot;", "&#39;", "&#x27;", "&amp;" }
			replace	:= []string{ "<",    ">",    `"`,     `"`,      `"`,      "'",     "'",      "&" }
			replacer, err := replaceHelper(find, replace)
			if err != nil {
				err := logError(err.Error())
				return reflect.Value{}, err
			}
			return reflect.ValueOf(replacer.Replace(value.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(htmlDecode))
}

/*
 func htmlEncode[T any](value T) (T, error)
Converts literal HTML special characters into safe, character-entity equivalents.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func htmlEncode(value reflect.Value) (reflect.Value, error) {
	sig := "htmlencode(value any)"
	
	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " values cannot be untyped nils")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			find	:= []string{ "<",    ">",    `"`,     `"`,      `"`,      "'",     "'",      "&" }
			replace	:= []string{ "&lt;", "&gt;", "&#34;", "&#x22;", "&quot;", "&#39;", "&#x27;", "&amp;" }
			replacer, err := replaceHelper(find, replace)
			if err != nil {
				err := logError(err.Error())
				return reflect.Value{}, err
			}
			return reflect.ValueOf(replacer.Replace(value.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(htmlEncode))
}

/*
 func iterable(value ...int) ([]int, error)
Creates an integer slice so as to spoof a `for` loop:
 {{ range $v := iterable 5 }} -> for v := 0; v < 5; v++
 {{ range $v := iterable 3 5 }} -> for v := 3; v < 5; v++
 {{ range $v := iterable 3 5 2 }} -> for v := 3; v < 5; v += 2
*/
func iterable(values ...int) ([]int, error) {
	sig := "iterable(values ...int)"

	if len(values) < 1 {
		err := logError(sig + " requires at least one value")
		return []int{}, err
	}

	start := 0
	end := 1
	increment := 1

	switch len(values) {
		case 1: 
			end			= values[0]
		case 2:
			start		= values[0]
			end			= values[1]
		default:
			start		= values[0]
			end			= values[1]
			increment	= values[2]
	}

	if increment == 0 {
		err := logError(sig + " increment value must not be zero")
		return []int{}, err
	}

	if start > end && increment > 0 {
		err := logError(sig + " if start > end, increment value must be negative")
		return []int{}, err
	}

	if end > start && increment < 0 {
		err := logError(sig + " if end > start, increment value must be positive")
		return []int{}, err
	}

	items := []int{}
	for i := start; i < end; i += increment {
		items = append(items, i)
	}

	return items, nil
}

/*
 func join(separator string, values any) (string, error)
Joins slice or map `values` together as a string spaced by the `separator`.
*/
func join(separator reflect.Value, values reflect.Value) (string, error) {
	sig := "join(separator string, values any)"

	values		= reflectHelperUnpackInterface(values)
	separator	= reflectHelperUnpackInterface(separator)

	str := ""

	if !values.IsValid() {
		err := logError(sig + " is trying to join an untyped nil value")
		return str, err
	}

	if !separator.IsValid() || separator.Kind() != reflect.String {
		err := logError(sig + " can only join using strings")
		return str, err
	}

	switch values.Kind() {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool:
			return fmt.Sprintf("%v", values), nil
		case reflect.Array, reflect.Slice:
			for i := 0; i < values.Len(); i++ {
				if i > 0 { str += separator.String() }
				recurse, _ := join(separator, values.Index(i))
				str += recurse
			}
		case reflect.Map:
			keys, err := reflectHelperMapSort(values)
			if err == nil {
				for i := 0; i < keys.Len(); i++ {
					if i > 0 { str += separator.String() }
					recurse, _ := join(separator, values.MapIndex(keys.Index(i)))
					str += recurse
				}
			} else {
				iter := values.MapRange()
				i := 0
				for iter.Next() {
					if i > 0 { str += separator.String() }
					recurse, _ := join(separator, iter.Value())
					str += recurse
					i++
				}
			}
		case reflect.Struct:
			for i := 0; i < values.NumField(); i++ {
				if i > 0 { str += separator.String() }
				recurse, _ := join(separator, values.Field(i))
				str += recurse
			}
		default:
			err := logError(sig + " can't join items of type %s", values.Type())
			return str, err
	}

	return str, nil
}

/*
 func jsonDecode(value any) string
Decodes any JSON string value to a map.
*/
func jsonDecode(value string) (any, error) {
	sig := "jsondecode(value string)"

	var result any
	err := json.Unmarshal([]byte(value), &result)

	if err != nil {
		err = logError(sig + " " + err.Error())
	}

	return result, err
}

/*
 func jsonEncode(value any) (string, error)
Encodes any value to a JSON string.
*/
func jsonEncode(value any) (string, error) {
	sig := "jsonencode(value any)"

	result, err := json.Marshal(value)

	if err != nil {
		err = logError(sig + " " + err.Error())
		return "", err
	}

	return string(result), nil
}

/*
 func keyFn(input ...any) (any, error)
Accepts any number of nested keys and returns the result of indexing its final argument by them.
For strings this returns a byte value.
The indexed item must be a string, map, slice, or array.
*/
func keyFn(input ...reflect.Value) (reflect.Value, error) {
	sig := "key(indexes ...any, value any)"

	if len(input) < 2 {
		err := logError(sig + " requires at least two arguments")
		return reflect.Value{}, err
	}

	value	:= reflectHelperUnpackInterface(input[len(input) - 1])
	indexes	:= input[:len(input) - 1]
	var nilPointer bool
	var err error

	if !value.IsValid() {
		err := logError(sig + " is trying to access an untyped nil value")
		return reflect.Value{}, err
	}

	for _, index := range indexes {
		index = reflectHelperUnpackInterface(index)

		if value, nilPointer = reflectHelperCheckNilPointers(value); nilPointer {
			err := logError(sig + " is trying to access an index of a nil pointer")
			return reflect.Value{}, err
		}

		switch value.Kind() {
			case reflect.Array, reflect.Slice, reflect.String:
				origKind := value.Kind()
				value, err = reflectHelperGetSliceValue(value, index)
				if err != nil {
					err := logError(sig + " " + err.Error())
					return value, err
				}
				if origKind == reflect.String && value.Kind() == reflect.Uint8 {
					value = reflect.ValueOf(string(uint8(value.Uint())))
				}
			case reflect.Map:
				value, err = reflectHelperGetMapValue(value, index)
				if err != nil {
					err := logError(sig + " " + err.Error())
					return value, err
				}
			case reflect.Struct:
				value, err = reflectHelperGetStructValue(value, index)
				if err != nil {
					err := logError(sig + " " + err.Error())
					return value, err
				}
			default:
				err := logError(sig + " can't index item of type %s", value.Type())
				return reflect.Value{}, err
		}
	}

	return value, nil
}

/*
 func keys(value slice|map|struct) ([]any, error)
Returns the keys of a slice / array / map / struct
*/
func keys(value reflect.Value) (reflect.Value, error) {
	sig := "keys(value any)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " is trying to access an untyped nil value")
		return reflect.ValueOf([]int{}), err
	}

	switch value.Kind() {
		case reflect.Slice, reflect.Array:
			slice := []int{}
			for i := 0; i < value.Len(); i++ {
				slice = append(slice, i)
			}
			return reflect.ValueOf(slice), nil
		case reflect.Map:
			t := value.Type().Key()
			t = reflect.SliceOf(t)
			slice := reflect.New(t).Elem()
			keys, err := reflectHelperMapSort(value)
			if err == nil {
				for i := 0; i < keys.Len(); i++ {
					slice = reflect.Append(slice, keys.Index(i))
				}
			} else {
				iter := value.MapRange()
				for iter.Next() {
					slice = reflect.Append(slice, iter.Key())
				}
			}
			return slice, nil
		case reflect.Struct:
			slice := []any{}
			for i := 0; i < value.NumField(); i++ {
				slice = append(slice, value.Type().Field(i).Name)	
			}
			return reflect.ValueOf(slice), nil
	}

	err := logWarning(sig + " being called on a non-[slice|array|map|struct]")
	return reflect.ValueOf([]int{}), err
}

/*
 func kind[T any](value T) (string, error)
Returns a string representation of the reflection Kind
*/
func kind(value reflect.Value) (string, error) {
	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		return "invalid", nil
	}

	return value.Kind().String(), nil
}

/*
 func last(value string|slice|array) (any, error)
Gets the last value from slices / arrays or the last word from strings.
*/
func last(value reflect.Value) (reflect.Value, error) {
	sig := "last(value string|slice)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " value cannot be an untyped nil value")
		return reflect.Value{}, err
	}

	switch value.Kind() {
		case reflect.String:
			if value.Len() > 0 {
				sl := strings.Split(strings.TrimRight(value.String(), " \n\r\t"), " ")
				str := sl[len(sl) - 1]
				return reflect.ValueOf(str), nil
			}
		case reflect.Array, reflect.Slice:
			if value.Len() > 0 {
				return value.Index(value.Len() - 1), nil
			}
		case reflect.Map:
			if value.Len() > 0 {
				keys, err := reflectHelperMapSort(value)
				if err == nil {
					return value.MapIndex(keys.Index(keys.Len() - 1)), nil
				} else {
					iter := value.MapRange()
					v := reflect.Value{}
					for iter.Next() {
						v = iter.Value()
					}
					return v, nil
				}
			}
		case reflect.Struct:
			value, err := reflectHelperGetStructValue(value, reflect.ValueOf(value.NumField() - 1))
			if err != nil {
				err := logError(sig + " " + err.Error())
				return reflect.Value{}, err
			}
			return value, nil
	}

	err := logError(sig + " can't handle items of type %s", value.Type())
	return reflect.Value{}, err
}

/*
 func length(value any) (int, error)
Gets the length of any type without panics.
*/
func length(value reflect.Value) (int, error) {
	sig := "length(value any)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logWarning(sig + " value cannot be an untyped nil value")
		return 0, err
	}

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return len(fmt.Sprintf("%v", value)), nil
		case reflect.Float32, reflect.Float64:
			return len(fmt.Sprintf("%#v", value)), nil
		case reflect.Bool:
			return 1, nil
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			return value.Len(), nil
		case reflect.Struct:
			return value.NumField(), nil
	}

	err := logError(sig + " can't handle items of type %s", value.Type())
	return 0, err
}

/*
 func lessthan(value1 any, value2 any) (bool, error)
Determines if `value2` is less than `value1`
*/
func lessThan(value1 reflect.Value, value2 reflect.Value) (bool, error) {
	sig := "lto(value any, value any)"
	
	value1 = reflectHelperUnpackInterface(value1)
	if !value1.IsValid() {
		err := logError(sig + " values cannot be untyped nil values")
		return false, err
	}

	value2 = reflectHelperUnpackInterface(value2)
	if !value2.IsValid() {
		err := logError(sig + " values cannot be untyped nil values")
		return false, err
	}

	err := reflectHelperVeryLooseTypeCompatibility(value1, value2)
	if err != nil {
		err := logError(sig + " values of dramatically different types cannot be compared")
		return false, err
	}

	switch value1.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			val1, _ := reflectHelperConvertToFloat64(value1)
			val2, _ := reflectHelperConvertToFloat64(value2)

			if val2 < val1 {
				return true, nil
			}
	}

	err = logError(sig + " values cannot be type %s", value1.Type())
	return false, err
}

/*
 func lessthanequal(value1 any, value2 any) (bool, error)
Determines if `value2` is less than or equal to `value1`
*/
func lessThanEqual(value1 reflect.Value, value2 reflect.Value) (bool, error) {
	sig := "lte(value any, value any)"
	
	value1 = reflectHelperUnpackInterface(value1)
	if !value1.IsValid() {
		err := logError(sig + " values cannot be untyped nil values")
		return false, err
	}

	value2 = reflectHelperUnpackInterface(value2)
	if !value2.IsValid() {
		err := logError(sig + " values cannot be untyped nil values")
		return false, err
	}

	err := reflectHelperVeryLooseTypeCompatibility(value1, value2)
	if err != nil {
		err := logError(sig + " values of dramatically different types cannot be compared")
		return false, err
	}

	switch value1.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			val1, _ := reflectHelperConvertToFloat64(value1)
			val2, _ := reflectHelperConvertToFloat64(value2)

			if val2 < val1 || equalFloats(val1, val2) {
				return true, err
			}
	}

	err = logError(sig + " values cannot be type %s", value1.Type())
	return false, err
}

/*
 func list(values ...any) ([]any, error)
Creates a slice from any number of values
*/
func list(values ...reflect.Value) ([]reflect.Value, error) {
	return values, nil
}

/*
 func localtime(location string|time.Location, t time.Time) (time.Time, error)
Localises a time.Time object to display local times / dates.
*/
func localtime(location any, t time.Time) (time.Time, error) {
	var tz *time.Location
	switch v := location.(type) {
		case time.Location:
			tz = &v
		case string:
			tmp, err := time.LoadLocation(v)
			if err != nil {
				return t, err
			}
			tz = tmp
	}

	return t.In(tz), nil
}

/*
 func lower[T any](value T) (T, error)
Converts string text to lower case.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func lower(value reflect.Value) (reflect.Value, error) {
	sig := "lower(value string)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logWarning(sig + " cannot accept an untyped nil value")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.ToLower(value.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(lower))
}

/*
 func ltrim[T any](remove string, value T) (T, error)
Removes the passed characters from the left end of string values.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func ltrim(remove reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "ltrim(remove string, value any)"

	remove	= reflectHelperUnpackInterface(remove)
	value	= reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logWarning(sig + " cannot accept an untyped nil value")
		return value, err
	}

	if remove.Kind() != reflect.String {
		err := logError(sig + " remove can only be a string")
		return reflect.Value{}, err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.TrimLeft(value.String(), remove.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(ltrim), remove)
}

/*
 func md5(input any) (string, error)
Computes an md5 hash of the input.
*/
func md5Fn(value reflect.Value) (string, error) {
	sig := "md5(input any)"
	
	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " cannot accept untyped nil values")
		return "", err
	}

	switch value.Kind() {
		case reflect.String:
			hash := md5.Sum([]byte(value.String()))
			return hex.EncodeToString(hash[:]), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool:
				hash := md5.Sum([]byte(fmt.Sprintf("%v", value)))
				return hex.EncodeToString(hash[:]), nil
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
				hash := md5.Sum([]byte(fmt.Sprintf("%T%v", value.Interface(), value)))
				return hex.EncodeToString(hash[:]), nil
	}

	err := logError(sig + " values of type: %s cannot be hashed", value.Type())
	return "", err
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
func mktime(params ...string) (time.Time, error) {
	sig	:= "mktime(params ...string)"
	t	:= time.Now()

	if len(params) == 1 {
		tmp, err := time.Parse(time.RFC3339, params[0])
		if err != nil {
			err := logError(sig + " Invalid RFC3339 (\"" + time.RFC3339 + "\") passed: mktime(\"" + params[0] + "\")")
			return t.In(dateLocalTimezone), err
		}
		t = tmp
	} else if len(params) == 2 {
		tmp, err := time.Parse(dateFormatHelper(params[0]), params[1])
		if err != nil {
			err := logError(sig + "Invalid date / format passed: mktime(\"" + params[0] + "\", \"" + params[1] + "\")")
			return t.In(dateLocalTimezone), err
		}
		t = tmp
	}

	return t.In(dateLocalTimezone), nil
}

/*
 func multiply[T any](multiplier int|float, value T) (T, error)
Multiplies the `value` by the `multiplier`.
If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain.
*/
func multiply(multiplier reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "multiply(multiplier int, value any)"

	multiplier	= reflectHelperUnpackInterface(multiplier)
	value		= reflectHelperUnpackInterface(value)

	if !multiplier.IsValid() {
		err := logError(sig + " multiplier cannot be an untyped nil value")
		return value, err
	}

	if !reflectHelperIsNumeric(multiplier) {
		err := logError(sig + " multiplier must be numeric, not %s", value.Type())
		return value, err
	}

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64:
			val, _ := reflectHelperConvertToFloat64(value)
			mul, _ := reflectHelperConvertToFloat64(multiplier)
			op := val * mul
			return reflect.ValueOf(int64(roundFloat(op, 0))).Convert(value.Type()), nil
		case reflect.Float32, reflect.Float64:
			val, _ := reflectHelperConvertToFloat64(value)
			mul, _ := reflectHelperConvertToFloat64(multiplier)
			op := val * mul
			return reflect.ValueOf(op).Convert(value.Type()), nil
		case reflect.String, reflect.Bool:
			err := logWarning(sig + " trying to multiply a %s", value.Type())
			return value, err
	}

	return recursiveHelper(value, reflect.ValueOf(multiply), multiplier)
}

/*
 func nl2br[T any](value T) (T, error)
Replaces all instances of "\n" (new line) with instances of "<br>" within `value`.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func nl2br(value reflect.Value) (reflect.Value, error) {
	value, _ = replaceAll(reflect.ValueOf("\r\n"), reflect.ValueOf("\n"), value)
	value, _ = replaceAll(reflect.ValueOf("\r"), reflect.ValueOf("\n"), value)
	return replaceAll(reflect.ValueOf("\n"), reflect.ValueOf("<br>"), value)
}

/*
 func notequal(values ...any) (bool, error)
Determines whether any values are not equal.
*/
func notequal(values ...reflect.Value) (bool, error) {
	eq, err := equal(values...)
	return !eq, err
}

/*
 func now() (time.Time, error)
Returns the current `time.Time` value
*/
func now() (time.Time, error) {
	return time.Now().In(dateLocalTimezone), nil
}

/*
 func ol(value any) (string, error)
Converts slices, arrays or maps into an HTML ordered list.
*/
func ol(value reflect.Value) (string, error) {
	return listHelper(value, "ol")
}

/*
 func ordinal[T int|float64|string](value T) (string, error)
Suffixes a number with the correct English ordinal
*/
func ordinal(value reflect.Value) (string, error) {
	sig := "ordinal(value int)"

	value = reflectHelperUnpackInterface(value)

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
			return strconv.Itoa(integer) + suffix, nil
	}

	err := logError(sig + " attempting an ordinal conversion on a non numeric type")
	return "", err
}

/*
 func paragraph[T any](value T) (T, error)
Replaces all instances of "\n+" (multiple new lines) with paragraphs and instances of "\n" (new line) with instances of "<br>" within `value`
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func paragraph(value reflect.Value) (reflect.Value, error) {
	value, _ = replaceAll(reflect.ValueOf("\r\n"), reflect.ValueOf("\n"), value)
	value, _ = replaceAll(reflect.ValueOf("\r"), reflect.ValueOf("\n"), value)
	value, _ = regexpReplaceAll(reflect.ValueOf("(\\s*\\n\\s*){2,}"), reflect.ValueOf("</p><p>"), value)
	value, _ = replaceAll(reflect.ValueOf("\n"), reflect.ValueOf("<br>"), value)
	value, _ = wrap(reflect.ValueOf("<p>"), reflect.ValueOf("</p>"), value)

	return value, nil
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
func pluralise(values ...any) (string, error) {
	if len(values) < 1 {
		err := logError("pluralise(): called without argument")
		return "", err
	}

	num := 1
	suffixSingular := ""
	suffixPlural := "s"

	if len(values) == 1 {
		switch v := values[0].(type) {
			case int: num = v
			default: 
				err := logError("pluralise(num int): single value should be an integer")
				return "", err
		}
	} else if len(values) == 2 {
		switch v := values[0].(type) {
			case string: suffixPlural = v
			default: 
				err := logError("pluralise(suffix string, num int): first value should be a string")
				return "", err
		}

		switch v := values[1].(type) {
			case int: num = v
			default: 
				err := logError("pluralise(suffix string, num int): final value should be an integer")
				return "", err
		}
	} else if len(values) == 3 {
		switch v := values[0].(type) {
			case string: suffixSingular = v
			default: 
				err := logError("pluralise(suffixSingular string, suffixPlural string, num int): first value should be a string")
				return "", err
		}

		switch v := values[1].(type) {
			case string: suffixPlural = v
			default: 
				err := logError("pluralise(suffixSingular string, suffixPlural string, num int): second value should be a string")
				return "", err
		}

		switch v := values[2].(type) {
			case int: num = v
			default: 
				err := logError("pluralise(suffixSingular string, suffixPlural string, num int): final value should be an integer")
				return "", err
		}
	}

	if num == 1 {
		return suffixSingular, nil
	}

	return suffixPlural, nil
}

/*
 func prefix[T any](prefixes ...string, value T) (T, error)
Prefixes all strings within `value` with `prefixes` (in order)
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func prefix(values ...reflect.Value) (reflect.Value, error) {
	sig := "prefix(suffixes ...string, value any)"

	if len(values) < 1 {
		err := logError(sig + " requires at least 2 values")
		return reflect.ValueOf(""), err
	}

	if len(values) < 2 {
		err := logError(sig + " requires at least 2 values")
		return reflect.ValueOf(values[0]), err
	}

	prefixes	:= values[:len(values) - 1]
	value		:= values[len(values) - 1:][0]
	value		 = reflectHelperUnpackInterface(value)

	switch value.Kind() {
		case reflect.String:
			str := ""
			for _, prefix := range prefixes {
				prefix = reflectHelperUnpackInterface(prefix)
				pref, err := reflectHelperConvertToString(prefix)
				if err != nil {
					err := logError(sig + " can only prefix values that can be converted into strings")
					return value, err
				}
				str += pref
			}
			str += value.String()
			return reflect.ValueOf(str), nil
	}

	return recursiveHelper(value, reflect.ValueOf(prefix), prefixes...)
}

/*
 func query[T any](name string, value any, link T) (T, error)
Adds / replaces a query parameter with `name` and `value` in the provided `link`
If `link` is a slice, array or map it will apply this conversion to any string elements that it contains.
*/
func query(name reflect.Value, value reflect.Value, link reflect.Value) (reflect.Value, error) {
	sig := "query(name string, value string, link any)"

	name	= reflectHelperUnpackInterface(name)
	value	= reflectHelperUnpackInterface(value)
	link	= reflectHelperUnpackInterface(link)

	if !name.IsValid() {
		err := logError(sig + " variable name cannot be an untyped nil value")
		return link, err
	}

	if !link.IsValid() {
		err := logError(sig + " link cannot be an untyped nil value")
		return link, err
	}

	if name.Kind() != reflect.String {
		err := logError(sig + " can only append string parameter names")
		return link, err
	}

	if !value.IsValid() {
		value = reflect.ValueOf("")
	}

	switch link.Kind() {
		case reflect.String:
			uri, err := url.Parse(link.String())
			if err != nil {
				err := logError(sig + " invalid URL passed: `" + link.String() + "`")
				return reflect.ValueOf(""), err
			}

			query := uri.Query()
			replaceFind := "TEMPLATE_MANAGER_QUERY_FUNCTION_REPLACE=1"
			replaceValue := ""

			// Remove the variable that we are concerned with completely
			query.Del(name.String())
			for k := range query {
				if strings.HasPrefix(k, name.String() + "[") {
					query.Del(k)
				}
			}

			switch value.Kind() {
				case reflect.String:
					query.Set(name.String(), value.String())
				case reflect.Bool:
					val, _ := reflectHelperConvertToString(value)
					query.Set(name.String(), val)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					val, _ := reflectHelperConvertToString(value)
					query.Set(name.String(), val)
				case reflect.Float32, reflect.Float64:
					val, _ := reflectHelperConvertToString(value)
					query.Set(name.String(), val)
				case reflect.Slice, reflect.Array:
					for i := 0; i < value.Len(); i++ {
						val, err := reflectHelperConvertToString(value.Index(i))
						if err == nil {
							if i > 0 {
								replaceValue += "&"
							}
							encode, _ := urlEncode(reflect.ValueOf(val))
							replaceValue += name.String() + "[]=" + encode.String()
						}
					}
				case reflect.Map:
					iter := value.MapRange()
					i := 0
					for iter.Next() {
						val, err := reflectHelperConvertToString(iter.Value())
						if err == nil {
							if i > 0 {
								replaceValue += "&"
							}
							encode, _ := urlEncode(reflect.ValueOf(val))
							replaceValue += name.String() + "[" + iter.Key().String() + "]=" + encode.String()
						}
						i++
					}
				case reflect.Struct:
					for i := 0; i < value.NumField(); i++ {
						val, err := reflectHelperConvertToString(value.Field(i))
						if err == nil {
							if i > 0 {
								replaceValue += "&"
							}
							encode, _ := urlEncode(reflect.ValueOf(val))
							replaceValue += name.String() + "[" + value.Type().Field(i).Name + "]=" + encode.String()
						}
					}
			}

			if len(replaceValue) > 0 {
				query.Set("TEMPLATE_MANAGER_QUERY_FUNCTION_REPLACE", "1")
			}

			uri.RawQuery = query.Encode()

			// Remove unnecessary encoding applied by the package
			url := uri.String();
			for k := range query {
				if strings.Contains(k, "[") {
					encode, _ := urlEncode(reflect.ValueOf(k))
					url = strings.ReplaceAll(url, encode.String(), k)
				}
			}
			
			if len(replaceValue) > 0 {
				url = strings.ReplaceAll(url, replaceFind, replaceValue)
			}

			return reflect.ValueOf(url), nil
	}

	return recursiveHelper(link, reflect.ValueOf(query), name, value)
}

/*
Generates random numbers
 random()                 // Returns a random number between 0 and 10000
 random(limit int)        // Returns a random number between 0 and `limit`
 random(min int, max int) // Returns a random number between `min` and `max`
*/
func random(values ...int) (int, error) {
	rand.Seed(time.Now().UnixNano())

	if len(values) < 1 {
		return rand.Intn(10000), nil
	} else if len(values) == 1 {
		return rand.Intn(values[0]), nil
	}

	min := values[0]
	max := values[1]
	if min > max {
		min = values[1]
		max = values[0]
	} else if min == max {
		return min, nil
	}

	return rand.Intn(max - min) + min, nil
}

/*
 func regexpFindAll(find string, value string) ([][]string, error)
Finds all instances of `find` regexp within `value`.
It ONLY acts on strings
*/
func regexpFindAll(find reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "regexp(find string, value string)"

	find	= reflectHelperUnpackInterface(find)
	value	= reflectHelperUnpackInterface(value)

	if find.Kind() != reflect.String {
		err := logError(sig + " can only find string values")
		return reflect.ValueOf([][]string{}), err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept untyped nil values")
		return value, err
	}

	findRegexp, err := regexp.Compile(find.String())
	if err != nil {
		err := logError(sig + " invalid regexp: " + find.String())
		return reflect.ValueOf([][]string{}), err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(findRegexp.FindAllStringSubmatch(value.String(), -1)), nil
	}

	return reflect.ValueOf([][]string{}), nil
}

/*
 func regexpReplaceAll[T any](find string, replace string, value T) (T, error)
Replaces all instances of `find` regexp with instances of `replace` within `value`
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func regexpReplaceAll(find reflect.Value, replace reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "regexpreplace(find string, replace string, value any)"

	find	= reflectHelperUnpackInterface(find)
	replace	= reflectHelperUnpackInterface(replace)
	value	= reflectHelperUnpackInterface(value)

	if find.Kind() != reflect.String {
		err := logError(sig + " can only find string values")
		return value, err
	}

	if replace.Kind() != reflect.String {
		err := logError(sig + " can only replace string values")
		return value, err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept untyped nil values")
		return value, err
	}

	findRegexp, err := regexp.Compile(find.String())
	if err != nil {
		err := logError(sig + " invalid regexp: " + find.String())
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(findRegexp.ReplaceAllString(value.String(), replace.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(regexpReplaceAll), find, replace)
}

/*
 func replaceAll[T any](find string, replace string, value T) (T, error)
Replaces all instances of `find` with instances of `replace` within `value`
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func replaceAll(find reflect.Value, replace reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "replace(find string, replace string, value any)"

	find	= reflectHelperUnpackInterface(find)
	replace	= reflectHelperUnpackInterface(replace)
	value	= reflectHelperUnpackInterface(value)

	if find.Kind() != reflect.String {
		err := logError(sig + " can only find string values")
		return value, err
	}

	if replace.Kind() != reflect.String {
		err := logError(sig + " can only replace string values")
		return value, err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept untyped nil values")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.ReplaceAll(value.String(), find.String(), replace.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(replaceAll), find, replace)
}

/*
 func round[T any](precision int, value T) (T, error)
Rounds any floats to the required precision.
If `value` is a slice, array or map it will apply this conversion to any float elements that they contain.
*/
func round(precision reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "round(precision int, value any)"

	precision	= reflectHelperUnpackInterface(precision)
	value		= reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " cannot accept untyped nil values")
		return value, err
	}

	prec, err := reflectHelperConvertToUint(precision)
	if err != nil {
		err := logError(sig + " precision can only be an integer")
		return value, err
	}

	switch value.Kind() {
		case reflect.Float32, reflect.Float64:
			val, _ := reflectHelperConvertToFloat64(value)
			return reflect.ValueOf(roundFloat(val, prec)).Convert(value.Type()), nil
	}

	return recursiveHelper(value, reflect.ValueOf(round), precision)
}

/*
 func rtrim[T any](remove string, value T) (T, error)
Removes the passed characters from the right end of string values.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func rtrim(remove reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "rtrim(remove string, value any)"

	remove	= reflectHelperUnpackInterface(remove)
	value	= reflectHelperUnpackInterface(value)

	if remove.Kind() != reflect.String {
		err := logError(sig + " remove can only be a string")
		return reflect.Value{}, err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept untyped nil values")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.TrimRight(value.String(), remove.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(rtrim), remove)
}

/*
 func sha1(input any) (string, error)
Computes a SHA1 hash of the input.
*/
func sha1Fn(value reflect.Value) (string, error) {
	sig := "sha1(input any)"
	
	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " values cannot be untyped nil values")
		return "", err
	}

	switch value.Kind() {
		case reflect.String:
			hash := sha1.Sum([]byte(value.String()))
			return hex.EncodeToString(hash[:]), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool:
			hash := sha1.Sum([]byte(fmt.Sprintf("%v", value)))
			return hex.EncodeToString(hash[:]), nil
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
			hash := sha1.Sum([]byte(fmt.Sprintf("%T%v", value.Interface(), value)))
			return hex.EncodeToString(hash[:]), nil
	}

	err := logError(sig + " values of type: %s cannot be hashed", value.Type())
	return "", err
}

/*
 func sha256(input any) (string, error)
Computes a SHA256 hash of the input.
*/
func sha256Fn(value reflect.Value) (string, error) {
	sig := "sha256(input any)"
	
	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " values cannot be untyped nil values")
		return "", err
	}

	switch value.Kind() {
		case reflect.String:
			hash := sha256.Sum256([]byte(value.String()))
			return hex.EncodeToString(hash[:]), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool:
			hash := sha256.Sum256([]byte(fmt.Sprintf("%v", value)))
			return hex.EncodeToString(hash[:]), nil
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
				hash := sha256.Sum256([]byte(fmt.Sprintf("%T%v", value.Interface(), value)))
				return hex.EncodeToString(hash[:]), nil
	}

	err := logError(sig + " values of type: %s cannot be hashed", value.Type())
	return "", err
}

/*
 func sha512(input any) (string, error)
Computes a SHA512 hash of the input.
*/
func sha512Fn(value reflect.Value) (string, error) {
	sig := "sha256(input any)"
	
	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " values cannot be untyped nil values")
		return "", err
	}

	switch value.Kind() {
		case reflect.String:
			hash := sha512.Sum512([]byte(value.String()))
			return hex.EncodeToString(hash[:]), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool:
			hash := sha512.Sum512([]byte(fmt.Sprintf("%v", value)))
			return hex.EncodeToString(hash[:]), nil
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
				hash := sha512.Sum512([]byte(fmt.Sprintf("%T%v", value.Interface(), value)))
				return hex.EncodeToString(hash[:]), nil
	}

	err := logError(sig + " values of type: %s cannot be hashed", value.Type())
	return "", err
}

/*
 func split(separator string, value string) ([]string, error)
Splits strings on the `separator` value and returns a slice of the pieces.
*/
func split(separator reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "split(separator string, value any)"

	separator	= reflectHelperUnpackInterface(separator)
	value		= reflectHelperUnpackInterface(value)

	if separator.Kind() != reflect.String {
		err := logError(sig + " separator can only be a string")
		return reflect.Value{}, err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept untyped nil values")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			tmp := []string{}
			for _, val := range strings.Split(value.String(), separator.String()) {
				if val != "" {
					tmp = append(tmp, val)
				}
			}
			return reflect.ValueOf(tmp), nil
	}

	err := logError(sig + " can only split strings, type: %s given", value.Type())
	return reflect.Value{}, err
}

/*
 func startswith(find any, value any) (bool, error)
Determines if a string starts with a certain value.
*/
func startswith(find reflect.Value, value reflect.Value) (bool, error) {
	sig := "startswith(find any, value any)"

	find	= reflectHelperUnpackInterface(find)
	value	= reflectHelperUnpackInterface(value)

	if !find.IsValid() || find.Kind() != reflect.String {
		err := logError(sig + " can only be used to find strings")
		return false, err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return false, err
	}

	switch value.Kind() {
		case reflect.String:
			return strings.HasPrefix(value.String(), find.String()), nil	
	}

	err := logError(sig + " can't handle items of type %s", value.Type())
	return false, err
}

/*
 func stripTags[T any](value T) (T, error)
Strips HTML tags from strings.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func stripTags(value reflect.Value) (reflect.Value, error) {
	sig := "striptags(value any)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strip.StripTags(value.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(stripTags))
}

/*
 func substr[T any](offset int, length int, value T) (T, error)
Extracts a substring from a `value` starting at the specified `offset` and including `length` runes from that point.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func substr(offset reflect.Value, length reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "substr(offset int, length int, value any)"

	offset	= reflectHelperUnpackInterface(offset)
	length	= reflectHelperUnpackInterface(length)
	value	= reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " value acted upon must be a string or an object that can contain strings")
		return value, err
	}

	if !reflectHelperIsNumeric(offset) || !offset.IsValid() || offset.Int() < 0 {
		err := logError(sig + " offset must be a positive number")
		return value, err
	}

	if !reflectHelperIsNumeric(length) || !length.IsValid() {
		err := logError(sig + " length must be a number")
		return value, err
	}

	length, _ = reflectHelperConvertUnderlying(length, reflect.Int64)

	str, err := reflectHelperConvertToString(value)
	if err == nil {
		runes := []rune(str)

		if length.Int() == 0 {
			str = string(runes[offset.Int():])
		} else if length.Int() < 0 {
			end := length.Int() + int64(len(runes))
			if end < offset.Int() {
				end = offset.Int() 
			}
			str = string(runes[offset.Int():end])
		} else {
			end := length.Int() + offset.Int()
			if end > int64(len(runes)) {
				end = int64(len(runes))
			}
			str = string(runes[offset.Int():end])
		}

		ret, _ := reflectHelperConvertUnderlying(reflect.ValueOf(str), value.Kind())

		return ret, nil
	}

	return recursiveHelper(value, reflect.ValueOf(substr), offset, length)
}

/*
 func subtract[T any](value T, to T) (T, error)
Removes a value from the existing item.
For numeric items this is a simple subtraction. For other types this is removed as appropriate.
*/
func subtract(value reflect.Value, from reflect.Value) (reflect.Value, error) {
	sig := "subtract(value any, from any)"

	value	= reflectHelperUnpackInterface(value)
	from	= reflectHelperUnpackInterface(from)

	if !value.IsValid() {
		err := logError(sig + " value subtracted cannot be an untyped nil value")
		return from, err
	}

	if !from.IsValid() {
		err := logError(sig + " value subtracted from cannot be an untyped nil value")
		return from, err
	}

	// It's a simple type, do it recursively
	if reflectHelperIsNumeric(value) || value.Kind() == reflect.String {
		switch from.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
			reflect.Uint64:
				subVal, _ := reflectHelperConvertToFloat64(value)
				fromVal, _ := reflectHelperConvertToFloat64(from)

				return reflect.ValueOf(int64(roundFloat(fromVal - subVal, 0))).Convert(from.Type()), nil
			case reflect.Float32, reflect.Float64:
				subVal, _ := reflectHelperConvertToFloat64(value)
				fromVal, _ := reflectHelperConvertToFloat64(from)

				return reflect.ValueOf(fromVal - subVal).Convert(from.Type()), nil
			case reflect.String:
				subVal, _ := reflectHelperConvertToString(value)

				return cut(reflect.ValueOf(subVal), from)
		}

		return recursiveHelper(from, reflect.ValueOf(subtract), value)
	}

	if err := reflectHelperLooseTypeCompatibility(value, from); err != nil {
		err := logError(sig + " the value and subtraction must have the same types; trying to remove %s from %s", value.Type(), from.Type())
		return from, err
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
			return slice, nil
		case reflect.Map:
			tmp := reflect.MakeMap(from.Type())
			iter := from.MapRange()
			for iter.Next() {
				if val := value.MapIndex(iter.Key()); val.IsValid() {
					recurse, _ := subtract(val, iter.Value())
					subtracted := recurse
					if !reflectHelperIsEmpty(subtracted) {
						tmp.SetMapIndex(iter.Key(), subtracted)
					}
				} else {
					tmp.SetMapIndex(iter.Key(), iter.Value())
				}
			}
			return tmp, nil
	}

	return from, nil
}

/*
 func suffix[T any](suffixes ...string, value T) (T, error)
Suffixes all strings within `value` with `suffixes` (in order)
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func suffix(values ...reflect.Value) (reflect.Value, error) {
	sig := "suffix(suffixes ...string, value any)"

	if len(values) < 1 {
		err := logError(sig + " requires at least 2 values")
		return reflect.ValueOf(""), err
	}

	if len(values) < 2 {
		err := logError(sig + " requires at least 2 values")
		return reflect.ValueOf(values[0]), err
	}

	suffixes	:= values[:len(values) - 1]
	value		:= values[len(values) - 1:][0]
	value		 = reflectHelperUnpackInterface(value)

	switch value.Kind() {
		case reflect.String:
			str := value.String()
			for _, suffix := range suffixes {
				suffix = reflectHelperUnpackInterface(suffix)
				suff, err := reflectHelperConvertToString(suffix)
				if err != nil {
					err := logError(sig + " can only suffix values which can be converted to strings")
					return value, err
				}
				str += suff
			}
			return reflect.ValueOf(str), nil
	}

	return recursiveHelper(value, reflect.ValueOf(suffix), suffixes...)
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
func timeFn(params ...any) (string, error) {
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
					err := logError(sig + " Invalid RFC3339 date passed: time(\"%s\", \"%s\")", f, val)
					return "", err
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
			err := logError(sig + " Invalid date / format passed to time in template: time(\"%s\", \"%s\", \"%s\")\n%s", f, params[1].(string), params[2].(string), err.Error())
			return "", err
		}
		if strings.Contains(l, "MST") {
			location, err := time.LoadLocation(tmp.Location().String())
			if err != nil {
				err := logError(sig + " " + err.Error())
				return "", err
			}
			tmp, _ = time.ParseInLocation(l, params[2].(string), location)
		}
		t = tmp
	}

	return t.In(dateLocalTimezone).Format(f), nil
}

/*
 func timeSince(t time.Time) (map[string]int, error)
Calculates the approximate duration since the `time.Time` value.
The map of integers contains the keys: `years`, `weeks`, `days`, `hours`, `minutes`, `seconds`
*/
func timeSince(t time.Time) (map[string]int, error) {
	return formatDuration(time.Since(t))
}

/*
 func timeUntil(t time.Time) (map[string]int, error)
Calculates the approximate duration until the `time.Time` value.
The map of integers contains the keys: `years`, `weeks`, `days`, `hours`, `minutes`, `seconds`
*/
func timeUntil(t time.Time) (map[string]int, error) {
	return formatDuration(time.Until(t))
}

/*
 func title[T any](value T) (T, error)
Converts string text to title case.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func title(value reflect.Value) (reflect.Value, error) {
	sig := "title(value any)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.Title(strings.ToLower(value.String()))), nil
	}

	return recursiveHelper(value, reflect.ValueOf(title))
}

/*
 func toBool(value any) (bool, error)
Attempts to convert any `value` to a boolean. If this is impossible, the nil value (false) will be returned.
*/
func toBool(value reflect.Value) (bool, error) {
	sig := "bool(value any)"

	value = reflectHelperUnpackInterface(value)

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
			val, err := reflectHelperConvertToBool(value)
			if err != nil {
				err := logError(sig + " could not convert %T with value: `%v` to a bool", value.Interface(), value)
				return false, err
			}
			return val, nil
		case reflect.Invalid:
			return false, nil
	}

	err := logError(sig + " can only convert simple types to booleans, not a %T", value.Interface())
	return false, err
}

/*
 func toFloat(value any) (float64, error)
Attempts to convert any `value` to a float64. If this is impossible, the nil value (0.0) will be returned.
*/
func toFloat(value reflect.Value) (float64, error) {
	sig := "float(value any)"

	value = reflectHelperUnpackInterface(value)

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
			val, err := reflectHelperConvertToFloat64(value)
			if err != nil {
				err := logError(sig + " could not convert %T with value: `%v` to a float", value.Interface(), value)
				return 0.0, err
			}
			return val, nil
		case reflect.Invalid:
			return 0.0, nil
	}

	err := logError(sig + " can only convert simple types to floats, not a %T", value.Interface())
	return 0.0, err
}

/*
 func toInt(value any) (int, error)
Attempts to convert any `value` to an integer. If this is impossible, the nil value (0) will be returned.
*/
func toInt(value reflect.Value) (int, error) {
	sig := "int(value any)"

	value = reflectHelperUnpackInterface(value)

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
			val, err := reflectHelperConvertToInt(value)
			if err != nil {
				err := logError(sig + " could not convert %T with value: `%v` to an int", value.Interface(), value)
				return 0, err
			}
			return val, nil
		case reflect.Invalid:
			return 0, nil
	}

	err := logError(sig + " can only convert simple types to integers, not a %T", value.Interface())
	return 0, err
}

/*
 func toString(value any) (string, error)
Attempts to convert any `value` to a string. If this is impossible, the nil value ("") will be returned.
*/
func toString(value reflect.Value) (string, error) {
	sig := "string(value any)"

	value = reflectHelperUnpackInterface(value)

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
			val, err := reflectHelperConvertToString(value)
			if err != nil {
				err := logError(sig + " could not convert %T with value: `%v` to a string", value.Interface(), value)
				return "", err
			}
			return val, nil
		case reflect.Invalid:
			return "", nil
	}

	err := logError(sig + " can only convert simple types to strings, not a %T", value.Interface())
	return "", err
}

/*
 func trim[T any](remove string, value T) (T, error)
Removes the passed characters from the ends of string values.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func trim(remove reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "trim(remove string, value any)"
	
	remove	= reflectHelperUnpackInterface(remove)
	value	= reflectHelperUnpackInterface(value)
	
	if !remove.IsValid() || remove.Kind() != reflect.String {
		err := logError(sig + " `remove` can only be a string")
		return value, err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.Trim(value.String(), remove.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(trim), remove)
}

/*
 func truncate[T any](length int, value T) (T, error)
Truncates strings to a certain number of characters. It is multi-byte safe.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func truncate(length reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "truncate(length int, value any)"
	
	length	= reflectHelperUnpackInterface(length)
	value	= reflectHelperUnpackInterface(value)
	
	if !length.IsValid() || !reflectHelperIsNumeric(length) {
		err := logError(sig + " length can only be a number")
		return value, err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return value, err
	}

	intLength, _ := reflectHelperConvertToInt(length)

	switch value.Kind() {
		case reflect.String:
			if intLength <= 0 {
				return reflect.ValueOf(""), nil
			}

			runes := []rune(value.String())
			
			if stringLength := len(runes); intLength >= stringLength {
				return value, nil
			}

			if !strings.Contains(value.String(), "<") && !strings.Contains(value.String(), "&") {
				return reflect.ValueOf(string(runes[:intLength])), nil
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

			return reflect.ValueOf(output), nil
	}

	return recursiveHelper(value, reflect.ValueOf(truncate), length)
}

/*
 func truncatewords[T any](length int, value T) (T, error)
Truncates strings to a certain number of words. It is multi-byte safe.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func truncatewords(length reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "truncatewords(length int, value any)"
	
	length	= reflectHelperUnpackInterface(length)
	value	= reflectHelperUnpackInterface(value)
	
	if !length.IsValid() || !reflectHelperIsNumeric(length) {
		err := logError(sig + " length can only be a number")
		return value, err
	}

	if !value.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return value, err
	}

	intLength, _ := reflectHelperConvertToInt(length)

	switch value.Kind() {
		case reflect.String:
			if intLength <= 0 {
				return reflect.ValueOf(""), nil
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

			return reflect.ValueOf(output), nil
	}

	return recursiveHelper(value, reflect.ValueOf(truncatewords), length)
}

/*
 func type[T any](value T) (string, error)
Returns a string representation of the reflection Type
*/
func typeFn(value reflect.Value) (string, error) {
	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		return "invalid", nil
	}

	return value.Type().String(), nil
}

/*
 func ul(value any) (string, error)
Converts slices, arrays or maps into an HTML unordered list.
*/
func ul(value reflect.Value) (string, error) {
	return listHelper(value, "ul")
}

/*
 func upper[T any](value T) (T, error)
Converts string text to upper case.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func upper(value reflect.Value) (reflect.Value, error) {
	sig := "upper(value any)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return value, err
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.ToUpper(value.String())), nil
	}

	return recursiveHelper(value, reflect.ValueOf(upper))
}

/*
 func urlDecode[T any](url T) (T, error)
Converts URL character-entity equivalents back into their literal, URL-unsafe forms.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func urlDecode(url reflect.Value) (reflect.Value, error) {
	sig := "urldecode(value any)"

	url = reflectHelperUnpackInterface(url)

	if !url.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return url, err
	}

	switch url.Kind() {
		case reflect.String:
			find	:= []string{ "%21", "%2A", "%27", "%28", "%29", "%3B", "%3A", "%40", "%26", "%3D", "%2B", "%24", "%2C", "%2F", "%3F", "%25", "%23", "%5B", "%5D" }
			replace	:= []string{ "!",   "*",   "'",   "(",   ")",   ";",   ":",   "@",   "&",   "=",   "+",   "$",   ",",   "/",   "?",   "%",   "#",   "[",   "]" }
			replacer, err := replaceHelper(find, replace)
			if err != nil {
				err := logError(err.Error())
				return reflect.Value{}, err
			}
			return reflect.ValueOf(replacer.Replace(url.String())), nil
	}
	
	return recursiveHelper(url, reflect.ValueOf(urlDecode))
}

/*
 func urlEncode[T any](url T) T
Converts URL-unsafe characters into character-entity equivalents to allow the string to be used as part of a URL.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func urlEncode(url reflect.Value) (reflect.Value, error) {
	sig := "urlencode(value any)"

	url = reflectHelperUnpackInterface(url)

	if !url.IsValid() {
		err := logError(sig + " cannot accept an untyped nil value")
		return url, err
	}

	switch url.Kind() {
		case reflect.String:
			find	:= []string{ "!",   "*",   "'",   "(",   ")",   ";",   ":",   "@",   "&",   "=",   "+",   "$",   ",",   "/",   "?",   "%",   "#",   "[",   "]" }
			replace	:= []string{ "%21", "%2A", "%27", "%28", "%29", "%3B", "%3A", "%40", "%26", "%3D", "%2B", "%24", "%2C", "%2F", "%3F", "%25", "%23", "%5B", "%5D" }
			replacer, err := replaceHelper(find, replace)
			if err != nil {
				err := logError(err.Error())
				return reflect.Value{}, err
			}
			return reflect.ValueOf(replacer.Replace(url.String())), nil
	}

	return recursiveHelper(url, reflect.ValueOf(urlEncode))
}

/*
 func values(value slice|map|struct) ([]any, error)
Returns the values of a slice / array / map / struct
*/
func values(value reflect.Value) (reflect.Value, error) {
	sig := "values(value any)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logError(sig + " is trying to access an untyped nil value")
		return reflect.ValueOf([]int{}), err
	}

	switch value.Kind() {
		case reflect.Slice, reflect.Array:
			t := value.Type().Elem()
			t = reflect.SliceOf(t)
			slice := reflect.New(t).Elem()
			for i := 0; i < value.Len(); i++ {
				slice = reflect.Append(slice, value.Index(i))
			}
			return slice, nil
		case reflect.Map:
			t := value.Type().Elem()
			t = reflect.SliceOf(t)
			slice := reflect.New(t).Elem()
			keys, err := reflectHelperMapSort(value)
			if err == nil {
				for i := 0; i < keys.Len(); i++ {
					slice = reflect.Append(slice, value.MapIndex(keys.Index(i)))
				}
			} else {
				iter := value.MapRange()
				for iter.Next() {
					slice = reflect.Append(slice, iter.Value())
				}
			}
			return slice, nil
		case reflect.Struct:
			slice := []any{}
			for i := 0; i < value.NumField(); i++ {
				switch value.Field(i).Kind() {
					case reflect.Bool:
						bool, _ := reflectHelperConvertToBool(value.Field(i))
						slice = append(slice, bool)
					case reflect.Int:
						num, _ := reflectHelperConvertToInt(value.Field(i))
						slice = append(slice, num)
					case reflect.Int8:
						num, _ := reflectHelperConvertToInt8(value.Field(i))
						slice = append(slice, num)
					case reflect.Int16:
						num, _ := reflectHelperConvertToInt16(value.Field(i))
						slice = append(slice, num)
					case reflect.Int32:
						num, _ := reflectHelperConvertToInt32(value.Field(i))
						slice = append(slice, num)
					case reflect.Int64:
						num, _ := reflectHelperConvertToInt64(value.Field(i))
						slice = append(slice, num)
					case reflect.Uint:
						num, _ := reflectHelperConvertToUint(value.Field(i))
						slice = append(slice, num)
					case reflect.Uint8:
						num, _ := reflectHelperConvertToUint8(value.Field(i))
						slice = append(slice, num)
					case reflect.Uint16:
						num, _ := reflectHelperConvertToUint16(value.Field(i))
						slice = append(slice, num)
					case reflect.Uint32:
						num, _ := reflectHelperConvertToUint32(value.Field(i))
						slice = append(slice, num)
					case reflect.Uint64:
						num, _ := reflectHelperConvertToUint64(value.Field(i))
						slice = append(slice, num)
					case reflect.Float32:
						float, _ := reflectHelperConvertToFloat32(value.Field(i))
						slice = append(slice, float)
					case reflect.Float64:
						float, _ := reflectHelperConvertToFloat64(value.Field(i))
						slice = append(slice, float)
					case reflect.String:
						str, _ := reflectHelperConvertToString(value.Field(i))
						slice = append(slice, str)
					default:
						slice = append(slice, value.Field(i))
				}	
			}
			return reflect.ValueOf(slice), nil
	}

	err := logWarning(sig + " being called on a non-[slice|array|map|struct]")
	return reflect.ValueOf([]int{}), err
}

/*
 func wordcount(value string) (int, error)
Counts the number of words (excluding HTML, numbers and special characters) in a string.
*/
func wordcount(value reflect.Value) (int, error) {
	sig := "wordcount(value string)"

	value = reflectHelperUnpackInterface(value)

	if !value.IsValid() {
		err := logWarning(sig + " cannot accept an untyped nil value")
		return 0, err
	}

	switch value.Kind() {
		case reflect.String:
			tmp := strip.StripTags(value.String())
			decoded, _ := urlDecode(reflect.ValueOf(tmp))
			tmp = decoded.String()
			strip := map[string]string{
				"!": " ", "*": " ", "'": " ", "(": " ", ")": " ", ";": " ", ":": " ", "@": " ", "&": " ", "=": " ", "+": " ", "$": " ", ",": " ", "/": " ", "?": " ", 
				"%": " ", "#": " ", "[": " ", "]": " ", "`": " ", "¬": " ", `"`: " ", "£": " ", "^": " ", "-": " ", "_": " ", "{": " ", "}": " ", ".": " ", "~": " ", 
				"\\": " ", "<": " ", ">": " ", "|": " ", "0": " ", "1": " ", "2": " ", "3": " ", "4": " ", "5": " ", "6": " ", "7": " ", "8": " ", "9": " ", 
			}
			replacer, err := replaceHelper(strip)
			if err != nil {
				err := logError(err.Error())
				return 0, err
			}
			tmp = replacer.Replace(tmp)
			words := strings.Fields(tmp)
			return len(words), nil
	}

	err := logWarning(sig + " being called on a none string variable")
	return 0, err
}

/*
 func wrap[T any](prefix string, suffix string, value T) (T, error)
Wraps all strings within `value` with a prefix and suffix
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.
*/
func wrap(prefix reflect.Value, suffix reflect.Value, value reflect.Value) (reflect.Value, error) {
	sig := "wrap(prefix string, suffix string, value any)"

	value = reflectHelperUnpackInterface(value)

	pref, err := reflectHelperConvertToString(prefix)
	if err != nil {
		err := logError(sig + " can only prefix values that can be converted to a string")
		return value, err
	}

	suff, err := reflectHelperConvertToString(suffix)
	if err != nil {
		err := logError(sig + " can only suffix values that can be converted to a string")
		return value, err
	}

	if !value.IsValid() {
		logWarning(sig + " using an untyped nil `value` - assigning an empty string")
		value = reflect.ValueOf("")
	}

	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(pref + value.String() + suff), nil
	}

	return recursiveHelper(value, reflect.ValueOf(wrap), prefix, suffix)
}

/*
 func year(times nil|time.Time) (int, error)
Returns an integer year from a `time.Time` input, or the current year if no time is provided.
*/
func year(times ...time.Time) (int, error) {
	t := time.Now().In(dateLocalTimezone)
	if len(times) > 0 {
		t = times[0]
	}
	year, _, _ := t.Date()

	return year, nil
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
func yesno(values ...reflect.Value) (string, error) {
	sig		:= "yesno(values ...any)"

	test	:= reflect.Value{}
	yes		:= reflect.ValueOf("Yes")
	no		:= reflect.ValueOf("No")
	maybe	:= reflect.ValueOf("No")

	if len(values) < 1 {
		err := logError(sig + " requires at least one argument")
		return no.String(), err
	} else if len(values) == 1 {
		test	= reflectHelperUnpackInterface(values[0])
	} else if len(values) == 2 {
		yes		= reflectHelperUnpackInterface(values[0])
		test	= reflectHelperUnpackInterface(values[1])
	} else if len(values) == 3 {
		yes		= reflectHelperUnpackInterface(values[0])
		no		= reflectHelperUnpackInterface(values[1])
		maybe	= reflectHelperUnpackInterface(values[1])
		test	= reflectHelperUnpackInterface(values[2])
	} else if len(values) == 4 {
		yes		= reflectHelperUnpackInterface(values[0])
		no		= reflectHelperUnpackInterface(values[1])
		maybe	= reflectHelperUnpackInterface(values[2])
		test	= reflectHelperUnpackInterface(values[3])
	}

	if !no.IsValid() || no.Kind() != reflect.String {
		err := logError(sig + " value for `No` must be a string")
		return "No", err
	}

	if !yes.IsValid() || yes.Kind() != reflect.String {
		err := logError(sig + " value for `Yes` must be a string")
		return no.String(), err
	}

	if !maybe.IsValid() || maybe.Kind() != reflect.String {
		err := logError(sig + " value for `Maybe` must be a string")
		return no.String(), err
	}

	if !test.IsValid() {
		err := logError(sig + " must not pass an untyped nil value")
		return no.String(), err
	}

	switch test.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			if test.Len() > 0 {
				return yes.String(), nil
			}
			return no.String(), nil
		case reflect.Bool:
			if test.Bool() {
				return yes.String(), nil
			}
			return no.String(), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64:
			if integer, err := reflectHelperConvertToInt(test); err == nil {
				if integer > 0 {
					return yes.String(), nil
				} else if integer != 0 {
					return maybe.String(), nil
				}
				return no.String(), nil
			}
			return maybe.String(), nil
		case reflect.Struct:
			if !reflectHelperIsEmptyStruct(test) {
				return yes.String(), nil
			}
			return no.String(), nil
		default:
			return no.String(), nil
	}
}