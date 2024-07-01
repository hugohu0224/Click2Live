package main

import (
	"go.uber.org/zap"
	"tap2live/internal/initinal"
)

func main() {
	// logger
	initinal.InitLogger()
	zap.S().Info("logger initialized")

	// viper
	initinal.InitConfigByViper()
	zap.S().Info("config initialized")

	// routers
	Router := initinal.InitRouters()
	zap.S().Info("router initialized")

	// manager
	initinal.InitHubManager()
	zap.S().Info("HubManager initialized")

	err := Router.Run(":8080")
	if err != nil {
		zap.S().Panicf("fail to start web server")
	}
}
