package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"sync"
)

func ServeWs(conn *websocket.Conn, hub *Hub) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer conn.Close()

	client := &Client{
		id:   "",
		hub:  hub,
		conn: conn,
		send: make(chan *BroadcastScore),
		ps:   &PlayerScore{},
	}

	hub.clientManager.AddClient(client)
	defer hub.clientManager.RemoveClient(client)

	// start
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		client.readPump(ctx)
	}()

	go func() {
		defer wg.Done()
		client.writePump(ctx)
	}()

	wg.Wait()
}
