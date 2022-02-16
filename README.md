# 为什么要使用ginrbac？
正常开发权限管理系统 gin代码：
```go
router := gin.New()
//获取用户列表
router.GET("/user/list", ...gin.HandlerFunc)
router.Run()
```
//然后还要手动插入数据库，方便使用人员授权时勾选，下面只写个sql语句意思下，大家都懂的...
```go
insert into `permissions` ('name', 'url', 'method') values ('获取用户列表', '/user/list', 'GET')
```

是不是感觉我们开发人员好辛苦，写路由时写了一遍注释，还得插入数据库，再写一遍注释...

使用ginrbac开发权限管理，注释直接写在路由函数最后面：
```go
a := app.NewApp()
a.Router.Get("/user/list", ...app.HandlerFunc, "获取用户列表")
a.Router.Get("/permissions", func(c *app.Context){
  //获取所有路由及注释
  routes := a.GetRoutes()
  c.JSON(200, app.H{
    "routes": routes,
  })
}, "获取权限列表")
a.Run()
```
#用户登录，简单示例
```go
a.Router.Post("/sign", func(c *app.Context){
  credentials := map[string]interfade{}{
    "email": c.PostForm("email"),
    "password": c.PostForm("password"),
  }
  //在gin.Context包装了一层加入了 Auth方法
  signed := c.Auth().Attempt(credentials, true)
  if signed {
    c.JSON(200, app.H{
      "msg": "登录成功",
    })
    return
  }
  c.JSON(200, app.H{
    "msg": "登录失败",
  })
}, "ajax用户登录")
```
