package html

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
)

//go:embed templates/*.html
//go:embed static/*.css
var embeddedFiles embed.FS

func ParseTemplates() (*template.Template, error) {
	tmpl := template.New("")

	files, err := fs.Glob(embeddedFiles, "templates/*.html")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		content, err := embeddedFiles.ReadFile(file)
		if err != nil {
			return nil, err
		}
		_, err = tmpl.New(filepath.Base(file)).Parse(string(content))
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

func StaticFileServer() http.Handler {
	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic("failed to create static sub FS: " + err.Error())
	}
	return http.FileServer(http.FS(staticFS))
}
