package rendering

import "html/template"

type Options struct {
	funcs             template.FuncMap
	pagesDir          string
	templateDir       string
	templateExtension string
}

type Option func(o *Options)

func WithFuncs(funcs template.FuncMap) Option {
	return func(o *Options) {
		o.funcs = funcs
	}
}

func PagesDir(dir string) Option {
	return func(o *Options) {
		o.pagesDir = dir
	}
}

func TemplateDir(dir string) Option {
	return func(o *Options) {
		o.templateDir = dir
	}
}

func TemplateExtension(ext string) Option {
	return func(o *Options) {
		o.templateExtension = ext
	}
}
