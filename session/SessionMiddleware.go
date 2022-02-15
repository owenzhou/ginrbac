package session

import (
	"ginrbac/bootstrap/contracts"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

type SessionMiddleware struct {
	App contracts.IApplication
}

func (s *SessionMiddleware) Middleware() {
	store := cookie.NewStore([]byte("cookie_store"))
	s.App.GetEngine().Use(sessions.Sessions("session_id", store))
}
