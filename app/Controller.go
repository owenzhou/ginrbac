package app

import (
	"ginrbac/bootstrap/contracts"
	"net/http"
)

//restful api 控制器接口
type IController interface {
	Init(*Context)
	Index(*Context)
	Show(*Context)
	Create(*Context)
	Store(*Context)
	Edit(*Context)
	Update(*Context)
	Destroy(*Context)
}

type Controller struct {
	App contracts.IApplication
}

//初始化，在其它方法之前调用
func (ctrl *Controller) Init(c *Context) {}

func (ctrl *Controller) NotFound(c *Context) {
	c.JSON(http.StatusNotFound, H{
		"status":  404,
		"message": "method not found",
	})
}

//列表
func (ctrl *Controller) Index(c *Context) {
	ctrl.NotFound(c)
}

//显示指定数据
func (ctrl *Controller) Show(c *Context) {
	ctrl.NotFound(c)
}

//创建
func (ctrl *Controller) Create(c *Context) {
	ctrl.NotFound(c)
}

//创建保存
func (ctrl *Controller) Store(c *Context) {
	ctrl.NotFound(c)
}

//编辑
func (ctrl *Controller) Edit(c *Context) {
	ctrl.NotFound(c)
}

//编辑保存
func (ctrl *Controller) Update(c *Context) {
	ctrl.NotFound(c)
}

//删除指定数据
func (ctrl *Controller) Destroy(c *Context) {
	ctrl.NotFound(c)
}
