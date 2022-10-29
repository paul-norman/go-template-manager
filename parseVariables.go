package templateManager

/*
Functions dedicated to the parsing of variables embedded in templates (as strings) into their actual, strongly typed forms
*/

import (
	"regexp"
	"strings"
	"strconv"
)

// Determines a variable's (likely) basic type from a string representation of it
func getVariableType(value string, b ...string) (string, any) {
	bias := ""
	if len(b) > 0 {
		bias = b[0]
	}

	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		return "map", nil
	} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		return "slice", nil
	} else if val, err := strconv.Atoi(value); bias != "float" && err == nil {
		return "int", val
	} else if val, err := strconv.ParseFloat(value, 64); err == nil {
		return "float", val
	} else if val, err := strconv.ParseBool(strings.ToLower(value)); err == nil {
		return "bool", val
	}
	
	return "string", value
}

// Generic helper to parse supported map type combinations from a map of strings to a map of their actual types
func parseVariableMap[K string|int|float64|bool, V string|int|float64|bool](values map[string]string) map[K]V {
	tmp := map[K]V{}
	for key, val := range values {
		_, key := getVariableType(key)
		_, val := getVariableType(val)
		tmp[key.(K)] = val.(V)
	}

	return tmp
}

// Generic helper to parse supported slice types from a slice of strings to a slice of their actual type
func parseVariableSlice[V string|int|float64|bool](values []string) []V {
	t := ""
	var emptyV V
	var emptyVInterface any = emptyV
	if _, ok := emptyVInterface.(int); ok {
		t = "int"
	} else if _, ok := emptyVInterface.(float64); ok {
		t = "float"
	}
	
	tmp := []V{}
	for _, val := range values {
		_, v := getVariableType(val, t)
		tmp = append(tmp, v.(V))
	}

	return tmp
}

// Generic helper to parse supported nested slice types from a slice of strings to a slice of their actual type
func parseNestedVariableSlice[V string|int|float64|bool](values []string) [][]V {
	t := ""
	var emptyV V
	var emptyVInterface any = emptyV
	if _, ok := emptyVInterface.(int); ok {
		t = "int"
	} else if _, ok := emptyVInterface.(float64); ok {
		t = "float"
	}

	tmp := [][]V{}
	for _, val := range values {
		sub := []V{}
		tmp2 := prepareSlice(val)

		for _, subval := range tmp2 {
			_, v := getVariableType(subval, t)
			sub = append(sub, v.(V))
		}
		tmp = append(tmp, sub)
	}

	return tmp
}

// Parses a string representation of a map into string values ready for type detection
func prepareMap(value string) map[string]string {
	value = value[1:len(value) - 1]
	value = value + ","
	m := make(map[string]string)

	findStringStringMap, _ := regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")
	if findStringStringMap.MatchString(value) {
		matches := findStringStringMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[match[1]] = match[2]
		}

		return m
	}
	
	findStringNumericMap, _ := regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*([\\-\\.\\d]+)\\s*,")
	if findStringNumericMap.MatchString(value) {
		matches := findStringNumericMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[match[1]] = match[2]
		}

		return m
	}

	findStringBoolMap, _ := regexp.Compile("(?i)[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*(true|false)\\s*,")
	if findStringBoolMap.MatchString(value) {
		matches := findStringBoolMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[match[1]] = strings.ToLower(match[2])
		}

		return m
	}

	findNumericStringMap, _ := regexp.Compile("([\\-\\.\\d]+)\\s*:\\s*[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")
	if findNumericStringMap.MatchString(value) {
		matches := findNumericStringMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[match[1]] = match[2]
		}

		return m
	}

	findNumericNumericMap, _ := regexp.Compile(`([\-\.\d]+)\s*:\s*([\-\.\d]+)\s*,`)
	if findNumericNumericMap.MatchString(value) {
		matches := findNumericNumericMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[match[1]] = match[2]
		}

		return m
	}

	findNumericBoolMap, _ := regexp.Compile(`(?i)([\-\.\d]+)\s*:\s*(true|false)\s*,`)
	if findNumericBoolMap.MatchString(value) {
		matches := findNumericBoolMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[match[1]] = strings.ToLower(match[2])
		}

		return m
	}

	findBoolStringMap, _ := regexp.Compile("(?i)(true|false)\\s*:\\s*[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")
	if findBoolStringMap.MatchString(value) {
		matches := findBoolStringMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[strings.ToLower(match[1])] = match[2]
		}

		return m
	}

	findBoolNumericMap, _ := regexp.Compile(`(?i)(true|false)\s*:\s*([\-\.\d]+)\s*,`)
	if findBoolNumericMap.MatchString(value) {
		matches := findBoolNumericMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[strings.ToLower(match[1])] = match[2]
		}

		return m
	}

	findBoolBoolMap, _ := regexp.Compile(`(?i)(true|false)\s*:\s*(true|false)\s*,`)
	if findBoolBoolMap.MatchString(value) {
		matches := findBoolBoolMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[strings.ToLower(match[1])] = strings.ToLower(match[2])
		}

		return m
	}

	/* 
	findSliceMap, _ := regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*(\\[.*?\\])\\s*,")

	if findSliceMap.MatchString(value) {
		matches := findSliceMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[match[1]] = match[2] 
		}
	}
	*/

	return m
}

// Parses a string representation of a slice of booleans into a slice of strings ready for type detection
func prepareBooleanSlice(value string) []string {
	value = value + ","
	slice := []string{}

	findBooleanSlice, _	:= regexp.Compile(`(?i)\s*(true|false)\s*,`)
	if findBooleanSlice.MatchString(value) {
		matches := findBooleanSlice.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			slice = append(slice, strings.ToLower(match[1]))
		}
	}

	return slice
}

// Parses a string representation of a slice of numbers into a slice of strings ready for type detection
func prepareNumericSlice(value string) []string {
	value = value + ","
	slice := []string{}

	findNumericSlice, _	:= regexp.Compile(`\s*([\-\d\.]+)\s*,`)
	if findNumericSlice.MatchString(value) {
		matches := findNumericSlice.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			slice = append(slice, match[1])
		}
	}

	return slice
}

// Parses a string representation of a slice into a slice of string values ready for type detection
func prepareSlice(value string) []string {
	value = value[1:len(value) - 1]
	var values = []string{}

	if value[0:1] == "[" {
		values = prepareSliceSlice(value)
	} else if strings.Contains(value, `"`) || strings.Contains(value, `'`) || strings.Contains(value, "`") {
		values = prepareStringSlice(value)
	} else {
		values = prepareNumericSlice(value)
	}

	if len(values) == 0 {
		values = prepareBooleanSlice(value)
	}

	return values
}

// Parses a string representation of a slice of slices into a slice of strings ready for type detection
func prepareSliceSlice(value string) []string {
	value = value + ","
	slice := []string{}

	findSliceSlice, _ := regexp.Compile(`(\[.*?\])\s*,`)

	if findSliceSlice.MatchString(value) {
		matches := findSliceSlice.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			slice = append(slice, match[1])
		}
	}

	return slice
}

// Parses a string representation of a slice of strings into an actual slice of strings ready for type detection
func prepareStringSlice(value string) []string {
	value = value + ","
	slice := []string{}

	// No back-references in GoLang's RE2 regexp engine :-(
	findStringSlice, _	:= regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")

	if findStringSlice.MatchString(value) {
		matches := findStringSlice.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			val := strings.Replace(match[1], "\\`", "`", -1)
			val = strings.Replace(val, `\"`, `"`, -1)
			val = strings.Replace(val, "\\'", "'", -1)

			slice = append(slice, val)
		}
	}

	return slice
}