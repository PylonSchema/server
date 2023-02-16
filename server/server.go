package server

import (
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/PylonSchema/server/api/gateway"
	"github.com/PylonSchema/server/auth"
	githubAuth "github.com/PylonSchema/server/auth/github"
	"github.com/PylonSchema/server/database"
	"github.com/PylonSchema/server/store"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var SecretKey *secret

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// load config form conf.toml
	var conf conf
	_, err := toml.DecodeFile("conf.toml", &conf)
	if err != nil {
		panic(err)
	}

	// database setting
	d, err := database.New(conf.Database.Username, conf.Database.Password, conf.Database.Address, conf.Database.Port)
	if err != nil {
		panic(err)
	}
	err = d.AutoMigration() // auto migration, check table is Exist, if not create
	if err != nil {
		panic(err)
	}

	//redis setting
	store, err := store.New(&redis.Options{
		Addr:     conf.Store.Address,
		Password: conf.Store.Password, // no password set
		DB:       conf.Store.Db,       // use default DB
	})
	if err != nil {
		panic(err)
	}

	// MiddleWare setting, server/middleware.go
	setMiddleWare(r, &conf)

	jwtAuth := &auth.JwtAuth{
		Secret: conf.Secret.Jwtkey,
		DB:     d,
		Store:  store,
	}

	auth := &auth.Auth{
		JwtAuth: jwtAuth,
	}

	// github Oauth router
	githubRouter := githubAuth.Github{
		DB:      d,
		JwtAuth: jwtAuth,
		OAuthConfig: &oauth2.Config{
			ClientID:     conf.Oauth["github"].Id,
			ClientSecret: conf.Oauth["github"].Secret,
			RedirectURL:  conf.Oauth["github"].Redirect,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		}}

	r.GET("/", func(c *gin.Context) {
	})

	gateway := &gateway.Gateway{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool { // origin check for dev, allow all origin
				return true
			},
		},
		JwtAuth: jwtAuth,
	}

	r.GET("/gateway", gateway.OpenGateway)

	messageRouter := r.Group("/message").Use(jwtAuth.AuthorizeRequired())
	{
		messageRouter.POST("/")
	}

	authRouter := r.Group("/auth")
	{
		sse := authRouter.Group("/sse")
		{
			sse.GET("/login")
			sse.POST("/create")
		}
		github := authRouter.Group("/github")
		{
			github.GET("/login", githubRouter.Login)
			github.GET("/callback", githubRouter.Callback)
		}
		authRouter.Use(jwtAuth.AuthorizeRequired()).GET("/token", func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"message": "token is vaild",
			})
		})
		authRouter.GET("/blacklist", auth.Blacklist)
	}

	return r
}
