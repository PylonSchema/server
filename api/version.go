package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var version = map[string]string{
	"auth":    "0.0.1",
	"gateway": "0.0.1",
	"message": "0.0.1",
	"channel": "0.0.1",
}

var endpoint = map[string]string{
	"auth":    "/auth",
	"gateway": "/gateway",
	"message": "/message",
	"channel": "/channel",
}

type VersionAPI struct {
	Host string
}

func NewVersionAPI(host string, port int) *VersionAPI {
	if port != 80 {
		host = fmt.Sprintf("%s:%d", host, port)
	}
	return &VersionAPI{
		Host: host,
	}
}

func (v *VersionAPI) VersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"host":     v.Host,
		"version":  version,
		"endpoint": endpoint,
	})
}
