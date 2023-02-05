package auth

import "github.com/gin-gonic/gin"

var (
	AccessDenied = gin.H{
		"message": "access denied",
	}
	InternalServerError = gin.H{
		"message": "server error, try again",
	}
)
