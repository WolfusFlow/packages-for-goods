package html

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed templates/*.html
//go:embed static/*.css
var embeddedFiles embed.FS

func ParseTemplates() (*template.Template, error) {
	return template.New("").ParseFS(embeddedFiles, "templates/*.html")
}

func StaticFileServer() http.Handler {
	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic("failed to create static sub FS: " + err.Error())
	}
	return http.FileServer(http.FS(staticFS))
}
