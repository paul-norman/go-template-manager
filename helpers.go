package templateManager

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
Rounds floats to integers for numeric conversion.
*/
func roundFloat(float float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))

	return math.Round(float * ratio) / ratio
}

/*
Checks if two floats are equal (to a very small tolerance).
*/
func equalFloats(float1 float64, float2 float64) bool {
	diff := math.Abs(float1 - float2)

	return diff < 0.00000000001
}

/*
A helper that converts Python and PHP date formats to Go format
*/
func dateFormatHelper(date string) string {
	if !strings.Contains(date, "06") {

		// Convenience formats

		// Must be separate otherwise replacer will remove the shorter versions leaving Z characters
		replacements := map[string]string{
			"ISO8601Z": "X-m-d\\TH:i:sP",   // "2006-01-02T15:04:05Z07:00"
			"RFC1123Z": "D, d M Y H:i:s O", // "Mon, 02 Jan 2006 15:04:05 -07:00"
			"RFC822Z": "D, d M y H:i:s O",  // "Mon, 02 Jan 06 15:04:05 -07:00" should be?: "02 Jan 06 15:04 -07:00" ??
		}
		replacer, _ := replaceHelper(replacements)
		date = replacer.Replace(date)

		replacements = map[string]string{
			"ISO8601": 	"Y-m-d\\TH:i:sO",	// "2006-01-02T15:04:05-07:00"
			"RFC822": 	"D, d M y H:i:s T", // "Mon, 02 Jan 06 15:04:05 MST" should be?: "02 Jan 06 15:04 MST" ??
			"RFC850": 	"l, d-M-y H:i:s T", // "Monday, 02-Jan-06 15:04:05 MST"
			"RFC1036": 	"D, d M y H:i:s O",	// "02 Jan 06 15:04 -07:00"
			"RFC1123": 	"D, d M Y H:i:s T", // "Mon, 02 Jan 2006 15:04:05 MST"
			"RFC2822": 	"D, d M Y H:i:s O",	// "Mon, 02 Jan 2006 15:04:05 -07:00"
			"RFC3339": 	"Y-m-d\\TH:i:sP", 	// "2006-01-02T15:04:05Z07:00"
			"W3C": 		"Y-m-d\\TH:i:sP",	// "2006-01-02T15:04:05Z07:00"
			"ATOM": 	"Y-m-d\\TH:i:sP",	// "2006-01-02T15:04:05Z07:00"
			"COOKIE": 	"l, d-M-Y H:i:s T",	// "Monday, 02-Jan-2006 15:04:05 MST"
			"RSS": 		"D, d M Y H:i:s O",	// "Mon, 02 01 2006 15:04:05 +00:00"
			"MYSQL": 	"Y-m-d H:i:s",		// "2006-01-02 15:04:05"
			"UNIX": 	"D M _j H:i:s T Y", // "Mon Jan _2 15:04:05 MST 2006"
			"RUBY": 	"D M d H:i:s o Y", 	// "Mon Jan 02 15:04:05 -0700 2006"
			"ANSIC": 	"D M _j H:i:s Y", 	// "Mon Jan _2 15:04:05 2006"
		}
		replacer, _ = replaceHelper(replacements)
		date = replacer.Replace(date)

		if strings.Contains(date, "%") {
			// Python syntax support
			replacements = map[string]string{
				"%Y": "2006",
				"%y": "06",
				"%d": "02",
				"%I": "3",
				"%H": "15",
				"%M": "04",
				"%S": "05",
				"%m": "01",
				"%p": "PM",
				"%b": "Jan",
				"%B": "January",
				"%a": "Mon",
				"%A": "Monday",
				"%Z": "MST",
				"%z": "-07:00",
				"%f": "000",	// what is GO's microseconds?
			}
		} else {
			// PHP syntax support
			date = strings.ReplaceAll(date, "\\T", "@1")

			replacements = map[string]string{
				"Y": "2006",
				"y": "06",
				"d": "02",
				"j": "2",
				"g": "3",
				"H": "15",
				"i": "04",
				"s": "05",
				"n": "1",
				"m": "01",
				"a": "pm",
				"A": "PM",
				"M": "Jan",
				"F": "January",
				"D": "Mon",
				"l": "Monday",
				"T": "MST",
				"t": "-07:00",
				"P": "Z07:00",
				"O": "-07:00",
				"o": "-0700",
				"v": "000",		// what is GO's microseconds?
				"X": "2006",	// What even is this?
			}
		}

		replacer, _ = replaceHelper(replacements)
		date = replacer.Replace(date)
		date = strings.ReplaceAll(date, "@1", "T")
	}

	return date
}

