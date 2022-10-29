package templateManager

import (
	"testing"
)

func TestAAParseVariablesSetup(tester  *testing.T) {
	testsShowDetails	= true
	testsShowSuccessful = false
	logErrors			= false
	logWarnings			= false

	testFormatTitle("parseVariables")
}

func TestGetVariableType(tester *testing.T) {
	tests := map[string]struct{Type string; Value any}{
		"input": {"string", "input"},
		"1,000": {"string", "1,000"},
		"1": {Type: "int", Value: 1},
		"0": {Type: "int", Value: 0},
		"6540": {Type: "int", Value: 6540},
		"-3250": {Type: "int", Value: -3250},
		"1.0": {Type: "float", Value: 1.0},
		"0.0": {Type: "float", Value: 0.0},
		"6540.0": {Type: "float", Value: 6540.0},
		"-3250.0": {Type: "float", Value: -3250.0},
		"true": {Type: "bool", Value: true},
		"false": {Type: "bool", Value: false},
		"True": {Type: "bool", Value: true},
		"False": {Type: "bool", Value: false},
		"TRUE": {Type: "bool", Value: true},
		"FALSE": {Type: "bool", Value: false},
		"truE": {Type: "bool", Value: true},
		"falsE": {Type: "bool", Value: false},
		"[slice]": {"slice", nil},
		"{map}": {"map", nil},
	}

	passed, failed := 0, 0
	for input, expected := range tests {
		typ, val := getVariableType(input)

		if typ == expected.Type && val == expected.Value {
			passed++
		} else {
			tester.Errorf("\033[31mFAIL: \033[36m%s(%s\033[36m)\033[0m:\n\t\033[31mProduced: \033[33m%#v \033[36m%T\033[0m\n\t\033[31mExpected: \033[33m%#v \033[36m%T\033[0m", "getVariableType", input, val, val, expected.Value, expected.Value)
			failed++
		}
	}

	testFormatPassFail("getVariableType", passed, failed)
}

func TestParseVariableMap(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{map[string]string{"string": "test"}}, parseVariableMap[string, string](map[string]string{"string": "test"}), map[string]string{"string": "test"} },
		{ []any{map[string]string{"string": "12"}}, parseVariableMap[string, int](map[string]string{"string": "12"}), map[string]int{"string": 12} },
		{ []any{map[string]string{"string": "-12.32"}}, parseVariableMap[string, float64](map[string]string{"string": "-12.32"}), map[string]float64{"string": -12.32} },
		{ []any{map[string]string{"string": "TrUe"}}, parseVariableMap[string, bool](map[string]string{"string": "TrUe"}), map[string]bool{"string": true} },

		{ []any{map[string]string{"12": "test"}}, parseVariableMap[int, string](map[string]string{"12": "test"}), map[int]string{12: "test"} },
		{ []any{map[string]string{"12": "12"}}, parseVariableMap[int, int](map[string]string{"12": "12"}), map[int]int{12: 12} },
		{ []any{map[string]string{"12": "-12.32"}}, parseVariableMap[int, float64](map[string]string{"12": "-12.32"}), map[int]float64{12: -12.32} },
		{ []any{map[string]string{"12": "TrUe"}}, parseVariableMap[int, bool](map[string]string{"12": "TrUe"}), map[int]bool{12: true} },

		{ []any{map[string]string{"12.5": "test"}}, parseVariableMap[float64, string](map[string]string{"12.5": "test"}), map[float64]string{12.5: "test"} },
		{ []any{map[string]string{"12.5": "12"}}, parseVariableMap[float64, int](map[string]string{"12.5": "12"}), map[float64]int{12.5: 12} },
		{ []any{map[string]string{"12.5": "-12.32"}}, parseVariableMap[float64, float64](map[string]string{"12.5": "-12.32"}), map[float64]float64{12.5: -12.32} },
		{ []any{map[string]string{"12.5": "TrUe"}}, parseVariableMap[float64, bool](map[string]string{"12.5": "TrUe"}), map[float64]bool{12.5: true} },

		{ []any{map[string]string{"FaLsE": "test"}}, parseVariableMap[bool, string](map[string]string{"FaLsE": "test"}), map[bool]string{false: "test"} },
		{ []any{map[string]string{"FaLsE": "12"}}, parseVariableMap[bool, int](map[string]string{"FaLsE": "12"}), map[bool]int{false: 12} },
		{ []any{map[string]string{"FaLsE": "-12.32"}}, parseVariableMap[bool, float64](map[string]string{"FaLsE": "-12.32"}), map[bool]float64{false: -12.32} },
		{ []any{map[string]string{"FaLsE": "TrUe"}}, parseVariableMap[bool, bool](map[string]string{"FaLsE": "TrUe"}), map[bool]bool{false: true} },
	}

	testRunTests("parseVariableMap", tests, tester)
}

