package handlers

import (
	"net/http"

	"github.com/nelsonleduc/calmanbot/service"
)

var groupmeService service.Service

func init() {
	groupmeService = *service.NewService("groupme")
}

func HandleBotHook(w http.ResponseWriter, r *http.Request) {

	message := groupmeService.MessageFromJSON(r.Body)
	monitor, _ := groupmeService.ServiceMonitor()
	cache := NewSmartCache(monitor)

	HandleCalman(message, groupmeService, cache)
}
