package gateway

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeDeadline = 2 * time.Second
)

type pipe interface {
	Inject(c *Client) error
	Remove(c *Client) error
}

type Client struct {
	conn         *websocket.Conn
	once         sync.Once
	writeChannel chan *interface{}
	gatewayPipe  pipe
}

// close socket connection & remove client from gateway
func (c *Client) closeConnection() {
	c.once.Do(func() {
		c.conn.Close()
	})
}

func (c *Client) readHandler(pongTimeout time.Duration) {
	defer c.closeConnection()
	c.conn.SetReadLimit(64)
	c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
		return nil
	})

	for {
		if _, _, err := c.conn.NextReader(); err != nil { // 2th value will be used later, message
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
			err := c.conn.WriteJSON(&message)
			if err != nil {
				return
			}
		case <-pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}

}
