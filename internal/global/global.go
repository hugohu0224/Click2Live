package global

import (
	"github.com/spf13/viper"
	"tap2live/internal/ws"
)

var (
	HubManager *ws.HubManager
	Dv         *viper.Viper // default viper
	Gv         *viper.Viper // google client viper
)
