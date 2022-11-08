package templateManager

import (
	"reflect"
	"testing"
	"time"
)

func TestAAFunctionsSetup(tester  *testing.T) {
	testsShowDetails	= true
	testsShowSuccessful = false
	logErrors			= false
	logWarnings			= false

	testFormatTitle("functions")
}

func TestAdd(tester *testing.T) {
	tests := []struct{ input1, input2, expected any } {
		{ false, true, true },
		{ true, false, false },
		{ int8(10), int8(20), int8(30) },
		{ 10, 20, 30 },
		{ -10, 20, 10 },
		{ -10, -20, -30 },
		{ 10.1, 20.2, 30.299999999999997 },
		{ -10.1, 20.1, 10.000000000000002 },
		{ -10.1, -20.2, -30.299999999999997 },
		{ "add", "to", "toadd" },
		{ []int{5}, []int{10}, []int{10, 5} },
		{ 5, []int{10, 5}, []int{15, 10} },
		{ []string{"add"}, []string{"to"}, []string{"to", "add"} },
		{ "add", []string{"to"}, []string{"toadd"} },
		{ [1]string{"add"}, [1]string{"to"}, [2]string{"to", "add"} },
		{ "add", [1]string{"to"}, [1]string{"toadd"} },
		{ map[string]string{"add": "add-value"}, map[string]string{"to": "to-value"}, map[string]string{"add": "add-value", "to": "to-value"} },
		{ "add-value", map[string]string{"to": "to-value"}, map[string]string{"to": "to-valueadd-value"} },
		{ struct{ str string }{"add"}, struct{ string }{}, struct{ string }{} },
		{ struct{ Str string }{"add"}, struct{ string }{"to"}, struct{ string }{"to"} },
	}

	testRunArgTests(add, tests, tester)
}

func TestCapfirst(tester *testing.T) {
	tests := []struct{ input1, expected any } {
		{ true, true },
		{ false, false },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ "hello world", "Hello world" },
		{ "123 hello world", "123 Hello world" },
		{ []string{"hello world"}, []string{"Hello world"} },
		{ []string{"123 hello world"}, []string{"123 Hello world"} },
		{ map[string]string{"test": "hello world"}, map[string]string{"test": "Hello world"} },
		{ map[string]string{"test": "123 hello world"}, map[string]string{"test": "123 Hello world"} },
		{ struct{Str string}{"test"}, struct{Str string}{"Test"} },
		{ struct{str string}{"test"}, struct{str string}{""} },
	}

	testRunArgTests(capfirst, tests, tester)
}

func TestCollection(tester *testing.T) {
	tests := []struct{ inputs []any; result any; expected any } {
		{ []any{}, collection(), map[string]any{} },
		{ []any{0}, collection(0), map[string]any{} },
		{ []any{0, 0}, collection(0, 0), map[string]any{} },
		{ []any{"var", 0}, collection("var", 0), map[string]any{ "var": 0 } },
		{ []any{"var1", 0, "var2", true}, collection("var1", 0, "var2", true), map[string]any{ "var1": 0, "var2": true } },
	}

	testRunTests("collection", tests, tester)
}

func TestContains(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ "world", "hello world", true },
		{ "World", "hello world", false },
		{ "hello world", []string{"hello world"}, true },
		{ "world", []string{"hello world"}, false },
		{ "hello world", map[string]string{"test": "hello world"}, true },
		{ "world", map[string]string{"test": "hello world"}, false },
		{ "hello world", struct{ Str1, Str2 string }{ "test", "hello world" }, true },
		{ "world", struct{ Str1, Str2 string }{ "test", "hello world" }, false },
		{ "hello world", struct{ str1, str2 string }{ "test", "hello world" }, true },
		{ "world", struct{ str1, str2 string }{ "test", "hello world" }, false },
	}

	testRunArgTests(contains, tests, tester)
}

func TestCut(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ "world", "hello world", "hello " },
		{ "World", "hello world", "hello world" },
		{ "hello world", []string{"hello world"}, []string{""} },
		{ "world", []string{"hello world"}, []string{"hello "} },
		{ "hello world", map[string]string{"test": "hello world"}, map[string]string{"test": ""} },
		{ "world", map[string]string{"test": "hello world"}, map[string]string{"test": "hello "} },
		{ "hello world", struct{ Str1, Str2 string }{"test", "hello world"}, struct{ Str1, Str2 string }{"test", ""} },
		{ "world", struct{ Str1, Str2 string }{"test", "hello world"}, struct{ Str1, Str2 string }{"test", "hello "} },
		{ "hello world", struct{ str1, str2 string }{"test", "hello world"}, struct{ str1, str2 string }{"", ""} },
		{ "world", struct{ str1, str2 string }{"test", "hello world"}, struct{ str1, str2 string }{"", ""} },
	}

	testRunArgTests(cut, tests, tester)
}

func TestDate(tester *testing.T) {
	currentTime			:= time.Now().In(dateLocalTimezone)
	testTimeISO8601Z	:= "2019-04-23T11:30:21+01:00"
	testTimeISO8601		:= "2019-04-23T11:30:21+01:00"
	testTimeRFC822Z		:= "Tue, 23 Apr 19 11:30:21 +01:00"
	testTimeRFC822		:= "Tue, 23 Apr 19 11:30:21 CET"
	testTimeRFC850		:= "Tuesday, 23-Apr-19 11:30:21 CET"
	testTimeRFC1036		:= "Tue, 23 Apr 19 11:30:21 +01:00"
	testTimeRFC1123Z	:= "Tue, 23 Apr 2019 11:30:21 +01:00"
	testTimeRFC1123		:= "Tue, 23 Apr 2019 11:30:21 CET"
	testTimeRFC2822		:= "Tue, 23 Apr 2019 11:30:21 +01:00"
	testTimeRFC3339		:= "2019-04-23T11:30:21+01:00"
	testTimeWC3			:= "2019-04-23T11:30:21+01:00"
	testTimeATOM		:= "2019-04-23T11:30:21+01:00"
	testTimeCOOKIE		:= "Tuesday, 23-Apr-2019 11:30:21 CET"
	testTimeRSS			:= "Tue, 23 Apr 2019 11:30:21 +01:00"
	testTimeMYSQL		:= "2019-04-23 11:30:21"
	testTimeUNIX		:= "Tue Apr 23 11:30:21 CET 2019"
	testTimeRUBY		:= "Tue Apr 23 11:30:21 +0100 2019"
	testTimeANSIC		:= "Tue Apr 23 11:30:21 2019"

	testTime, _ := time.Parse(time.RFC3339, testTimeRFC3339)
	testTime = testTime.In(dateLocalTimezone)

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, date(), currentTime.Format("02/01/2006") },

		{ []any{testTime}, date(testTime), testTime.Format("02/01/2006") },
		{ []any{1556015421}, date(1556015421), testTime.Format("02/01/2006") },

		{ []any{"02-01-2006"}, date("02-01-2006"), currentTime.Format("02-01-2006") },
		{ []any{"d-m-Y"}, date("d-m-Y"), currentTime.Format("02-01-2006") },
		{ []any{"%d-%m-%Y"}, date("%d-%m-%Y"), currentTime.Format("02-01-2006") },
		{ []any{"Mon 02 Jan 06"}, date("Mon 02 Jan 06"), currentTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y"}, date("D d M y"), currentTime.Format("Mon 02 Jan 06") },
		{ []any{"%a %d %b %y"}, date("%a %d %b %y"), currentTime.Format("Mon 02 Jan 06") },

		{ []any{"Mon 02 Jan 06", 1556015421}, date("Mon 02 Jan 06", 1556015421), testTime.Format("Mon 02 Jan 06") },

		{ []any{"Mon 02 Jan 06", testTime}, date("Mon 02 Jan 06", testTime), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", testTime}, date("D d M y", testTime), testTime.Format("Mon 02 Jan 06") },
		{ []any{"%a %d %b %y", testTime}, date("%a %d %b %y", testTime), testTime.Format("Mon 02 Jan 06") },

		{ []any{"Mon 02 Jan 06", testTimeRFC3339}, date("Mon 02 Jan 06", testTimeRFC3339), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", testTimeRFC3339}, date("D d M y", testTimeRFC3339), testTime.Format("Mon 02 Jan 06") },
		{ []any{"%a %d %b %y", testTimeRFC3339}, date("%a %d %b %y", testTimeRFC3339), testTime.Format("Mon 02 Jan 06") },

		{ []any{"Mon 02 Jan 06", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, date("Mon 02 Jan 06", "2006-01-02T15:04:05Z07:00", testTimeRFC3339), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, date("D d M y", "2006-01-02T15:04:05Z07:00", testTimeRFC3339), testTime.Format("Mon 02 Jan 06") },
		{ []any{"%a %d %b %y", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, date("%a %d %b %y", "2006-01-02T15:04:05Z07:00", testTimeRFC3339), testTime.Format("Mon 02 Jan 06") },

		{ []any{"D d M y", "ISO8601Z", testTimeISO8601Z}, date("D d M y", "ISO8601Z", testTimeISO8601Z), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "ISO8601", testTimeISO8601}, date("D d M y", "ISO8601", testTimeISO8601), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC822Z", testTimeRFC822Z}, date("D d M y", "RFC822Z", testTimeRFC822Z), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC822", testTimeRFC822}, date("D d M y", "RFC822", testTimeRFC822), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC850", testTimeRFC850}, date("D d M y", "RFC850", testTimeRFC850), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC1036", testTimeRFC1036}, date("D d M y", "RFC1036", testTimeRFC1036), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC1123Z", testTimeRFC1123Z}, date("D d M y", "RFC1123Z", testTimeRFC1123Z), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC1123", testTimeRFC1123}, date("D d M y", "RFC1123", testTimeRFC1123), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC2822", testTimeRFC2822}, date("D d M y", "RFC2822", testTimeRFC2822), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC3339", testTimeRFC3339}, date("D d M y", "RFC3339", testTimeRFC3339), testTime.Format("Mon 02 Jan 06") },

		{ []any{"D d M y", "ATOM", testTimeATOM}, date("D d M y", "ATOM", testTimeATOM), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "W3C", testTimeWC3}, date("D d M y", "W3C", testTimeWC3), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "COOKIE", testTimeCOOKIE}, date("D d M y", "COOKIE", testTimeCOOKIE), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RSS", testTimeRSS}, date("D d M y", "RSS", testTimeRSS), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "MYSQL", testTimeMYSQL}, date("D d M y", "MYSQL", testTimeMYSQL), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "UNIX", testTimeUNIX}, date("D d M y", "UNIX", testTimeUNIX), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RUBY", testTimeRUBY}, date("D d M y", "RUBY", testTimeRUBY), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "ANSIC", testTimeANSIC}, date("D d M y", "ANSIC", testTimeANSIC), testTime.Format("Mon 02 Jan 06") },
	}

	testRunTests("date", tests, tester)
}

