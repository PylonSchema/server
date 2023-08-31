package gateway

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

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
				c.syncErrorMessageWrite(0, "message marshal error, internal server error")
				return
			}
			err = c.conn.WriteMessage(websocket.TextMessage, command)
			if err != nil {
				c.syncErrorMessageWrite(0, "message write error, internal server error")
			}
		case <-pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			pingMessage := Message{
				Op: 0,
				D:  nil,
			}
			command, err := json.Marshal(pingMessage)
			if err != nil {
				c.syncErrorMessageWrite(0, "message marshal error, internal server error")
				return
			}
			err = c.conn.WriteMessage(websocket.TextMessage, command)
			if err != nil {
				c.syncErrorMessageWrite(0, "message write error, internal server error")
				return
			}
		}
	}

}

func (c *Client) syncErrorMessageWrite(code int, errorMessage string) error {
	err := c.syncMessageWrite(&map[string]interface{}{
		"op": MessageError,
		"d":  map[string]interface{}{"code": code, "data": errorMessage},
	})
	return err
}

func (c *Client) syncMessageWrite(data *map[string]interface{}) error {
	command, err := json.Marshal(data)
	if err != nil {
		return err
	}
	c.conn.WriteMessage(websocket.TextMessage, command)
	return nil
}
