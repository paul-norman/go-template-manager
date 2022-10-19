# Basic Usage of `text/templates` Package

Go has a built in templating engine in the standard library, the `text/template` package ([docs](https://pkg.go.dev/text/template)).

This short guide does not cover how to configure these templates *(what `templateManager` simplifies)*, but rather just how to use them.

## Contents

- [Basic Syntax](#basic-syntax)
  - [Whitespace](#whitespace)
  - [Comments](#comments)
  - [Output in Delimiters](#output-in-delimiters)
  - [Local Variables](#local-variables)
  - [Global Variables](#global-variables)
  - [if / else](#if--else)
  - [with](#with)
  - [range](#range)
  - [template](#template)
  - [define](#define)
  - [block](#block)
- [Built-in Functions](#built-in-functions)
  - [Boolean Operators](#boolean-operators)
  - [Equality Operators](#equality-operators)
  - [Relational Operators](#relational-operators)
  - [Formatting Functions](#formatting-functions)
  - [Text Escaping Functions](#text-escaping-functions)
  - [General Utility Functions](#general-utility-functions)
- [Examples](#examples)

## Basic Syntax

There are only 10 keywords in the package:

```
block, break, define, else, end, if, range, nil, template, with
```

and the special `.` (dot) operator which allows access to variables in ay given "context" *(more on that later)*.

To use any of these keywords, special delimiters are wrapped around them:

```django
{{ template "layout.html" . }}
```

Some statements also require an `{{ end }}` tag:

```django
{{ define "name" }} content {{ end }}
```

### Whitespace

It is often clearer to write statements such as the above `define` statement with space before and after the content. However, this will result in physical whitespace characters being added to the content *(often undesirably)*.

To trigger stripping of whitespace, the delimiters may have a hyphen attached to them to trigger that whitespace will be removed from that direction:

```django
some {{ block "name1" }} content {{ end }} here
<!-- some  content  here-->

some {{ block "name2" -}} ltrimmed content {{ end }} here
<!-- some ltrimmed content  here-->

some {{ block "name2" }} rtrimmed content {{- end }} here
<!-- some  ltrimmed content here-->

some {{ block "name3" -}} trimmed content {{- end }} here
<!-- some trimmed content here-->

some {{- block "name4" -}} squashed content {{- end -}} here
<!-- somesquashed contenthere-->
```

### Comments

Code comments are useful, but when working with HTML, native comments are visible in the source code. For this situation, it is possible to add comments to delimiters such that they do not appear in the HTML source:

```django
{{/* A comment */}}
```

It is sometimes a good idea to strip whitespace from around them too as they have likely been inserted between page tags:

```django
{{- /* A comment */ -}}
```

### Output in Delimiters

Literal values (e.g. ints or strings) may be placed inside delimiters so as to allow whitespace stripping or for their value to be used as an input to functions.

```django
{{ "hello" }} <!-- hello -->
```

Or all whitespace before and after this word can be removed:

```django
hello {{- "world" -}} ! <!-- helloworld! -->
```

Whether it is a literal or a variable, the output of delimited tags is equivalent to writing:

```go
fmt.Print(variable_name)
```

 *(i.e. can output anything)*.

### Local Variables

Local variables may be assigned / reassigned within a single template file.

```django
{{ $test := "World" }}
Hello {{ $test }} <!-- Hello World -->
{{ $test = "Bob" }}
Hello {{ $test }} <!-- Hello Bob -->
```

These are limited to the current file context unless specifically passed in to another template / block.

### Global Variables

Global variables are passed in to the template when it is "Executed" (rendered). At the root level of the template the `.` (dot) character represents **ALL** of these variables. So if a string is passed in as the only parameter, `{{.}}` will represent a string. If a map of strings (or a struct) is passed in, then `{{.}}` will represent the whole object and their keys will be available after the dot, for example: `{{ .Key1 }}`.

*(N.B. Inside loops, or in nested templates, the dot can be remapped to the root of that "context". e.g. the value of the looped variable will be reassigned to `{{.}}`)*

### `if` / `else`

In Go, if statements require a strict boolean evaluation to work. However within these templates, the condition is whether the value is empty or populated (i.e. not the default for a give variable type). For example:

```django
{{ if .String }}
	<p>String is: {{ .String }}</p>
{{ end }}

{{ if .Int }}
	<p>Int is not zero</p>
{{ else }}
	<p>Int is zero</p>
{{ end }}
```

Approximately, something like this happens for `if` (and `with`) statements:

```go
true                      => true
false                     => false

"string"                  => true
""                        => false

1                         => true
0                         => false

1.0                       => true
0.0                       => false

[]int{5}                  => true
[]int{}                   => false

map[string]int{"Test": 1} => true
map[string]int{}          => false

struct{}                  => true
nil                       => false
```

### `with`

`with` is very similar to `if` in that checks to see if a variable is empty and if not executes its block. However, `with` reassigns the `.` (dot) to that variable's value for the duration of the block.

```django
{{ with .Variable }}
	Dot is now the value of ".Variable" {{.}}
{{ end }}
```

Like `if`, this can also be paired with an optional `else` statement:

```django
{{ with .Variable }}
	{{.}} now means ".Variable" 
{{ else }}
	Dot is still as it was outside this loop
{{ end }}
```

`with` blocks may be nested, so the context of `{{.}}` may change several times!

If `.` (dot) has been reassigned, the global variables can always be accessed using the `$.` prefix:

```django
{{ with .Variable }}
	Dot is now the value of ".Variable" {{.}} and so is {{ $.Variable }}
{{ end }}
```

### `range`

If the `range` keyword is passed an array, slice, map, or channel it will iterate over it, redefining the `.` (dot) to its current value. It may also use an `else` block that is executed if the value is empty:

```django
<ul>
{{ range .Slice }}
	<li>{{ . }}</li>
{{ end }}
</ul>

{{ range .Slice }}
	<ul>
	{{ range .Slice }}
		<li>{{ . }}</li>
	{{ end }}
	</ul>
{{ else }}
	<p>.Slice is empty</p>
{{ end }}
```
It is also possible to get the index alongside the value from the `range` function:

```django
<dl>
{{ range $index, $value := .Map }}
	<dt>{{ $index }}</dt>
	<dd>{{ $value }}</dd>
{{ end }}
</dl>
```

It's possible to break a range loop early by calling `{{ break }}` within the loop. If there are nested ranges, this only stops the first loop that directly contains it. 

### `template`

The `template` block is an instruction to execute the named template and return its data to this position. It may accept either a single argument or two arguments.

```django
{{ template "template_name" }}
```

This renders the "template_name" template without any variables passed to it (and it must not require any). If the template requires some / all of the original input data it's possible to pass it through to the nested template.

This renders the "template_name" template with ALL of the current context's data passed to it:

```django
{{ template "template_name" . }}
```
However, any single variable can be passed in if that's all that's needed:

```django
{{ template "template_name" .MyMap }}
```

### `define`

The `define` block creates a new template that can then be called using the `template` (or `block`) commands.

```django
{{ define "template_name" }}
	This is the template content
{{ end }}
```

### `block`

The `block` block is shorthand for creating a `define` block and then immediately rendering it using the `template` command. If the content is declared elsewhere in a `define` block, that content overrides it.

The following would immediately run the "template_name" template:

```django
{{ block "template_name" }}
	This is the template content
{{ end }}
```

But it could equally be overwritten using any define block:

```django
{{ define "template_name" }}
	This content will be output
{{ end }}

{{ block "template_name" }}
	This content will not be
{{ end }}
```

## Built-in Functions

To support the logic of the templates, there are a small selection of useful functions. The format of these is that the functions name is written, followed by the arguments that it accepts:

```go
{{ function_name "arg1" "arg2" }}
```

Functions can be chained together using the `|` (pipe) operator. The result of the prior function (or literal statement) is passed to the next in the chain as the **final** argument:

```go
{{ "input" | func1 "arg1" "arg2" | func2 }}
```

This is the equivalent of running:

```go
result1 := func1("arg1", "arg2", "input")
result2 := func2(result1)
fmt.Print(result2)
```

All template functions are written in this manner.

### Boolean Operators

There are no simple boolean operators such as `&&`, `||` or `!`, so they have become functions:

#### `and`
```go 
if and x y // true if both x and y are true
```

#### `or`
```go 
if or x y // true if either x or y is true
```

#### `not`
```go 
if not x // true if either x is false
```

### Equality Operators

There are no simple equality operators such as `==` and `!=`, so they have become functions:

#### `eq`

```go 
if eq x y // if x == y
```

#### `ne`

```go 
if ne x y // if x != y
```

### Relational Operators

There are no simple relational operators such as `>`, `<`, `>=` and `<=` so they have become functions:

#### `lt`

```go 
if lt x y // if x < y
```

#### `gt`

```go 
if gt x y // if x > y
```

#### `le`

```go 
if le x y // if x <= y
```

#### `ge`

```go
if ge x y // if x >= y
```

### Formatting Functions

These functions allow effective and safe output of all data types.

#### `print` 

Alias for `fmt.Sprint`

#### `printf`

Alias for `fmt.Sprintf`

#### `println`

Alias for `fmt.Sprintln`

### Text Escaping Functions

In order to make certain strings safe for use on the web, these functions are available: *(Oddly, no reverse functions are provided?)*

#### `html`

Converts literal HTML special characters into safe, character-entity equivalents.

```django
{{ html "<p>HTML <em>string</em></p>" }}
<!-- &lt;p&gt;HTML &lt;em&gt;string&lt;/em&gt;&lt;/p&gt; -->
```

#### `js`

Returns the escaped JavaScript equivalent of the textual representation of its arguments.

```django
{{ js "<script>let test = 1;</script>" }}
<!-- \u003Cscript\u003Elet test \u003D 1;\u003C/script\u003E -->
```

#### `urlquery`

Converts URL-unsafe characters into character-entity equivalents to allow the string to be used as part of a URL.

```django
{{ urlquery "https://www.example.com/?v=cat+and+dogs" }}
<!-- https%3A%2F%2Fwww.example.com%2F%3Fv%3Dcat%2Band%2Bdogs -->
```

### General Utility Functions

#### `len`

Accepts a single argument and returns the integer length of its input.

```django
{{ len "test" }}
<!-- 4 -->
```

#### `slice`

Returns the result of slicing its first argument by the remaining arguments. The first argument must be a string, slice, or array.

```django
<!-- Given .Array = [1, 2, 3, 4] -->

{{ slice .Array }}
<!-- .Array[:] => [1, 2, 3, 4] -->

{{ slice .Array 1 }}
<!-- .Array[1:] => [2, 3, 4] -->

{{ slice .Array 1 3 }}
<!-- .Array[1:3] => [2, 3] -->

{{ slice .Array 1 3 6 }}
<!-- .Array[1:3:6] => [2, 3] (capacity 5) -->
```

*Chainability of this function is poor as it receives its main argument as its first argument (rather than its last)*

*Can cause panics, use with caution!*

#### `index`

Returns the result of "indexing" the first argument by the keys in the later arguments. The first argument must be a string, slice, or array. If used on strings, it returns the character code.

```django
<!-- Given .Map = {"test": [[1, 2, 3], [4, 5, 6], [7, 8, 9]]} -->

{{ index .Map }}
<!-- .Map => {"test": [[1, 2, 3], [4, 5, 6], [7, 8, 9]]} -->

{{ index .Map "test" }}
<!-- .Map["test"] => [[1, 2, 3], [4, 5, 6], [7, 8, 9]] -->

{{ index .Map "test" 1 }}
<!-- .Map["test"][1] => [4, 5, 6] -->

{{ index .Map "test" 1 2 }}
<!-- .Map["test"][1][2] => 6 -->

{{ index .Map "test" 1 2 3 }}
<!-- .Map["test"][1][2][3] => PANIC! -->

{{ index .Map "missing" }}
<!-- .Map["missing"] => "" (empty string) -->
```

*Chainability of this function is poor as it receives its main argument as its first argument (rather than its last)*

*Can cause panics, use with caution!*

#### `call`

Calls the function in the first argument returning its value. Passes all subsequent arguments to it. First argument must be a function.
 
```django
{{ call .Func "arg1" "arg2" }}
<!-- .Func("arg1" "arg2")--->
```

## Examples