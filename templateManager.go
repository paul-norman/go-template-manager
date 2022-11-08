/*
Package templateManager simplifies the use of Go's standard library `text/template` for use with HTML templates. 
It DOES NOT use the `html/template` package to avoid over-zealous HTML escaping (good for security, terrible for usability).

It automates the process of choosing which files to group together for parsing, and builds a cache of each entry 
template file complete with all of its dependencies.

It does this by allowing files to "extend" their layout templates without manually specifying all bundle files (or 
including everything via the built-in `template` function. Instead it adds a new tag:

 {{ extends "layouts/main.html" }}

which automatically defines the specified file as part of the bundle. It can then follow all instances of the template 
tag in these two files until all files are in the bundle. This ensures that only the correct blocks exist.

It also defines a second new tag to allow VERY basic variables to be defined within the templates themself and allows them to be overridden via 
a simple hierarchy:

 {{ var "int1" }} 123 {{ end }}

templateManager also comes with a small set of extra convenience functions which may be used or ignored.
*/
package templateManager

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"text/template"

	"golang.org/x/exp/slices"

	"github.com/paul-norman/go-template-manager/fsWalk"
)

// Holds all templates and variables along with all required settings
type TemplateManager struct {
	templates 			map[string]*template.Template
	params				map[string]map[string]any
	descendants			map[string][]string
	delimiterLeft		string
	delimiterRight		string
	fileSystem			http.FileSystem
	directory			string
	extension			string
	excludedDirectories	[]string
	functions			map[string]any
	mutex				sync.RWMutex
	debug				bool
	reload				bool
	parsed				bool
}

// Convenience type allowing any variables types to be passed in
type Params map[string]any

// Creates a new `TemplateManager` struct instance
func Init(directory string, extension string) *TemplateManager {
	templateManager := &TemplateManager{
		templates:				make(map[string]*template.Template),
		params:					make(map[string]map[string]any),
		descendants:			make(map[string][]string),
		delimiterLeft:			"{{",
		delimiterRight:			"}}",
		directory:				directory,
		extension:				extension,
		excludedDirectories:	[]string{"layouts", "partials"},
		functions:				make(map[string]any),
		debug:					false,
		reload:					false,
		parsed:					false,
	}

	templateManager.addDefaultFunctions()

	return templateManager
}

func InitEmbed(fileSystem embed.FS, directory string, extension string) *TemplateManager {
	templateManager := &TemplateManager{
		templates:				make(map[string]*template.Template),
		params:					make(map[string]map[string]any),
		descendants:			make(map[string][]string),
		delimiterLeft:			"{{",
		delimiterRight:			"}}",
		directory:				"/" + strings.TrimLeft(directory, " /"),
		fileSystem:				http.FS(fileSystem),
		extension:				extension,
		excludedDirectories:	[]string{"layouts", "partials"},
		functions:				make(map[string]any),
		debug:					false,
		reload:					false,
		parsed:					false,
	}

	templateManager.addDefaultFunctions()

	return templateManager
}

