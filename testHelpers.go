package templateManager

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// Call the reflection-based function that has 1 argument
func testCall1Arg[T reflect.Value|any](tester *testing.T, fn func(reflect.Value) (T, error), value1 any, expected any) bool {
	result, _ := fn(reflect.ValueOf(value1))
	arguments := fmt.Sprintf("(\033[33m%#v\033[36m)", value1)
	return testDeepEqual(tester, arguments, result, expected, reflect.ValueOf(fn))
}

// Call the reflection-based function that has 2 arguments
func testCall2Args[T reflect.Value|any](tester *testing.T, fn func(reflect.Value, reflect.Value) (T, error), value1 any, value2 any, expected any) bool {
	result, _ := fn(reflect.ValueOf(value1), reflect.ValueOf(value2))
	arguments := fmt.Sprintf("(\033[33m%#v\033[0m, \033[33m%#v\033[36m)", value1, value2)
	return testDeepEqual(tester, arguments, result, expected, reflect.ValueOf(fn))
}

// Call the reflection-based function that has 3 arguments 
func testCall3Args[T reflect.Value|any](tester *testing.T, fn func(reflect.Value, reflect.Value, reflect.Value) (T, error), value1 any, value2 any, value3 any, expected any) bool {
	result, _ := fn(reflect.ValueOf(value1), reflect.ValueOf(value2), reflect.ValueOf(value3))
	arguments := fmt.Sprintf("(\033[33m%#v\033[0m, \033[33m%#v\033[0m, \033[33m%#v\033[36m)", value1, value2, value3)
	return testDeepEqual(tester, arguments, result, expected, reflect.ValueOf(fn))
}

// Call the reflection-based function that has variable arguments
func testCallVarArgs[T reflect.Value|any](tester *testing.T, fn func(...reflect.Value) (T, error), values []any, expected any) bool {
	tmp := []reflect.Value{}
	arguments := "("
	for i, v := range values {
		tmp = append(tmp, reflect.ValueOf(v))
		if i > 0 { arguments += ", " }
		arguments += fmt.Sprintf("\033[33m%#v\033[36m", v)
	}
	arguments += ")"
	result, _ := fn(tmp...)
	return testDeepEqual(tester, arguments, result, expected, reflect.ValueOf(fn))
}

// Test that two values (`reflect.Value` or `any`) are equal
func testDeepEqual[T reflect.Value|any](tester *testing.T, arguments string, result T, expected any, fn reflect.Value) bool {
	test := any(result)
	switch v := test.(type) {
		case reflect.Value:
			if v.IsValid() {
				test = v.Interface()
			} else {
				test = nil
			}
		default: test = v
	}
	
	if !reflect.DeepEqual(test, expected) {
		tmp := strings.Split(runtime.FuncForPC(fn.Pointer()).Name(), ".")
		name := tmp[len(tmp) - 1]
		tester.Errorf("\033[31mFAIL: \033[36m%s%s:\n\t\033[31mProduced: \033[33m%#v \033[36m%T\033[0m\n\t\033[31mExpected: \033[33m%#v \033[36m%T\033[0m", name, arguments,  test, test, expected, expected)

		return false
	} else {
		if testsShowSuccessful {
			tmp := strings.Split(runtime.FuncForPC(fn.Pointer()).Name(), ".")
			name := tmp[len(tmp) - 1]
			fmt.Printf("\t\033[32mPASSED: \033[36m%s(%s\033[36m)\033[0m:\n\t\tProduced: \033[33m%#v \033[36m%T\033[0m\n", name, arguments, test, test)
		}
	}

	return true
}

// Neatly outputs whether a certain test passed or failed
func testFormatPassFail(name string, passed int, failed int) {
	if testsShowDetails {
		title := fmt.Sprintf("Running %s() Tests:", name)
		fmt.Printf("\033[36m%-50s\033[0m ", title)
		if failed == 0 {
			fmt.Printf("\033[32mPASSED: %d/%d\033[0m\n", passed, passed + failed)
			return
		}
		fmt.Printf("\033[33mPASSED: %d, \033[31mFAILED: %d\033[0m\n", passed, failed)
	}
}

// Outputs the test title
func testFormatTitle(name string) {
	if testsShowDetails {
		fmt.Printf("\nRunning: %s_test.go:\n\n", name)
	}
}

// Runs tests with 1 `reflect.Value` input and a single `reflect.Value` / interface output
func testRun1ArgTests[T reflect.Value|any](fn func(reflect.Value) (T, error), tests []struct{ input1, expected any }, tester *testing.T) {
	passed, failed := 0, 0
	for _, test := range tests {
		if testCall1Arg(tester, fn, test.input1, test.expected) {
			passed++
		} else { failed++ }
	}

	tmp := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	name := tmp[len(tmp) - 1]

	testFormatPassFail(name, passed, failed)
}

