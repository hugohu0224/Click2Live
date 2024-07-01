package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"tap2live/internal/global"
	"tap2live/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  64,
	WriteBufferSize: 64,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsEndpoint(c *gin.Context) {
	zap.S().Infof("start to upgrade connection")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.S().Errorf("failed to upgrade connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
		return
	}
	zap.S().Infof("successfully upgraded connection to WebSocket")

	hubId := c.Query("room")
	if len(hubId) == 0 {
		zap.S().Errorf("failed to get room")
		return
	}

	hub, exists := global.HubManager.GetHubById(hubId)
	if !exists {
		zap.S().Errorf("failed to find hub by id: %s", hubId)
	}

	ws.ServeWs(conn, hub)
}

func GetRoomPage(c *gin.Context) {
	c.HTML(http.StatusOK, "room_selection.html", gin.H{})
}

func GetGamePage(c *gin.Context) {
	room := c.Query("room")
	userId := uuid.New().String()

	c.HTML(http.StatusOK, "game.html", gin.H{
		"Room":   room,
		"UserID": userId,
	})
}
