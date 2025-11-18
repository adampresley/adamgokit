package rendering

import "html/template"

type BaseViewModel struct {
	IsError            bool
	IsHtmx             bool
	IsWarning          bool
	JavascriptIncludes []JavascriptInclude
	Message            template.HTML
}