func TestDatetime(tester *testing.T) {
	currentTime			:= time.Now().In(dateLocalTimezone)
	testTimeISO8601Z	:= "2019-04-23T11:30:21+01:00"
	testTimeISO8601		:= "2019-04-23T11:30:21+01:00"
	testTimeRFC822Z		:= "Tue, 23 Apr 19 11:30:21 +01:00"
	testTimeRFC822		:= "Tue, 23 Apr 19 11:30:21 CET"
	testTimeRFC850		:= "Tuesday, 23-Apr-19 11:30:21 CET"
	testTimeRFC1036		:= "Tue, 23 Apr 19 11:30:21 +01:00"
	testTimeRFC1123Z	:= "Tue, 23 Apr 2019 11:30:21 +01:00"
	testTimeRFC1123		:= "Tue, 23 Apr 2019 11:30:21 CET"
	testTimeRFC2822		:= "Tue, 23 Apr 2019 11:30:21 +01:00"
	testTimeRFC3339		:= "2019-04-23T11:30:21+01:00"
	testTimeWC3			:= "2019-04-23T11:30:21+01:00"
	testTimeATOM		:= "2019-04-23T11:30:21+01:00"
	testTimeCOOKIE		:= "Tuesday, 23-Apr-2019 11:30:21 CET"
	testTimeRSS			:= "Tue, 23 Apr 2019 11:30:21 +01:00"
	testTimeMYSQL		:= "2019-04-23 10:30:21"
	testTimeUNIX		:= "Tue Apr 23 11:30:21 CET 2019"
	testTimeRUBY		:= "Tue Apr 23 11:30:21 +0100 2019"
	testTimeANSIC		:= "Tue Apr 23 10:30:21 2019"

	testTime, _ := time.Parse(time.RFC3339, testTimeRFC3339)
	testTime = testTime.In(dateLocalTimezone)

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, datetime(), currentTime.Format("02/01/2006 15:04") },

		{ []any{testTime}, datetime(testTime), testTime.Format("02/01/2006 15:04") },
		{ []any{1556015421}, datetime(1556015421), testTime.Format("02/01/2006 15:04") },

		{ []any{"02-01-2006 15:04"}, datetime("02-01-2006 15:04"), currentTime.Format("02-01-2006 15:04") },
		{ []any{"d-m-Y H:i"}, datetime("d-m-Y H:i"), currentTime.Format("02-01-2006 15:04") },
		{ []any{"%d-%m-%Y %H:%M"}, datetime("%d-%m-%Y %H:%M"), currentTime.Format("02-01-2006 15:04") },
		{ []any{"Mon 02 Jan 06 15:04"}, datetime("Mon 02 Jan 06 15:04"), currentTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i"}, datetime("D d M y H:i"), currentTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M"}, datetime("%a %d %b %y %H:%M"), currentTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", 1556015421}, datetime("Mon 02 Jan 06 15:04", 1556015421), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", testTime}, datetime("Mon 02 Jan 06 15:04", testTime), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", testTime}, datetime("D d M y H:i", testTime), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", testTime}, datetime("%a %d %b %y %H:%M", testTime), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", testTimeRFC3339}, datetime("Mon 02 Jan 06 15:04", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", testTimeRFC3339}, datetime("D d M y H:i", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", testTimeRFC3339}, datetime("%a %d %b %y %H:%M", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, datetime("Mon 02 Jan 06 15:04", "2006-01-02T15:04:05Z07:00", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, datetime("D d M y H:i", "2006-01-02T15:04:05Z07:00", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, datetime("%a %d %b %y %H:%M", "2006-01-02T15:04:05Z07:00", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"D d M y H:i", "ISO8601Z", testTimeISO8601Z}, datetime("D d M y H:i", "ISO8601Z", testTimeISO8601Z), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "ISO8601", testTimeISO8601}, datetime("D d M y H:i", "ISO8601", testTimeISO8601), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC822Z", testTimeRFC822Z}, datetime("D d M y H:i", "RFC822Z", testTimeRFC822Z), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC822", testTimeRFC822}, datetime("D d M y H:i", "RFC822", testTimeRFC822), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC850", testTimeRFC850}, datetime("D d M y H:i", "RFC850", testTimeRFC850), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1036", testTimeRFC1036}, datetime("D d M y H:i", "RFC1036", testTimeRFC1036), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1123Z", testTimeRFC1123Z}, datetime("D d M y H:i", "RFC1123Z", testTimeRFC1123Z), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1123", testTimeRFC1123}, datetime("D d M y H:i", "RFC1123", testTimeRFC1123), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC2822", testTimeRFC2822}, datetime("D d M y H:i", "RFC2822", testTimeRFC2822), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC3339", testTimeRFC3339}, datetime("D d M y H:i", "RFC3339", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"D d M y H:i", "ATOM", testTimeATOM}, datetime("D d M y H:i", "ATOM", testTimeATOM), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "W3C", testTimeWC3}, datetime("D d M y H:i", "W3C", testTimeWC3), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "COOKIE", testTimeCOOKIE}, datetime("D d M y H:i", "COOKIE", testTimeCOOKIE), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RSS", testTimeRSS}, datetime("D d M y H:i", "RSS", testTimeRSS), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "MYSQL", testTimeMYSQL}, datetime("D d M y H:i", "MYSQL", testTimeMYSQL), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "UNIX", testTimeUNIX}, datetime("D d M y H:i", "UNIX", testTimeUNIX), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RUBY", testTimeRUBY}, datetime("D d M y H:i", "RUBY", testTimeRUBY), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "ANSIC", testTimeANSIC}, datetime("D d M y H:i", "ANSIC", testTimeANSIC), testTime.Format("Mon 02 Jan 06 15:04") },
	}

	testRunTests("datetime", tests, tester)
}

func TestDefault(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ false, true, true },
		{ true, false, true },
		{ false, false, false },
		{ true, "", true },
		{ true, "test", "test" },
		{ 0, 0, 0 },
		{ 20, 0, 20 },
		{ 20, 10, 10 },
		{ -20, -10, -10},
		{ 3.5, 0, 3.5},
		{ 3.5, 6.4, 6.4},
		{ "default val", "test val", "test val"},
		{ []string{"default val"}, []string{}, []string{"default val"}},
		{ []string{"default val"}, []string{"test val"}, []string{"test val"}},
		{ map[string]string{"test": "default val"}, map[string]string{}, map[string]string{"test": "default val"} },
		{ map[string]string{"test": "default val"}, map[string]string{"test": "test val"}, map[string]string{"test": "test val"} },
		{ struct{ string string }{"default val"}, struct{ string string }{}, struct{ string string }{"default val"} },
		{ struct{ string string }{"default val"}, struct{ string string }{"test val"}, struct{ string string }{"test val"} },
	}

	testRunArgTests(defaultVal, tests, tester)
}

func TestDivide(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ "string", 10, 10 },
		{ true, 10, 10 },
		{ 0, 10, 10 },
		{ 2, 10, 5 },
		{ -2, 10, -5 },
		{ 2, -10, -5 },
		{ 2, 10.3, 5.15 },
		{ 3.3, -104.3, -31.606060606060606 },
		{ 5, "test", "test" },
		{ 5, []string{"test"}, []string{"test"} },
		{ 5, []int{10, 20}, []int{2, 4} },
		{ 5.1, []int{10, 20}, []int{2, 4} },
		{ 3.2, []float64{10, 20}, []float64{3.125, 6.25} },
		{ 5, map[string]int{"val1": 10, "val2": 20}, map[string]int{"val1": 2, "val2": 4} },
		{ 5, struct{ Num1, Num2 int }{10, 20}, struct{ Num1, Num2 int }{2, 4} },
		{ 5, struct{ num1, num2 int }{10, 20}, struct{ num1, num2 int }{0, 0} },
	}

	testRunArgTests(divide, tests, tester)
}

func TestDivisibleBy(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ 0, 10, false },
		{ 2, 10, true },
		{ 10, 2, false },
		{ -2, 10, true },
		{ 2, -10, true },
		{ 3, 10, false },
		{ 2, 10.3, false },
		{ 2.5, 10.0, true },
		{ 0.8, 2.4, true },
		{ 3.3, -104.3, false },
		{ 1000000000, 9999999999, false },
		{ 5, []int{10, 20}, false },
		{ 5.1, []int{10, 20}, false },
		{ 3.2, []float64{10, 20}, false },
		{ 5, map[string]int{"val1": 10, "val2": 20}, false },
		{ 5, struct{ Num1, Num2 int }{10, 20}, false },
		{ 5, struct{ num1, num2 int }{10, 20}, false },
	}

	testRunArgTests(divisibleBy, tests, tester)
}

func TestDl(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ false, "false" },
		{ true, "true" },
		{ 0, "0" },
		{ 1, "1" },
		{ -2, "-2" },
		{ 3.5, "3.5" },
		{ input1: -4.6, expected: "-4.6" },
		{ "test string", "test string" },
		{ []int{1, 2}, "<dl><dd>1</dd><dd>2</dd></dl>" },
		{ [2]int{1, 2}, "<dl><dd>1</dd><dd>2</dd></dl>" },
		{ []string{"value1", "value2"}, "<dl><dd>value1</dd><dd>value2</dd></dl>" },
		{ [2]string{"value1", "value2"}, "<dl><dd>value1</dd><dd>value2</dd></dl>" },
		{ map[int]string{1: "value1", 2: "value2"}, "<dl><dt>1</dt><dd>value1</dd><dt>2</dt><dd>value2</dd></dl>" },
		{ map[string]string{"title1": "value1", "title2": "value2"}, "<dl><dt>title1</dt><dd>value1</dd><dt>title2</dt><dd>value2</dd></dl>" },
		{ [][]string{{"subvalue1", "subvalue2"}, {"subvalue3", "subvalue4"}}, "<dl><dd><dl><dd>subvalue1</dd><dd>subvalue2</dd></dl></dd><dd><dl><dd>subvalue3</dd><dd>subvalue4</dd></dl></dd></dl>" },
		{ map[string]map[string]string{"title1": {"nested1": "subvalue1", "sub2": "subvalue2"}}, "<dl><dt>title1</dt><dd><dl><dt>nested1</dt><dd>subvalue1</dd><dt>sub2</dt><dd>subvalue2</dd></dl></dd></dl>" },
	}

	testRunArgTests(dl, tests, tester)
}

func TestEqual(tester *testing.T) {
	tests2 := []struct { input1, input2, expected any } {
		{ true, true, true },
		{ false, false, true },
		{ true, false, false },
		{ 1, 1, true },
		{ uint8(1), int64(1), true },
		{ 1, 1.0, true },
		{ 1, 1.000000000001, true },
		{ "1", 1, false },
		{ "hello world", "hello world", true },
		{ "Hello world", "hello world", false },
		{ "hello world", []string{"hello world"}, false },
		{ []string{"hello world"}, []string{"hello world"}, true },
		{ map[string]string{"test": "hello world"}, map[string]string{"test": "hello world"}, true },
		{ map[string]string{"test1": "hello world"}, map[string]string{"test": "hello world"}, false },
	}

	tests3 := []struct { input1, input2, input3, expected any } {
		{ true, true, true, true },
		{ false, false, true, false },
		{ true, false, false, false },
		{ 1, 1, 1, true },
		{ uint8(1), int64(1), float32(1), true },
		{ 1, 1.0, int16(1), true },
		{ 1, 1.000000000001, 0.999999999999, true },
		{ "1", 1, 1, false },
		{ "hello world", "hello world", "hello world", true },
		{ "hello world", "hello world", []string{"hello world"}, false },
		{ []string{"hello world"}, []string{"hello world"}, []string{"hello world"}, true },
		{ map[string]string{"test": "hello world"}, map[string]string{"test": "hello world"}, map[string]string{"test": "hello world"}, true },
		{ map[string]string{"test1": "hello world"}, map[string]string{"test1": "hello world"}, map[string]string{"test": "hello world"}, false },
	}

	passed, failed := 0, 0
	for _, test := range tests2 {
		if testCallVarArgs(tester, equal, []any{test.input1, test.input2}, test.expected) {
			passed++
		} else { failed++ }
	}
	for _, test := range tests3 {
		if testCallVarArgs(tester, equal, []any{test.input1, test.input2, test.input3}, test.expected) {
			passed++
		} else { failed++ }
	}
	
	testFormatPassFail("equal", passed, failed)
}

func TestFirst(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ nil, nil },
		{ "", nil },
		{ 0, nil },
		{ 10, nil },
		{ -10, nil },
		{ 10.34, nil },
		{ -10.34, nil },
		{ true, nil },
		{ "hello world", "hello" },
		{ "hello", "hello" },
		{ " hello ", "hello" },
		{ " 123 hello ", "123" },
		{ []int{}, nil },
		{ []int{1}, 1 },
		{ []int{1, 2, 3}, 1 },
		{ []string{}, nil },
		{ []string{"hello", "world", "how", "are", "you?"}, "hello" },
		{ [0]string{}, nil },
		{ [5]string{"hello", "world", "how", "are", "you?"}, "hello" },
		{ [][]string{}, nil },
		{ [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, []string{"hello", "world"} },
		{ map[int]string{1: "test", 2: "test"}, nil },
		{ struct{ Str1, Str2 string } {"first", "last"}, "first" },
		{ struct{ str1, str2 string } {"first", "last"}, "first" },
	}

	testRunArgTests(first, tests, tester)
}

