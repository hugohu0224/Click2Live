package initinal

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"tap2live/internal/global"
)

func InitConfigByViper() {
	// default viper
	global.Dv = viper.New()
	global.Dv.SetConfigName("config")
	global.Dv.SetConfigType("yaml")
	global.Dv.AddConfigPath(".")
	global.Dv.AddConfigPath("../")
	global.Dv.AddConfigPath("../../")

	// google client viper
	global.Gv = viper.New()
	global.Gv.SetConfigName("client_secret")
	global.Gv.SetConfigType("json")
	global.Gv.AddConfigPath(".")
	global.Gv.AddConfigPath("../")
	global.Gv.AddConfigPath("../../")

	// Read default config
	if err := global.Dv.ReadInConfig(); err != nil {
		zap.S().Fatalf("error reading default config file: %v", err)
	}
	zap.S().Infof("default config file used: %s", global.Dv.ConfigFileUsed())

	// Read Google client config
	if err := global.Gv.ReadInConfig(); err != nil {
		zap.S().Fatalf("error reading Google client config file: %v", err)
	}
	zap.S().Infof("google client config file used: %s", global.Gv.ConfigFileUsed())
}
