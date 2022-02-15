package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type route struct {
	Url      string   `json:"url"`
	Comment  string   `json:"comment"`
	Method   string   `json:"method"`
	Group    string   `json:"group"`
	Children []*route `json:"children"`
}

//添加String方法，方便打印出数据
func (r *route) String() string {
	return fmt.Sprintf("{Url:%s, Comment:%s, Method:%s, Group:%s, Children: %v}", r.Url, r.Comment, r.Method, r.Group, r.Children)
}

//初始化设置根路由
var Routes = []*route{
	{Url: "/", Comment: "根路由", Method: "GROUP", Group: "/"},
}

//往 Routes 里添加路由
func addRoute(routes []*route, route *route) {

	groupExists := false
	//查找中由组是否存在
	for _, v := range Routes {
		if v.Group == route.Group+route.Url {
			groupExists = true
			break
		}

		if route.Method != "GROUP" && route.Group == v.Group {
			if route.Group != "/" {
				route.Url = route.Group + route.Url
			}
			v.Children = append(v.Children, route)
			break
		}
	}
	if !groupExists && route.Method == "GROUP" {
		if route.Group == "/" {
			route.Group = route.Url
		} else {
			route.Group = route.Group + route.Url
			route.Url = route.Group
		}
		Routes = append(Routes, route)
	}
}

//-------------重写 gin 路由方法---------------

type IRoutes interface {
	Use(...interface{}) IRoutes
	Handle(string, string, ...interface{}) IRoutes
	Any(string, ...interface{}) IRoutes
	Get(string, ...interface{}) IRoutes
	Post(string, ...interface{}) IRoutes
	Delete(string, ...interface{}) IRoutes
	Patch(string, ...interface{}) IRoutes
	Put(string, ...interface{}) IRoutes
	Options(string, ...interface{}) IRoutes
	Head(string, ...interface{}) IRoutes
}

type IRouter interface {
	IRoutes
	Group(string, ...interface{}) *RouterGroup
}

type RouterGroup struct {
	App     *App
	GRouter *gin.RouterGroup
	Comment string
	root    bool
}

func (r *RouterGroup) returnObj() IRoutes {
	if r.root {
		return &r.App.Router
	}
	return r
}

func (r *RouterGroup) Resource(relativePath string, ctrl IController, comment string) {
	r.Get("/"+relativePath, ctrl.Index, comment+"-"+"列表")
	r.Get("/"+relativePath+"/:id", ctrl.Show, comment+"-"+"获取")
	r.Get("/"+relativePath+"/create", ctrl.Create, comment+"-"+"创建")
	r.Post("/"+relativePath, ctrl.Store, comment+"-"+"保存")
	r.Get("/"+relativePath+"/:id/edit", ctrl.Edit, comment+"-"+"编辑")
	r.Patch("/"+relativePath+"/:id", ctrl.Update, comment+"-"+"更新")
	r.Delete("/"+relativePath+"/:id", ctrl.Destroy, comment+"-"+"删除")
}

func (r *RouterGroup) Group(relativePath string, args ...interface{}) *RouterGroup {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "GROUP", Group: r.GRouter.BasePath()})
	group := &RouterGroup{
		App:     r.App,
		GRouter: r.GRouter.Group(relativePath, handlers...),
		Comment: comment,
	}

	return group
}

func (r *RouterGroup) Use(args ...interface{}) IRoutes {
	handlers, _ := r.App.getComment(args...)
	r.GRouter.Use(handlers...)
	return r.returnObj()
}

func (r *RouterGroup) Handle(httpMethod, relativePath string, args ...interface{}) IRoutes {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "HANDLE", Group: r.GRouter.BasePath()})
	r.GRouter.Handle(httpMethod, relativePath, handlers...)
	return r.returnObj()
}

func (r *RouterGroup) Any(relativePath string, args ...interface{}) IRoutes {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "ANY", Group: r.GRouter.BasePath()})
	r.GRouter.Any(relativePath, handlers...)
	return r.returnObj()
}

func (r *RouterGroup) Get(relativePath string, args ...interface{}) IRoutes {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "GET", Group: r.GRouter.BasePath()})
	r.GRouter.GET(relativePath, handlers...)
	return r.returnObj()
}

func (r *RouterGroup) Post(relativePath string, args ...interface{}) IRoutes {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "POST", Group: r.GRouter.BasePath()})
	r.GRouter.POST(relativePath, handlers...)
	return r.returnObj()
}

func (r *RouterGroup) Delete(relativePath string, args ...interface{}) IRoutes {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "DELETE", Group: r.GRouter.BasePath()})
	r.GRouter.DELETE(relativePath, handlers...)
	return r.returnObj()
}

func (r *RouterGroup) Patch(relativePath string, args ...interface{}) IRoutes {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "PATCH", Group: r.GRouter.BasePath()})
	r.GRouter.PATCH(relativePath, handlers...)
	return r.returnObj()
}

func (r *RouterGroup) Put(relativePath string, args ...interface{}) IRoutes {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "PUT", Group: r.GRouter.BasePath()})
	r.GRouter.PUT(relativePath, handlers...)
	return r.returnObj()
}

func (r *RouterGroup) Options(relativePath string, args ...interface{}) IRoutes {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "OPTIONS", Group: r.GRouter.BasePath()})
	r.GRouter.OPTIONS(relativePath, handlers...)
	return r.returnObj()
}

func (r *RouterGroup) Head(relativePath string, args ...interface{}) IRoutes {
	handlers, comment := r.App.getComment(args...)
	addRoute(Routes, &route{Url: relativePath, Comment: comment, Method: "HEAD", Group: r.GRouter.BasePath()})
	r.GRouter.HEAD(relativePath, handlers...)
	return r.returnObj()
}

//检测RouterGroup是否实现了IRouter接口
var _ IRouter = &RouterGroup{}