func TestParseVariableSlice(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{[]string{"string", "test"}}, parseVariableSlice[string]([]string{"string", "test"}), []string{"string", "test"} },
		{ []any{[]string{"10", "-40"}}, parseVariableSlice[int]([]string{"10", "-40"}), []int{10, -40} },
		{ []any{[]string{"10.0", "-40"}}, parseVariableSlice[float64]([]string{"10.0", "-40"}), []float64{10.0, -40} },
		{ []any{[]string{"10", "-40"}}, parseVariableSlice[float64]([]string{"10", "-40"}), []float64{10, -40} },
		{ []any{[]string{"true", "FaLsE"}}, parseVariableSlice[bool]([]string{"true", "FaLsE"}), []bool{true, false} },
	}

	testRunTests("parseVariableSlice", tests, tester)
}

func TestParseNestedVariableSlice(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{[]string{`["string"]`, `["test"]`}}, parseNestedVariableSlice[string]([]string{`["string"]`, `["test"]`}), [][]string{{"string"}, {"test"}} },
		{ []any{[]string{`["10"]`, `["-40"]`}}, parseNestedVariableSlice[int]([]string{`["10"]`, `["-40"]`}), [][]int{{10}, {-40}} },
		{ []any{[]string{`["10.0"]`, `["-40"]`}}, parseNestedVariableSlice[float64]([]string{`["10.0"]`, `["-40"]`}), [][]float64{{10}, {-40}} },
		{ []any{[]string{`["10"]`, `["-40"]`}}, parseNestedVariableSlice[float64]([]string{`["10"]`, `["-40"]`}), [][]float64{{10}, {-40}} },
		{ []any{[]string{`["true"]`, `["FaLsE"]`}}, parseNestedVariableSlice[bool]([]string{`["true"]`, `["FaLsE"]`}), [][]bool{{true}, {false}} },
	}

	testRunTests("parseNestedVariableSlice", tests, tester)
}

func TestPrepareMap(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{`{"key1": "value1", "key2": "value2"}`}, prepareMap(`{"key1": "value1", "key2": "value2"}`), map[string]string{"key1": "value1", "key2": "value2"} },
		{ []any{`{"key1": 1, "key2": -2}`}, prepareMap(`{"key1": 1, "key2": -2}`), map[string]string{"key1": "1", "key2": "-2"} },
		{ []any{`{"key1": 1.0, "key2": -2.0}`}, prepareMap(`{"key1": 1.0, "key2": -2.0}`), map[string]string{"key1": "1.0", "key2": "-2.0"} },
		{ []any{`{"key1": true, "key2": false}`}, prepareMap(`{"key1": true, "key2": false}`), map[string]string{"key1": "true", "key2": "false"} },

		{ []any{`{1: "value1", 2: "value2"}`}, prepareMap(`{1: "value1", 2: "value2"}`), map[string]string{"1": "value1", "2": "value2"} },
		{ []any{`{1: 1, 2: -2}`}, prepareMap(`{1: 1, 2: -2}`), map[string]string{"1": "1", "2": "-2"} },
		{ []any{`{1: 1.0, 2: -2.0}`}, prepareMap(`{1: 1.0, 2: -2.0}`), map[string]string{"1": "1.0", "2": "-2.0"} },
		{ []any{`{1: true, 2: false}`}, prepareMap(`{1: true, 2: false}`), map[string]string{"1": "true", "2": "false"} },

		{ []any{`{1.0: "value1", 2.0: "value2"}`}, prepareMap(`{1.0: "value1", 2.0: "value2"}`), map[string]string{"1.0": "value1", "2.0": "value2"} },
		{ []any{`{1.0: 1, 2.0: -2}`}, prepareMap(`{1.0: 1, 2.0: -2}`), map[string]string{"1.0": "1", "2.0": "-2"} },
		{ []any{`{1.0: 1.0, 2.0: -2.0}`}, prepareMap(`{1.0: 1.0, 2.0: -2.0}`), map[string]string{"1.0": "1.0", "2.0": "-2.0"} },
		{ []any{`{1.0: true, 2.0: false}`}, prepareMap(`{1.0: true, 2.0: false}`), map[string]string{"1.0": "true", "2.0": "false"} },

		{ []any{`{true: "value1", false: "value2"}`}, prepareMap(`{true: "value1", false: "value2"}`), map[string]string{"true": "value1", "false": "value2"} },
		{ []any{`{true: 1, false: -2}`}, prepareMap(`{true: 1, false: -2}`), map[string]string{"true": "1", "false": "-2"} },
		{ []any{`{true: 1.0, false: -2.0}`}, prepareMap(`{true: 1.0, false: -2.0}`), map[string]string{"true": "1.0", "false": "-2.0"} },
		{ []any{`{true: true, false: false}`}, prepareMap(`{true: true, false: false}`), map[string]string{"true": "true", "false": "false"} },
	}

	testRunTests("prepareMap", tests, tester)
}

