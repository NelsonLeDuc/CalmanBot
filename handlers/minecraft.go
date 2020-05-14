package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/whatupdave/mcping"
)

func MonitorMinecraft(address string, minuteCadence int) {
	for ; true; <-time.Tick(time.Duration(minuteCadence) * time.Minute) {
		status, err := mcping.Ping(address)
		fmt.Printf("[MC] Minecraft server status for %s: %v %v\n", address, status.Version, err)
	}
}

func HandleMinecraft(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("addr")
	if query == "" {
		return
	}

	status, err := mcping.Ping(query)
	if err == nil {
		output := map[string]string{
			"description": "Server is Online with " + strconv.Itoa(status.Online) + "/" + strconv.Itoa(status.Max),
		}
		json, _ := json.Marshal(output)
		w.Write(json)
	} else {
		output := map[string]string{
			"description": "Server is Offline",
		}
		json, _ := json.Marshal(output)
		w.Write(json)
	}
}