func TestFirstOf(tester *testing.T) {
	tests1 := []struct { input1, expected any } {
		{ nil, nil },
		{ "", nil },
		{ 0, nil },
		{ 0.0, nil },
		{ 10, 10 },
		{ -10, -10 },
		{ 10.34, 10.34 },
		{ -10.34, -10.34 },
		{ true, true },
		{ "hello world", "hello world" },
		{ []int{}, nil },
		{ []int{1}, []int{1} },
		{ []int{1, 2, 3}, []int{1, 2, 3} },
		{ []string{}, nil },
		{ []string{"hello", "world", "how", "are", "you?"}, []string{"hello", "world", "how", "are", "you?"} },
		{ [0]string{}, nil },
		{ [5]string{"hello", "world", "how", "are", "you?"}, [5]string{"hello", "world", "how", "are", "you?"} },
		{ [][]string{}, nil },
		{ [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}} },
		{ map[int]string{}, nil },
		{ map[int]string{1: "test", 2: "test"}, map[int]string{1: "test", 2: "test"} },
		{ struct{ Str string }{}, nil },
		{ struct{ Str string }{ "test" }, struct{ Str string }{ "test" } },
	}

	tests2 := []struct { input1, input2, expected any } {
		{ nil, nil, nil },
		{ nil, 1, 1 },
		{ 2, 1, 2 },
	}

	tests3 := []struct { input1, input2, input3, expected any } {
		{ nil, nil, nil, nil },
		{ nil, nil, 1, 1 },
		{ nil, 2, 1, 2 },
		{ 3, 2, 1, 3 },
	}

	passed, failed := 0, 0
	for _, test := range tests1 {
		if testCallVarArgs(tester, firstOf, []any{test.input1}, test.expected) {
			passed++
		} else { failed++ }
	}
	for _, test := range tests2 {
		if testCallVarArgs(tester, firstOf, []any{test.input1, test.input2}, test.expected) {
			passed++
		} else { failed++ }
	}
	for _, test := range tests3 {
		if testCallVarArgs(tester, firstOf, []any{test.input1, test.input2, test.input3}, test.expected) {
			passed++
		} else { failed++ }
	}
	
	testFormatPassFail("firstof", passed, failed)
}

func TestFormatTime(tester *testing.T) {
	testTimeRFC3339 := "2019-04-23T11:30:21+01:00"
	testTime, _ := time.Parse(time.RFC3339, testTimeRFC3339)

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{testTime}, formattime("02/01/2006 15:04", testTime), testTime.Format("02/01/2006 15:04") },
		{ []any{testTime}, formattime("d/m/Y H:i", testTime), testTime.Format("02/01/2006 15:04") },
		{ []any{testTime}, formattime("%d/%m/%Y %H:%M", testTime), testTime.Format("02/01/2006 15:04") },
	}

	testRunTests("formattime", tests, tester)
}

func TestGreaterThan(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ false, true, false},
		{ false, 5, false},
		{ 5, true, false},
		{ 5, 5, false},
		{ 10, 5, false},
		{ 5, 10, true},
		{ 5.0, 5, false},
		{ 10.0, 5, false},
		{ 5, 10.0, true},
		{ 5.1, 5.1, false},
		{ 10.1, 5.1, false},
		{ 5.1, 10.1, true},
		{ "test", 10, false},
		{ 10, "test", false},
		{ "test1", "test2", false},
		{ []int{10}, 5, false},
		{ []int{5}, 10, false},
		{ 5, []int{10}, false},
		{ 10, []int{5}, false},
	}

	testRunArgTests(greaterThan, tests, tester)
}

func TestGreaterThanEqual(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ false, true, false},
		{ false, 5, false},
		{ 5, true, false},
		{ 5, 5, true},
		{ 10, 5, false},
		{ 5, 10, true},
		{ 5.0, 5, true},
		{ 10.0, 5, false},
		{ 5, 10.0, true},
		{ 5.1, 5.1, true},
		{ 10.1, 5.1, false},
		{ 5.1, 10.1, true},
		{ "test", 10, false},
		{ 10, "test", false},
		{ "test1", "test2", false},
		{ []int{10}, 5, false},
		{ []int{5}, 10, false},
		{ 5, []int{10}, false},
		{ 10, []int{5}, false},
	}

	testRunArgTests(greaterThanEqual, tests, tester)
}

func TestHtmlDecode(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ false, false },
		{ true, true },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ 10.1, 10.1 },
		{ -10.1, -10.1 },
		{ "string without html", "string without html" },
		{ "string <strong>with</strong> html", "string <strong>with</strong> html" },
		{ "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#x27; &amp;amp; other &#34;nasty&#x22; stuff", "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff" },
		{ []string{"string without html"}, []string{"string without html"} },
		{ []string{"string <strong>with</strong> html"}, []string{"string <strong>with</strong> html"} },
		{ []string{"safe string", "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#x27; &amp;amp; other &#34;nasty&#x22; stuff"}, []string{"safe string", "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff"} },
		{ map[int]string{1: "string without html"}, map[int]string{1: "string without html"} },
		{ map[int]string{1: "string <strong>with</strong> html"}, map[int]string{1: "string <strong>with</strong> html"} },
		{ map[int]string{1: "safe string", 2: "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#x27; &amp;amp; other &#34;nasty&#x22; stuff"}, map[int]string{1: "safe string", 2: "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff"} },
		{ struct{ String1, String2 string }{"string without html", "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#x27; &amp;amp; other &#34;nasty&#x22; stuff"}, struct{ String1, String2 string }{"string without html", "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff"} },
		{ struct{ string1, string2 string }{"string without html", "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#x27; &amp;amp; other &#34;nasty&#x22; stuff"}, struct{ string1, string2 string }{"", ""} },
	}

	testRunArgTests(htmlDecode, tests, tester)
}

func TestHtmlEncode(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ false, false },
		{ true, true },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ 10.1, 10.1 },
		{ -10.1, -10.1 },
		{ "string without html", "string without html" },
		{ "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff", "&#34;string&#34; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#39; &amp;amp; other &#34;nasty&#34; stuff" },
		{ []string{"string without html"}, []string{"string without html"} },
		{ []string{"safe string", "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff"}, []string{"safe string", "&#34;string&#34; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#39; &amp;amp; other &#34;nasty&#34; stuff"} },
		{ map[int]string{1: "string without html"}, map[int]string{1: "string without html"} },
		{ map[int]string{1: "safe string", 2: "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff"}, map[int]string{1: "safe string", 2: "&#34;string&#34; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#39; &amp;amp; other &#34;nasty&#34; stuff"} },
		{ struct{ String1, String2 string }{"string without html", "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff"}, struct{ String1, String2 string }{"string without html", "&#34;string&#34; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#39; &amp;amp; other &#34;nasty&#34; stuff"} },
		{ struct{ string1, string2 string }{"string without html", "&#34;string&#34; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#39; &amp;amp; other &#34;nasty&#34; stuff"}, struct{ string1, string2 string }{"", ""} },
	}

	testRunArgTests(htmlEncode, tests, tester)
}

func TestJoin(tester *testing.T) {
	tests := []struct{ inputs []any; result any; expected any } {
		{ []any{", ", ""}, join(", ", reflect.ValueOf("")), "" },
		{ []any{", ", nil}, join(", ", reflect.ValueOf(nil)), "" },
		{ []any{", ", 0}, join(", ", reflect.ValueOf(0)), "0" },
		{ []any{", ", -1}, join(", ", reflect.ValueOf(-1)), "-1" },
		{ []any{", ", 1}, join(", ", reflect.ValueOf(1)), "1" },
		{ []any{", ", 0.0}, join(", ", reflect.ValueOf(0.0)), "0" },
		{ []any{", ", 1.0}, join(", ", reflect.ValueOf(1.0)), "1" },
		{ []any{", ", 0.1}, join(", ", reflect.ValueOf(0.1)), "0.1" },
		{ []any{", ", 1.1}, join(", ", reflect.ValueOf(1.1)), "1.1" },
		{ []any{", ", true}, join(", ", reflect.ValueOf(true)), "true" },
		{ []any{", ", false}, join(", ", reflect.ValueOf(false)), "false" },
		{ []any{", ", "string value"}, join(", ", reflect.ValueOf("string value")), "string value" },
		{ []any{", ", []string{"string", "value"}}, join(", ", reflect.ValueOf([]string{"string", "value"})), "string, value" },
		{ []any{", ", []int{1, 2}}, join(", ", reflect.ValueOf([]int{1, 2})), "1, 2" },
		{ []any{", ", []float64{0.0, 1.1, 2.2}}, join(", ", reflect.ValueOf([]float64{0.0, 1.1, 2.2})), "0, 1.1, 2.2" },
		{ []any{", ", []bool{true, false, true}}, join(", ", reflect.ValueOf([]bool{true, false, true})), "true, false, true" },
		{ []any{", ", map[int]string{1: "first", 2: "second"}}, join(", ", reflect.ValueOf(map[int]string{1: "first", 2: "second"})), "first, second" },
		{ []any{", ", map[int][]string{1: {"first", "second"}, 2: {"third"}}}, join(", ", reflect.ValueOf(map[int][]string{1: {"first", "second"}, 2: {"third"}})), "first, second, third" },
		{ []any{", ", struct{ first string; second int; third float64 } {"first", 1, 1.1}}, join(", ", reflect.ValueOf(struct{ first string; second int; third float64 } {"first", 1, 1.1})), "first, 1, 1.1" },
	}

	testRunTests("join", tests, tester)
}

func TestJsonDecode(tester *testing.T) {
	tests := []struct{ inputs []any; result any; expected any } {
		{ []any{""}, jsonDecode(""), nil },
		{ []any{"null"}, jsonDecode("null"), nil },
		{ []any{"{}"}, jsonDecode("{}"), map[string]any{} },
		{ []any{"[]"}, jsonDecode("{}"), map[string]any{} },
		{ []any{"1"}, jsonDecode("1"), 1.0 },
		{ []any{"1"}, jsonDecode("1"), 1.0 },
		{ []any{"-1.5"}, jsonDecode("-1.5"), -1.5 },
		{ []any{"true"}, jsonDecode("true"), true },
		{ []any{"false"}, jsonDecode("false"), false },
		{ []any{"string"}, jsonDecode("string"), nil },
		{ []any{`"string"`}, jsonDecode(`"string"`), "string" },
		{ []any{`["string","value"]`}, jsonDecode(`["string","value"]`), []any{"string", "value"} },
		{ []any{"[1,2]"}, jsonDecode("[1,2]"), []any{1.0, 2.0} },
		{ []any{"[0,1.1,2.2]"}, jsonDecode("[0,1.1,2.2]"), []any{0.0, 1.1, 2.2} },
		{ []any{"[true,false,true]"}, jsonDecode("[true,false,true]"), []any{true, false, true} },
		{ []any{`{"1":"first","2":"second"}`}, jsonDecode(`{"1":"first","2":"second"}`), map[string]any{"1":"first", "2":"second"} },
		{ []any{`{"1":["first","second"],"2":["third"]}`}, jsonDecode(`{"1":["first","second"],"2":["third"]}`), map[string]any{"1":[]any{"first", "second"}, "2":[]any{"third"}} },
		{ []any{`{"First":"first","Second":1,"Third":1.1}`}, jsonDecode(`{"First":"first","Second":1,"Third":1.1}`), map[string]any{"First":"first", "Second":1.0, "Third":1.1} },
	}

	testRunTests("jsonDecode", tests, tester)
}

func TestJsonEncode(tester *testing.T) {
	tests := []struct{ inputs []any; result any; expected any } {
		{ []any{""}, jsonEncode(""), `""` },
		{ []any{nil}, jsonEncode(nil), "null" },
		{ []any{0}, jsonEncode(0), "0" },
		{ []any{-1}, jsonEncode(-1), "-1" },
		{ []any{1}, jsonEncode(1), "1" },
		{ []any{0.0}, jsonEncode(0.0), "0" },
		{ []any{1.0}, jsonEncode(1.0), "1" },
		{ []any{0.1}, jsonEncode(0.1), "0.1" },
		{ []any{1.1}, jsonEncode(1.1), "1.1" },
		{ []any{true}, jsonEncode(true), "true" },
		{ []any{false}, jsonEncode(false), "false" },
		{ []any{"string value"}, jsonEncode("string value"), `"string value"` },
		{ []any{[]string{"string", "value"}}, jsonEncode([]string{"string", "value"}), `["string","value"]` },
		{ []any{[]int{1, 2}}, jsonEncode([]int{1, 2}), "[1,2]" },
		{ []any{[]float64{0.0, 1.1, 2.2}}, jsonEncode([]float64{0.0, 1.1, 2.2}), "[0,1.1,2.2]" },
		{ []any{[]bool{true, false, true}}, jsonEncode([]bool{true, false, true}), "[true,false,true]" },
		{ []any{map[int]string{1: "first", 2: "second"}}, jsonEncode(map[int]string{1: "first", 2: "second"}), `{"1":"first","2":"second"}` },
		{ []any{map[int][]string{1: {"first", "second"}, 2: {"third"}}}, jsonEncode(map[int][]string{1: {"first", "second"}, 2: {"third"}}), `{"1":["first","second"],"2":["third"]}` },
		{ []any{struct{ first string; second int; third float64 } {"first", 1, 1.1}}, jsonEncode(struct{ first string; second int; third float64 } {"first", 1, 1.1}), "{}" },
		{ []any{struct{ First string; Second int; Third float64 } {"first", 1, 1.1}}, jsonEncode(struct{ First string; Second int; Third float64 } {"first", 1, 1.1}), `{"First":"first","Second":1,"Third":1.1}` },
	}

	testRunTests("jsonEncode", tests, tester)
}