/*
A helper that parses a `time.Duration` field into a map of integers containing the keys:

`years`, `weeks`, `days`, `hours`, `minutes`, `seconds`
*/
func formatDuration(duration time.Duration) map[string]int {
	const (
		Day		= 24 * time.Hour
		Week	= 7 * Day
		Year	= 8766 * time.Hour
	)

	years := (duration / Year)
	duration = duration % Year

	weeks := (duration / Week)
	duration = duration % Week

	days := (duration / Day)
	duration = duration % Day

	hours := duration / time.Hour
	duration = duration % time.Hour

	minutes := duration / time.Minute
	duration = duration % time.Minute

	seconds := duration / time.Second

	return map[string]int{"years": int(years), "weeks": int(weeks), "days": int(days), "hours": int(hours), "minutes": int(minutes), "seconds": int(seconds) }
}

/*
A helper that will look at all values in an `input` argument and run function `call` on each passing in
the `arguments` and the `input` value as the final argument 
*/
func recursiveHelper(input reflect.Value, call reflect.Value, arguments ...reflect.Value) reflect.Value {
	input = reflectHelperUnpackInterface(input)

	if !input.IsValid() {
		return reflect.Value{}
	}

	for i, v := range arguments {
		arguments[i] = reflect.ValueOf(reflectHelperUnpackInterface(v))
	}

	switch input.Kind() {
		case reflect.String:
			arguments = append(arguments, reflect.ValueOf(input))
			return call.Call(arguments)[0]
		case reflect.Slice, reflect.Array:
			tmp, _ := reflectHelperCreateEmptySlice(input)
			for i := 0; i < input.Len(); i++ {
				arguments = append(arguments, reflect.ValueOf(input.Index(i)))
				tmp = reflect.Append(tmp, call.Call(arguments)[0].Interface().(reflect.Value))
				arguments = arguments[:len(arguments) - 1]
			}
			if input.Kind() == reflect.Array {
				tmp, _ = reflectHelperConvertSliceToArray(tmp)
			}
			return tmp
		case reflect.Map:
			tmp := reflect.MakeMap(input.Type())
			iter := input.MapRange()
			for iter.Next() {
				arguments = append(arguments, reflect.ValueOf(iter.Value()))
				tmp.SetMapIndex(iter.Key(), call.Call(arguments)[0].Interface().(reflect.Value))
				arguments = arguments[:len(arguments) - 1]
			}
			return tmp
		case reflect.Struct:
			tmp := reflect.New(input.Type()).Elem()
			for i := 0; i < tmp.NumField(); i++ {
				if tmp.Field(i).CanSet() {
					arguments = append(arguments, reflect.ValueOf(input.Field(i)))
					tmp.Field(i).Set(call.Call(arguments)[0].Interface().(reflect.Value))
					arguments = arguments[:len(arguments) - 1]
				}
			}
			return tmp
		default:
			return input
	}
}

/*
Simple helper to perform the logic for ul, ol and dl functions
*/
func listHelper(value reflect.Value, tag string) string {
	sig		:= tag + "(value any, tag string)"
	value	= reflectHelperUnpackInterface(value)
	list	:= ""
	li		:= "li"
	if tag == "dl" {
		li = "dd"
	}

	if !value.IsValid() {
		logError(sig + " is trying to list an untyped nil value")
		return ""
	}

	switch value.Kind() {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, 
		reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool:
			return fmt.Sprintf("%v", value)
		case reflect.Array, reflect.Slice:
			list += "<" + tag + ">"
			for i := 0; i < value.Len(); i++ {
				list += "<" + li + ">" + listHelper(value.Index(i), tag) + "</" + li + ">"
			}
			list += "</" + tag + ">"
		case reflect.Map:
			keys, err := reflectHelperMapSort(value)
			list += "<" + tag + ">"
			if err == nil {
				for i := 0; i < keys.Len(); i++ {
					if tag == "dl" {
						list += "<dt>" + listHelper(keys.Index(i), tag) + "</dt>"
					}
					list += "<" + li + ">" + listHelper(value.MapIndex(keys.Index(i)), tag) + "</" + li + ">"
				}
			} else {
				iter := value.MapRange()
				for iter.Next() {
					if tag == "dl" {
						list += "<dt>" + listHelper(iter.Key(), tag) + "</dt>"
					}
					list += "<" + li + ">" + listHelper(iter.Value(), tag) + "</" + li + ">"
				}	
			}
			list += "</" + tag + ">"
		case reflect.Invalid:
			logError(sig + " invalid value passed")
			return ""
		default:
			logError(sig + fmt.Sprintf(" can't list items of type %s", value.Type()))
			return ""
	}

	return list
}

