package initinal

import (
	"go.uber.org/zap"
)

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Panic(err.Error())
	}
	zap.ReplaceGlobals(logger)
}
