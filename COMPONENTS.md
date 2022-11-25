# Components in `templateManager`

It is often convenient to define a template for a certain, re-usable piece of code. However, passing only the required variables to `text/template` blocks can get a little messy and passing in sub-templates is awkward. It may be a cleaner experience if these are presented in an HTML format so that it is clear what is being passed to them and trivial to pass in other templates as arguments.

## Contents

- [Configuration](#configuration)
- [Explanatory Example](#explanatory-example)
- [Passing Data to Components](#passing-data-to-components)
- [Nested Components](#nested-components)

## Configuration

`templateManager` components require initialisation. At a minimum the component location where components are kept must be set **(within the templates directory)** *(default: "components")*. Individual component files must match your declared templates file type.

```go
tm.AddComponentDirectory("components")
// OR
tm.AddComponentDirectories([]string{"components", "partials"})

// Remove previously declared component directories
tm.RemoveComponentDirectory("components")
// OR
tm.RemoveComponentDirectories([]string{"components", "partials"})
```

The filenames of the component files will lend their names to their matching HTML components, so `Youtube.html` would allow use of a `<Youtube>` component.

## Explanatory Example

As a simple example, a Youtube component could replace a template:

```go
{{ template "components/Youtube.html" .Src }}
```

could become:

```html
<Youtube Src="{{ .Src }}">
```

This isn't much of an improvement yet, but if several options were required to be passed to the template the benefits can be more clearly seen:

```go
{{ template "components/Youtube.html" .PreOrganisedVariables }}
// OR
{{ template "components/Youtube.html" collection "Id" .Id "Language" .Site.Language "Subtitles" 1 }}
```

might become:

```html
<Youtube Id="{{ .Id }}" Language="{{ .Site.Language }}" Subtitles=1>
```

Simplifying the markup whilst making it clear what is being sent to the component.

## Passing Data to Components

Components are completely normal `text/template` files. They **cannot** be extended *(like normal `templateManager` files can)*, and do not support internal variables either. When components are used, a new template is created from the file so that any number of unique instances of each component may be used in a single file.

Within the component, the attribute names are mapped to their corresponding variables. So in the example [above](#explanatory-example), the variables `Id`, `Language` and `Subtitles` would be known to the component as `.Id`, `.Language` and `.Subtitles` respectively.

Attributes may be numeric or string values. If quotes are used, the value will be interpreted as a `string`, and if they are omitted it will either be a `float64` or `int` depending upon whether a decimal point is included.

Each component will also be assigned a unique identifier (uuid), which is available as `.ComponentUuid` and can be used for many purposes.

So the example `Youtube` component above might look like this:

```html
<iframe src="https://www.youtube-nocookie.com/embed/{{ .Id }}&hl={{ .Language }}&cc_lang_pref={{ .Language }}&cc_load_policy={{ .Subtitles }}" loading="lazy" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
```

Components may call other components if desirable. So if you had a `Vimeo` component and a `Youtube` component, both could call a `VideoIframe` component internally. For example the `Youtube` component could be refactored to something like:

```html
{{- $src := "https://www.youtube-nocookie.com/embed/" | suffix .Id -}}
{{- with .Language -}} {{ $src = printf "%v&hl=%v&cc_lang_pref=%v" $src . . }} {{- end -}}
{{- with .Subtitles -}} {{ $src = printf "%v&cc_load_policy=%v" $src . }} {{- end -}}
<VideoIframe src="{{ $src }}">
```

which could support missing attributes.

### Capturing Wrapped Content

HTML tags are designed to wrap content, and so are components. It may be that you would prefer to create a Youtube component that wraps the Youtube source rather than using an attribute, for example: 

```html
<Youtube>{{ .Src }}</Youtube>
```

In this case the wrapped content is defined as its own template and the name of this template is passed to the component file as `.ComponentContent`. Sadly the `text/template` package will not render a template from a variable, but there is a replacement function, `render`, to do this for us:

```go
{{ render .ComponentContent . }}
```

this will output the wrapped content with all variables known to the template available to it *(of course you may share only the variables that are needed if you wish)*.

This may also be captured as a variable if preferred:

```go
{{ $content := render .ComponentContent . }}
```

## Nested Components

Nesting components is allowed and can make sense in many cases, for example a `Slideshow` component might have many `Slides`:

```html
<Slideshow Id="test">
	<Slide Active=1><img src="slide1.png"></Slide>
	<Slide><img src="slide2.png"></Slide>
</Slideshow>
```

In this case all components will be individually rendered with **no information shared between parent and children**. The `render` function will still need to be called within the `Slideshow` component using the `.ComponentContent` variable to display the rendered `Slide` items. 

### Collecting Nesting Components

It is not always desirable to render wrapped information in its entirety or in the order in which it was declared. It may be that the information collected needs to be used in different parts of the master template. An example might be a tabbed interface where a set of tabs and a set of corresponding content windows must be collected, but the tabs need to be rendered at the top of the markup with the windows lower down. Conceptually:

```html
<Tabset>
	<Tab>Tab 1</Tab>
	<Tab>Tab 2</Tab>

	<TabContent>Content 1</TabContent>
	<TabContent Active=1>Content 2</TabContent>
</Tabset>
```

It is likely that the `Tab` items will need to be wrapped as will the `TabContent` items and they could not simply be rendered as is *(or at least not without heavy JavaScript support)*. To solve this, the nested components may be prefixed with an `x-` string to signify that we wish to collect them into a variable for later output:

```html
<Tabset>
	<x-Tab>Tab 1</x-Tab>
	<x-Tab>Tab 2</x-Tab>

	<x-TabContent>Content 1</x-TabContent>
	<x-TabContent Active=1>Content 2</x-TabContent>
</Tabset>
```

Now there will be two extra variables available to the parent template `.Tab` and `.TabContent` *(in this example)* which will contain a string slice of the rendered templates in the order that they were defined.

This could now be re-written as:

```html
<Tabset>
	<x-Tab>Tab 1</x-Tab>
	<x-TabContent>Content 1</x-TabContent>

	<x-Tab>Tab 2</x-Tab>
	<x-TabContent Active=1>Content 2</x-TabContent>
</Tabset>
```

and function identically.

In addition to all standard data passed to these collected, nested templates, they are also assigned two additional variables: `.ParentUuid` and `.ParentPosition`.

`.ParentUuid` is the uuid of the parent component, and `.ParentPosition` is the index that the specific component occupies in the string slice passed to the parent component *(i.e. the position of the item within the `.Tab` slice)*.

These variables should allow tricks such as the CSS checkbox hack to be implemented without assigning names / ids to all nested items, keeping the code as clean as possible.