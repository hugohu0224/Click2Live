package ws

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
	"tap2live/internal/models"
)

var (
	gs        = &models.GlobalScore{}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan *models.GlobalScore)
	clientMu  sync.Mutex
)

func UpdateScores(gs *models.GlobalScore, ps *models.PlayerScore, fire, water, food int) {
	// global score
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	// player score
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()

	// update global score
	gs.Fire += fire
	gs.Water += water
	gs.Food += food

	// update global player score
	ps.Fire += fire
	ps.Water += water
	ps.Food += food
}

func HandleConnection(conn *websocket.Conn) {
	var ps = &models.PlayerScore{}

	defer conn.Close()
	clients[conn] = true

	// init
	gs.Mutex.Lock()
	initialScore := models.ClickMessage{
		Score: models.Score{
			Fire:  gs.Fire,
			Water: gs.Water,
			Food:  gs.Food,
		},
	}

	gs.Mutex.Unlock()
	if err := conn.WriteJSON(initialScore); err != nil {
		zap.S().Errorf("error initial state: %v", err)
		clientMu.Lock()
		delete(clients, conn)
		clientMu.Unlock()
		return
	}

	// start to update
	for {
		var msg models.ClickMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			zap.S().Errorf("error reading message: %v", err)
			clientMu.Lock()
			delete(clients, conn)
			clientMu.Unlock()
			break
		}

		UpdateScores(gs, ps, msg.Fire, msg.Water, msg.Food)

		broadcast <- gs
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			go func(client *websocket.Conn, gs *models.GlobalScore) {
				err := client.WriteJSON(gs)
				if err != nil {
					zap.S().Errorf("error sending message: %v", err)
					client.Close()
					clientMu.Lock()
					delete(clients, client)
					clientMu.Unlock()
				}
			}(client, msg)
		}
	}
}
