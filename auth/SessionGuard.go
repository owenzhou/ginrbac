package auth

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/owenzhou/ginrbac/contracts"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionGuard struct {
	name     string
	provider contracts.IUserProvider
	session  sessions.Session
	context  *gin.Context
}

func (a *SessionGuard) Check() bool {
	return !a.User().IsEmpty()
}

func (a *SessionGuard) Guest() bool {
	return !a.Check()
}

func (a *SessionGuard) ID() string {
	return a.User().GetAuthIdentifier()
}

func (a *SessionGuard) Validate(credentials map[string]interface{}) bool {
	user := a.provider.RetrieveByCredentials(credentials)
	return a.hasValidCredentials(user, credentials)
}

func (a *SessionGuard) SetUser(user contracts.IUser) {
	//将用户信息写进gin.Context
	a.context.Set(a.getAuthUserName(), user)
}

func (a *SessionGuard) User() contracts.IUser {
	var user contracts.IUser
	//如果gin.context里有数据
	u, exists := a.context.Get(a.getAuthUserName())
	if exists {
		user = u.(contracts.IUser)
		return user
	}
	//如果服务器不存在session
	id := a.session.Get(a.getName())
	if id == nil {
		if a.getRemember() != "" { //如果客户端存在cookie
			user = a.provider.RetrieveByToken(a.getID(), a.getToken())
			if user != nil {
				a.updateSession(user.GetAuthIdentifier())
				//将用户信息写进gin.Context
				a.SetUser(user)
			} else {
				user = &GenericUser{}
			}
		} else { //客户端不存在cookie，则返回空
			user = &GenericUser{}
		}
		return user
	}

	user = a.provider.RetrieveById(id.(string))
	a.SetUser(user)
	return user
}

//session, cookie 会用到context
func (a *SessionGuard) WithContext(c *gin.Context) contracts.IAuth {
	a.session = sessions.Default(c)
	a.context = c
	return a
}

//尝试登录
func (a *SessionGuard) Attempt(credentials map[string]interface{}, remember ...bool) bool {
	user := a.provider.RetrieveByCredentials(credentials)
	if a.hasValidCredentials(user, credentials) {
		a.Login(user, remember...)

		return true
	}
	return false
}

//使用IUser登录
func (a *SessionGuard) Login(user contracts.IUser, remember ...bool) {

	a.updateSession(user.GetAuthIdentifier())

	if len(remember) > 0 && remember[0] {
		a.ensureRememberTokenIsSet(user)
		a.setRemember(user)
	}

	a.SetUser(user)

}

func (a *SessionGuard) Logout() {
	user := a.User()

	//清除gin.Context里的信息
	a.context.Set(a.getAuthUserName(), nil)
	//清除session
	a.session.Delete(a.getName())
	a.session.Save()
	//清除cookie
	if a.getRemember() != "" {
		name := a.getRecallerName()
		value := ""
		expires := -7200
		a.setCookie(name, value, expires)
	}

	//更新数据库remember_token
	if !user.IsEmpty() {
		a.cycleRememberToken(user)
	}
}

// --------------------------------------

func (a *SessionGuard) getAuthUserName() string {
	return a.name + "_auth_user"
}

func (a *SessionGuard) hasValidCredentials(user contracts.IUser, credentials map[string]interface{}) bool {
	if user.IsEmpty() {
		return false
	}
	return a.provider.ValidateCredentials(user, credentials)
}

func (a *SessionGuard) updateSession(id string) {
	a.session.Set(a.getName(), id)
	a.session.Save()
}

func (a *SessionGuard) ensureRememberTokenIsSet(user contracts.IUser) {
	if user.GetRememberToken() == "" {
		a.cycleRememberToken(user)
	}
}

func (a *SessionGuard) random(l int) string {
	baseStr := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		str += string(baseStr[r.Intn(len(baseStr))])
	}
	str = base64.StdEncoding.EncodeToString([]byte(str))
	replacer := strings.NewReplacer("/", "", "+", "", "=", "")
	str = replacer.Replace(str)
	return str
}

func (a *SessionGuard) cycleRememberToken(user contracts.IUser) {
	token := a.random(60)
	user.SetRememberToken(token)
	a.provider.UpdateRememberToken(user, token)
}

func (a *SessionGuard) setRemember(user contracts.IUser) {
	name := a.getRecallerName()
	value := user.GetAuthIdentifier() + "|" + user.GetRememberToken() + "|" + user.GetAuthPassword()
	expires := 7200
	a.setCookie(name, value, expires)
}

func (a *SessionGuard) setCookie(name, value string, expires int) {
	path := "/"
	domain := ""
	secure := false   //设置这个 Cookie 是否仅仅通过安全的 HTTPS 连接传给客户端
	httponly := false //设置成 true，Cookie 仅可通过 HTTP 协议访问。 这意思就是 Cookie 无法通过类似 JavaScript 这样的脚本语言访问。 要有效减少 XSS 攻击时的身份窃取行为，可建议用此设置

	a.context.SetCookie(name, value, expires, path, domain, secure, httponly)
}

func (a *SessionGuard) getRemember() string {
	c, _ := a.context.Cookie(a.getRecallerName())
	return c
}

func (a *SessionGuard) getID() string {
	return strings.SplitN(a.getRemember(), "|", 3)[0]
}

func (a *SessionGuard) getToken() string {
	return strings.SplitN(a.getRemember(), "|", 3)[1]
}

func (a *SessionGuard) getName() string {
	return "login_" + a.name + "_" + fmt.Sprintf("%x", sha1.New().Sum([]byte("session")))
}

func (a *SessionGuard) getRecallerName() string {
	return "remember_" + a.name + "_" + fmt.Sprintf("%x", sha1.New().Sum([]byte("session")))
}
