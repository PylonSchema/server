package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (g *Gateway) CreateGatewayHandler(c *gin.Context) {
	conn, err := g.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
			"error":   err,
		})
		return
	}
	client := &Client{
		conn:         conn,
		gatewayPipe:  g,
		writeChannel: make(chan *Message),
		username:     "",
		uuid:         uuid.UUID{},
		token:        "",
		refreshToken: "",
	}

	go client.readHandler(pongTimeout)
	go client.writeHandler(pingTick)
}
