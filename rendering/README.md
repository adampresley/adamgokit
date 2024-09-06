# Rendering

This package provides components that can be used to render content.

## Go Template Renderer

The component `GoTemplateRenderer` (which implements the `TemplateRenderer`
interface) is used to render Go templates to an `io.Writer`. This is
particularly useful for rendering web applications. Here is a basic example
that uses this component for rendering an HTTP route. This example makes
a few assumtions:

- You have a directory named "templates" in the root of your application
- You have the files _layout.tmpl_ and _index.tmpl_

```go
import (
  "embed"
  "net/http"

  "github.com/adampresley/adamgokit/rendering"
  "github.com/adampresley/adamgokit/mux"
)

var (
  //go:embed templates/*
  templateFS embed.FS

  renderer rendering.TemplateRenderer
)

func main() {
  renderer = rendering.NewGoTemplateRenderer(rendering.GoTemplateRendererConfig{
    TemplateDir:       "templates",
    TemplateExtension: ".tmpl",
    TemplateFS:        templateFS,
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

  renderer.Render("index", "layout", viewData, w)
}
```

### Configuration

The basic requirements of configuring the Go template renderer are:

- Template directory - The directory where Go template files live
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
  AdditionalFunc: moreFuncs,
  TemplateDir:       "templates",
  TemplateExtension: ".tmpl",
  TemplateFS:        templateFS,
})
```
