package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kisielk/sqlstruct"
	"github.com/nelsonleduc/calmanbot/config"
	"github.com/nelsonleduc/calmanbot/service"
	"github.com/whatupdave/mcping"
)

const defaultMonitorInterval = time.Duration(20) * time.Second

type minecraftServer struct {
	Address string  `sql:"address"`
	Name    *string `sql:"name"`
}

func monitorIntervalSeconds() time.Duration {
	if configValue := config.Configuration().MonitorIntervalSeconds(); configValue > 0 {
		return time.Duration(configValue) * time.Second
	}
	return defaultMonitorInterval
}

func storedAddresses() []minecraftServer {
	queryStr := fmt.Sprintf("SELECT %s FROM minecraft_servers", sqlstruct.Columns(minecraftServer{}))
	rows, err := config.DB().Query(queryStr)
	if err != nil {
		return []minecraftServer{}
	}
	defer rows.Close()

	if config.Configuration().SuperVerboseMode() {
		fmt.Println("Loaded server list:")
	}
	servers := []minecraftServer{}
	for rows.Next() {
		var server minecraftServer
		err := sqlstruct.Scan(&server, rows)
		if err != nil {
			continue
		}
		if config.Configuration().SuperVerboseMode() {
			fmt.Printf("   %+v\n", server)
		}
		servers = append(servers, server)
	}
	return servers
}

var trackedServers map[string]bool

func init() {
	trackedServers = make(map[string]bool)
}

func minecraftState(address string, tickInterval time.Duration) (mcping.PingResponse, bool) {
	status, err := mcping.Ping(address)
	if config.Configuration().SuperVerboseMode() {
		fmt.Printf("[MC: %v] Minecraft server status for %s: %v %v err: %v\n", tickInterval, address, status.Version, status.Online, err)
	}
	return status, err == nil
}

func MonitorMinecraft() {
	for _, server := range storedAddresses() {
		if trackedServers[server.Address] {
			continue
		}
		trackedServers[server.Address] = true
		address := server.Address
		name := server.Name
		identifierString := name
		addressStr := address
		if strings.Contains(address, ":25565") {
			addressStr = strings.ReplaceAll(address, ":25565", "")
		}
		if name == nil || len(*name) == 0 {
			identifierString = &addressStr
		}
		go func() {
			var prevState *bool
			tickInterval := monitorIntervalSeconds()
			fmt.Printf("[MC: %v] Monitoring Minecraft server status for %s\n", tickInterval, address)
			for ; true; <-time.Tick(tickInterval) {
				status, currentState := minecraftState(address, tickInterval)
				if config.Configuration().SuperVerboseMode() {
					fmt.Printf("[MC: %v] Minecraft server status for %s: %v %v\n", tickInterval, address, status.Version, status.Online)
				}

				if prevState != nil && *prevState != currentState {
					status, currentState = minecraftState(address, tickInterval)
					if prevState != nil && *prevState != currentState {
						statusText := *identifierString + " is now offline!"
						if currentState {
							statusText = *identifierString + " is now online!"
						}
						fmt.Printf("[MC: %v] Changed status for Minecraft server status for %s: %v %v\n", tickInterval, address, status.Version, status.Online)
						post := service.Post{"", statusText, statusText, service.PostTypeText, 0}
						service.FanoutTrigger(address, post)
					}
				}
				prevState = &currentState
			}
		}()
	}
}

func HandleMinecraft(w http.ResponseWriter, r *http.Request) {
	if !config.Configuration().EnableMinecraft() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	query := r.URL.Query().Get("addr")
	if query == "" {
		return
	}

	split := strings.Split(query, ":")
	address := split[0]
	status, err := mcping.Ping(query)
	if err == nil {
		output := map[string]string{
			"description": address + " is Online with " + strconv.Itoa(status.Online) + "/" + strconv.Itoa(status.Max),
		}
		json, _ := json.Marshal(output)
		w.Write(json)
	} else {
		output := map[string]string{
			"description": address + " is Offline",
		}
		json, _ := json.Marshal(output)
		w.Write(json)
	}
}

func HandleTrackMinecraft(w http.ResponseWriter, r *http.Request) {
	if !config.Configuration().EnableMinecraft() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	query := r.URL.Query().Get("addr")
	name := r.URL.Query().Get("name")
	if query == "" {
		return
	}

	if len(name) > 0 {
		queryStr := "INSERT INTO minecraft_servers(address, name) VALUES($1, $2) ON CONFLICT (address) DO UPDATE SET name = $2"
		config.DB().Exec(queryStr, query, name)
	} else {
		queryStr := "INSERT INTO minecraft_servers(address, name) VALUES($1, NULL) ON CONFLICT DO NOTHING"
		config.DB().Exec(queryStr, query)
	}
	go MonitorMinecraft()
}
