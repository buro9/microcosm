package templates

import (
	"fmt"
	"html/template"
	"sync"
	"embed"

	"github.com/buro9/microcosm/web/templates/funcs"
)

var compileTemplatesOnce sync.Once

// templates is the map of compiled templates is populated at init via load.go
// from the definitions held in definitions.go, the definitions themselves use
// the Templates slice
var templates map[string]*template.Template

//go:embed templates/*/*.html.tmpl
var templateFS embed.FS

// Compile compiles templates and is expected to be called by main.go
// as we require that the flags are parsed first to obtain the value of
// *opts.FilesPath
func Compile() {
	compileTemplatesOnce.Do(
		func() {
			if templates == nil {
				templates = make(map[string]*template.Template)
			}

			pathFormat := "templates/%s/%s.html.tmpl"

			for _, t := range Templates {
				// Gather a list of all files required by this template
				var paths []string
				paths = append(
					paths,
					fmt.Sprintf(pathFormat, "base", t.Base),
				)
				paths = append(
					paths,
					fmt.Sprintf(pathFormat, "pages", t.Page),
				)
				for _, include := range t.Includes {
					paths = append(
						paths,
						fmt.Sprintf(pathFormat, "includes", include),
					)
				}

				// MustCompile all templates as a compile error is more preferable than a
				// runtime error
				templates[t.Page] =
					template.Must(
						template.New(t.Base).Funcs(funcs.FuncMap).ParseFS(
							templateFS,
							paths...,
						),
					)
			}
		},
	)
}
