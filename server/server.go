package server

import (
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/PylonSchema/server/api"
	"github.com/PylonSchema/server/api/gateway"
	"github.com/PylonSchema/server/auth"
	githubAuth "github.com/PylonSchema/server/auth/github"
	pylonAuth "github.com/PylonSchema/server/auth/origin"
	"github.com/PylonSchema/server/database"
	"github.com/PylonSchema/server/store"
	"github.com/gin-gonic/gin"
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

	jwtAuth := auth.NewJwtAuth(d, store, conf.Secret.Jwtkey)

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
		},
	}

	channelAPI := api.ChannelAPI{
		DB: d,
	}

	r.GET("/", func(c *gin.Context) {
	})

	gateway := gateway.New(jwtAuth, d)

	r.GET("/gateway", gateway.OpenGateway)

	messageAPI := api.NewMessageAPI(gateway, d)

	messageRouter := r.Group("/message").Use(jwtAuth.AuthorizeRequired())
	{
		messageRouter.POST("/", messageAPI.CreateMessage)
	}

	userRouter := r.Group("/user").Use(jwtAuth.AuthorizeRequired())
	{
		userRouter.GET("/channel")
	}

	channelRouter := r.Group("/channel").Use(jwtAuth.AuthorizeRequired())
	{
		channelRouter.GET("/", channelAPI.GetChannelIds)        // get channel ids
		channelRouter.POST("/", channelAPI.CreateChannel)       // create channel
		channelRouter.DELETE("/", channelAPI.RemoveChannel)     // delete channel
		channelRouter.POST("/join/:id", channelAPI.JoinChannel) // join channel
	}

	pylonAuthAPI := pylonAuth.New(d, jwtAuth)

	authRouter := r.Group("/auth")
	{
		pylon := authRouter.Group("/pylon")
		{
			pylon.POST("/login", pylonAuthAPI.LoginAccountHandler)
			pylon.POST("/create", pylonAuthAPI.CreateAccountHandler)
		}
		github := authRouter.Group("/github")
		{
			github.GET("/login", githubRouter.Login)
			github.GET("/callback", githubRouter.Callback)
		}
		authRouter.Use(jwtAuth.AuthorizeRequired()).GET("/token", func(ctx *gin.Context) {
			t, e := ctx.Cookie("token")
			if e != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"token": "internal server error",
				})
			}
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"token": t,
			})
		})
		authRouter.GET("/blacklist", auth.Blacklist)
	}

	return r
}
