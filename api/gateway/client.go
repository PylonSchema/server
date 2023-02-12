package gateway

import "github.com/gorilla/websocket"

type Client struct {
	conn *websocket.Conn
}
