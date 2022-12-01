# Functions in `templateManager`

All functions in `templateManager` accept their principle argument **last** to allow simple chaining. *Efforts have been made to output clear errors and return suitable empty values rather than cause panics (a problem in several `text/template` functions)*.

Contents: [`add`](#add), [`bool`](#bool), [`capfirst`](#capfirst), [`collection`](#collection), [`concat`](#concat), [`contains`](#contains), [`cut`](#cut), [`date`](#date), [`datetime`](#datetime), [`default`](#default), [`divide`](#divide), [`divideceil`](#divideceil), [`dividefloor`](#dividefloor), [`divisibleby`](#divisibleby), [`dl`](#dl), [`endswith`](#endswith), [`equal`](#equal), [`first`](#first), [`firstof`](#firstof), [`float`](#float), [`formattime`](#formattime), [`gto`](#gto-greater-than), [`gte`](#gte-greater-than-equal), [`htmldecode`](#htmldecode), [`htmlencode`](#htmlencode), [`int`](#int), [`iterable`](#iterable), [`join`](#join), [`jsondecode`](#jsondecode), [`jsonencode`](#jsonencode), [`key`](#key), [`keys`](#keys), [`kind`](#kind), [`last`](#last), [`length`](#length), [`list`](#list), [`lto`](#lto-less-than), [`lte`](#lte-less-than-equal), [`localtime`](#localtime), [`lower`](#lower), [`ltrim`](#ltrim), [`md5`](#md5), [`mktime`](#mktime), [`multiply`](#multiply), [`nl2br`](#nl2br), [`notequal`](#notequal), [`now`](#now), [`ol`](#ol), [`ordinal`](#ordinal), [`paragraph`](#paragraph), [`pluralise`](#pluralise), [`prefix`](#prefix), [`query`](#query), [`random`](#random), [`regexp`](#regexp), [`regexpreplace`](#regexpreplace), [`render`](#render), [`replace`](#replace), [`round`](#round), [`rtrim`](#rtrim), [`sha1`](#sha1), [`sha256`](#sha256), [`sha512`](#sha512), [`split`](#split), [`startswith`](#startswith), [`string`](#string), [`striptags`](#striptags), [`substr`](#substr), [`subtract`](#subtract), [`suffix`](#suffix), [`time`](#time), [`timesince`](#timesince), [`timeuntil`](#timeuntil), [`title`](#title), [`trim`](#trim), [`truncate`](#truncate), [`truncatewords`](#truncatewords), [`type`](#type), [`ul`](#ul), [`upper`](#upper), [`urldecode`](#urldecode), [`urlencode`](#urlencode), [`uuid`](#uuid), [`values`](#values), [`wordcount`](#wordcount), [`wrap`](#wrap), [`year`](#year), [`yesno`](#yesno)

## `add`

```go
func add[T any](add T, to T) T
```

Adds a value to the existing item. If the added value is a simple numeric type, this will be treated as a simple addition **using floats** and applying rounding for integers *(recursively on all possible items)*. If the added value is a string, it will be appended to string values as a suffix *(recursively on all possible items)*. If the added value is a more complex type (e.g. slice / map), then it is appended / merged as appropriate in a similar manner to Django's add function. Unsupported types (e.g. booleans and structs are ignored and passed through).

Returns new variable of the original `to` data type.

```django
<!-- Integers: .Test is 10 -->
{{ add 25 .Test }} <!-- 35 -->
{{ add -30 .Test }} <!-- -20 -->
{{ add 2.5 .Test }} <!-- 13 -->
{{ add 2.4 .Test }} <!-- 12 -->
{{ add "5" .Test }} <!-- 15 -->
{{ add "5.5" .Test }} <!-- 16 -->
{{ add "string" .Test }} <!-- 10 -->
{{ add .Test "string" }} <!-- string -->

<!-- Floats: .Test is 10.0 -->
{{ add 25 .Test }} <!-- 35.0 -->
{{ add -30 .Test }} <!-- -20.0 -->
{{ add 2.5 .Test }} <!-- 12.5 -->
{{ add 2.4 .Test }} <!-- 12.4 -->
{{ add "5" .Test }} <!-- 15.0 -->
{{ add "5.5" .Test }} <!-- 15.5 -->
{{ add "string" .Test }} <!-- 10.0 -->
{{ add .Test "string" }} <!-- string -->

<!-- Strings: .Test is "test" -->
{{ add 25 .Test }} <!-- test -->
{{ add "5" .Test }} <!-- test5 -->
{{ add "5.5" .Test }} <!-- test5.5 -->
{{ add "string" .Test }} <!-- teststring -->
{{ add .Test "string" }} <!-- stringtest -->

<!-- Recursive Slices / Arrays: .Test is [1, 2, 3] -->
{{ add 25 .Test }} <!-- [26, 27, 28] (see Integers for examples) -->
{{ add "string" .Test }} <!-- [1, 2, 3] -->
<!-- Recursive Slices / Arrays: .Test is ["string", "slice"] -->
{{ add "test" .Test }} <!-- ["stringtest", "slicetest"] -->

<!-- APPEND - slices / arrays must be of the same type as added element -->
<!-- Slices / Arrays: .Test is [1, 2, 3], .Add is [4, 5, 6] -->
{{ add .Add .Test }} <!-- [1, 2, 3, 4, 5, 6] -->
<!-- Slices / Arrays: .Test is ["string", "slice"], .Add is ["addition"] -->
{{ add .Add .Test }} <!-- ["string", "slice", "addition"] -->

<!-- Recursive Maps: .Test is ["first": 1, "second": 2] -->
{{ add 25 .Test }} <!-- ["first": 26, "second": 27] (see Integers for examples) -->
{{ add "string" .Test }} <!-- ["first": 1, "second": 2] -->
<!-- Recursive Maps: .Test is ["first": "one", "second": "two"] -->
{{ add "test" .Test }} <!-- ["first": "onetest", "second": "twotest"] -->

<!-- APPEND - map values must be of the same type as added element -->
<!-- Maps: .Test is ["first": 1, "second": 2], .Add is ["third": 3, "fourth": 4] -->
{{ add .Add .Test }} <!-- ["first": 1, "second": 2, "third": 3, "fourth": 4] -->
```

## `bool`

```go
func bool(value any) bool
```

Attempts to convert any `value` to a boolean. If this is impossible, the nil value (false) will be returned. Positive integers (> 0) and floats (>= 1) are treated as true, as are several string values: "true", "1", "yes", "y". More complex types area always false.

```django
{{ bool "yes" }}
<!-- true -->

{{ bool 0 }}
<!-- false -->

{{ bool 1 }}
<!-- true -->

{{ bool 10 }}
<!-- true -->

{{ bool -1.23 }}
<!-- false -->
```

## `capfirst`

```go
func capfirst[T any](value T) T
```

Capitalises the first letter of strings. Does not alter any other letters. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain, ignoring other types *(passed through)*.

Returns new variable of the original `value` data type.

```django
{{ capfirst "this string. has two sentences." }}
<!-- This string. has two sentences. -->
```

## `collection`

```go
func collection(key1 string, value1 any, key2 string, value2 any, ...) map[string]any
```

Allows any number of pairs of keys / values to be grouped together into a map. This map can be passed to templates / blocks.

```django
{{ template "partials/requiresMultipleVars.html" collection "var1" .Var1 "var2" .Var2 }}
```

## `concat`

```go
func concat(values ...any) string
```

Concatenates any number of string-able values together in the order that they were declared.

```django
{{ concat "first" "second" 3 4.567 }}
<!-- firstsecond34.567 -->

<!-- .Test is: []string{"test", "this"}
{{ concat "first" .Test "last" }}
<!-- firsttestthislast -->
```

## `contains`

```go
func contains(find any, within any) bool
```

Determines whether the `find` value is contained in the `within` value. The `find` value can act on strings, slices, arrays and maps, but contained types must match. For more complex types, the `find` variable must match the *whole* value.

Returns a boolean value (always false for incompatible types).

```django
{{ contains "test" "A string containing test" }} <!-- true -->
{{ contains "test" "A string containing Test" }} <!-- false -->

<!-- Slices / Arrays: .Test is ["hello world", "goodbye world"] -->
{{ contains "world" .Test }} <!-- false -->
{{ contains "hello world" .Test }} <!-- true -->

<!-- Maps: .Test is ["hello": "hello world", "world": "goodbye world"] -->
{{ contains "world" .Test }} <!-- false -->
{{ contains "hello world" .Test }} <!-- true -->

<!-- Structs: .Test is [name: "hello world", Name: "goodbye world"] -->
{{ contains "world" .Test }} <!-- false -->
{{ contains "hello world" .Test }} <!-- true -->
```

## `cut`

```go
func cut[T any](remove string, from T) T
```

Will `remove` a string value that is contained in the `from` value. If `from` is a slice, array or map it will apply this conversion to any string elements that they contain, ignoring other types *(passed through)*.

Returns new variable of the original `from` data type.

```django
{{ cut "test" "A string containing test" }}
<!-- A string containing  -->
{{ cut "test" "A string containing Test" }}
<!-- A string containing Test -->

<!-- Slices / Arrays: .Test is ["hello world", "goodbye world"] -->
{{ cut "world" .Test }}
<!-- ["hello ", "goodbye "] -->
{{ cut "hello world" .Test }}
<!-- ["", "goodbye world"] -->

<!-- Maps: .Test is ["hello": "hello world", "world": "goodbye world"] -->
{{ cut "world" .Test }}
<!-- ["hello": "hello", "world": "goodbye"] -->
{{ cut "hello world" .Test }}
<!-- ["hello": "", "world": "goodbye world"] -->

<!-- Structs: .Test is [name: "hello world", Name: "goodbye world"] -->
{{ cut "world" .Test }}
<!-- [name: "", Name: "goodbye"] -->
{{ cut "hello world" .Test }}
<!-- [name: "", Name: "goodbye world"] -->
```

## `date`

```go
func date(input ...any) string
```

Parses dates to return a simple date string (by default: "d/m/Y"). Supports Go, Python and PHP formatting standards *(for input / output formatting)*. The last parameter is always the date input.

It can accept various parameter combinations:

```django
<!-- Current date and default output format -->
{{ date }}
<!-- 13/10/2022 -->

<!-- Passed `time.Time` object (default output format) -->
{{ date .Time }}
<!-- 15/02/2020 -->

<!-- Passed Unix time (default output format) -->
{{ date 1556015421 }}
<!-- 23/04/2019 -->

<!-- Current date and Go formatting string -->
{{ date "02 Jan 2006" }}
<!-- 13 Oct 2022 -->

<!-- Current date and PHP formatting string -->
{{ date "d M Y" }}
<!-- 13 Oct 2022 -->

<!-- Passed Go formatting string and `time.Time` object -->
{{ date "02 Jan 2006" .Time }}
<!-- 15 Feb 2020 -->

<!--  Passed Go formatting string and Unix time -->
{{ date "02 Jan 2006" 1556015421 }}
<!-- 23 Apr 2019 -->

<!-- Passed Python formatting string and `time.RFC3339` string -->
{{ date "%d %b %Y" "2020-02-15T11:30:12Z" }}
<!-- 15 Feb 2020 -->

<!-- Passed defined output formatting string (MYSQL), defined input layout (ATOM) and matching date string -->
{{ date "MYSQL", "ATOM" "2020-02-15T11:30:12Z00:00" }}
<!-- 2020-02-15 11:30:12 -->
```

Date and time functions support various pre-defined formats for simplicity: 

```go
`ISO8601Z`: "X-m-d\\TH:i:sP",   // "2006-01-02T15:04:05Z07:00"
`ISO8601`:  "Y-m-d\\TH:i:sO",   // "2006-01-02T15:04:05-07:00"
`RFC822Z`:  "D, d M y H:i:s O", // "Mon, 02 Jan 06 15:04 -07:00"
`RFC822`:   "D, d M y H:i:s T", // "Mon, 02 Jan 06 15:04 MST"
`RFC850`:   "l, d-M-y H:i:s T", // "Monday, 02-Jan-06 15:04:05 MST"
`RFC1036`:  "D, d M y H:i:s O", // "02 Jan 06 15:04 -07:00"
`RFC1123Z`: "D, d M Y H:i:s O", // "Mon, 02 Jan 2006 15:04:05 -07:00"
`RFC1123`:  "D, d M Y H:i:s T", // "Mon, 02 Jan 2006 15:04:05 MST"
`RFC2822`:  "D, d M Y H:i:s O", // "Mon, 02 Jan 2006 15:04:05 -07:00"
`RFC3339`:  "Y-m-d\\TH:i:sP",   // "2006-01-02T15:04:05Z07:00"
`W3C`:      "Y-m-d\\TH:i:sP",   // "2006-01-02T15:04:05Z07:00"
`ATOM`:     "Y-m-d\\TH:i:sP",   // "2006-01-02T15:04:05Z07:00"
`COOKIE`:   "l, d-M-Y H:i:s T", // "Monday, 02-Jan-2006 15:04:05 MST"
`RSS`:      "D, d M Y H:i:s O", // "Mon, 02 01 2006 15:04:05 -07:00"
`MYSQL`:    "Y-m-d H:i:s",      // "2006-01-02 15:04:05"
`UNIX`:     "D M _j H:i:s T Y", // "Mon Jan _2 15:04:05 MST 2006"
`RUBY`:     "D M d H:i:s o Y",  // "Mon Jan 02 15:04:05 -0700 2006"
`ANSIC`:    "D M _j H:i:s Y",   // "Mon Jan _2 15:04:05 2006"
```

## `datetime`

```go
func datetime(input ...any) string
```

Parses dates to return a simple date and time string (by default: "d/m/Y H:i"). Supports Go, Python and PHP formatting standards *(for input / output formatting)*. The last parameter is always the date input.

It can accept various parameter combinations:

```django
<!-- Current date and default output format -->
{{ datetime }}
<!-- 13/10/2022 12:30 -->

<!-- Passed `time.Time` object (default output format) -->
{{ datetime .Time }}
<!-- 15/02/2020 11:30 -->

<!-- Passed Unix time (default output format) -->
{{ datetime 1556015421 }}
<!-- 23/04/2019 11:30 -->

<!-- Current date and Passed Go formatting string -->
{{ datetime "02 Jan 2006 15:04" }}
<!-- 13 Oct 2022 12:30 -->

<!-- Current date and PHP formatting string -->
{{ datetime "d M Y H:i" }}
<!-- 13 Oct 2022 12:30 -->

<!-- Passed Go formatting string and `time.Time` object -->
{{ datetime "02 Jan 2006 15:04" .Time }}
<!-- 15 Feb 2020 11:30 -->

<!-- Passed Go formatting string and Unix time -->
{{ datetime "02 Jan 2006 15:04" 1556015421 }}
<!-- 23 Apr 2019 11:30 -->

<!-- Passed Python formatting string and `time.RFC3339` string -->
{{ datetime "%d %b %Y %H:%M" "2020-02-15T11:30:12Z" }}
<!-- 15 Feb 2020 11:30 -->

<!-- Passed defined output formatting string (MYSQL), defined input layout (ATOM) and matching date string -->
{{ datetime "MYSQL", "ATOM" "2020-02-15T11:30:12Z00:00" }}
<!-- 2020-02-15 11:30:12 -->
```

Date and time functions support various pre-defined formats for simplicity, see [`date`](#date).

## `default`

```go
func default(def any, test any) any
```

Will return the second `test` value if it is not empty, else return the `def` value.

```django
{{ default "default" .Empty }} <!-- default -->
{{ default "default" "Not Empty" }} <!-- Not Empty -->
```

## `divide`

```go
func divide[D int|float64, T any](divisor D, value T) T
```

Divides the `value` by the `divisor`. If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain. All values are first converted to floats, the operation is performed and then any **rounding is applied as necessary to return the item to its original type**.

Returns new variable of the original `value` data type.

```django
<!-- Integers: .Test is 10 -->
{{ divide 5 .Test }} <!-- 2 -->
{{ divide -5 .Test }} <!-- -2 -->
{{ divide 2.5 .Test }} <!-- 4 -->
{{ divide 2.4 .Test }} <!-- 4 -->
{{ divide 2.6 .Test }} <!-- 4 -->
{{ divide "5" .Test }} <!-- 2 -->
{{ divide "5.5" .Test }} <!-- 2 -->
{{ divide "string" .Test }} <!-- 10 -->
{{ divide .Test "string" }} <!-- string -->

<!-- Floats: .Test is 10.0 -->
{{ divide 5 .Test }} <!-- 2.0 -->
{{ divide -2 .Test }} <!-- -2.0 -->
{{ divide 2.5 .Test }} <!-- 4.0 -->
{{ divide 2.4 .Test }} <!-- 4.1666666 -->
{{ divide 2.6 .Test }} <!-- 3.8461538 -->
{{ divide "5" .Test }} <!-- 2.0 -->
{{ divide "5.5" .Test }} <!-- 1.81818181 -->
{{ divide "string" .Test }} <!-- 10.0 -->
{{ divide .Test "string" }} <!-- string -->

<!-- Slices / Arrays: .Test is [10, 20, 30] -->
{{ divide 2 .Test }} <!-- [5, 10, 15] (see above for examples) -->
{{ divide "string" .Test }} <!-- [10, 20, 30] -->

<!-- Maps: .Test is ["first": 10, "second": 20] -->
{{ divide 2 .Test }} <!-- ["first": 5, "second": 10] (see above for examples) -->
{{ divide "string" .Test }} <!-- ["first": 10, "second": 20] -->
```

## `divideceil`

```go
func divideceil[D int|float64, T any](divisor D, value T) T
```

Divides the `value` by the `divisor`. If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain. All values are first converted to floats, the operation is performed and then any **rounding (ceil) is applied as necessary to return the item to its original type**.

Returns new variable of the original `value` data type.

```django
<!-- Integers: .Test is 10 -->
{{ divideceil 3 .Test }} <!-- 4 -->
```

See [`divide`](#divide) for further examples.

## `dividefloor`

```go
func dividefloor[D int|float64, T any](divisor D, value T) T
```

Divides the `value` by the `divisor`. If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain. All values are first converted to floats, the operation is performed and then any **rounding (floor) is applied as necessary to return the item to its original type**.

Returns new variable of the original `value` data type.

```django
<!-- Integers: .Test is 10 -->
{{ dividefloor 4 .Test }} <!-- 2 -->
```

See [`divide`](#divide) for further examples.

## `divisibleby`

```go
func divisibleby[T any](divisor int, value T) bool
```

Determines if the `value` is exactly divisible by the `divisor`. Non-numeric values return false.

```django
{{ divisibleby 2 20 }} <!-- true -->
{{ divisibleby 2 19 }} <!-- false -->
{{ divisibleby 2.5 20 }} <!-- true -->
{{ divisibleby 2.6 20 }} <!-- false -->
{{ divisibleby 0.8 2.4 }} <!-- true -->
{{ divisibleby "2" 20 }} <!-- true -->
{{ divisibleby 2 "20" }} <!-- false -->
{{ divisibleby 2 true }} <!-- false -->
{{ divisibleby true 20 }} <!-- false -->
```

## `dl`

```go
func dl(value any) string
```

Converts slices, arrays or maps into an HTML definition list. For maps this will use the keys as the dt elements.

Other data types will just return a string representation of themselves.

```django
<!-- .Test is slice: [1, 2, 3] -->
{{ dl .Test }}
<!-- produces: -->
<dl>
	<dd>1</dd>
	<dd>2</dd>
	<dd>3</dd>
</dl>

<!-- .Test is map: ["first": "first-content", "second": "second-content"] -->
{{ dl .Test }}
<!-- produces: -->
<dl>
	<dt>first</dt>
	<dd>first-content</dd>
	<dt>second</dt>
	<dd>second-content</dd>
</dl>

<!-- .Test is map: ["first": ["slice", "one"], "second": ["slice", "two"]] -->
{{ dl .Test }}
<!-- produces: -->
<dl>
	<dt>first</dt>
	<dd>
		<dl>
			<dd>slice</dd>
			<dd>one</dd>
		</dl>
	</dd>
	<dt>second</dt>
	<dd>
		<dl>
			<dd>slice</dd>
			<dd>two</dd>
		</dl>
	</dd>
</dl>
```

## `endswith`

```go
func endswith(find string, value string) bool
```

Tests whether a string `value` ends with another string, `find`. Returns false if either value is not a string.

```django
{{ if endswith "fox" "the quick brown fox" }} ends with fox {{ end }}
<!-- ends with fox -->

<!-- OR more idiomatically...  -->
{{ if "the quick brown fox" | endswith "fox" }} ends with fox {{ end }}
<!-- ends with fox -->
```

## `equal`

```go
func equal(values ...any) bool
```

Tests whether any number of variables are equal. A safer alternative *(no panics)* for the `text/template` [`eq`](BASICS.md#`eq`) function. For numeric variables this is **loose** equality, in that types will be ignored (all values converted to `float64`) and compared with a small tolerance. This function is not safe for very large `uint64` values.

```django
{{ equal "hello" "hello" }} <!-- true -->
{{ equal 1 1.0 }} <!-- true -->
{{ equal 1 1.0000000000001 }} <!-- true -->
{{ equal (divide 0.8 2.4) 3 }} <!-- true -->
{{ equal 1 1.0 1.0000000000001 }} <!-- true -->
{{ equal 1 1.0 1.0000000000002 }} <!-- false -->

<!-- .Test1 and .Test2 are both [1, 2, 3] -->
{{ equal .Test1 .Test2 }} <!-- true -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`eq`](BASICS.md#eq) function)*

## `first`

```go
func first(value any) any
```

Gets the first value from slices / arrays or the first word from strings. **All other data types return an empty variable.**

```django
{{ first "my test string" }} <!-- my -->

<!-- .Test is [1, 2, 3] -->
{{ first .Test }} <!-- 1 -->
```

## `firstof`

```go
func firstof(values ...any) any
```

Accepts any number of values and returns the first one of them that exists and is not empty. If none are found it returns an empty value.

```django
{{ firstof .Empty .AlsoEmpty .NotEmpty }} <!-- .NotEmpty -->
{{ firstof .Empty "" 0 .NotEmpty .AlsoEmpty }} <!-- .NotEmpty -->
```

## `float`

```go
func float(value any) float64
```

Attempts to convert any `value` to a float64. If this is impossible, the nil value (0.0) will be returned. More complex types are always 0.0.

```django
{{ float "3" }}
<!-- 3.0 -->

{{ float 1 }}
<!-- 1.0 -->

{{ float true }}
<!-- 1.0 -->

{{ float -1.23 }}
<!-- -1.23 -->
```

## `formattime`

```go
func formattime(format string, t time.Time) string
```

Formats a time.Time object for display.

```django
{{ now | formattime "d/m/y H:i:s" }}
```

Date and time functions support various pre-defined formats for simplicity, see [`date`](#date).

## `gto` (greater than)

```go
func gto(value1 any, value2 any) bool
```

Determines if `value2` is greater than `value1` with loose numeric types. Only works on numeric types, anything else will return false. **Arguments are reversed for pipelining.**

```django
{{ 2.4 | gto 2 }} <!-- true -->
{{ 2.4 | gto 2.4 }} <!-- false -->

<!-- Opposite to normal variable order for pipelining -->
{{ gto 2 2.4 }} <!-- true -->
{{ gto 2.4 2 }} <!-- false -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`gt`](BASICS.md#gt) function)*

## `gte` (greater than equal)

```go
func gte(value1 any, value2 any) bool
```

Determines if `value2` is greater than or equal to `value1` with loose numeric types. Only works on numeric types, anything else will return false. **Arguments are reversed for pipelining.**

```django
{{ 2.4 | gte 2 }} <!-- true -->
{{ 2.4 | gte 2.4 }} <!-- true -->

<!-- Opposite to normal variable order for pipelining -->
{{ gte 2 2.4 }} <!-- true -->
{{ gte 2.4 2 }} <!-- false -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`ge`](BASICS.md#ge) function)*

## `htmldecode`

```go
func htmldecode[T any](value T) T
```

Converts HTML character-entity equivalents back into their literal, usable forms. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ htmldecode "&quot;string&quot; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#x27; &amp;amp; other &#34;nasty&#x22; stuff" }}
<!-- "string" <strong>with</strong> 'html entities' &amp; other "nasty" stuff -->
```

## `htmlencode`

```go
func htmlencode[T any](value T) T
```

Converts literal HTML special characters into safe, character-entity equivalents. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ htmlencode `"string" <strong>with</strong> 'html entities' &amp; other "nasty" stuff` }}
<!-- &#34;string&#34; &lt;strong&gt;with&lt;/strong&gt; &#39;html entities&#39; &amp;amp; other &#34;nasty&#34; stuff -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`html`](BASICS.md#html) function)*

## `int`

```go
func int(value any) int
```

Attempts to convert any `value` to an int. Floats will be rounded. If conversion is impossible, the nil value (0) will be returned. More complex types are always 0.

```django
{{ int "3" }}
<!-- 3 -->

{{ int 1 }}
<!-- 1 -->

{{ int true }}
<!-- 1 -->

{{ int 4.67 }}
<!-- 5 -->

{{ int -1.23 }}
<!-- -1 -->
```

## `iterable`

```go
func iterable(options ...int) []int
```

There is no `for` loop in the `text/template` package, but one can be simulated using the `range` function and an integer slice. The `iterable` function allows a slice to be created for this purpose.

The possible options are best illustrated with examples:

```django
<!-- for i := 0; i < 5; i++ -->
{{ range $i := iterable 5 }} {{ $i }} {{ end }}

<!-- for i := 3; i < 5; i++ -->
{{ range $i := iterable 3 5 }} {{ $i }} {{ end }}

<!-- for i := 3; i < 5; i += 2 -->
{{ range $i := iterable 3 5 2 }} {{ $i }} {{ end }}
```

## `join`

```go
func join(separator string, values any) string
```

Joins slice or map `values` together as a string spaced by the `separator`. Strings are returned unaltered and numeric types are coerced into strings. Maps are sorted alphabetically / numerically by their keys first for predictable results. 

```django
<!-- Slice / Array: .Test is [1, 2, 3] -->
{{ join "::" .Test }} <!-- 1::2::3 -->

<!-- Map: .Test is [2: "second", 1: "first"] -->
{{ join ", " .Test }} <!-- first, second -->
```

## `jsondecode`

```go
func jsondecode(value any) any
```

Decodes any JSON string to an `interface{}`. This usually produces a type: `map[string]any`, but may result in other types (e.g. `[]any`) or simple types (e.g. `bool`, `string`) for trivial data sources. ALL numbers will be `float64`.

**Use only with caution / testing.**

## `jsonencode`

```go
func jsonencode(value any) string
```

Encodes any value to a JSON string.

```django
<!-- .Test is map: ["first": ["slice", "one"], "second": ["slice", "two"]] -->
{{ jsonencode .Test }}
<!-- { "first": ["slice", "one"], "second": ["slice", "two"] } -->
```

## `key`

```go
func key(input ...any) any
```

Very similar to the in-built `text/template` [`index`](BASICS.md#general-utility-functions) function, `key` accepts any number of nested keys and returns the result of indexing its **final argument** by them. For strings this returns individual letters. The indexed item must be a string, map, slice, array or struct.

```django
{{ key 2 "string" }} <!-- r -->

<!-- Slices / Arrays: .Test is ["first", "second", "third"] -->
{{ key 2 .Test }} <!-- third -->
{{ key 2 2 .Test }} <!-- i -->
{{ key 2 2 0 .Test }} <!-- i -->
{{ key 2 2 2 .Test }} <!-- -->

<!-- Maps: .Test is ["first": ["nested": "nested-value"]] -->
{{ key "first" .Test }} <!-- ["nested": "nested-value"] -->
{{ key "first" "nested" .Test }} <!-- nested-value -->
{{ key "first" "nested" 3 .Test }} <!-- t -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`index`](BASICS.md#index) function)*

## `keys`

```go
func keys(value any) []any
```

Returns the keys of a slice / array / map / struct.

```django
{{ keys .Slice }}
<!-- []int{0, 1, 2} -->

{{ keys .Map }}
<!-- []string{"key1", "key2", "key3"} -->
```

## `kind`

```go
func kind(value any) string
```

Returns the value's `reflect.Kind` as a string. Mainly useful for testing.

```django
{{ kind 3.14159 }} <!-- float64 -->
{{ kind "test" }} <!-- string -->

<!-- Slice: .Test is []int{1, 2, 3} -->
{{ kind .Test }} <!-- slice -->

<!-- Array: .Test is [3]int{1, 2, 3} -->
{{ kind .Test }} <!-- array -->

<!-- Map: .Test is map[int]string{2: "second", 1: "first"} -->
{{ kind .Test }} <!-- map -->
```

## `last`

```go
func last(value any) any
```

Gets the last value from slices / arrays or the last word from strings. **All other data types return an empty variable.**

```django
{{ last "my test string" }} <!-- string -->

<!-- .Test is [1, 2, 3] -->
{{ last .Test }} <!-- 3 -->
```

## `length`

```go
func length(value any) int
```

Gets length of any type. Unlike the `text/template` function: `len`, the `length` function will work without panics on numeric types and booleans.

```django
{{ length "my test string" }} <!-- 14 -->
{{ length 12 }} <!-- 2 -->
{{ length -3.14159 }} <!-- 8 -->
{{ length true }} <!-- 1 -->

<!-- .Test is [1, 2, 3] -->
{{ length .Test }} <!-- 3 -->

<!-- .Test is [1:"first", 2:"second"] -->
{{ length .Test }} <!-- 2 -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`len`](BASICS.md#len) function)*

## `list`

```go
func list(values ...any) []any
```

Accepts any number of arguments / types and returns a slice of interfaces that may be iterated over.

```django
{{ $list := list 1 "string" true }}
<!-- [1, "string", true] -->
```

## `lto` (less than)

```go
func lto(value1 any, value2 any) bool
```

Determines if `value2` is less than `value1` with loose numeric types. Only works on numeric types, anything else will return false. **Arguments are reversed for pipelining.**

```django
{{ 2.4 | lto 3 }} <!-- true -->
{{ 2.4 | lto 2 }} <!-- false -->
{{ 2.4 | lto 2.4 }} <!-- false -->

<!-- Opposite to normal variable order for pipelining -->
{{ lto 2 2.4 }} <!-- false -->
{{ lto 2.4 2 }} <!-- true -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`lt`](BASICS.md#lt) function)*

## `lte` (less than equal)

```go
func lte(value1 any, value2 any) bool
```

Determines if `value2` is less than or equal to `value1` with loose numeric types. Only works on numeric types, anything else will return false. **Arguments are reversed for pipelining.**

```django
{{ 2.4 | lte 3 }} <!-- true -->
{{ 2.4 | lte 2 }} <!-- false -->
{{ 2.4 | lte 2.4 }} <!-- true -->

<!-- Opposite to normal variable order for pipelining -->
{{ lte 2 2.4 }} <!-- false -->
{{ lte 2.4 2 }} <!-- true -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`le`](BASICS.md#le) function)*

## `localtime`

```go
func localtime(location string|time.Location, t time.Time) time.Time
```

Localises a time object to display local times / dates. Localisation strings are system dependant.

```django
{{ now | localtime "PST" | formattime "d/m/y H:i:s" }}
```

Date and time functions support various pre-defined formats for simplicity, see [`date`](#date).

## `lower`

```go
func lower[T any](value T) T
```

Converts string text to lower case. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ lower "This string. Has TWO sentences." }}
<!-- this string. has two sentences. -->
```

## `ltrim`

```go
func ltrim[T any](remove string, value T) T
```

Removes the passed characters from the left end of string values. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ ltrim " " "  This string. Has TWO sentences." }}
<!-- This string. Has TWO sentences. -->

{{ ltrim " hT" "  This string. Has TWO sentences." }}
<!-- is string. Has TWO sentences. -->

{{ ltrim "hT" "  This string. Has TWO sentences." }}
<!-- This string. Has TWO sentences. -->
```

## `md5`

```go
func md5(value any) string
```

The `md5` function creates an MD5 hash for whatever variable is passed to it. Since hashing algorithms only accept strings, any other type of variable is first given a string value and then hashed.

```django
{{ md5 "anything" }}
<!-- f0e166dc34d14d6c228ffac576c9a43c -->
```

## `mktime`

```go
func mktime(values ...any) time.Time
```

The `mktime` function creates new `time.Time` struct from simple time strings. Returns the current time if an invalid input is given. Supports Go, Python and PHP formatting standards.

It can accept various parameter combinations:

```django
<!-- Current time -->
{{ mktime }}

<!-- Parse from a `time.RFC3339` string -->
{{ mktime "2020-02-15T11:30:12Z" }}

<!-- Parse from a custom Go layout string -->
{{ mktime "2006-01-02T15:04:05Z07:00", "2020-02-15T11:30:12Z00:00" }}

<!-- Parse from a custom PHP layout string -->
{{ mktime "Y-m-d\\TH:i:sP", "2020-02-15T11:30:12Z00:00" }}

<!-- Parse from a pre-defined layout string -->
{{ mktime "ATOM", "2020-02-15T11:30:12Z00:00" }}
```

Date and time functions support various pre-defined formats for simplicity, see [`date`](#date).

## `multiply`

```go
func multiply[M int|float64, T any](multiplier M, value T) T
```

Multiplies the `value` by the `multiplier`. If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain. All values are first converted to floats, the operation is performed and then any **rounding is applied as necessary to return the item to its original type**.

```django
<!-- Integers: .Test is 10 -->
{{ multiply 5 .Test }} <!-- 50 -->
{{ multiply -5 .Test }} <!-- -50 -->
{{ multiply 2.5 .Test }} <!-- 25 -->
{{ multiply 2.4 .Test }} <!-- 24 -->
{{ multiply 2.6 .Test }} <!-- 26 -->
{{ multiply "5" .Test }} <!-- 50 -->
{{ multiply "5.5" .Test }} <!-- 55 -->
{{ multiply "string" .Test }} <!-- 10 -->
{{ multiply .Test "string" }} <!-- string -->

<!-- Floats: .Test is 10.0 -->
{{ multiply 5 .Test }} <!-- 50.0 -->
{{ multiply -5 .Test }} <!-- -50.0 -->
{{ multiply 2.5 .Test }} <!-- 25.0 -->
{{ multiply 2.4 .Test }} <!-- 24.0 -->
{{ multiply 2.6 .Test }} <!-- 26.0 -->
{{ multiply "5" .Test }} <!-- 50.0 -->
{{ multiply "5.5" .Test }} <!-- 55.0 -->
{{ multiply "string" .Test }} <!-- 10.0 -->
{{ multiply .Test "string" }} <!-- string -->

<!-- Slices / Arrays: .Test is [10, 20, 30] -->
{{ multiply 2 .Test }} <!-- [20, 40, 60] (see above for examples) -->
{{ multiply "string" .Test }} <!-- [10, 20, 30] -->

<!-- Maps: .Test is ["first": 10, "second": 20] -->
{{ multiply 2 .Test }} <!-- ["first": 20, "second": 40] (see above for examples) -->
{{ multiply "string" .Test }} <!-- ["first": 10, "second": 20] -->
```

## `nl2br`

```go
func nl2br[T any](value T) T
```

Replaces all instances of `\n` (new line) with instances of `<br>` within `value`. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain. [`paragraph`](#paragraph) will perform a similar task in a more intelligent manner.

Returns new variable of the original `value` data type.

```django
{{ nl2br "test\nstring" }} <!-- test<br>string -->
```

## `notequal`

```go
func notequal(values ...any) bool
```

Tests whether any number of variables are not equal. A safer alternative *(no panics)* for the `text/template` [`neq`](BASICS.md#`neq`) function. For numeric variables this is **loose** equality, in that types will be ignored (all values converted to `float64`) and compared with a small tolerance. This function is not safe for very large `uint64` values.

```django
{{ notequal "hello" "hello" }} <!-- false -->
{{ notequal "hello" "Hello" }} <!-- true -->
{{ notequal 1 2 }} <!-- true -->
{{ notequal 1 1.0 }} <!-- false -->
{{ notequal 1 1.0000000000001 }} <!-- false -->
{{ notequal (divide 0.8 2.4) 3 }} <!-- false -->
{{ notequal 1 1.0 1.0000000000001 }} <!-- false -->
{{ notequal 1 1.0 1.0000000000002 }} <!-- true -->

<!-- .Test1 and .Test2 are both [1, 2, 3] -->
{{ notequal .Test1 .Test2 }} <!-- false -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`ne`](BASICS.md#ne) function)*

## `now`

```go
func now() time.Time
```

Returns the current `time.Time` value.

```django
{{ now | localtime "PST" | formattime "d/m/y H:i:s" }}
```

## `ol`

```go
func ol(value any) string
```

Converts slices, arrays or maps into an HTML ordered list.

Other data types will just return a string representation of themselves.

```django
<!-- .Test is slice: [1, 2, 3] -->
{{ ol .Test }}
<!-- produces: -->
<ol>
	<li>1</li>
	<li>2</li>
	<li>3</li>
</ol>

<!-- .Test is map: ["first": "first-content", "second": "second-content"] -->
{{ ol .Test }}
<!-- produces: -->
<ol>
	<li>first-content</li>
	<li>second-content</li>
</ol>

<!-- .Test is map: ["first": ["slice", "one"], "second": ["slice", "two"]] -->
{{ ol .Test }}
<!-- produces: -->
<ol>
	<li>
		<ol>
			<li>slice</li>
			<li>one</li>
		</ol>
	</li>
	<li>
		<ol>
			<li>slice</li>
			<li>two</li>
		</ol>
	</li>
</ol>
```

## `ordinal`

```go
func ordinal[T int|float64|string](value T) string
```

Suffixes a number with the correct, English ordinal. If `value` is not numeric or a valid numeric string, an empty string is returned.

```django
{{ ordinal 1 }} <!-- 1st -->
{{ ordinal 112 }} <!-- 112th -->
{{ ordinal 1122 }} <!-- 1022nd -->
```

## `paragraph`

```go
func paragraph[T any](value T) T
```

Replaces all string instances of `\n+` (multiple new lines) with paragraph tags (`</p><p>`) and instances of `\n` (new line) with instances of `<br>` within `value`. Finally wraps the string in paragraph tags. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ paragraph "test\nstring" }} <!-- <p>test<br>string</p> -->
{{ paragraph "test\n\nstring" }} <!-- <p>test</p><p>string</p> -->
```

## `pluralise`

```go
func pluralise(values ...any) string
```

Allows pluralisation of word endings. Allows basic customisation of the possible singular / plural forms. The default singular suffix is empty and the default plural suffix is "s".

It can accept various parameter combinations:

```django
1 cat{{ pluralise 1 }}
<!-- 1 cat -->

2 cat{{ pluralise 2 }}
<!-- 2 cats -->

0 cat{{ pluralise 0 }}
<!-- 0 cats -->

1 mattress{{ pluralise "es" 1 }}
<!-- 1 mattress -->

2 mattress{{ pluralise "es" 2 }}
<!-- 2 mattresses -->

1 cherr{{ pluralise "y" "ies" 1 }}
<!-- 1 cherry -->

2 cherr{{ pluralise "y" "ies" 2 }}
<!-- 2 cherries -->
```

## `prefix`

```go
func prefix[T any](prefix ...string, value T) T
```

Prefixes all strings within `value` with all strings within `prefix`. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ prefix "prefix " "value" }} <!-- prefix value -->
{{ prefix "prefix1 " "prefix2 " "value" }} <!-- prefix1 prefix2 value -->

<!-- Slices / Arrays: .Test is ["string1", "string2"] -->
{{ prefix "prefix " .Test }}
<!-- ["prefix string1", "prefix string2"] -->

<!-- Maps: .Test is [1: "string1", 2: "string2"] -->
{{ prefix "prefix " .Test }}
<!-- [1: "prefix string1", 2: "prefix string2"] -->
```

## `query`

```go
func query[T any](name string, value any, link T) T
```

Acts upon URL query strings and allows parameters to be added / replaced within them *(the parameter being acted upon, `name`, will always be removed prior to the new value being added)*. If value is a slice / array, then an array variable will be created. If value is a map or struct, then a map variable will be created.

Returns new variable of the original `value` data type.

```django
{{ query "test" "value" "/" }}
<!-- /?test=value -->

{{ query "test" "value" "/?test=1" }}
<!-- /?test=value -->

{{ query "test" "value" "/?existing=1" }}
<!-- /?existing=1&test=value -->

<!-- .Slice is: []string{"value1", "value2"} -->
{{ query "test" .Slice "/" }}
<!-- /?test[]=value1&test[]=value2 -->

<!-- .Map is: map[string]string{"name1": "value1", "name2": "value2"} -->
{{ query "test" .Map "/" }}
<!-- /?test[name1]=value1&test[name2]=value2 -->

<!-- .Struct is: struct{ name1, name2 string }{"value1", "value2"} -->
{{ query "test" .Struct "/" }}
<!-- /?test[name1]=value1&test[name2]=value2 -->
```

## `random`

```go
func random(...int) int
```

Generates random numbers.

It can accept various parameter combinations:

```django
<!-- Returns a random number between 0 and 10000 -->
{{ random }}

<!-- Returns a random number between 0 and 100 -->
{{ random 100 }}

<!-- Returns a random number between 200 and 1000 -->
{{ random 200, 1000 }}

<!-- Returns a random number between 200 and 1000 -->
{{ random 1000, 200 }}

<!-- Returns a random number between -50 and 50 -->
{{ random -50, 50 }}
```

## `regexp`

```go
func regexp(find string, value string) [][]string
```

Finds all instances of `find` regexp within `value` using [`regexp.FindAllStringSubmatch`](https://pkg.go.dev/regexp#Regexp.FindAllStringSubmatch). It only acts on strings, returning an empty string slice for any other values.

```django
{{ regexp "(?:[^ ]*?rk)" "bark clock lark hark block" }}
<!-- [["bark"], ["lark"], ["hark"]]-->
```

## `regexpreplace`

```go
func regexpreplace[T any](find string, replace string, value T) T
```

Replaces all instances of `find` regexp with instances of `replace` within `value` using [`regexp.ReplaceAllString`](https://pkg.go.dev/regexp#Regexp.ReplaceAllString). If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ regexpreplace "\n{2,}", "\n", "test\n\n\nstring" }}
<!-- test\nstring -->

{{ regexpreplace "[^ ]in", "replace", "hard to find it in" }}
<!-- hard to replaced it in -->
```

## `render`

```go
func render(name string, parameters any) string
```

Renders a template defined within the current template scope. This differs from the `text/template` defined `template` keyword in that its scope is more limited (to templates defined within the template itself), but will also accept a variable as the input variable.

This function is mainly used with [components](COMPONENTS.md#passing-data-to-components) when passing wrapped values.

```django
{{ define "template123" }} Template content here {{ end }}
{{ $template := "template123" }}
{{ render $template . }}
<!-- Template content here -->
```

## `replace`

```go
func replace[T any](find string, replace string, value T) T
```

Replaces all instances of `find` with instances of `replace` within `value`.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ replace "find", "replace", "test find string where find is replaced" }}
<!-- test replace string where replace is replaced -->
```

## `round`

```go
func round[T any](precision int, value T) T
```

Rounds floats to the passed number of decimal places (`precision`). If `value` is a slice, array or map it will apply this conversion to any float elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ round 3 3.14159 }}
<!-- 3.1416 -->
```

## `rtrim`

```go
func rtrim[T any](remove string, value T) T
```

Removes the passed characters from the right end of string values. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ rtrim " " "This string. Has TWO sentences.  " }}
<!-- This string. Has TWO sentences. -->

{{ rtrim " ." "This string. Has TWO sentences.  " }}
<!-- This string. Has TWO sentences -->

{{ rtrim "hT" "  This string. Has TWO sentences.  " }}
<!-- This string. Has TWO sentences.   -->
```

## `sha1`

```go
func sha1(value any) string
```

The `sha1` function creates a SHA1 hash for whatever variable is passed to it. Since hashing algorithms only accept strings, any other type of variable is first given a string value and then hashed.

```django
{{ sha1 "anything" }}
<!-- 8867c88b56e0bfb82cffaf15a66bc8d107d6754a -->
```

## `sha256`

```go
func sha256(value any) string
```

The `sha256` function creates a SHA256 hash for whatever variable is passed to it. Since hashing algorithms only accept strings, any other type of variable is first given a string value and then hashed.

```django
{{ sha256 "anything" }}
<!-- ee0874170b7f6f32b8c2ac9573c428d35b575270a66b757c2c0185d2bd09718d -->
```

## `sha512`

```go
func sha512(value any) string
```

The `sha512` function creates a SHA512 hash for whatever variable is passed to it. Since hashing algorithms only accept strings, any other type of variable is first given a string value and then hashed.

```django
{{ sha512 "anything" }}
<!-- cc27d84e5fdb68439143b6143f80ba4021e4b05380ba412b3652d56ec5ef86824da18eae36caab4a2f2aaddef32dea3058848c75f3415a0ea664d847d8e94b94 -->
```

## `split`

```go
func split(separator string, value string) []string
```

Splits strings on the `separator` value and returns a slice of the pieces. It only works on strings and returns an empty slice for all other data types.

```django
{{ split " " "a test string" }}
<!-- ["a", "test", "string"] -->

{{ split "::" "some::joined::data" }}
<!-- ["some", "joined", "data"] --> 
```

## `startswith`

```go
func startswith(find string, value string) bool
```

Tests whether a string `value` starts with another string, `find`. Returns false if either value is not a string.

```django
{{ if startswith "the" "the quick brown fox" }} starts with the {{ end }}
<!-- starts with the -->

<!-- OR more idiomatically...  -->
{{ if "the quick brown fox" | startswith "the" }} starts with the {{ end }}
<!-- starts with the -->
```

## `string`

```go
func string(value any) int
```

Attempts to convert any `value` to a string. If conversion is impossible, the nil value ("") will be returned. More complex types are always "".

```django
{{ string "3" }}
<!-- 3 -->

{{ string 3 }}
<!-- 3 -->
```

## `striptags`

```go
func striptags[T any](value T) T
```

Strips HTML tags from strings. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ striptags "<p>Remove <strong>all</strong> HTML tags</p>" }}
<!-- Remove all HTML tags -->
```

## `substr`

```go
func substr[T any](offset int, length int, value T) T
```

Extracts a substring from a `value` starting at the specified `offset` and including `length` runes from that point. If `length` is zero, the substring will be taken to the end of the `value`. If `length` is negative, it is the number of runes from the end of the `value`.

If `value` is a slice, array or map it will apply this conversion to any valid elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ substr 3 13 "<p>Remove P Tags</p>" }}
<!-- Remove P Tags -->

