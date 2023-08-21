package gateway

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

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
			err = c.defineClient(&message)
			if err != nil {
				c.syncErrorMessageWrite(0, "authorized error")
				return
			}
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
		case MessageClose:
			command, err := json.Marshal(&Message{
				Op: 10,
				D:  nil,
			})
			if err != nil {
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, command)
			return
		}
	}
}
