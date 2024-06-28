package ws

import (
	"go.uber.org/zap"
	"sync"
)

type Hub struct {
	id            string
	clientManager *ClientManager
	broadcast     chan *BroadcastScore
	register      chan *Client
	unregister    chan *Client
	mu            sync.RWMutex
}

type HubManager struct {
	Hubs     map[*Hub]bool
	HubsById map[string]*Hub
	Mu       sync.RWMutex
}

func NewHub(id string) *Hub {
	return &Hub{
		id:            id,
		clientManager: NewClientManager(),
		broadcast:     make(chan *BroadcastScore),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		mu:            sync.RWMutex{},
	}
}

func (hm *HubManager) AddNewHub(h *Hub) {
	hm.Mu.Lock()
	hm.Hubs[h] = true
	hm.HubsById[h.id] = h
	hm.Mu.Unlock()
}

func (hm *HubManager) GetHubById(id string) (*Hub, bool) {
	hm.Mu.RLock()
	defer hm.Mu.RUnlock()
	hub, exists := hm.HubsById[id]
	return hub, exists
}

func (h *Hub) Run() {
	zap.S().Infof("game hub %s is running", h.id)
	for {
		select {
		// register
		case client := <-h.register:
			h.clientManager.AddClient(client)

		// unregister
		case client := <-h.unregister:
			h.clientManager.RemoveClient(client)

		// broadcast
		case message := <-h.broadcast:
			for client := range h.clientManager.GetClients() {
				select {
				case client.send <- message:
				// no receiver
				default:
					close(client.send)
					h.clientManager.RemoveClient(client)
				}
			}
		}
	}
}

func (hm *HubManager) StartAllHubs() {
	hm.Mu.RLock()
	defer hm.Mu.RUnlock()
	for hub := range hm.Hubs {
		go func(h *Hub) {
			h.Run()
		}(hub)
	}
}
