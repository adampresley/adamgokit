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

type TemplateRenderer interface {
	/*
	   Render renders a template file into a layout template file using the provided
	   data to an io.Writer.
	*/
	Render(templateName, layoutName string, data any, w io.Writer)
	/*
	   RenderString renders a Go template string with a set of data to an io.Writer.
	*/
	RenderString(templateString string, data any, w io.Writer)
}

type GoTemplateRendererConfig struct {
	AdditionalFuncs   template.FuncMap
	TemplateDir       string
	TemplateExtension string
	TemplateFS        fs.FS
}

type GoTemplateRenderer struct {
	additionalFuncs   template.FuncMap
	templateDir       string
	templateExtension string
	templateFS        fs.FS
}

func NewGoTemplateRenderer(config GoTemplateRendererConfig) GoTemplateRenderer {
	return GoTemplateRenderer{
		additionalFuncs:   config.AdditionalFuncs,
		templateFS:        config.TemplateFS,
		templateExtension: config.TemplateExtension,
		templateDir:       config.TemplateDir,
	}
}

func (tr GoTemplateRenderer) getFuncs() template.FuncMap {
	templateFuncs := template.FuncMap{
		"join":                strings.Join,
		"isSet":               templateFuncIsSet,
		"isLastItem":          tr.isLastItem,
		"containsString":      containsString,
		"stringSliceContains": sliceContains[string],
		"uintSliceContains":   sliceContains[uint],
		"stringNotEmpty":      stringNotEmpty,
	}

	if tr.additionalFuncs != nil {
		for k, v := range tr.additionalFuncs {
			templateFuncs[k] = v
		}
	}

	return templateFuncs
}

/*
Render renders a Go template file into a layout template file using the provided
data to an io.Writer.
*/
func (tr GoTemplateRenderer) Render(templateName, layoutName string, data any, w io.Writer) {
	var (
		err  error
		tmpl *template.Template
	)

	templateFuncs := tr.getFuncs()
	templates := []string{
		fmt.Sprintf(
			"%s/%s%s",
			normalizeTemplateDir(tr.templateDir),
			normalizeTemplateName(templateName),
			normalizeTemplateExt(tr.templateExtension),
		),
	}

	if layoutName != "" {
		templates = append(templates, fmt.Sprintf(
			"%s/%s%s",
			normalizeTemplateDir(tr.templateDir),
			normalizeTemplateName(layoutName),
			normalizeTemplateExt(tr.templateExtension),
		))
	}

	if tmpl, err = template.New(templateName+".tmpl").Funcs(templateFuncs).ParseFS(tr.templateFS, templates...); err != nil {
		slog.Error("error parsing template", "error", err, "templateName", templateName, "layoutName", layoutName)
		fmt.Fprintf(w, "error parsing template '%s' (layout '%s'): %s", templateName, layoutName, err.Error())
		return
	}

	if err = tmpl.Execute(w, data); err != nil {
		slog.Error("error executing template", "error", err, "templateName", templateName, "layoutName", layoutName)
		fmt.Fprintf(w, "error executing template '%s' (layout '%s'): %s", templateName, layoutName, err.Error())
	}
}

/*
RenderString renders a Go template string with a set of data to an io.Writer.
*/
func (tr GoTemplateRenderer) RenderString(templateString string, data any, w io.Writer) {
	var (
		err  error
		tmpl *template.Template
	)

	templateFuncs := tr.getFuncs()

	if tmpl, err = template.New("raw").Funcs(templateFuncs).Parse(templateString); err != nil {
		slog.Error("error parsing template", "error", err)
		fmt.Fprintf(w, "error parsing template: %s", err.Error())
		return
	}

	if err = tmpl.Execute(w, data); err != nil {
		slog.Error("error executing template", "error", err)
		fmt.Fprintf(w, "error executing template: %s", err.Error())
	}
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

func (tr GoTemplateRenderer) isLastItem(index, length int) bool {
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
