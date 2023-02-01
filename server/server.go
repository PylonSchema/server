package server

import (
	"github.com/BurntSushi/toml"
	githubAuth "github.com/devhoodit/sse-chat/auth/github"
	"github.com/devhoodit/sse-chat/database"
	"github.com/gin-gonic/gin"
)

type conf struct {
	Database *databaseInfo
}

type databaseInfo struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Address  string `toml:"address"`
	Port     string `toml:"port"`
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// MiddleWare setting, server/middleware.go
	setMiddleWare(r)

	var conf conf
	_, err := toml.DecodeFile("conf.toml", &conf)
	if err != nil {
		panic(err)
	}

	d, err := database.Connect(conf.Database.Username, conf.Database.Password, conf.Database.Address, conf.Database.Port)

	if err != nil {
		panic(err)
	}

	githubRouter := githubAuth.Github{DB: d}

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
			github.GET("/login", githubRouter.RenderAuthView)
			github.GET("/callback", githubRouter.Authenticate)
		}
	}

	return r
}
