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
	MessageHeartbeat  = 0
	MessageAuthorized = 1
	MessageData       = 2
)

type pipe interface {
	Inject(c *Client) error
	Remove(c *Client) error
}

type Client struct {
	conn         *websocket.Conn
	once         sync.Once
	writeChannel chan *Message
	gatewayPipe  pipe
}

// close socket connection & remove client from gateway
func (c *Client) closeConnection() {
	c.once.Do(func() {
		fmt.Println("connection closed")
		c.conn.Close()
	})
}

func (c *Client) readHandler(pongTimeout time.Duration) {
	defer c.closeConnection()
	c.conn.SetReadLimit(1024)
	c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
		return nil
	})

	for {
		var message Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			return // need socket read error
		}
		switch message.Op {
		case MessageHeartbeat:
			c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
		case MessageData:
			c.writeChannel <- &message // reply test
		case MessageAuthorized:
			// message authorized implements
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
			err := c.conn.WriteJSON(&message)
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
