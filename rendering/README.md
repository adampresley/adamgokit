# Rendering

This package provides components that can be used to render content.

## Go Template Renderer

The component `GoTemplateRenderer` (which implements the `TemplateRenderer`
interface) is used to render Go templates to an `io.Writer`. This is
particularly useful for rendering web applications. Here is a basic example
that uses this component for rendering an HTTP route. This example makes
a few assumtions:

- You have a directory named "app" in the root of your application
- You have a subdirectory under "app" named "static"
- You have the files _layouts/layout.html_ and _pages/index.html_ under "app"

Here are the layout and index pages respectively.

```html
{{- define "layouts/layout"}}
<!DOCTYPE html>
<html lang="en">

<head>
   <meta charset="UTF-8" />
   <meta name="viewport" content="width=device-width, initial-scale=1.0" />
   <meta name="color-scheme" content="light dark" />
   <title>{{template "title" .}}</title>
   <link type="text/css" rel="stylesheet" media="screen" href="/static/css/pico.min.css" />
   {{stylesheetIncludes "Stylesheets" .}}
</head>

<body>
   <header class="grid top-header">
      <nav>
         <ul>
            <li>
               <h1><a href="/">My App</a></h1>
            </li>
         </ul>

         <ul>
            <li>
               <a href="/logout">Log Out</a>
            </li>
         </ul>
      </nav>
   </header>
   <!-- END NAV -->

   <main class="container">
      {{template "content" .}}
   </main>

   <footer id="page-footer">
      <p>&copy; Mine!</p>
   </footer>

   {{javascriptIncludes "JavascriptIncludes" .}}
</body>

</html>
{{end}}
```

```html
{{template "layouts/layout" .}}
{{define "title"}}Home{{end}}
{{define "content"}}

<h2>Home Page</h2>
<p>
	This is a sample home page.
</p>

{{end}}
```

```go
import (
  "embed"
  "net/http"

  "github.com/adampresley/adamgokit/rendering"
  "github.com/adampresley/adamgokit/mux"
)

var (
  //go:embed app
  appFS embed.FS

  renderer rendering.TemplateRenderer
)

func main() {
  var (
    err error
  )

  shutdownCtx, stopApp := context.WithCancel(context.Background())

  // INITIALIZE the renderer using default values
  if renderer, err = rendering.NewGoTemplateRenderer(appFS); err != nil {
    panic(err)
  }

  routes := []mux.Route{
    {Path: "GET /", Handler: http.HandlerFunc(getIndex)},
  }

  mux := mux2.Setup(
    &config, // This would be a struct that embeds mux.MuxConfig
    routes,
    shutdownCtx,
    stopApp,

    mux2.WithStaticContent("app", "/static/", appFS),
  )

  mux.Start()
}

func getIndex(w http.ResponseWriter, r *http.Request) {
  viewData := struct{
    Title string
  }{
    Title: "Hello!",
  }

  renderer.Render("pages/index", viewData, w)
}
```

### Configuration

The basic requirements of configuring the Go template renderer are:

- Template directory - The directory where Go template files live. 
- Template extension - The extension for Go template files. The default for this is **.html**
- Template filesystem - A filesystem reference for Go template files
- Pages directory - The name of the directory that houses pages. Pages can reference other templates like layouts, components, and partials. Those should be in a different directory. The default for this is **pages**

A good layout might look like this:

```
cmd/your-app/                  # Main application
├── app/                       # Frontend assets
│   ├── layouts/               # HTML layout templates
│   ├── pages/                 # Page templates
│   ├── components/            # Reusable components
│   └── static/                # CSS, JS, images
```


The Go template renderer includes a small set of handy functions that
you can use in your templates. These include:

- `join` - Allows you to join two strings
- `isSet` - Checks for a variable's existance. `{{if (isSet "Stylesheets" .}}`
- `isLastItem` - Returns true if the index is equal to the index of the last
  item. `{{if (isLastItem $index (len .MyArray))}}`
- `stringNotEmpty` - Returns true if the provided string isn't empty. This
  function also handles HTML templates, and automatically trims spaces.
  `{{if (stringNotEmpty .SomeString)}}`

### Additional Options

You can customize some behaviors of this component during initalization by passing option functions. Here is an example of adding additional functions to the renderer.

```go
moreFuncs := template.FuncMap{
  "newFunc": // func goes here...
}

renderer, err = rendering.NewGoTemplateRenderer(
  appFS,
  rendering.WithFuncs(moreFuncs),
)
```

#### Template Functions

Call `WithFuncs(funcs)`. **funcs** is a _template.FuncMap_.

#### Pages directory

Call `PagesDir(dir)`. **dir** is a string with the relative subdirectory of your root filesystem directory.

#### Template directory

Call `TemplateDir(dir)`. **dir** is a string with the relative subdirectory of your root filesystem directory.

#### Template extension

Call `TemplateExtension(ext)`. **ext** is a string with the extension to use for template files (including the period).