/*
Initialises a `strings.Replacer` from either: 

- A map of key / value (find / replace) pairs

- Two slices, the first containing the strings to find, the second with what to replace them with
*/
func replaceHelper(init ...any) (*strings.Replacer, error) {
	if len(init) < 1 {
		return strings.NewReplacer(), fmt.Errorf("replaceHelper(): you must pass at least one argument to initialise a string replacer")
	}
	
	// Initialise with a map of keys = find, value = replace pairs
	if len(init) == 1 {
		m := init[0].(map[string]string)
		replacerSlice := []string{}
		for find, replace := range m {
			replacerSlice = append(replacerSlice, find, replace)
		}

		return strings.NewReplacer(replacerSlice...), nil
	}

	// Initialise with two separate slices
	find	:= init[0].([]string)
	replace := init[1].([]string)

	if len(find) == len(replace) {
		replacerSlice := make([]string, 2 * len(find))
		for index, value := range find {
			replacerSlice[index * 2] = value
			replacerSlice[index * 2 + 1] = replace[index]
		}

		return strings.NewReplacer(replacerSlice...), nil
	}

	return strings.NewReplacer(), fmt.Errorf("replaceHelper(): must pass an even number of find and replace variables to a string replacer when using two slices")
}

/*
Creates the absolute value of an integer
*/
func interfaceHelperAbs(value any) (int64, error) {
	switch v := value.(type) {
		case int:
			val := int64(v)
			if v < 0 { return -val, nil }
			return val, nil
		case int8:
			val := int64(v)
			if v < 0 { return -val, nil }
			return val, nil
		case int16:
			val := int64(v)
			if v < 0 { return -val, nil }
			return val, nil
		case int32:
			val := int64(v)
			if v < 0 { return -val, nil }
			return val, nil
		case int64:
			if v < 0 { return -v, nil }
			return v, nil
		case uint:
			return int64(v), nil
		case uint8:
			return int64(v), nil
		case uint16:
			return int64(v), nil
		case uint32:
			return int64(v), nil
		case uint64:
			return int64(v), nil
		case uintptr:
			return int64(v), nil
		case float32:
			return int64(math.Abs(roundFloat(float64(v), 0))), nil
		case float64:
			return int64(math.Abs(roundFloat(v, 0))), nil
		case bool:
			if v {
				return int64(1), nil
			} 
			return int64(0), nil
		case string:
			tmp, err := strconv.Atoi(v)
			if err != nil {
				return int64(0), fmt.Errorf("can't convert type string to an int")
			}
			val := int64(tmp)
			if val < 0 { return -val, nil }
			return val, nil
	}

	return int64(0), fmt.Errorf("can't convert type %T to an int", value)
}

/*
Converts an `interface{}` to an `int64` if possible.
*/
func interfaceHelperConvertToInt64(value any) (int64, error) {
	var intValue int64

	switch v := value.(type) {
		case int:
			intValue = int64(v)
		case int8:
			intValue = int64(v)
		case int16:
			intValue = int64(v)
		case int32:
			intValue = int64(v)
		case int64:
			intValue = int64(v)
		case uint:
			intValue = int64(v)
		case uint8:
			intValue = int64(v)
		case uint16:
			intValue = int64(v)
		case uint32:
			intValue = int64(v)
		case uint64:
			intValue = int64(v)
		case uintptr:
			intValue = int64(v)
		case float32:
			intValue = int64(roundFloat(float64(v), 0))
		case float64:
			intValue = int64(roundFloat(v, 0))
		case bool:
			if v {
				intValue = int64(1)
			} else {
				intValue = int64(0)
			}
		case string:
			tmp, err := strconv.Atoi(v)
			if err != nil {
				return int64(0), fmt.Errorf("can't convert type string to an int")
			}
			intValue = int64(tmp)
		default:
			return int64(0), fmt.Errorf("can't convert type %T to an int", value)
	}

	return intValue, nil
}

