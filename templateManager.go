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
	//"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"golang.org/x/exp/slices"
)

type TemplateManager struct {
	templates 			map[string]*template.Template
	params				map[string]map[string]any
	descendants			map[string][]string
	delimiterLeft		string
	delimiterRight		string
	directory			string
	extension			string
	excludedDirectories	[]string
	layout				string
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
		layout:					"embed",
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

// Enable re-rebuilding of the template bundle upon every page load (for development)
func (tm *TemplateManager) Reload(reload bool) *TemplateManager {
	tm.reload = reload

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

// Removes all functions currently assigned to the instance of `TemplateManager`.
// Useful if you do not want the default functions included
func (tm *TemplateManager) RemoveAllFunctions() *TemplateManager {
	tm.mutex.Lock()
	tm.functions = make(map[string]any)
	tm.mutex.Unlock()

	return tm
}

// Triggers scanning of files and bundling of all templates
func (tm *TemplateManager) Parse() error {
	if tm.parsed {
		return nil
	}

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if tm.debug {
		logWarning("Parsing all templates...")
	}

	err := filepath.WalkDir(tm.directory, func(path string, info fs.DirEntry, err error) error {
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
	})

	if err == nil {
		tm.parsed = true

		if tm.debug {
			logSuccess("All templates parsed and ready to use")
		}
	}

	return err
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
	dependencies, err := getFileDependencies(path, directory)
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
	buffer, err := os.ReadFile(path)
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

	k := "" 
	v := ""
	for k, v = range values {
		break
	}

	tK, _ := getVariableType(k)
	tV, _ := getVariableType(v)

	switch tK + tV {
		case "stringstring": 
			tmp := map[string]string{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(string)] = val.(string)
			}
			tm.AddParam(template, name, tmp)
		case "stringint":
			tmp := map[string]int{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(string)] = val.(int)
			}
			tm.AddParam(template, name, tmp)
		case "stringfloat":
			tmp := map[string]float64{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(string)] = val.(float64)
			}
			tm.AddParam(template, name, tmp)
		case "stringbool":
			tmp := map[string]bool{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(string)] = val.(bool)
			}
			tm.AddParam(template, name, tmp)
		case "intstring": 
			tmp := map[int]string{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(int)] = val.(string)
			}
			tm.AddParam(template, name, tmp)
		case "intint":
			tmp := map[int]int{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(int)] = val.(int)
			}
			tm.AddParam(template, name, tmp)
		case "intfloat":
			tmp := map[int]float64{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(int)] = val.(float64)
			}
			tm.AddParam(template, name, tmp)
		case "intbool":
			tmp := map[int]bool{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(int)] = val.(bool)
			}
			tm.AddParam(template, name, tmp)
		case "floatstring": 
			tmp := map[float64]string{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(float64)] = val.(string)
			}
			tm.AddParam(template, name, tmp)
		case "floatint":
			tmp := map[float64]int{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(float64)] = val.(int)
			}
			tm.AddParam(template, name, tmp)
		case "floatfloat":
			tmp := map[float64]float64{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(float64)] = val.(float64)
			}
			tm.AddParam(template, name, tmp)
		case "floatbool":
			tmp := map[float64]bool{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(float64)] = val.(bool)
			}
			tm.AddParam(template, name, tmp)
		case "boolstring": 
			tmp := map[bool]string{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(bool)] = val.(string)
			}
			tm.AddParam(template, name, tmp)
		case "boolint":
			tmp := map[bool]int{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(bool)] = val.(int)
			}
			tm.AddParam(template, name, tmp)
		case "boolfloat":
			tmp := map[bool]float64{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(bool)] = val.(float64)
			}
			tm.AddParam(template, name, tmp)
		case "boolbool":
			tmp := map[bool]bool{}
			for key, val := range values {
				_, key := getVariableType(key)
				_, val := getVariableType(val)
				tmp[key.(bool)] = val.(bool)
			}
			tm.AddParam(template, name, tmp)
	}
}

// Parses a slice variable declared in a template file.
// Empty slices are string slices.
// (this system is very limited to preserve actual types / avoid interfaces and reflection)
func (tm *TemplateManager) parseVariableSlice(template string, name string, value string) {
	values, _ := prepareSlice(value)

	if len(values) == 0 {
		tm.AddParam(template, name, values)
		return
	}

	value = values[0]
	t, _ := getVariableType(value)

	switch t {
		case "int": 
			tmp := []int{}
			for _, val := range values {
				val, err := strconv.Atoi(val)
				if err != nil {
					logError(err.Error())
					return	
				}
				tmp = append(tmp, val)
			}
			tm.AddParam(template, name, tmp)
		case "float":
			tmp := []float64{}
			for _, val := range values {
				val, err := strconv.ParseFloat(val, 64)
				if err != nil {
					logError(err.Error())
					return	
				}
				tmp = append(tmp, val)
			}
			tm.AddParam(template, name, tmp)
		case "bool":
			tmp := []bool{}
			for _, val := range values {
				val, err := strconv.ParseBool(val)
				if err != nil {
					logError(err.Error())
					return	
				}
				tmp = append(tmp, val)
			}
			tm.AddParam(template, name, tmp)
		case "slice":
			nestedType := "string"
			if len(values[0]) > 0 {
				tmp, _ := prepareSlice(values[0])
				nestedType, _ = getVariableType(tmp[0])
			}

			switch nestedType {
				case "int":
					tmp := [][]int{}
					for _, val := range values {
						sub := []int{}
						tmp2, _ := prepareSlice(val)
						for _, subval := range tmp2 {
							subval, err := strconv.Atoi(subval)
							if err != nil {
								logError(err.Error())
								return	
							}
							sub = append(sub, subval)
						}
						tmp = append(tmp, sub)
					}
					tm.AddParam(template, name, tmp)
				case "float":
					tmp := [][]float64{}
					for _, val := range values {
						sub := []float64{}
						tmp2, _ := prepareSlice(val)
						for _, subval := range tmp2 {
							subval, err := strconv.ParseFloat(subval, 64)
							if err != nil {
								logError(err.Error())
								return	
							}
							sub = append(sub, subval)
						}
						tmp = append(tmp, sub)
					}
					tm.AddParam(template, name, tmp)
				case "bool":
					tmp := [][]bool{}
					for _, val := range values {
						sub := []bool{}
						tmp2, _ := prepareSlice(val)
						for _, subval := range tmp2 {
							subval, err := strconv.ParseBool(subval)
							if err != nil {
								logError(err.Error())
								return	
							}
							sub = append(sub, subval)
						}
						tmp = append(tmp, sub)
					}
					tm.AddParam(template, name, tmp)
				case "string":
					tmp := [][]string{}
					for _, val := range values {
						sub := []string{}
						tmp2, _ := prepareSlice(val)
						sub = append(sub, tmp2...)
						tmp = append(tmp, sub)
					}
					tm.AddParam(template, name, tmp)
			}

		case "string":
			tm.AddParam(template, name, values)
	}
}

