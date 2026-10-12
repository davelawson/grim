package server

import (
	"reflect"
	"slices"
	"testing"

	"github.com/davelawson/grim/api"
)

func TestSelfLeaveLobby(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(3)
	path := api.LeaveLobbyPath(id)
	for _, actor := range []string{"", "invalid"} {
		assertErrorCode(t, f.call(401, api.LeaveLobby.Method, path, actor, ""), "unauthorized")
	}
	f.ready(id, 3)
	before := f.lobbyView(id, "owner")
	assertErrorCode(t, f.call(409, "DELETE", path, "owner", ""), "owner_must_transfer")
	assertErrorCode(t, f.call(403, "DELETE", path, "outsider", ""), "forbidden")
	// A member cannot remove another member through the explicit-target endpoint.
	f.call(403, api.RemoveUserFromLobby.Method, api.RemoveUserFromLobbyPath(id, f.ids["third"]), "member", "")
	if !reflect.DeepEqual(before, f.lobbyView(id, "owner")) {
		t.Fatal("failed departure changed lobby")
	}
	// A supplied target cannot change which user leaves: identity comes from auth.
	w := f.call(200, "DELETE", path+"?user_id="+f.ids["third"], "member", `{"userid":"`+f.ids["third"]+`"}`)
	if w.Body.Len() != 0 {
		t.Fatal("leave returned a body")
	}
	view := f.lobbyView(id, "owner")
	if slices.Contains(view.Members, f.ids["member"]) || !slices.Contains(view.Members, f.ids["third"]) {
		t.Fatal(view.Members)
	}
	for _, ready := range view.Readiness {
		if ready {
			t.Fatal("departure retained readiness")
		}
	}
	assertErrorCode(t, f.call(403, "DELETE", path, "member", ""), "forbidden")
	// Both compatibility modes still work: owner removal and explicit self-removal.
	f.call(200, "POST", api.JoinLobbyPath(id), "member", "")
	f.call(200, "DELETE", api.RemoveUserFromLobbyPath(id, f.ids["member"]), "owner", "")
	f.call(200, "POST", api.JoinLobbyPath(id), "member", "")
	f.call(200, "DELETE", api.RemoveUserFromLobbyPath(id, f.ids["member"]), "member", "")
}

func TestSelfLeaveUnavailableLobby(t *testing.T) {
	for _, state := range []string{"missing", "deleted", "closed"} {
		t.Run(state, func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(2)
			status, code := 404, "not_found"
			switch state {
			case "missing":
				id = "missing"
			case "deleted":
				f.call(200, "DELETE", api.DeleteLobbyPath(id), "owner", "")
			case "closed":
				f.ready(id, 2)
				f.launch(id)
				status, code = 409, "lobby_closed"
			}
			for _, actor := range []string{"owner", "member"} {
				assertErrorCode(t, f.call(status, "DELETE", api.LeaveLobbyPath(id), actor, ""), code)
			}
		})
	}
}
