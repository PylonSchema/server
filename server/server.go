package server

import (
	"github.com/BurntSushi/toml"
	githubAuth "github.com/devhoodit/sse-chat/auth/github"
	"github.com/devhoodit/sse-chat/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type conf struct {
	Database   *databaseInfo
	Sentry     *sentryInfo
	githubAuth *oauth2Info
}

type databaseInfo struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Address  string `toml:"address"`
	Port     string `toml:"port"`
}

type sentryInfo struct {
	Dsn string `toml:"dsn"`
}

type oauth2Info struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURL  string `toml:"redirect_url"`
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	var conf conf
	_, err := toml.DecodeFile("conf.toml", &conf)
	if err != nil {
		panic(err)
	}

	d, err := database.Connect(conf.Database.Username, conf.Database.Password, conf.Database.Address, conf.Database.Port)
	if err != nil {
		panic(err)
	}

	err = d.AutoMigration() // auto migration, check table is Exist, if not create
	if err != nil {
		panic(err)
	}

	// MiddleWare setting, server/middleware.go
	setMiddleWare(r, &conf)

	// make router

	// github Oauth router
	githubRouter := githubAuth.Github{
		DB: d,
		OAuthConfig: &oauth2.Config{
			ClientID:     conf.githubAuth.ClientID,
			ClientSecret: conf.githubAuth.ClientSecret,
			RedirectURL:  conf.githubAuth.RedirectURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		}}

	r.GET("/", func(c *gin.Context) {
	})

	auth := r.Group("/auth")
	{
		sse := auth.Group("/sse")
		{
			sse.POST("/login")
			sse.POST("/create")
		}
		github := auth.Group("/github")
		{
			github.GET("/login", githubRouter.Login)
			github.GET("/callback", githubRouter.Callback)
		}
	}

	return r
}
