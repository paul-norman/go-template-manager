package templateManager

import (
	"fmt"
	"io"
	TT "text/template"
	HT "html/template"
)

func NewTemplate(typ string, name string) *Template {
	if typ == "html" {
		return NewHtmlTemplate(name)
	}
	return NewTextTemplate(name)
}
func NewTextTemplate(name string) *Template {
	template := &Template{ nil, nil }
	template.InitText(name)
	return template
}
func NewHtmlTemplate(name string) *Template {
	template := &Template{ nil, nil }
	template.InitHtml(name)
	return template
}

/*
Allows the use of either HTML or Text template packages
*/
type Template struct {
	html *HT.Template
	text *TT.Template
}
func (t *Template) InitText(name string) *TT.Template {
	t.text = TT.New(name)
	t.html = nil
	return t.text
}
func (t *Template) InitHtml(name string) *HT.Template {
	t.html = HT.New(name)
	t.text = nil
	return t.html
}
func (t *Template) NewSubTemplate(name string) *Template {
	if t.Type() == "text" {
		template := t.text.New(name)
		return &Template{ nil, template }
	} else if t.Type() == "html" {
		template := t.html.New(name)
		return &Template{ template, nil }
	}
	return t
}
func (t *Template) Type() string {
	if t.text == nil && t.html == nil {
		return "invalid"
	} else if t.text == nil {
		return "html"
	}
	return "text"
}
func (t *Template) Delims(left string, right string) *Template {
	if t.Type() == "text" {
		t.text.Delims(left, right)
	} else if t.Type() == "html" {
		t.html.Delims(left, right)
	}
	return t
}
func (t *Template) Funcs(funcs map[string]any) *Template {
	if t.Type() == "text" {
		t.text.Funcs(funcs)
	} else if t.Type() == "html" {
		t.html.Funcs(funcs)
	}
	return t
}
func (t *Template) Execute(wr io.Writer, data any) error {
	if t.Type() == "text" {
		return t.text.Execute(wr, data)
	} else if t.Type() == "html" {
		return t.html.Execute(wr, data)
	}
	return fmt.Errorf("templateManager Template not loaded correctly")
}
func (t *Template) ExecuteTemplate(wr io.Writer, name string, data any) error {
	if t.Type() == "text" {
		return t.text.ExecuteTemplate(wr, name, data)
	} else if t.Type() == "html" {
		return t.html.ExecuteTemplate(wr, name, data)
	}
	return fmt.Errorf("templateManager Template not loaded correctly")
}
func (t *Template) Parse(text string) (*Template, error) {
	var err error
	if t.Type() == "text" {
		_, err = t.text.Parse(text)
	} else if t.Type() == "html" {
		_, err = t.html.Parse(text)
	} else {
		err = fmt.Errorf("templateManager Template not loaded correctly")
	}
	return t, err
}
func (t *Template) Lookup(name string) *Template {
	if t.Type() == "text" {
		template := t.text.Lookup(name)
		if template == nil {
			return nil
		}
		return &Template{ nil, template }
	} else if t.Type() == "html" {
		template := t.html.Lookup(name)
		if template == nil {
			return nil
		}
		return &Template{ template, nil }
	}
	return nil
}