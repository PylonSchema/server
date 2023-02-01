package server

import (
	githubAuth "github.com/devhoodit/sse-chat/auth/github"
	// "github.com/devhoodit/sse-chat/database"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// MiddleWare setting, server/middleware.go
	setMiddleWare(r)

	// d, err := database.Connect()

	// if err != nil {
	// 	panic(err)
	// }

	r.GET("/", func(c *gin.Context) {
	})

	auth := r.Group("/auth")
	{
		sse := auth.Group("/sse")
		{
			sse.POST("/login")
		}
		github := auth.Group("/github")
		{
			github.GET("/login", githubAuth.RenderAuthView)
			github.GET("/callback", githubAuth.Authenticate)
		}
	}

	return r
}