// Runs tests with 2 `reflect.Value` inputs and a single `reflect.Value` / interface output
func testRun2ArgTests[T reflect.Value|any](fn func(reflect.Value, reflect.Value) (T, error), tests []struct{ input1, input2, expected any }, tester *testing.T) {
	passed, failed := 0, 0
	for _, test := range tests {
		if testCall2Args(tester, fn, test.input1, test.input2, test.expected) {
			passed++
		} else { failed++ }
	}

	tmp := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	name := tmp[len(tmp) - 1]

	testFormatPassFail(name, passed, failed)
}

// Runs tests with 3 `reflect.Value` inputs and a single `reflect.Value` / interface output
func testRun3ArgTests[T reflect.Value|any](fn func(reflect.Value, reflect.Value, reflect.Value) (T, error), tests []struct{ input1, input2, input3, expected any }, tester *testing.T) {
	passed, failed := 0, 0
	for _, test := range tests {
		if testCall3Args(tester, fn, test.input1, test.input2, test.input3, test.expected) {
			passed++
		} else { failed++ }
	}

	tmp := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	name := tmp[len(tmp) - 1]

	testFormatPassFail(name, passed, failed)
}

// Runs tests with unlimited `reflect.Value` inputs and a single `reflect.Value` / interface output
func testRunVarArgTests[T reflect.Value|any](fn func(...reflect.Value) (T, error), tests []struct{ inputs []any; expected any }, tester *testing.T) {
	passed, failed := 0, 0
	for _, test := range tests {
		if testCallVarArgs(tester, fn, test.inputs, test.expected) {
			passed++
		} else { failed++ }
	}

	tmp := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	name := tmp[len(tmp) - 1]

	testFormatPassFail(name, passed, failed)
}

// Automatically determine which tests should fun for the variables passed in
func testRunArgTests(fn any, tests any, tester *testing.T) {
	if real, ok := fn.(func(reflect.Value) (reflect.Value, error)); ok {
		testRun1ArgTests(real, tests.([]struct{ input1, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value) (bool, error)); ok {
		testRun1ArgTests(real, tests.([]struct{ input1, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value) (int, error)); ok {
		testRun1ArgTests(real, tests.([]struct{ input1, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value) (float64, error)); ok {
		testRun1ArgTests(real, tests.([]struct{ input1, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value) (string, error)); ok {
		testRun1ArgTests(real, tests.([]struct{ input1, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value) (reflect.Value, error)); ok {
		testRun2ArgTests(real, tests.([]struct{ input1, input2, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value) (bool, error)); ok {
		testRun2ArgTests(real, tests.([]struct{ input1, input2, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value) (int, error)); ok {
		testRun2ArgTests(real, tests.([]struct{ input1, input2, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value) (float64, error)); ok {
		testRun2ArgTests(real, tests.([]struct{ input1, input2, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value) (string, error)); ok {
		testRun2ArgTests(real, tests.([]struct{ input1, input2, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value, reflect.Value) (reflect.Value, error)); ok {
		testRun3ArgTests(real, tests.([]struct{ input1, input2, input3, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value, reflect.Value) (bool, error)); ok {
		testRun3ArgTests(real, tests.([]struct{ input1, input2, input3, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value, reflect.Value) (int, error)); ok {
		testRun3ArgTests(real, tests.([]struct{ input1, input2, input3, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value, reflect.Value) (float64, error)); ok {
		testRun3ArgTests(real, tests.([]struct{ input1, input2, input3, expected any }), tester)
	} else if real, ok := fn.(func(reflect.Value, reflect.Value, reflect.Value) (string, error)); ok {
		testRun3ArgTests(real, tests.([]struct{ input1, input2, input3, expected any }), tester)
	} else if real, ok := fn.(func(...reflect.Value) (reflect.Value, error)); ok {
		testRunVarArgTests(real, tests.([]struct{ inputs []any; expected any }), tester)
	} else {
		tester.Errorf("\033[31mFAIL: \033[0mTests could not run!!")	
	}
}

// Runs simple tests with any inputs and a single output (result must be calculated prior to passing in)
func testRunTests(name string, tests []struct { inputs []any; result any; expected any }, tester *testing.T) {
	passed, failed := 0, 0
	for _, test := range tests {
		arguments := ""
		for i, input := range test.inputs {
			if i > 0 { arguments += ", " }
			arguments += fmt.Sprintf("\033[33m%#v\033[0m", input)
		}

		if !reflect.DeepEqual(test.result, test.expected) {
			tester.Errorf("\033[31mFAIL: \033[36m%s(%s\033[36m)\033[0m:\n\t\033[31mProduced: \033[33m%#v \033[36m%T\033[0m\n\t\033[31mExpected: \033[33m%#v \033[36m%T\033[0m", name, arguments, test.result, test.result, test.expected, test.expected)
			failed++
		} else {
			if testsShowSuccessful {
				fmt.Printf("\t\033[32mPASSED: \033[36m%s(%s\033[36m)\033[0m:\n\t\tProduced: \033[33m%#v \033[36m%T\033[0m\n", name, arguments, test.result, test.result)
			}
			passed++
		}
	}

	testFormatPassFail(name, passed, failed)
}