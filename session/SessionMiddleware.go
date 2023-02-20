package session

import (
	"github.com/owenzhou/ginrbac/contracts"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/owenzhou/ginrbac/support/facades"
)

type SessionMiddleware struct {
	App contracts.IApplication
}

func (s *SessionMiddleware) Middleware() {
	secret := "cookie_store"
	if sessionSecret := facades.Config.SessionSecret; sessionSecret != "" {
		secret = sessionSecret
	}
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{MaxAge: 7200})
	s.App.GetEngine().Use(sessions.Sessions("session_id", store))
}