func TestKey(tester *testing.T) {
	tests1 := []struct { input1, expected any } {
		{ nil, nil },
		{ "", "" },
		{ 0, 0 },
		{ 0.0, 0.0 },
		{ 10, 10 },
		{ -10, -10 },
		{ 10.34, 10.34 },
		{ -10.34, -10.34 },
		{ true, true },
		{ "hello world", "hello world" },
		{ []int{}, []int{} },
		{ []int{1}, []int{1} },
		{ []int{1, 2, 3}, []int{1, 2, 3} },
		{ []string{}, []string{} },
		{ []string{"hello", "world", "how", "are", "you?"}, []string{"hello", "world", "how", "are", "you?"} },
		{ [0]string{}, [0]string{} },
		{ [5]string{"hello", "world", "how", "are", "you?"}, [5]string{"hello", "world", "how", "are", "you?"} },
		{ [][]string{}, [][]string{} },
		{ [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}} },
		{ map[int]string{}, map[int]string{} },
		{ map[int]string{1: "test", 2: "test"}, map[int]string{1: "test", 2: "test"} },
		{ struct{value string} {"test"}, struct{value string} {"test"} },
		{ struct{Value string} {"test"}, struct{Value string} {"test"} },
	}

	tests2 := []struct { input1, input2, expected any } {
		{ nil, nil, nil },
		{ nil, 1, nil },
		{ 1, nil, nil },
		{ 2, "string", "r" },
		{ "r", "string", nil },
		{ 0, []int{}, nil },
		{ 0, []int{1}, 1 },
		{ 1, []int{1, 2, 3}, 2 },
		{ 1, []string{}, nil },
		{ 2, []string{"hello", "world", "how", "are", "you?"}, "how" },
		{ 0, [0]string{}, nil },
		{ 4, [5]string{"hello", "world", "how", "are", "you?"}, "you?" },
		{ 0, [][]int{}, nil },
		{ 0, [][]int{{1, 2}, {3, 4}, {5}}, []int{1, 2} },
		{ 2, [][]int{{1, 2}, {3, 4}, {5}}, []int{5} },
		{ 3, [][]int{{1, 2}, {3, 4}, {5}}, nil },
		{ 0, [][]string{}, nil },
		{ 0, [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, []string{"hello", "world"} },
		{ 3, [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, nil },
		{ 0, map[int]string{}, nil },
		{ 0, map[int]string{1: "test", 2: "test"}, nil },
		{ 1, map[int]string{1: "test", 2: "test"}, "test" },
		{ 0, struct{value string} {"test"}, "test" },
		{ 1, struct{value string} {"test"}, nil },
		{ "value", struct{value string} {"test"}, "test" },
		{ "unknown", struct{value string} {"test"}, nil },
		{ 0, struct{Value string} {"test"}, "test" },
		{ 1, struct{Value string} {"test"}, nil },
		{ "Value", struct{Value string} {"test"}, "test" },
		{ "value", struct{value []string} {[]string{"test"}}, []string{"test"} },
		{ "Value", struct{Value []string} {[]string{"test"}}, []string{"test"} },
		{ "value", struct{value []int} {[]int{1, 2}}, []int{1, 2} },
		{ "Value", struct{Value []int} {[]int{1, 2}}, []int{1, 2} },
		{ "value", struct{value []float64} {[]float64{1, 2}}, []float64{1, 2} },
		{ "value", struct{value []bool} {[]bool{true, false}}, []bool{true, false} },
		{ "value", struct{value [][]string} {[][]string{{"string", "slice"}, {"in", "private", "field"}}}, [][]string{{"string", "slice"}, {"in", "private", "field"}} },
		{ "value", struct{value map[int]string} {map[int]string{1:"string", 2:"map", 3:"in", 4:"private", 5:"field"}}, map[int]string{1:"string", 2:"map", 3:"in", 4:"private", 5:"field"} },
	}

	tests3 := []struct { input1, input2, input3, expected any } {
		{ 2, 1, "string", nil },
		{ 2, 0, "string", "r" },
		{ "r", 2, "string", nil },
		{ 0, 1, []int{}, nil },
		{ 0, 1, []int{1}, nil },
		{ 1, 1, []int{1, 2, 3}, nil },
		{ 1, 1, []string{}, nil },
		{ 2, 2, []string{"hello", "world", "how", "are", "you?"}, "w" },
		{ 0, 1, [0]string{}, nil },
		{ 4, 1, [5]string{"hello", "world", "how", "are", "you?"}, "o" },
		{ 0, 1, [][]int{}, nil },
		{ 0, 1, [][]int{{1, 2}, {3, 4}, {5}}, 2 },
		{ 2, 0, [][]int{{1, 2}, {3, 4}, {5}}, 5 },
		{ 2, 1, [][]int{{1, 2}, {3, 4}, {5}}, nil },
		{ 3, 0, [][]int{{1, 2}, {3, 4}, {5}}, nil },
		{ 0, 1, [][]string{}, nil },
		{ 0, 1, [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, "world" },
		{ 2, 0, [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, "you?" },
		{ 3, 0, [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, nil },
		{ 0, 1, map[int]string{}, nil },
		{ 1, 17, map[int]string{1: "test", 2: "test"}, nil },
		{ 1, 1, map[int]string{1: "test", 2: "test"}, "e" },
		{ 0, 1, struct{value string} {"test"}, "e" },
		{ 0, 1, struct{Value string} {"test"}, "e" },
		{ "value", 1, struct{value []bool} {[]bool{true, false}}, false },
		{ "value", 1, struct{value [][]string} {[][]string{{"string", "slice"}, {"in", "private", "field"}}}, []string{"in", "private", "field"} },
		{ "value", 2, struct{value map[int]string} {map[int]string{1:"string", 2:"map", 3:"in", 4:"private", 5:"field"}}, "map" },
	}

	passed, failed := 0, 0
	for _, test := range tests1 {
		if testCallVarArgs(tester, keyFn, []any{test.input1}, test.expected) {
			passed++
		} else { failed++ }
	}
	for _, test := range tests2 {
		if testCallVarArgs(tester, keyFn, []any{test.input1, test.input2}, test.expected) {
			passed++
		} else { failed++ }
	}
	for _, test := range tests3 {
		if testCallVarArgs(tester, keyFn, []any{test.input1, test.input2, test.input3}, test.expected) {
			passed++
		} else { failed++ }
	}
	
	testFormatPassFail("key", passed, failed)
}

func TestKind(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true,  "bool" },
		{ "anything", "string" },
		{ -10, "int" },
		{ int8(10), "int8" },
		{ uint32(10), "uint32" },
		{ 10.45, "float64" },
		{ []string{"hello world"}, "slice" },
		{ [1]string{"hello world"}, "array" },
		{ map[string]string{"test": "hello world"}, "map" },
		{ struct{Str string}{"hello world"}, "struct" },
		{ struct{str string}{"hello world"}, "struct" },
	}

	testRunArgTests(kind, tests, tester)
}

func TestLast(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ nil, nil },
		{ "", nil },
		{ 0, nil },
		{ 10, nil },
		{ -10, nil },
		{ 10.34, nil },
		{ -10.34, nil },
		{ true, nil },
		{ "hello world", "world" },
		{ "hello", "hello" },
		{ " hello ", "hello" },
		{ " 123 hello ", "hello" },
		{ []int{}, nil },
		{ []int{1}, 1 },
		{ []int{1, 2, 3}, 3 },
		{ []string{}, nil },
		{ []string{"hello", "world", "how", "are", "you?"}, "you?" },
		{ [0]string{}, nil },
		{ [5]string{"hello", "world", "how", "are", "you?"}, "you?" },
		{ [][]string{}, nil },
		{ [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, []string{"you?"} },
		{ map[int]string{1: "test", 2: "test"}, nil },
		{ struct{ Str1, Str2 string } {"first", "last"}, "last" },
		{ struct{ str1, str2 string } {"first", "last"}, "last" },
	}

	testRunArgTests(last, tests, tester)
}

func TestLength(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, 1 },
		{ "anything", 8 },
		{ -10, 3 },
		{ int8(10), 2 },
		{ uint(10), 2 },
		{ uint32(10), 2 },
		{ 10.45, 5 },
		{ 1.66666666667, 13 },
		{ []string{"hello world"}, 1 },
		{ [2]string{"hello", "world"}, 2 },
		{ map[string]string{"test": "hello world"}, 1 },
		{ struct{Str string}{"hello world"}, 1 },
		{ struct{str string}{"hello world"}, 1 },
	}

	testRunArgTests(length, tests, tester)
}

func TestLessThan(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ false, true, false},
		{ false, 5, false},
		{ 5, true, false},
		{ 5, 5, false},
		{ 10, 5, true},
		{ 5, 10, false},
		{ 5.0, 5, false},
		{ 10.0, 5, true},
		{ 5, 10.0, false},
		{ 5.1, 5.1, false},
		{ 10.1, 5.1, true},
		{ 5.1, 10.1, false},
		{ "test", 10, false},
		{ 10, "test", false},
		{ "test1", "test2", false},
		{ []int{10}, 5, false},
		{ []int{5}, 10, false},
		{ 5, []int{10}, false},
		{ 10, []int{5}, false},
	}

	testRunArgTests(lessThan, tests, tester)
}

func TestLessThanEqual(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ false, true, false},
		{ false, 5, false},
		{ 5, true, false},
		{ 5, 5, true},
		{ 10, 5, true},
		{ 5, 10, false},
		{ 5.0, 5, true},
		{ 10.0, 5, true},
		{ 5, 10.0, false},
		{ 5.1, 5.1, true},
		{ 10.1, 5.1, true},
		{ 5.1, 10.1, false},
		{ "test", 10, false},
		{ 10, "test", false},
		{ "test1", "test2", false},
		{ []int{10}, 5, false},
		{ []int{5}, 10, false},
		{ 5, []int{10}, false},
		{ 10, []int{5}, false},
	}

	testRunArgTests(lessThanEqual, tests, tester)
}

func TestLocaltime(tester *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2019-04-23T11:30:21+01:00")

	utc, _ := time.LoadLocation("UTC")
	lon, _ := time.LoadLocation("Europe/London")
	est, _ := time.LoadLocation("EST")

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{testTime}, localtime("UTC", testTime), testTime.In(utc) },
		{ []any{testTime}, localtime("Europe/London", testTime), testTime.In(lon) },
		{ []any{testTime}, localtime("EST", testTime), testTime.In(est) },
	}

	testRunTests("localtime", tests, tester)
}

func TestLower(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, true },
		{ false, false },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ "hello World", "hello world" },
		{ "123 hello World", "123 hello world" },
		{ []string{"HEllo wORld"}, []string{"hello world"} },
		{ []string{"123 hello world"}, []string{"123 hello world"} },
		{ map[string]string{"test": "HELLO world"}, map[string]string{"test": "hello world"} },
		{ map[string]string{"test": "123 HELLO world"}, map[string]string{"test": "123 hello world"} },
		{ struct{Str string}{"Test"}, struct{Str string}{"test"} },
		{ struct{str string}{"test"}, struct{str string}{""} },
	}

	testRunArgTests(lower, tests, tester)
}

func TestLtrim(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, true, nil },
		{ "anything", false, false },
		{ " ", 0, 0 },
		{ " ", 10, 10 },
		{ " ", -10, -10 },
		{ " ", "hello world ", "hello world " },
		{ " 123", "123 hello world 123 ", "hello world 123 " },
		{ " ", []string{"   hello world "}, []string{"hello world "} },
		{ " 123", []string{"123 hello world 123"}, []string{"hello world 123"} },
		{ " ", map[string]string{"test": " hello world "}, map[string]string{"test": "hello world "} },
		{ " 123", map[string]string{"test": " 123 hello world 123 "}, map[string]string{"test": "hello world 123 "} },
		{ " ", struct{Str string}{" Test "}, struct{Str string}{"Test "} },
		{ " ", struct{str string}{" test "}, struct{str string}{""} },
	}

	testRunArgTests(ltrim, tests, tester)
}

