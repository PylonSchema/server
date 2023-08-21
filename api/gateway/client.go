package gateway

import (
	"sync"
	"time"

	"github.com/PylonSchema/server/auth"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeDeadline = 2 * time.Second
)

const (
	MessageHeartbeat      = 0  // check websocket alive
	MessageAuthentication = 1  // check authentication
	MessageData           = 2  // payload is data (message etc.)
	MessageEvent          = 8  // event message (change of user or user check notification, authorized etc.)
	MessageError          = 9  // error occur in several reason (authentication error, websocket write length err etc.)
	MessageClose          = 10 // websocket close
)

type pipe interface {
	Inject(c *Client) error
	Remove(c *Client) error
	Auth(tokenString string) (*auth.AuthTokenClaims, error)
}

type Client struct {
	conn         *websocket.Conn
	once         sync.Once
	writeChannel chan *Message
	gatewayPipe  pipe
	username     string    // client username
	uuid         uuid.UUID // client uuid
}

// close socket connection & remove client from gateway
func (c *Client) closeConnection() {
	c.once.Do(func() {
		d := map[string]interface{}{"type": "close connection"}
		command, _ := json.Marshal(&Message{
			Op: 10,
			D:  d,
		})
		c.conn.WriteMessage(websocket.TextMessage, command)
		c.conn.Close()
	})
}

func (c *Client) defineClient(message *Message) error {
	claims, err := c.gatewayPipe.Auth((message.D["token"]).(string))
	if err != nil {
		return err
	}
	c.username = claims.Username
	c.uuid = claims.UserUUID
	return nil
}

func (c *Client) GatewayInject() error {
	return c.gatewayPipe.Inject(c)
}

func (c *Client) GatewayRemove() {
	c.gatewayPipe.Remove(c)
}
