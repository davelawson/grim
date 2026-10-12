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
