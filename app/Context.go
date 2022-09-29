package app

import (
	"strings"

	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support/facades"
	"github.com/owenzhou/ginrbac/utils"

	"github.com/gin-gonic/gin"
)

//避免在控制器引入gin包，重定义gin.H
type H = gin.H

type Context struct {
	share map[string]interface{} //视图共享数据
	*gin.Context
}

//对含有*号的路由转为参数
//如 /user/*action ( /user/region/china/page/3 map[string]string{"region":"china", "page":"3"})
func (c *Context) ParseParams(key string) map[string]string {
	paramsMap := make(map[string]string)
	params := strings.Split(strings.Trim(c.Param(key), "/"), "/")
	if len(params)%2 == 0 {
		for i := 0; i < len(params); i += 2 {
			paramsMap[params[i]] = params[i+1]
		}
	}
	return paramsMap
}

//重写param函数，获取带.html, .htm, .xhtml等后缀的参数
func (c *Context) Param(key string) string {
	param := c.Context.Param(key)
	if i := strings.LastIndex(param, ".htm"); i > 0 {
		param = param[0:i]
	}
	return param
}


//共享视图数据
func (c *Context) Share(key string, data interface{}) {
	if c.share == nil {
		c.share = make(map[string]interface{})
	}
	c.share[key] = data
}

//获取共享数据
func (c *Context) Shared(key string) interface{} {
	if c.share == nil {
		return nil
	}
	if val, ok := c.share[key]; ok {
		return val
	}
	return ""
}

//重写HTML，加入share
func (c *Context) HTML(code int, name string, obj interface{}) {
	if data, ok := obj.(H); ok && c.share != nil {
		for k, v := range data {
			c.share[k] = v
		}
		c.Context.HTML(code, name, c.share)
		return
	}
	c.Context.HTML(code, name, obj)
}

//用户登录
func (c *Context) Auth(guard ...string) contracts.IAuth {
	auth := facades.App.Make("auth").(contracts.IAuthManager)
	return auth.Guard(guard...).WithContext(c.Context)
}

//分页
func (c *Context) NewPagination(total int64) *pagination {
	page := newPagination(total, c.Request.RequestURI)
	if p, ok := c.GetQuery(page.PageLabel); ok { //从query里查找分页标签
		page.CurrentPage = utils.Str2Int(p)
		page.Apply()
	} else if p := c.Param(page.PageLabel); p != "" { //从param里查找分页标签
		page.CurrentPage = utils.Str2Int(p)
		page.Apply()
	} else { //从当前url查找分页标签
		page.CurrentPage = page.FindCurrentPage()
		page.Apply()
	}

	return page
}

//判断是否是ajax请求
func (c *Context) IsAjax() (isAjax bool) {
	XRequestedWith := c.Request.Header.Get("X-Requested-With")
	if strings.HasPrefix(XRequestedWith, "XMLHttpRequest") {
		isAjax = true
	}
	return isAjax
}