{{ "<p>Remove P Tags</p>" | substr 3 -4 }}
<!-- Remove P Tags -->

{{ "Do nothing" | substr 3 0 }}
<!-- nothing -->
```

## `subtract`

```go
func subtract[T any](subtract T, from T) T
```

Subtracts a value from the existing item. If the subtracted value is a simple numeric type this will be treated as a simple maths *(recursively on all possible items)*. If the removed value is a string, it will be [`cut`](#cut) from string values *(recursively on all possible items)*. If the removed value is a more complex type (e.g. slice / map), then it is removed on a per key bases as appropriate in a similar manner to Django's subtract function.

```django
<!-- Integers: .Test is 10 -->
{{ subtract 5 .Test }} <!-- 5 -->
{{ subtract -5 .Test }} <!-- 15 -->
{{ subtract 2.5 .Test }} <!-- 8 -->
{{ subtract 2.4 .Test }} <!-- 8 -->
{{ subtract 2.6 .Test }} <!-- 7 -->
{{ subtract "5" .Test }} <!-- 5 -->
{{ subtract "5.5" .Test }} <!-- 5 -->
{{ subtract "string" .Test }} <!-- 10 -->
{{ subtract .Test "string" }} <!-- string -->

<!-- Floats: .Test is 10.0 -->
{{ subtract 5 .Test }} <!-- 5.0 -->
{{ subtract -5 .Test }} <!-- 15.0 -->
{{ subtract 2.5 .Test }} <!-- 7.5 -->
{{ subtract 2.4 .Test }} <!-- 7.6 -->
{{ subtract 2.6 .Test }} <!-- 7.4 -->
{{ subtract "5" .Test }} <!-- 5.0 -->
{{ subtract "5.5" .Test }} <!-- 4.5 -->
{{ subtract "string" .Test }} <!-- 10.0 -->
{{ subtract .Test "string" }} <!-- string -->