/*
Converts an `interface{}` to an `int` if possible.
*/
func interfaceHelperConvertToInt(value any) (int, error) {
	int64Value, err := interfaceHelperConvertToInt64(value)

	return int(int64Value), err
}

/*
Converts an `interface{}` to an `int64` if possible.
*/
func interfaceHelperConvertToUint64(value any) (uint64, error) {
	var intValue uint64

	switch v := value.(type) {
		case int:
			tmp, _ := interfaceHelperAbs(v)
			intValue = uint64(tmp)
		case int8:
			tmp, _ := interfaceHelperAbs(v)
			intValue = uint64(tmp)
		case int16:
			tmp, _ := interfaceHelperAbs(v)
			intValue = uint64(tmp)
		case int32:
			tmp, _ := interfaceHelperAbs(v)
			intValue = uint64(tmp)
		case int64:
			tmp, _ := interfaceHelperAbs(v)
			intValue = uint64(tmp)
		case uint:
			intValue = uint64(v)
		case uint8:
			intValue = uint64(v)
		case uint16:
			intValue = uint64(v)
		case uint32:
			intValue = uint64(v)
		case uint64:
			intValue = uint64(v)
		case uintptr:
			intValue = uint64(v)
		case float32:
			intValue = uint64(math.Abs(roundFloat(float64(v), 0)))
		case float64:
			intValue = uint64(math.Abs(roundFloat(v, 0)))
		case bool:
			tmp, _ := interfaceHelperAbs(v)
			intValue = uint64(tmp)
		case string:
			tmp, _ := interfaceHelperAbs(v)
			intValue = uint64(tmp)
		default:
			return uint64(0), fmt.Errorf("can't convert type %T to a uint", value)
	}

	return intValue, nil
}

/*
If the `reflect.Value` is an `interface{}` unpack it to its concrete value.

If it is `nil`, returns a new `reflect.Value`
*/
func reflectHelperUnpackInterface(value reflect.Value) reflect.Value {
	if value.Kind() != reflect.Interface {
		return value
	}

	if value.IsNil() {
		return reflect.Value{}
	}

	return value.Elem()
}

/*
Checks if a `reflect.Value` is a pointer and checks it for `nil` value
*/
func reflectHelperCheckNilPointers(value reflect.Value) (reflect.Value, bool) {
	value = reflectHelperUnpackInterface(value)

	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return value, true
		}
	}

	return value, false
}

/*
Checks to see if `value` can be used as an argument of type `typ`.
Converts an invalid value to a zero value of the appropriate type if possible.
*/
func reflectHelperPrepareValue(value reflect.Value, typ reflect.Type) (reflect.Value, error) {
	if !value.IsValid() {
		if !reflectHelperCanBeNil(typ) {
			return reflect.Value{}, fmt.Errorf("value is nil; should be of type %s", typ)
		}
		value = reflect.Zero(typ)
	}

	if value.Type().AssignableTo(typ) {
		return value, nil
	}

	if reflectHelperIsInteger(value) && reflectHelperIsIntegerType(typ) && value.Type().ConvertibleTo(typ) {
		value = value.Convert(typ)

		return value, nil
	}

	if reflectHelperIsFloat(value) && reflectHelperIsFloatType(typ) && value.Type().ConvertibleTo(typ) {
		value = value.Convert(typ)

		return value, nil
	}

	return reflect.Value{}, fmt.Errorf("value has type %s; should be %s", value.Type(), typ)
}