func TestMktime(tester *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2019-04-23T11:30:21+01:00")
	testTime = testTime.In(dateLocalTimezone)

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, mktime(), now() },
		{ []any{"invalid"}, mktime("invalid"), now() },
		{ []any{"2019-04-23T11:30:21+01:00"}, mktime("2019-04-23T11:30:21+01:00"), testTime },
		{ []any{"ATOM", "2019-04-23T11:30:21+01:00"}, mktime("ATOM", "2019-04-23T11:30:21+01:00"), testTime },
	}

	testRunTests("mktime", tests, tester)
}

func TestMultiply(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ "string", 10, 10 },
		{ true, 10, 10 },
		{ 0, 10, 0 },
		{ 2, 10, 20 },
		{ -2, 10, -20 },
		{ 2, -10, -20 },
		{ 2, 10.3, 20.6 },
		{ 3.3, -104.3, -344.19 },
		{ 5, "test", "test" },
		{ 5, []string{"test"}, []string{"test"} },
		{ 5, []int{10, 20}, []int{50, 100} },
		{ 5.1, []int{10, 20}, []int{51, 102} },
		{ 5.15, []float64{10, 20}, []float64{51.5, 103} },
		{ 5, map[string]int{"val1": 10, "val2": 20}, map[string]int{"val1": 50, "val2": 100} },
		{ 5, struct{ Num1, Num2 int }{10, 20}, struct{ Num1, Num2 int }{50, 100} },
		{ 5, struct{ num1, num2 int }{10, 20}, struct{ num1, num2 int }{0, 0} },
	}

	testRunArgTests(multiply, tests, tester)
}

func TestNl2br(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, true },
		{ false, false },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ "hello world", "hello world" },
		{ "hello\nworld", "hello<br>world" },
		{ []string{"hello world"}, []string{"hello world"} },
		{ []string{"hello\nworld"}, []string{"hello<br>world"} },
		{ map[string]string{"test": "hello world"}, map[string]string{"test": "hello world"} },
		{ map[string]string{"test": "hello\nworld"}, map[string]string{"test": "hello<br>world"} },
		{ struct{Str string}{"hello world"}, struct{Str string}{"hello world"} },
		{ struct{Str string}{"hello\nworld"}, struct{Str string}{"hello<br>world"} },
		{ struct{str string}{"hello world"}, struct{str string}{""} },
	}

	testRunArgTests(nl2br, tests, tester)
}

func TestNow(tester *testing.T) {
	testTime := time.Now().In(dateLocalTimezone)

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, now(), testTime },
	}

	testRunTests("now", tests, tester)
}

func TestOl(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ false, "false" },
		{ true, "true" },
		{ 0, "0" },
		{ 1, "1" },
		{ -2, "-2" },
		{ 3.5, "3.5" },
		{ input1: -4.6, expected: "-4.6" },
		{ "test string", "test string" },
		{ []int{1, 2}, "<ol><li>1</li><li>2</li></ol>" },
		{ [2]int{1, 2}, "<ol><li>1</li><li>2</li></ol>" },
		{ []string{"value1", "value2"}, "<ol><li>value1</li><li>value2</li></ol>" },
		{ [2]string{"value1", "value2"}, "<ol><li>value1</li><li>value2</li></ol>" },
		{ map[int]string{1: "value1", 2: "value2"}, "<ol><li>value1</li><li>value2</li></ol>" },
		{ map[string]string{"title1": "value1", "title2": "value2"}, "<ol><li>value1</li><li>value2</li></ol>" },
		{ [][]string{{"subvalue1", "subvalue2"}, {"subvalue3", "subvalue4"}}, "<ol><li><ol><li>subvalue1</li><li>subvalue2</li></ol></li><li><ol><li>subvalue3</li><li>subvalue4</li></ol></li></ol>" },
		{ map[string]map[string]string{"title1": {"nested1": "subvalue1", "sub2": "subvalue2"}}, "<ol><li><ol><li>subvalue1</li><li>subvalue2</li></ol></li></ol>" },
	}

	testRunArgTests(ol, tests, tester)
}

func TestOrdinal(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{0}, ordinal(reflect.ValueOf(0)), "0th" },
		{ []any{1}, ordinal(reflect.ValueOf(1)), "1st" },
		{ []any{2}, ordinal(reflect.ValueOf(2)), "2nd" },
		{ []any{3}, ordinal(reflect.ValueOf(3)), "3rd" },
		{ []any{4}, ordinal(reflect.ValueOf(4)), "4th" },
		{ []any{5}, ordinal(reflect.ValueOf(5)), "5th" },
		{ []any{10}, ordinal(reflect.ValueOf(10)), "10th" },
		{ []any{11}, ordinal(reflect.ValueOf(11)), "11th" },
		{ []any{12}, ordinal(reflect.ValueOf(12)), "12th" },
		{ []any{13}, ordinal(reflect.ValueOf(13)), "13th" },
		{ []any{20}, ordinal(reflect.ValueOf(20)), "20th" },
		{ []any{21}, ordinal(reflect.ValueOf(21)), "21st" },
		{ []any{22}, ordinal(reflect.ValueOf(22)), "22nd" },
		{ []any{23}, ordinal(reflect.ValueOf(23)), "23rd" },
		{ []any{101}, ordinal(reflect.ValueOf(101)), "101st" },
		{ []any{102}, ordinal(reflect.ValueOf(102)), "102nd" },
		{ []any{103}, ordinal(reflect.ValueOf(103)), "103rd" },
		{ []any{111}, ordinal(reflect.ValueOf(111)), "111th" },
		{ []any{112}, ordinal(reflect.ValueOf(112)), "112th" },
		{ []any{113}, ordinal(reflect.ValueOf(113)), "113th" },
		{ []any{121}, ordinal(reflect.ValueOf(121)), "121st" },
		{ []any{122}, ordinal(reflect.ValueOf(122)), "122nd" },
		{ []any{123}, ordinal(reflect.ValueOf(123)), "123rd" },
		{ []any{1001}, ordinal(reflect.ValueOf(1001)), "1001st" },
		{ []any{1002}, ordinal(reflect.ValueOf(1002)), "1002nd" },
		{ []any{1003}, ordinal(reflect.ValueOf(1003)), "1003rd" },
		{ []any{1011}, ordinal(reflect.ValueOf(1011)), "1011th" },
		{ []any{1012}, ordinal(reflect.ValueOf(1012)), "1012th" },
		{ []any{1013}, ordinal(reflect.ValueOf(1013)), "1013th" },
		{ []any{1021}, ordinal(reflect.ValueOf(1021)), "1021st" },
		{ []any{1022}, ordinal(reflect.ValueOf(1022)), "1022nd" },
		{ []any{1023}, ordinal(reflect.ValueOf(1023)), "1023rd" },
	}

	testRunTests("ordinal", tests, tester)
}

func TestParagraph(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, true },
		{ false, false },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ "hello world", "<p>hello world</p>" },
		{ "hello\nworld", "<p>hello<br>world</p>" },
		{ "hello\rworld", "<p>hello<br>world</p>" },
		{ "hello\r\nworld", "<p>hello<br>world</p>" },
		{ "hello\n\nworld", "<p>hello</p><p>world</p>" },
		{ "hello\n\n\nworld", "<p>hello</p><p>world</p>" },
		{ "hello \n \n \n world", "<p>hello</p><p>world</p>" },
		{ "hello\r\n \r\nworld", "<p>hello</p><p>world</p>" },
		{ []string{"hello world"}, []string{"<p>hello world</p>"} },
		{ []string{"hello\nworld"}, []string{"<p>hello<br>world</p>"} },
		{ []string{"hello\n\nworld"}, []string{"<p>hello</p><p>world</p>"} },
		{ map[string]string{"test": "hello world"}, map[string]string{"test": "<p>hello world</p>"} },
		{ map[string]string{"test": "hello\nworld"}, map[string]string{"test": "<p>hello<br>world</p>"} },
		{ map[string]string{"test": "hello\n\nworld"}, map[string]string{"test": "<p>hello</p><p>world</p>"} },
		{ struct{Str string}{"hello world"}, struct{Str string}{"<p>hello world</p>"} },
		{ struct{Str string}{"hello\nworld"}, struct{Str string}{"<p>hello<br>world</p>"} },
		{ struct{Str string}{"hello\n\nworld"}, struct{Str string}{"<p>hello</p><p>world</p>"} },
		{ struct{str string}{"hello world"}, struct{str string}{""} },
	}

	testRunArgTests(paragraph, tests, tester)
}

func TestPluralise(tester *testing.T) {
	tests := []struct { inputs []any; result any; expected any } {
		{ []any{0}, pluralise(0), "s" },
		{ []any{1}, pluralise(1), "" },
		{ []any{2}, pluralise(2), "s" },
		{ []any{"es", 0}, pluralise("es", 0), "es" },
		{ []any{"es", 1}, pluralise("es", 1), "" },
		{ []any{"es", 2}, pluralise("es", 2), "es" },
		{ []any{"y", "ies", 0}, pluralise("y", "ies", 0), "ies" },
		{ []any{"y", "ies", 1}, pluralise("y", "ies", 1), "y" },
		{ []any{"y", "ies", 2}, pluralise("y", "ies", 2), "ies" },
		{ []any{1.5}, pluralise(1.5), "" },
		{ []any{false}, pluralise(false), "" },
		{ []any{[]string{"test"}}, pluralise([]string{"test"}), "" },
		{ []any{map[int]string{1: "test"}}, pluralise(map[int]string{1: "test"}), "" },
		{ []any{struct{ Str string }{"test"}}, pluralise(struct{ Str string }{"test"}), "" },
	}

	testRunTests("pluralise", tests, tester)
}

func TestPrefix(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ "prefix", 10, 10 },
		{ true, 10, 10 },
		{ 0, 10, 10 },
		{ "prefix", "test", "prefixtest" },
		{ "prefix", []string{"test"}, []string{"prefixtest"} },
		{ "prefix", []string{"test", "strings"}, []string{"prefixtest", "prefixstrings"} },
		{ 5, []int{10, 20}, []int{10, 20} },
		{ "prefix", []int{10, 20}, []int{10, 20} },
		{ 5, map[int]string{1: "val1", 2: "val2"}, map[int]string{1: "val1", 2: "val2"} },
		{ "prefix", map[int]string{1: "val1", 2: "val2"}, map[int]string{1: "prefixval1", 2: "prefixval2"} },
		{ 5, struct{ Str1, Str2 string }{"val1", "val2"}, struct{ Str1, Str2 string }{"val1", "val2"} },
		{ "prefix", struct{ Str1, Str2 string }{"val1", "val2"}, struct{ Str1, Str2 string }{"prefixval1", "prefixval2"} },
		{ 5, struct{ str1, str2 string }{"val1", "val2"}, struct{ str1, str2 string }{"val1", "val2"} },
		{ "prefix", struct{ str1, str2 string }{"val1", "val2"}, struct{ str1, str2 string }{"", ""} },
	}

	testRunArgTests(prefix, tests, tester)
}

func TestRandom(tester *testing.T) {
	passed, failed := 0, 0
	for i := 0; i < 1000; i++ {
		num := random()
		if num >= 0 && num <= 10000 {
			passed++
		} else { failed++ }
	}

	for i := 0; i < 1000; i++ {
		num := random(500)
		if num >= 0 && num <= 500 {
			passed++
		} else { failed++ }
	}

	for i := 0; i < 1000; i++ {
		num := random(-50, 50)
		if num >= -50 && num <= 50 {
			passed++
		} else { failed++ }
	}

	testFormatPassFail("random", passed, failed)
}

func TestRegexpFindAll(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, 10, [][]string{} },
		{ 0, 10, [][]string{} },
		{ "test", 10, [][]string{} },
		{ "test", "we have test", [][]string{{"test"}} },
		{ "ab", "ab ab ba ab", [][]string{{"ab"}, {"ab"}, {"ab"}} },
		{ "([^ ]*?rk)", "bark clock lark hark block", [][]string{{"bark", "bark"}, {"lark", "lark"}, {"hark", "hark"}} },
		{ "(?:[^ ]*?rk)", "bark clock lark hark block", [][]string{{"bark"}, {"lark"}, {"hark"}} },
		{ "(?:[^ ]*?rk) (?:[^ ]*?ck)", "bark clock lark hark block", [][]string{{"bark clock"}, {"hark block"}} },
		{ "(https?://){0,1}([^/ ?]+)([^ ?]+)*([^ ]*)", "https://www.test.com/page?var=1", [][]string{{"https://www.test.com/page?var=1", "https://", "www.test.com", "/page", "?var=1"}} },
	}

	testRunArgTests(regexpFindAll, tests, tester)
}

