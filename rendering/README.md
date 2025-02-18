# Rendering

This package provides components that can be used to render content.

## Go Template Renderer

The component `GoTemplateRenderer` (which implements the `TemplateRenderer`
interface) is used to render Go templates to an `io.Writer`. This is
particularly useful for rendering web applications. Here is a basic example
that uses this component for rendering an HTTP route. This example makes
a few assumtions:

- You have a directory named "app" in the root of your application
- You have the files _layouts/layout.html_ and _pages/index.html_

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
  renderer = rendering.NewGoTemplateRenderer(rendering.GoTemplateRendererConfig{
    TemplateDir:       "app",
    LayoutsDir:        "layouts",
    TemplateExtension: ".html",
    TemplateFS:        appFS,
  })

  routes := []mux.Route{
    {Path: "GET /", Handler: http.HandlerFunc(getIndex)},
  }

  routerConfig := mux.RouterConfig{
    Address: "localhost:8080",
  }

  m := mux.SetupRouter(routerConfig, routes)
  httpServer, quit := mux.SetupServer(routerConfig, m)

  <-quit
  mux.Shutdown(httpServer)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
  viewData := struct{
    Title string
  }{
    Title: "Hello!",
  }

  renderer.Render("index", viewData, w)
}
```

### Configuration

The basic requirements of configuring the Go template renderer are:

- Template directory - The directory where Go template files live
- Layout directory - The subdirectory (under template directory) where layouts live
- Template extension - The extension for Go template files
- Template filesystem - A filesystem reference for Go template files

The Go template renderer includes a small set of handy functions that
you can use in your templates. These include:

- `join` - Allows you to join two strings
- `isSet` - Checks for a variable's existance. `{{if (isSet "Stylesheets" .}}`
- `isLastItem` - Returns true if the index is equal to the index of the last
  item. `{{if (isLastItem $index (len .MyArray))}}`
- `stringNotEmpty` - Returns true if the provided string isn't empty. This
  function also handles HTML templates, and automatically trims spaces.
  `{{if (stringNotEmpty .SomeString)}}`

If you wish to include additional functions, you can add it to the
renderer configuration.

```go
moreFuncs := template.FuncMap{
  "newFunc": // func goes here...
}

renderer = rendering.NewGoTemplateRenderer(rendering.GoTemplateRendererConfig{
  AdditionalFunc:    moreFuncs,
  TemplateDir:       "app",
  LayoutsDir:        "layouts",
  TemplateExtension: ".html",
  TemplateFS:        appFS,
})
```
