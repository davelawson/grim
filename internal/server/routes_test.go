package server

import (
	"main/api"
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
		{"login", api.LoginMethod, api.LoginRoute, api.LoginRoute, 400},
		{"create user", api.CreateUserMethod, api.CreateUserRoute, api.CreateUserRoute, 400},
		{"get user", api.GetUserByEmailMethod, api.GetUserByEmailRoute, api.GetUserByEmailRoute, 401},
		{"create lobby", api.CreateLobbyMethod, api.CreateLobbyRoute, api.CreateLobbyRoute, 401},
		{"delete lobby", api.DeleteLobbyMethod, api.DeleteLobbyRoute, api.DeleteLobbyPath(lobbyID), 401},
		{"get lobby", api.GetLobbyMethod, api.GetLobbyRoute, api.GetLobbyPath(lobbyID), 401},
		{"update lobby", api.UpdateLobbyMethod, api.UpdateLobbyRoute, api.UpdateLobbyPath(lobbyID), 401},
		{"add user", api.AddUserToLobbyMethod, api.AddUserToLobbyRoute, api.AddUserToLobbyPath(lobbyID), 401},
		{"remove user", api.RemoveUserFromLobbyMethod, api.RemoveUserFromLobbyRoute, api.RemoveUserFromLobbyPath(lobbyID, userID), 401},
		{"ready", api.SetReadyMethod, api.SetReadyRoute, api.SetReadyPath(lobbyID), 401},
		{"launch", api.LaunchMatchMethod, api.LaunchMatchRoute, api.LaunchMatchPath(lobbyID), 401},
		{"get match", api.GetMatchMethod, api.GetMatchRoute, api.GetMatchPath(lobbyID), 401},
		{"command", api.SubmitMatchCommandMethod, api.SubmitMatchCommandRoute, api.SubmitMatchCommandPath(lobbyID), 401},
		{"catalogue", api.GetCatalogueMethod, api.GetCatalogueRoute, api.GetCataloguePath(lobbyID), 401},
		{"end match", api.EndMatchMethod, api.EndMatchRoute, api.EndMatchPath(lobbyID), 401},
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
