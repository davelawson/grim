package server

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/davelawson/grim/api"
	"github.com/google/uuid"
)

func createLobbyAs(t *testing.T, f *lifecycleFixture, actor string) string {
	t.Helper()
	w := f.call(200, api.CreateLobby.Method, api.CreateLobby.Route, actor, `{"name":"Public game"}`)
	var response api.CreateLobbyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	return response.Id
}

func assertErrorCode(t *testing.T, w *httptest.ResponseRecorder, code string) {
	t.Helper()
	var response api.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.Error == nil || response.Error.Code != code {
		t.Fatalf("wanted %s: %s", code, w.Body.String())
	}
}

func publicLobbies(t *testing.T, f *lifecycleFixture, actor string) []api.PublicLobbySummary {
	t.Helper()
	w := f.call(200, api.ListPublicLobbies.Method, api.ListPublicLobbies.Route, actor, "")
	var response api.ListPublicLobbiesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Lobbies == nil {
		t.Fatal("null public list")
	}
	return response.Lobbies
}

func TestPublicLobbyDiscovery(t *testing.T) {
	f := newLifecycleFixture(t)
	for _, actor := range []string{"", "invalid"} {
		assertErrorCode(t, f.call(401, api.ListPublicLobbies.Method, api.ListPublicLobbies.Route, actor, ""), "unauthorized")
	}
	if got := publicLobbies(t, f, "owner"); len(got) != 0 {
		t.Fatal(got)
	}
	id := f.createLobby(4)
	if got := publicLobbies(t, f, "outsider"); len(got) != 0 {
		t.Fatal(got)
	}
	f.call(200, "DELETE", api.RemoveUserFromLobbyPath(id, f.ids["fourth"]), "fourth", "")
	second := createLobbyAs(t, f, "fourth")
	want := []api.PublicLobbySummary{
		{ID: id, Name: "Test match", Owner: f.ids["owner"], MemberCount: 3},
		{ID: second, Name: "Public game", Owner: f.ids["fourth"], MemberCount: 1},
	}
	slices.SortFunc(want, func(a, b api.PublicLobbySummary) int { return strings.Compare(a.ID, b.ID) })
	// Discovery remains available to users with an open lobby, including admins.
	for _, actor := range []string{"owner", "outsider", "admin"} {
		if got := publicLobbies(t, f, actor); !reflect.DeepEqual(got, want) {
			t.Fatalf("%+v != %+v", got, want)
		}
	}
	w := f.call(200, "GET", api.ListPublicLobbies.Route, "outsider", "")
	var shape map[string][]map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &shape); err != nil {
		t.Fatal(err)
	}
	if len(shape) != 1 {
		t.Fatal(shape)
	}
	for _, item := range shape["lobbies"] {
		if len(item) != 4 {
			t.Fatalf("unexpected summary fields: %v", item)
		}
	}
	f.call(403, "GET", api.GetLobbyPath(id), "outsider", "")
	f.call(200, "DELETE", api.DeleteLobbyPath(second), "fourth", "")
	f.ready(id, 3)
	f.launch(id)
	if got := publicLobbies(t, f, "outsider"); len(got) != 0 {
		t.Fatal(got)
	}
}

func TestPublicJoinAndLaunch(t *testing.T) {
	f := newLifecycleFixture(t)
	id := createLobbyAs(t, f, "owner")
	path := api.JoinLobbyPath(id)
	view := f.lobbyView(id, "owner")
	if view.Owner != f.ids["owner"] || len(view.Members) != 1 || view.Readiness[f.ids["owner"]] {
		t.Fatal(view)
	}
	for _, actor := range []string{"", "invalid"} {
		assertErrorCode(t, f.call(401, "POST", path, actor, ""), "unauthorized")
	}
	assertErrorCode(t, f.call(404, "POST", api.JoinLobbyPath("missing"), "member", ""), "not_found")
	f.call(403, "PUT", api.SetReadyPath(id), "member", `{"ready":true}`)
	f.call(200, "PUT", api.SetReadyPath(id), "owner", `{"ready":true}`)
	w := f.call(200, "POST", path, "member", "")
	var joined api.GetLobbyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &joined); err != nil {
		t.Fatal(err)
	}
	view = joined.Lobby
	if len(view.Members) != 2 || view.Readiness[f.ids["owner"]] || view.Readiness[f.ids["member"]] {
		t.Fatal(view)
	}
	if got := publicLobbies(t, f, "third"); len(got) != 1 || got[0].MemberCount != 2 {
		t.Fatal(got)
	}
	f.ready(id, 2)
	before := f.lobbyView(id, "member")
	f.call(200, "POST", path, "member", "")
	if !reflect.DeepEqual(before, f.lobbyView(id, "member")) {
		t.Fatal("join retry changed readiness")
	}
	f.call(200, "PUT", api.SetReadyPath(id), "member", `{"ready":false}`)
	assertErrorCode(t, f.call(409, "POST", api.LaunchMatchPath(id), "owner", launchBody(uuid.NewString())), "lobby_not_ready")
	f.call(200, "PUT", api.SetReadyPath(id), "member", `{"ready":true}`)
	assertErrorCode(t, f.call(403, "POST", api.LaunchMatchPath(id), "member", launchBody(uuid.NewString())), "forbidden")
	match, requestID, original := f.launch(id)
	if match.Status != "setup" || len(match.Participants) != 2 {
		t.Fatal(match)
	}
	if f.call(201, "POST", api.LaunchMatchPath(id), "owner", launchBody(requestID)).Body.String() != original {
		t.Fatal("launch retry changed response")
	}
	for _, actor := range []string{"owner", "third"} {
		assertErrorCode(t, f.call(409, "POST", path, actor, ""), "lobby_closed")
	}
	// Closed memberships do not block the next lobby.
	createLobbyAs(t, f, "owner")
	next := createLobbyAs(t, f, "third")
	f.call(200, "POST", api.JoinLobbyPath(next), "member", "")
	f.call(200, "DELETE", api.DeleteLobbyPath(next), "third", "")
	assertErrorCode(t, f.call(404, "POST", api.JoinLobbyPath(next), "member", ""), "not_found")
	createLobbyAs(t, f, "member")
}

