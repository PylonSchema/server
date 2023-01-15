package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setMiddleWare(r *gin.Engine) {
	// cookie-based session store
	r.Use(setSession())
}

func setSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte("secret"))
	return sessions.Sessions("mySession", store)
}
