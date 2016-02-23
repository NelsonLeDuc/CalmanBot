package handlers

import (
	"net/http"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/service"
	_ "github.com/nelsonleduc/calmanbot/service/groupme"
)

var groupmeService service.Service

func init() {
	groupmeService = *service.NewService("groupme")
}

func HandleBotHook(w http.ResponseWriter, r *http.Request) {

	message := groupmeService.MessageFromJSON(r.Body)
	monitor, _ := groupmeService.ServiceMonitor()
	cache := cache.NewSmartCache(monitor)

	HandleCalman(message, groupmeService, cache)
}