<!-- Strings: .Test is "test string" -->
{{ subtract "test" .Test }} <!--  string -->
{{ subtract " string" .Test }} <!-- test -->

<!-- Recursive Slices / Arrays: .Test is [1, 2, 3] -->
{{ subtract 5 .Test }} <!-- [-4, -3, -2] (see Integers for examples) -->
{{ subtract "string" .Test }} <!-- [1, 2, 3] -->
<!-- Recursive Slices / Arrays: .Test is ["string", "test"] -->
{{ subtract "test" .Test }} <!-- ["string", ""] -->

<!-- REMOVE - slices / arrays must be of the same type as added element -->
<!-- Slices / Arrays: .Test is [1, 2, 3], .Remove is [2, 3, 4] -->
{{ subtract .Remove .Test }} <!-- [1] -->
<!-- Slices / Arrays: .Test is ["string value", "slice"], .Remove is ["slice", "string"] -->
{{ subtract .Remove .Test }} <!-- ["string value"] -->

<!-- Recursive Maps: .Test is ["first": 1, "second": 2] -->
{{ subtract 5 .Test }} <!-- ["first": -4, "second": -3] (see Integers for examples) -->
{{ subtract "second" .Test }} <!-- ["first": 1, "second": 2] -->
<!-- Recursive Maps: .Test is ["first": "one", "second": "two"] -->
{{ subtract "two" .Test }} <!-- ["first": "one", "second": ""] -->

