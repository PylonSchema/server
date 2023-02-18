package gateway

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/PylonSchema/server/auth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	pingTick    = 10 * time.Second
	pongTimeout = (pingTick * 19) / 10
)

type Database interface {
	GetChannelsByUserUUID(uuid string)
}

type Gateway struct {
	Upgrader websocket.Upgrader
	m        *sync.RWMutex
	channels map[string][]*Client
	JwtAuth  *auth.JwtAuth
	db       Database
}

func New(jwtAuth *auth.JwtAuth) *Gateway {
	return &Gateway{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool { // origin check for dev, allow all origin
				return true
			},
		},
		JwtAuth:  jwtAuth,
		channels: make(map[string][]*Client),
		m:        new(sync.RWMutex),
	}
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
		conn:         conn,
		gatewayPipe:  g,
		writeChannel: make(chan *Message),
		username:     "",
		uuid:         "",
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