func TestRegexpReplaceAll(tester *testing.T) {
	tests := []struct { input1, input2, input3, expected any } {
		{ true, 10, 10, 10 },
		{ 0, 10, [][]string{}, [][]string{}},
		{ "find", "replace", "in", "in" },
		{ "find", "replace", "hard to find it in", "hard to replace it in" },
		{ "[^ ]in", "replace", "hard to find it in", "hard to replaced it in" },
		{ "(hard) to (find) it in", "$1 $2", "hard to find it in", "hard find" },
		{ "\n{2,}", "\n", []string{"test string", "test\nstring", "test\n\nstring"}, []string{"test string", "test\nstring", "test\nstring"} },
		{ "\n{2,}", "\n", map[int]string{1:"test string", 2:"test\nstring", 3:"test\n\nstring"}, map[int]string{1:"test string", 2:"test\nstring", 3:"test\nstring"} },
		{ "\n{2,}", "\n", struct{ Str1, Str2, Str3 string }{"test string", "test\nstring", "test\n\nstring"}, struct{ Str1, Str2, Str3 string }{"test string", "test\nstring", "test\nstring"} },
		{ "\n{2,}", "\n", struct{ str1, str2, str3 string }{"test string", "test\nstring", "test\n\nstring"}, struct{ str1, str2, str3 string }{"", "", ""} },
	}

	testRunArgTests(regexpReplaceAll, tests, tester)
}

func TestReplaceAll(tester *testing.T) {
	tests := []struct { input1, input2, input3, expected any } {
		{ true, 10, 10, 10 },
		{ 0, 10, [][]string{}, [][]string{}},
		{ "find", "replace", "in", "in" },
		{ "find", "replace", "hard to find it in", "hard to replace it in" },
		{ "find", "replace", "hard to find it in find", "hard to replace it in replace" },
		{ "find", "replace", []string{"test string", "find string", "find another find string"}, []string{"test string", "replace string", "replace another replace string"} },
		{ "find", "replace", map[int]string{1:"test string", 2:"find string", 3:"find another find string"}, map[int]string{1:"test string", 2:"replace string", 3:"replace another replace string"} },
		{ "find", "replace", struct{ Str1, Str2, Str3 string }{"test string", "find string", "find another find string"}, struct{ Str1, Str2, Str3 string }{"test string", "replace string", "replace another replace string"} },
		{ "find", "replace", struct{ str1, str2, str3 string }{"test string", "find string", "find another find string"}, struct{ str1, str2, str3 string }{"", "", ""} },
	}

	testRunArgTests(replaceAll, tests, tester)
}

func TestRound(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, true, true },
		{ "anything", false, false },
		{ 0, 0, 0 },
		{ 0, 10, 10 },
		{ 0, -10, -10 },
		{ 1, 10, 10 },
		{ 0, 10.5, 11.0 },
		{ 1, 10.5, 10.5 },
		{ 2, 10.5, 10.5 },
		{ 2, 10.6666666, 10.67 },
		{ 2, float32(10.6666666), float32(10.67) },
		{ 2, "test", "test" },
		{ 2, []string{"hello world"}, []string{"hello world"} },
		{ 2, []int{1, 2, 3}, []int{1, 2, 3} },
		{ 2, []float64{1, 2.2, 3.33, 4.444, 5.5555}, []float64{1, 2.2, 3.33, 4.44, 5.56} },
		{ 2, map[string]float64{"test": 3.14159}, map[string]float64{"test": 3.14} },
		{ 2, struct{Val float64}{3.14159}, struct{Val float64}{3.14} },
		{ 2, struct{val float64}{3.14159}, struct{val float64}{0} },
	}

	testRunArgTests(round, tests, tester)
}

func TestRtrim(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, true, nil },
		{ "anything", false, false },
		{ " ", 0, 0 },
		{ " ", 10, 10 },
		{ " ", -10, -10 },
		{ " ", "hello world ", "hello world" },
		{ " 123", "123 hello world 123 ", "123 hello world" },
		{ " ", []string{"   hello world "}, []string{"   hello world"} },
		{ " 123", []string{"123 hello world 123"}, []string{"123 hello world"} },
		{ " ", map[string]string{"test": " hello world "}, map[string]string{"test": " hello world"} },
		{ " 123", map[string]string{"test": " 123 hello world 123 "}, map[string]string{"test": " 123 hello world"} },
		{ " ", struct{Str string}{" Test "}, struct{Str string}{" Test"} },
		{ " ", struct{str string}{" test "}, struct{str string}{""} },
	}

	testRunArgTests(rtrim, tests, tester)
}

func TestSplit(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, true, nil },
		{ "anything", false, nil },
		{ " ", 0, nil },
		{ " ", 10, nil },
		{ " ", -10, nil },
		{ " ", "hello world ", []string{"hello", "world"} },
		{ " ", "    hello world   ", []string{"hello", "world"} },
		{ "123", "123 hello world 123 ", []string{" hello world ", " "} },
		{ " ", []string{"   hello world "}, nil },
		{ " ", map[string]string{"test": " hello world "}, nil },
		{ " ", struct{Str string}{" Test "}, nil },
	}

	testRunArgTests(split, tests, tester)
}

func TestStripTags(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, true },
		{ false, false },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ "hello world", "hello world" },
		{ "<p>hello <strong class=\"test classes\">world</p>", "hello world" },
		{ "test <script>alert('nasty');</script>", "test " },
		{ []string{"<p>hello <strong class=\"test classes\">world</p>"}, []string{"hello world"} },
		{ map[string]string{"test": "<p>hello <strong class=\"test classes\">world</p>"}, map[string]string{"test": "hello world"} },
		{ struct{Str string}{"<p>hello <strong class=\"test classes\">world</p>"}, struct{Str string}{"hello world"} },
		{ struct{str string}{"<p>hello <strong class=\"test classes\">world</p>"}, struct{str string}{""} },
	}

	testRunArgTests(stripTags, tests, tester)
}

func TestSubtract(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ false, true, true },
		{ true, false, false },
		{ int8(10), int8(20), int8(10) },
		{ 10, 20, 10 },
		{ -10, 20, 30 },
		{ -10, -20, -10 },
		{ 10.1, 20.2, 10.1 },
		{ -10.1, 20.1, 30.200000000000003 },
		{ -10.1, -20.2, -10.1},
		{ "remove", 20, 20 },
		{ "remove", "from", "from" },
		{ "remove", "from remove", "from " },
		{ []int{5}, []int{10, 5}, []int{10} },
		{ 5, []int{10, 5}, []int{5, 0} },
		{ []string{"remove"}, []string{"from", "remove"}, []string{"from"} },
		{ "remove", []string{"from", "remove"}, []string{"from", ""} },
		{ [1]string{"remove"}, [2]string{"from", "remove"}, [1]string{"from"} },
		{ "remove", [2]string{"from", "remove it"}, [2]string{"from", " it"} },
		{ map[string]string{"remove": "remove-value"}, map[string]string{"from": "from-value", "remove": "remove-value"}, map[string]string{"from": "from-value"} },
		{ "remove", map[string]string{"from": "from-value", "remove": "remove-value"}, map[string]string{"from": "from-value", "remove": "-value"} },
		{ map[string]string{"remove": "value"}, map[string]string{"from": "from-value", "remove": "remove-value"}, map[string]string{"from": "from-value", "remove": "remove-"} },
		{ struct{ Str string }{"remove"}, struct{ Str string }{}, struct{ Str string }{} },
		{ struct{ Str string }{"remove"}, struct{ Str string }{"from"}, struct{ Str string }{"from"} },
		{ struct{ Str string }{"remove"}, struct{ Str string }{"remove-value"}, struct{ Str string }{"remove-value"} },
	}

	testRunArgTests(subtract, tests, tester)
}

func TestSuffix(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ "suffix", 10, 10 },
		{ true, 10, 10 },
		{ 0, 10, 10 },
		{ "suffix", "test", "testsuffix" },
		{ "suffix", []string{"test"}, []string{"testsuffix"} },
		{ "suffix", []string{"test", "strings"}, []string{"testsuffix", "stringssuffix"} },
		{ 5, []int{10, 20}, []int{10, 20} },
		{ "suffix", []int{10, 20}, []int{10, 20} },
		{ 5, map[int]string{1: "val1", 2: "val2"}, map[int]string{1: "val1", 2: "val2"} },
		{ "suffix", map[int]string{1: "val1", 2: "val2"}, map[int]string{1: "val1suffix", 2: "val2suffix"} },
		{ 5, struct{ Str1, Str2 string }{"val1", "val2"}, struct{ Str1, Str2 string }{"val1", "val2"} },
		{ "suffix", struct{ Str1, Str2 string }{"val1", "val2"}, struct{ Str1, Str2 string }{"val1suffix", "val2suffix"} },
		{ 5, struct{ str1, str2 string }{"val1", "val2"}, struct{ str1, str2 string }{"val1", "val2"} },
		{ "suffix", struct{ str1, str2 string }{"val1", "val2"}, struct{ str1, str2 string }{"", ""} },
	}

	testRunArgTests(suffix, tests, tester)
}

func TestTime(tester *testing.T) {
	currentTime			:= time.Now().In(dateLocalTimezone)
	testTimeISO8601Z	:= "2019-04-23T11:30:21+01:00"
	testTimeISO8601		:= "2019-04-23T11:30:21+01:00"
	testTimeRFC822Z		:= "Tue, 23 Apr 19 11:30:21 +01:00"
	testTimeRFC822		:= "Tue, 23 Apr 19 11:30:21 CET"
	testTimeRFC850		:= "Tuesday, 23-Apr-19 11:30:21 CET"
	testTimeRFC1036		:= "Tue, 23 Apr 19 11:30:21 +01:00"
	testTimeRFC1123Z	:= "Tue, 23 Apr 2019 11:30:21 +01:00"
	testTimeRFC1123		:= "Tue, 23 Apr 2019 11:30:21 CET"
	testTimeRFC2822		:= "Tue, 23 Apr 2019 11:30:21 +01:00"
	testTimeRFC3339		:= "2019-04-23T11:30:21+01:00"
	testTimeWC3			:= "2019-04-23T11:30:21+01:00"
	testTimeATOM		:= "2019-04-23T11:30:21+01:00"
	testTimeCOOKIE		:= "Tuesday, 23-Apr-2019 11:30:21 CET"
	testTimeRSS			:= "Tue, 23 Apr 2019 11:30:21 +01:00"
	testTimeMYSQL		:= "2019-04-23 10:30:21"
	testTimeUNIX		:= "Tue Apr 23 11:30:21 CET 2019"
	testTimeRUBY		:= "Tue Apr 23 11:30:21 +0100 2019"
	testTimeANSIC		:= "Tue Apr 23 10:30:21 2019"

	testTime, _ := time.Parse(time.RFC3339, testTimeRFC3339)
	testTime = testTime.In(dateLocalTimezone)

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, timeFn(), currentTime.Format("15:04") },

		{ []any{testTime}, timeFn(testTime), testTime.Format("15:04") },
		{ []any{1556015421}, timeFn(1556015421), testTime.Format("15:04") },

		{ []any{"02-01-2006 15:04"}, timeFn("02-01-2006 15:04"), currentTime.Format("02-01-2006 15:04") },
		{ []any{"d-m-Y H:i"}, timeFn("d-m-Y H:i"), currentTime.Format("02-01-2006 15:04") },
		{ []any{"%d-%m-%Y %H:%M"}, timeFn("%d-%m-%Y %H:%M"), currentTime.Format("02-01-2006 15:04") },
		{ []any{"Mon 02 Jan 06 15:04"}, timeFn("Mon 02 Jan 06 15:04"), currentTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i"}, timeFn("D d M y H:i"), currentTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M"}, timeFn("%a %d %b %y %H:%M"), currentTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", 1556015421}, timeFn("Mon 02 Jan 06 15:04", 1556015421), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", testTime}, timeFn("Mon 02 Jan 06 15:04", testTime), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", testTime}, timeFn("D d M y H:i", testTime), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", testTime}, timeFn("%a %d %b %y %H:%M", testTime), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", testTimeRFC3339}, timeFn("Mon 02 Jan 06 15:04", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", testTimeRFC3339}, timeFn("D d M y H:i", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", testTimeRFC3339}, timeFn("%a %d %b %y %H:%M", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, timeFn("Mon 02 Jan 06 15:04", "2006-01-02T15:04:05Z07:00", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, timeFn("D d M y H:i", "2006-01-02T15:04:05Z07:00", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, timeFn("%a %d %b %y %H:%M", "2006-01-02T15:04:05Z07:00", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"D d M y H:i", "ISO8601Z", testTimeISO8601Z}, timeFn("D d M y H:i", "ISO8601Z", testTimeISO8601Z), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "ISO8601", testTimeISO8601}, timeFn("D d M y H:i", "ISO8601", testTimeISO8601), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC822Z", testTimeRFC822Z}, timeFn("D d M y H:i", "RFC822Z", testTimeRFC822Z), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC822", testTimeRFC822}, timeFn("D d M y H:i", "RFC822", testTimeRFC822), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC850", testTimeRFC850}, timeFn("D d M y H:i", "RFC850", testTimeRFC850), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1036", testTimeRFC1036}, timeFn("D d M y H:i", "RFC1036", testTimeRFC1036), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1123Z", testTimeRFC1123Z}, timeFn("D d M y H:i", "RFC1123Z", testTimeRFC1123Z), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1123", testTimeRFC1123}, timeFn("D d M y H:i", "RFC1123", testTimeRFC1123), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC2822", testTimeRFC2822}, timeFn("D d M y H:i", "RFC2822", testTimeRFC2822), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC3339", testTimeRFC3339}, timeFn("D d M y H:i", "RFC3339", testTimeRFC3339), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"D d M y H:i", "ATOM", testTimeATOM}, timeFn("D d M y H:i", "ATOM", testTimeATOM), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "W3C", testTimeWC3}, timeFn("D d M y H:i", "W3C", testTimeWC3), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "COOKIE", testTimeCOOKIE}, timeFn("D d M y H:i", "COOKIE", testTimeCOOKIE), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RSS", testTimeRSS}, timeFn("D d M y H:i", "RSS", testTimeRSS), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "MYSQL", testTimeMYSQL}, timeFn("D d M y H:i", "MYSQL", testTimeMYSQL), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "UNIX", testTimeUNIX}, timeFn("D d M y H:i", "UNIX", testTimeUNIX), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RUBY", testTimeRUBY}, timeFn("D d M y H:i", "RUBY", testTimeRUBY), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "ANSIC", testTimeANSIC}, timeFn("D d M y H:i", "ANSIC", testTimeANSIC), testTime.Format("Mon 02 Jan 06 15:04") },
	}

	testRunTests("time", tests, tester)
}

