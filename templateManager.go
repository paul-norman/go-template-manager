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
	"bytes"
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
	"strconv"
	"sync"

	"golang.org/x/exp/slices"
	"github.com/google/uuid"

	"github.com/paul-norman/go-template-manager/fsWalk"
)

// Holds all templates and variables along with all required settings
type TemplateManager struct {
	templateType			string
	templates 				map[string]*Template
	params					map[string]map[string]any
	descendants				map[string][]string
	componentDirectories	[]string
	components				map[string]string
	delimiterLeft			string
	delimiterRight			string
	fileSystem				http.FileSystem
	directory				string
	extension				string
	excludedDirectories		[]string
	functions				map[string]any
	mutex					sync.RWMutex
	debug					bool
	reload					bool
	parsed					bool
}

// Convenience type allowing any variables types to be passed in
type Params map[string]any

// Allow regexps to be pre-compiled
var regexps map[string]*regexp.Regexp

// Creates a new `TemplateManager` struct instance
func Init(directory string, extension string) *TemplateManager {
	templateManager := &TemplateManager{
		templateType:			"text",
		templates:				make(map[string]*Template),
		params:					make(map[string]map[string]any),
		descendants:			make(map[string][]string),
		componentDirectories:	[]string{"components"},
		components:				make(map[string]string),
		delimiterLeft:			"{{",
		delimiterRight:			"}}",
		directory:				directory,
		extension:				extension,
		excludedDirectories:	[]string{"layouts", "partials", "components"},
		functions:				make(map[string]any),
		debug:					false,
		reload:					false,
		parsed:					false,
	}

	initRegexps()
	templateManager.initRegexps()
	templateManager.addDefaultFunctions()

	return templateManager
}

// Creates a new `TemplateManager` struct instance using an embedded filesystem
func InitEmbed(fileSystem embed.FS, directory string, extension string) *TemplateManager {
	templateManager := &TemplateManager{
		templateType:			"text",
		templates:				make(map[string]*Template),
		params:					make(map[string]map[string]any),
		descendants:			make(map[string][]string),
		componentDirectories:	[]string{"components"},
		components:				make(map[string]string),
		delimiterLeft:			"{{",
		delimiterRight:			"}}",
		directory:				"/" + strings.TrimLeft(directory, " /"),
		fileSystem:				http.FS(fileSystem),
		extension:				extension,
		excludedDirectories:	[]string{"layouts", "partials", "components"},
		functions:				make(map[string]any),
		debug:					false,
		reload:					false,
		parsed:					false,
	}

	initRegexps()
	templateManager.initRegexps()
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
				templateDirectory	= strings.Split(name, "/")[0]
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

// Adds multiple directories that contain components.
// (Must be within the templates directory)
func (tm *TemplateManager) AddComponentDirectories(directories []string) *TemplateManager {
	for _, directory := range directories {
		tm.AddComponentDirectory(directory)
	}

	return tm
}

// Adds a directory that contains components.
// (Must be within the templates directory)
func (tm *TemplateManager) AddComponentDirectory(directory string) *TemplateManager {
	if !slices.Contains(tm.componentDirectories, directory) {
		tm.mutex.Lock()
		tm.componentDirectories = append(tm.componentDirectories, directory)
		tm.mutex.Unlock()
	}

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

	consoleErrors	= true
	consoleWarnings	= true

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
		tm.mutex.Lock()
		tm.excludedDirectories = append(tm.excludedDirectories, directory)
		tm.mutex.Unlock()
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
		logWarning("Parsing all components...")
	}

	tm.parseComponents()

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

		tm.templates[name] = tm.configureNewTemplate(NewTemplate(tm.templateType, ""))
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

// Removes multiple directories that contain components from component parsing.
func (tm *TemplateManager) RemoveComponentDirectories(directories []string) *TemplateManager {
	for _, directory := range directories {
		tm.RemoveComponentDirectory(directory)
	}

	return tm
}

// Removes a directory that contains components from component parsing.
func (tm *TemplateManager) RemoveComponentDirectory(directory string) *TemplateManager {
	if slices.Contains(tm.componentDirectories, directory) {
		tm.mutex.Lock()
		for k, v := range tm.componentDirectories {
			if v == directory {
				tm.componentDirectories = append(tm.componentDirectories[:k], tm.componentDirectories[k + 1:]...)
				break
			}
		}
		tm.mutex.Unlock()
	}

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
	for name := range names {
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
			err = logError(err.Error())
			return err
		}
	}

	if tm.reload {
		err := tm.reParseIndividualTemplate(tm.directory + "/" + name)
		if err != nil {
			err = logError(err.Error())
			return err
		}
	}

	tmpl, err := tm.find(name)
	if err != nil {
		err = logError(err.Error())
		return err
	}

	params = tm.buildParams(name, params)
	
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, params)
	if err != nil {
		err = logError("FATAL: " + err.Error())
		return err
	}

	buf.WriteTo(writer)
	return nil
}

