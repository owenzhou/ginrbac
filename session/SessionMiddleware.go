package session

import (
	"github.com/owenzhou/ginrbac/contracts"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

type SessionMiddleware struct {
	App contracts.IApplication
}

func (s *SessionMiddleware) Middleware() {
	store := cookie.NewStore([]byte("cookie_store"))
	store.Options(sessions.Options{MaxAge: 7200})
	s.App.GetEngine().Use(sessions.Sessions("session_id", store))
}
