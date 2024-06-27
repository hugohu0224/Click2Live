package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (cm *ClientManager) AddClient(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[conn] = true
}

func (cm *ClientManager) RemoveClient(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.clients, conn)
}

func (cm *ClientManager) GetClients() map[*websocket.Conn]bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	// return a copy of the clients map to avoid concurrency issues
	clientsCopy := make(map[*websocket.Conn]bool)
	for k, v := range cm.clients {
		clientsCopy[k] = v
	}
	return clientsCopy
}
