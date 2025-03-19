package rendering

import (
	"html/template"
	"io"
)

type TemplateRenderer interface {
	/*
	   Render renders a template file into a layout template file using the provided
	   data to an io.Writer.
	*/
	Render(templateName string, data any, w io.Writer)

	/*
		RenderWithFuncs renders a template file using the provided data
		to an io.Writer with ad-hoc template functions.
	*/
	RenderWithFuncs(templateName string, data any, funcs template.FuncMap, w io.Writer)
	/*
	   RenderString renders a Go template string with a set of data to an io.Writer.
	*/
	RenderString(templateString string, data any, w io.Writer)
}
