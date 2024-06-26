package main

import (
	"go.uber.org/zap"
	"tap2live/internal/handlers/ws"
	"tap2live/internal/initinal"
)

func main() {
	// logger
	initinal.InitLogger()
	zap.S().Info("logger initialized")

	// routers
	Router := initinal.InitRouters()
	zap.S().Info("router initialized")

	// broadcast
	go ws.HandleMessages()

	err := Router.Run(":8080")
	if err != nil {
		zap.S().Panicf("fail to start web server")
	}
}
