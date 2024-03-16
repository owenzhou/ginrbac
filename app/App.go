package app

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	bootstrapconfig "github.com/owenzhou/ginrbac/bootstrap"
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/render"
	"github.com/owenzhou/ginrbac/support/facades"
	"github.com/owenzhou/ginrbac/utils"

	"github.com/gin-gonic/gin"
)

type HandlerFunc = func(*Context)

type Action map[IController]string

func NewApp(views ...fs.FS) *App {

	app := &App{
		bindings:  make(map[string]interface{}),
		instances: make(map[string]interface{}),
		shareMap:  make(map[string]bool),
	}

	//注册基础的服务及门面
	app.Register(new(bootstrapconfig.Providers))
	app.Register(new(bootstrapconfig.Facades))

	if facades.Config != nil {
		//设置运行模式
		if facades.Config.Debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		location, err := time.LoadLocation(facades.Config.Timezone)
		if err != nil {
			panic(err)
		}
		time.Local = location
	}

	//不使用logger，只使用recovery
	engine := gin.New()
	engine.Use(gin.Recovery())

	app.Engine = engine

	if facades.Config != nil {
		//设置默认全局中间件,必需在 app.IRouter 各方法前面
		app.Middlewares(new(bootstrapconfig.Middlewares))
	}

	//设置模板及静态文件
	if len(views) > 0 {
		//全局变量Views赋值
		facades.Views = views[0]
		app.SetFuncMap(template.FuncMap{})
		app.SetHTMLTemplate(views[0])
	}

	app.Router = RouterGroup{
		App:     app,
		GRouter: engine.Group("/"),
		root:    true,
	}

	//设置自定义中间件，必需要 app.Router 之后
	//app.Middlewares(new(customconfig.Middlewares))

	app.pool.New = func() interface{} {
		return app.allocateContext()
	}

	return app
}

type App struct {
	Engine    *gin.Engine
	Router    RouterGroup
	FuncMap   template.FuncMap
	bindings  map[string]interface{}
	instances map[string]interface{}
	shareMap  map[string]bool
	pool      sync.Pool
}

func (a *App) allocateContext() *Context {
	return &Context{}
}

// 包装进自己的context
func (a *App) bindContext(handler interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := a.pool.Get().(*Context)
		ctx.Context = c
		ctx.share = nil

		if handle, ok := handler.(HandlerFunc); ok {
			handle(ctx)
		} else if handle, ok := handler.(Action); ok {
			for iCtrl, actionStr := range handle {
				//调用init方法
				iCtrl.Init(ctx)

				controller := reflect.ValueOf(iCtrl)
				action := controller.MethodByName(actionStr)
				if action.Kind() != reflect.Invalid {
					action.Call([]reflect.Value{reflect.ValueOf(ctx)})
				} else {
					panic(controller.Type().String() + " has no method " + actionStr)
				}
			}
		} else {
			panic("bindcontext error: except gin.HandlerFunc, app.HandlerFunc or app.Action")
		}

		a.pool.Put(ctx)
	}
}

// 私有方法，判断路由是否有注解，并绑定自定义context
func (a *App) getComment(args ...interface{}) ([]gin.HandlerFunc, string) {
	//取接口最后一个元素，判断是否为字符串注解
	var comment string
	var handlers = make([]gin.HandlerFunc, 0)

	if len(args) == 0 {
		return handlers, comment
	}

	op := len(args) - 1
	last := args[op]
	if reflect.TypeOf(last).Kind().String() == "string" {
		comment = last.(string)
		args = args[:op]
	}

	for _, v := range args {
		if handler, ok := v.(HandlerFunc); ok {
			handlers = append(handlers, a.bindContext(handler))
			continue
		}

		if handler, ok := v.(Action); ok {
			handlers = append(handlers, a.bindContext(handler))
			continue
		}

		if handler, ok := v.(func(*gin.Context)); ok {
			handlers = append(handlers, handler)
		}
	}

	return handlers, comment
}

