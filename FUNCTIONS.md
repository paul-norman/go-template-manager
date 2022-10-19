# Functions in `templateManager`

All functions in `templateManager` accept their principle argument last to allow simple chaining. *Efforts have been made to output clear errors and return suitable empty values rather than cause panics (a problem in several `text/template` functions)*.

Contents: [`add`](#add), [`capfirst`](#capfirst), [`collection`](#collection), [`contains`](#contains), [`cut`](#cut), [`date`](#date), [`datetime`](#datetime), [`default`](#default), [`divide`](#divide), [`divisibleby`](#divisibleby), [`dl`](#dl), [`first`](#first), [`firstof`](#firstof), [`formattime`](#formattime), [`htmldecode`](#htmldecode), [`htmlencode`](#htmlencode), [`join`](#join), [`jsondecode`](#jsondecode), [`jsonencode`](#jsonencode), [`key`](#key), [`last`](#last), [`localtime`](#localtime), [`lower`](#lower), [`ltrim`](#ltrim), [`mktime`](#mktime), [`multiply`](#multiply), [`nl2br`](#nl2br), [`now`](#now), [`ol`](#ol), [`ordinal`](#ordinal), [`paragraph`](#paragraph), [`pluralise`](#pluralise), [`prefix`](#prefix), [`random`](#random), [`regexp`](#regexp), [`regexpreplace`](#regexpreplace), [`replace`](#replace), [`rtrim`](#rtrim), [`split`](#split), [`striptags`](#striptags), [`subtract`](#subtract), [`suffix`](#suffix), [`time`](#time), [`timesince`](#timesince), [`timeuntil`](#timeuntil), [`title`](#title), [`trim`](#trim), [`truncate`](#truncate), [`truncatewords`](#truncatewords), [`ul`](#ul), [`upper`](#upper), [`urldecode`](#urldecode), [`urlencode`](#urlencode), [`wordcount`](#wordcount), [`wrap`](#wrap), [`year`](#year), [`yesno`](#yesno)

## `add`

```go
func add[T any](add T, to T) T
```

Adds a value to the existing item. If the added value is a simple numeric type this will be treated as a simple addition *(recursively on all possible items)*. If the added value is a string, it will be appended to string values *(recursively on all possible items)*. If the added value is a more complex type (e.g. slice / map), then it is appended / merged as appropriate in a similar manner to Django's add function.

Returns new variable of the original `to` data type.

## `capfirst`

```go
func capfirst[T any](value T) T
```

Capitalises the first letter of strings. Does not alter any other letters. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `collection`

```go
func collection(key1 string, value1 any, key2 string, value2 any, ...) map[string]any
```

Allows any number of pairs of keys / values to be grouped together into a map. This map can be passed to templates / blocks.

```django
{{ template "partials/requiresMultipleVars.html" collection "var1" .Var1 "var2" .Var2 }}
```

## `contains`

```go
func contains(find any, within any) bool
```

Determines whether the `find` value is contained in the `within` value. The `find` value can act on strings, slices, arrays and maps, but contained types must match.

Returns a boolean value (always false for incompatible types).

## `cut`

```go
func cut[T any](remove string, from T) T
```

Will `remove` a string value that is contained in the `from` value. If `from` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `from` data type.

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

Will return the second `test` value if it is not empty, else return the `def` value

## `divide`

```go
func divide[D int|float64, T any](divisor D, value T) T
```

Divides the `value` by the `divisor`. If `value` is a slice, array or map it will apply this conversion to any numeric elements that they contain. All values are first converted to floats, the operation is performed and then any **rounding is applied as necessary to return the item to its original type**.

## `divisibleby`

```go
func divisibleby[T any](divisor int, value T) bool
```

Determines if the `value` is exactly divisible by the `divisor`.

## `dl`

```go
func dl(value any) string
```

Converts slices, arrays or maps into an HTML definition list. For maps this will use the keys as the dt elements.

Other data types will just return a string representation of themselves.

## `first`

```go
func first(value any) any
```

Gets the first value from slices / arrays or the first word from strings. **All other data types return an empty variable.**

## `firstof`

```go
func firstof(values ...any) any
```

Accepts any number of values and returns the first one of them that exists and is not empty. If none are found it returns an empty value.

## `formattime`

```go
func formattime(format string, t time.Time) string
```

Formats a time.Time object for display.

```django
{{ now | formattime "d/m/y H:i:s" }}
```

Date and time functions support various pre-defined formats for simplicity, see [`date`](#date).

## `htmldecode`

```go
func htmldecode[T any](value T) T
```

Converts HTML character-entity equivalents back into their literal, usable forms.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `htmlencode`

```go
func htmlencode[T any](value T) T
```

Converts literal HTML special characters into safe, character-entity equivalents.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `join`

```go
func join(separator string, values any) string
```

Joins slice or map `values` together as a string spaced by the `separator`. Strings are returned unaltered and numeric types are coerced into strings.

## `jsondecode`

```go
func jsondecode(value any) any
```

Decodes any JSON string to an `interface{}`. This usually produces a type: `map[string]any`, but may result in other types (e.g. `[]any`) or simple types (e.g. `bool`, `string`) for trivial data sources. ALL numbers will be `float64`. Use only with caution / testing.

## `jsonencode`

```go
func jsonencode(value any) string
```

Encodes any value to a JSON string.

## `key`

```go
func key(input ...any) any
```

Very similar to the `text/template` [`slice`](BASICS.md#general-utility-functions) function, `key` accepts any number of nested keys and returns the result of indexing its **final argument** by them. For strings this returns individual letters. The indexed item must be a string, map, slice, or array.

## `last`

```go
func last(value any) any
```

Gets the last value from slices / arrays or the last word from strings. **All other data types return an empty variable.**

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

## `ltrim`

```go
func ltrim[T any](remove string, value T) T
```

Removes the passed characters from the left end of string values. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

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

## `nl2br`

```go
func nl2br[T any](value T) T
```

Replaces all instances of `\n` (new line) with instances of `<br>` within `value`. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain. [`paragraph`](#paragraph) will perform a similar task in a more intelligent manner.

Returns new variable of the original `value` data type.

## `now`

```go
func now() time.Time
```

Returns the current `time.Time` value.

## `ol`

```go
func ol(value any) string
```

Converts slices, arrays or maps into an HTML ordered list.

Other data types will just return a string representation of themselves.

## `ordinal`

```go
func ordinal[T int|float64|string](value T) string
```

Suffixes a number with the correct English ordinal. If `value` is not numeric or a valid numeric string, an empty string is returned.

## `paragraph`

```go
func paragraph[T any](value T) T
```

Replaces all string instances of `\n+` (multiple new lines) with paragraph tags (`</p><p>`) and instances of `\n` (new line) with instances of `<br>` within `value`. Finally wraps the string in paragraph tags. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

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
func prefix[T any](prefix string, value T) T
```

Prefixes all strings within `value` with `prefix`. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `random`

```go
func random(...int) int
```

Generates positive random numbers.

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
```

## `regexp`

```go
func regexp(find string, value string) [][]string
```

Finds all instances of `find` regexp within `value` using [`regexp.FindAllStringSubmatch`](https://pkg.go.dev/regexp#Regexp.FindAllStringSubmatch). It only acts on strings, returning an empty string slice for any other values.

## `regexpreplace`

```go
func regexpreplace[T any](find string, replace string, value T) T
```

Replaces all instances of `find` regexp with instances of `replace` within `value` using [`regexp.ReplaceAllString`](https://pkg.go.dev/regexp#Regexp.ReplaceAllString). If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `replace`

```go
func replaceAll[T any](find string, replace string, value T) T
```

Replaces all instances of `find` with instances of `replace` within `value`.
If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `rtrim`

```go
func rtrim[T any](remove string, value T) T
```

Removes the passed characters from the right end of string values. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `split`

```go
func split(separator string, value string) []string
```

Splits strings on the `separator` value and returns a slice of the pieces. It only works on strings and returns an empty slice for all other data types.

## `striptags`

```go
func stripTags[T any](value T) T
```

Strips HTML tags from strings. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `subtract`

```go
func subtract[T any](subtract T, from T) T
```

Subtracts a value from the existing item. If the subtracted value is a simple numeric type this will be treated as a simple maths *(recursively on all possible items)*. If the removed value is a string, it will be [`cut`](#cut) from string values *(recursively on all possible items)*. If the removed value is a more complex type (e.g. slice / map), then it is removed on a per key bases as appropriate in a similar manner to Django's subtract function.

## `suffix`

```go
func suffix[T any](suffix string, value T) T
```

Suffixes all strings within `value` with `suffix`. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

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

## `trim`

```go
func trim[T any](remove string, value T) T
```

Removes the passed characters from the ends of string values. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `truncate`

```go
func truncate[T any](length int, value T) T
```

Truncates strings to a certain number of characters. Is multi-byte safe and HTML aware. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `truncatewords`

```go
func truncatewords[T any](length int, value T) T
```

Truncates strings to a certain number of words. Is multi-byte safe and HTML aware. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `ul`

```go
func ul(value any) string
```

Converts slices, arrays or maps into an HTML unordered list.

Other data types will just return a string representation of themselves.

## `upper`

```go
func upper[T any](value T) T
```

Converts string text to upper case. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `urldecode`

```go
func urlDecode[T any](url T) T
```

Converts URL character-entity equivalents back into their literal, URL-unsafe forms. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `urlencode`

```go
func urlEncode[T any](url T) T
```

Converts URL-unsafe characters into character-entity equivalents to allow the string to be used as part of a URL. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

## `wordcount`

```go
func wordcount(value string) int
```

Counts the number of words (excluding HTML, numbers and special characters) in a string. Only works on strings and returns 0 for any other value.

## `wrap`

```go
func wrap[T any](prefix string, suffix string, value T) T
```

Wraps all strings within `value` with a prefix and suffix. If `value` is a slice, array or map it will apply this conversion to any string elements that they contain.

Returns new variable of the original `value` data type.

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