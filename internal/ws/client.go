package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
)

type Score struct {
	mu    sync.Mutex
	Fire  int `json:"fire"`
	Water int `json:"water"`
	Food  int `json:"food"`
}

func (s *Score) Update(fire, water, food int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Fire += fire
	s.Water += water
	s.Food += food
}

type ClickMessage struct {
	UserId uuid.UUID `json:"userId"`
	Score
}

type PlayerScore struct {
	Id uuid.UUID `json:"id"`
	Score
}

type GlobalScore struct {
	Score
}

type BroadcastScore struct {
	UserId uuid.UUID    `json:"userId"`
	Ps     *PlayerScore `json:"ps"`
	Gs     *GlobalScore `json:"gs"`
}

type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan *BroadcastScore
	ps   *PlayerScore
}

func (c *Client) readPump() {
	for {
		var msg ClickMessage
		// get score from player's click
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			zap.S().Errorf("error reading message: %v", err)
			c.hub.clientManager.RemoveClient(c)
			break
		}

		// inject id for broadcast process
		c.ps.Id = msg.UserId

		// update
		c.hub.clientManager.gs.Update(msg.Fire, msg.Water, msg.Food)
		c.ps.Update(msg.Fire, msg.Water, msg.Food)

		// start to broadcast
		bs := &BroadcastScore{
			UserId: msg.UserId,
			Ps:     c.ps,
			Gs:     c.hub.clientManager.gs,
		}
		c.hub.broadcast <- bs
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteJSON(message)
		}
	}
}

type ClientManager struct {
	clients map[*Client]bool
	gs      *GlobalScore
	mu      sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[*Client]bool),
		gs:      &GlobalScore{},
		mu:      sync.RWMutex{},
	}
}

func (cm *ClientManager) AddClient(client *Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[client] = true
}

func (cm *ClientManager) RemoveClient(client *Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.clients, client)
}

func (cm *ClientManager) GetClients() map[*Client]bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.clients
}
