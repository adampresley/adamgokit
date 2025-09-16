package rendering

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"maps"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

type GoTemplateRendererConfig struct {
	AdditionalFuncs   template.FuncMap
	PagesDir          string
	TemplateDir       string
	TemplateExtension string
	TemplateFS        fs.FS
}

type GoTemplateRenderer struct {
	funcs             template.FuncMap
	pagesDir          string
	templateDir       string
	templateExtension string
	templateFS        fs.FS
	componentsDir     string
	layoutsDir        string

	allTemplates *template.Template
}

func NewGoTemplateRenderer(config GoTemplateRendererConfig) (*GoTemplateRenderer, error) {
	renderer := &GoTemplateRenderer{
		funcs:             config.AdditionalFuncs,
		pagesDir:          config.PagesDir,
		templateDir:       config.TemplateDir,
		templateExtension: config.TemplateExtension,
		templateFS:        config.TemplateFS,
	}

	if renderer.pagesDir == "" {
		return nil, fmt.Errorf("pages directory not provided. this is required to render pages. this directory should contain your HTML templates that rely on layouts and other components")
	}

	renderer.pagesDir = filepath.Clean(renderer.pagesDir)

	if renderer.templateExtension == "" {
		renderer.templateExtension = ".html"
	}

	if err := renderer.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return renderer, nil
}

func (tr *GoTemplateRenderer) loadTemplates() error {
	funcs := getFuncs(tr.funcs)
	tr.allTemplates = template.New("").Funcs(funcs)

	return fs.WalkDir(tr.templateFS, tr.templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, tr.templateExtension) {
			return nil
		}

		templateName := tr.getTemplateName(path)

		/*
		 * Skip pages in the pages directory, as we'll load them individually when needed.
		 */
		if strings.HasPrefix(templateName, tr.pagesDir) {
			return nil
		}

		content, err := fs.ReadFile(tr.templateFS, path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		tmpl := tr.allTemplates.New(templateName)

		if _, err := tmpl.Parse(string(content)); err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		slog.Debug("parsed template", "name", templateName)
		return nil
	})
}

func (tr *GoTemplateRenderer) getTemplateName(path string) string {
	relPath, _ := filepath.Rel(tr.templateDir, path)
	ext := filepath.Ext(relPath)
	return strings.TrimSuffix(relPath, ext)
}

func (tr *GoTemplateRenderer) Render(templateName string, data any, w io.Writer) error {
	slog.Debug("executing template", "name", templateName)

	/*
	 * For templates in the pages directory, we need to create a custom template that includes
	 * external templates. We do this by cloning allTemplates, then rendering the page template separately.
	 */
	if strings.HasPrefix(templateName, tr.pagesDir) {
		return tr.renderPageWithLayout(templateName, data, w)
	}

	clonedTemplates, err := tr.allTemplates.Clone()

	if err != nil {
		return fmt.Errorf("failed to clone shared templates: %w", err)
	}

	tmpl := clonedTemplates.Lookup(templateName)

	if tmpl == nil {
		return fmt.Errorf("template %s not found", templateName)
	}

	if err := clonedTemplates.ExecuteTemplate(w, templateName, data); err != nil {
		slog.Error("failed to execute template", "name", templateName, "error", err)
		return tr.renderError(templateName, err, w)
	}

	return nil
}

func (tr *GoTemplateRenderer) renderPageWithLayout(templateName string, data any, w io.Writer) error {
	pageTemplate, err := tr.allTemplates.Clone()

	if err != nil {
		return fmt.Errorf("failed to clone shared templates: %w", err)
	}

	pagePath := filepath.Join(tr.templateDir, templateName+tr.templateExtension)
	pageContent, err := fs.ReadFile(tr.templateFS, pagePath)

	if err != nil {
		return fmt.Errorf("failed to read page template %s: %w", pagePath, err)
	}

	if _, err := pageTemplate.New(templateName).Parse(string(pageContent)); err != nil {
		return tr.renderError(templateName, err, w)
	}

	if err := pageTemplate.ExecuteTemplate(w, templateName, data); err != nil {
		return tr.renderError(templateName, err, w)
	}

	return nil
}

func (tr *GoTemplateRenderer) RenderString(templateString string, data any, w io.Writer) error {
	funcs := getFuncs(tr.funcs)
	tmpl, err := template.New("inline").Funcs(funcs).Parse(templateString)

	if err != nil {
		return fmt.Errorf("failed to parse template string: %w", err)
	}

	for _, t := range tr.allTemplates.Templates() {
		if _, err := tmpl.AddParseTree(t.Name(), t.Tree); err != nil {
			return fmt.Errorf("failed to add template %s: %w", t.Name(), err)
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("failed to execute template string", "error", err)
		return tr.renderError("inline template", err, w)
	}

	return nil
}

func (tr *GoTemplateRenderer) renderError(templateName string, err error, w io.Writer) error {
	templateString := `<html><body>{{.Content}}</body></html>`

	data := map[string]any{
		"Content": template.HTML(fmt.Sprintf(`
<h2>Rendering Error</h2>

<article style="background-color: #AF291D; color: white; padding: 1rem; border-radius: 5px;">
	<p>
		An error occurred while rendering the page '%s': %s
	</p>
</article>
`, templateName, err.Error())),
	}

	tmpl, err := template.New("error-inline").Parse(templateString)

	if err != nil {
		slog.Error("failed to parse nice error template", "error", err)
		return fmt.Errorf("failed to parse nice error template: %w", err)
	}

	if err = tmpl.Execute(w, data); err != nil {
		slog.Error("failed to execute nice error template", "error", err)
		return fmt.Errorf("failed to execute nice error template: %w", err)
	}

	return nil
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

	maps.Copy(templateFuncs, additionalFuncs)
	return templateFuncs
}

func templateFuncIsSet(name string, data any) bool {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Pointer {
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
	return slices.Contains(slice, item)
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

	if v.Kind() == reflect.Pointer {
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

	if v.Kind() == reflect.Pointer {
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

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return ""
	}

	var result strings.Builder
	length := v.Len()

	for i := range length {
		result.WriteString(fmt.Sprintf("%v", v.Index(i).Interface()))

		if i < length-1 {
			result.WriteString(sep)
		}
	}

	return result.String()
}
