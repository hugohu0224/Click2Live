package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
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
	ws.HandleConnection(conn)
}

func GetGamePage(c *gin.Context) {
	c.HTML(http.StatusOK, "game.html", gin.H{
		"UserID": uuid.New().String(),
	})
}