<!-- REMOVE - map values must be of the same type as added element -->
<!-- Maps: .Test is ["first": 1, "second": 2], .Remove is ["second": 2, "third": 3] -->
{{ subtract .Remove .Test }} <!-- ["first": 1] -->
```

## `suffix`

```go
func suffix[T any](suffix ...string, value T) T
```

Suffixes all strings within `value` with all strings within `suffix`. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ suffix " suffix" "value" }} <!-- value suffix -->
{{ suffix " suffix1" " suffix2" "value" }} <!-- value suffix1 suffix2 -->

<!-- Slices / Arrays: .Test is ["string1", "string2"] -->
{{ suffix " suffix" .Test }}
<!-- ["string1 suffix", "string2 suffix"] -->

<!-- Maps: .Test is [1: "string1", 2: "string2"] -->
{{ suffix " suffix" .Test }}
<!-- [1: "string1 suffix", 2: "string2 suffix"] -->
```

## `time`

```go
func time(input ...any) string
```

Parses dates to return a simple time string (by default: "HH:MM"). Supports Go, Python and PHP formatting standards *(for input / output formatting)*. The last parameter is always the date input.

It can accept various parameter combinations:

```django
<!-- Current date and default output format -->
{{ time }}
<!-- 12:30 -->

<!-- Passed `time.Time` object (default output format) -->
{{ time .Time }}
<!-- 11:30 -->

<!-- Passed Unix time (default output format) -->
{{ datetime 1556015421 }}
<!-- 11:30 -->

<!-- Current time and Go formatting string -->
{{ time "15:04:05" }}
<!-- 12:30:45 -->

<!-- Current time and PHP formatting string -->
{{ time "H:i:s" }}
<!-- 12:30:45 -->

<!-- Passed Go formatting string and `time.Time` object -->
{{ time "15:04:05" .Time }}
<!-- 11:30:12 -->

<!-- Passed Go formatting string and Unix time -->
{{ datetime "15:04:05" 1556015421 }}
<!-- 11:30:12 -->

<!-- Passed Python formatting string and `time.RFC3339` string -->
{{ time "%H:%M:%S" "2020-02-15T11:30:12Z" }}
<!-- 11:30:12 -->

<!-- Passed defined output formatting string (MYSQL), defined input layout (ATOM) and matching date string -->
{{ time "MYSQL", "ATOM" "2020-02-15T11:30:12Z00:00" }}
<!-- 2020-02-15 11:30:12 -->
```

