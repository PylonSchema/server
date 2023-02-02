package github

import "github.com/gin-gonic/gin"

var (
	AccessDenied = gin.H{
		"message": "access denied",
	}
)
