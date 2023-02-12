package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setMiddleWare(r *gin.Engine, config *conf) {
	// cookie-based session store
	r.Use(setSession(config.Secret.Session))

	// set sentry
	// r.Use(setSentry(config.Sentry.Dsn))
}

func setSession(key string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(key))
	return sessions.Sessions("mySession", store)
}

// func setSentry(dsn string) gin.HandlerFunc {
// 	if err := sentry.Init(sentry.ClientOptions{
// 		Dsn:           dsn,
// 		EnableTracing: true,
// 		// Set TracesSampleRate to 1.0 to capture 100%
// 		// of transactions for performance monitoring.
// 		// We recommend adjusting this value in production,
// 		TracesSampleRate: 1.0,
// 	}); err != nil {
// 		fmt.Printf("Sentry initialization failed: %v\n", err)
// 	}
// 	return sentrygin.New(sentrygin.Options{})
// }
