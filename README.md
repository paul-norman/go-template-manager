# TemplateManager

Package `templateManager` simplifies the use of Go's standard libraries: `text/template` / `html/template` for use with HTML templates.

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

`templateManager` generally just simplifies the usage of the `text/templates` or `html/templates` package. A basic usage guide to these [Go Templates](https://pkg.go.dev/text/template) is provided [here](BASICS.md).

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

`templateManager` requires initialisation. At a minimum it requires you to tell it where the templates are kept, and what file types the template files are:

```go
var tm = templateManager.Init("templates", ".html", ".htm")

// OR for embedded templates

//go:embed templates/*
var embeddedTemplates embed.FS
var tm = templateManager.InitEmbed(embeddedTemplates, ".html", ".htm")
```

During parsing, `templateManager` only needs to know about the location of the "entry" template files as it will scan each for its dependencies. It is therefore best that it does not know about layout files, partials or components and these should be kept separate and excluded from parsing:

```go
tm.ExcludeDirectories([]string{"layouts", "partials", "components"})
```

Once these variables have been configured, the template parsing can be triggered to build the cache.

```go
tm.Parse()
```

To run a template, use the `Render()` method:

```go
err := tm.Render("test.html", TM.Params{"Title": "Test"}, ioWriter)
```

If an error is encountered, no output will be written to `ioWriter` allowing you to display a custom error page of your choosing.

## Customisation Options

All customisation options are chainable for neat declaration.

### Engine Choice

Templates may be rendered with either the `text/template` package *(default)* or the `html/template` package:

```go
tm.TemplateEngine("html")
```

### Delimiters

The default delimiters set in the `text/template` package are `{{` and `}}` which wrap all commands. These can be customised if desired before templates are parsed:

```go
tm.Delimiters("{%", "%}")
```

### Debugging

During development it's often useful to see what is happening. Enabling debug mode outputs console entries showing what happens upon parsing *(and forces warnings and errors to be shown in the console)*:

```go
tm.Debug(true)
```

### Reload Each Time

During development it's often useful to be able to change template files and not have to restart the server to test them. Enabling reload mode forces a bundle rebuild for the selected template upon each `Render()` method call:

```go
tm.Reload(true)
```

*(N.B. this does not work with an embedded file system as the changes are not picked up until the next build)*

### Excluding Directories

It is most efficient if the parser only runs over "entry" templates *(i.e those which will be called directly)*. For this reason it's best to exclude all directories *(within the designated templates folder)* which do not contain entry templates from this process.

By default, three directories *("layouts", "partials" and "components")* are already excluded.

```go
// Remove directories from parsing
tm.ExcludeDirectory("layouts")
// OR
tm.ExcludeDirectories([]string{"layouts", "partials"})

// Re-adds a previously excluded directory so that it will be parsed
tm.RemoveExcludedDirectory("layouts")
```

### Functions

See the [Built-in Functions](#built-in-functions) section.

## Global Options

There are also several customisation options that apply globally to `templateManager` functions / use. These should be set **prior** to initialisation of the main store.

```go
import (
	"time"
	TM "github.com/paul-norman/go-template-manager"
)

// Control whether errors will cause rendering of the template to abort (default: true)
TM.SetHaltOnErrors(false)

// Control whether warnings will cause rendering of the template to abort (default: false)
TM.SetHaltOnWarnings(true)

// Control whether errors are written to the log (default: true)
TM.SetConsoleErrors(false)

// Control whether warnings are written to the log (default: true)
TM.SetConsoleWarnings(false)

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

These are parsed by regular expressions and turned into a "best-guess" version of what they represent *(leading / trailing spaces are always trimmed)*.

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
<!-- creates a type map[int]int -->
{{ var "MapIntInt" }} {1: 1, 10: -10, 100: 100} {{ end }}

<!-- creates a type map[int]float64  -->
{{ var "MapIntFloat" }} {1: 1.0, 10: -10.0, 100: 100.0} {{ end }}

<!-- creates a type map[int]bool -->
{{ var "MapIntBool" }} {0: false, 1: true} {{ end }}

<!-- creates a type map[int]string -->
{{ var "MapIntString" }} {1: "string 1", 2: "string 2"} {{ end }}

<!-- creates a type map[float64]int -->
{{ var "MapFloatInt" }} {1.1: 1, 10.10: -10, 100.100: 100} {{ end }}

<!-- creates a type map[float64]float64 -->
{{ var "MapFloatFloat" }} {1.0: 1.0, 10.0: -10.0, 100.0: 100.0} {{ end }}

<!-- creates a type map[float64]bool -->
{{ var "MapFloatBool" }} {0.0: false, 1.0: true} {{ end }}

<!-- creates a type map[float64]string -->
{{ var "MapFloatString" }} {1.0: "string 1.0", 2.0: "string 2.0"} {{ end }}

<!-- creates a type map[bool]int -->
{{ var "MapBoolInt" }} {true: 1, false: 0} {{ end }}

<!-- creates a type map[bool]float64 -->
{{ var "MapBoolFloat" }} {true: 1.0, false: 0.0} {{ end }}

<!-- creates a type map[bool]bool -->
{{ var "MapBoolBool" }} {false: true, true: false} {{ end }}

<!-- creates a type map[bool]string -->
{{ var "MapBoolString" }} {false: "no", true: "yes"} {{ end }}

<!-- creates a type map[string]int -->
{{ var "MapStringInt" }} {"key1": 1, "key2": 2, "key3": 3} {{ end }}

<!-- creates a type map[string]float64 -->
{{ var "MapStringFloat" }} {"key1": 1.0, "key2": 2.0, "key3": 3.0} {{ end }}

<!-- creates a type map[string]bool -->
{{ var "MapStringBool" }} {"key0": false, "key1": true} {{ end }}

<!-- creates a type map[string]string -->
{{ var "MapStringString" }} {"key1": "string 1", "key2": "string 2"} {{ end }}
```

### Attaching Variables to Templates

As an alternative to creating variables in the templates directly *(or at the `Render()` stage)*, variables can be directly assigned to any template **before `Parse()` is called** *(and they will be picked up by all bundles which use the file)*. This offers more freedom to define variable types.

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

**N.B. This method has a lower precedence in the hierarchy than defining variables directly into the template files *(and at `Render()` time)* and the same variables defined there will override these.**

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

In the `text/template` system, if a function receives invalid input, it will stop execution of the entire template at the point that the error occurs, so it's important to write flexible functions, or be careful as to how / where they are called.

It might be safer to rewrite the above test function (`add()`) using the `reflect` package *(which is used throughout `text/templates` and `templateManager`)* to allow for type checking and sensible return values for unsupported types *(e.g. do nothing, or return 0)*.

## Built-in Functions

A selection of useful functions have been created to use in the templates to compliment those already built in to `text/template`. These are all optimised for "pipeline" use *(i.e. receive their principle argument last)*. They are documented in their own [guide](FUNCTIONS.md), quick links:

[`add`](FUNCTIONS.md#add), [`bool`](FUNCTIONS.md#bool), [`capfirst`](FUNCTIONS.md#capfirst), [`collection`](FUNCTIONS.md#collection), [`concat`](FUNCTIONS.md#concat), [`contains`](FUNCTIONS.md#contains), [`cut`](FUNCTIONS.md#cut), [`date`](FUNCTIONS.md#date), [`datetime`](FUNCTIONS.md#datetime), [`default`](FUNCTIONS.md#default), [`divide`](FUNCTIONS.md#divide), [`divideceil`](FUNCTIONS.md#divideceil), [`dividefloor`](FUNCTIONS.md#dividefloor), [`divisibleby`](FUNCTIONS.md#divisibleby), [`dl`](FUNCTIONS.md#dl), [`endswith`](FUNCTIONS.md#endswith), [`equal`](FUNCTIONS.md#equal), [`first`](FUNCTIONS.md#first), [`firstof`](FUNCTIONS.md#firstof), [`float`](FUNCTIONS.md#float), [`formattime`](FUNCTIONS.md#formattime), [`gto`](FUNCTIONS.md#gto-greater-than), [`gte`](FUNCTIONS.md#gte-greater-than-equal), [`htmldecode`](FUNCTIONS.md#htmldecode), [`htmlencode`](FUNCTIONS.md#htmlencode), [`int`](FUNCTIONS.md#int), [`iterable`](FUNCTIONS.md#iterable), [`join`](FUNCTIONS.md#join), [`jsondecode`](FUNCTIONS.md#jsondecode), [`jsonencode`](FUNCTIONS.md#jsonencode), [`key`](FUNCTIONS.md#key), [`keys`](FUNCTIONS.md#keys), [`kind`](FUNCTIONS#kind), [`last`](FUNCTIONS.md#last), [`length`](FUNCTIONS.md#length), [`list`](FUNCTIONS.md#list), [`lto`](FUNCTIONS.md#lto-less-than), [`lte`](FUNCTIONS.md#lte-less-than-equal), [`localtime`](FUNCTIONS.md#localtime), [`lower`](FUNCTIONS.md#lower), [`lpad`](FUNCTIONS.md#lpad), [`ltrim`](FUNCTIONS.md#ltrim), [`md5`](FUNCTIONS.md#md5), [`mktime`](FUNCTIONS.md#mktime), [`multiply`](FUNCTIONS.md#multiply), [`nl2br`](FUNCTIONS.md#nl2br), [`notequal`](FUNCTIONS.md#notequal), [`now`](FUNCTIONS.md#now), [`ol`](FUNCTIONS.md#ol), [`ordinal`](FUNCTIONS.md#ordinal), [`paragraph`](FUNCTIONS.md#paragraph), [`pluralise`](FUNCTIONS.md#pluralise), [`prefix`](FUNCTIONS.md#prefix), [`query`](FUNCTIONS.md#query), [`random`](FUNCTIONS.md#random), [`regexp`](FUNCTIONS.md#regexp), [`regexpreplace`](FUNCTIONS.md#regexpreplace), [`render`](FUNCTIONS.md#render), [`replace`](FUNCTIONS.md#replace), [`round`](FUNCTIONS.md#round), [`rpad`](FUNCTIONS.md#rpad), [`rtrim`](FUNCTIONS.md#rtrim), [`sha1`](FUNCTIONS.md#sha1), [`sha256`](FUNCTIONS.md#sha256), [`sha512`](FUNCTIONS.md#sha512), [`split`](FUNCTIONS.md#split), [`startswith`](FUNCTIONS.md#startswith), [`string`](FUNCTIONS.md#string), [`striptags`](FUNCTIONS.md#striptags), [`substr`](FUNCTIONS.md#substr), [`subtract`](FUNCTIONS.md#subtract), [`suffix`](FUNCTIONS.md#suffix), [`time`](FUNCTIONS.md#time), [`timesince`](FUNCTIONS.md#timesince), [`timeuntil`](FUNCTIONS.md#timeuntil), [`title`](FUNCTIONS.md#title), [`trim`](FUNCTIONS.md#trim), [`truncate`](FUNCTIONS.md#truncate), [`truncatewords`](FUNCTIONS.md#truncatewords), [`type`](FUNCTIONS.md#type), [`ul`](FUNCTIONS.md#ul), [`upper`](FUNCTIONS.md#upper), [`urldecode`](FUNCTIONS.md#urldecode), [`urlencode`](FUNCTIONS.md#urlencode), [`uuid`](FUNCTIONS.md#uuid), [`values`](FUNCTIONS.md#values), [`wordcount`](FUNCTIONS.md#wordcount), [`wrap`](FUNCTIONS.md#wrap), [`year`](FUNCTIONS.md#year), [`yesno`](FUNCTIONS.md#yesno)

They are all added by default, but can be removed or renamed if necessary *(e.g. before adding any functions of your own)*:

```go
// Remove functions
tm.RemoveAllFunctions()
tm.RemoveFunction("striptags")
tm.RemoveFunctions([]string{"yesno", "year"})

// Rename functions
tm.RenameFunction("concat", "cat")
tm.RenameFunctions(map[string]string{ "equal": "equals", "lpad": "pad_left" })
```

*(N.B. this does not remove / rename the functions built in to `text/template` - [see guide](BASICS.md))*

### Overloading `text/template` Functions

Many of the built-in `text/template` functions throw errors which halt template execution as they are encountered, and are not optimised for their own pipelining system *(i.e. they receive their principle argument first, not last)*. For this reason those functions can be replaced by their equivalents from `templateManager`:

```go
tm.OverloadFunctions()
```

This will replace: `eq`, `ge`, `len`, `index`, `lt`, `le`, `ne`, `html` and `urlquery`.

This can also be undone for all or some of the functions overloaded:

```go
tm.RemoveOverloadFunctions()
tm.RemoveFunction("eq")
tm.RemoveFunctions([]string{"ge", "le"})
```

## Error Handling

The `text/template` package allows functions to return errors which will halt execution immediately, mid way through a document. This is often undesirable.

For this reason, `templateManager` only outputs any data if no errors are encountered. In production environments, this is useful so as to be able to display custom 500 pages. However, in development it is often desirable to have errors in the console, but allow execution to complete anyway *(so as better to see what is happening)*.

To achieve this, all `templateManager` functions generate errors in a manner where they can be controlled. There are two types of `error`, a warning and a full error. The latter is designed to alert the developer to them doing something likely undesirable *(e.g. trying to divide by a `nil` variable or zero)*, while a warning is designed to draw attention to something that only *might* be incorrect *(e.g. testing if an unset variable is equal to 1)*.

These special errors may be instructed to display in the console and also to halt the program if they are encountered. These are global options:

```go
import (
	TM "github.com/paul-norman/go-template-manager"
)

// Control whether errors will cause rendering of the template to abort (default: true)
TM.SetHaltOnErrors(false)

// Control whether warnings will cause rendering of the template to abort (default: false)
TM.SetHaltOnWarnings(true)

// Control whether errors are written to the log (default: true)
TM.SetConsoleErrors(false)

// Control whether warnings are written to the log (default: true)
TM.SetConsoleWarnings(false)
```

Sadly, there is no means of getting detailed information *(i.e. template name, line number etc)* for errors, so this is only provided if an error bubbles back up to the `Render` function and is caught there. *(N.B. if haltOnError is disabled, these bubbled errors are missing, but that shouldn't be a problem during development as they will be generated as pages are used).*

## Simple Example

To illustrate `templateManager` usage, a trivial example with 5 files can be used:

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
		ExcludeDirectories([]string{"layouts", "partials"})
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

{{ var "LanguageCode" }} en-GB {{ end }}

{{- define "title" -}}{{ .Title }} Title {{- end -}}

{{- define "content" -}}
<section>
	<div>
		<h1>{{ .Title }}</h1>
	</div>
</section>
{{- end -}}
```

`templates/test.html`
```django
{{ extends "layouts/public.html" }}

{{- define "description" -}}{{ .Title }} Description{{- end -}}

{{- define "content" -}}
<section>
	<div>
		<h1>{{ .Title }}</h1>
	</div>
</section>
{{- end -}}
```

`templates/layouts/public.html`
```django
{{ var "LanguageCode" }} en-US {{ end }}
<!DOCTYPE html>
<html lang="{{ .LanguageCode }}">
<head>
	<title>{{ block "title" . -}} default title {{- end }}</title>
	<meta name="description" value="{{ block "description" "" -}} default description {{- end }}">
	{{ template "partials/meta.html" . }}
</head>
<body>
	{{- block "content" . -}} default content {{- end -}}
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

It is possible to embed the templates within the package *(so that the files can be accessed using the embedded filesystem)* using the `InitEmbed` method. This accepts an extra parameter that is the embedded files (`embed.FS`).

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
		ExcludeDirectories([]string{"layouts", "partials"})
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