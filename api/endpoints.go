package api

import (
	"net/http"
	"net/url"
	"strings"
)

// Method and Route constants describe server registration. Route parameters use
// :name syntax; clients use the corresponding Path helpers to substitute IDs.
// Paths are relative to the server origin and do not include a base URL.
const (
	LoginMethod = http.MethodPost
	LoginRoute  = "/login"

	CreateUserMethod     = http.MethodPost
	CreateUserRoute      = "/user"
	GetUserByEmailMethod = http.MethodPost
	GetUserByEmailRoute  = "/user/getbyemail"

	CreateLobbyMethod         = http.MethodPost
	CreateLobbyRoute          = "/lobby"
	DeleteLobbyMethod         = http.MethodDelete
	DeleteLobbyRoute          = "/lobby/:id"
	GetLobbyMethod            = http.MethodGet
	GetLobbyRoute             = "/lobby/:id"
	UpdateLobbyMethod         = http.MethodPut
	UpdateLobbyRoute          = "/lobby/:id"
	AddUserToLobbyMethod      = http.MethodPost
	AddUserToLobbyRoute       = "/lobby/:id/user"
	RemoveUserFromLobbyMethod = http.MethodDelete
	RemoveUserFromLobbyRoute  = "/lobby/:id/user/:user_id"
	SetReadyMethod            = http.MethodPut
	SetReadyRoute             = "/lobby/:id/ready"
	LaunchMatchMethod         = http.MethodPost
	LaunchMatchRoute          = "/lobby/:id/launch"

	GetMatchMethod           = http.MethodGet
	GetMatchRoute            = "/match/:id"
	SubmitMatchCommandMethod = http.MethodPost
	SubmitMatchCommandRoute  = "/match/:id/commands"
	GetCatalogueMethod       = http.MethodGet
	GetCatalogueRoute        = "/match/:id/catalogue"
	EndMatchMethod           = http.MethodPost
	EndMatchRoute            = "/admin/match/:id/end"
)

func DeleteLobbyPath(id string) string    { return idPath(DeleteLobbyRoute, id) }
func GetLobbyPath(id string) string       { return idPath(GetLobbyRoute, id) }
func UpdateLobbyPath(id string) string    { return idPath(UpdateLobbyRoute, id) }
func AddUserToLobbyPath(id string) string { return idPath(AddUserToLobbyRoute, id) }

func RemoveUserFromLobbyPath(lobbyID, userID string) string {
	return strings.NewReplacer(
		":id", url.PathEscape(lobbyID),
		":user_id", url.PathEscape(userID),
	).Replace(RemoveUserFromLobbyRoute)
}

func SetReadyPath(id string) string           { return idPath(SetReadyRoute, id) }
func LaunchMatchPath(id string) string        { return idPath(LaunchMatchRoute, id) }
func GetMatchPath(id string) string           { return idPath(GetMatchRoute, id) }
func SubmitMatchCommandPath(id string) string { return idPath(SubmitMatchCommandRoute, id) }
func GetCataloguePath(id string) string       { return idPath(GetCatalogueRoute, id) }
func EndMatchPath(id string) string           { return idPath(EndMatchRoute, id) }

func idPath(route, id string) string {
	return strings.Replace(route, ":id", url.PathEscape(id), 1)
}
