package api

import (
	"net/http"
	"net/url"
	"strings"
)

// Endpoint describes an HTTP method and server route. Route parameters use
// :name syntax; clients use the corresponding Path helpers to substitute IDs.
// Routes are relative to the server origin and do not include a base URL.
type Endpoint struct {
	Method string
	Route  string
}

// Endpoint values are shared definitions and should be treated as read-only.
var (
	ListMatches       = Endpoint{Method: http.MethodGet, Route: "/matches"}
	ListLobbies       = Endpoint{Method: http.MethodGet, Route: "/lobbies"}
	ListPublicLobbies = Endpoint{Method: http.MethodGet, Route: "/lobbies/public"}
	JoinLobby         = Endpoint{Method: http.MethodPost, Route: "/lobby/:id/join"}

	Login = Endpoint{Method: http.MethodPost, Route: "/login"}

	CreateUser     = Endpoint{Method: http.MethodPost, Route: "/user"}
	GetUserByEmail = Endpoint{Method: http.MethodPost, Route: "/user/getbyemail"}

	CreateLobby         = Endpoint{Method: http.MethodPost, Route: "/lobby"}
	DeleteLobby         = Endpoint{Method: http.MethodDelete, Route: "/lobby/:id"}
	GetLobby            = Endpoint{Method: http.MethodGet, Route: "/lobby/:id"}
	UpdateLobby         = Endpoint{Method: http.MethodPut, Route: "/lobby/:id"}
	AddUserToLobby      = Endpoint{Method: http.MethodPost, Route: "/lobby/:id/user"}
	LeaveLobby          = Endpoint{Method: http.MethodDelete, Route: "/lobby/:id/leave"}
	RemoveUserFromLobby = Endpoint{Method: http.MethodDelete, Route: "/lobby/:id/user/:user_id"}
	SetReady            = Endpoint{Method: http.MethodPut, Route: "/lobby/:id/ready"}
	LaunchMatch         = Endpoint{Method: http.MethodPost, Route: "/lobby/:id/launch"}

	GetMatch           = Endpoint{Method: http.MethodGet, Route: "/match/:id"}
	SubmitMatchCommand = Endpoint{Method: http.MethodPost, Route: "/match/:id/commands"}
	GetCatalogue       = Endpoint{Method: http.MethodGet, Route: "/match/:id/catalogue"}
	EndMatch           = Endpoint{Method: http.MethodPost, Route: "/admin/match/:id/end"}
)

func LeaveLobbyPath(id string) string { return idPath(LeaveLobby.Route, id) }

func JoinLobbyPath(id string) string { return idPath(JoinLobby.Route, id) }

func DeleteLobbyPath(id string) string    { return idPath(DeleteLobby.Route, id) }
func GetLobbyPath(id string) string       { return idPath(GetLobby.Route, id) }
func UpdateLobbyPath(id string) string    { return idPath(UpdateLobby.Route, id) }
func AddUserToLobbyPath(id string) string { return idPath(AddUserToLobby.Route, id) }

func RemoveUserFromLobbyPath(lobbyID, userID string) string {
	return strings.NewReplacer(
		":id", url.PathEscape(lobbyID),
		":user_id", url.PathEscape(userID),
	).Replace(RemoveUserFromLobby.Route)
}

func SetReadyPath(id string) string           { return idPath(SetReady.Route, id) }
func LaunchMatchPath(id string) string        { return idPath(LaunchMatch.Route, id) }
func GetMatchPath(id string) string           { return idPath(GetMatch.Route, id) }
func SubmitMatchCommandPath(id string) string { return idPath(SubmitMatchCommand.Route, id) }
func GetCataloguePath(id string) string       { return idPath(GetCatalogue.Route, id) }
func EndMatchPath(id string) string           { return idPath(EndMatch.Route, id) }

func idPath(route, id string) string {
	return strings.Replace(route, ":id", url.PathEscape(id), 1)
}