func TestPrepareBooleanSlice(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{`true, false`}, prepareBooleanSlice(`true, false`), []string{"true", "false"} },
		{ []any{`True, FaLsE`}, prepareBooleanSlice(`True, FaLsE`), []string{"true", "false"} },
	}

	testRunTests("prepareBooleanSlice", tests, tester)
}

func TestPrepareNumericSlice(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{`30, -20`}, prepareNumericSlice(`30, -20`), []string{"30", "-20"} },
		{ []any{`1, -3.14`}, prepareNumericSlice(`1, -3.14`), []string{"1", "-3.14"} },
	}

	testRunTests("prepareNumericSlice", tests, tester)
}

func TestPrepareSlice(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{`["test", "values"]`}, prepareSlice(`["test", "values"]`), []string{"test", "values"} },
		{ []any{`['test', 'values']`}, prepareSlice(`['test', 'values']`), []string{"test", "values"} },
		{ []any{`["test', 'values"]`}, prepareSlice(`["test', 'values"]`), []string{"test", "values"} },
		{ []any{"[`test`, `values`]"}, prepareSlice("[`test`, `values`]"), []string{"test", "values"} },

		{ []any{`[30, -20]`}, prepareSlice(`[30, -20]`), []string{"30", "-20"} },
		{ []any{`[1, -3.14]`}, prepareSlice(`[1, -3.14]`), []string{"1", "-3.14"} },

		{ []any{`[true, false]`}, prepareSlice(`[true, false]`), []string{"true", "false"} },
		{ []any{`[True, FaLsE]`}, prepareSlice(`[True, FaLsE]`), []string{"true", "false"} },

		{ []any{`[["test", "slice"], ["nesting", "here"]]`}, prepareSlice(`[["test", "slice"], ["nesting", "here"]]`), []string{`["test", "slice"]`, `["nesting", "here"]`} },
		{ []any{`[[30, -20], [10, -40]]`}, prepareSlice(`[[30, -20], [10, -40]]`), []string{"[30, -20]", "[10, -40]"} },
		{ []any{`[[30.0, -20.4], [10.234, -40.987]]`}, prepareSlice(`[[30.0, -20.4], [10.234, -40.987]]`), []string{"[30.0, -20.4]", "[10.234, -40.987]"} },
		{ []any{`[[true, false], [TruE, FALSE]]`}, prepareSlice(`[[true, false], [TruE, FALSE]]`), []string{`[true, false]`, `[TruE, FALSE]`} },
	}

	testRunTests("prepareSlice", tests, tester)
}

func TestPrepareSliceSlice(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{`["test", "slice"], ["nesting", "here"]`}, prepareSliceSlice(`["test", "slice"], ["nesting", "here"]`), []string{`["test", "slice"]`, `["nesting", "here"]`} },
		{ []any{`[30, -20], [10, -40]`}, prepareSliceSlice(`[30, -20], [10, -40]`), []string{"[30, -20]", "[10, -40]"} },
		{ []any{`[30.0, -20.4], [10.234, -40.987]`}, prepareSliceSlice(`[30.0, -20.4], [10.234, -40.987]`), []string{"[30.0, -20.4]", "[10.234, -40.987]"} },
		{ []any{`[true, false], [TruE, FALSE]`}, prepareSliceSlice(`[true, false], [TruE, FALSE]`), []string{`[true, false]`, `[TruE, FALSE]`} },
	}

	testRunTests("prepareSliceSlice", tests, tester)
}

func TestPrepareStringSlice(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{`"test", "values"`}, prepareStringSlice(`"test", "values"`), []string{"test", "values"} },
		{ []any{`'test', 'values'`}, prepareStringSlice(`'test', 'values'`), []string{"test", "values"} },
		{ []any{`"test', 'values"`}, prepareStringSlice(`"test', 'values"`), []string{"test", "values"} },
		{ []any{"`test`, `values`"}, prepareStringSlice("`test`, `values`"), []string{"test", "values"} },
	}

	testRunTests("prepareStringSlice", tests, tester)
}