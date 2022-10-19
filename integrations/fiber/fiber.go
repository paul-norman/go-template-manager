package fiber

import (
	"fmt"
	"io"

	TM "github.com/paul-norman/go-template-manager"
)

type Engine struct {
	tm TM.TemplateManager
}

func New(directory string, extension string) *Engine {
	return &Engine{*TM.Init(directory, extension)}
}

/*
func NewFileSystem(fs http.FileSystem, extension string) *TM.TemplateManager {
	return &Engine{
		TM.Init(directory, extension)
	}
}
*/

// Not required
func (e *Engine) Layout(key string) *Engine {
	return e
}

func (e *Engine) Delims(left string, right string) *Engine {
	e.tm.Delimiters(left, right)

	return e
}

func (e *Engine) AddFunc(name string, function interface{}) *Engine {
	e.tm.AddFunction(name, function)

	return e
}

func (e *Engine) AddFuncMap(functions map[string]interface{}) *Engine {
	e.tm.AddFunctions(functions)

	return e
}

func (e *Engine) Reload(enabled bool) *Engine {
	e.tm.Reload(enabled)

	return e
}

func (e *Engine) Debug(enabled bool) *Engine {
	e.tm.Debug(enabled)

	return e
}

func (e *Engine) Parse() error {
	return fmt.Errorf("Parse() is deprecated, please use Load() instead")
}

func (e *Engine) Load() error {
	return e.tm.Parse()
}

func (e *Engine) Render(out io.Writer, template string, params TM.Params, layout ...string) error {
	return e.tm.Render(template, params, out)
}