package routers

import (
	"github.com/gin-gonic/gin"
	"tap2live/internal/api"
)

func InitGameRouter(router *gin.RouterGroup) {
	{
		Router := router.Group("/game")
		Router.Static("/static", "./internal/static")
		Router.GET("/ws", api.WsEndpoint)
		Router.GET("/room", api.GetRoomPage)
		Router.GET("/gamepage", api.GetGamePage)

	}
}