func reflectHelperIsEmpty(value reflect.Value) bool {
    return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

/*
Determines if a struct is empty
*/
func reflectHelperIsEmptyStruct(value reflect.Value) bool {
	empty := reflect.New(value.Type()).Elem().Interface()
	return reflect.DeepEqual(value.Interface(), empty)
}

/*
Returns the value of a struct element at the specified index.
*/
func reflectHelperGetStructValue(structValue reflect.Value, index reflect.Value) (reflect.Value, error) {
	var value reflect.Value
	empty := reflect.New(structValue.Type()).Elem()

	switch index.Kind() {
		case reflect.String:
			value = structValue.FieldByName(index.String())
			test := empty.FieldByName(index.String())
			if test.CanSet() {
				return value, nil
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field, _ := reflectHelperConvertToInt(index)
			if structValue.NumField() > field {
				test := empty.Field(field)
				value = structValue.Field(field)
				if test.CanSet() {
					return value, nil
				}
			}
		default:
			return reflect.Value{}, fmt.Errorf("unsupported index type")
	}

	return reflectHelperDeepCopy(value)
}

func reflectHelperDeepCopy(value reflect.Value) (reflect.Value, error) {
	switch value.Kind() {
		case reflect.String:
			return reflect.ValueOf(value.String()), nil
		case reflect.Int:
			integer, _ := reflectHelperConvertToInt(value)
			return reflect.ValueOf(integer), nil
		case reflect.Int8:
			integer, _ := reflectHelperConvertToInt64(value)
			return reflect.ValueOf(int8(integer)), nil
		case reflect.Int16:
			integer, _ := reflectHelperConvertToInt64(value)
			return reflect.ValueOf(int16(integer)), nil
		case reflect.Int32:
			integer, _ := reflectHelperConvertToInt64(value)
			return reflect.ValueOf(int32(integer)), nil
		case reflect.Int64:
			integer, _ := reflectHelperConvertToInt64(value)
			return reflect.ValueOf(integer), nil
		case reflect.Uint:
			integer, _ := reflectHelperConvertToUint(value)
			return reflect.ValueOf(integer), nil
		case reflect.Uint8:
			integer, _ := reflectHelperConvertToUint64(value)
			return reflect.ValueOf(uint8(integer)), nil
		case reflect.Uint16:
			integer, _ := reflectHelperConvertToUint64(value)
			return reflect.ValueOf(uint16(integer)), nil
		case reflect.Uint32:
			integer, _ := reflectHelperConvertToUint64(value)
			return reflect.ValueOf(uint32(integer)), nil
		case reflect.Uint64:
			integer, _ := reflectHelperConvertToUint64(value)
			return reflect.ValueOf(integer), nil
		case reflect.Float32:
			float, _ := reflectHelperConvertToFloat(value)
			return reflect.ValueOf(float32(float)), nil
		case reflect.Float64:
			float, _ := reflectHelperConvertToFloat(value)
			return reflect.ValueOf(float), nil
		case reflect.Bool:
			return reflect.ValueOf(value.Bool()), nil
		case reflect.Slice, reflect.Array:
			tmp, _ := reflectHelperSliceCopy(value)
			return tmp, nil
		case reflect.Map:
			tmp, _ := reflectHelperMapCopy(value)
			return tmp, nil
		case reflect.Struct:
			tmp, _ := reflectHelperStructCopy(value)
			return tmp, nil
	}

	return reflect.Value{}, fmt.Errorf("can't copy variable %v", value)
}

func reflectHelperStructCopy(value reflect.Value) (reflect.Value, error) {
	tmp := reflect.New(value.Type()).Elem()
	for i := 0; i < tmp.NumField(); i++ {
		if tmp.Field(i).CanSet() {
			tmp.Field(i).Set(value.Field(i))
		}
	}
	return tmp, nil
}

func reflectHelperSliceCopy(value reflect.Value) (reflect.Value, error) {
	tmp, _ := reflectHelperCreateEmptySlice(value)
	for i := 0; i < value.Len(); i++ {
		val, _ := reflectHelperDeepCopy(value.Index(i))
		tmp = reflect.Append(tmp, val)
	}
	return tmp, nil
}

func reflectHelperMapCopy(value reflect.Value) (reflect.Value, error) {
	tmp := reflect.MakeMap(value.Type())
	iter := value.MapRange()
	for iter.Next() {
		key, _ := reflectHelperDeepCopy(iter.Key())
		val, _ := reflectHelperDeepCopy(iter.Value())
		tmp.SetMapIndex(key, val)
	}
	return tmp, nil
}

/*
Returns the value of a map element at the specified index.
*/
func reflectHelperGetMapValue(mapValue reflect.Value, index reflect.Value) (reflect.Value, error) {
	if mapValue.Len() == 0 {
		return reflect.Value{}, fmt.Errorf("can't access a map index on a zero length map")
	}
	
	index, err := reflectHelperPrepareValue(index, mapValue.Type().Key())
	if err != nil {
		return reflect.Value{}, err
	}

	if val := mapValue.MapIndex(index); val.IsValid() {
		return val, nil
	}
	
	//return reflect.Zero(mapValue.Type().Elem()), nil
	return reflect.Value{}, nil
}

/*
Returns a string representation a map element's value type.
*/
func reflectHelperGetMapType(m reflect.Value) string {
	typ := reflectHelperGetTypeString(m.Type())[4:]

	open := 1
	tmp := ""
	for i := range typ {
		tmp += string(typ[i])

		if string(typ[i]) == "]" {
			open -= 1
		} else if string(typ[i]) == "[" {
			open += 1
		}
		
		if open == 0 { break }
	}
	typ = typ[len(tmp):]

	return typ
}

/*
Sorts a map's key in default order for that type.
*/
func reflectHelperMapSort(value reflect.Value) (reflect.Value, error) {
	switch value.Kind() {
		case reflect.Map:
			keys := value.MapKeys()
			switch value.Type().Key().Kind() {
				case reflect.String:
					tmp := []string{}
					for _, key := range keys {
						tmp = append(tmp, key.String())
					}
					sort.Strings(tmp)
					return reflect.ValueOf(tmp), nil
				case reflect.Int:
					tmp := []int{}
					for _, key := range keys {
						tmp = append(tmp, int(key.Int()))
					}
					sort.Ints(tmp)
					return reflect.ValueOf(tmp), nil
				case reflect.Float64:
					tmp := []float64{}
					for _, key := range keys {
						tmp = append(tmp, key.Float())
					}
					sort.Float64s(tmp)
					return reflect.ValueOf(tmp), nil
			}
			return reflect.Value{}, fmt.Errorf("map key type unrecognised")
	}
	
	return reflect.Value{}, fmt.Errorf("value was not a map")
}

/*
Checks if the `reflect.Kind` is numeric
*/
func reflectHelperIsNumericKind(kind reflect.Kind) bool {
	if reflectHelperIsIntegerKind(kind) || reflectHelperIsFloatKind(kind) {
		return true
	}

	return false
}

/*
Checks if the `reflect.Type` is numeric
*/
func reflectHelperIsNumericType(typ reflect.Type) bool {
	return reflectHelperIsNumericKind(typ.Kind())
}

/*
Checks if the `reflect.Value` is numeric
*/
func reflectHelperIsNumeric(value reflect.Value) bool {
	return reflectHelperIsNumericKind(value.Kind())
}

/*
Checks if the `reflect.Kind` is an integer
*/
func reflectHelperIsIntegerKind(kind reflect.Kind) bool {
	switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return true
	}

	return false
}

/*
Checks if the `reflect.Type` is an integer
*/
func reflectHelperIsIntegerType(typ reflect.Type) bool {
	return reflectHelperIsIntegerKind(typ.Kind())
}

/*
Checks if the `reflect.Value` is an integer
*/
func reflectHelperIsInteger(value reflect.Value) bool {
	return reflectHelperIsIntegerKind(value.Kind())
}

/*
Checks if the `reflect.Kind` is a float
*/
func reflectHelperIsFloatKind(kind reflect.Kind) bool {
	switch kind {
		case reflect.Float32, reflect.Float64:
			return true
	}

	return false
}

/*
Checks if the `reflect.Type` is a float
*/
func reflectHelperIsFloatType(typ reflect.Type) bool {
	return reflectHelperIsFloatKind(typ.Kind())
}

/*
Checks if the `reflect.Value` is a float
*/
func reflectHelperIsFloat(value reflect.Value) bool {
	return reflectHelperIsFloatKind(value.Kind())
}

/*
Checks if the `reflect.Type` is allowed to be `nil`
*/
func reflectHelperCanBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			return true
		case reflect.Struct:
			return typ == reflect.TypeOf((*reflect.Value)(nil)).Elem()
	}
	return false
}