// Attempts to automatically configure a `TemplateManager` instance
func AutomaticInit() (*TemplateManager, error) {
	scan, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	templateDirectory := ""
	templateExtension := ""
	allowedExtensions := []string{".go.html", ".go.tmpl", ".html.tmpl", ".html.tpl", ".html", ".tmpl", ".tpl", ".gohtml", ".gotmpl", ".gotpl", ".thtml", ".htm"}

	err = filepath.WalkDir(scan, func(path string, info fs.DirEntry, err error) error {
		if err != nil || info == nil {
			return err
		}

		if path == scan {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		name, err := cleanPath(path, scan)
		if err != nil {
			return err
		}

		if strings.HasPrefix(name, "templates/") || strings.HasPrefix(name, "views/") {
			extension, err := checkAllowedExtension(name, allowedExtensions)
			if err == nil {
				templateDirectory		= strings.Split(name, "/")[0]
				templateExtension	= extension
				return io.EOF
			}
		} else if strings.Contains(name, "/") {
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil && err != io.EOF {
		return nil, err
	}

	if len(templateDirectory) < 1 {
		return nil, fmt.Errorf("templates were not manually initialised and automatic detection failed")
	}

	return Init(templateDirectory, templateExtension), nil
}

// Adds a custom function for use in all templates within the instance of `TemplateManager`
func (tm *TemplateManager) AddFunction(name string, function any) *TemplateManager {
	tm.mutex.Lock()
	tm.functions[name] = function
	tm.mutex.Unlock()

	return tm
}

// Adds multiple custom functions for use in all templates within the instance of `TemplateManager`
// Function names are the map keys.
func (tm *TemplateManager) AddFunctions(functions map[string]any) *TemplateManager {
	tm.mutex.Lock()
	for name, function := range functions {
		tm.functions[name] = function
	}
	tm.mutex.Unlock()

	return tm
}

// Adds a single variable (`name`) with value `value` that will always be available in the `templateName` template
func (tm *TemplateManager) AddParam(templateName string, name string, value any) *TemplateManager {
	if _, ok := tm.params[templateName]; !ok {
		tm.params[templateName] = make(Params)
	}
	
	tm.params[templateName][name] = value

	return tm
}

// Adds several variables (`params`) that will always be available in the `templateName` template
func (tm *TemplateManager) AddParams(templateName string, params Params) *TemplateManager {
	for name, param := range params {
		tm.AddParam(templateName, name, param)
	}

	return tm
}

// Sets the delimiters used by `text/template` (Default: "{{" and "}}")
func (tm *TemplateManager) Delimiters(left string, right string) *TemplateManager {
	tm.delimiterLeft	= left
	tm.delimiterRight	= right

	return tm
}

// Enable debugging of the template build process
func (tm *TemplateManager) Debug(debug bool) *TemplateManager {
	tm.debug = debug

	logErrors	= true
	logWarnings	= true

	return tm
}

// Excludes multiple directories from the build scanning process (which only wants entry files).
// This does not prevent files in these directories from being included via `template`.
// Typically, directories containing base layouts and partials should be excluded.
func (tm *TemplateManager) ExcludeDirectories(directories []string) *TemplateManager {
	for _, directory := range directories {
		tm.ExcludeDirectory(directory)
	}

	return tm
}

// Exclude a directory from the build scanning process (which only wants entry files).
// This does not prevent files in this directory from being included via `template`.
// Typically, directories containing base layouts and partials should be excluded.
func (tm *TemplateManager) ExcludeDirectory(directory string) *TemplateManager {
	if !slices.Contains(tm.excludedDirectories, directory) {
		tm.excludedDirectories = append(tm.excludedDirectories, directory)
	}

	return tm
}

// Replaces standard `text/template` functions with the `TemplateManager` alternatives
func (tm *TemplateManager) OverloadFunctions() *TemplateManager {
	tm.AddFunctions(getOverloadFunctions())

	return tm
}

// Triggers scanning of files and bundling of all templates
func (tm *TemplateManager) Parse() error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if tm.parsed {
		return nil
	}

	if tm.debug {
		logWarning("Parsing all templates...")
	}

	var err error

	walk := func(path string, info fs.DirEntry, err error) error {
		if err != nil || info == nil {
			return err
		}

		if len(tm.extension) >= len(path) || path[len(path) - len(tm.extension):] != tm.extension {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		name, _ := cleanPath(path, tm.directory)

		for _, excludedDirectory := range tm.excludedDirectories {
			if strings.HasPrefix(name, excludedDirectory + "/") {
				return nil
			}
		}

		tm.templates[name] = tm.configureNewTemplate(template.New(""))
		err = tm.parseFileDependencies(path, name, tm.directory, tm.templates[name])
		if err != nil {
			return err
		}

		return nil
	}

	if tm.fileSystem != nil {
		err = fsWalk.WalkDir(tm.fileSystem, "/", walk)
	} else {
		err = filepath.WalkDir(tm.directory, walk)
	}

	if err == nil {
		tm.parsed = true

		if tm.debug {
			logSuccess("All templates parsed and ready to use")
		}
	}

	return err
}

// Enable re-rebuilding of the template bundle upon every page load (for development)
func (tm *TemplateManager) Reload(reload bool) *TemplateManager {
	tm.reload = reload

	return tm
}

// Removes all functions currently assigned to the instance of `TemplateManager`.
// Useful if you do not want the default functions included
func (tm *TemplateManager) RemoveAllFunctions() *TemplateManager {
	tm.mutex.Lock()
	tm.functions = make(map[string]any)
	tm.mutex.Unlock()

	return tm
}

// Removes a directory that was previously excluded to allow it to feature in the build scanning process (which only wants entry files).
func (tm *TemplateManager) RemoveExcludedDirectory(directory string) *TemplateManager {
	if slices.Contains(tm.excludedDirectories, directory) {
		index := slices.Index(tm.excludedDirectories, directory)
		tm.excludedDirectories = append(tm.excludedDirectories[:index], tm.excludedDirectories[index + 1:]...)
	}

	return tm
}

// Removes the function named (does not affect built-in `template/text` functions)
func (tm *TemplateManager) RemoveFunction(name string) *TemplateManager {
	tm.mutex.Lock()
	delete(tm.functions, name)
	tm.mutex.Unlock()

	return tm
}

// Removes the functions named (does not affect built-in `template/text` functions)
func (tm *TemplateManager) RemoveFunctions(names []string) *TemplateManager {
	tm.mutex.Lock()
	for _, name := range names {
		delete(tm.functions, name)
	}
	tm.mutex.Unlock()

	return tm
}

// Replaces standard `text/template` functions with the `TemplateManager` alternatives
func (tm *TemplateManager) RemoveOverloadFunctions() *TemplateManager {
	names := getOverloadFunctions()
	tm.mutex.Lock()
	for name, _ := range names {
		delete(tm.functions, name)
	}
	tm.mutex.Unlock()

	return tm
}

// Executes a single template (`name`)
func (tm *TemplateManager) Render(name string, params Params, writer io.Writer) error {
	if ! tm.parsed {
		err := tm.Parse()
		if err != nil {
			logError(err.Error())
			return err
		}
	}

	if tm.reload {
		err := tm.reParseIndividualTemplate(tm.directory + "/" + name)
		if err != nil {
			logError(err.Error())
			return err
		}
	}

	tmpl, err := tm.find(name)
	if err != nil {
		logError(err.Error())
		return err
	}

	params = tm.buildParams(name, params)
	
	err = tmpl.Execute(writer, params)
	if err != nil {
		logError(err.Error())
		return err
	}

	return err
}

// Adds the default functions to the `TemplateManager` instance
func (tm *TemplateManager) addDefaultFunctions() *TemplateManager {
	tm.AddFunctions(getDefaultFunctions())

	return tm
}

// Adds a `descendant` template to the `templateName` bundle
func (tm *TemplateManager) addDescendant(templateName string, descendant string) *TemplateManager {
	if _, ok := tm.descendants[templateName]; !ok {
		tm.descendants[templateName] = []string{}
	}
	tm.descendants[templateName] = append(tm.descendants[templateName], descendant)

	return tm
}

// Readies the parameters for a single template using nested variables from descendants
// Uses a simple hierarchy to determine which should be included
func (tm *TemplateManager) buildParams(name string, params Params) Params {
	var descendants = []string{name}
	if dependencies, ok := tm.descendants[name]; ok {
		descendants = append(descendants, dependencies...)
	}

	for _, descendant := range descendants {
		if templateParams, ok := tm.params[descendant]; ok {
			for key, value := range templateParams {
				if _, ok := params[key]; !ok {
					params[key] = value
				}
			}
		}
	}

	return params
}

// Finds an individual template bundle from the `TemplateManager`
func (tm *TemplateManager) find(file string) (*template.Template, error) {
	if tmpl, ok := tm.templates[file]; ok {
		tmpl = tmpl.Lookup(file)
		if tmpl == nil {
			return nil, fmt.Errorf("template %s not found", file)
		}

		return tmpl, nil
	}

	return nil, fmt.Errorf("template %s not found", file) 
}

// Re-parses an individual template file (if reload is enabled)
func (tm *TemplateManager) reParseIndividualTemplate(path string) error {
	name, err := cleanPath(path, tm.directory)
	if err != nil {
		return err
	}

	tm.params[name]			= make(Params)
	tm.descendants[name]	= []string{}

	if tm.debug {
		logWarning(fmt.Sprintf("Re-Parsing: %s (Path: %s)\n", name, path))
	}

	tm.templates[name] = tm.configureNewTemplate(template.New(""))
	err = tm.parseFileDependencies(path, name, tm.directory, tm.templates[name])
	if err != nil {
		return err
	}

	return nil
}

// Handles parsing an individual file
func (tm *TemplateManager) parseFileDependencies(path string, name string, directory string, tmpl *template.Template) error {
	dependencies, err := tm.getFileDependencies(path, directory)
	if err != nil {
		return err
	}

	for _, dependency := range dependencies {
		dependencyName, err := cleanPath(dependency, directory)
		if err != nil {
			return err
		}

		err = tm.addTemplate(dependency, dependencyName, directory, tmpl)
		if err != nil {
			return err
		}

		tm.addDescendant(name, dependencyName)
	}

	err = tm.addTemplate(path, name, directory, tmpl)
	if err != nil {
		return err
	}
	
	if tm.debug {
		logInformation(fmt.Sprintf("Parsed template: %s (Path: %s)\n", name, path))

		if len(dependencies) > 0 {
			logInformation("\tDependencies:")
			for _, dependency := range dependencies {
				dependencyName, _ := cleanPath(dependency, directory)
				fmt.Printf("\t\tTemplate: %s (Path: %s)\n", dependencyName, dependency)
			}
		}
		
		var vars = ""
		var descendants = []string{name}
		if dependencies, ok := tm.descendants[name]; ok {
			descendants = append([]string{name}, dependencies...)
		}
		for _, descendant := range descendants {
			if templateParams, ok := tm.params[descendant]; ok {
				vars = vars + fmt.Sprintf("\t\tSet in: %s\n", descendant)
				keys := []string{}
  
				for k := range templateParams{
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					vars = vars + fmt.Sprintf("\t\t\t%s (%T) = %v\n", k, templateParams[k], templateParams[k])
				}
			}
		}
		if len(vars) > 0 {
			logInformation("\tVariables:")
			fmt.Print(vars)
		}
	}

	return nil
}

// Configures a `template.Template` instance to use the `TemplateManager` settings
func (tm *TemplateManager) configureNewTemplate(tmpl *template.Template) *template.Template {
	tmpl.Delims(tm.delimiterLeft, tm.delimiterRight)
	tmpl.Funcs(tm.functions)

	return tmpl
}

// Adds the file contents to the bundle
func (tm *TemplateManager) addTemplate(path string, name string, directory string, tmpl *template.Template) error {
	contents, err := tm.getFileContents(path, directory)
	if err != nil {
		return err
	}

	into := tm.configureNewTemplate(tmpl.New(name))

	for _, content := range contents {
		_, err := into.Parse(content)
		if err != nil {
			return err
		}
	}

	return nil
}

// Reads a file's contents and gets any extended templates too
func (tm *TemplateManager) getFileContents(path string, directory string) ([]string, error) {
	//buffer, err := os.ReadFile(path)
	buffer, err := fsWalk.ReadFile(path, tm.fileSystem)
	if err != nil {
		return []string{}, err
	}
	content  := string(buffer)
	contents := []string{}

	findVars, _	:= regexp.Compile("(?s)\\s*{{\\s*var\\s*[\"`]{1}\\s*([^\"]+)\\s*[\"`]{1}.*?}}\\s*(.*?)\\s*{{\\s*end\\s*}}\\s*")

	if findVars.Match(buffer) {
		matches := findVars.FindAllSubmatch(buffer, -1)
		name, _ := cleanPath(path, directory)

		for _, match := range matches {
			content = strings.Replace(content, string(match[0]), "", 1)
			varName := string(match[1])
			if _, ok := tm.params[name][varName]; !ok {
				tm.parseVariable(name, varName, string(match[2]))
			}
		}
		buffer = []byte(content)
	}

	findExtends, _ := regexp.Compile("^\\s*{{\\s*extends\\s*[\"`]{1}([^\"]+)[\"`]{1}.*?}}\\s*")

	if findExtends.Match(buffer) {
		matches := findExtends.FindAllSubmatch(buffer, -1)
		content = strings.Replace(content, string(matches[0][0]), "", 1)
		contents, err = tm.getFileContents(directory + "/" + string(matches[0][1]), directory)
		if err != nil {
			return []string{}, err
		}
	}

	return append(contents, content), nil
}

// Recursively finds all file dependencies
func (tm *TemplateManager) getFileDependencies(path string, directory string) ([]string, error) {
	//buffer, err := os.ReadFile(path)
	buffer, err := fsWalk.ReadFile(path, tm.fileSystem)
	if err != nil {
		return []string{}, err
	}
	dependencies := []string{}

	findExtends, _		:= regexp.Compile("^\\s*{{\\s*extends\\s*[\"`]{1}([^\"]+)[\"`]{1}.*}}\\s*")
	findTemplates, _	:= regexp.Compile("{{\\s*template\\s*[\"`]{1}([^\"]+)[\"`]{1}.*?}}")

	if findExtends.Match(buffer) {
		matches := findExtends.FindAllSubmatch(buffer, -1)
		dependencies = append(dependencies, directory + "/" + string(matches[0][1]))
	}

	if findTemplates.Match(buffer) {
		matches := findTemplates.FindAllSubmatch(buffer, -1)
		for _, match := range matches {
			dependencies = append(dependencies, directory + "/" + string(match[1]))
		}
	}

	tmp := dependencies
	for _, dependency := range tmp {
		subDependencies, err := tm.getFileDependencies(dependency, directory)
		if err != nil {
			return []string{}, err
		}

		for _, subDependency := range subDependencies {
			if !slices.Contains(dependencies, subDependency) {
				dependencies = append(dependencies, subDependency)
			}
		}
	}
	
	return dependencies, nil
}

// Parses a variable declared in a template file
// (this system is very limited to preserve actual types / avoid interfaces and reflection)
func (tm *TemplateManager) parseVariable(template string, name string, value string) {
	t, val := getVariableType(value)

	switch t {
		case "map": tm.parseVariableMap(template, name, value)
		case "slice": tm.parseVariableSlice(template, name, value)
		case "int": tm.AddParam(template, name, val.(int))
		case "float": tm.AddParam(template, name, val.(float64))
		case "bool": tm.AddParam(template, name, val.(bool))
		case "string": tm.AddParam(template, name, val.(string))
	}
}

// Parses a map variable declared in a template file
// (this system is very limited to preserve actual types / avoid interfaces and reflection)
func (tm *TemplateManager) parseVariableMap(template string, name string, value string) {
	values := prepareMap(value)

	if len(values) == 0 {
		tm.AddParam(template, name, values)
		return
	}

	k, v := "", ""
	for k, v = range values {
		break
	}

	tK, _ := getVariableType(k)
	tV, _ := getVariableType(v)

	switch tK + tV {
		case "stringstring": 
			tm.AddParam(template, name, parseVariableMap[string, string](values))
		case "stringint":
			tm.AddParam(template, name, parseVariableMap[string, int](values))
		case "stringfloat":
			tm.AddParam(template, name, parseVariableMap[string, float64](values))
		case "stringbool":
			tm.AddParam(template, name, parseVariableMap[string, bool](values))
		case "intstring": 
			tm.AddParam(template, name, parseVariableMap[int, string](values))
		case "intint":
			tm.AddParam(template, name, parseVariableMap[int, int](values))
		case "intfloat":
			tm.AddParam(template, name, parseVariableMap[int, float64](values))
		case "intbool":
			tm.AddParam(template, name, parseVariableMap[int, bool](values))
		case "floatstring": 
			tm.AddParam(template, name, parseVariableMap[float64, string](values))
		case "floatint":
			tm.AddParam(template, name, parseVariableMap[float64, int](values))
		case "floatfloat":
			tm.AddParam(template, name, parseVariableMap[float64, float64](values))
		case "floatbool":
			tm.AddParam(template, name, parseVariableMap[float64, bool](values))
		case "boolstring": 
			tm.AddParam(template, name, parseVariableMap[bool, string](values))
		case "boolint":
			tm.AddParam(template, name, parseVariableMap[bool, int](values))
		case "boolfloat":
			tm.AddParam(template, name, parseVariableMap[bool, float64](values))
		case "boolbool":
			tm.AddParam(template, name, parseVariableMap[bool, bool](values))
	}
}

// Parses a slice variable declared in a template file.
// Empty slices are string slices.
// (this system is very limited to preserve actual types / avoid interfaces and reflection)
func (tm *TemplateManager) parseVariableSlice(template string, name string, value string) {
	values := prepareSlice(value)

	if len(values) == 0 {
		tm.AddParam(template, name, values)
		return
	}

	value = values[0]
	t, _ := getVariableType(value)

	switch t {
		case "string":
			tm.AddParam(template, name, values)
		case "int":
			tm.AddParam(template, name, parseVariableSlice[int](values))
		case "float":
			tm.AddParam(template, name, parseVariableSlice[float64](values))
		case "bool":
			tm.AddParam(template, name, parseVariableSlice[bool](values))
		case "slice":
			nestedType := "string"
			if len(values[0]) > 0 {
				tmp := prepareSlice(values[0])
				nestedType, _ = getVariableType(tmp[0])
			}

			switch nestedType {
				case "int":
					tm.AddParam(template, name, parseNestedVariableSlice[int](values))
				case "float":
					tm.AddParam(template, name, parseNestedVariableSlice[float64](values))
				case "bool":
					tm.AddParam(template, name, parseNestedVariableSlice[bool](values))
				case "string":
					tm.AddParam(template, name, parseNestedVariableSlice[string](values))
			}
	}
}

// Check that the template file has a standard file extension
// Only used for automatic initialisation  
func checkAllowedExtension(file string, allowedExtensions []string) (string, error) {
	for _, extension := range allowedExtensions {
		if strings.HasSuffix(file, extension) {
			return extension, nil
		}
	}

	return "", fmt.Errorf("no match")
}

// Cleans a file path to be relative
func cleanPath(path string, directory string) (string, error) {
	file, err := filepath.Rel(directory, path)
	if err != nil {
		return "", err
	}

	return filepath.ToSlash(file), nil
}