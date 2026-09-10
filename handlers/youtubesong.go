package handlers

import (
	"encoding/json"
	"fmt"
	"io"
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

func playlistForGroup(serverID, groupName string, create bool) *spotifyPlaylist {
	log.Printf(
		"[spotify] playlistForGroup: serverID=%q groupName=%q create=%v",
		serverID,
		groupName,
		create,
	)

	queryStr := fmt.Sprintf(
		"SELECT %s FROM spotify_playlists WHERE group_id = $1 LIMIT 1",
		sqlstruct.Columns(spotifyPlaylist{}),
	)

	rows, err := config.DB().Query(queryStr, serverID)
	if err != nil {
		log.Printf(
			"[spotify] playlistForGroup: database query failed serverID=%q: %v",
			serverID,
			err,
		)
		return nil
	}
	defer rows.Close()

	var playlist spotifyPlaylist

	if rows.Next() {
		if err := sqlstruct.Scan(&playlist, rows); err != nil {
			log.Printf(
				"[spotify] playlistForGroup: failed to scan playlist serverID=%q: %v",
				serverID,
				err,
			)
			return nil
		}

		log.Printf(
			"[spotify] playlistForGroup: existing playlist found serverID=%q playlistID=%q",
			serverID,
			playlist.PlaylistID,
		)

		return &playlist
	}

	if err := rows.Err(); err != nil {
		log.Printf(
			"[spotify] playlistForGroup: row iteration error serverID=%q: %v",
			serverID,
			err,
		)
		return nil
	}

	if !create {
		log.Printf(
			"[spotify] playlistForGroup: no playlist found and create=false serverID=%q",
			serverID,
		)
		return &playlist
	}

	log.Printf(
		"[spotify] playlistForGroup: no playlist found; attempting to create one serverID=%q",
		serverID,
	)

	user, err := client.CurrentUser()
	if err != nil {
		log.Printf(
			"[spotify] playlistForGroup: CurrentUser failed serverID=%q: %v",
			serverID,
			err,
		)
		return nil
	}

	name := serverID
	if len(groupName) > 0 {
		name = groupName
	}

	playlistName := name + " (CalmanBot)"

	log.Printf(
		"[spotify] playlistForGroup: creating Spotify playlist name=%q userID=%q",
		playlistName,
		user.ID,
	)

	spotifyPlaylistCreated, err := client.CreatePlaylistForUser(
		user.ID,
		playlistName,
		"",
		false,
	)
	if err != nil {
		log.Printf(
			"[spotify] playlistForGroup: CreatePlaylistForUser failed serverID=%q name=%q: %v",
			serverID,
			playlistName,
			err,
		)
		return nil
	}

	if spotifyPlaylistCreated == nil {
		log.Printf(
			"[spotify] playlistForGroup: CreatePlaylistForUser returned nil playlist serverID=%q",
			serverID,
		)
		return nil
	}

	playlistID := spotifyPlaylistCreated.ID.String()

	log.Printf(
		"[spotify] playlistForGroup: Spotify playlist created serverID=%q playlistID=%q",
		serverID,
		playlistID,
	)

	queryStr = `
		INSERT INTO spotify_playlists(group_id, playlist_id)
		VALUES($1, $2)
		ON CONFLICT DO NOTHING
	`

	if _, err := config.DB().Exec(queryStr, serverID, playlistID); err != nil {
		log.Printf(
			"[spotify] playlistForGroup: failed to save playlist to database serverID=%q playlistID=%q: %v",
			serverID,
			playlistID,
			err,
		)
		return nil
	}

	log.Printf(
		"[spotify] playlistForGroup: playlist saved to database serverID=%q playlistID=%q",
		serverID,
		playlistID,
	)

	playlist = spotifyPlaylist{
		GroupID:    serverID,
		PlaylistID: playlistID,
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

	log.Printf(
		"[spotify] SetupSpotify: starting redirectURL=%q oauthTokenPresent=%v",
		redirectURL,
		oauthToken != "",
	)

	auth = spotify.NewAuthenticator(
		redirectURL,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPrivate,
	)

	authURL := auth.AuthURL(state)
	fmt.Printf("spotify auth url: %+v\n", authURL)

	var token oauth2.Token
	b := []byte(oauthToken)

	if err := json.Unmarshal(b, &token); err != nil {
		log.Printf("[spotify] SetupSpotify: failed to decode spotify_oauth: %v", err)
	}

	client = auth.NewClient(&token)

	log.Printf(
		"[spotify] SetupSpotify: Spotify client initialized tokenExpiry=%v",
		token.Expiry,
	)
}

func HandleSpotifyRedirect(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(state, r)
	if err != nil {
		log.Printf("[spotify] HandleSpotifyRedirect: failed to get token: %v", err)
		http.Error(w, "Spotify authorization failed", http.StatusInternalServerError)
		return
	}

	client := auth.NewClient(token)

	tokenOutput, err := client.Token()
	if err != nil {
		log.Printf("[spotify] HandleSpotifyRedirect: failed to get client token: %v", err)
		return
	}

	jsonVersion, err := json.Marshal(tokenOutput)
	if err != nil {
		log.Printf("[spotify] HandleSpotifyRedirect: failed to marshal token: %v", err)
		return
	}

	output := string(jsonVersion)
	output = strings.Replace(output, "\"", "\\\"", -1)

	fmt.Printf("oauth token: \"%+v\"\n", output)
}

func processSpotify(groupID, serverID, spotifyID string, groupName string) {
	log.Printf(
		"[spotify] processSpotify: START groupID=%q serverID=%q spotifyID=%q groupName=%q",
		groupID,
		serverID,
		spotifyID,
		groupName,
	)

	if !config.Configuration().EnableSpotify() {
		log.Printf(
			"[spotify] processSpotify: SKIP Spotify disabled groupID=%q serverID=%q",
			groupID,
			serverID,
		)
		return
	}

	log.Printf(
		"[spotify] processSpotify: checking trigger groupID=%q serverID=%q",
		groupID,
		serverID,
	)

	hasTrigger := service.TriggerExists(
		"spotifyPlaylist",
		serverID,
		groupID,
	)

	if !hasTrigger {
		log.Printf(
			"[spotify] processSpotify: SKIP no spotifyPlaylist trigger groupID=%q serverID=%q",
			groupID,
			serverID,
		)
		return
	}

	log.Printf(
		"[spotify] processSpotify: trigger found groupID=%q serverID=%q",
		groupID,
		serverID,
	)

	playlist := playlistForGroup(serverID, groupName, true)

	if playlist == nil {
		log.Printf(
			"[spotify] processSpotify: ERROR playlistForGroup returned nil groupID=%q serverID=%q",
			groupID,
			serverID,
		)
		return
	}

	if playlist.PlaylistID == "" {
		log.Printf(
			"[spotify] processSpotify: ERROR playlist has empty PlaylistID groupID=%q serverID=%q",
			groupID,
			serverID,
		)
		return
	}

	if spotifyID == "" {
		log.Printf(
			"[spotify] processSpotify: ERROR Spotify track ID is empty groupID=%q serverID=%q playlistID=%q",
			groupID,
			serverID,
			playlist.PlaylistID,
		)
		return
	}

	log.Printf(
		"[spotify] processSpotify: adding track spotifyID=%q playlistID=%q",
		spotifyID,
		playlist.PlaylistID,
	)

	snapshotID, err := client.AddTracksToPlaylist(
		spotify.ID(playlist.PlaylistID),
		spotify.ID(spotifyID),
	)
	if err != nil {
		log.Printf(
			"[spotify] processSpotify: ERROR AddTracksToPlaylist failed spotifyID=%q playlistID=%q: %v",
			spotifyID,
			playlist.PlaylistID,
			err,
		)
		return
	}

	log.Printf(
		"[spotify] processSpotify: SUCCESS added spotifyID=%q playlistID=%q snapshotID=%q",
		spotifyID,
		playlist.PlaylistID,
		snapshotID,
	)
}

func HandlePlaylistRequest(w http.ResponseWriter, r *http.Request) {
	if !config.Configuration().EnableSpotify() {
		outputData := map[string]string{
			"output": "This feature is not enabled!",
		}
		json, _ := json.Marshal(outputData)
		w.Write(json)
		return
	}

	serverID := r.URL.Query().Get("serverid")
	if serverID == "" {
		return
	}

	groupName := r.URL.Query().Get("groupName")

	hasTrigger := service.TriggerExists("spotifyPlaylist", serverID, "")
	playlist := playlistForGroup(serverID, groupName, hasTrigger)

	var outputData map[string]string

	if playlist != nil && len(playlist.PlaylistID) > 0 {
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
	serverID := r.URL.Query().Get("serverid")
	groupName := r.URL.Query().Get("groupName")

	url := "https://api.ml.jadquir.com/v1/lookup?q=" + url.QueryEscape(query)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("musiclink request error: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("musiclink_api_key"))
	req.Header.Set("User-Agent", "CalmanBot/2.5.3")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("musiclink error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	content, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Bad Response: %d length: %d", resp.StatusCode, len(content))
		return
	}

	var stuff map[string]interface{}

	err = json.Unmarshal(content, &stuff)
	if err != nil {
		fmt.Printf("musiclink json error: %v\n", err)
		return
	}

	data := stuff["data"].([]interface{})
	if len(data) == 0 {
		fmt.Printf("musiclink no results\n")
		return
	}

	track := data[0].(map[string]interface{})
	pageURL := track["url"].(string)
	links := track["links"].(map[string]interface{})

	spotifyLink, hasSpotify := links["spotify"].(string)
	_, hasAppleMusic := links["apple_music"].(string)

	if !hasAppleMusic || !hasSpotify {
		fmt.Printf(
			"musiclink missing platform apple: %v spotify: %v\n",
			hasAppleMusic,
			hasSpotify,
		)
		return
	}

	spotifyEntityID := strings.TrimPrefix(
		spotifyLink,
		"https://open.spotify.com/track/",
	)

	if index := strings.Index(spotifyEntityID, "?"); index >= 0 {
		spotifyEntityID = spotifyEntityID[:index]
	}

	log.Printf(
		"[spotify] HandleYoutubeLinkt: resolved track groupID=%q serverID=%q spotifyLink=%q spotifyID=%q",
		groupID,
		serverID,
		spotifyLink,
		spotifyEntityID,
	)

	if len(groupID) > 0 && hasSpotify {
		log.Printf(
			"[spotify] HandleYoutubeLinkt: launching processSpotify groupID=%q serverID=%q spotifyID=%q",
			groupID,
			serverID,
			spotifyEntityID,
		)

		go processSpotify(groupID, serverID, spotifyEntityID, groupName)
	} else {
		log.Printf(
			"[spotify] HandleYoutubeLinkt: NOT launching processSpotify groupID=%q hasSpotify=%v",
			groupID,
			hasSpotify,
		)
	}

	outputData := map[string]string{
		"pageUrl": pageURL,
	}

	json, _ := json.Marshal(outputData)
	w.Write(json)
}
