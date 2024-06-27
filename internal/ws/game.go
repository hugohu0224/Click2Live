package ws

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
	"tap2live/pkg/global"
	"tap2live/pkg/models"
)

var (
	cm        = NewClientManager()
	broadcast = make(chan *models.BroadcastScore)
)

func UpdateScores(gs *models.GlobalScore, ps *models.PlayerScore, fire, water, food int) {
	// lock
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()

	// update global score
	gs.Fire += fire
	gs.Water += water
	gs.Food += food

	// update player score
	ps.Fire += fire
	ps.Water += water
	ps.Food += food
}

func HandleConnection(conn *websocket.Conn) {
	var ps = &models.PlayerScore{}
	defer conn.Close()
	cm.AddClient(conn)

	// initial for join the game after started
	global.Gs.Mutex.Lock()
	initialScore := models.BroadcastScore{
		Ps: nil,
		Gs: &models.GlobalScore{
			Score: models.Score{
				Fire:  global.Gs.Fire,
				Water: global.Gs.Water,
				Food:  global.Gs.Food,
			},
			Mutex: sync.Mutex{},
		},
	}
	global.Gs.Mutex.Unlock()

	if err := conn.WriteJSON(initialScore); err != nil {
		zap.S().Errorf("error initial state: %v", err)
		cm.RemoveClient(conn)
		return
	}

	// start to update score
	for {
		var msg models.ClickMessage
		// get score from player's click
		err := conn.ReadJSON(&msg)
		if err != nil {
			zap.S().Errorf("error reading message: %v", err)
			cm.RemoveClient(conn)
			break
		}

		// inject id for broadcast process
		ps.Id = msg.UserId

		// update
		UpdateScores(global.Gs, ps, msg.Fire, msg.Water, msg.Food)

		bs := &models.BroadcastScore{
			UserId: msg.UserId,
			Ps:     ps,
			Gs:     global.Gs,
		}

		broadcast <- bs

	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range cm.GetClients() {
			go func(client *websocket.Conn, bs *models.BroadcastScore) {
				err := client.WriteJSON(bs)
				if err != nil {
					zap.S().Errorf("error sending message: %v", err)
					client.Close()
					cm.RemoveClient(client)
				}
			}(client, msg)
		}
	}

}
