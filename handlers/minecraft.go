package handlers

import (
	"encoding/json"
	"errors"
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
			var trackedState *bool
			tickInterval := monitorIntervalSeconds()
			fmt.Printf("[MC: %v] Monitoring Minecraft server status for %s\n", tickInterval, address)
			stateStack := NewStateStack(10)
			ticker := time.NewTicker(tickInterval)
			for ; true; <-ticker.C {
				status, currentState := minecraftState(address, tickInterval)
				stateStack.PushState(currentState)
				if config.Configuration().SuperVerboseMode() {
					fmt.Printf("[MC: %v] Minecraft server status for %s: %v %v  %+v\n", tickInterval, address, status.Version, status.Online, stateStack)
				}

				if trackedState == nil {
					trackedState = &currentState
				} else if currentState != *trackedState {
					if stateStack.LastNStatesMatch(2, currentState) {
						trackedState = &currentState
						statusText := *identifierString + " is now offline!"
						if currentState {
							statusText = *identifierString + " is now online!"
						}
						fmt.Printf("[MC: %v] Changed status for Minecraft server status for %s: %v %v %+v\n", tickInterval, address, status.Version, status.Online, stateStack)
						post := service.Post{Key: "", Text: statusText, RawText: statusText, Type: service.PostTypeText, CacheID: 0}
						service.FanoutTrigger(address, post)
					}
				}
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

type stateStack struct {
	capacity int
	history  []bool
}

func NewStateStack(capacity int) stateStack {
	return stateStack{capacity: capacity, history: []bool{}}
}

func (s *stateStack) PushState(state bool) {
	s.history = append(s.history, state)
	if len(s.history) > s.capacity {
		s.history = s.history[1:]
	}
}

func (s stateStack) Len() int {
	return len(s.history)
}

func (s stateStack) Capacity() int {
	return s.capacity
}

func (s stateStack) LastState() (bool, error) {
	if len(s.history) == 0 {
		return false, errors.New("stack is empty")
	}
	return s.history[len(s.history)-1], nil
}

func (s stateStack) LastNStatesMatch(n int, state bool) bool {
	if len(s.history) < n {
		return false
	}
	for i := len(s.history) - 1; i >= len(s.history)-n; i-- {
		if s.history[i] != state {
			return false
		}
	}
	return true
}
