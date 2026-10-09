package server

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/davelawson/grim/api"
)

func TestMembershipLists(t *testing.T) {
	f := newLifecycleFixture(t)
	// Insert deliberately unreadable snapshots: listings must use metadata only.
	for _, item := range []struct{ id, status string }{{"match-c", "finished"}, {"match-a", "setup"}, {"match-b", "playing"}, {"match-deleted", "setup"}, {"match-outsider", "setup"}} {
		lobbyID := "lobby-" + item.id
		f.exec("insert into lobbies(id,name,owner_id) values(?,?,?)", lobbyID, "private name", f.ids["owner"])
		var outcome any
		if item.status == "finished" {
			outcome = "abandoned"
		}
		f.exec("insert into matches(id,lobby_id,name,status,revision,outcome,snapshot) values(?,?,?,?,0,?,?)", item.id, lobbyID, "private name", item.status, outcome, `{"schemaversion":999}`)
		f.exec("update lobbies set status='closed',match_id=? where id=?", item.id, lobbyID)
		actor := "owner"
		if item.id == "match-outsider" {
			actor = "outsider"
		}
		f.exec("insert into players(player_id,match_id) values(?,?)", f.ids[actor], item.id)
		if actor == "owner" {
			f.exec("insert into players(player_id,match_id) values(?,?)", f.ids["member"], item.id)
		}
	}
	f.exec("update matches set deleted_at=datetime('now') where id='match-deleted'")
	for _, id := range []string{"lobby-a", "lobby-c", "lobby-deleted", "lobby-departed", "lobby-outsider"} {
		f.exec("insert into lobbies(id,name,owner_id) values(?,?,?)", id, "private lobby", f.ids["owner"])
	}
	for _, id := range []string{"lobby-a", "lobby-c", "lobby-deleted", "lobby-departed", "lobby-match-a"} {
		f.exec("insert into lobby_users(lobby_id,user_id) values(?,?)", id, f.ids["owner"])
	}
	f.exec("insert into lobby_users(lobby_id,user_id) values('lobby-a',?)", f.ids["member"])
	f.exec("insert into lobby_users(lobby_id,user_id) values('lobby-outsider',?)", f.ids["outsider"])
	f.exec("update lobbies set deleted_at=datetime('now') where id='lobby-deleted'")
	f.exec("delete from lobby_users where lobby_id='lobby-departed' and user_id=?", f.ids["owner"])
	check := func(path, key, actor string, want []map[string]string) {
		t.Helper()
		response := f.call(http.StatusOK, http.MethodGet, path, actor, "")
		var envelope map[string][]map[string]string
		if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
			t.Fatal(err)
		}
		if len(envelope) != 1 || !reflect.DeepEqual(envelope[key], want) {
			t.Fatalf("%s as %s: %s", path, actor, response.Body.String())
		}
	}
	check(api.ListMatchesRoute, "matches", "owner", []map[string]string{{"id": "match-a", "status": "setup"}, {"id": "match-b", "status": "playing"}, {"id": "match-c", "status": "finished"}})
	check(api.ListLobbiesRoute, "lobbies", "owner", []map[string]string{{"id": "lobby-a", "status": "open"}, {"id": "lobby-c", "status": "open"}, {"id": "lobby-match-a", "status": "closed"}})
	check(api.ListMatchesRoute, "matches", "outsider", []map[string]string{{"id": "match-outsider", "status": "setup"}})
	check(api.ListLobbiesRoute, "lobbies", "outsider", []map[string]string{{"id": "lobby-outsider", "status": "open"}})
	check(api.ListLobbiesRoute, "lobbies", "member", []map[string]string{{"id": "lobby-a", "status": "open"}})
	for _, actor := range []string{"third", "admin"} {
		check(api.ListMatchesRoute, "matches", actor, []map[string]string{})
		check(api.ListLobbiesRoute, "lobbies", actor, []map[string]string{})
	}
	for _, path := range []string{api.ListMatchesRoute, api.ListLobbiesRoute} {
		for _, actor := range []string{"", "invalid"} {
			f.call(http.StatusUnauthorized, http.MethodGet, path, actor, "")
		}
	}
	// Database failures use the established generic API error envelope.
	f.exec("drop table players")
	failure := f.call(http.StatusInternalServerError, api.ListMatchesMethod, api.ListMatchesRoute, "owner", "")
	var result api.ErrorResponse
	if err := json.Unmarshal(failure.Body.Bytes(), &result); err != nil || result.Error == nil || result.Error.Code != "internal_error" {
		t.Fatalf("unexpected failure: %s", failure.Body.String())
	}
}
