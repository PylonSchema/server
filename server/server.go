package server

import (
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
			github.POST("/login")
			github.GET("/callback")
		}
	}

	return r
}