func (app *App) GetEngine() *gin.Engine {
	return app.Engine
}

// 使用中间件
func (app *App) Use(middlewares ...interface{}) {
	handlers, _ := app.getComment(middlewares...)
	app.Engine.Use(handlers...)
}

// 页面404
func (app *App) NoRoute(handlers ...interface{}) {
	handlers2, _ := app.getComment(handlers...)
	app.Engine.NoRoute(handlers2...)
}

// 获取应用名称
func (app *App) Name() string {
	return facades.Config.AppName
}

// 调用默认的全局中间件
func (app *App) Middlewares(middlewares interface{}) {
	v := reflect.Indirect(reflect.ValueOf(app.InjectApp(middlewares)))
	for i := 0; i < v.NumField(); i++ {
		middleware := v.Field(i)
		if middleware.MethodByName("Middleware").Kind() != reflect.Invalid {
			middleware.Interface().(contracts.IMiddleware).Middleware()
		}
	}
}

// 往指定的结构体注入app对象
func (app *App) InjectApp(providers interface{}) interface{} {
	v := reflect.Indirect(reflect.ValueOf(providers))

	//如果是找到App interfaces.IApplication退出递归
	if v.Type().Kind() == reflect.Interface {
		return reflect.ValueOf(app).Interface()
	}

	for i := 0; i < v.NumField(); i++ {
		provider := v.Field(i)

		if provider.IsZero() {
			var instance reflect.Value
			if provider.Type().Kind() == reflect.Interface {
				instance = reflect.New(provider.Type())
			} else if provider.Type().Kind() == reflect.Ptr {
				instance = reflect.New(provider.Type().Elem())
			} else {
				instance = reflect.New(provider.Type()).Elem()
			}
			if !provider.CanSet() {
				continue
			}
			provider.Set(reflect.ValueOf(app.InjectApp(instance.Interface())))
		}
	}

	return v.Addr().Interface()
}

func (app *App) Register(providers interface{}) {
	v := reflect.Indirect(reflect.ValueOf(app.InjectApp(providers)))

	for i := 0; i < v.NumField(); i++ {
		provider := v.Field(i)
		if provider.MethodByName("Boot").Kind() != reflect.Invalid {
			provider.Interface().(contracts.IServiceProvider).Boot()
		}
		if provider.MethodByName("Register").Kind() != reflect.Invalid {
			provider.Interface().(contracts.IServiceProvider).Register()
		}

		if provider.MethodByName("GetFacadeAccessor").Kind() != reflect.Invalid {
			provider.Interface().(contracts.IFacade).GetFacadeAccessor()
		}
	}
}

// 只实例化一次
func (app *App) Singleton(abstract string, closure contracts.ServiceFunc) {
	app.Bind(abstract, closure, true)
}

// 每次实例化
func (app *App) Bind(abstract string, closure contracts.ServiceFunc, share ...bool) {

	if app.bindings != nil {
		delete(app.bindings, abstract)
	}

	app.bindings[abstract] = closure
	if len(share) > 0 {
		app.shareMap[abstract] = share[0]
	} else {
		app.shareMap[abstract] = false
	}
}

// 创建实例
func (app *App) Make(abstract string) interface{} {
	if _, ok := app.instances[abstract]; ok {
		return app.instances[abstract]
	}
	if _, ok := app.bindings[abstract]; !ok {
		return nil
	}
	instance := app.bindings[abstract].(contracts.ServiceFunc)(app)
	if share, ok := app.shareMap[abstract]; ok && share {
		app.instances[abstract] = instance
	}
	return instance
}