Date and time functions support various pre-defined formats for simplicity, see [`date`](#date).

## `timesince`

```go
func timesince(t time.Time) map[string]int
```

Calculates the approximate duration since the `time.Time` value. The map of integers contains the keys: `years`, `weeks`, `days`, `hours`, `minutes`, `seconds`.

## `timeuntil`

```go
func timeuntil(t time.Time) map[string]int
```

Calculates the approximate duration until the `time.Time` value. The map of integers contains the keys: `years`, `weeks`, `days`, `hours`, `minutes`, `seconds`.

## `title`

```go
func title[T any](value T) T
```

Converts string text to title case. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ title "This string. Has TWO sentences." }}
<!-- This String. Has Two Sentences. -->
```

## `trim`

```go
func trim[T any](remove string, value T) T
```

Removes the passed characters from the ends of string values. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ trim " " "  This string. Has TWO sentences.  " }}
<!-- This string. Has TWO sentences. -->

{{ trim " .hT" "This string. Has TWO sentences.  " }}
<!-- is string. Has TWO sentences -->

{{ trim "h." "  This string. Has TWO sentences.  " }}
<!--   This string. Has TWO sentences.   -->
```

## `truncate`

```go
func truncate[T any](length int, value T) T
```

Truncates strings to a certain number of characters. Is multi-byte safe and HTML aware. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ truncate 5 "hello world"}}
<!-- hello -->