/*
Returns the value of an array / slice element at the specified index.
*/
func reflectHelperGetSliceValue(slice reflect.Value, index reflect.Value) (reflect.Value, error) {
	if slice.Len() == 0 {
		return reflect.Value{}, fmt.Errorf("can't access a slice index on a zero length slice")
	}

	key, err := reflectHelperCleanSliceIndex(index, slice.Len())
	if err != nil {
		return reflect.Value{}, err
	}

	if val := slice.Index(key); val.IsValid() {
		return val, nil
	}

	switch slice.Kind() {
		case reflect.String:
			return reflect.Zero(slice.Type()), nil
	}

	return reflect.Value{}, nil
}

/*
Returns a string representation of the `reflect.Type`
*/
func reflectHelperGetTypeString(typ reflect.Type) string {
	return fmt.Sprint(typ)
}

/*
Returns a string representation of the type that a slice / array contains
*/
func reflectHelperGetSliceType(slice reflect.Value) string {
	typ := reflectHelperGetTypeString(slice.Type())

	tmp := ""
	for i := range typ {
		tmp += string(typ[i])
		if string(typ[i]) == "]" { break }
	}
	typ = typ[len(tmp):]

	if typ[:1] == "[" {
		tmp = ""
		for i := range typ {
			tmp += string(typ[i])
			if string(typ[i]) == "]" { break }
		}
		if len(tmp) == 2 {
			return "slice"
		}

		return "array"
	}

	return typ
}

