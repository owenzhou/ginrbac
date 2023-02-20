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
	lifetime := 7200
	if sessionSecret := facades.Config.Session.SecretKey; sessionSecret != "" {
		secret = sessionSecret
	}
	if lt := facades.Config.Session.LifeTime; lt != 0 {
		lifetime = lt
	}
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{MaxAge: lifetime})
	s.App.GetEngine().Use(sessions.Sessions("session_id", store))
}