{{ truncate 5 `<a href="#test"><strong>hello world</strong></a>` }}
<!-- <a href="#test"><strong>hello</strong></a> -->
```

## `truncatewords`

```go
func truncatewords[T any](length int, value T) T
```

Truncates strings to a certain number of words. Is multi-byte safe and HTML aware. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ truncatewords 3 "hello world how are you?"}}
<!-- hello world how -->

{{ truncatewords 3 `<a href="#test"><strong>hello world</strong></a> how are you?` }}
<!-- <a href="#test"><strong>hello world</strong></a> how -->
```

## `type`

```go
func type(value any) string
```

Returns the value's `reflect.Type` as a string. Mainly useful for testing.

```django
{{ type 3.14159 }} <!-- float64 -->
{{ type "test" }} <!-- string -->

<!-- Slice: .Test is []int{1, 2, 3} -->
{{ type .Test }} <!-- []int -->

<!-- Array: .Test is [3]int{1, 2, 3} -->
{{ type .Test }} <!-- [3]int -->

<!-- Map: .Test is map[int]string{2: "second", 1: "first"} -->
{{ type .Test }} <!-- map[int]string -->
```

## `ul`

```go
func ul(value any) string
```

Converts slices, arrays or maps into an HTML unordered list.

Other data types will just return a string representation of themselves.

