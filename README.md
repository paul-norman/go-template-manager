# TemplateManager

Package `templateManager` simplifies the use of Go's standard library `text/template` for use with HTML templates.

It automates the process of choosing which files to group together for parsing (creating bundles), and builds a store of each entry template file complete with all of its dependencies. It also allows simple variables to be globally defined in the templates themselves and the use of simple components.

## Contents

- [Installation](#installation)
- [Main Features](#main-features)
- [Basic Usage](#basic-usage)
- [Customisation Options](#customisation-options)
- [Setting Variables](#setting-variables)
- [Creating Functions](#creating-functions)
- [Built-in Functions](#built-in-functions)
- [Error Handling](#error-handling)
- [Simple Example](#simple-example)
- [Integrations](#integrations)

## Installation

Install / update `templateManager` using go get:

```
go get -u github.com/paul-norman/go-template-manager
```

or import it as:

```go
import "github.com/paul-norman/go-template-manager"
```

## Main Features

`templateManager` generally just simplifies the usage of the `text/templates` package. A basic usage guide to these [Go Templates](https://pkg.go.dev/text/template) is provided [here](BASICS.md).

### `extends` keyword

`templateManager` adds support for an `extends` keyword to automate bundling. It allows entry files to specify which layout that they "extend" directly in the template without manually specifying all bundle files in Go *(nor manually including files via the built-in `template` function)*:

```go
{{ extends "layouts/main.html" }}
```

`templateManager` will then follow all instances of the `template` tag in these two files until all necessary files for the entry template are known to the bundle. This ensures that only the correct / required blocks exist in any bundle.

### Variables in Templates

`templateManager` adds a second new keyword, `var`, to allow VERY basic variables to be defined within the templates themselves. It also allows these variables to be overridden via a simple hierarchy based on load depth:

```go
{{ var "int1" }} 123 {{ end }}
```
*(see [variables](#setting-variables) section for details)*

### Components

`templateManager` allows the user to create their own HTML-style web components to simplify their source code. For example:

```html
<Youtube id="QH2-TGUlwu4">
```

or

```html
<Tabset>
	<x-Tab>Tab 1</x-Tab>
	<x-TabContent>Tab 1 Content</x-TabContent>

	<x-Tab>Tab 2</x-Tab>
	<x-TabContent Checked=1>
		Tab 2 Content
		<Youtube id="QH2-TGUlwu4">
	</x-TabContent>
</Tabset>
```

*(These require a more in-depth explanation, so have been moved to their own file - see [components](COMPONENTS.md) for details)*

### Convenience Functions

`templateManager` comes with a small set of convenience functions which may be used or removed.

*(see [functions](FUNCTIONS.md) file for details)*

## Basic Usage

`templateManager` requires initialisation. At a minimum it requires you to tell it where the templates are kept, and what file type the template files are:

```go
var tm = templateManager.Init("templates", ".html")
```

During parsing, `templateManager` only needs to know about the location of the "entry" template files as it will scan each for its dependencies. It is therefore best that it does not know about layout files or partials and these should be kept separate and excluded from parsing:

```go
tm.ExcludeDirectories([]string{"layouts", "partials"})
```

Once these variables have been configured, the template parsing can be triggered to build the cache.

```go
tm.Parse()
```

To run a template, use the `Render()` method:

```go
tm.Render("test.html", TM.Params{"Title": "Test"}, ioWriter)
```

## Customisation Options

All customisation options are chainable for neat declaration.

### Delimiters

The default delimiters set in the `text/template` package are `{{` and `}}` which wrap all commands. These can be customised if desired before templates are parsed:

```go
tm.Delimiters("{%", "%}")
```

### Debugging

During development it's often useful to see what is happening. Enabling debug mode outputs a log showing what happens upon parsing *(and forces warnings and errors to be shown in the log)*:

```go
tm.Debug(true)
```

### Reload Each Time

During development it's often useful to be able to change template files and not have to restart the server to test them. Enabling reload mode forces a bundle rebuild for the selected template upon each `Render()` method call:

```go
tm.Reload(true)
```

*(N.B. this does not work with an embedded file system as the changes are not picked up until next build)*

### Excluding Directories

It is most efficient if the parser only runs over "entry" templates (that is those that will be called directly). For this reason it's best to exclude all directories *(within the designated templates folder)* from this process.

By default, two directories *("layouts" and "partials")* are already excluded.

```go
// Remove directories from parsing
tm.ExcludeDirectory("layouts")
// OR
tm.ExcludeDirectories([]string{"layouts", "partials"})

// Re-adds a previously excluded directory so that it will be parsed
tm.RemoveExcludedDirectory("layouts")
```

## Global Options

There are also several customisation options that apply globally to `templateManager` functions / use. These should be set **prior** to initialisation of the main store.

```go
import (
	"time"
	TM "github.com/paul-norman/go-template-manager"
)

// Control whether errors are written to the log
TM.SetErrors(false)

// Control whether warnings are written to the log
TM.SetWarnings(warnings bool)

// Sets the default format for the `date` function (default: d/m/Y)
// May be in Go, PHP or Python format
TM.SetDefaultDateFormat("d/M/Y")

// Sets the default format for the `datetime` function (default: d/m/Y H:i)
// May be in Go, PHP or Python format
TM.SetDefaultDatetimeFormat("d/m/Y H:i")

// Sets the default format for the `time` function (default: H:i)
// May be in Go, PHP or Python format
TM.SetDefaultTimeFormat("H:i")

// Sets the default timezone `time.Location` used by date / time functions (default: UTC)
location, _ := time.LoadLocation("ICT")
TM.SetTimezoneLocation(location)

// Sets the default timezone location used by date / time functions from a string (default: UTC)
TM.SetTimezoneLocationString("ICT")

// Sets the default timezone location used by date / time functions to a fixed numeric offset (default: UTC)
TM.SetTimezoneFixed("ICT", 7 * 60 * 60) // UTC +7
```

## Setting Variables

Variables can be set at various levels and load based on a hierarchy, with those defined in the `Render()` method having top priority, those being defined in the entry file having secondary priority and those being defined in lower templates reducing in priority based on how deeply defined they are.

Variables are still called in the code according to `text/template` syntax, i.e. `{{ .VarName }}`.

### `Render()` Variables

The `Render()` method accepts a `templateManager.Params` variable as its second argument. This type is an alias to a map of type `map[string]any`.

```go
params := templateManager.Params{
	"Title": "Test",
	"Slice": []int{1, 2, 3}
}
tm.Render("test.html", params, ioWriter)
```

Variables defined this way have the highest priority.

### Creating Variables in Templates

To keep all front end management in one place it is possible to define ***simple*** variables in templates directly. The syntax for doing so is:

```django
{{ var "varName" }} var value {{ end }}
```

These are parsed by regular expressions and turned into a "best" guess version of what they represent *(leading / trailing spaces are always trimmed)*.

**The following represents all possible types that may be declared this way:** *(no deeper nesting than is explicitly shown is currently possible)*

`simple types`
```django
<!-- creates a type int-->
{{ var "Int1" }} 1 {{ end }}
{{ var "Int2" }} -36 {{ end }}

<!-- creates a type float64 -->
{{ var "Float1" }} 3.14 {{ end }}
{{ var "Float2" }} -1.23 {{ end }}

<!-- creates a type bool -->
{{ var "Bool1" }} true {{ end }}
{{ var "Bool2" }} False {{ end }}

<!-- creates a type string -->
{{ var "String1" }} this is a string {{ end }}
{{ var "String2" }} May have "quotes" of 'various' `types` {{ end }}
{{ var "String3" }} this is an <strong>HTML</strong> string {{ end }}
```

`slices` *(no deeper nesting is supported)*
```django
<!-- creates a type []int -->
{{ var "SliceInt" }} [1, 2, 3, -4, -5, -6] {{ end }}

<!-- creates a type []float64 -->
{{ var "SliceFloat" }} [1.1, 2.2, 3.3, -4.4, -5.5, -6.6] {{ end }}

<!-- creates a type []bool -->
{{ var "SliceBool" }} [true, True, FALSE, trUE] {{ end }}

<!-- creates a type []string -->
{{ var "SliceString1" }} ["this", "is", "a", "string", "slice"] {{ end }}
{{ var "SliceString2" }} ["this \"is\"", 'a "more"', `"complex"`, "<span class=\"test\">string</span>", "slice"] {{ end }}

<!-- creates a type [][]int -->
{{ var "SliceSliceInt" }} [[1, 2], [3, -4], [-5, -6]] {{ end }}

<!-- creates a type [][]float64 -->
{{ var "SliceSliceFloat" }} [[1.1, 2.2], [3.3, -4.4], [-5.5, -6.6]] {{ end }}

<!-- creates a type [][]bool -->
{{ var "SliceSliceBool" }} [[true, True], [FALSE, trUE]] {{ end }}

<!-- creates a type [][]string -->
{{ var "SliceSliceString" }} [["this", "is"], ["a", "string"], ["slice"]] {{ end }}
```

`maps` *(Does not support nested maps nor maps of slices)*
```django
<!-- creates a type map[int]int (N.B. Quotes) -->
{{ var "MapIntInt" }} {"1": 1, "10": -10, "100": 100} {{ end }}

<!-- creates a type map[int]float64 (N.B. Quotes) -->
{{ var "MapIntFloat" }} {"1": 1.0, "10": -10.0, "100": 100.0} {{ end }}

<!-- creates a type map[int]bool (N.B. Quotes) -->
{{ var "MapIntBool" }} {"0": false, "1": true} {{ end }}

<!-- creates a type map[int]string (N.B. Quotes) -->
{{ var "MapIntString" }} {"1": "string 1", "2": "string 2"} {{ end }}

<!-- creates a type map[float64]int (N.B. Quotes) -->
{{ var "MapFloatInt" }} {"1.1": 1, "10.10": -10, "100.100": 100} {{ end }}

<!-- creates a type map[float64]float64 (N.B. Quotes) -->
{{ var "MapFloatFloat" }} {"1.0": 1.0, "10.0": -10.0, "100.0": 100.0} {{ end }}

<!-- creates a type map[float64]bool (N.B. Quotes) -->
{{ var "MapFloatBool" }} {"0.0": false, "1.0": true} {{ end }}

<!-- creates a type map[float64]string (N.B. Quotes) -->
{{ var "MapFloatString" }} {"1.0": "string 1.0", "2.0": "string 2.0"} {{ end }}

<!-- creates a type map[bool]int (N.B. Quotes) -->
{{ var "MapBoolInt" }} {"true": 1, "false": 0} {{ end }}

<!-- creates a type map[bool]float64 (N.B. Quotes) -->
{{ var "MapBoolFloat" }} {"true": 1.0, "false": 0.0} {{ end }}

<!-- creates a type map[bool]bool (N.B. Quotes) -->
{{ var "MapBoolBool" }} {"false": true, "true": false} {{ end }}

<!-- creates a type map[bool]string (N.B. Quotes) -->
{{ var "MapBoolString" }} {"false": "no", "true": "yes"} {{ end }}

<!-- creates a type map[string]int (N.B. Quotes) -->
{{ var "MapStringInt" }} {"key1": 1, "key2": 2, "key3": 3} {{ end }}

<!-- creates a type map[string]float64 (N.B. Quotes) -->
{{ var "MapStringFloat" }} {"key1": 1.0, "key2": 2.0, "key3": 3.0} {{ end }}

<!-- creates a type map[string]bool (N.B. Quotes) -->
{{ var "MapStringBool" }} {"key0": false, "key1": true} {{ end }}

<!-- creates a type map[string]string (N.B. Quotes) -->
{{ var "MapStringString" }} {"key1": "string 1", "key2": "string 2"} {{ end }}
```

### Attaching Variables to Templates

As an alternative to creating variables in the templates directly, variables can be directly assigned to any template **before `Parse()` is called** *(and they will be picked up by all bundles which use the file)*. This offers more freedom to define variable types.

```go
tm.AddParam("test.html", "MyMap", map[string][]int{"test": []int{1, 2, 3})
// OR
tm.AddParams("test.html", templateManager.Params{"MyInt": 123, "MyFloat": -42})
```

this can be done more globally too by attaching them to a layout:

```go
tm.AddParam("layout/public.html", "MyMap", map[string][]int{"test": []int{1, 2, 3})
// OR
tm.AddParams("layout/public.html", templateManager.Params{"MyInt": 123, "MyFloat": -42})
```

this way, all files that extend that layout will have these variables.

**N.B. This method has a lower precedence in the hierarchy than defining variables directly into the template files and the same variables defined there will override these**

## Creating Functions

Functions to manipulate variables may be created and passed to the templates. At present, functions are passed to ALL templates and cannot be passed to only a select few.

```go
tm.AddFunction("add", add)
// OR
tm.AddFunctions(map[string]any{"add": add})

func add[T int|float64](numbers ...T) T {
	var result T
	for _, number := range numbers {
		result += number
	}
	return result
}
```

It is not currently possible to return errors from functions, so any errors are output directly to STDOUT.

If a function receives invalid input, it will stop execution of the entire template at the point that the error occurs, so it's important to write flexible functions, or be careful as to where they are called.

It might be safer to rewrite the above test function (`add()`) using the `reflect` package (used throughout `text/templates`) to allow for type checking and sensible return values for unsupported types *(e.g. do nothing, or return 0)*.

## Built-in Functions

A selection of useful functions have been created to use in the templates to compliment those already built in to `text/template`. They are documented in their own [guide](FUNCTIONS.md), quick links:

[`add`](FUNCTIONS.md#add), [`capfirst`](FUNCTIONS.md#capfirst), [`collection`](FUNCTIONS.md#collection), [`concat`](FUNCTIONS.md#concat), [`contains`](FUNCTIONS.md#contains), [`cut`](FUNCTIONS.md#cut), [`date`](FUNCTIONS.md#date), [`datetime`](FUNCTIONS.md#datetime), [`default`](FUNCTIONS.md#default), [`divide`](FUNCTIONS.md#divide), [`divideceil`](FUNCTIONS.md#divideceil), [`dividefloor`](FUNCTIONS.md#dividefloor), [`divisibleby`](FUNCTIONS.md#divisibleby), [`dl`](FUNCTIONS.md#dl), [`endswith`](FUNCTIONS.md#endswith), [`equal`](FUNCTIONS.md#equal), [`first`](FUNCTIONS.md#first), [`firstof`](FUNCTIONS.md#firstof), [`formattime`](FUNCTIONS.md#formattime), [`gto`](FUNCTIONS.md#gto-greater-than), [`gte`](FUNCTIONS.md#gte-greater-than-equal), [`htmldecode`](FUNCTIONS.md#htmldecode), [`htmlencode`](FUNCTIONS.md#htmlencode), [`iterable`](FUNCTIONS.md#iterable), [`join`](FUNCTIONS.md#join), [`jsondecode`](FUNCTIONS.md#jsondecode), [`jsonencode`](FUNCTIONS.md#jsonencode), [`key`](FUNCTIONS.md#key), [`kind`](FUNCTIONS#kind), [`last`](FUNCTIONS.md#last), [`length`](FUNCTIONS.md#length), [`list`](FUNCTIONS.md#list), [`lto`](FUNCTIONS.md#lto-less-than), [`lte`](FUNCTIONS.md#lte-less-than-equal), [`localtime`](FUNCTIONS.md#localtime), [`lower`](FUNCTIONS.md#lower), [`ltrim`](FUNCTIONS.md#ltrim), [`md5`](FUNCTIONS.md#md5), [`mktime`](FUNCTIONS.md#mktime), [`multiply`](FUNCTIONS.md#multiply), [`nl2br`](FUNCTIONS.md#nl2br), [`notequal`](FUNCTIONS.md#notequal), [`now`](FUNCTIONS.md#now), [`ol`](FUNCTIONS.md#ol), [`ordinal`](FUNCTIONS.md#ordinal), [`paragraph`](FUNCTIONS.md#paragraph), [`pluralise`](FUNCTIONS.md#pluralise), [`prefix`](FUNCTIONS.md#prefix), [`query`](FUNCTIONS.md#query), [`random`](FUNCTIONS.md#random), [`regexp`](FUNCTIONS.md#regexp), [`regexpreplace`](FUNCTIONS.md#regexpreplace), [`render`](FUNCTIONS.md#render), [`replace`](FUNCTIONS.md#replace), [`round`](FUNCTIONS.md#round), [`rtrim`](FUNCTIONS.md#rtrim), [`sha1`](FUNCTIONS.md#sha1), [`sha256`](FUNCTIONS.md#sha256), [`sha512`](FUNCTIONS.md#sha512), [`split`](FUNCTIONS.md#split), [`startswith`](FUNCTIONS.md#startswith), [`striptags`](FUNCTIONS.md#striptags), [`substr`](FUNCTIONS.md#substr), [`subtract`](FUNCTIONS.md#subtract), [`suffix`](FUNCTIONS.md#suffix), [`time`](FUNCTIONS.md#time), [`timesince`](FUNCTIONS.md#timesince), [`timeuntil`](FUNCTIONS.md#timeuntil), [`title`](FUNCTIONS.md#title), [`trim`](FUNCTIONS.md#trim), [`truncate`](FUNCTIONS.md#truncate), [`truncatewords`](FUNCTIONS.md#truncatewords), [`type`](FUNCTIONS.md#type), [`ul`](FUNCTIONS.md#ul), [`upper`](FUNCTIONS.md#upper), [`urldecode`](FUNCTIONS.md#urldecode), [`urlencode`](FUNCTIONS.md#urlencode), [`uuid`](FUNCTIONS.md#uuid), [`wordcount`](FUNCTIONS.md#wordcount), [`wrap`](FUNCTIONS.md#wrap), [`year`](FUNCTIONS.md#year), [`yesno`](FUNCTIONS.md#yesno)

They are all added by default, but can be removed if necessary *(e.g. before adding any functions of your own)*:

```go
tm.RemoveAllFunctions()
tm.RemoveFunction("striptags")
tm.RemoveFunctions([]string{"yesno", "year"})
```

*(N.B. this does not remove the functions built in to `text/template` - [see guide](BASICS.md))*

### Overloading `text/template` Functions

Many of the built-in `text/template` functions throw errors which halt template execution as they are encountered, and are not optimised for their own pipelining system *(i.e. they receive their principle argument first, not last)*. For this reason those functions can be replaced by their equivalents from `templateManager`:

```go
tm.OverloadFunctions()
```

This will replace: `eq`, `ge`, `len`, `index`, `lt`, `le`, `ne`, `html` and `urlquery`.

This can also be undone:

```go
tm.RemoveOverloadFunctions()
tm.RemoveFunction("eq")
tm.RemoveFunctions([]string{"ge", "le"})
```

## Error Handling

TODO add error handling!! 

## Simple Example

To illustrate the `templateManager` usage, a trivial example with 5 files can be used:

`main.go`
```go
package main

import (
	"net/http"

	TM "github.com/paul-norman/go-template-manager"
)

var tm *TM.TemplateManager

func main() {
	tm = TM.Init("templates", ".html").
		ExcludeFolders([]string{"layouts", "partials"})
	err := tm.Parse()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/test", test)
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tm.Render("home.html", TM.Params{"Title": "Home"}, w)
}

func test(w http.ResponseWriter, r *http.Request) {
	tm.Render("test.html", TM.Params{"Title": "Test"}, w)
}
```

`templates/home.html`
```django
{{ extends "layouts/public.html" }}

{{ var "languageCode" }} en-GB {{ end }}

{{ define "title" }}{{ .Title }} Title{{ end }}

{{ define "content" }}
<section>
	<div>
		<h1>{{ .Title }}</h1>
	</div>
</section>
{{ end }}
```

`templates/test.html`
```django
{{ extends "layouts/public.html" }}

{{ define "description" }}{{ .Title }} Description{{ end }}

{{ define "content" }}
<section>
	<div>
		<h1>{{ .Title }}</h1>
	</div>
</section>
{{ end }}
```

`templates/layouts/public.html`
```django
{{ var "LanguageCode" }} en-US {{ end }}
<!DOCTYPE html lang="{{ .LanguageCode }}">
<html>
<head>
	<title>{{ block "title" "" }}default title{{ end }}</title>
	<meta name="description" value="{{ block "description" "" }}default description{{ end }}">
	{{ template "partials/meta.html" . }}
</head>
<body>
	{{- block "content" . }}default content{{ end -}}
</body>
</html>
```

`templates/partials/meta.html`
```django
{{ var "LanguageCode" }} en {{ end }}
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="language" content="{{ .LanguageCode }}">
```

Running the server and visiting `http://127.0.0.1:8080` in a browser would load the `home.html` template with a title displaying: "Home Title", a description displaying: "default description" *(fallback)*, and the `LanguageCode` displaying as "en-GB" throughout.

Visiting `http://127.0.0.1:8080/test` would load the `test.html` template with a title displaying: "default title" *(fallback)*, a description displaying: "Test description", and the `LanguageCode` set as "en-US" throughout.

The `LanguageCode` defined in `meta.html` is never needed, but would be used throughout if neither the layout (`public.html`) nor the entry (e.g. `home.html`) templates defined it. This allows flexible fallback variables to be set in the templates themselves.

## Embedded Example

It is possible to embed the templates within the package (so that the files can be accessed using the embedded filesystem) using the `InitEmbed` method. This accepts an extra parameter that is the embedded files (`embed.FS`). The template folder directory value is optional, but to make the template naming behave identically to the `Init` version it is required to match the template directory name. Without it, all templates must use longer paths in block *(e.g. `/templates/partials/layout.html` and not `partials/layout.html`)*

Just altering the main file from the simple example:

`main.go`
```go
package main

import (
	"embed"
	"net/http"

	TM "github.com/paul-norman/go-template-manager"
)

var tm *TM.TemplateManager

//go:embed templates/*
var embeddedTemplates embed.FS

func main() {
	tm = TM.InitEmbed(embeddedTemplates, "templates", ".html").
		ExcludeFolders([]string{"layouts", "partials"})
	err := tm.Parse()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/test", test)
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tm.Render("home.html", TM.Params{"Title": "Home"}, w)
}

func test(w http.ResponseWriter, r *http.Request) {
	tm.Render("test.html", TM.Params{"Title": "Test"}, w)
}
```

## Integrations

Currently `templateManager` has an integrations for:

- [Fiber](https://gofiber.io/) in its [own repository](https://github.com/paul-norman/go-template-manager-fiber) with a basic example.
- [Echo](https://echo.labstack.com/) in its [own repository](https://github.com/paul-norman/go-template-manager-echo) with a basic example.
- [Gin](https://gin-gonic.com/) in its [own repository](https://github.com/paul-norman/go-template-manager-gin) with a basic example.

[Martini](https://github.com/go-martini/martini) doesn't require an integration as it can be implemented as in the trivial example [above](#simple-example).