// Sets whether `TemplateManager` should use the `text/template` package or the `html/template` package
func (tm *TemplateManager) TemplateEngine(engine string) *TemplateManager {
	engine = strings.ToLower(engine)
	if engine == "text" || engine == "text/template" {
		tm.templateType = "text"
	} else if engine == "html" || engine == "html/template" {
		tm.templateType = "html"
	} else {
		panic("invalid template engine chosen: " + engine)
	}

	return tm
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
func (tm *TemplateManager) find(file string) (*Template, error) {
	if tmpl, ok := tm.templates[file]; ok {
		tmpl = tmpl.Lookup(file)
		if tmpl == nil {
			return nil, fmt.Errorf("template %s not found", file)
		}

		return tmpl, nil
	}

	return nil, fmt.Errorf("template %s not found", file) 
}

// Initialises all regexps required for detected components
func (tm *TemplateManager) initComponentRegexps() {
	if len(tm.components) > 0 {
		findGeneratedDefines, _		:= regexp.Compile(`(?s)` + tm.delimiterLeft + `- define "content-([^"]{36})" -` + tm.delimiterRight + `(.+?)` + tm.delimiterLeft + `- end -` + tm.delimiterRight)
		findCollectionComponents, _	:= regexp.Compile(`(?s)` + tm.delimiterLeft + ` ([^ ]+) x\-(render "[^ ]+-([^"]{36})" \(collection "ComponentUuid" .*?)` + tm.delimiterRight)

		regexps["findGeneratedDefines"]		= findGeneratedDefines
		regexps["findCollectionComponents"] = findCollectionComponents

		for component := range tm.components {
			findComponentsDouble, _				:= regexp.Compile(`(?s)<` + component + `(\s+[^>]*)?\s*>(.*?)</` + component + `>`)
			findComponentsSingle, _				:= regexp.Compile(`(?s)<` + component + `(\s+[^>]*)?\s*>`)
			findComponentsCollectedDouble, _	:= regexp.Compile(`(?s)<x-` + component + `(\s+[^>]*)?\s*>(.*?)</x-` + component + `>`)
			findComponentsCollectedSingle, _	:= regexp.Compile(`(?s)<x-` + component + `(\s+[^>]*)?\s*>`)

			regexps[component + "_findComponentsDouble"]				= findComponentsDouble
			regexps[component + "_findComponentsSingle"]				= findComponentsSingle
			regexps[component + "_findComponentsCollectedDouble"]	= findComponentsCollectedDouble
			regexps[component + "_findComponentsCollectedSingle"]	= findComponentsCollectedSingle
		}
	}
}

// Initialises the regexps required by the file scanning
func (tm *TemplateManager) initRegexps() {
	findVars, _					:= regexp.Compile("(?s)\\s*" + tm.delimiterLeft + "(?:- )?(?:\\/\\*)?\\s*var\\s*[\"`]{1}\\s*([^\"]+)\\s*[\"`]{1}.*?" + tm.delimiterRight + "\\s*(.*?)\\s*" + tm.delimiterLeft + "\\s*end\\s*(?:\\*\\/)?(?: -)?" + tm.delimiterRight + "\\s*")
	findExtends, _				:= regexp.Compile("^\\s*" + tm.delimiterLeft + "(?:- )?(?:\\/\\*)?\\s*extends\\s*[\"`]{1}([^\"`]+)[\"`]{1}\\s*(?:\\*\\/)?(?: -)?" + tm.delimiterRight + "\\s*")
	findTemplates, _			:= regexp.Compile(tm.delimiterLeft + "\\-?\\s*template\\s*[\"`]{1}([^\"`]+)[\"`]{1}.*?\\-?" + tm.delimiterRight)

	regexps["findVars"]			= findVars
	regexps["findExtends"]		= findExtends
	regexps["findTemplates"]	= findTemplates
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

	tm.templates[name] = tm.configureNewTemplate(NewTemplate(tm.templateType, ""))
	err = tm.parseFileDependencies(path, name, tm.directory, tm.templates[name])
	if err != nil {
		return err
	}

	return nil
}

func (tm *TemplateManager) parseComponents() error {
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

		componentPath, _ := cleanPath(path, tm.directory)
		index := strings.LastIndex(componentPath, "/")
		name := stripExtension(componentPath[index + 1:])
		tm.components[name] = componentPath

		return nil
	}

	for _, componentDirectory := range tm.componentDirectories {
		if tm.fileSystem != nil {
			err = fsWalk.WalkDir(tm.fileSystem, "/" + componentDirectory, walk)
		} else {
			err = filepath.WalkDir(tm.directory + "/" + componentDirectory, walk)
		}
	}

	tm.initComponentRegexps()

	if err == nil {
		tm.parsed = true

		if tm.debug {
			logSuccess("All components parsed and ready to use")
		}
	}

	return err
}

