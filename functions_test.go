package templateManager

import (
	"testing"
	"time"
)

func TestAAFunctionsSetup(tester  *testing.T) {
	testsShowDetails	= true
	testsShowSuccessful = false
	consoleErrors		= false
	consoleWarnings		= false
	haltOnErrors		= false
	haltOnWarnings		= false

	initRegexps()
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
	fn := func(m map[string]any, _ error) map[string]any { return m }

	tests := []struct{ inputs []any; result any; expected any } {
		{ []any{}, fn(collection()), map[string]any{} },
		{ []any{0}, fn(collection(0)), map[string]any{} },
		{ []any{0, 0}, fn(collection(0, 0)), map[string]any{} },
		{ []any{"var", 0}, fn(collection("var", 0)), map[string]any{ "var": 0 } },
		{ []any{"var1", 0, "var2", true}, fn(collection("var1", 0, "var2", true)), map[string]any{ "var1": 0, "var2": true } },
	}

	testRunTests("collection", tests, tester)
}

func TestConcat(tester *testing.T) {
	tests := []struct { inputs []any; expected any } {
		{ []any{2, "string", 1.2345}, "2string1.2345" },
		{ []any{true, "string", []string{"one", "two"}}, "truestringonetwo" },
		{ []any{[]float64{1.23, 4.56}, "string", [][]string{{"one", "two"}, {"three"}}}, "1.234.56stringonetwothree" },
		{ []any{map[int]string{1: "one", 2: "two", 3: "three"}}, "onetwothree" },
		{ []any{struct{ num1, num2 string}{"one", "two"}}, "onetwo" },
	}

	testRunArgTests(concat, tests, tester)
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

	fn := func(d string, _ error) string { return d }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, fn(date()), currentTime.Format("02/01/2006") },

		{ []any{testTime}, fn(date(testTime)), testTime.Format("02/01/2006") },
		{ []any{1556015421}, fn(date(1556015421)), testTime.Format("02/01/2006") },

		{ []any{"02-01-2006"}, fn(date("02-01-2006")), currentTime.Format("02-01-2006") },
		{ []any{"d-m-Y"}, fn(date("d-m-Y")), currentTime.Format("02-01-2006") },
		{ []any{"%d-%m-%Y"}, fn(date("%d-%m-%Y")), currentTime.Format("02-01-2006") },
		{ []any{"Mon 02 Jan 06"}, fn(date("Mon 02 Jan 06")), currentTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y"}, fn(date("D d M y")), currentTime.Format("Mon 02 Jan 06") },
		{ []any{"%a %d %b %y"}, fn(date("%a %d %b %y")), currentTime.Format("Mon 02 Jan 06") },

		{ []any{"Mon 02 Jan 06", 1556015421}, fn(date("Mon 02 Jan 06", 1556015421)), testTime.Format("Mon 02 Jan 06") },

		{ []any{"Mon 02 Jan 06", testTime}, fn(date("Mon 02 Jan 06", testTime)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", testTime}, fn(date("D d M y", testTime)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"%a %d %b %y", testTime}, fn(date("%a %d %b %y", testTime)), testTime.Format("Mon 02 Jan 06") },

		{ []any{"Mon 02 Jan 06", testTimeRFC3339}, fn(date("Mon 02 Jan 06", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", testTimeRFC3339}, fn(date("D d M y", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"%a %d %b %y", testTimeRFC3339}, fn(date("%a %d %b %y", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06") },

		{ []any{"Mon 02 Jan 06", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, fn(date("Mon 02 Jan 06", "2006-01-02T15:04:05Z07:00", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, fn(date("D d M y", "2006-01-02T15:04:05Z07:00", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"%a %d %b %y", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, fn(date("%a %d %b %y", "2006-01-02T15:04:05Z07:00", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06") },

		{ []any{"D d M y", "ISO8601Z", testTimeISO8601Z}, fn(date("D d M y", "ISO8601Z", testTimeISO8601Z)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "ISO8601", testTimeISO8601}, fn(date("D d M y", "ISO8601", testTimeISO8601)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC822Z", testTimeRFC822Z}, fn(date("D d M y", "RFC822Z", testTimeRFC822Z)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC822", testTimeRFC822}, fn(date("D d M y", "RFC822", testTimeRFC822)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC850", testTimeRFC850}, fn(date("D d M y", "RFC850", testTimeRFC850)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC1036", testTimeRFC1036}, fn(date("D d M y", "RFC1036", testTimeRFC1036)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC1123Z", testTimeRFC1123Z}, fn(date("D d M y", "RFC1123Z", testTimeRFC1123Z)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC1123", testTimeRFC1123}, fn(date("D d M y", "RFC1123", testTimeRFC1123)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC2822", testTimeRFC2822}, fn(date("D d M y", "RFC2822", testTimeRFC2822)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RFC3339", testTimeRFC3339}, fn(date("D d M y", "RFC3339", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06") },

		{ []any{"D d M y", "ATOM", testTimeATOM}, fn(date("D d M y", "ATOM", testTimeATOM)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "W3C", testTimeWC3}, fn(date("D d M y", "W3C", testTimeWC3)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "COOKIE", testTimeCOOKIE}, fn(date("D d M y", "COOKIE", testTimeCOOKIE)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RSS", testTimeRSS}, fn(date("D d M y", "RSS", testTimeRSS)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "MYSQL", testTimeMYSQL}, fn(date("D d M y", "MYSQL", testTimeMYSQL)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "UNIX", testTimeUNIX}, fn(date("D d M y", "UNIX", testTimeUNIX)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "RUBY", testTimeRUBY}, fn(date("D d M y", "RUBY", testTimeRUBY)), testTime.Format("Mon 02 Jan 06") },
		{ []any{"D d M y", "ANSIC", testTimeANSIC}, fn(date("D d M y", "ANSIC", testTimeANSIC)), testTime.Format("Mon 02 Jan 06") },
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

	fn := func(d string, _ error) string { return d }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, fn(datetime()), currentTime.Format("02/01/2006 15:04") },

		{ []any{testTime}, fn(datetime(testTime)), testTime.Format("02/01/2006 15:04") },
		{ []any{1556015421}, fn(datetime(1556015421)), testTime.Format("02/01/2006 15:04") },

		{ []any{"02-01-2006 15:04"}, fn(datetime("02-01-2006 15:04")), currentTime.Format("02-01-2006 15:04") },
		{ []any{"d-m-Y H:i"}, fn(datetime("d-m-Y H:i")), currentTime.Format("02-01-2006 15:04") },
		{ []any{"%d-%m-%Y %H:%M"}, fn(datetime("%d-%m-%Y %H:%M")), currentTime.Format("02-01-2006 15:04") },
		{ []any{"Mon 02 Jan 06 15:04"}, fn(datetime("Mon 02 Jan 06 15:04")), currentTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i"}, fn(datetime("D d M y H:i")), currentTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M"}, fn(datetime("%a %d %b %y %H:%M")), currentTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", 1556015421}, fn(datetime("Mon 02 Jan 06 15:04", 1556015421)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", testTime}, fn(datetime("Mon 02 Jan 06 15:04", testTime)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", testTime}, fn(datetime("D d M y H:i", testTime)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", testTime}, fn(datetime("%a %d %b %y %H:%M", testTime)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", testTimeRFC3339}, fn(datetime("Mon 02 Jan 06 15:04", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", testTimeRFC3339}, fn(datetime("D d M y H:i", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", testTimeRFC3339}, fn(datetime("%a %d %b %y %H:%M", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, fn(datetime("Mon 02 Jan 06 15:04", "2006-01-02T15:04:05Z07:00", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, fn(datetime("D d M y H:i", "2006-01-02T15:04:05Z07:00", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, fn(datetime("%a %d %b %y %H:%M", "2006-01-02T15:04:05Z07:00", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"D d M y H:i", "ISO8601Z", testTimeISO8601Z}, fn(datetime("D d M y H:i", "ISO8601Z", testTimeISO8601Z)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "ISO8601", testTimeISO8601}, fn(datetime("D d M y H:i", "ISO8601", testTimeISO8601)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC822Z", testTimeRFC822Z}, fn(datetime("D d M y H:i", "RFC822Z", testTimeRFC822Z)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC822", testTimeRFC822}, fn(datetime("D d M y H:i", "RFC822", testTimeRFC822)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC850", testTimeRFC850}, fn(datetime("D d M y H:i", "RFC850", testTimeRFC850)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1036", testTimeRFC1036}, fn(datetime("D d M y H:i", "RFC1036", testTimeRFC1036)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1123Z", testTimeRFC1123Z}, fn(datetime("D d M y H:i", "RFC1123Z", testTimeRFC1123Z)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1123", testTimeRFC1123}, fn(datetime("D d M y H:i", "RFC1123", testTimeRFC1123)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC2822", testTimeRFC2822}, fn(datetime("D d M y H:i", "RFC2822", testTimeRFC2822)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC3339", testTimeRFC3339}, fn(datetime("D d M y H:i", "RFC3339", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"D d M y H:i", "ATOM", testTimeATOM}, fn(datetime("D d M y H:i", "ATOM", testTimeATOM)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "W3C", testTimeWC3}, fn(datetime("D d M y H:i", "W3C", testTimeWC3)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "COOKIE", testTimeCOOKIE}, fn(datetime("D d M y H:i", "COOKIE", testTimeCOOKIE)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RSS", testTimeRSS}, fn(datetime("D d M y H:i", "RSS", testTimeRSS)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "MYSQL", testTimeMYSQL}, fn(datetime("D d M y H:i", "MYSQL", testTimeMYSQL)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "UNIX", testTimeUNIX}, fn(datetime("D d M y H:i", "UNIX", testTimeUNIX)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RUBY", testTimeRUBY}, fn(datetime("D d M y H:i", "RUBY", testTimeRUBY)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "ANSIC", testTimeANSIC}, fn(datetime("D d M y H:i", "ANSIC", testTimeANSIC)), testTime.Format("Mon 02 Jan 06 15:04") },
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

func TestDivideCeil(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ "string", 10, 10 },
		{ true, 10, 10 },
		{ 0, 10, 10 },
		{ 3, 10, 4 },
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

	testRunArgTests(divideCeil, tests, tester)
}

func TestDivideFloor(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ "string", 10, 10 },
		{ true, 10, 10 },
		{ 0, 10, 10 },
		{ 3, 10, 3 },
		{ -2, 10, -5 },
		{ 2, -10, -5 },
		{ 2, 10.3, 5.15 },
		{ 3.3, -104.3, -31.606060606060606 },
		{ 5, "test", "test" },
		{ 5, []string{"test"}, []string{"test"} },
		{ 5, []int{10, 20}, []int{2, 4} },
		{ 5.1, []int{10, 20}, []int{1, 3} },
		{ 3.2, []float64{10, 20}, []float64{3.125, 6.25} },
		{ 5, map[string]int{"val1": 10, "val2": 20}, map[string]int{"val1": 2, "val2": 4} },
		{ 5, struct{ Num1, Num2 int }{10, 20}, struct{ Num1, Num2 int }{2, 4} },
		{ 5, struct{ num1, num2 int }{10, 20}, struct{ num1, num2 int }{0, 0} },
	}

	testRunArgTests(divideFloor, tests, tester)
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

func TestEndswith(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, "anything", false },
		{ "word", "word is first", false },
		{ input1: "first", input2: "word is first", expected: true },
		{ input1: "dog", input2: []string{"word is first"}, expected: false },
	}

	testRunArgTests(endswith, tests, tester)
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
		{ map[int]string{1: "test1", 2: "test2"}, "test1" },
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

	fn := func(t string, _ error) string { return t }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{testTime}, fn(formattime("02/01/2006 15:04", testTime)), testTime.Format("02/01/2006 15:04") },
		{ []any{testTime}, fn(formattime("d/m/Y H:i", testTime)), testTime.Format("02/01/2006 15:04") },
		{ []any{testTime}, fn(formattime("%d/%m/%Y %H:%M", testTime)), testTime.Format("02/01/2006 15:04") },
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
		{ "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff", "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &apos;html entities&apos; &amp;amp; other &quot;nasty&quot; stuff" },
		{ []string{"string without html"}, []string{"string without html"} },
		{ []string{"safe string", "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff"}, []string{"safe string", "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &apos;html entities&apos; &amp;amp; other &quot;nasty&quot; stuff"} },
		{ map[int]string{1: "string without html"}, map[int]string{1: "string without html"} },
		{ map[int]string{1: "safe string", 2: "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff"}, map[int]string{1: "safe string", 2: "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &apos;html entities&apos; &amp;amp; other &quot;nasty&quot; stuff"} },
		{ struct{ String1, String2 string }{"string without html", "\"string\" <strong>with</strong> 'html entities' &amp; other \"nasty\" stuff"}, struct{ String1, String2 string }{"string without html", "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &apos;html entities&apos; &amp;amp; other &quot;nasty&quot; stuff"} },
		{ struct{ string1, string2 string }{"string without html", "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &apos;html entities&apos; &amp;amp; other &quot;nasty&quot; stuff"}, struct{ string1, string2 string }{"", ""} },
	}

	testRunArgTests(htmlEncode, tests, tester)
}

func TestIterable(tester *testing.T) {
	fn := func(i []int, _ error) []int { return i }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{0}, fn(iterable(0)), []int{} },
		{ []any{1}, fn(iterable(1)), []int{0} },
		{ []any{5}, fn(iterable(5)), []int{0, 1, 2, 3, 4} },
		{ []any{-5}, fn(iterable(-5)), []int{} },
		{ []any{1, 3}, fn(iterable(1, 3)), []int{1, 2} },
		{ []any{5, 9}, fn(iterable(5, 9)), []int{5, 6, 7, 8} },
		{ []any{-5, -1}, fn(iterable(-5, -1)), []int{-5, -4, -3, -2} },
		{ []any{1, 3, 2}, fn(iterable(1, 3, 2)), []int{1} },
		{ []any{5, 9, 2}, fn(iterable(5, 9, 2)), []int{5, 7} },
		{ []any{5, 9, 20}, fn(iterable(5, 9, 20)), []int{5} },
		{ []any{-5, -1, -1}, fn(iterable(-5, -1, -1)), []int{} },
		{ []any{-1, -5, -1}, fn(iterable(-1, -5, -1)), []int{} },
	}

	testRunTests("iterable", tests, tester)
}

func TestJoin(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ ", ", "", "" },
		{ ", ", nil, "" },
		{ ", ", 0, "0" },
		{ ", ", -1, "-1" },
		{ ", ", 1, "1" },
		{ ", ", 0.0, "0" },
		{ ", ", 1.0, "1" },
		{ ", ", 0.1, "0.1" },
		{ ", ", 1.1, "1.1" },
		{ ", ", true, "true" },
		{ ", ", false, "false" },
		{ ", ", "string value", "string value" },

		{ ", ", []string{"string", "value"}, "string, value" },
		{ ", ", []int{1, 2}, "1, 2" },
		{ ", ", []float64{0.0, 1.1, 2.2}, "0, 1.1, 2.2" },
		{ ", ", []bool{true, false, true}, "true, false, true" },
		{ ", ", map[int]string{1: "first", 2: "second"}, "first, second" },
		{ ", ", map[int][]string{1: {"first", "second"}, 2: {"third"}}, "first, second, third" },
		{ ", ", struct{ first string; second int; third float64 } {"first", 1, 1.1}, "first, 1, 1.1" },
	}

	testRunArgTests(join, tests, tester)
}

func TestJsonDecode(tester *testing.T) {
	fn := func(j any, _ error) any { return j }

	tests := []struct{ inputs []any; result any; expected any } {
		{ []any{""}, fn(jsonDecode("")), nil },
		{ []any{"null"}, fn(jsonDecode("null")), nil },
		{ []any{"{}"}, fn(jsonDecode("{}")), map[string]any{} },
		{ []any{"[]"}, fn(jsonDecode("{}")), map[string]any{} },
		{ []any{"1"}, fn(jsonDecode("1")), 1.0 },
		{ []any{"1"}, fn(jsonDecode("1")), 1.0 },
		{ []any{"-1.5"}, fn(jsonDecode("-1.5")), -1.5 },
		{ []any{"true"}, fn(jsonDecode("true")), true },
		{ []any{"false"}, fn(jsonDecode("false")), false },
		{ []any{"string"}, fn(jsonDecode("string")), nil },
		{ []any{`"string"`}, fn(jsonDecode(`"string"`)), "string" },
		{ []any{`["string","value"]`}, fn(jsonDecode(`["string","value"]`)), []any{"string", "value"} },
		{ []any{"[1,2]"}, fn(jsonDecode("[1,2]")), []any{1.0, 2.0} },
		{ []any{"[0,1.1,2.2]"}, fn(jsonDecode("[0,1.1,2.2]")), []any{0.0, 1.1, 2.2} },
		{ []any{"[true,false,true]"}, fn(jsonDecode("[true,false,true]")), []any{true, false, true} },
		{ []any{`{"1":"first","2":"second"}`}, fn(jsonDecode(`{"1":"first","2":"second"}`)), map[string]any{"1":"first", "2":"second"} },
		{ []any{`{"1":["first","second"],"2":["third"]}`}, fn(jsonDecode(`{"1":["first","second"],"2":["third"]}`)), map[string]any{"1":[]any{"first", "second"}, "2":[]any{"third"}} },
		{ []any{`{"First":"first","Second":1,"Third":1.1}`}, fn(jsonDecode(`{"First":"first","Second":1,"Third":1.1}`)), map[string]any{"First":"first", "Second":1.0, "Third":1.1} },
	}

	testRunTests("jsonDecode", tests, tester)
}

func TestJsonEncode(tester *testing.T) {
	fn := func(j string, _ error) string { return j }

	tests := []struct{ inputs []any; result any; expected any } {
		{ []any{""}, fn(jsonEncode("")), `""` },
		{ []any{nil}, fn(jsonEncode(nil)), "null" },
		{ []any{0}, fn(jsonEncode(0)), "0" },
		{ []any{-1}, fn(jsonEncode(-1)), "-1" },
		{ []any{1}, fn(jsonEncode(1)), "1" },
		{ []any{0.0}, fn(jsonEncode(0.0)), "0" },
		{ []any{1.0}, fn(jsonEncode(1.0)), "1" },
		{ []any{0.1}, fn(jsonEncode(0.1)), "0.1" },
		{ []any{1.1}, fn(jsonEncode(1.1)), "1.1" },
		{ []any{true}, fn(jsonEncode(true)), "true" },
		{ []any{false}, fn(jsonEncode(false)), "false" },
		{ []any{"string value"}, fn(jsonEncode("string value")), `"string value"` },
		{ []any{[]string{"string", "value"}}, fn(jsonEncode([]string{"string", "value"})), `["string","value"]` },
		{ []any{[]int{1, 2}}, fn(jsonEncode([]int{1, 2})), "[1,2]" },
		{ []any{[]float64{0.0, 1.1, 2.2}}, fn(jsonEncode([]float64{0.0, 1.1, 2.2})), "[0,1.1,2.2]" },
		{ []any{[]bool{true, false, true}}, fn(jsonEncode([]bool{true, false, true})), "[true,false,true]" },
		{ []any{map[int]string{1: "first", 2: "second"}}, fn(jsonEncode(map[int]string{1: "first", 2: "second"})), `{"1":"first","2":"second"}` },
		{ []any{map[int][]string{1: {"first", "second"}, 2: {"third"}}}, fn(jsonEncode(map[int][]string{1: {"first", "second"}, 2: {"third"}})), `{"1":["first","second"],"2":["third"]}` },
		{ []any{struct{ first string; second int; third float64 } {"first", 1, 1.1}}, fn(jsonEncode(struct{ first string; second int; third float64 } {"first", 1, 1.1})), "{}" },
		{ []any{struct{ First string; Second int; Third float64 } {"first", 1, 1.1}}, fn(jsonEncode(struct{ First string; Second int; Third float64 } {"first", 1, 1.1})), `{"First":"first","Second":1,"Third":1.1}` },
	}

	testRunTests("jsonEncode", tests, tester)
}

func TestKey(tester *testing.T) {
	tests1 := []struct { input1, expected any } {
		{ nil, nil },
		{ "", nil },
		{ 0, nil },
		{ 0.0, nil },
		{ 10, nil },
		{ -10, nil },
		{ 10.34, nil },
		{ -10.34, nil },
		{ true, nil },
		{ "hello world", nil },
		{ []int{}, nil },
		{ []int{1}, nil },
		{ []int{1, 2, 3}, nil },
		{ []string{}, nil },
		{ []string{"hello", "world", "how", "are", "you?"}, nil },
		{ [0]string{}, nil },
		{ [5]string{"hello", "world", "how", "are", "you?"}, nil },
		{ [][]string{}, nil },
		{ [][]string{{"hello", "world"}, {"how", "are"}, {"you?"}}, nil },
		{ map[int]string{}, nil },
		{ map[int]string{1: "test", 2: "test"}, nil },
		{ struct{value string} {"test"}, nil },
		{ struct{Value string} {"test"}, nil },
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

func TestKeys(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true,  []int{} },
		{ 0,  []int{} },
		{ 1.23,  []int{} },
		{ "invalid",  []int{} },
		{ []int{ 1, 3, 5 },  []int{ 0, 1, 2 } },
		{ [2]string{ "one", "two" },  []int{ 0, 1 } },
		{ [][]string{ {"one", "two"}, {"three"} },  []int{ 0, 1 } },
		{ map[int]string{ 1: "one", 2: "two" },  []int{ 1, 2 } },
		{ map[string]string{ "one": "one", "two": "two" },  []string{ "one", "two" } },
		{ map[float64]int{ 1.23: 123, 4.56: 456 },  []float64{ 1.23, 4.56 } },
		{ struct{ num int; String string; float float64 }{ 1, "two", 3.0 },  []any{ "num", "String", "float" } },
	}

	testRunArgTests(keys, tests, tester)
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
		{ map[int]string{1: "test1", 2: "test2"}, "test2" },
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

	fn := func(t time.Time, _ error) time.Time { return t }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{testTime}, fn(localtime("UTC", testTime)), testTime.In(utc) },
		{ []any{testTime}, fn(localtime("Europe/London", testTime)), testTime.In(lon) },
		{ []any{testTime}, fn(localtime("EST", testTime)), testTime.In(est) },
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

func TestLpad(tester *testing.T) {
	tests := []struct { input1, input2, input3, expected any } {
		{ "10", "A", "test", "test" },
		{ false, "A", "test", "test" },
		{ 10.5, "A", "test", "test" },
		{ 10, 1, "test", "test" },
		{ 10, "A", 1, 1 },
		{ 2, "A", "test", "test" },
		{ 10, "A", "test", "AAAAAAtest" },
		{ 10, "AB", "test", "ABABABtest" },
		{ 10, "&amp;", "test", "&amp;&amp;&amp;&amp;&amp;&amp;test" },
		{ 9, "AB", "test", "BABABtest" },
		{ 9, "&amp;&amp;", "test", "&amp;&amp;&amp;&amp;&amp;test" },
		{ 9, "A&amp;", "test", "&amp;A&amp;A&amp;test" },
		{ 9, "&amp;B", "test", "B&amp;B&amp;Btest" },
		{ 9, "&&", "test", "&&&&&test" },
		{ 9, "&amp;&", "test", "&&amp;&&amp;&test" },
		{ 10, "&amp;&", "test", "&amp;&&amp;&&amp;&test" },
		{ 10, "&am", "test", "&am&amtest" },
		{ 10, "A", []string{"test1", "test2"}, []string{"AAAAAtest1", "AAAAAtest2"} },
	}

	testRunArgTests(lpad, tests, tester)
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

func TestMd5(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, "b326b5062b2f0e69046810717534cb09" },
		{ "anything", "f0e166dc34d14d6c228ffac576c9a43c" },
		{ -10, "1b0fd9efa5279c4203b7c70233f86dbf" },
		{ int8(10), "d3d9446802a44259755d38e6d163e820" },
		{ uint(10), "d3d9446802a44259755d38e6d163e820" },
		{ uint32(10), "d3d9446802a44259755d38e6d163e820" },
		{ 10.45, "f389f08b2d1aca50e981a1e91286169d" },
		{ 1.66666666667, "11f53b8d65f3b30d17267d71cd5d9142" },
		{ []string{"hello world"}, "685d85b3c8e128e36e3252f73eb8bfd5" },
		{ [2]string{"hello", "world"}, "04d32f084c0dfcaaf4b084d4c8862a28" },
		{ map[string]string{"test": "hello world"}, "1e8598410deeb6961913d063ca3e72de" },
		{ struct{Str string}{"hello world"}, "9e4cb5afcc8f0a06f1c246c4db2aa3bb" },
		{ struct{str string}{"hello world"}, "c9f807f9cf913138441f31063601a907" },
	}

	testRunArgTests(md5Fn, tests, tester)
}

func TestMktime(tester *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2019-04-23T11:30:21+01:00")
	testTime = testTime.In(dateLocalTimezone)

	fn := func(t time.Time, _ error) time.Time { return t }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, fn(mktime()), fn(now()) },
		{ []any{"invalid"}, fn(mktime("invalid")), fn(now()) },
		{ []any{"2019-04-23T11:30:21+01:00"}, fn(mktime("2019-04-23T11:30:21+01:00")), testTime },
		{ []any{"ATOM", "2019-04-23T11:30:21+01:00"}, fn(mktime("ATOM", "2019-04-23T11:30:21+01:00")), testTime },
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

	fn := func(t time.Time, _ error) time.Time { return t }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, fn(now()), testTime },
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
	tests := []struct { input1, expected any } {
		{ 0, "0th" },
		{ 1, "1st" },
		{ 2, "2nd" },
		{ 3, "3rd" },
		{ 4, "4th" },
		{ 5, "5th" },
		{ 10, "10th" },
		{ 11, "11th" },
		{ 12, "12th" },
		{ 13, "13th" },
		{ 20, "20th" },
		{ 21, "21st" },
		{ 22, "22nd" },
		{ 23, "23rd" },
		{ 101, "101st" },
		{ 102, "102nd" },
		{ 103, "103rd" },
		{ 111, "111th" },
		{ 112, "112th" },
		{ 113, "113th" },
		{ 121, "121st" },
		{ 122, "122nd" },
		{ 123, "123rd" },
		{ 1001, "1001st" },
		{ 1002, "1002nd" },
		{ 1003, "1003rd" },
		{ 1011, "1011th" },
		{ 1012, "1012th" },
		{ 1013, "1013th" },
		{ 1021, "1021st" },
		{ 1022, "1022nd" },
		{ 1023, "1023rd" },
	}

	testRunArgTests(ordinal, tests, tester)
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
	fn := func(s string, _ error) string { return s }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{0}, fn(pluralise(0)), "s" },
		{ []any{1}, fn(pluralise(1)), "" },
		{ []any{2}, fn(pluralise(2)), "s" },
		{ []any{"es", 0}, fn(pluralise("es", 0)), "es" },
		{ []any{"es", 1}, fn(pluralise("es", 1)), "" },
		{ []any{"es", 2}, fn(pluralise("es", 2)), "es" },
		{ []any{"y", "ies", 0}, fn(pluralise("y", "ies", 0)), "ies" },
		{ []any{"y", "ies", 1}, fn(pluralise("y", "ies", 1)), "y" },
		{ []any{"y", "ies", 2}, fn(pluralise("y", "ies", 2)), "ies" },
		{ []any{1.5}, fn(pluralise(1.5)), "" },
		{ []any{false}, fn(pluralise(false)), "" },
		{ []any{[]string{"test"}}, fn(pluralise([]string{"test"})), "" },
		{ []any{map[int]string{1: "test"}}, fn(pluralise(map[int]string{1: "test"})), "" },
		{ []any{struct{ Str string }{"test"}}, fn(pluralise(struct{ Str string }{"test"})), "" },
	}

	testRunTests("pluralise", tests, tester)
}

func TestPrefix(tester *testing.T) {
	tests := []struct { inputs []any; expected any } {
		{ []any{"prefix", 10}, 10 },
		{ []any{true, 10}, 10 },
		{ []any{0, 10}, 10 },
		{ []any{"prefix", "test"}, "prefixtest" },
		{ []any{"prefix", []string{"test"}}, []string{"prefixtest"} },
		{ []any{"prefix", []string{"test", "strings"}}, []string{"prefixtest", "prefixstrings"} },
		{ []any{5, []int{10, 20}}, []int{10, 20} },
		{ []any{"prefix", []int{10, 20}}, []int{10, 20} },
		{ []any{5, map[int]string{1: "val1", 2: "val2"}}, map[int]string{1: "5val1", 2: "5val2"} },
		{ []any{"prefix", map[int]string{1: "val1", 2: "val2"}}, map[int]string{1: "prefixval1", 2: "prefixval2"} },
		{ []any{5, struct{ Str1, Str2 string }{"val1", "val2"}}, struct{ Str1, Str2 string }{"5val1", "5val2"} },
		{ []any{"prefix", struct{ Str1, Str2 string }{"val1", "val2"}}, struct{ Str1, Str2 string }{"prefixval1", "prefixval2"} },
		{ []any{5, struct{ str1, str2 string }{"val1", "val2"}}, struct{ str1, str2 string }{"", ""} },
		{ []any{"prefix", struct{ str1, str2 string }{"val1", "val2"}}, struct{ str1, str2 string }{"", ""} },
	}

	testRunArgTests(prefix, tests, tester)
}

func TestQuery(tester *testing.T) {
	tests := []struct { input1, input2, input3, expected any } {
		{ "test", "value", 26, 26},
		{ true, "value", "/", "/"},
		{ 12, "value", "/", "/"},
		{ "test", "value", "/", "/?test=value" },
		{ "test", "value", "/longer", "/longer?test=value" },
		{ "test", "value", "http://www.example.com/longer", "http://www.example.com/longer?test=value" },
		{ "test", "value", "/?test=1", "/?test=value" },
		{ "test", "value", "/?existing[]=value1&existing[]=value2", "/?existing[]=value1&existing[]=value2&test=value" },
		{ "test", "value", "/?existing[one]=value1&existing[two]=value2", "/?existing[one]=value1&existing[two]=value2&test=value" },
		{ "test", []string{"value1", "value2"}, "/", "/?test[]=value1&test[]=value2" },
		{ "test", map[string]string{"name1": "value1", "name2": "value2"}, "/", "/?test[name1]=value1&test[name2]=value2" },
		{ "test", struct{name1, Name2 string}{"value1", "value2"}, "/", "/?test[name1]=value1&test[Name2]=value2" },

		{ "existing", "value", "/?existing[]=1&existing[]=2", "/?existing=value" },
		{ "existing", "value", "/?existing[one]=value1&existing[two]=value2", "/?existing=value" },
	}

	testRunArgTests(query, tests, tester)
}

func TestRandom(tester *testing.T) {
	passed, failed := 0, 0
	for i := 0; i < 1000; i++ {
		num, _ := random()
		if num >= 0 && num <= 10000 {
			passed++
		} else { failed++ }
	}

	for i := 0; i < 1000; i++ {
		num, _ := random(500)
		if num >= 0 && num <= 500 {
			passed++
		} else { failed++ }
	}

	for i := 0; i < 1000; i++ {
		num, _ := random(-50, 50)
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

func TestRpad(tester *testing.T) {
	tests := []struct { input1, input2, input3, expected any } {
		{ "10", "A", "test", "test" },
		{ false, "A", "test", "test" },
		{ 10.5, "A", "test", "test" },
		{ 10, 1, "test", "test" },
		{ 10, "A", 1, 1 },
		{ 10, "A", "test", "testAAAAAA" },
		{ 10, "AA", "test", "testAAAAAA" },
		{ 10, "&amp;", "test", "test&amp;&amp;&amp;&amp;&amp;&amp;" },
		{ 9, "AB", "test", "testABABA" },
		{ 9, "&amp;", "test", "test&amp;&amp;&amp;&amp;&amp;" },
		{ 9, "&amp;&amp;", "test", "test&amp;&amp;&amp;&amp;&amp;" },
		{ 9, "A&amp;", "test", "testA&amp;A&amp;A" },
		{ 9, "&amp;B", "test", "test&amp;B&amp;B&amp;" },
		{ 10, "A", []string{"test1", "test2"}, []string{"test1AAAAA", "test2AAAAA"} },
	}

	testRunArgTests(rpad, tests, tester)
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

func TestSha1(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, "5ffe533b830f08a0326348a9160afafc8ada44db" },
		{ "anything", "8867c88b56e0bfb82cffaf15a66bc8d107d6754a" },
		{ -10, "35c0ba310bf18ad1a4c2544a19cee254ca5d900f" },
		{ int8(10), "b1d5781111d84f7b3fe45a0852e59758cd7a87e5" },
		{ uint(10), "b1d5781111d84f7b3fe45a0852e59758cd7a87e5" },
		{ uint32(10), "b1d5781111d84f7b3fe45a0852e59758cd7a87e5" },
		{ 10.45, "31b5223a130b55497399a8ad81f9a138f0838e90" },
		{ 1.66666666667, "f18397626c3caa2085f245711190ffab00919993" },
		{ []string{"hello world"}, "7ec7019f32634da15cb51a60f4eae40eae86b254" },
		{ [2]string{"hello", "world"}, "37ea102be6932905e2ee81a96c6080984d121280" },
		{ map[string]string{"test": "hello world"}, "28803e9112f466f37d6e10974091ea111dfb023e" },
		{ struct{Str string}{"hello world"}, "2d4e2f0a92b8566a511653a81a87bf9540bc81e4" },
		{ struct{str string}{"hello world"}, "4e0c001c2ab91a196cadc0b542a181eb68e9caa8" },
	}

	testRunArgTests(sha1Fn, tests, tester)
}

func TestSha256(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, "b5bea41b6c623f7c09f1bf24dcae58ebab3c0cdd90ad966bc43a45b44867e12b" },
		{ "anything", "ee0874170b7f6f32b8c2ac9573c428d35b575270a66b757c2c0185d2bd09718d" },
		{ -10, "c171d4ec282b23db89a99880cd624e9ba2940c1d894783602edab5d7481dc1ea" },
		{ int8(10), "4a44dc15364204a80fe80e9039455cc1608281820fe2b24f1e5233ade6af1dd5" },
		{ uint(10), "4a44dc15364204a80fe80e9039455cc1608281820fe2b24f1e5233ade6af1dd5" },
		{ uint32(10), "4a44dc15364204a80fe80e9039455cc1608281820fe2b24f1e5233ade6af1dd5" },
		{ 10.45, "e6e843038cbc70663cdb09ca925d58d6a5c4dfb11bd6dd6314dca195360f5e4d" },
		{ 1.66666666667, "19fb7f024e7722d904a47fd9face8a4878bbe159df3972e2909a2ab9cba3fd81" },
		{ []string{"hello world"}, "61ab313128f7ee161d50ba5e1bba01fe7bfcceb524ecc6278e5a8c0f2be4280b" },
		{ [2]string{"hello", "world"}, "c68328a2c14c686164c662e06e0601959b398748eddd4847c52eceec2d28cc33" },
		{ map[string]string{"test": "hello world"}, "e08e788e0c24d0fef07d84a47ec8e57474e28a8ba09c13a73ad1442d165049fb" },
		{ struct{Str string}{"hello world"}, "1067722804f9932b9a9696ef7dc212c988f2cc74d58aba3008d45b62ab036fd1" },
		{ struct{str string}{"hello world"}, "2271e2072f92be152a104838db7dca3cf7bd55419eb21dbe0a0713ec8f764ca9" },
	}

	testRunArgTests(sha256Fn, tests, tester)
}

func TestSha512(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true, "9120cd5faef07a08e971ff024a3fcbea1e3a6b44142a6d82ca28c6c42e4f852595bcf53d81d776f10541045abdb7c37950629415d0dc66c8d86c64a5606d32de" },
		{ "anything", "cc27d84e5fdb68439143b6143f80ba4021e4b05380ba412b3652d56ec5ef86824da18eae36caab4a2f2aaddef32dea3058848c75f3415a0ea664d847d8e94b94" },
		{ -10, "9513a4b36d647a77c13858c6e1020d12549810cf481be796613bffd4f8ade008b1b03680db56945bb25e186a1f643aa297ab06ca4318d9b3ebc2b1b0529c473b" },
		{ int8(10), "3c11e4f316c956a27655902dc1a19b925b8887d59eff791eea63edc8a05454ec594d5eb0f40ae151df87acd6e101761ecc5bb0d3b829bf3a85f5432493b22f37" },
		{ uint(10), "3c11e4f316c956a27655902dc1a19b925b8887d59eff791eea63edc8a05454ec594d5eb0f40ae151df87acd6e101761ecc5bb0d3b829bf3a85f5432493b22f37" },
		{ uint32(10), "3c11e4f316c956a27655902dc1a19b925b8887d59eff791eea63edc8a05454ec594d5eb0f40ae151df87acd6e101761ecc5bb0d3b829bf3a85f5432493b22f37" },
		{ 10.45, "bd40229fa3cec786fff6af255e3cad67d5e32306b698e5a329be8c3d9d9361784a0adddbc4a160fc3baf1c0f717123975b919b9becef386f1dccb3b7f8303ab1" },
		{ 1.66666666667, "6d03ffed9e655f676603c9e82e3e8748fbe2d050dbf228aa4863b2394e13b7cfa288e7c29441b3b7ad1e84e7519f2c4b6cae54cf1d4f5b7e38387ff6b0a6defa" },
		{ []string{"hello world"}, "5b60dce05933d8182793214e7cf93dbd28e53fbb631fa812b399adf9069a49ff2560cadb0a3f9bf1792866355e802603629b308fc7b32501f7cc80cdd956301d" },
		{ [2]string{"hello", "world"}, "89bffb15105e6863348d60e4a5acd9d1573eecef97c89f9fa7222ee63d2b9c8e1f3f756b43504f39543fa620fe82c87b56844e43463a6460cbec0643b80a51cf" },
		{ map[string]string{"test": "hello world"}, "afdaca2d26d87b7e101f184823f3cdabfcade88d4d04fb330c6567410c74a4f8aeef6cd0638ae55ee5ae73f5dde0898a7312685a0d91e89876e76127f65a51a8" },
		{ struct{Str string}{"hello world"}, "c9db125ae03c58f5eca61cf60a706492036b2604839aa62ba35688af9279d07255b6c5cee29236b21704b3e3330b3c580617f29b2c7b48c8ee100257d01913b7" },
		{ struct{str string}{"hello world"}, "d77372b9c1d0e943f018ac9c7c03e6b798fc3c5de11b702eb2205ef1e26b2e3fef24448d2639e9c08c147b9bcc8adb198db90c8d6dac9284b9d0cfa8d00c084a" },
	}

	testRunArgTests(sha512Fn, tests, tester)
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

