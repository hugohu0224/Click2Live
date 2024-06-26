package routers

import (
	"github.com/gin-gonic/gin"
	"tap2live/internal/handlers/api"
)

func InitWsRouter(router *gin.RouterGroup) {
	{
		Router := router.Group("/ws")
		Router.GET("/game", api.WsEndpoint)
	}
}