// Cleans a file path to be relative
func cleanPath(path string, directory string) (string, error) {
	file, err := filepath.Rel(directory, path)
	if err != nil {
		return "", err
	}

	return filepath.ToSlash(file), nil
}

// Recursively finds all file dependencies
func getFileDependencies(path string, directory string) ([]string, error) {
	buffer, err := os.ReadFile(path)
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
		subDependencies, err := getFileDependencies(dependency, directory)
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

// Determines a variable's (likely) basic type from a string representation of it
func getVariableType(value string) (string, any) {
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		return "map", nil
	} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		return "slice", nil
	} else if val, err := strconv.Atoi(value); err == nil {
		return "int", val
	} else if val, err := strconv.ParseFloat(value, 64); err == nil {
		return "float", val
	} else if val, err := strconv.ParseBool(value); err == nil {
		return "bool", val
	}
	
	return "string", value
}

// Parses a string representation of a map into string values ready for type detection
func prepareMap(value string) map[string]string {
	value = value[1:len(value) - 1]
	value = value + ","
	m := make(map[string]string)

	findStringMap, _ := regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")

	if findStringMap.MatchString(value) {
		matches := findStringMap.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			m[match[1]] = match[2] 
		}
	} else {
		findNumericMap, _ := regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*([\\-\\.\\d]+)\\s*,")

		if findNumericMap.MatchString(value) {
			matches := findNumericMap.FindAllStringSubmatch(value, -1)
			for _, match := range matches {
				m[match[1]] = match[2] 
			}
		}
		/* 
		else {
			findSliceMap, _ := regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*(\\[.*?\\])\\s*,")
	
			if findSliceMap.MatchString(value) {
				matches := findSliceMap.FindAllStringSubmatch(value, -1)
				for _, match := range matches {
					m[match[1]] = match[2] 
				}
			}
		}
		*/
	}

	return m
}

// Parses a string representation of a slice into a slice of string values ready for type detection
func prepareSlice(value string) ([]string, error) {
	value = value[1:len(value) - 1]
	var values = []string{}

	if value[0:1] == "[" {
		values, _ = prepareSliceSlice(value)
	} else if strings.Contains(value, `"`) || strings.Contains(value, `'`) || strings.Contains(value, "`") {
		values, _ = prepareStringSlice(value)
	} else {
		values, _ = prepareNumericSlice(value)
	}

	return values, nil
}

// Parses a string representation of a slice of slices into a slice of strings ready for type detection
func prepareSliceSlice(value string) ([]string, error) {
	value = value + ","
	slice := []string{}

	findSliceSlice, _ := regexp.Compile(`(\[.*?\])\s*,`)

	if findSliceSlice.MatchString(value) {
		matches := findSliceSlice.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			slice = append(slice, match[1])
		}
	}

	return slice, nil
}

// Parses a string representation of a slice of strings into an actual slice of strings ready for type detection
func prepareStringSlice(value string) ([]string, error) {
	value = value + ","
	slice := []string{}

	// No backreferences in GoLang's RE2 regexp engine :-(
	findStringSlice, _	:= regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")

	if findStringSlice.MatchString(value) {
		matches := findStringSlice.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			val := strings.Replace(match[1], "\\`", "`", -1)
			val = strings.Replace(val, `\"`, `"`, -1)
			val = strings.Replace(val, "\\'", "'", -1)

			slice = append(slice, val)
		}
	}

	return slice, nil
}

// Parses a string representation of a slice of numbers into a slice of strings ready for type detection
func prepareNumericSlice(value string) ([]string, error) {
	value = value + ","
	slice := []string{}

	findNumericSlice, _	:= regexp.Compile(`\s*([\-\d\.]+)\s*,`)
	if findNumericSlice.MatchString(value) {
		matches := findNumericSlice.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			slice = append(slice, match[1])
		}
	}

	return slice, nil
}

// Logs error messages
func logError(err string) {
	if logErrors {
		fmt.Println("\033[31m" + err + "\033[0m")
	}
}

// Logs warning messages
func logWarning(warning string) {
	if logWarnings {
		fmt.Println("\033[33m" + warning + "\033[0m")
	}
}

// Logs informational messages
func logInformation(information string) {
	fmt.Println("\033[36m" + information + "\033[0m")
}

// Logs success messages
func logSuccess(success string) {
	fmt.Println("\033[32m" + success + "\033[0m")
}