func TestOpenLobbyMembershipLimits(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	other := createLobbyAs(t, f, "third")
	f.ready(id, 2)
	before := f.lobbyView(id, "owner")
	for _, actor := range []string{"owner", "member"} {
		assertErrorCode(t, f.call(409, "POST", api.CreateLobby.Route, actor, `{"name":"Another"}`), "already_in_open_lobby")
		assertErrorCode(t, f.call(409, "POST", api.JoinLobbyPath(other), actor, ""), "already_in_open_lobby")
	}
	assertErrorCode(t, f.call(409, "POST", api.AddUserToLobbyPath(id), "owner", fmt.Sprintf(`{"userid":%q}`, f.ids["third"])), "already_in_open_lobby")
	if f.count("lobbies") != 2 || !reflect.DeepEqual(before, f.lobbyView(id, "owner")) {
		t.Fatal("failed membership change wrote state")
	}
	// Voluntary departure releases membership and clears the old roster's readiness.
	f.call(200, "DELETE", api.RemoveUserFromLobbyPath(id, f.ids["member"]), "member", "")
	f.call(200, "POST", api.JoinLobbyPath(other), "member", "")
	if f.lobbyView(id, "owner").Readiness[f.ids["owner"]] {
		t.Fatal("departure retained readiness")
	}
	// Deletion preserves history but releases all open memberships.
	f.call(200, "DELETE", api.DeleteLobbyPath(other), "third", "")
	f.call(200, "POST", api.AddUserToLobbyPath(id), "owner", fmt.Sprintf(`{"userid":%q}`, f.ids["third"]))
	createLobbyAs(t, f, "member")
	f.call(200, "POST", api.JoinLobbyPath(id), "fourth", "")
	f.call(200, "POST", api.JoinLobbyPath(id), "outsider", "")
	for _, actor := range []string{"owner", "third", "fourth", "outsider"} {
		f.call(200, "PUT", api.SetReadyPath(id), actor, `{"ready":true}`)
	}
	before = f.lobbyView(id, "owner")
	assertErrorCode(t, f.call(409, "POST", api.JoinLobbyPath(id), "admin", ""), "lobby_full")
	assertErrorCode(t, f.call(409, "POST", api.AddUserToLobbyPath(id), "owner", fmt.Sprintf(`{"userid":%q}`, f.ids["admin"])), "lobby_full")
	f.call(200, "POST", api.JoinLobbyPath(id), "owner", "")
	if !reflect.DeepEqual(before, f.lobbyView(id, "owner")) {
		t.Fatal("full lobby retries changed state")
	}
}

// raceCalls starts requests together without assertions in worker goroutines.
func raceCalls(calls ...func() *httptest.ResponseRecorder) []*httptest.ResponseRecorder {
	start := make(chan struct{})
	var wg sync.WaitGroup
	results := make([]*httptest.ResponseRecorder, len(calls))
	for i, call := range calls {
		wg.Add(1)
		go func(i int, call func() *httptest.ResponseRecorder) { defer wg.Done(); <-start; results[i] = call() }(i, call)
	}
	close(start)
	wg.Wait()
	return results
}

