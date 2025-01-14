package rendering

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"reflect"
	"slices"
	"strings"
)

type GoTemplateRendererConfig struct {
	AdditionalFuncs   template.FuncMap
	TemplateDir       string
	TemplateExtension string
	TemplateFS        fs.FS
}

type GoTemplateRenderer struct {
	funcs             template.FuncMap
	templateDir       string
	templateExtension string
	templateFS        fs.FS
	templates         *template.Template
}

func NewGoTemplateRenderer(config GoTemplateRendererConfig) *GoTemplateRenderer {
	ext := config.TemplateExtension

	if ext == "" {
		ext = ".html"
	}

	funcs := getFuncs(config.AdditionalFuncs)
	normalizedTemplatePath := fmt.Sprintf("%s/*%s", normalizeTemplateDir(config.TemplateDir), normalizeTemplateExt(config.TemplateExtension))

	result := &GoTemplateRenderer{
		funcs:             funcs,
		templateFS:        config.TemplateFS,
		templateExtension: normalizeTemplateExt(ext),
		templateDir:       normalizeTemplateDir(config.TemplateDir),
		templates:         template.Must(template.New("").Funcs(funcs).ParseFS(config.TemplateFS, normalizedTemplatePath)),
	}

	return result
}

/*
Render renders a Go template file into a layout template file using the provided
data to an io.Writer.
*/
func (tr *GoTemplateRenderer) Render(templateName string, data any, w io.Writer) {
	normalizedTemplateName := fmt.Sprintf("%s%s", normalizeTemplateName(templateName), tr.templateExtension)
	templateNameAndDir := fmt.Sprintf("%s/%s", tr.templateDir, normalizedTemplateName)

	tmpl := template.Must(tr.templates.Clone())
	tmpl = template.Must(tmpl.ParseFS(tr.templateFS, templateNameAndDir))

	if err := tmpl.ExecuteTemplate(w, normalizedTemplateName, data); err != nil {
		slog.Error("error executing template", "error", err, "templateName", normalizedTemplateName)
		fmt.Fprintf(w, "error executing template '%s': %s", templateName, err.Error())
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
