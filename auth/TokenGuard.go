package auth

import (
	"encoding/base64"
	"ginrbac/bootstrap/contracts"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenGuard struct {
	name     string
	provider contracts.IUserProvider
	context  *gin.Context
}

func (t *TokenGuard) Check() bool {
	return !t.User().IsEmpty()
}

func (t *TokenGuard) Guest() bool {
	return !t.Check()
}

func (t *TokenGuard) ID() string {
	return t.User().GetAuthIdentifier()
}

func (t *TokenGuard) Validate(credentials map[string]interface{}) bool {
	user := t.provider.RetrieveByCredentials(credentials)
	return t.hasValidCredentials(user, credentials)
}

func (t *TokenGuard) User() contracts.IUser {
	var user contracts.IUser
	u, exists := t.context.Get(t.getAuthUserName())
	//如果gin.Context里存在用户信息
	if exists {
		user = u.(contracts.IUser)
		return user
	}
	//如果request里不存在api_token
	token := t.getTokenFromRequest()
	if token == "" {
		user = &GenericUser{}
		return user
	}
	credentials := map[string]interface{}{
		t.getTokenName(): token,
	}
	user = t.provider.RetrieveByCredentials(credentials)
	t.SetUser(user)
	return user
}

func (t *TokenGuard) SetUser(user contracts.IUser) {
	//将用户信息写进gin.Context
	t.context.Set(t.getAuthUserName(), user)
}

func (t *TokenGuard) WithContext(c *gin.Context) contracts.IAuth {
	t.context = c
	return t
}

func (t *TokenGuard) Attempt(credentials map[string]interface{}, remember ...bool) bool {
	user := t.provider.RetrieveByCredentials(credentials)
	if t.hasValidCredentials(user, credentials) {
		t.Login(user, remember...)

		return true
	}
	return false
}

func (t *TokenGuard) Login(user contracts.IUser, remember ...bool) {
	if !user.IsEmpty() {
		t.updateToken(user)
	}
	t.SetUser(user)
}

func (t *TokenGuard) Logout() {
	user := t.User()
	if !user.IsEmpty() {
		t.cycleApiToken(user)
	}
}

//------------------------------

func (t *TokenGuard) getAuthUserName() string {
	return t.name + "_auth_user"
}

func (t *TokenGuard) getTokenName() string {
	return "api_token"
}

func (t *TokenGuard) cycleApiToken(user contracts.IUser) {
	token := t.random(60)
	user.SetRememberToken(token)
	t.provider.UpdateApiToken(user, token)
}

func (t *TokenGuard) hasValidCredentials(user contracts.IUser, credentials map[string]interface{}) bool {
	if user.IsEmpty() {
		return false
	}
	return t.provider.ValidateCredentials(user, credentials)
}

func (t *TokenGuard) getTokenFromRequest() string {
	if token := t.context.Query(t.getTokenName()); token != "" {
		return token
	}
	if token := t.context.PostForm(t.getTokenName()); token != "" {
		return token
	}
	return ""
}

func (t *TokenGuard) random(l int) string {
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

func (t *TokenGuard) updateToken(user contracts.IUser) {
	token := t.random(60)
	user.SetApiToken(token)
	t.provider.UpdateApiToken(user, token)
}
