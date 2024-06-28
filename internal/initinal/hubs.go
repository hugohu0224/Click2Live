package initinal

import (
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
	h1 := ws.NewHub("a")
	h2 := ws.NewHub("b")
	hm.AddNewHub(h1)
	hm.AddNewHub(h2)

	global.HubManager = hm
	global.HubManager.StartAllHubs()
}