// 注册事件
func (app *App) RegistEvent(listens interface{}) {
	//反射创建listens
	v := reflect.ValueOf(listens).Elem()
	for i := 0; i < v.NumField(); i++ {
		event := v.Field(i)
		if event.Kind() != reflect.Ptr {
			continue
		}
		if event.IsZero() {
			instance := reflect.New(event.Type().Elem())
			event.Set(instance)
		}

		listeners := event.Elem()
		for j := 0; j < listeners.NumField(); j++ {
			listener := listeners.Field(j)
			if listener.Kind() != reflect.Ptr {
				continue
			}
			if listener.IsZero() {
				instance := reflect.New(listener.Type().Elem())
				listener.Set(instance)
			}

			facades.Event.Attach(event.Type().String(), listener)
		}
	}
}

// 将打包的模板文件整理并且可以嵌套使用
func (app *App) SetHTMLTemplate(views fs.FS) {
	render := app.loadTemplate(views)
	app.Engine.HTMLRender = render

	//设置静态资源
	if gin.IsDebugging() {
		app.Engine.Static("/assets", "./views/layouts/assets")
		return
	}
	fs, _ := fs.Sub(views, "views/layouts/assets")
	app.Engine.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/assets/") {
			c.Header("Cache-Control", "max-age=5184000")
		}
		c.Next()
	})
	app.Engine.StaticFS("/assets", http.FS(fs))
}

// 设置模板函数
func (app *App) SetFuncMap(funcMap template.FuncMap) {
	app.FuncMap = utils.FuncMap
	for key, value := range funcMap {
		if _, ok := utils.FuncMap[key]; ok {
			continue
		}
		app.FuncMap[key] = value
	}
	app.FuncMap["widget"] = utils.Widget
}

// 加载模板
func (app *App) loadTemplate(views fs.FS) render.Renderer {
	r := render.New()
	//获取模板配置
	conf := facades.Config.Template
	//循环注册各模块模板
	for _, module := range conf {
		moduleLayout, ok := module["layout"]
		if !ok || moduleLayout == "" {
			log.Printf("Template Error: %s %s\n", module["modulename"], "'layout' tag not find or empty.")
			continue
		}
		layout, err := fs.Glob(views, moduleLayout)
		if err != nil {
			log.Printf("Template Error: %s %s\n", module["modulename"], err)
			continue
		}
		if len(layout) > 0 {
			contentViewPath, ok := module["viewpath"]
			if !ok || contentViewPath == "" {
				log.Printf("Template Error: %s %s\n", module["modulename"], "'viewpath' tag not find or empty.")
				continue
			}
			var contents []string
			contentSplit := strings.Split(contentViewPath, ";")
			for _, split := range contentSplit {
				contentGlob, err := fs.Glob(views, split)
				if err != nil {
					log.Printf("Template Error: %s %s %s\n", module["modulename"], "'viewpath' tag", err)
					continue
				}
				contents = append(contents, contentGlob...)
			}

			for _, content := range contents {
				layoutCopy := make([]string, len(layout))
				copy(layoutCopy, layout)
				files := append(layoutCopy, content)
				name := content[strings.Index(content, "/"):strings.LastIndex(content, ".")]
				r.AddFromFSFuncs(views, name, app.FuncMap, files...)
			}
		}
	}

	return r
}

// 获取路由树
func (app *App) GetRoutesTree() interface{} {
	return routes
}

// 获取路由数组，数组每个元素都是一个路由组
func (app *App) GetRoutesGroup() (result interface{}) {
	result = generateGroup(routes)
	return
}

// 递归生成路由组
func generateGroup(node *route) (result []*route) {
	//退出递归
	if len(node.Children) == 0 {
		return
	}
	//深拷贝node数据
	routeCopy := *node
	routeCopy.Children = []*route{}
	result = append(result, &routeCopy)
	for _, v := range node.Children {
		if v.Method != "GROUP" {
			routeCopy.Children = append(routeCopy.Children, v)
			continue
		}
		result = append(result, generateGroup(v)...)
	}

	return
}

// 开始程序
func (app *App) Run(addr ...string) {
	app.Engine.Run(addr...)
}
