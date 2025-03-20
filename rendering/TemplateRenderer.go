package rendering

import (
	"io"
)

type TemplateRenderer interface {
	/*
	   Render renders a template file into a layout template file using the provided
	   data to an io.Writer.
	*/
	Render(templateName string, data any, w io.Writer)

	/*
	   RenderString renders a Go template string with a set of data to an io.Writer.
	*/
	RenderString(templateString string, data any, w io.Writer)
}
