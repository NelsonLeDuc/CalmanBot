package main

import (
	"net/http"

	"github.com/nelsonleduc/calmanbot/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"CalmanRespond",
		"POST",
		"/bots/botHook",
		handlers.BotHook(handlers.HandleCalman),
	},
	Route{
		"GoogleImage",
		"GET",
		"/animated",
		handlers.HandleGoogleImage,
	},
	Route{
		"JeffSum",
		"GET",
		"/jeff",
		handlers.HandleJeffFetch,
	},
	Route{
		"Minecraft",
		"GET",
		"/minecraftStatus",
		handlers.HandleMinecraft,
	},
	Route{
		"MinecraftAdd",
		"GET",
		"/minecraftTrack",
		handlers.HandleTrackMinecraft,
	},
	Route{
		"YoutubeSong",
		"GET",
		"/youtubeSong",
		handlers.HandleYoutubeLinkt,
	},
	Route{
		"SpotifyAuth",
		"GET",
		"/spotifyRedirect",
		handlers.HandleSpotifyRedirect,
	},
	Route{
		"Playlist",
		"GET",
		"/playlistGet",
		handlers.HandlePlaylistRequest,
	},
}
