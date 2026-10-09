package api_test

import (
	"main/api"
	"net/url"
	"strings"
	"testing"
)

func TestEndpointPaths(t *testing.T) {
	tests := []struct {
		name string
		path func(string) string
		want string
	}{
		{"delete lobby", api.DeleteLobbyPath, "/lobby/ID"},
		{"get lobby", api.GetLobbyPath, "/lobby/ID"},
		{"update lobby", api.UpdateLobbyPath, "/lobby/ID"},
		{"add user", api.AddUserToLobbyPath, "/lobby/ID/user"},
		{"remove user lobby ID", func(id string) string { return api.RemoveUserFromLobbyPath(id, "user") }, "/lobby/ID/user/user"},
		{"remove user user ID", func(id string) string { return api.RemoveUserFromLobbyPath("lobby", id) }, "/lobby/lobby/user/ID"},
		{"ready", api.SetReadyPath, "/lobby/ID/ready"},
		{"launch", api.LaunchMatchPath, "/lobby/ID/launch"},
		{"get match", api.GetMatchPath, "/match/ID"},
		{"command", api.SubmitMatchCommandPath, "/match/ID/commands"},
		{"catalogue", api.GetCataloguePath, "/match/ID/catalogue"},
		{"end match", api.EndMatchPath, "/admin/match/ID/end"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parameters resembling route placeholders must remain literal values.
			for _, id := range []string{"abc-123", "a/b c?d#e%f", "雪", ":user_id", ":id"} {
				want := strings.Replace(tt.want, "ID", url.PathEscape(id), 1)
				got := tt.path(id)
				if got != want {
					t.Errorf("path(%q) = %q, want %q", id, got, want)
				}
				parsed, err := url.Parse("https://localhost:8080" + got)
				if err != nil {
					t.Fatal(err)
				}
				if parsed.RawQuery != "" || parsed.Fragment != "" || parsed.EscapedPath() != got {
					t.Errorf("parameter changed URL structure: %s", parsed)
				}
			}
		})
	}
}

func TestRemoveUserPathEscapesBothParameters(t *testing.T) {
	got := api.RemoveUserFromLobbyPath("lobby/one", "user/two?#")
	if want := "/lobby/lobby%2Fone/user/user%2Ftwo%3F%23"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}