func TestConcurrentPublicMembership(t *testing.T) {
	t.Run("last slot including owner additions", func(t *testing.T) {
		f := newLifecycleFixture(t)
		id := f.createLobby(3)
		results := raceCalls(
			func() *httptest.ResponseRecorder { return f.request("POST", api.JoinLobbyPath(id), "fourth", "") },
			func() *httptest.ResponseRecorder {
				return f.request("POST", api.AddUserToLobbyPath(id), "owner", fmt.Sprintf(`{"userid":%q}`, f.ids["outsider"]))
			},
		)
		codes := []int{results[0].Code, results[1].Code}
		slices.Sort(codes)
		if !slices.Equal(codes, []int{200, 409}) || len(f.lobbyView(id, "owner").Members) != 4 {
			t.Fatal(codes)
		}
		for _, result := range results {
			if result.Code == 409 {
				assertErrorCode(t, result, "lobby_full")
			}
		}
	})
	for _, mode := range []string{"join", "create", "owner add"} {
		t.Run("one open lobby vs "+mode, func(t *testing.T) {
			f := newLifecycleFixture(t)
			first := createLobbyAs(t, f, "owner")
			second := createLobbyAs(t, f, "third")
			results := raceCalls(
				func() *httptest.ResponseRecorder { return f.request("POST", api.JoinLobbyPath(first), "member", "") },
				func() *httptest.ResponseRecorder {
					switch mode {
					case "create":
						return f.request("POST", api.CreateLobby.Route, "member", `{"name":"Race"}`)
					case "owner add":
						return f.request("POST", api.AddUserToLobbyPath(second), "third", fmt.Sprintf(`{"userid":%q}`, f.ids["member"]))
					default:
						return f.request("POST", api.JoinLobbyPath(second), "member", "")
					}
				},
			)
			codes := []int{results[0].Code, results[1].Code}
			slices.Sort(codes)
			if !slices.Equal(codes, []int{200, 409}) {
				t.Fatal(codes)
			}
			for _, result := range results {
				if result.Code == 409 {
					assertErrorCode(t, result, "already_in_open_lobby")
				}
			}
			var count int
			if err := f.db.QueryRow(`select count(*) from lobby_users m join lobbies l on l.id=m.lobby_id where m.user_id=? and l.status='open' and l.deleted_at is null`, f.ids["member"]).Scan(&count); err != nil || count != 1 {
				t.Fatalf("memberships %d: %v", count, err)
			}
		})
	}
}

func TestJoinRacesWithLaunchOrDelete(t *testing.T) {
	for _, operation := range []string{"launch", "delete"} {
		t.Run(operation, func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(2)
			f.ready(id, 2)
			results := raceCalls(
				func() *httptest.ResponseRecorder { return f.request("POST", api.JoinLobbyPath(id), "third", "") },
				func() *httptest.ResponseRecorder {
					if operation == "launch" {
						return f.request("POST", api.LaunchMatchPath(id), "owner", launchBody(uuid.NewString()))
					}
					return f.request("DELETE", api.DeleteLobbyPath(id), "owner", "")
				},
			)
			join, mutation := results[0], results[1]
			if operation == "launch" {
				view := f.lobbyView(id, "owner")
				if mutation.Code == 201 {
					if join.Code != 409 || len(view.Members) != 2 || f.count("matches") != 1 {
						t.Fatalf("launch won: join %d; %+v", join.Code, view)
					}
					assertErrorCode(t, join, "lobby_closed")
				} else {
					if mutation.Code != 409 || join.Code != 200 || view.Status != "open" || len(view.Members) != 3 || f.count("matches") != 0 {
						t.Fatalf("join won: %d/%d; %+v", join.Code, mutation.Code, view)
					}
					assertErrorCode(t, mutation, "lobby_not_ready")
					for _, ready := range view.Readiness {
						if ready {
							t.Fatal("join retained readiness")
						}
					}
				}
			} else {
				if mutation.Code != 200 || (join.Code != 200 && join.Code != 404) || !f.deletedAt("lobbies", id).Valid {
					t.Fatalf("delete race: %d/%d", join.Code, mutation.Code)
				}
				if join.Code == 404 {
					assertErrorCode(t, join, "not_found")
				}
				createLobbyAs(t, f, "third")
			}
		})
	}
}

func TestPublicLobbyFailuresRollback(t *testing.T) {
	for _, failure := range []struct{ name, sql string }{
		{"readiness reset", `create trigger fail before update of ready on lobby_users begin select raise(abort,'private failure'); end`},
		{"commit", `create table deferred_failure (user_id text references users(id) deferrable initially deferred);
create trigger fail after insert on lobby_users begin insert into deferred_failure values('missing'); end`},
	} {
		t.Run(failure.name, func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(1)
			f.ready(id, 1)
			before := f.lobbyView(id, "owner")
			f.exec(failure.sql)
			assertErrorCode(t, f.call(500, "POST", api.JoinLobbyPath(id), "member", ""), "internal_error")
			if !reflect.DeepEqual(before, f.lobbyView(id, "owner")) {
				t.Fatal("failed join changed lobby")
			}
			f.exec("drop trigger fail")
			f.call(200, "POST", api.JoinLobbyPath(id), "member", "")
		})
	}
	f := newLifecycleFixture(t)
	f.exec("drop table lobby_users")
	assertErrorCode(t, f.call(500, "GET", api.ListPublicLobbies.Route, "owner", ""), "internal_error")
	assertErrorCode(t, f.call(500, "POST", api.CreateLobby.Route, "owner", `{"name":"Failure"}`), "internal_error")
}
