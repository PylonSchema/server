package gateway

import (
	"sync"
	"time"

	"github.com/PylonSchema/server/auth"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

const (
	writeDeadline = 2 * time.Second
)

const (
	MessageHeartbeat      = 0
	MessageAuthentication = 1
	MessageData           = 2
	MessageClose          = 10
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
	username     string // client username
	uuid         string // client uuid
}

// close socket connection & remove client from gateway
func (c *Client) closeConnection() {
	c.once.Do(func() {
		d := map[string]interface{}{"type": "authorized error"}
		command, _ := json.Marshal(&Message{
			Op: 10,
			D:  d,
		})
		c.conn.WriteMessage(websocket.TextMessage, command)
		c.conn.Close()
	})
}

func (c *Client) defineClient(message *Message) {
	claims, err := c.gatewayPipe.Auth((message.D["token"]).(string))
	if err != nil {
		c.closeConnection()
		return
	}
	c.username = claims.Username
	c.uuid = claims.UserUUID
}

func (c *Client) GatewayInject() {
	c.gatewayPipe.Inject(c)
}

func (c *Client) GatewayRemove() {
	c.gatewayPipe.Remove(c)
}

func (c *Client) readHandler(pongTimeout time.Duration) {
	defer c.closeConnection()
	c.conn.SetReadLimit(2048)
	c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
	c.conn.SetPongHandler(func(_ string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
		return nil
	})

	// only allow auth, heartbeat, close connection messages
	for {
		var message Message
		var isNext = false
		err := c.conn.ReadJSON(&message)
		if err != nil {
			// need error handle
			// like websocket.ErrlimitRead
			return
		}
		switch message.Op {
		case MessageHeartbeat:
			c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
		case MessageAuthentication:
			// message authorized implements
			c.defineClient(&message)
			isNext = true
		case MessageClose:
			command, err := json.Marshal(&Message{
				Op: 10,
				D:  nil,
			})
			if err != nil {
				return // need websocket write error
			}
			c.conn.WriteMessage(websocket.TextMessage, command)
			return
		}
		if isNext {
			break
		}
	}

	c.GatewayInject() // inject client in gateway

	// implement except only authentication
	for {
		var message Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			// need error handle
			// like websocket.ErrlimitRead
			return
		}
		switch message.Op {
		case MessageHeartbeat:
			c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
		case MessageData:
			c.writeChannel <- &message // reply test
			// boardcast to channel implements need
		case MessageClose:
			command, err := json.Marshal(&Message{
				Op: 10,
				D:  nil,
			})
			if err != nil {
				return // need websocket write error
			}
			c.conn.WriteMessage(websocket.TextMessage, command)
			return
		}
	}
}

func (c *Client) writeHandler(pingTick time.Duration) {
	pingTicker := time.NewTicker(pingTick)
	defer func() {
		pingTicker.Stop()
		c.closeConnection()
	}()

	for {
		select {
		case message := <-c.writeChannel:
			c.conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			command, err := json.Marshal(message)
			if err != nil {
				return // need websocket write error
			}
			err = c.conn.WriteMessage(websocket.TextMessage, command)
			if err != nil {
				return // need websocket write error
			}
		case <-pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			pingMessage := Message{
				Op: 0,
				D:  nil,
			}
			command, err := json.Marshal(pingMessage)
			if err != nil {
				return // need websocket write error
			}
			err = c.conn.WriteMessage(websocket.TextMessage, command)
			if err != nil {
				return // need websocket write error
			}
		}
	}

}
