package handlers

import (
	"net/http"

	"github.com/nelsonleduc/calmanbot/service/groupme"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"github.com/nelsonleduc/calmanbot/service"
)

var groupmeService groupme.GMService

func init() {
	groupmeService = groupme.GMService{}
}

func BotHook(calman func(service.Message, service.Service, cache.QueryCache, models.Repo)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		message := groupmeService.MessageFromJSON(r.Body)
		monitor, _ := groupmeService.ServiceMonitor()
		cache := cache.NewSmartCache(monitor)

		calman(message, groupmeService, cache, models.PostGresRepo())
	}
}
