package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kisielk/sqlstruct"
	"github.com/nelsonleduc/calmanbot/config"
	"github.com/nelsonleduc/calmanbot/service"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const state = "state"

var auth spotify.Authenticator
var client spotify.Client

type spotifyPlaylist struct {
	GroupID    string `sql:"group_id"`
	PlaylistID string `sql:"playlist_id"`
}

func playlistForGroup(groupID, groupName string, create bool) *spotifyPlaylist {
	queryStr := fmt.Sprintf("SELECT %s FROM spotify_playlists LIMIT 1", sqlstruct.Columns(spotifyPlaylist{}))
	rows, err := config.DB().Query(queryStr)
	if err != nil {
		return nil
	}
	defer rows.Close()

	if config.Configuration().SuperVerboseMode() {
		fmt.Println("Loaded server list:")
	}
	var playlist spotifyPlaylist
	if rows.Next() {
		err = sqlstruct.Scan(&playlist, rows)
	} else if create {
		user, _ := client.CurrentUser()
		name := groupID
		if len(groupName) > 0 {
			name = groupName
		}
		spotifyPlaylistCreated, _ := client.CreatePlaylistForUser(user.ID, name+" (CalmanBot)", "", false)
		queryStr := "INSERT INTO spotify_playlists(group_id, playlist_id) VALUES($1, $2) ON CONFLICT DO NOTHING"
		config.DB().Exec(queryStr, groupID, spotifyPlaylistCreated.ID.String())
		playlist = spotifyPlaylist{groupID, spotifyPlaylistCreated.ID.String()}
	}

	return &playlist
}

func SetupSpotify() {
	var redirectURL string
	if config.Configuration().LocalSpotifyAuth() {
		redirectURL = "http://localhost:4000/spotifyRedirect"
	} else {
		redirectURL = "https://calmanbot-production.herokuapp.com/spotifyRedirect"
	}

	oauthToken := os.Getenv("spotify_oauth")
	auth = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPrivate)

	url := auth.AuthURL(state)
	fmt.Printf("spotify auth url: %+v\n", url)

	var token oauth2.Token
	b := []byte(oauthToken)
	json.Unmarshal(b, &token)
	client = auth.NewClient(&token)
	fmt.Printf("oauth spotify client: %+v\n", client)
}

func HandleSpotifyRedirect(w http.ResponseWriter, r *http.Request) {
	token, _ := auth.Token(state, r)
	client := auth.NewClient(token)

	tokenOutput, _ := client.Token()
	jsonVersion, _ := json.Marshal(tokenOutput)
	output := string(jsonVersion)
	output = strings.Replace(output, "\"", "\\\"", -1)
	fmt.Printf("oauth token: \"%+v\"\n", output)
}

func processSpotify(groupID string, spotifyID string, groupName string) {
	if !config.Configuration().EnableSpotify() {
		return
	}

	hasTrigger := service.TriggerExists("spotifyPlaylist", groupID)
	if !hasTrigger {
		return
	}

	playlist := playlistForGroup(groupID, groupName, true)
	client.AddTracksToPlaylist(spotify.ID(playlist.PlaylistID), spotify.ID(spotifyID))
}

func HandlePlaylistRequest(w http.ResponseWriter, r *http.Request) {
	if !config.Configuration().EnableSpotify() {
		outputData := map[string]string{
			"output": "This feature is not enabled!",
		}
		json, _ := json.Marshal(outputData)
		w.Write(json)
	}

	groupID := r.URL.Query().Get("groupid")
	if groupID == "" {
		return
	}

	groupName := r.URL.Query().Get("groupName")

	hasTrigger := service.TriggerExists("spotifyPlaylist", groupID)
	playlist := playlistForGroup(groupID, groupName, hasTrigger)

	var outputData map[string]string
	if len(playlist.PlaylistID) > 0 {
		spotifyLink := "https://open.spotify.com/playlist/" + playlist.PlaylistID

		outputData = map[string]string{
			"output": spotifyLink,
		}
	} else {
		outputData = map[string]string{
			"output": "There is no playlist for this server",
		}
	}
	json, _ := json.Marshal(outputData)
	w.Write(json)
}

func HandleYoutubeLinkt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("link")
	if query == "" {
		return
	}

	groupID := r.URL.Query().Get("groupid")
	groupName := r.URL.Query().Get("groupName")

	url := "https://api.song.link/v1-alpha.1/links?url=" + url.QueryEscape(query)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "CalmanBot/2.5.3")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Bad Response: %d length: %d", resp.StatusCode, len(content))
		return
	}

	var stuff map[string]interface{}

	err = json.Unmarshal(content, &stuff)
	if err != nil {
		return
	}

	pageURL := stuff["pageUrl"].(string)
	links := stuff["linksByPlatform"].(map[string]interface{})
	if len(links) <= 2 {
		return
	}

	spotifyLinkPayload := links["spotify"].(map[string]interface{})
	hasSpotify := spotifyLinkPayload != nil
	hasAppleMusic := links["appleMusic"] != nil

	if !hasAppleMusic || !hasSpotify {
		return
	}

	spotifyEntityID := (spotifyLinkPayload["entityUniqueId"].(string))[14:]

	if len(groupID) > 0 {
		go processSpotify(groupID, spotifyEntityID, groupName)
	}

	outputData := map[string]string{
		"pageUrl": pageURL,
	}

	json, _ := json.Marshal(outputData)
	w.Write(json)
}
