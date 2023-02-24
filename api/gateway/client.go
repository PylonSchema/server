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

func (c *Client) GatewayInject() error {
	return c.gatewayPipe.Inject(c)
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

	err := c.GatewayInject() // inject client in gateway
	if err != nil {
		c.syncMessageWrite(&map[string]interface{}{
			"Op": MessageError,
			"d": map[string]interface{}{
				"code":    0, // this is sample code, need change follow protocol rule ---------------------------------------
				"message": "gateway inject error",
			},
		})
		return
	}
	// authorized
	c.syncMessageWrite(&map[string]interface{}{
		"Op": MessageEvent,
		"d":  "",
	})

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
			c.writeChannel <- &message // reply test, message boardcast implement will be add in POST api request
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
		case message := <-c.writeChannel: // this channel message will be triggered by message POST api
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

func (c *Client) syncMessageWrite(data *map[string]interface{}) error {
	command, err := json.Marshal(data)
	if err != nil {
		return err
	}
	c.conn.WriteMessage(websocket.TextMessage, command)
	return nil
}
