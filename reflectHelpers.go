package templateManager

/*
Functions to assist with reflection.
*/

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

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

/*
Checks if a value is empty.
*/
func reflectHelperIsEmpty(value reflect.Value) bool {
    return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

/*
Determines if a struct is empty.
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

/*
Makes a copy of the reflected element.
*/
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
			float, _ := reflectHelperConvertToFloat64(value)
			return reflect.ValueOf(float32(float)), nil
		case reflect.Float64:
			float, _ := reflectHelperConvertToFloat64(value)
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

/*
Makes a copy of the reflected struct.
*/
func reflectHelperStructCopy(value reflect.Value) (reflect.Value, error) {
	tmp := reflect.New(value.Type()).Elem()
	for i := 0; i < tmp.NumField(); i++ {
		if tmp.Field(i).CanSet() {
			tmp.Field(i).Set(value.Field(i))
		}
	}
	return tmp, nil
}

/*
Makes a copy of the reflected slice.
*/
func reflectHelperSliceCopy(value reflect.Value) (reflect.Value, error) {
	tmp, _ := reflectHelperCreateEmptySlice(value)
	for i := 0; i < value.Len(); i++ {
		val, _ := reflectHelperDeepCopy(value.Index(i))
		tmp = reflect.Append(tmp, val)
	}
	return tmp, nil
}

/*
Makes a copy of the reflected map.
*/
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
Returns a string representation of a map element's value type.
*/
func reflectHelperGetMapType(m reflect.Value) string {
	typ := m.Type().String()[4:]

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
Returns a string representation of the type that a slice / array contains
*/
func reflectHelperGetSliceType(slice reflect.Value) string {
	typ := slice.Type().String()

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
Converts the underlying type to another and returns a `reflect.Value` for it
Only works with simple types (i.e. numbers / bools / strings)
*/
func reflectHelperConvertUnderlying(value reflect.Value, to reflect.Kind) (reflect.Value, error) {
	t := value.Kind()
	if t == to {
		return value, nil
	}

	switch to {
		case reflect.Bool:
			bool, err := reflectHelperConvertToBool(value)
			return reflect.ValueOf(bool), err
		case reflect.Int:
			num, err := reflectHelperConvertToInt(value)
			return reflect.ValueOf(num), err
		case reflect.Int8:
			num, err := reflectHelperConvertToInt8(value)
			return reflect.ValueOf(num), err
		case reflect.Int16:
			num, err := reflectHelperConvertToInt16(value)
			return reflect.ValueOf(num), err
		case reflect.Int32:
			num, err := reflectHelperConvertToInt32(value)
			return reflect.ValueOf(num), err
		case reflect.Int64:
			num, err := reflectHelperConvertToInt64(value)
			return reflect.ValueOf(num), err
		case reflect.Uint:
			num, err := reflectHelperConvertToUint(value)
			return reflect.ValueOf(num), err
		case reflect.Uint8:
			num, err := reflectHelperConvertToUint8(value)
			return reflect.ValueOf(num), err
		case reflect.Uint16:
			num, err := reflectHelperConvertToUint16(value)
			return reflect.ValueOf(num), err
		case reflect.Uint32:
			num, err := reflectHelperConvertToUint32(value)
			return reflect.ValueOf(num), err
		case reflect.Uint64:
			num, err := reflectHelperConvertToUint64(value)
			return reflect.ValueOf(num), err
		case reflect.Float32:
			float, err := reflectHelperConvertToFloat32(value)
			return reflect.ValueOf(float), err
		case reflect.Float64:
			float, err := reflectHelperConvertToFloat64(value)
			return reflect.ValueOf(float), err
		case reflect.String:
			str, err := reflectHelperConvertToString(value)
			return reflect.ValueOf(str), err
		case reflect.Slice:
			switch t {
				case reflect.Bool:
					slice, err := reflectHelperConvertToSlice[bool](value)
					return reflect.ValueOf(slice), err
				case reflect.Int:
					slice, err := reflectHelperConvertToSlice[int](value)
					return reflect.ValueOf(slice), err
				case reflect.Int8:
					slice, err := reflectHelperConvertToSlice[int8](value)
					return reflect.ValueOf(slice), err
				case reflect.Int16:
					slice, err := reflectHelperConvertToSlice[int16](value)
					return reflect.ValueOf(slice), err
				case reflect.Int32:
					slice, err := reflectHelperConvertToSlice[int32](value)
					return reflect.ValueOf(slice), err
				case reflect.Int64:
					slice, err := reflectHelperConvertToSlice[int64](value)
					return reflect.ValueOf(slice), err
				case reflect.Uint:
					slice, err := reflectHelperConvertToSlice[uint](value)
					return reflect.ValueOf(slice), err
				case reflect.Uint8:
					slice, err := reflectHelperConvertToSlice[uint8](value)
					return reflect.ValueOf(slice), err
				case reflect.Uint16:
					slice, err := reflectHelperConvertToSlice[uint16](value)
					return reflect.ValueOf(slice), err
				case reflect.Uint32:
					slice, err := reflectHelperConvertToSlice[uint32](value)
					return reflect.ValueOf(slice), err
				case reflect.Uint64:
					slice, err := reflectHelperConvertToSlice[uint64](value)
					return reflect.ValueOf(slice), err
				case reflect.Float32:
					slice, err := reflectHelperConvertToSlice[float32](value)
					return reflect.ValueOf(slice), err
				case reflect.Float64:
					slice, err := reflectHelperConvertToSlice[float64](value)
					return reflect.ValueOf(slice), err
				case reflect.String:
					slice, err := reflectHelperConvertToSlice[string](value)
					return reflect.ValueOf(slice), err
				case reflect.Array:
					slice, err := reflectHelperConvertArrayToSlice(value)
					return reflect.ValueOf(slice), err
			}	
		case reflect.Array:
			switch t {
				case reflect.Bool:
					array, err := reflectHelperConvertToArray[bool](value)
					return reflect.ValueOf(array), err
				case reflect.Int:
					array, err := reflectHelperConvertToArray[int](value)
					return reflect.ValueOf(array), err
				case reflect.Int8:
					array, err := reflectHelperConvertToArray[int8](value)
					return reflect.ValueOf(array), err
				case reflect.Int16:
					array, err := reflectHelperConvertToArray[int16](value)
					return reflect.ValueOf(array), err
				case reflect.Int32:
					array, err := reflectHelperConvertToArray[int32](value)
					return reflect.ValueOf(array), err
				case reflect.Int64:
					array, err := reflectHelperConvertToArray[int64](value)
					return reflect.ValueOf(array), err
				case reflect.Uint:
					array, err := reflectHelperConvertToArray[uint](value)
					return reflect.ValueOf(array), err
				case reflect.Uint8:
					array, err := reflectHelperConvertToArray[uint8](value)
					return reflect.ValueOf(array), err
				case reflect.Uint16:
					array, err := reflectHelperConvertToArray[uint16](value)
					return reflect.ValueOf(array), err
				case reflect.Uint32:
					array, err := reflectHelperConvertToArray[uint32](value)
					return reflect.ValueOf(array), err
				case reflect.Uint64:
					array, err := reflectHelperConvertToArray[uint64](value)
					return reflect.ValueOf(array), err
				case reflect.Float32:
					array, err := reflectHelperConvertToArray[float32](value)
					return reflect.ValueOf(array), err
				case reflect.Float64:
					array, err := reflectHelperConvertToArray[float64](value)
					return reflect.ValueOf(array), err
				case reflect.String:
					array, err := reflectHelperConvertToArray[string](value)
					return reflect.ValueOf(array), err
				case reflect.Slice:
					array, err := reflectHelperConvertSliceToArray(value)
					return reflect.ValueOf(array), err
			}
	}

	return reflect.Value{}, fmt.Errorf("could not convert value of type: %v to type: %v", t, to)
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
Converts a `reflect.Value` to an `int` if possible.
*/
func reflectHelperConvertToInt(value reflect.Value) (int, error) {
	intValue, err := reflectHelperConvertToInt64(value)
	return int(intValue), err
}

/*
Converts a `reflect.Value` to an `int32` if possible.
*/
func reflectHelperConvertToInt32(value reflect.Value) (int32, error) {
	intValue, err := reflectHelperConvertToInt64(value)
	return int32(intValue), err
}

/*
Converts a `reflect.Value` to an `int16` if possible.
*/
func reflectHelperConvertToInt16(value reflect.Value) (int16, error) {
	intValue, err := reflectHelperConvertToInt64(value)
	return int16(intValue), err
}

/*
Converts a `reflect.Value` to an `int8` if possible.
*/
func reflectHelperConvertToInt8(value reflect.Value) (int8, error) {
	intValue, err := reflectHelperConvertToInt64(value)
	return int8(intValue), err
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
func reflectHelperConvertToUint(value reflect.Value) (uint, error) {
	intValue, err := reflectHelperConvertToUint64(value)
	return uint(intValue), err
}

/*
Converts a `reflect.Value` to a `uint32` if possible.
*/
func reflectHelperConvertToUint32(value reflect.Value) (uint32, error) {
	intValue, err := reflectHelperConvertToUint64(value)
	return uint32(intValue), err
}

/*
Converts a `reflect.Value` to a `uint16` if possible.
*/
func reflectHelperConvertToUint16(value reflect.Value) (uint16, error) {
	intValue, err := reflectHelperConvertToUint64(value)
	return uint16(intValue), err
}

/*
Converts a `reflect.Value` to a `uint8` if possible.
*/
func reflectHelperConvertToUint8(value reflect.Value) (uint8, error) {
	intValue, err := reflectHelperConvertToUint64(value)
	return uint8(intValue), err
}

/*
Converts a `reflect.Value` to a `float64` if possible.
*/
func reflectHelperConvertToFloat64(value reflect.Value) (float64, error) {
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
Converts a `reflect.Value` to a `float32` if possible.
*/
func reflectHelperConvertToFloat32(value reflect.Value) (float32, error) {
	floatValue, err := reflectHelperConvertToFloat64(value)
	return float32(floatValue), err
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
			stringValue = fmt.Sprintf("%v", value.Float())
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
Converts a `reflect.Value` to a `string` if possible flattening out slices, arrays, maps etc.
*/
func reflectHelperConvertAnythingToString(value reflect.Value) (string, error) {
	var stringValue string

	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			stringValue = strconv.Itoa(int(value.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			stringValue = strconv.Itoa(int(value.Uint()))
		case reflect.Float32, reflect.Float64:
			stringValue = fmt.Sprintf("%v", value.Float())
		case reflect.Bool:
			stringValue = fmt.Sprintf("%v", value.Bool())
		case reflect.String:
			stringValue = value.String()
		case reflect.Slice, reflect.Array:
			stringValue = ""
			for i := 0; i < value.Len(); i++ {
				val, err := reflectHelperConvertAnythingToString(value.Index(i))
				if err == nil {
					stringValue += val
				}
			}
		case reflect.Map:
			stringValue = ""
			iter := value.MapRange()
			for iter.Next() {
				val, err := reflectHelperConvertAnythingToString(iter.Value())
				if err == nil {
					stringValue += val
				}
			}
		case reflect.Struct:
			stringValue = ""
			for i := 0; i < value.NumField(); i++ {
				val, err := reflectHelperConvertAnythingToString(value.Field(i))
				if err == nil {
					stringValue += val
				}
			}
		case reflect.Invalid:
			return "", fmt.Errorf("can't convert type nil to a string")
		default:
			return "", fmt.Errorf("can't convert type %s to a string", value.Type())
	}

	return stringValue, nil
}

/*
Converts a `reflect.Value` to a `bool` if possible.
*/
func reflectHelperConvertToBool(value reflect.Value) (bool, error) {
	switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() > 0 {
				return true, nil
			}
			return false, nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if value.Uint() > 0 {
				return true, nil
			}
			return false, nil
		case reflect.Float32, reflect.Float64:
			if value.Float() >= 1 {
				return true, nil
			}
			return false, nil
		case reflect.Bool:
			return value.Bool(), nil
		case reflect.String:
			str := strings.ToLower(value.String())
			if str == "true" || str == "1" || str == "t" || str == "y" || str == "yes" {
				return true, nil
			} else if str == "false" || str == "0" || str == "f" || str == "n" || str == "no" {
				return false, nil
			}
			return false, fmt.Errorf("can't convert unknown string value to a bool")
		case reflect.Invalid:
			return false, fmt.Errorf("can't convert type nil to a bool")
		default:
			return false, fmt.Errorf("can't convert type %s to a bool", value.Type())
	}
}

/*
Converts a simple `reflect.Value` to a `slice` if possible.
*/
func reflectHelperConvertToSlice[T int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|float32|float64|string|bool](value reflect.Value) ([]T, error) {
	val, ok := value.Interface().(T)

	if !ok {
		return []T{}, fmt.Errorf("cannot convert type %T to a slice", value.Interface())
	}

	return []T{val}, nil
}

/*
Converts a simple `reflect.Value` to an `array` if possible.
*/
func reflectHelperConvertToArray[T int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|float32|float64|string|bool](value reflect.Value) ([1]T, error) {
	val, ok := value.Interface().(T)

	if !ok {
		return [1]T{}, fmt.Errorf("cannot convert type %T to an array", value.Interface())
	}

	return [1]T{val}, nil
}

/*
Converts an `array` to a `slice`.
*/
func reflectHelperConvertArrayToSlice(array reflect.Value) (reflect.Value, error) {
	if array.Kind() != reflect.Array {
		return reflect.Value{}, fmt.Errorf("can't convert a type %s to a slice", array.Type())
	}

	t := array.Type().Elem()
	t = reflect.SliceOf(t)
	slice := reflect.New(t).Elem()

	for i := 0; i < array.Len(); i++ {
		slice.Index(i).Set(array.Index(i))
	}

	return slice, nil
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
Checks that the two values are of loosely compatible types (e.g. any type of loosely matching numeric, or roughly the right type of slice / array)
*/
func reflectHelperLooseTypeCompatibility(value1 reflect.Value, value2 reflect.Value) error {
	compatible := false

	switch value1.Kind() {
		case reflect.Bool:
			switch value2.Kind() {
				case reflect.Bool:
					compatible = true
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			switch value2.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					compatible = true
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			switch value2.Kind() {
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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

/*
Checks that the two values are of very loosely compatible types (e.g. any type numeric, or roughly the right type of slice / array)
*/
func reflectHelperVeryLooseTypeCompatibility(value1 reflect.Value, value2 reflect.Value) error {
	compatible := false

	switch value1.Kind() {
		case reflect.Bool:
			switch value2.Kind() {
				case reflect.Bool:
					compatible = true
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			switch value2.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
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