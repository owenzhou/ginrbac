package auth

import (
	"fmt"
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support/facades"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWTGuard struct {
	name     string
	provider contracts.IUserProvider
	context  *gin.Context
}

func (j *JWTGuard) ID() string {
	return j.User().GetAuthIdentifier()
}

func (j *JWTGuard) Check() bool {
	return !j.User().IsEmpty()
}

func (j *JWTGuard) Guest() bool {
	return !j.Check()
}

func (j *JWTGuard) Validate(credentials map[string]interface{}) bool {
	user := j.provider.RetrieveByCredentials(credentials)
	return j.hasValidCredentials(user, credentials)
}

func (j *JWTGuard) WithContext(c *gin.Context) contracts.IAuth {
	j.context = c
	return j
}

func (j *JWTGuard) SetUser(user contracts.IUser) {
	j.context.Set(j.getAuthUserName(), user)
}

func (j *JWTGuard) User() contracts.IUser {
	var user contracts.IUser
	u, exists := j.context.Get(j.getAuthUserName())
	//如果gin.Context里存在用户信息
	if exists {
		user = u.(contracts.IUser)
		return user
	}
	//如果request里不存在api_token
	token := j.getTokenFromRequest()
	if token == "" {
		user = &GenericUser{}
		return user
	}
	cliams, err := j.parseToken(token)
	if err != nil {
		facades.Log.Error("login err: ", err)
		user = &GenericUser{}
	} else {
		user = cliams.User
		j.SetUser(user)
	}
	return user
}

func (j *JWTGuard) Attempt(credentials map[string]interface{}, remember ...bool) bool {
	user := j.provider.RetrieveByCredentials(credentials)
	if j.hasValidCredentials(user, credentials) {
		j.Login(user, remember...)

		return true
	}
	return false
}

func (j *JWTGuard) Login(user contracts.IUser, remember ...bool) {
	if !user.IsEmpty() {
		j.updateToken(user)
	}
	j.SetUser(user)
}

func (j *JWTGuard) Logout() {

}

//-------------------------

func (j *JWTGuard) updateToken(user contracts.IUser) {
	token, _ := j.createToken(user)
	user.SetApiToken(token)
	j.provider.UpdateApiToken(user, token)
}

func (j *JWTGuard) getAuthUserName() string {
	return j.name + "_auth_user"
}

func (j *JWTGuard) getTokenName() string {
	return "api_token"
}

func (j *JWTGuard) getTokenFromRequest() string {
	//从url获取
	if token := j.context.Query(j.getTokenName()); token != "" {
		return token
	}
	//从表单获取
	if token := j.context.PostForm(j.getTokenName()); token != "" {
		return token
	}
	//从cookie获取
	if token, err := j.context.Cookie(j.getTokenName()); err == nil {
		return token
	}
	//从header获取
	if token := j.context.Request.Header.Get("Authorization"); token != "" {
		return strings.Replace(token, "Bearer ", "", 1)
	}
	return ""
}

func (j *JWTGuard) hasValidCredentials(user contracts.IUser, credentials map[string]interface{}) bool {
	if user.IsEmpty() {
		return false
	}
	return j.provider.ValidateCredentials(user, credentials)
}

type CustomCliams struct {
	User *GenericUser
	*jwt.StandardClaims
}

func (j *JWTGuard) createToken(iUser contracts.IUser) (string, error) {
	user := iUser.(*GenericUser)
	delete(user.Attributes, "password")
	delete(user.Attributes, "remember_token")
	delete(user.Attributes, "api_token")
	// Create the Claims
	claims := CustomCliams{
		user,
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(facades.Config.JWT.ExpiresTime) * time.Second).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(facades.Config.JWT.SignKey))
	if err != nil {
		return "", err
	}
	return ss, err
}

func (j *JWTGuard) parseToken(tokenString string) (*CustomCliams, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomCliams{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method:%v", t.Header["alg"])
		}
		return []byte(facades.Config.JWT.SignKey), nil
	})
	if cliams, ok := token.Claims.(*CustomCliams); ok && token.Valid {
		return cliams, nil
	} else {
		return nil, err
	}
}
