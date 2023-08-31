package server

import (
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/PylonSchema/server/api"
	"github.com/PylonSchema/server/api/gateway"
	"github.com/PylonSchema/server/auth"
	githubAuth "github.com/PylonSchema/server/auth/github"
	pylonAuth "github.com/PylonSchema/server/auth/origin"
	"github.com/PylonSchema/server/database"
	"github.com/PylonSchema/server/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var SecretKey *secret

func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{"http://localhost:5500"},
			AllowMethods:     []string{"POST"},
			AllowHeaders:     []string{"Origin", "content-type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

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

	auth := auth.New(jwtAuth, d)

	// github Oauth router
	githubAuthRouter := githubAuth.Github{
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

	versionAPI := api.NewVersionAPI("localhost", 8080)

	r.GET("/version", versionAPI.VersionHandler)

	gateway := gateway.New(jwtAuth, d)

	r.GET("/gateway", gateway.CreateGatewayHandler)

	messageAPI := api.NewMessageAPI(gateway, d)

	messageRouter := r.Group("/message").Use(jwtAuth.AuthorizeRequiredMiddleware())
	{
		messageRouter.POST("/", messageAPI.CreateMessageHandler)
	}

	channelAPI := api.ChannelAPI{
		DB: d,
	}

	channelRouter := r.Group("/channel").Use(jwtAuth.AuthorizeRequiredMiddleware())
	{
		channelRouter.GET("/", channelAPI.GetChannelIdsHandler)        // get channel ids
		channelRouter.POST("/", channelAPI.CreateChannelHandler)       // create channel
		channelRouter.DELETE("/", channelAPI.RemoveChannelHandler)     // delete channel
		channelRouter.POST("/join/:id", channelAPI.JoinChannelHandler) // join channel
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
			github.GET("/login", githubAuthRouter.LoginHandler)
			github.GET("/callback", githubAuthRouter.CallbackHandler)
		}

		authRouter.Use(jwtAuth.AuthorizeRequiredMiddleware()).GET("/token", func(ctx *gin.Context) {
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
		authRouter.GET("/blacklist", auth.BlacklistHandler).Use(jwtAuth.AuthorizeRequiredMiddleware())
		authRouter.GET("/usertokens", auth.GetTokenHandler).Use(jwtAuth.AuthorizeRequiredMiddleware())
	}

	return r
}
