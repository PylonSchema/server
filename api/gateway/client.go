package gateway

import (
	"fmt"
	"sync"
	"time"

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
	Auth(tokenString string) error
}

type Client struct {
	conn         *websocket.Conn
	once         sync.Once
	writeChannel chan *Message
	gatewayPipe  pipe
	username     string // client username
	uuid         string // client uuid
	authorized   bool   // client authorized?, client username & uuid is defined after authorized, if not authorized can't do anything
}

// close socket connection & remove client from gateway
func (c *Client) closeConnection() {
	c.once.Do(func() {
		fmt.Println("connection closed")
		c.conn.Close()
	})
}

func (c *Client) defineClient(message *Message) {
	if c.authorized {
		return
	}
	err := c.gatewayPipe.Auth((message.D["token"]).(string))
	if err != nil {
		d := map[string]interface{}{"type": "authorized error"}
		command, err := json.Marshal(&Message{
			Op: 10,
			D:  d,
		})
		if err != nil {
			return // need websocket write error
		}
		c.conn.WriteMessage(websocket.TextMessage, command)
		c.closeConnection()
		return
	}
	c.authorized = true
	c.username = "client username"
	c.uuid = "client uuid"
}

func (c *Client) readHandler(pongTimeout time.Duration) {
	defer c.closeConnection()
	c.conn.SetReadLimit(2048)
	c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
		return nil
	})

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
		case MessageAuthentication:
			// message authorized implements
			c.defineClient(&message)
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