func TestStartswith(tester *testing.T) {
	tests := []struct { input1, input2, expected any } {
		{ true, "anything", false },
		{ "word", "word is first", true },
		{ input1: "dog", input2: "word is first", expected: false },
		{ input1: "dog", input2: []string{"word is first"}, expected: false },
	}

	testRunArgTests(startswith, tests, tester)
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

func TestSubstr(tester *testing.T) {
	tests := []struct { input1, input2, input3, expected any } {
		{ 2, 3, 12345, 345 },
		{ 2, 3, 1.2345, 234.0 },
		{ 2, 3, 12.345, 0.34 },
		{ 0, 2, 1.2345, 1.0 },
		{ 2, 3, "hello world", "llo" },
		{ 2, 0, "hello world", "llo world" },
		{ 2, -3, "hello world", "llo wo" },
		{ 2, 50, "hello world", "llo world" },
		{ 1, 2, "世界世界世界", "界世" },
		{ 2, 3, []string{"hello", "world"}, []string{"llo", "rld"} },
	}

	testRunArgTests(substr, tests, tester)
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
	tests := []struct { inputs []any; expected any } {
		{ []any{"suffix", 10}, 10 },
		{ []any{true, 10}, 10 },
		{ []any{0, 10}, 10 },
		{ []any{"suffix", "test"}, "testsuffix" },
		{ []any{"suffix", []string{"test"}}, []string{"testsuffix"} },
		{ []any{"suffix", []string{"test", "strings"}}, []string{"testsuffix", "stringssuffix"} },
		{ []any{5, []int{10, 20}}, []int{10, 20} },
		{ []any{"suffix", []int{10, 20}}, []int{10, 20} },
		{ []any{5, map[int]string{1: "val1", 2: "val2"}}, map[int]string{1: "val15", 2: "val25"} },
		{ []any{"suffix", map[int]string{1: "val1", 2: "val2"}}, map[int]string{1: "val1suffix", 2: "val2suffix"} },
		{ []any{5, struct{ Str1, Str2 string }{"val1", "val2"}}, struct{ Str1, Str2 string }{"val15", "val25"} },
		{ []any{"suffix", struct{ Str1, Str2 string }{"val1", "val2"}}, struct{ Str1, Str2 string }{"val1suffix", "val2suffix"} },
		{ []any{5, struct{ str1, str2 string }{"val1", "val2"}}, struct{ str1, str2 string }{"", ""} },
		{ []any{"suffix", struct{ str1, str2 string }{"val1", "val2"}}, struct{ str1, str2 string }{"", ""} },
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

	fn := func(s string, _ error) string { return s }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, fn(timeFn()), currentTime.Format("15:04") },

		{ []any{testTime}, fn(timeFn(testTime)), testTime.Format("15:04") },
		{ []any{1556015421}, fn(timeFn(1556015421)), testTime.Format("15:04") },

		{ []any{"02-01-2006 15:04"}, fn(timeFn("02-01-2006 15:04")), currentTime.Format("02-01-2006 15:04") },
		{ []any{"d-m-Y H:i"}, fn(timeFn("d-m-Y H:i")), currentTime.Format("02-01-2006 15:04") },
		{ []any{"%d-%m-%Y %H:%M"}, fn(timeFn("%d-%m-%Y %H:%M")), currentTime.Format("02-01-2006 15:04") },
		{ []any{"Mon 02 Jan 06 15:04"}, fn(timeFn("Mon 02 Jan 06 15:04")), currentTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i"}, fn(timeFn("D d M y H:i")), currentTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M"}, fn(timeFn("%a %d %b %y %H:%M")), currentTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", 1556015421}, fn(timeFn("Mon 02 Jan 06 15:04", 1556015421)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", testTime}, fn(timeFn("Mon 02 Jan 06 15:04", testTime)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", testTime}, fn(timeFn("D d M y H:i", testTime)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", testTime}, fn(timeFn("%a %d %b %y %H:%M", testTime)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", testTimeRFC3339}, fn(timeFn("Mon 02 Jan 06 15:04", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", testTimeRFC3339}, fn(timeFn("D d M y H:i", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", testTimeRFC3339}, fn(timeFn("%a %d %b %y %H:%M", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"Mon 02 Jan 06 15:04", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, fn(timeFn("Mon 02 Jan 06 15:04", "2006-01-02T15:04:05Z07:00", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, fn(timeFn("D d M y H:i", "2006-01-02T15:04:05Z07:00", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"%a %d %b %y %H:%M", "2006-01-02T15:04:05Z07:00", testTimeRFC3339}, fn(timeFn("%a %d %b %y %H:%M", "2006-01-02T15:04:05Z07:00", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"D d M y H:i", "ISO8601Z", testTimeISO8601Z}, fn(timeFn("D d M y H:i", "ISO8601Z", testTimeISO8601Z)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "ISO8601", testTimeISO8601}, fn(timeFn("D d M y H:i", "ISO8601", testTimeISO8601)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC822Z", testTimeRFC822Z}, fn(timeFn("D d M y H:i", "RFC822Z", testTimeRFC822Z)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC822", testTimeRFC822}, fn(timeFn("D d M y H:i", "RFC822", testTimeRFC822)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC850", testTimeRFC850}, fn(timeFn("D d M y H:i", "RFC850", testTimeRFC850)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1036", testTimeRFC1036}, fn(timeFn("D d M y H:i", "RFC1036", testTimeRFC1036)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1123Z", testTimeRFC1123Z}, fn(timeFn("D d M y H:i", "RFC1123Z", testTimeRFC1123Z)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC1123", testTimeRFC1123}, fn(timeFn("D d M y H:i", "RFC1123", testTimeRFC1123)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC2822", testTimeRFC2822}, fn(timeFn("D d M y H:i", "RFC2822", testTimeRFC2822)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RFC3339", testTimeRFC3339}, fn(timeFn("D d M y H:i", "RFC3339", testTimeRFC3339)), testTime.Format("Mon 02 Jan 06 15:04") },

		{ []any{"D d M y H:i", "ATOM", testTimeATOM}, fn(timeFn("D d M y H:i", "ATOM", testTimeATOM)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "W3C", testTimeWC3}, fn(timeFn("D d M y H:i", "W3C", testTimeWC3)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "COOKIE", testTimeCOOKIE}, fn(timeFn("D d M y H:i", "COOKIE", testTimeCOOKIE)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RSS", testTimeRSS}, fn(timeFn("D d M y H:i", "RSS", testTimeRSS)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "MYSQL", testTimeMYSQL}, fn(timeFn("D d M y H:i", "MYSQL", testTimeMYSQL)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "UNIX", testTimeUNIX}, fn(timeFn("D d M y H:i", "UNIX", testTimeUNIX)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "RUBY", testTimeRUBY}, fn(timeFn("D d M y H:i", "RUBY", testTimeRUBY)), testTime.Format("Mon 02 Jan 06 15:04") },
		{ []any{"D d M y H:i", "ANSIC", testTimeANSIC}, fn(timeFn("D d M y H:i", "ANSIC", testTimeANSIC)), testTime.Format("Mon 02 Jan 06 15:04") },
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

		m, _ := timeSince(testTime)
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

		m, _ := timeUntil(testTime)
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

func TestToBool(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ 0, false },
		{ 1, true },
		{ 12, true },
		{ -47, false },
		{ "0", false },
		{ "1", true },
		{ "5", false },
		{ "true", true },
		{ "yes", true },
		{ "y", true },
		{ "-795", false },
		{ -7.95, false },
		{ 0.2, false },
		{ 2.2, true },
		{ true, true },
		{ false, false },
		{ "string", false },
		{ []string{"string"}, false },
		{ []int{123}, false },
		{ map[int]string{1: "string"}, false },
		{ struct{ name string }{"string"}, false },
	}

	testRunArgTests(toBool, tests, tester)
}

func TestToFloat(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ 0, 0.0 },
		{ 1, 1.0},
		{ 12, 12.0 },
		{ -47, -47.0 },
		{ "0", 0.0 },
		{ "1", 1.0 },
		{ "-795", -795.0 },
		{ -7.95, -7.95 },
		{ 2.2, 2.2 },
		{ true, 1.0 },
		{ false, 0.0 },
		{ "string", 0.0 },
		{ []string{"string"}, 0.0 },
		{ []int{123}, 0.0 },
		{ map[int]string{1: "string"}, 0.0 },
		{ struct{ name string }{"string"}, 0.0 },
	}

	testRunArgTests(toFloat, tests, tester)
}

func TestToInt(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ 0, 0 },
		{ 1, 1 },
		{ 12, 12 },
		{ -47, -47 },
		{ "0", 0 },
		{ "1", 1 },
		{ "-795", -795 },
		{ -7.95, -8 },
		{ 2.2, 2 },
		{ true, 1 },
		{ false, 0 },
		{ "string", 0 },
		{ []string{"string"}, 0 },
		{ []int{123}, 0 },
		{ map[int]string{1: "string"}, 0 },
		{ struct{ name string }{"string"}, 0 },
	}

	testRunArgTests(toInt, tests, tester)
}

