package contracts

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

//每个服务需实现该接口
type IServiceProvider interface {
	Register()
	Boot()
}

//绑定app的函数
type ServiceFunc func(IApplication) interface{}

//事件监听函数
type ListenFunc func(IEvent) bool

//app接口暴露给注入app的结构体
type IApplication interface {
	Bind(string, ServiceFunc, ...bool)
	Singleton(string, ServiceFunc)
	Make(string) interface{}
	Register(interface{})
	InjectApp(interface{}) interface{}
	Name() string
	RegistEvent(interface{})
	GetRoutes() []map[string]interface{}
	GetEngine() *gin.Engine
	Use(...interface{})
	NoRoute(...interface{})
}

//中间件接口
type IMiddleware interface {
	Middleware()
}

//facade接口
type IFacade interface {
	GetFacadeAccessor()
}

//事件接口
type IEvent interface {
	Data() map[string]interface{}
}

//观察者接口
type IObserver interface {
	Handle(IEvent) bool
}

//被观察者接口
type ISubject interface {
	Attach(string, reflect.Value)
	Detach(string, string)
	Fire(IEvent)
}

//用户接口
type IUser interface {
	Get(string) interface{}
	GetAuthIdentifierName() string
	GetAuthIdentifier() string
	GetAuthPassword() string
	GetRememberToken() string
	SetRememberToken(interface{})
	GetRememberTokenName() string
	SetApiToken(interface{})
	GetApiToken() string
	GetApiTokenName() string
	IsEmpty() bool
}

//用户提供者
type IUserProvider interface {
	RetrieveById(string) IUser
	RetrieveByToken(string, string) IUser
	UpdateRememberToken(IUser, string)
	UpdateApiToken(IUser, string)
	RetrieveByCredentials(map[string]interface{}) IUser
	ValidateCredentials(IUser, map[string]interface{}) bool
}

//登录管理接口
type IAuthManager interface {
	Guard(...string) IAuth
}

//用户登录接口
type IAuth interface {
	Check() bool
	Guest() bool
	User() IUser
	ID() string
	Attempt(map[string]interface{}, ...bool) bool
	Validate(map[string]interface{}) bool
	Login(IUser, ...bool)
	Logout()
	SetUser(IUser)
	WithContext(*gin.Context) IAuth
}

//加密接口
type IHash interface {
	Check(string, string, ...string) bool
	Make(string, ...map[string]int) string
}