```django
<!-- .Test is slice: [1, 2, 3] -->
{{ ul .Test }}
<!-- produces: -->
<ul>
	<li>1</li>
	<li>2</li>
	<li>3</li>
</ul>

<!-- .Test is map: ["first": "first-content", "second": "second-content"] -->
{{ ul .Test }}
<!-- produces: -->
<ul>
	<li>first-content</li>
	<li>second-content</li>
</ul>

<!-- .Test is map: ["first": ["slice", "one"], "second": ["slice", "two"]] -->
{{ ul .Test }}
<!-- produces: -->
<ul>
	<li>
		<ul>
			<li>slice</li>
			<li>one</li>
		</ul>
	</li>
	<li>
		<ul>
			<li>slice</li>
			<li>two</li>
		</ul>
	</li>
</ul>
```

## `upper`

```go
func upper[T any](value T) T
```

Converts string text to upper case. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ upper "This string. Has TWO sentences." }}
<!-- THIS STRING. HAS TWO SENTENCES. -->
```

## `urldecode`

```go
func urlDecode[T any](url T) T
```

Converts URL character-entity equivalents back into their literal, URL-unsafe forms. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ urldecode "%21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D" }}
<!-- ! * ' ( ) ; : @ & = + $ , / ? % # [ ] -->
```

