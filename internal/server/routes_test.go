package server

import (
	"github.com/davelawson/grim/api"
	"net/http"
	"strings"
	"testing"
)

func TestSharedEndpointsMatchProductionRoutes(t *testing.T) {
	f := newLifecycleFixture(t)
	lobbyID, userID := f.ids["owner"], f.ids["member"]
	tests := []struct {
		name, method, route, path string
		wantStatus                int
	}{
		{"public lobbies", api.ListPublicLobbies.Method, api.ListPublicLobbies.Route, api.ListPublicLobbies.Route, 401},
		{"leave lobby", api.LeaveLobby.Method, api.LeaveLobby.Route, api.LeaveLobbyPath(lobbyID), 401},
		{"join lobby", api.JoinLobby.Method, api.JoinLobby.Route, api.JoinLobbyPath(lobbyID), 401},
		{"list matches", api.ListMatches.Method, api.ListMatches.Route, api.ListMatches.Route, 401},
		{"list lobbies", api.ListLobbies.Method, api.ListLobbies.Route, api.ListLobbies.Route, 401},
		{"login", api.Login.Method, api.Login.Route, api.Login.Route, 400},
		{"create user", api.CreateUser.Method, api.CreateUser.Route, api.CreateUser.Route, 400},
		{"get user", api.GetUserByEmail.Method, api.GetUserByEmail.Route, api.GetUserByEmail.Route, 401},
		{"create lobby", api.CreateLobby.Method, api.CreateLobby.Route, api.CreateLobby.Route, 401},
		{"delete lobby", api.DeleteLobby.Method, api.DeleteLobby.Route, api.DeleteLobbyPath(lobbyID), 401},
		{"get lobby", api.GetLobby.Method, api.GetLobby.Route, api.GetLobbyPath(lobbyID), 401},
		{"update lobby", api.UpdateLobby.Method, api.UpdateLobby.Route, api.UpdateLobbyPath(lobbyID), 401},
		{"add user", api.AddUserToLobby.Method, api.AddUserToLobby.Route, api.AddUserToLobbyPath(lobbyID), 401},
		{"remove user", api.RemoveUserFromLobby.Method, api.RemoveUserFromLobby.Route, api.RemoveUserFromLobbyPath(lobbyID, userID), 401},
		{"ready", api.SetReady.Method, api.SetReady.Route, api.SetReadyPath(lobbyID), 401},
		{"launch", api.LaunchMatch.Method, api.LaunchMatch.Route, api.LaunchMatchPath(lobbyID), 401},
		{"PlayCard", api.PlayCard.Method, api.PlayCard.Route, api.PlayCardPath(lobbyID), 401},
		{"Attack", api.Attack.Method, api.Attack.Route, api.AttackPath(lobbyID), 401},
		{"AttemptChallenge", api.AttemptChallenge.Method, api.AttemptChallenge.Route, api.AttemptChallengePath(lobbyID), 401},
		{"ActivateAbility", api.ActivateAbility.Method, api.ActivateAbility.Route, api.ActivateAbilityPath(lobbyID), 401},
		{"ResearchElement", api.ResearchElement.Method, api.ResearchElement.Route, api.ResearchElementPath(lobbyID), 401},
		{"ResearchSpell", api.ResearchSpell.Method, api.ResearchSpell.Route, api.ResearchSpellPath(lobbyID), 401},
		{"HealWithFocus", api.HealWithFocus.Method, api.HealWithFocus.Route, api.HealWithFocusPath(lobbyID), 401},
		{"EndTurn", api.EndTurn.Method, api.EndTurn.Route, api.EndTurnPath(lobbyID), 401},
		{"Maintenance", api.Maintenance.Method, api.Maintenance.Route, api.MaintenancePath(lobbyID, "choice"), 401},
		{"Kneel", api.Kneel.Method, api.Kneel.Route, api.KneelPath(lobbyID, "choice"), 401},
		{"ResolveEffect", api.ResolveEffect.Method, api.ResolveEffect.Route, api.ResolveEffectPath(lobbyID, "choice"), 401},
		{"SkipEffect", api.SkipEffect.Method, api.SkipEffect.Route, api.SkipEffectPath(lobbyID, "choice"), 401},
		{"UndertakeHunt", api.UndertakeHunt.Method, api.UndertakeHunt.Route, api.UndertakeHuntPath(lobbyID, "choice"), 401},
		{"ChooseResearch", api.ChooseResearch.Method, api.ChooseResearch.Route, api.ChooseResearchPath(lobbyID, "choice"), 401},
		{"EndDiscard", api.EndDiscard.Method, api.EndDiscard.Route, api.EndDiscardPath(lobbyID, "choice"), 401},
		{"HandDiscard", api.HandDiscard.Method, api.HandDiscard.Route, api.HandDiscardPath(lobbyID, "choice"), 401},
		{"AssetLoss", api.AssetLoss.Method, api.AssetLoss.Route, api.AssetLossPath(lobbyID, "choice"), 401},
		{"Scry", api.Scry.Method, api.Scry.Route, api.ScryPath(lobbyID, "choice"), 401},
		{"React", api.React.Method, api.React.Route, api.ReactPath(lobbyID), 401},
		{"Defend", api.Defend.Method, api.Defend.Route, api.DefendPath(lobbyID), 401},
		{"Pass", api.Pass.Method, api.Pass.Route, api.PassPath(lobbyID), 401},
		{"get match", api.GetMatch.Method, api.GetMatch.Route, api.GetMatchPath(lobbyID), 401},
		{"get draft", api.GetDraft.Method, api.GetDraft.Route, api.GetDraftPath(lobbyID), 401},
		{"submit draft", api.SubmitDraft.Method, api.SubmitDraft.Route, api.SubmitDraftPath(lobbyID), 401},
		{"catalogue", api.GetCatalogue.Method, api.GetCatalogue.Route, api.GetCataloguePath(lobbyID), 401},
		{"end match", api.EndMatch.Method, api.EndMatch.Route, api.EndMatchPath(lobbyID), 401},
	}
	routes := f.router.Routes()
	if len(routes) != len(tests) {
		t.Fatalf("registered %d routes, expected %d", len(routes), len(tests))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, route := range routes {
				if route.Method == tt.method && route.Path == tt.route {
					found = true
				}
			}
			if !found {
				t.Fatalf("missing production route %s %s", tt.method, tt.route)
			}
			// Dispatch the concrete client path through production middleware.
			w := f.call(tt.wantStatus, tt.method, tt.path, "", "")
			if tt.wantStatus == http.StatusUnauthorized && !strings.Contains(w.Body.String(), `"code":"unauthorized"`) {
				t.Fatalf("unexpected authentication response: %s", w.Body.String())
			}
		})
	}
}
