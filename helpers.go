package templateManager

import (
	"fmt"
	"math"
	"reflect"
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
Floors floats to integers for numeric conversion.
*/
func floorFloat(float float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))

	return math.Floor(float * ratio) / ratio
}

/*
Ceils floats to integers for numeric conversion.
*/
func ceilFloat(float float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))

	return math.Ceil(float * ratio) / ratio
}

/*
Checks if two floats are equal (to a very small tolerance).
*/
func equalFloats(float1 float64, float2 float64) bool {
	diff := math.Abs(float1 - float2)

	return diff < 0.00000000001
}

/*
A helper that powers the 3 divide functions.
*/
func divideHelper(roundMethod reflect.Value, divisor reflect.Value, value reflect.Value) reflect.Value {
	sig		:= "divide" + roundMethod.String() + "(divisor int, value any)"
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
			switch roundMethod.String() {
				case "ceil": op = ceilFloat(op, 0)
				case "floor": op = floorFloat(op, 0)
				case "round": op = roundFloat(op, 0)
			}
			return reflect.ValueOf(int64(op)).Convert(value.Type())
		case reflect.Float32, reflect.Float64:
			val, _ := reflectHelperConvertToFloat64(value)
			op := val / div
			return reflect.ValueOf(op).Convert(value.Type())
		case reflect.String, reflect.Bool:
			logWarning(sig + fmt.Sprintf(" trying to divide a %s", value.Type()))
			return value
	}

	return recursiveHelper(value, reflect.ValueOf(divideHelper), roundMethod, divisor)
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