/*
Checks if a `reflect.Value` can be used as an index, and converts it to an `int` if possible.
*/
func reflectHelperCleanSliceIndex(index reflect.Value, length int) (int, error) {
	key, err := reflectHelperConvertToInt(index)
	if err != nil {
		return 0, err
	}

	// Confirm the key is in range of the slice / array
	if key < 0 || key > length - 1 {
		return 0, fmt.Errorf("index out of range: %d", key)
	}

	return key, nil
}

/*
Converts a `reflect.Value` to an `int64` if possible.
*/
func reflectHelperConvertToInt64(value reflect.Value) (int64, error) {
	var intValue int64

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue = value.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			intValue = int64(value.Uint())
		case reflect.Float32, reflect.Float64:
			intValue = int64(roundFloat(value.Float(), 0))
		case reflect.Bool:
			if value.Bool() {
				intValue = int64(1)
			} else {
				intValue = int64(0)
			}
		case reflect.String:
			str := value.String()
			tmp, err := strconv.Atoi(str)
			if err != nil {
				return int64(0), fmt.Errorf("can't convert type string to an int")
			}
			intValue = int64(tmp)
		case reflect.Invalid:
			return int64(0), fmt.Errorf("can't convert type nil to an int")
		default:
			return int64(0), fmt.Errorf("can't convert type %s to an int", value.Type())
	}

	return intValue, nil
}

/*
Converts a `reflect.Value` to a `uint64` if possible.
*/
func reflectHelperConvertToUint64(value reflect.Value) (uint64, error) {
	var uintValue uint64

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			uintValue = uint64(value.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			uintValue = value.Uint()
		case reflect.Float32, reflect.Float64:
			uintValue = uint64(roundFloat(value.Float(), 0))
		case reflect.Bool:
			if value.Bool() {
				uintValue = uint64(1)
			} else {
				uintValue = uint64(0)
			}
		case reflect.String:
			str := value.String()
			tmp, err := strconv.Atoi(str)
			if err != nil {
				return uint64(0), fmt.Errorf("can't convert type string to an int")
			}
			uintValue = uint64(tmp)
		case reflect.Invalid:
			return uint64(0), fmt.Errorf("can't convert type nil to an int")
		default:
			return uint64(0), fmt.Errorf("can't convert type %s to an int", value.Type())
	}

	return uintValue, nil
}

/*
Converts a `reflect.Value` to a `uint` if possible.
*/
func reflectHelperConvertToUint(value reflect.Value) (int, error) {
	intValue, err := reflectHelperConvertToInt64(value)
	return int(intValue), err
}

/*
Converts a `reflect.Value` to an `int` if possible.
*/
func reflectHelperConvertToInt(value reflect.Value) (int, error) {
	intValue, err := reflectHelperConvertToInt64(value)
	return int(intValue), err
}

