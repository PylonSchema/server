package gateway

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	pingTick    = 10 * time.Second
	pongTimeout = (pingTick * 19) / 10
)

type Gateway struct {
	Upgrader websocket.Upgrader
	m        sync.RWMutex
	channels map[string][]*Client
}

func (g *Gateway) OpenGateway(c *gin.Context) {
	conn, err := g.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "internal server error",
			"trace code": "create upgrade connection",
			"error":      err,
		})
		return
	}
	client := &Client{
		conn:        conn,
		gatewayPipe: g,
	}

	go client.readHandler(pongTimeout)
	go client.writeHandler(pingTick)
}

func (g *Gateway) Inject(c *Client) error { // inject client to channel
	return nil
}

func (g *Gateway) Remove(c *Client) error { //  remove client from channel
	return nil
}