func TestToString(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ 0, "0" },
		{ 1, "1" },
		{ 12, "12" },
		{ -47, "-47" },
		{ "0", "0" },
		{ "1", "1" },
		{ "-795", "-795" },
		{ -7.95, "-7.95" },
		{ 2.2, "2.2" },
		{ true, "true" },
		{ false, "false" },
		{ "string", "string" },
		{ []string{"string"}, "" },
		{ []int{123}, "" },
		{ map[int]string{1: "string"}, "" },
		{ struct{ name string }{"string"}, "" },
	}

	testRunArgTests(toString, tests, tester)
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


func TestValues(tester *testing.T) {
	tests := []struct { input1, expected any } {
		{ true,  []int{} },
		{ 0,  []int{} },
		{ 1.23,  []int{} },
		{ "invalid",  []int{} },
		{ []int{ 1, 3, 5 },  []int{ 1, 3, 5 } },
		{ [2]string{ "one", "two" },  []string{ "one", "two" } },
		{ [][]string{ {"one", "two"}, {"three"} },  [][]string{ {"one", "two"}, {"three"} } },
		{ map[int]string{ 1: "one", 2: "two" },  []string{ "one", "two" } },
		{ map[string]string{ "one": "one", "two": "two" },  []string{ "one", "two" } },
		{ map[float64]int{ 1.23: 123, 4.56: 456 },  []int{ 123, 456 } },
		{ struct{ num int; String string; float float64 }{ 1, "two", 3.0 }, []any{ 1, "two", 3.0 } },
	}

	testRunArgTests(values, tests, tester)
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
		{ "prefix", 5, map[int]string{1: "val1", 2: "val2"}, map[int]string{1: "prefixval15", 2: "prefixval25"} },
		{ "prefix", "suffix", map[int]string{1: "val1", 2: "val2"}, map[int]string{1: "prefixval1suffix", 2: "prefixval2suffix"} },
		{ "prefix", 5, struct{ Str1, Str2 string }{"val1", "val2"}, struct{ Str1, Str2 string }{"prefixval15", "prefixval25"} },
		{ "prefix", "suffix", struct{ Str1, Str2 string }{"val1", "val2"}, struct{ Str1, Str2 string }{"prefixval1suffix", "prefixval2suffix"} },
		{ "prefix", 5, struct{ str1, str2 string }{"val1", "val2"}, struct{ str1, str2 string }{"", ""} },
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

	wrap := func(n int, _ error) int { return n }

	tests := []struct { inputs []any; result any; expected any } {
		{ []any{}, wrap(year()), currentYear },
		{ []any{testTime}, wrap(year(testTime)), testYear },
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