/*
Converts a `reflect.Value` to a `float64` if possible.
*/
func reflectHelperConvertToFloat(value reflect.Value) (float64, error) {
	var floatValue float64

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			floatValue = float64(value.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			floatValue = float64(int64(value.Uint()))
		case reflect.Float32, reflect.Float64:
			floatValue = value.Float()
		case reflect.String:
			str := value.String()
			var err error
			floatValue, err = strconv.ParseFloat(str, 64)
			if err != nil {
				return 0, fmt.Errorf("can't convert type string to a float")
			}
		case reflect.Invalid:
			return 0, fmt.Errorf("can't convert type nil to a float")
		default:
			return 0, fmt.Errorf("can't convert type %s to a float", value.Type())
	}

	return floatValue, nil
}

/*
Converts a `reflect.Value` to a `string` if possible.
*/
func reflectHelperConvertToString(value reflect.Value) (string, error) {
	var stringValue string

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			stringValue = strconv.Itoa(int(value.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			stringValue = strconv.Itoa(int(value.Uint()))
		case reflect.Float32, reflect.Float64:
			stringValue = fmt.Sprintf("%.2f", value.Float())
		case reflect.Bool:
			stringValue = fmt.Sprintf("%v", value.Bool())
		case reflect.String:
			stringValue = value.String()
		case reflect.Invalid:
			return "", fmt.Errorf("can't convert type nil to a string")
		default:
			return "", fmt.Errorf("can't convert type %s to a string", value.Type())
	}

	return stringValue, nil
}

/*
Converts a `slice` to an `array`.
*/
func reflectHelperConvertSliceToArray(slice reflect.Value) (reflect.Value, error) {
	if slice.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("can't convert a type %s to an array", slice.Type())
	}

	t := slice.Type().Elem()
	t = reflect.ArrayOf(slice.Len(), t)
	arr := reflect.New(t).Elem()

	for i := 0; i < slice.Len(); i++ {
		v := arr.Index(i)
		v.Set(slice.Index(i))
	}

	return arr, nil
}

/*
Creates an empty slice to match the type of the value passed in.
*/
func reflectHelperCreateEmptySlice(value reflect.Value) (reflect.Value, error) {
	switch value.Kind() {
		case reflect.Array:
			return reflect.New(value.Type()).Elem().Slice(0, 0), nil
		case reflect.Slice:
			return reflect.MakeSlice(value.Type(), 0, 0), nil
	}

	return reflect.Value{}, fmt.Errorf("can't create slice from type %s", value.Type())
}

/*
Checks that the two values are of exactly the same types
*/
func reflectHelperStrictTypeCompatibility(value1 reflect.Value, value2 reflect.Value) error {
	if value1.Kind() != value2.Kind() || value1.Type() != value2.Type() {
		return fmt.Errorf("types do not match: %s vs %s", value1.Type(), value2.Type())
	}

	return nil
}

/*
Checks that the two values are of compatible types (e.g. any type of int, or roughly the right type of slice / array)
*/
func reflectHelperLooseTypeCompatibility(value1 reflect.Value, value2 reflect.Value) error {
	compatible := false

	switch value1.Kind() {
		case reflect.Bool:
			switch value2.Kind() {
				case reflect.Bool:
					compatible = true
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			switch value2.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					compatible = true
			}
		case reflect.Float32, reflect.Float64:
			switch value2.Kind() {
				case reflect.Float32, reflect.Float64:
					compatible = true
			}
		case reflect.String:
			switch value2.Kind() {
				case reflect.String:
					compatible = true
			}
		case reflect.Slice, reflect.Array:
			switch value2.Kind() {
				case reflect.Slice, reflect.Array:
					//if err := reflectHelperLooseTypeCompatibility(reflect.Zero(value1.Type().Elem()), reflect.Zero(value2.Type().Elem())); err == nil {
					if value1.Type().Elem() == value2.Type().Elem() {
						compatible = true
					}
			}
		case reflect.Map:
			switch value2.Kind() {
				case reflect.Map:
					if value1.Type().Key() == value2.Type().Key() {
						//if err := reflectHelperLooseTypeCompatibility(reflect.Zero(value1.Type().Elem()), reflect.Zero(value2.Type().Elem())); err == nil {
						if value1.Type().Elem() == value2.Type().Elem() {
							compatible = true
						}
					}
			}
		case reflect.Struct:
			// TODO 
			if value2.Kind() == reflect.Struct {
				compatible = true
			}
	}
	
	if !compatible {
		return fmt.Errorf("types do not match: %s vs %s", value1.Type(), value2.Type())
	}

	return nil
}