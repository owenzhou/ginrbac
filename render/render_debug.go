package render

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin/render"
)

type RenderDebug map[string]*template.Template

var (
	_ render.HTMLRender = Render{}
)

func NewDebug() RenderDebug {
	return make(RenderDebug)
}

func (r RenderDebug) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	if _, ok := r[name]; ok {
		//panic(fmt.Sprintf("template %s already exists", name))
		return
	}
	r[name] = tmpl
}

func (r RenderDebug) AddFromFS(fs fs.FS, name string, files ...string) *template.Template {
	tmpl := template.Must(template.ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

func (r RenderDebug) AddFromFSFuncs(fs fs.FS, name string, funcMap template.FuncMap, files ...string) *template.Template {
	tname := filepath.Base(files[0])
	tmpl := template.Must(template.New(tname).Funcs(funcMap).ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

func (r RenderDebug) Instance(name string, data interface{}) render.Render {
	if strings.Index(name, "/") > 0 {
		name = "/" + name
	}
	if _, ok := r[name]; !ok {
		panic("Template Error: view " + name + " not found.")
	}
	return render.HTML{
		Template: r[name],
		Data:     data,
	}
}
