package templateManager

/*
Functions to assist with interfaces
*/

import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"fmt"
)

/*
Gets the "type" of the interface (e.g. int8, map[string]int or templateManager.TemplateManager)
*/
func interfaceHelperType(value any) string {
	switch value.(type) {
		case int, *int:
			return "int"
		case int8, *int8:
			return "int8"
		case int16, *int16:
			return "int16"
		case int32, *int32:
			return "int32"
		case int64, *int64:
			return "int64"
		case uint, *uint:
			return "uint"
		case uint8, *uint8:
			return "uint8"
		case uint16, *uint16:
			return "uint16"
		case uint32, *uint32:
			return "uint32"
		case uint64, *uint64:
			return "uint64"
		case uintptr, *uintptr:
			return "uintptr"
		case float32, *float32:
			return "float32"
		case float64, *float64:
			return "float64"
		case bool, *bool:
			return "bool"
		case string, *string:
			return "string"
	}

	typ := fmt.Sprintf("%T", value)

	// Not concerned about pointers here
	if string(typ[0]) == "*" {
		typ = typ[1:]
	}

	return typ
}

/*
Gets the underlying "kind" of the interface (e.g. int8, map or struct)
*/
func interfaceHelperKind(value any) string {
	switch value.(type) {
		case int, *int:
			return "int"
		case int8, *int8:
			return "int8"
		case int16, *int16:
			return "int16"
		case int32, *int32:
			return "int32"
		case int64, *int64:
			return "int64"
		case uint, *uint:
			return "uint"
		case uint8, *uint8:
			return "uint8"
		case uint16, *uint16:
			return "uint16"
		case uint32, *uint32:
			return "uint32"
		case uint64, *uint64:
			return "uint64"
		case uintptr, *uintptr:
			return "uintptr"
		case float32, *float32:
			return "float32"
		case float64, *float64:
			return "float64"
		case bool, *bool:
			return "bool"
		case string, *string:
			return "string"
	}

	kind := fmt.Sprintf("%T", value)

	// Not concerned about pointers here
	if string(kind[0]) == "*" {
		kind = kind[1:]
	}

	if found, err := interfaceHelperParseKind(kind); err == nil {
		return found
	}

	// Custom type, need reflection
	kind = reflect.ValueOf(value).Kind().String()
	if kind == "ptr" {
		kind = reflect.ValueOf(value).Elem().Kind().String()
	}
	
	if found, err := interfaceHelperParseKind(kind); err == nil {
		return found
	}

	return kind
}

/*
Checks if the kind string is a more complex standard type and simplifies its representation
*/
func interfaceHelperParseKind(kind string) (string, error) {
	if strings.HasPrefix(kind, "map[") {
		return "map", nil
	} else if strings.HasPrefix(kind, "struct {") {
		return "struct", nil
	} else if string(kind[0]) == "[" {
		if string(kind[1]) != "]" {
			return "array", nil
		}
		return "slice", nil
	}

	return "", fmt.Errorf("not a standard complex type")
}

/*
Creates the absolute value of a float
*/
func interfaceHelperAbsFloat64(value any) (float64, error) {
	switch v := value.(type) {
		case int:
			val := float64(v)
			if v < 0 { return -val, nil }
			return val, nil
		case int8:
			val := float64(v)
			if v < 0 { return -val, nil }
			return val, nil
		case int16:
			val := float64(v)
			if v < 0 { return -val, nil }
			return val, nil
		case int32:
			val := float64(v)
			if v < 0 { return -val, nil }
			return val, nil
		case int64:
			val := float64(v)
			if v < 0 { return -val, nil }
			return val, nil
		case uint:
			return float64(v), nil
		case uint8:
			return float64(v), nil
		case uint16:
			return float64(v), nil
		case uint32:
			return float64(v), nil
		case uint64:
			return float64(v), nil
		case uintptr:
			return float64(v), nil
		case float32:
			return math.Abs(roundFloat(float64(v), 0)), nil
		case float64:
			return math.Abs(roundFloat(v, 0)), nil
		case bool:
			if v {
				return float64(1), nil
			} 
			return float64(0), nil
		case string:
			val, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return float64(0), fmt.Errorf("can't convert type string to a float")
			}
			if val < 0 { return -val, nil }
			return val, nil
	}

	return float64(0), fmt.Errorf("can't convert type %T to an float64", value)
}

/*
Creates the absolute value of an integer
*/
func interfaceHelperAbsInt64(value any) (int64, error) {
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
Converts an `interface{}` to a `uint64` if possible.
*/
func interfaceHelperConvertToUint64(value any) (uint64, error) {
	var intValue uint64

	switch v := value.(type) {
		case int:
			tmp, _ := interfaceHelperAbsInt64(v)
			intValue = uint64(tmp)
		case int8:
			tmp, _ := interfaceHelperAbsInt64(v)
			intValue = uint64(tmp)
		case int16:
			tmp, _ := interfaceHelperAbsInt64(v)
			intValue = uint64(tmp)
		case int32:
			tmp, _ := interfaceHelperAbsInt64(v)
			intValue = uint64(tmp)
		case int64:
			tmp, _ := interfaceHelperAbsInt64(v)
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
			tmp, _ := interfaceHelperAbsInt64(v)
			intValue = uint64(tmp)
		case string:
			tmp, _ := interfaceHelperAbsInt64(v)
			intValue = uint64(tmp)
		default:
			return uint64(0), fmt.Errorf("can't convert type %T to a uint", value)
	}

	return intValue, nil
}

/*
Converts an `interface{}` to a `uint` if possible.
*/
func interfaceHelperConvertToUint(value any) (uint, error) {
	int64Value, err := interfaceHelperConvertToUint64(value)

	return uint(int64Value), err
}

/*
Converts an `interface{}` to a `string` if possible.
*/
func interfaceHelperConvertToString(value any) (string, error) {
	var stringValue string

	switch v := value.(type) {
		case int:
			stringValue = strconv.Itoa(v)
		case int8:
			stringValue = strconv.FormatInt(int64(v), 10)
		case int16:
			stringValue = strconv.FormatInt(int64(v), 10)
		case int32:
			stringValue = strconv.FormatInt(int64(v), 10)
		case int64:
			stringValue = strconv.FormatInt(v, 10)
		case uint:
			stringValue = strconv.FormatUint(uint64(v), 10)
		case uint8:
			stringValue = strconv.FormatUint(uint64(v), 10)
		case uint16:
			stringValue = strconv.FormatUint(uint64(v), 10)
		case uint32:
			stringValue = strconv.FormatUint(uint64(v), 10)
		case uint64:
			stringValue = strconv.FormatUint(v, 10)
		case uintptr:
			stringValue = strconv.FormatUint(uint64(v), 10)
		case float32:
			stringValue = strconv.FormatFloat(float64(v), 'g', -1, 64)
		case float64:
			stringValue = strconv.FormatFloat(v, 'g', -1, 64)
		case bool:
			stringValue = strconv.FormatBool(v)
		case string:
			stringValue = v
		default:
			return "", fmt.Errorf("can't convert type %T to a string", value)
	}

	return stringValue, nil
}