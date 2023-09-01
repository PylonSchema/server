package gateway

import (
	"net/http"
	"sync"
	"time"

	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/database"
	"github.com/PylonSchema/server/model"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	pingTick    = 40 * time.Second
	pongTimeout = (pingTick * 19) / 10
)

type Database interface {
	GetChannelsByUserUUID(uuid uuid.UUID) (*[]model.ChannelMember, error)
}

type Gateway struct {
	Upgrader    websocket.Upgrader
	m           *sync.RWMutex
	channels    map[uint][]*Client
	JwtAuth     *auth.JwtAuth
	db          Database
	UserSockets *userSockets
}

type userSockets struct {
	m           *sync.RWMutex
	userSockets map[string][]*Client
}

func New(jwtAuth *auth.JwtAuth, db *database.Database) *Gateway {
	return &Gateway{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool { // origin check for dev, allow all origin
				return true
			},
		},
		JwtAuth:  jwtAuth,
		channels: make(map[uint][]*Client),
		m:        new(sync.RWMutex),
		db:       db,
		UserSockets: &userSockets{
			m:           new(sync.RWMutex),
			userSockets: make(map[string][]*Client),
		},
	}
}

func (g *Gateway) Inject(c *Client) error { // inject client to channel
	channels, err := g.db.GetChannelsByUserUUID(c.uuid)
	if err != nil {
		return err
	}

	g.m.Lock()
	for _, channel := range *channels {
		g.channels[channel.ChannelId] = append(g.channels[channel.ChannelId], c)
	}
	g.m.Unlock()

	userUUID := c.uuid.String()
	g.UserSockets.m.Lock()
	g.UserSockets.userSockets[userUUID] = append(g.UserSockets.userSockets[userUUID], c)
	g.UserSockets.m.Unlock()

	return nil
}

func (g *Gateway) Remove(c *Client) error { //  remove client from channel
	channelMembers, err := g.db.GetChannelsByUserUUID(c.uuid)
	if err != nil {
		return err
	}

	g.m.Lock()
	for _, channel := range *channelMembers {
		for i, client := range g.channels[channel.ChannelId] {
			if client != c {
				continue
			}
			g.channels[channel.ChannelId] = append(g.channels[channel.ChannelId][:i], g.channels[channel.ChannelId][i+1:]...)
			break
		}
	}
	g.m.Unlock()

	userUUID := c.uuid.String()
	g.UserSockets.m.Lock()
	userSockets := g.UserSockets.userSockets[userUUID]
	for i, userSocket := range userSockets {
		if userSocket.token == c.token {
			g.UserSockets.userSockets[userUUID] = append(g.UserSockets.userSockets[userUUID][:i], g.UserSockets.userSockets[userUUID][i+1:]...)
			break
		}
	}
	g.UserSockets.m.Unlock()
	return nil
}

func (g *Gateway) Boardcast(channelId uint, message *Message) error {
	g.m.RLock()
	defer g.m.RUnlock()
	clients, ok := g.channels[channelId]
	if !ok {
		return nil
	}
	for _, client := range clients {
		client.writeChannel <- message
	}
	return nil
}