## `urlencode`

```go
func urlEncode[T any](url T) T
```

Converts URL-unsafe characters into character-entity equivalents to allow the string to be used as part of a URL. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ urlencode "! * ' ( ) ; : @ & = + $ , / ? % # [ ]" }}
<!-- %21 %2A %27 %28 %29 %3B %3A %40 %26 %3D %2B %24 %2C %2F %3F %25 %23 %5B %5D -->
```

*(N.B. If [`OverloadFunctions()`](README.md#overloading-texttemplate-functions) has been used, this function will also replace the built in [`urlquery`](BASICS.md#urlquery) function)*

## `uuid`

```go
func uuid() string
```

Generates a unique ID using the `github.com/google/uuid` package based on RFC 4122 and DCE 1.1.

```django
{{ uuid }}
<!-- f47ac10b-58cc-0372-8567-0e02b2c3d479 -->
```

## `values`

```go
func values(value any) []any
```

Returns the values of a slice / array / map / struct.

```django
{{ values .Slice }}
<!-- []int{12, 44, 128} -->

{{ values .Map }}
<!-- []string{"value1", "value2", "value3"} -->
```

## `wordcount`

```go
func wordcount(value string) int
```

Counts the number of words (excluding HTML, numbers and special characters) in a string. Only works on strings and returns 0 for any other value.

```django
{{ wordcount "hello world"}}
<!-- 2 -->

{{ wordcount `" <a href="#test"><strong>hello world</strong></a> how 12 " are " you ?` }}
<!-- 5 -->
```

## `wrap`

```go
func wrap[T any](prefix string, suffix string, value T) T
```

Wraps all strings within `value` with a prefix and suffix. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

```django
{{ wrap "prefix " " suffix" "value" }} <!-- prefix value suffix -->

<!-- Slices / Arrays: .Test is ["string1", "string2"] -->
{{ wrap "prefix " " suffix" .Test }}
<!-- ["prefix string1 suffix", "prefix string2 suffix"] -->

<!-- Maps: .Test is [1: "string1", 2: "string2"] -->
{{ wrap "prefix " " suffix" .Test }}
<!-- [1: "prefix string1 suffix", 2: "prefix string2 suffix"] -->
```

## `year`

```go
func year(times nil|time.Time) int
```

Returns an integer year from a `time.Time` input, or the current year if no time is provided.

```django
<!-- Current date -->
{{ year }}
<!-- 2022 -->

<!-- Passed `time.Time` object -->
{{ year .Time }}
<!-- 2020 -->
```

## `yesno`

```go
func yesno(values ...any) string
```

Returns "Yes" for true values, "No" for false values and "Maybe" for empty values (`maybe` defaults to "No" unless maybe is specifically defined).

Return string options may be customised. If numeric arguments are used, it treats numeric zero as "No", positive numbers as "Yes" and negative numbers as "Maybe"
If string, slice, array or map arguments are used, it treats empty as "Maybe", and populated as "Yes".

Examples of use:

```django
<!-- Uses the default "Yes" / "No" returns -->
{{ yesno 1 }} 
<!-- Yes -->

{{ yesno 0 }} 
<!-- No -->

<!-- Customises the string used for "Yes" -->
{{ yesno "Yep" 1 }}
<!-- Yep -->

{{ yesno "Yep" 0 }}
<!-- No -->

<!-- Customises the strings used for "Yes" and "No" -->
{{ yesno "Yep" "Nope" 1 }}
<!-- Yep -->

{{ yesno "Yep" "Nope" 0 }}
<!-- Nope -->

{{ yesno "Yep" "Nope" -1 }}
<!-- Nope -->

<!-- Customises the strings used for "Yes" and "No" and enables "Maybe" -->
{{ yesno "Yep" "Nope" "Perhaps" 1 }}
<!-- Yep -->

{{ yesno "Yep" "Nope" "Perhaps" 0 }}
<!-- Nope -->

{{ yesno "Yep" "Nope" "Perhaps" -1 }}
<!-- Perhaps -->

<!-- Values do not need to be integers -->
{{ yesno "Yarp" "Narp" "Larp" [1, 2, 3] }}
<!-- Yarp -->

{{ yesno "Yarp" "Narp" "Larp" [] }}
<!-- Larp -->

{{ yesno "Yarp" "Narp" [] }}
<!-- Narp -->
```