package ws

import (
	"github.com/gorilla/websocket"
)

func ServeWs(conn *websocket.Conn, hub *Hub) {
	defer conn.Close()

	client := &Client{
		id:   "",
		hub:  hub,
		conn: conn,
		send: make(chan *BroadcastScore),
		ps:   &PlayerScore{},
	}

	hub.clientManager.AddClient(client)

	go client.readPump()
	go client.writePump()

	<-make(chan struct{})

}