func TestTimeSince(tester *testing.T) {
	testTime1, _ := time.Parse(time.RFC3339, "2019-04-23T11:30:21+01:00")
	testTime2, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00+00:00")
	testTime3, _ := time.Parse(time.RFC3339, "1974-06-24T13:45:19+01:00")

	testTimes := []time.Time{
		testTime1.In(dateLocalTimezone),
		testTime2.In(dateLocalTimezone),
		testTime3.In(dateLocalTimezone),
	}

	passed, failed := 0, 0
	for _, testTime := range testTimes {
		since := time.Since(testTime).Round(time.Second)
		seconds := int(since.Seconds())

		m := timeSince(testTime)
		check := m["seconds"] + (m["minutes"] * 60) + (m["hours"] * 60 * 60) + (m["days"] * 60 * 60 * 24) + (m["weeks"] * 7 * 60 * 60 * 24) + (m["years"] * 8766 * 60 * 60)

		if check - 1 <= seconds && check + 1 >= seconds {
			passed++
		} else { failed++ }
	}

	testFormatPassFail("timesince", passed, failed)
}

func TestTimeUntil(tester *testing.T) {
	testTime1, _ := time.Parse(time.RFC3339, "2046-04-23T11:30:21+01:00")
	testTime2, _ := time.Parse(time.RFC3339, "3000-01-01T00:00:00+00:00")
	testTime3, _ := time.Parse(time.RFC3339, "2029-06-24T13:45:19+01:00")

	testTimes := []time.Time{
		testTime1.In(dateLocalTimezone),
		testTime2.In(dateLocalTimezone),
		testTime3.In(dateLocalTimezone),
	}

	passed, failed := 0, 0
	for _, testTime := range testTimes {
		since := time.Until(testTime).Round(time.Second)
		seconds := int(since.Seconds())

		m := timeUntil(testTime)
		check := m["seconds"] + (m["minutes"] * 60) + (m["hours"] * 60 * 60) + (m["days"] * 60 * 60 * 24) + (m["weeks"] * 7 * 60 * 60 * 24) + (m["years"] * 8766 * 60 * 60)

		if check - 1 <= seconds && check + 1 >= seconds {
			passed++
		} else { failed++ }
	}

	testFormatPassFail("timeuntil", passed, failed)
}

func TestTitle(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, true },
		{ false, false },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ "hello World", "Hello World" },
		{ "123 hello World", "123 Hello World" },
		{ []string{"HEllo wORld"}, []string{"Hello World"} },
		{ []string{"123 hello world"}, []string{"123 Hello World"} },
		{ map[string]string{"test": "HELLO world"}, map[string]string{"test": "Hello World"} },
		{ map[string]string{"test": "123 HELLO world"}, map[string]string{"test": "123 Hello World"} },
		{ struct{Str string}{"test"}, struct{Str string}{"Test"} },
		{ struct{str string}{"test"}, struct{str string}{""} },
	}

	testRunArgTests(title, tests, tester)
}

func TestTrim(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, true, true },
		{ "anything", false, false },
		{ " ", 0, 0 },
		{ " ", 10, 10 },
		{ " ", -10, -10 },
		{ " ", "hello world ", "hello world" },
		{ " 123", "123 hello world 123 ", "hello world" },
		{ " ", []string{"   hello world "}, []string{"hello world"} },
		{ " 123", []string{"123 hello world 123"}, []string{"hello world"} },
		{ " ", map[string]string{"test": " hello world "}, map[string]string{"test": "hello world"} },
		{ " 123", map[string]string{"test": " 123 hello world 123 "}, map[string]string{"test": "hello world"} },
		{ " ", struct{Str string}{" Test "}, struct{Str string}{"Test"} },
		{ " ", struct{str string}{" test "}, struct{str string}{""} },
	}

	testRunArgTests(trim, tests, tester)
}

func TestTruncate(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, true, true },
		{ "anything", false, false },
		{ " ", 0, 0 },
		{ " ", 10, 10 },
		{ " ", -10, -10 },
		{ -10, "hello world ", "" },
		{ 0, "hello world ", "" },
		{ 5, "hello world ", "hello" },
		{ 30, "hello world ", "hello world " },
		{ 9, "123 hello world 123 ", "123 hello" },
		{ 9, "<a href=\"#test\">123</a> <strong>hello world</strong> 123", "<a href=\"#test\">123</a> <strong>hello</strong>" },
		{ 9, "<a href=\"#test\">123 <strong>hello world</strong> 123</a>", "<a href=\"#test\">123 <strong>hello</strong></a>" },
		{ 9, "<a href=\"#te>st\">123 <strong>hello world</strong> 123</a>", "<a href=\"#te>st\">123 <strong>hello</strong></a>" },
		{ 9, "<a href = \"#te>s\\\"t\">123 <strong>hello world</strong> 123</a>", "<a href = \"#te>s\\\"t\">123 <strong>hello</strong></a>" },
		{ 9, "<a href='#te>st'>123 <strong>hello world</strong> 123</a>", "<a href='#te>st'>123 <strong>hello</strong></a>" },
		{ 9, "<a href = '#te>st'>123 <strong>hello world</strong> 123</a>", "<a href = '#te>st'>123 <strong>hello</strong></a>" },
		{ 5, []string{"123 hello world 123"}, []string{"123 h"} },
		{ 5, map[string]string{"test": " 123 hello world 123 "}, map[string]string{"test": " 123 "} },
		{ 5, struct{Str string}{" Test "}, struct{Str string}{" Test"} },
		{ 5, struct{str string}{" test "}, struct{str string}{""} },
	}

	testRunArgTests(truncate, tests, tester)
}

func TestTruncateWords(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, true, true },
		{ "anything", false, false },
		{ " ", 0, 0 },
		{ " ", 10, 10 },
		{ " ", -10, -10 },
		{ -10, "hello world how are you?", "" },
		{ 0, "hello world how are you?", "" },
		{ 3, "hello world how are you?", "hello world how" },
		{ 7, "hello world how are you?", "hello world how are you?" },
		{ 2, "123 hello! world 123 ", "123 hello!" },
		{ 2, "hello world\n\n\n\thow are you?", "hello world" },
		{ 3, "hello world\n\n\n\thow are you?", "hello world\n\n\n\thow" },
		{ 2, "<a href=\"#test\">123</a> <strong>hello world</strong> 123 how are you?", "<a href=\"#test\">123</a> <strong>hello</strong>" },
		{ 1, []string{"hello world"}, []string{"hello"} },
		{ 1, map[string]string{"test": "hello world"}, map[string]string{"test": "hello"} },
		{ 1, struct{Str string}{"hello world"}, struct{Str string}{"hello"} },
		{ 1, struct{str string}{"hello world"}, struct{str string}{""} },
	}

	testRunArgTests(truncatewords, tests, tester)
}

func TestType(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true,  "bool" },
		{ "anything", "string" },
		{ -10, "int" },
		{ int8(10), "int8" },
		{ uint32(10), "uint32" },
		{ 10.45, "float64" },
		{ []string{"hello world"}, "[]string" },
		{ [1]string{"hello world"}, "[1]string" },
		{ map[string]string{"test": "hello world"}, "map[string]string" },
		{ struct{Str string}{"hello world"}, "struct { Str string }" },
		{ struct{str string}{"hello world"}, "struct { str string }" },
	}

	testRunArgTests(typeFn, tests, tester)
}

func TestUl(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ false, "false" },
		{ true, "true" },
		{ 0, "0" },
		{ 1, "1" },
		{ -2, "-2" },
		{ 3.5, "3.5" },
		{ input1: -4.6, expected: "-4.6" },
		{ "test string", "test string" },
		{ []int{1, 2}, "<ul><li>1</li><li>2</li></ul>" },
		{ [2]int{1, 2}, "<ul><li>1</li><li>2</li></ul>" },
		{ []string{"value1", "value2"}, "<ul><li>value1</li><li>value2</li></ul>" },
		{ [2]string{"value1", "value2"}, "<ul><li>value1</li><li>value2</li></ul>" },
		{ map[int]string{1: "value1", 2: "value2"}, "<ul><li>value1</li><li>value2</li></ul>" },
		{ map[string]string{"title1": "value1", "title2": "value2"}, "<ul><li>value1</li><li>value2</li></ul>" },
		{ [][]string{{"subvalue1", "subvalue2"}, {"subvalue3", "subvalue4"}}, "<ul><li><ul><li>subvalue1</li><li>subvalue2</li></ul></li><li><ul><li>subvalue3</li><li>subvalue4</li></ul></li></ul>" },
		{ map[string]map[string]string{"title1": {"nested1": "subvalue1", "sub2": "subvalue2"}}, "<ul><li><ul><li>subvalue1</li><li>subvalue2</li></ul></li></ul>" },
	}

	testRunArgTests(ul, tests, tester)
}

func TestUpper(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, true },
		{ false, false },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ "hello World", "HELLO WORLD" },
		{ "123 hello World", "123 HELLO WORLD" },
		{ []string{"HEllo wORld"}, []string{"HELLO WORLD"} },
		{ []string{"123 hello world"}, []string{"123 HELLO WORLD"} },
		{ map[string]string{"test": "HELLO world"}, map[string]string{"test": "HELLO WORLD"} },
		{ map[string]string{"test": "123 HELLO world"}, map[string]string{"test": "123 HELLO WORLD"} },
		{ struct{Str string}{"test"}, struct{Str string}{"TEST"} },
		{ struct{str string}{"test"}, struct{str string}{""} },
	}

	testRunArgTests(upper, tests, tester)
}

func TestUrlDecode(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ false, false },
		{ true, true },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ 10.1, 10.1 },
		{ -10.1, -10.1 },
		{ "string without entities", "string without entities" },
		{ " %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D ", " ! * ' ( ) ; : @ & = + $ , / ? % # [ ] " },
		{ []string{"string without entities"}, []string{"string without entities"} },
		{ []string{" %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D "}, []string{" ! * ' ( ) ; : @ & = + $ , / ? % # [ ] "} },
		{ map[int]string{1: "string without entities"}, map[int]string{1: "string without entities"} },
		{ map[int]string{1: " %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D "}, map[int]string{1: " ! * ' ( ) ; : @ & = + $ , / ? % # [ ] "} },
		{ map[int]string{1: "string without entities", 2: " %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D "}, map[int]string{1: "string without entities", 2: " ! * ' ( ) ; : @ & = + $ , / ? % # [ ] "} },
		{ struct{ String1, String2 string }{"string without entities", " %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D "}, struct{ String1, String2 string }{"string without entities", " ! * ' ( ) ; : @ & = + $ , / ? % # [ ] "} },
		{ struct{ string1, string2 string }{"string without entities", " %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D "}, struct{ string1, string2 string }{"", ""} },
	}

	testRunArgTests(urlDecode, tests, tester)
}

