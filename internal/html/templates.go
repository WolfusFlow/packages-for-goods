package html

import (
	"embed"
	"html/template"
)

//go:embed templates/*.html
var TemplateFS embed.FS

func ParseTemplates() (*template.Template, error) {
	return template.ParseFS(TemplateFS, "templates/*.html")
}
