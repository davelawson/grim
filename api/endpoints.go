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
	AddBotToLobby       = Endpoint{Method: http.MethodPost, Route: "/lobby/:id/bot"}
	RemoveBotFromLobby  = Endpoint{Method: http.MethodDelete, Route: "/lobby/:id/bot/:botid"}
	LeaveLobby          = Endpoint{Method: http.MethodDelete, Route: "/lobby/:id/leave"}
	RemoveUserFromLobby = Endpoint{Method: http.MethodDelete, Route: "/lobby/:id/user/:user_id"}
	SetReady            = Endpoint{Method: http.MethodPut, Route: "/lobby/:id/ready"}
	LaunchMatch         = Endpoint{Method: http.MethodPost, Route: "/lobby/:id/launch"}

	GetMatch     = Endpoint{Method: http.MethodGet, Route: "/match/:id"}
	GetDraft     = Endpoint{Method: http.MethodGet, Route: "/match/:id/draft"}
	SubmitDraft  = Endpoint{Method: http.MethodPost, Route: "/match/:id/draft"}
	GetCatalogue = Endpoint{Method: http.MethodGet, Route: "/match/:id/catalogue"}
	EndMatch     = Endpoint{Method: http.MethodPost, Route: "/admin/match/:id/end"}
)

func LeaveLobbyPath(id string) string { return idPath(LeaveLobby.Route, id) }

func JoinLobbyPath(id string) string { return idPath(JoinLobby.Route, id) }

func DeleteLobbyPath(id string) string    { return idPath(DeleteLobby.Route, id) }
func GetLobbyPath(id string) string       { return idPath(GetLobby.Route, id) }
func UpdateLobbyPath(id string) string    { return idPath(UpdateLobby.Route, id) }
func AddUserToLobbyPath(id string) string { return idPath(AddUserToLobby.Route, id) }

func AddBotToLobbyPath(id string) string { return idPath(AddBotToLobby.Route, id) }
func RemoveBotFromLobbyPath(lobbyID, botID string) string {
	return strings.NewReplacer(":id", url.PathEscape(lobbyID), ":botid", url.PathEscape(botID)).Replace(RemoveBotFromLobby.Route)
}

func RemoveUserFromLobbyPath(lobbyID, userID string) string {
	return strings.NewReplacer(
		":id", url.PathEscape(lobbyID),
		":user_id", url.PathEscape(userID),
	).Replace(RemoveUserFromLobby.Route)
}

func SetReadyPath(id string) string     { return idPath(SetReady.Route, id) }
func LaunchMatchPath(id string) string  { return idPath(LaunchMatch.Route, id) }
func GetMatchPath(id string) string     { return idPath(GetMatch.Route, id) }
func GetDraftPath(id string) string     { return idPath(GetDraft.Route, id) }
func SubmitDraftPath(id string) string  { return idPath(SubmitDraft.Route, id) }
func GetCataloguePath(id string) string { return idPath(GetCatalogue.Route, id) }
func EndMatchPath(id string) string     { return idPath(EndMatch.Route, id) }

func idPath(route, id string) string {
	return strings.Replace(route, ":id", url.PathEscape(id), 1)
}

// Gameplay endpoints accept narrow inputs and share the match response.
var (
	PlayCard         = Endpoint{Method: http.MethodPost, Route: "/match/:id/actions/play-card"}
	Attack           = Endpoint{Method: http.MethodPost, Route: "/match/:id/actions/attack"}
	AttemptChallenge = Endpoint{Method: http.MethodPost, Route: "/match/:id/actions/attempt-challenge"}
	ActivateAbility  = Endpoint{Method: http.MethodPost, Route: "/match/:id/actions/activate-ability"}
	ResearchElement  = Endpoint{Method: http.MethodPost, Route: "/match/:id/actions/research-element"}
	ResearchSpell    = Endpoint{Method: http.MethodPost, Route: "/match/:id/actions/research-spell"}
	HealWithFocus    = Endpoint{Method: http.MethodPost, Route: "/match/:id/actions/heal-with-focus"}
	EndTurn          = Endpoint{Method: http.MethodPost, Route: "/match/:id/actions/end-turn"}
	Maintenance      = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/maintenance"}
	Kneel            = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/kneel"}
	ResolveEffect    = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/resolve-effect"}
	SkipEffect       = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/skip-effect"}
	UndertakeHunt    = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/undertake-hunt"}
	ChooseResearch   = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/research"}
	EndDiscard       = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/end-discard"}
	HandDiscard      = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/hand-discard"}
	AssetLoss        = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/asset-loss"}
	Scry             = Endpoint{Method: http.MethodPost, Route: "/match/:id/choices/:choiceid/scry"}
	React            = Endpoint{Method: http.MethodPost, Route: "/match/:id/reactions/play-card"}
	Defend           = Endpoint{Method: http.MethodPost, Route: "/match/:id/reactions/defend"}
	Pass             = Endpoint{Method: http.MethodPost, Route: "/match/:id/reactions/pass"}
)

func PlayCardPath(id string) string         { return idPath(PlayCard.Route, id) }
func AttackPath(id string) string           { return idPath(Attack.Route, id) }
func AttemptChallengePath(id string) string { return idPath(AttemptChallenge.Route, id) }
func ActivateAbilityPath(id string) string  { return idPath(ActivateAbility.Route, id) }
func ResearchElementPath(id string) string  { return idPath(ResearchElement.Route, id) }
func ResearchSpellPath(id string) string    { return idPath(ResearchSpell.Route, id) }
func HealWithFocusPath(id string) string    { return idPath(HealWithFocus.Route, id) }
func EndTurnPath(id string) string          { return idPath(EndTurn.Route, id) }
func MaintenancePath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(Maintenance.Route)
}
func KneelPath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(Kneel.Route)
}
func ResolveEffectPath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(ResolveEffect.Route)
}
func SkipEffectPath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(SkipEffect.Route)
}
func UndertakeHuntPath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(UndertakeHunt.Route)
}
func ChooseResearchPath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(ChooseResearch.Route)
}
func EndDiscardPath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(EndDiscard.Route)
}
func HandDiscardPath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(HandDiscard.Route)
}
func AssetLossPath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(AssetLoss.Route)
}
func ScryPath(id, choiceID string) string {
	return strings.NewReplacer(":id", url.PathEscape(id), ":choiceid", url.PathEscape(choiceID)).Replace(Scry.Route)
}
func ReactPath(id string) string  { return idPath(React.Route, id) }
func DefendPath(id string) string { return idPath(Defend.Route, id) }
func PassPath(id string) string   { return idPath(Pass.Route, id) }
