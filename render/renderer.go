package render

import (
	"html/template"
	"io/fs"

	"github.com/gin-gonic/gin/render"
)

type Renderer interface {
	render.HTMLRender
	Add(name string, tmpl *template.Template)
	AddFromFS(fs fs.FS, name string, files ...string) *template.Template
	AddFromFSFuncs(fs fs.FS, name string, funcMap template.FuncMap, files ...string) *template.Template
}
