package rendering

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

type GoTemplateRendererConfig struct {
	AdditionalFuncs   template.FuncMap
	TemplateDir       string
	TemplateExtension string
	TemplateFS        fs.FS
	LayoutsDir        string
}

type GoTemplateRenderer struct {
	funcs             template.FuncMap
	templateDir       string
	templateExtension string
	templateFS        fs.FS
	layoutsDir        string

	templates map[string]*template.Template
}

func NewGoTemplateRenderer(config GoTemplateRendererConfig) *GoTemplateRenderer {
	var (
		err error
	)

	ext := config.TemplateExtension

	if ext == "" {
		ext = ".html"
	}

	funcs := getFuncs(config.AdditionalFuncs)
	tmpl := template.New("").Funcs(funcs)
	templates := map[string]*template.Template{}

	/*
	 * Process layouts first
	 */
	err = fs.WalkDir(config.TemplateFS, filepath.Join(config.TemplateDir, config.LayoutsDir), func(path string, d fs.DirEntry, err error) error {
		var (
			relativePath string
			content      []byte
		)

		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, normalizeTemplateExt(ext)) {
			return nil
		}

		if relativePath, err = filepath.Rel(config.TemplateDir, path); err != nil {
			return err
		}

		templateName := strings.TrimSuffix(relativePath, ext)

		if content, err = fs.ReadFile(config.TemplateFS, path); err != nil {
			return err
		}

		template.Must(tmpl.New(templateName).Parse(string(content)))
		slog.Debug("parsed layout", "templateName", templateName, "path", path)
		return nil
	})

	if err != nil {
		slog.Error("error parsing layouts. shutting down.", "error", err, "templateDir", config.TemplateDir, "ext", ext)
		os.Exit(1)
	}

	err = fs.WalkDir(config.TemplateFS, config.TemplateDir, func(path string, d fs.DirEntry, err error) error {
		var (
			relativePath string
			content      []byte
			tt           *template.Template
		)

		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, normalizeTemplateExt(ext)) {
			return nil
		}

		// Don't re-parse layouts
		if strings.HasPrefix(path, filepath.Join(config.TemplateDir, config.LayoutsDir)) {
			return nil
		}

		if relativePath, err = filepath.Rel(config.TemplateDir, path); err != nil {
			return err
		}

		templateName := strings.TrimSuffix(relativePath, ext)

		if content, err = fs.ReadFile(config.TemplateFS, path); err != nil {
			return err
		}

		tt = template.Must(template.Must(tmpl.Clone()).New(templateName).Parse(string(content)))
		templates[templateName] = tt
		slog.Debug("parsed template", "templateName", templateName, "path", path)

		return nil
	})

	if err != nil {
		slog.Error("error parsing templates. shutting down.", "error", err, "templateDir", config.TemplateDir, "ext", ext)
		os.Exit(1)
	}

	result := &GoTemplateRenderer{
		funcs:             funcs,
		templateFS:        config.TemplateFS,
		templateExtension: normalizeTemplateExt(ext),
		templateDir:       normalizeTemplateDir(config.TemplateDir),
		layoutsDir:        config.LayoutsDir,
		templates:         templates,
	}

	return result
}

/*
Render renders a Go template file into a layout template file using the provided
data to an io.Writer.
*/
func (tr *GoTemplateRenderer) Render(templateName string, data any, w io.Writer) {
	var (
		err error
		t   *template.Template
		ok  bool
	)

	if t, ok = tr.templates[templateName]; !ok {
		slog.Error("cannot find template", "error", err, "templateName", templateName)
		fmt.Fprintf(w, "cannot find template '%s'", templateName)
		return
	}

	if err = t.ExecuteTemplate(w, templateName, data); err != nil {
		slog.Error("error executing template", "error", err, "templateName", templateName)
		fmt.Fprintf(w, "cannot execute template '%s': %s", templateName, err.Error())
		return
	}
}