func TestUrlEncode(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ false, false },
		{ true, true },
		{ 0, 0 },
		{ 10, 10 },
		{ -10, -10 },
		{ 10.1, 10.1 },
		{ -10.1, -10.1 },
		{ "string without entities", "string without entities" },
		{ " ! * ' ( ) ; : @ & = + $ , / ? % # [ ] ", " %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D " },
		{ []string{"string without entities"}, []string{"string without entities"} },
		{ []string{" ! * ' ( ) ; : @ & = + $ , / ? % # [ ] "}, []string{" %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D "} },
		{ map[int]string{1: "string without entities"}, map[int]string{1: "string without entities"} },
		{ map[int]string{1: " ! * ' ( ) ; : @ & = + $ , / ? % # [ ] "}, map[int]string{1: " %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D "} },
		{ map[int]string{1: "string without entities", 2: " ! * ' ( ) ; : @ & = + $ , / ? % # [ ] "}, map[int]string{1: "string without entities", 2: " %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D "} },
		{ struct{ String1, String2 string }{"string without entities", " ! * ' ( ) ; : @ & = + $ , / ? % # [ ] "}, struct{ String1, String2 string }{"string without entities", " %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D "} },
		{ struct{ string1, string2 string }{"string without entities", " ! * ' ( ) ; : @ & = + $ , / ? % # [ ] "}, struct{ string1, string2 string }{"", ""} },
	}

	testRunArgTests(urlEncode, tests, tester)
}

func TestWordcount(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ false, 0 },
		{ true, 0 },
		{ true, 0 },
		{ 0, 0 },
		{ 10, 0 },
		{ -10, 0 },
		{ 10.1, 0 },
		{ -10.1, 0 },
		{ "hello world", 2 },
		{ " 12 \" complex ' string ' world,together", 4 }, 
		{ "<span>simple<span> test", 2 }, 
		{ []string{"hello world"}, 0 },
		{ map[int]string{1: "hello world"}, 0 },
		{ struct{ String1 string }{"hello world"}, 0 },
		{ struct{ string1 string }{"hello world"}, 0 },
	}

	testRunArgTests(wordcount, tests, tester)
}

func TestWrap(tester *testing.T) {
	tests := []struct { input1, input2, input3, expected any } {
		{ "prefix", "suffix", 10, 10 },
		{ true, false, 10, 10 },
		{ 5, 0, 10, 10 },
		{ "prefix", "suffix", "test", "prefixtestsuffix" },
		{ "prefix", "suffix", []string{"test"}, []string{"prefixtestsuffix"} },
		{ "prefix", "suffix", []string{"test", "strings"}, []string{"prefixtestsuffix", "prefixstringssuffix"} },
		{ "prefix", 5, []int{10, 20}, []int{10, 20} },
		{ "prefix", "suffix", []int{10, 20}, []int{10, 20} },
		{ "prefix", 5, map[int]string{1: "val1", 2: "val2"}, map[int]string{1: "val1", 2: "val2"} },
		{ "prefix", "suffix", map[int]string{1: "val1", 2: "val2"}, map[int]string{1: "prefixval1suffix", 2: "prefixval2suffix"} },
		{ "prefix", 5, struct{ Str1, Str2 string }{"val1", "val2"}, struct{ Str1, Str2 string }{"val1", "val2"} },
		{ "prefix", "suffix", struct{ Str1, Str2 string }{"val1", "val2"}, struct{ Str1, Str2 string }{"prefixval1suffix", "prefixval2suffix"} },
		{ "prefix", 5, struct{ str1, str2 string }{"val1", "val2"}, struct{ str1, str2 string }{"val1", "val2"} },
		{ "prefix", "suffix", struct{ str1, str2 string }{"val1", "val2"}, struct{ str1, str2 string }{"", ""} },
	}

	testRunArgTests(wrap, tests, tester)
}

func TestYear(tester *testing.T) {
	currentTime := time.Now().In(dateLocalTimezone)
	currentYear, _, _ := currentTime.Date()

	testTime, _ := time.Parse(time.RFC3339, "2019-04-23T11:30:21+01:00")
	testTime = testTime.In(dateLocalTimezone)
	testYear := 2019

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, year(), currentYear },
		{ []any{testTime}, year(testTime), testYear },
	}

	testRunTests("year", tests, tester)
}

func TestYesNo(tester *testing.T) {
	tests1 := []struct { input1, expected any } {
		{ nil, "No" },
		{ "", "No" },
		{ 0, "No" },
		{ 0.0, "No" },
		{ 10, "Yes" },
		{ -10, "No" },
		{ 10.34, "Yes" },
		{ -10.34, "No" },
		{ true, "Yes" },
		{ false, "No" },
		{ "hello world", "Yes" },
		{ []int{}, "No" },
		{ []int{1}, "Yes" },
		{ []int{1, 2, 3}, "Yes" },
		{ []string{}, "No" },
		{ []string{"hello", "world", "how", "are", "you?"}, "Yes" },
		{ [0]string{}, "No" },
		{ [5]string{"hello", "world", "how", "are", "you?"}, "Yes" },
		{ [][]string{}, "No" },
		{ [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, "Yes" },
		{ map[int]string{}, "No" },
		{ map[int]string{1: "test", 2: "test"}, "Yes" },
		{ struct{ Str string }{ "test" }, "Yes" },
		{ struct{ str string }{ "test" }, "Yes" },
		{ struct{ Str string }{ "" }, "No" },
		{ struct{ str string }{ "" }, "No" },
	}

	tests2 := []struct { input1, input2, expected any } {
		{ "Yeah", nil, "No" },
		{ "Yeah", "", "No" },
		{ "Yeah", 0, "No" },
		{ "Yeah", 0.0, "No" },
		{ "Yeah", 10, "Yeah" },
		{ "Yeah", -10, "No" },
		{ "Yeah", 10.34, "Yeah" },
		{ "Yeah", -10.34, "No" },
		{ "Yeah", true, "Yeah" },
		{ "Yeah", false, "No" },
		{ "Yeah", "hello world", "Yeah" },
		{ "Yeah", []int{}, "No" },
		{ "Yeah", []int{1}, "Yeah" },
		{ "Yeah", []int{1, 2, 3}, "Yeah" },
		{ "Yeah", []string{}, "No" },
		{ "Yeah", []string{"hello", "world", "how", "are", "you?"}, "Yeah" },
		{ "Yeah", [0]string{}, "No" },
		{ "Yeah", [5]string{"hello", "world", "how", "are", "you?"}, "Yeah" },
		{ "Yeah", [][]string{}, "No" },
		{ "Yeah", [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, "Yeah" },
		{ "Yeah", map[int]string{}, "No" },
		{ "Yeah", map[int]string{1: "test", 2: "test"}, "Yeah" },
		{ "Yeah", struct{ Str string }{ "test" }, "Yeah" },
		{ "Yeah", struct{ str string }{ "test" }, "Yeah" },
		{ "Yeah", struct{ Str string }{ "" }, "No" },
		{ "Yeah", struct{ str string }{ "" }, "No" },
	}

	tests3 := []struct { input1, input2, input3, expected any } {
		{ "Yeah", "Nah", nil, "Nah" },
		{ "Yeah", "Nah", "", "Nah" },
		{ "Yeah", "Nah", 0, "Nah" },
		{ "Yeah", "Nah", 0.0, "Nah" },
		{ "Yeah", "Nah", 10, "Yeah" },
		{ "Yeah", "Nah", -10, "Nah" },
		{ "Yeah", "Nah", 10.34, "Yeah" },
		{ "Yeah", "Nah", -10.34, "Nah" },
		{ "Yeah", "Nah", true, "Yeah" },
		{ "Yeah", "Nah", false, "Nah" },
		{ "Yeah", "Nah", "hello world", "Yeah" },
		{ "Yeah", "Nah", []int{}, "Nah" },
		{ "Yeah", "Nah", []int{1}, "Yeah" },
		{ "Yeah", "Nah", []int{1, 2, 3}, "Yeah" },
		{ "Yeah", "Nah", []string{}, "Nah" },
		{ "Yeah", "Nah", []string{"hello", "world", "how", "are", "you?"}, "Yeah" },
		{ "Yeah", "Nah", [0]string{}, "Nah" },
		{ "Yeah", "Nah", [5]string{"hello", "world", "how", "are", "you?"}, "Yeah" },
		{ "Yeah", "Nah", [][]string{}, "Nah" },
		{ "Yeah", "Nah", [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, "Yeah" },
		{ "Yeah", "Nah", map[int]string{}, "Nah" },
		{ "Yeah", "Nah", map[int]string{1: "test", 2: "test"}, "Yeah" },
		{ "Yeah", "Nah", struct{ Str string }{ "test" }, "Yeah" },
		{ "Yeah", "Nah", struct{ str string }{ "test" }, "Yeah" },
		{ "Yeah", "Nah", struct{ Str string }{ "" }, "Nah" },
		{ "Yeah", "Nah", struct{ str string }{ "" }, "Nah" },
	}

	tests4 := []struct { input1, input2, input3, input4, expected any } {
		{ "Yeah", "Nah", "Meh", nil, "Nah" },
		{ "Yeah", "Nah", "Meh", "", "Nah" },
		{ "Yeah", "Nah", "Meh", 0, "Nah" },
		{ "Yeah", "Nah", "Meh", 0.0, "Nah" },
		{ "Yeah", "Nah", "Meh", 10, "Yeah" },
		{ "Yeah", "Nah", "Meh", -10, "Meh" },
		{ "Yeah", "Nah", "Meh", 10.34, "Yeah" },
		{ "Yeah", "Nah", "Meh", -10.34, "Meh" },
		{ "Yeah", "Nah", "Meh", true, "Yeah" },
		{ "Yeah", "Nah", "Meh", false, "Nah" },
		{ "Yeah", "Nah", "Meh", "hello world", "Yeah" },
		{ "Yeah", "Nah", "Meh", []int{}, "Nah" },
		{ "Yeah", "Nah", "Meh", []int{1}, "Yeah" },
		{ "Yeah", "Nah", "Meh", []int{1, 2, 3}, "Yeah" },
		{ "Yeah", "Nah", "Meh", []string{}, "Nah" },
		{ "Yeah", "Nah", "Meh", []string{"hello", "world", "how", "are", "you?"}, "Yeah" },
		{ "Yeah", "Nah", "Meh", [0]string{}, "Nah" },
		{ "Yeah", "Nah", "Meh", [5]string{"hello", "world", "how", "are", "you?"}, "Yeah" },
		{ "Yeah", "Nah", "Meh", [][]string{}, "Nah" },
		{ "Yeah", "Nah", "Meh", [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, "Yeah" },
		{ "Yeah", "Nah", "Meh", map[int]string{}, "Nah" },
		{ "Yeah", "Nah", "Meh", map[int]string{1: "test", 2: "test"}, "Yeah" },
		{ "Yeah", "Nah", "Meh", struct{ Str string }{ "test" }, "Yeah" },
		{ "Yeah", "Nah", "Meh", struct{ str string }{ "test" }, "Yeah" },
		{ "Yeah", "Nah", "Meh", struct{ Str string }{ "" }, "Nah" },
		{ "Yeah", "Nah", "Meh", struct{ str string }{ "" }, "Nah" },
	}

	passed, failed := 0, 0
	for _, test := range tests1 {
		if testCallVarArgs(tester, yesno, []any{test.input1}, test.expected) {
			passed++
		} else { failed++ }
	}
	for _, test := range tests2 {
		if testCallVarArgs(tester, yesno, []any{test.input1, test.input2}, test.expected) {
			passed++
		} else { failed++ }
	}
	for _, test := range tests3 {
		if testCallVarArgs(tester, yesno, []any{test.input1, test.input2, test.input3}, test.expected) {
			passed++
		} else { failed++ }
	}
	for _, test := range tests4 {
		if testCallVarArgs(tester, yesno, []any{test.input1, test.input2, test.input3, test.input4}, test.expected) {
			passed++
		} else { failed++ }
	}
	
	testFormatPassFail("yesno", passed, failed)
}