package server

import (
	githubAuth "github.com/devhoodit/sse-chat/auth/github"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// MiddleWare setting, server/middleware.go
	setMiddleWare(r)

	auth := r.Group("/auth")
	{
		sse := auth.Group("/sse")
		{
			sse.POST("/login")
		}
		github := auth.Group("/github")
		{
			github.GET("/login", githubAuth.RenderAuthView)
			github.GET("/callback")
		}
	}

	return r
}