/*
RenderString renders a Go template string with a set of data to an io.Writer.
*/
func (tr *GoTemplateRenderer) RenderString(templateString string, data any, w io.Writer) {
	var (
		err  error
		tmpl *template.Template
	)

	if tmpl, err = template.New("raw").Funcs(tr.funcs).Parse(templateString); err != nil {
		slog.Error("error parsing template", "error", err)
		fmt.Fprintf(w, "error parsing template: %s", err.Error())
		return
	}

	if err = tmpl.Execute(w, data); err != nil {
		slog.Error("error executing template", "error", err)
		fmt.Fprintf(w, "error executing template: %s", err.Error())
	}
}

func getFuncs(additionalFuncs template.FuncMap) template.FuncMap {
	templateFuncs := template.FuncMap{
		"join":                join,
		"isSet":               templateFuncIsSet,
		"isLastItem":          isLastItem,
		"containsString":      containsString,
		"stringSliceContains": sliceContains[string],
		"uintSliceContains":   sliceContains[uint],
		"stringNotEmpty":      stringNotEmpty,
		"javascriptIncludes":  javascriptIncludes,
		"stylesheetIncludes":  stylesheetIncludes,
	}

	if additionalFuncs != nil {
		for k, v := range additionalFuncs {
			templateFuncs[k] = v
		}
	}

	return templateFuncs
}

func templateFuncIsSet(name string, data any) bool {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	return v.FieldByName(name).IsValid()
}

func sliceContains[T comparable](array []T, value T) bool {
	return slices.Index(array, value) > -1
}

func isLastItem(index, length int) bool {
	return index == length-1
}

func containsString(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}

	return false
}

func stringNotEmpty(s any) bool {
	s1, ok := s.(string)

	if ok {
		return strings.TrimSpace(s1) != ""
	}

	s2, ok := s.(template.HTML)

	if ok {
		return strings.TrimSpace(string(s2)) != ""
	}

	return false
}

func javascriptIncludes(keyName string, data any) template.HTML {
	var result strings.Builder

	if !templateFuncIsSet(keyName, data) {
		return ""
	}

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	field := v.FieldByName(keyName)

	if !field.IsValid() || field.Kind() != reflect.Slice {
		return ""
	}

	for i := 0; i < field.Len(); i++ {
		include, ok := field.Index(i).Interface().(JavascriptInclude)
		if !ok {
			slog.Error("tried to do a javascript include that is the wrong structure")
			return ""
		}

		result.WriteString(fmt.Sprintf(`<script type="%s" src="%s"></script>`, include.Type, include.Src))
	}

	return template.HTML(result.String())
}

func stylesheetIncludes(keyName string, data any) template.HTML {
	var result strings.Builder

	if !templateFuncIsSet(keyName, data) {
		return ""
	}

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	field := v.FieldByName(keyName)

	if !field.IsValid() || field.Kind() != reflect.Slice {
		return ""
	}

	for i := 0; i < field.Len(); i++ {
		include, ok := field.Index(i).Interface().(StylesheetInclude)
		if !ok {
			slog.Error("tried to do a stylesheet include that is the wrong structure")
			return ""
		}

		result.WriteString(fmt.Sprintf(`<link type="text/css" rel="stylesheet" media="%s" href="%s" />`, include.Media, include.Href))
	}

	return template.HTML(result.String())
}

func join(s any, sep string) string {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return ""
	}

	var result strings.Builder
	length := v.Len()

	for i := 0; i < length; i++ {
		result.WriteString(fmt.Sprintf("%v", v.Index(i).Interface()))

		if i < length-1 {
			result.WriteString(sep)
		}
	}

	return result.String()
}

func normalizeTemplateDir(templateDir string) string {
	result := ""

	if strings.HasPrefix(templateDir, "/") {
		result = templateDir[1:]
	} else {
		result = templateDir
	}

	if strings.HasSuffix(result, "/") {
		result = result[:len(result)-1]
	}

	return result
}

func normalizeTemplateExt(templateExt string) string {
	if strings.HasPrefix(templateExt, ".") {
		return templateExt
	}

	return "." + templateExt
}

func normalizeTemplateName(templateName string) string {
	if !strings.HasPrefix(templateName, "/") {
		return templateName
	}

	return templateName[1:]
}
