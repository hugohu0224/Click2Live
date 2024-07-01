package initinal

import (
	"github.com/spf13/viper"
	"sync"
	"tap2live/internal/global"
	"tap2live/internal/ws"
)

func InitHubManager() {
	hm := &ws.HubManager{
		Hubs:     map[*ws.Hub]bool{},
		HubsById: map[string]*ws.Hub{},
		Mu:       sync.RWMutex{},
	}
	h1 := ws.NewHub(viper.GetString("FIRST_ROOM_ID"))
	h2 := ws.NewHub(viper.GetString("SECOND_ROOM_ID"))
	hm.AddNewHub(h1)
	hm.AddNewHub(h2)

	global.HubManager = hm
	global.HubManager.RunHubs()
}