// Handles parsing an individual file
func (tm *TemplateManager) parseFileDependencies(path string, name string, directory string, tmpl *Template) error {
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
func (tm *TemplateManager) configureNewTemplate(tmpl *Template) *Template {
	tmpl.Delims(tm.delimiterLeft, tm.delimiterRight)
	tmpl.Funcs(tm.functions)
	tmpl.Funcs(map[string]any {
		"render": func(name string, args ...any) string {
			var data any = nil
			if len(args) > 0 {
				data = args[0]
			}
			buf := &bytes.Buffer{}
			err := tmpl.ExecuteTemplate(buf, name, data)
			if err != nil {
				return ""
			}
			return buf.String()
        },
	})

	return tmpl
}

// Adds the file contents to the bundle
func (tm *TemplateManager) addTemplate(path string, name string, directory string, tmpl *Template) error {
	contents, err := tm.getFileContents(path, directory)
	if err != nil {
		return err
	}

	into := tm.configureNewTemplate(tmpl.NewSubTemplate(name))

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
	buffer, err := fsWalk.ReadFile(path, tm.fileSystem)
	if err != nil {
		return []string{}, err
	}
	content  := string(buffer)
	contents := []string{}

	if regexps["findVars"].MatchString(content) {
		matches := regexps["findVars"].FindAllStringSubmatch(content, -1)
		name, _ := cleanPath(path, directory)

		for _, match := range matches {
			content = strings.Replace(content, match[0], "", 1)
			varName := string(match[1])
			if _, ok := tm.params[name][varName]; !ok {
				tm.parseVariable(name, varName, match[2])
			}
		}
	}

	if regexps["findExtends"].MatchString(content) {
		matches := regexps["findExtends"].FindAllStringSubmatch(content, -1)
		content = strings.Replace(content, matches[0][0], "", 1)
		contents, err = tm.getFileContents(directory + "/" + matches[0][1], directory)
		if err != nil {
			return []string{}, err
		}
	}

	content = tm.parseContentComponents(content, directory)

	return append(contents, content), nil
}

// Recursively finds all file dependencies
func (tm *TemplateManager) getFileDependencies(path string, directory string) ([]string, error) {
	buffer, err := fsWalk.ReadFile(path, tm.fileSystem)
	if err != nil {
		return []string{}, err
	}
	dependencies := []string{}

	if regexps["findExtends"].Match(buffer) {
		matches := regexps["findExtends"].FindAllSubmatch(buffer, -1)
		dependencies = append(dependencies, directory + "/" + string(matches[0][1]))
	}

	if regexps["findTemplates"].Match(buffer) {
		matches := regexps["findTemplates"].FindAllSubmatch(buffer, -1)
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

func (tm *TemplateManager) parseContentComponents(content string, directory string) string {
	if len(tm.components) > 0 {
		for component, componentPath := range tm.components {
			if strings.Contains(content, "<" + component) {	
				matches := [][]string{}

				if regexps[component + "_findComponentsDouble"].MatchString(content) {
					matches = regexps[component + "_findComponentsDouble"].FindAllStringSubmatch(content, -1)
				} else {
					if regexps[component + "_findComponentsSingle"].MatchString(content) {
						matches = regexps[component + "_findComponentsSingle"].FindAllStringSubmatch(content, -1)
					}
				}

				for _, match := range matches {
					random_id	:= uuid.NewString()
					find		:= match[0]
					replace		:= tm.delimiterLeft + ` block "` + componentPath + `-` + random_id + `"`
					
					attributes := match[1]
					tagContent := ""
					if len(match) > 2 {
						tagContent = match[2]
					}
					
					replace += ` collection "ComponentUuid" "` + random_id + `" "ComponentContent" "content-` + random_id + `"`
					if len(attributes) > 0 {
						attributes := regexps["findAttributes"].FindAllStringSubmatch(attributes, -1)
						for _, attribute := range attributes {
							bias := "numeric"
							if strings.HasPrefix(attribute[2], `"`) {
								bias = "string"
							}

							value := strings.Trim(attribute[2], `"`)
							if bias == "string" {
								if strings.HasPrefix(value, tm.delimiterLeft) && strings.HasSuffix(value, tm.delimiterRight) {
									value = strings.TrimRight(strings.TrimLeft(value, tm.delimiterLeft + " "), tm.delimiterRight + " ")
								} else {
									value = `"` + value + `"`
								}
							}

							replace += ` "` + strings.Trim(attribute[1], " ") + `" ` + value
						}
					} else {
						replace += ` "Null" ""`
					}

					if len(tagContent) > 0 {
						// {{- define "content-RANDOM_ID" -}} passed content {{- end -}}
						tagContent = tm.delimiterLeft + `- define "content-` + random_id + `" -` + tm.delimiterRight + tagContent + tm.delimiterLeft + `- end -` + tm.delimiterRight
					}

					componentContents, err := tm.getFileContents(directory + "/" + componentPath, directory)
					if err != nil {
						continue
					}
					componentContent := componentContents[0] // TODO - this is wrong, will fail if extended. Cannot extend?

					replace += ` -` + tm.delimiterRight + componentContent + tm.delimiterLeft + `- end ` + tm.delimiterRight

					content = tagContent + strings.Replace(content, find, replace, 1)
				}
			}
			if strings.Contains(content, "<x-" + component) {	
				matches := [][]string{}

				if regexps[component + "_findComponentsCollectedDouble"].MatchString(content) {
					matches = regexps[component + "_findComponentsCollectedDouble"].FindAllStringSubmatch(content, -1)
				} else {
					if regexps[component + "_findComponentsCollectedSingle"].MatchString(content) {
						matches = regexps[component + "_findComponentsCollectedSingle"].FindAllStringSubmatch(content, -1)
					}
				}
				for _, match := range matches {
					random_id	:= uuid.NewString()
					find		:= match[0]
					define		:= tm.delimiterLeft + ` define "` + componentPath + `-` + random_id + `"`
					
					attributes := match[1]
					tagContent := ""
					if len(match) > 2 {
						tagContent = match[2]
					}
					
					create := `(collection "ComponentUuid" "` + random_id + `" "ComponentContent" "content-` + random_id + `"`
					if len(attributes) > 0 {
						attributes := regexps["findAttributes"].FindAllStringSubmatch(attributes, -1)
						for _, attribute := range attributes {
							bias := "numeric"
							if strings.HasPrefix(attribute[2], `"`) {
								bias = "string"
							}

							value := strings.Trim(attribute[2], `"`)
							if bias == "string" {
								if strings.HasPrefix(value, tm.delimiterLeft) && strings.HasSuffix(value, tm.delimiterRight) {
									value = strings.TrimRight(strings.TrimLeft(value, tm.delimiterLeft + " "), tm.delimiterRight + " ")
								} else {
									value = `"` + value + `"`
								}
							}

							create += ` "` + strings.Trim(attribute[1], " ") + `" ` + value
						}
					} else {
						create += ` "Null" ""`
					}

					if len(tagContent) > 0 {
						// {{- define "content-RANDOM_ID" -}} passed content {{- end -}}
						tagContent = tm.delimiterLeft + `- define "content-` + random_id + `" -` + tm.delimiterRight + tagContent + tm.delimiterLeft + `- end -` + tm.delimiterRight
					}

					componentContents, err := tm.getFileContents(directory + "/" + componentPath, directory)
					if err != nil {
						continue
					}
					componentContent := componentContents[0] // TODO - this is wrong, will fail if extended. Cannot extend?

					define += ` -` + tm.delimiterRight + componentContent + tm.delimiterLeft + `- end ` + tm.delimiterRight

					content = tagContent + define + strings.Replace(content, find, tm.delimiterLeft + ` ` + component + ` x-render "` + componentPath + `-` + random_id + `" ` + create + `) ` + tm.delimiterRight, 1)
				}
			}
		}

		// Collect nested components
		if strings.Contains(content, `x-render "`) {			
			if regexps["findGeneratedDefines"].MatchString(content) {
				matches := regexps["findGeneratedDefines"].FindAllStringSubmatch(content, -1)
				for _, match := range matches {
					// Individual component content definition to act upon
					if strings.Contains(match[2], `x-render`) {
						collectedVars	:= map[string][]string{}
						wholeDefine		:= match[0]
						newDefine		:= match[0]
						defineId		:= match[1]
						defineContents	:= match[2]

						if regexps["findCollectionComponents"].MatchString(defineContents) {
							submatches := regexps["findCollectionComponents"].FindAllStringSubmatch(defineContents, -1)
							for _, submatch := range submatches {
								wholeComponent	:= submatch[0]
								componentName	:= submatch[1]
								componentRender	:= submatch[2]

								if _, ok := collectedVars[componentName]; !ok {
									collectedVars[componentName] = []string{}
								}
								collectedVars[componentName] = append(collectedVars[componentName], strings.Replace(componentRender, `collection "ComponentUuid"`, `collection "ParentUuid" "` + defineId + `" "ParentPosition" ` + strconv.Itoa(len(collectedVars[componentName])) + ` "ComponentUuid"`, 1))

								newDefine = strings.Replace(newDefine, wholeComponent, "", 1)
							}
						}

						content = strings.Replace(content, wholeDefine, newDefine, 1)

						// Add collected variables to the parent template definition
						findParentDefine, _ := regexp.Compile(`(?s)` + tm.delimiterLeft + ` block ".+?-` + defineId + `" (.*?) \-` + tm.delimiterRight)
						if findParentDefine.MatchString(content) {
							submatches := findParentDefine.FindAllStringSubmatch(content, -1)
							for _, submatch := range submatches {
								wholeMaster	:= submatch[0]
								collected	:= submatch[1]

								for varName, slice := range collectedVars {
									collected += ` "` + varName + `" (list`
									for _, value := range slice {
										collected += ` (` + value + `)`
									}
									collected += `)`
								}
								newMaster := strings.Replace(wholeMaster, submatch[1], collected, 1)

								content = strings.Replace(content, wholeMaster, newMaster, 1)
							}
						}
					}
				}
			}
		}
	}

	return content
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

func initRegexps() {
	findAttributes, _			:= regexp.Compile(`(?s)([^=\s]+)\s*=\s*("[^"]+"|[\d\.\-]+)`)
	findNumericSlice, _			:= regexp.Compile(`\s*([\-\d\.]+)\s*,`)
	findBooleanSlice, _			:= regexp.Compile(`(?i)\s*(true|false)\s*,`)
	findStringSlice, _			:= regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")
	findSliceSlice, _			:= regexp.Compile(`(\[.*?\])\s*,`)
	findBoolBoolMap, _			:= regexp.Compile(`(?i)(true|false)\s*:\s*(true|false)\s*,`)
	findBoolNumericMap, _		:= regexp.Compile(`(?i)(true|false)\s*:\s*([\-\.\d]+)\s*,`)
	findBoolStringMap, _		:= regexp.Compile("(?i)(true|false)\\s*:\\s*[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")
	findNumericBoolMap, _		:= regexp.Compile(`(?i)([\-\.\d]+)\s*:\s*(true|false)\s*,`)
	findNumericNumericMap, _	:= regexp.Compile(`([\-\.\d]+)\s*:\s*([\-\.\d]+)\s*,`)
	findNumericStringMap, _		:= regexp.Compile("([\\-\\.\\d]+)\\s*:\\s*[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")
	findStringBoolMap, _		:= regexp.Compile("(?i)[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*(true|false)\\s*,")
	findStringNumericMap, _		:= regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*([\\-\\.\\d]+)\\s*,")
	findStringStringMap, _		:= regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*,")
	findSliceMap, _				:= regexp.Compile("[\"`']{1}(.*?[^\\\\])[\"`']{1}\\s*:\\s*(\\[.*?\\])\\s*,")

	regexps = map[string]*regexp.Regexp{
		"findAttributes":			findAttributes,

		"findNumericSlice":			findNumericSlice,
		"findBooleanSlice":			findBooleanSlice,
		"findStringSlice":			findStringSlice,
		"findSliceSlice":			findSliceSlice,
		"findBoolBoolMap":			findBoolBoolMap,
		"findBoolNumericMap":		findBoolNumericMap,
		"findBoolStringMap":		findBoolStringMap,
		"findNumericBoolMap":		findNumericBoolMap, 
		"findNumericNumericMap":	findNumericNumericMap,
		"findNumericStringMap":		findNumericStringMap,
		"findStringBoolMap":		findStringBoolMap, 
		"findStringNumericMap":		findStringNumericMap,
		"findStringStringMap":		findStringStringMap,
		"findSliceMap":				findSliceMap,
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

// Removes the extension from a file
func stripExtension(file string) string {
	index := strings.LastIndex(file, ".")
	return file[:index] 
}