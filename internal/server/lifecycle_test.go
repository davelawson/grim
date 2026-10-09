package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/auth"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/lobby"
	"github.com/davelawson/grim/match"
	"github.com/davelawson/grim/model"
	"github.com/davelawson/grim/user"
	randv2 "math/rand/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var roster = []string{"owner", "member", "third", "fourth", "outsider"}

type lifecycleFixture struct {
	t      *testing.T
	path   string
	db     *sql.DB
	router *gin.Engine
	ids    map[string]string
}

func newLifecycleFixture(t *testing.T) *lifecycleFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)
	f := &lifecycleFixture{t: t, path: filepath.Join(t.TempDir(), "lifecycle.db"), ids: make(map[string]string)}
	f.open()
	t.Cleanup(func() { f.db.Close() })
	schema, err := os.ReadFile("../../sql/create-database.sql")
	if err != nil {
		t.Fatal(err)
	}
	f.exec(string(schema))
	for _, name := range append(slices.Clone(roster), "admin") {
		id := uuid.NewString()
		f.ids[name] = id
		f.exec("insert into users(id, email, name, password_hash, token, admin) values(?, ?, ?, ?, ?, ?)",
			id, name+"@example.com", name, []byte("unused"), "token-"+name, name == "admin")
	}
	f.wire()
	return f
}

func (f *lifecycleFixture) open() {
	f.t.Helper()
	db, err := sql.Open("sqlite3", f.path+"?_foreign_keys=on&_txlock=immediate&_busy_timeout=5000")
	if err != nil {
		f.t.Fatal(err)
	}
	f.db = db
}

func (f *lifecycleFixture) wire() {
	users := user.NewUserRepo(f.db)
	authentication := auth.NewServiceFacade(auth.NewService(users), f.db)
	matches := match.NewService(match.NewRepo(f.db), users)
	lobbies := lobby.NewService(lobby.NewLobbyRepo(f.db), users, matches)
	f.router = gin.New()
	AddAuthRoutes(f.router, auth.NewController(authentication))
	AddLobbyRoutes(authentication, f.router, lobby.NewController(lobby.NewServiceFacade(lobbies)))
	AddMatchRoutes(authentication, f.router, match.NewController(match.NewServiceFacade(matches)))
	AddUserRoutes(authentication, f.router, user.NewController(user.NewServiceFacade(user.NewService(users))))
}

func (f *lifecycleFixture) restart() {
	f.t.Helper()
	if err := f.db.Close(); err != nil {
		f.t.Fatal(err)
	}
	f.open()
	f.wire()
}

func (f *lifecycleFixture) exec(query string, args ...any) {
	f.t.Helper()
	if _, err := f.db.Exec(query, args...); err != nil {
		f.t.Fatal(err)
	}
}

func (f *lifecycleFixture) request(method, path, actor, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	if actor != "" {
		r.Header.Set("Authorization", "token-"+actor)
	}
	w := httptest.NewRecorder()
	f.router.ServeHTTP(w, r)
	return w
}

func (f *lifecycleFixture) call(status int, method, path, actor, body string) *httptest.ResponseRecorder {
	f.t.Helper()
	w := f.request(method, path, actor, body)
	if w.Code != status {
		f.t.Fatalf("%s %s as %s: status %d, want %d; %s", method, path, actor, w.Code, status, w.Body.String())
	}
	return w
}

func (f *lifecycleFixture) createLobby(size int) string {
	f.t.Helper()
	w := f.call(200, "POST", "/lobby", "owner", `{"name":"Test match"}`)
	var response api.CreateLobbyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		f.t.Fatal(err)
	}
	for _, actor := range roster[1:size] {
		f.call(200, "POST", "/lobby/"+response.Id+"/user", "owner", fmt.Sprintf(`{"userid":%q}`, f.ids[actor]))
	}
	return response.Id
}

func (f *lifecycleFixture) ready(id string, size int) {
	f.t.Helper()
	for _, actor := range roster[:size] {
		f.call(200, "PUT", "/lobby/"+id+"/ready", actor, `{"ready":true}`)
	}
}

func launchBody(requestID string) string { return fmt.Sprintf(`{"requestid":%q}`, requestID) }

func matchView(t *testing.T, w *httptest.ResponseRecorder) *game.View {
	t.Helper()
	var response api.MatchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Match == nil {
		t.Fatalf("missing match in %s", w.Body.String())
	}
	return response.Match
}

func (f *lifecycleFixture) launch(id string) (*game.View, string, string) {
	f.t.Helper()
	requestID := uuid.NewString()
	w := f.call(201, "POST", "/lobby/"+id+"/launch", "owner", launchBody(requestID))
	return matchView(f.t, w), requestID, w.Body.String()
}

func (f *lifecycleFixture) saved(id string) ([]byte, *game.State) {
	f.t.Helper()
	var snapshot []byte
	if err := f.db.QueryRow("select snapshot from matches where id = ?", id).Scan(&snapshot); err != nil {
		f.t.Fatal(err)
	}
	state, err := game.Decode(snapshot)
	if err != nil {
		f.t.Fatal(err)
	}
	return snapshot, state
}

func (f *lifecycleFixture) count(table string) int {
	f.t.Helper()
	var count int
	if err := f.db.QueryRow("select count(*) from " + table).Scan(&count); err != nil {
		f.t.Fatal(err)
	}
	return count
}

func (f *lifecycleFixture) deletedAt(table, id string) sql.NullString {
	f.t.Helper()
	var timestamp sql.NullString
	if err := f.db.QueryRow("select deleted_at from "+table+" where id = ?", id).Scan(&timestamp); err != nil {
		f.t.Fatal(err)
	}
	if timestamp.Valid {
		// The SQLite driver converts DATETIME values to RFC3339 when scanning strings.
		if parsed, err := time.Parse(time.RFC3339, timestamp.String); err != nil {
			f.t.Fatalf("invalid deletion timestamp %q: %v", timestamp.String, err)
		} else if _, offset := parsed.Zone(); offset != 0 {
			f.t.Fatalf("deletion timestamp is not UTC: %q", timestamp.String)
		}
	}
	return timestamp
}

func (f *lifecycleFixture) lobbyView(id, actor string) api.Lobby {
	f.t.Helper()
	w := f.call(200, "GET", "/lobby/"+id, actor, "")
	var response api.GetLobbyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		f.t.Fatal(err)
	}
	return response.Lobby
}

func TestMatchLifecycleAndRestart(t *testing.T) {
	for _, size := range []int{2, 3, 4} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(size)
			f.ready(id, size)
			view, requestID, original := f.launch(id)
			if f.deletedAt("lobbies", id).Valid || f.deletedAt("matches", view.ID).Valid {
				t.Fatal("new records are deleted")
			}
			if view.Status != "setup" || view.Revision != 0 || view.Outcome != nil || len(view.Participants) != size {
				t.Fatalf("unexpected setup: %+v", view)
			}
			if uuid.Validate(view.ID) != nil || view.Name != "Test match" || view.LobbyID != id {
				t.Fatalf("unexpected identity: %+v", view)
			}
			if len(view.Ring) != 0 || view.Turn != nil || view.You == nil || view.You.Draft.Step != "chantry" {
				t.Fatal("launch assigned seating or omitted the owner's draft")
			}
			for _, private := range []string{"seed", "random", "generator", "schemaversion", "admin"} {
				if strings.Contains(original, `"`+private+`"`) {
					t.Fatalf("private state exposed: %s", private)
				}
			}
			before, state := f.saved(view.ID)
			if state.RulesVersion != game.OpeningRulesVersion || state.CatalogueVersion != game.OpeningCatalogueVersion || state.Random.Generator != "chacha8-v1" {
				t.Fatal("unexpected setup binding/random algorithm")
			}
			seed := [32]byte(state.Random.Seed)
			generator := randv2.NewChaCha8(seed)
			resumed := randv2.NewChaCha8([32]byte{})
			if err := resumed.UnmarshalBinary(state.Random.State); err != nil {
				t.Fatal(err)
			}
			if generator.Uint64() != resumed.Uint64() {
				t.Fatal("launch consumed randomness before drafting")
			}
			closed := f.lobbyView(id, "member")
			if closed.Status != "closed" || closed.MatchID == nil || *closed.MatchID != view.ID {
				t.Fatalf("lobby not closed/linked: %+v", closed)
			}
			if f.count("matches") != 1 || f.count("players") != size || f.count("launch_receipts") != 1 {
				t.Fatal("incorrect launch row counts")
			}
			f.restart()
			after, _ := f.saved(view.ID)
			if !slices.Equal(before, after) {
				t.Fatal("restart changed state")
			}
			get := f.call(200, "GET", "/match/"+view.ID, "owner", "")
			if get.Body.String() != original {
				t.Fatal("lookup differs after restart")
			}
			if f.call(201, "POST", "/lobby/"+id+"/launch", "owner", launchBody(requestID)).Body.String() != original {
				t.Fatal("active launch retry failed after restart")
			}
			f.call(403, "GET", "/match/"+view.ID, "outsider", "")
			f.call(403, "GET", "/match/"+view.ID, "admin", "")
			f.call(403, "POST", "/admin/match/"+view.ID+"/end", "owner", "")
			if w := f.call(204, "POST", "/admin/match/"+view.ID+"/end", "admin", ""); w.Body.Len() != 0 {
				t.Fatal("termination returned a body")
			}
			deletedAt := f.deletedAt("matches", view.ID)
			if !deletedAt.Valid || f.deletedAt("lobbies", id).Valid {
				t.Fatal("termination did not delete only the match")
			}
			after, endedState := f.saved(view.ID)
			if !slices.Equal(before, after) || !reflect.DeepEqual(state, endedState) {
				t.Fatal("termination changed saved state")
			}
			var status string
			var revision int64
			var outcome sql.NullString
			if err := f.db.QueryRow("select status, revision, outcome from matches where id = ?", view.ID).Scan(&status, &revision, &outcome); err != nil {
				t.Fatal(err)
			}
			if status != "setup" || revision != 0 || outcome.Valid || f.count("matches") != 1 || f.count("players") != size || f.count("launch_receipts") != 1 || f.count("lobby_users") != size {
				t.Fatal("termination changed metadata or removed related rows")
			}
			f.restart()
			f.exec(`create trigger fail before update on matches begin select raise(abort, 'retry must not write'); end`)
			f.call(204, "POST", "/admin/match/"+view.ID+"/end", "admin", "")
			if f.deletedAt("matches", view.ID) != deletedAt {
				t.Fatal("repeated termination changed timestamp")
			}
			f.exec("drop trigger fail")
			for _, actor := range []string{"owner", "outsider", "admin"} {
				f.call(404, "GET", "/match/"+view.ID, actor, "")
			}
			f.call(404, "POST", "/lobby/"+id+"/launch", "owner", launchBody(requestID))
			if !reflect.DeepEqual(closed, f.lobbyView(id, "member")) {
				t.Fatal("termination changed linked lobby")
			}
			f.call(409, "POST", "/lobby/"+id+"/launch", "owner", launchBody(uuid.NewString()))
			f.call(403, "POST", "/lobby/"+id+"/launch", "member", launchBody(requestID))
			for _, mutation := range []struct{ method, suffix, actor, body string }{
				{"PUT", "", "owner", fmt.Sprintf(`{"name":"changed","owner":%q}`, f.ids["member"])},
				{"DELETE", "", "owner", ""},
				{"PUT", "/ready", "owner", `{"ready":false}`},
				{"POST", "/user", "owner", fmt.Sprintf(`{"userid":%q}`, f.ids["outsider"])},
				{"DELETE", "/user/" + f.ids["member"], "owner", ""},
				{"DELETE", "/user/" + f.ids["member"], "member", ""},
			} {
				f.call(409, mutation.method, "/lobby/"+id+mutation.suffix, mutation.actor, mutation.body)
			}
			f.exec("update users set admin = 0 where id = ?", f.ids["admin"])
			f.call(403, "POST", "/admin/match/"+view.ID+"/end", "admin", "")
		})
	}
}

func TestLaunchValidationAndAuthorization(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	path := "/lobby/" + id + "/launch"
	for _, actor := range []string{"", "invalid"} {
		f.call(401, "POST", path, actor, launchBody(uuid.NewString()))
	}
	for _, actor := range []string{"member", "outsider", "admin"} {
		f.call(403, "POST", path, actor, launchBody(uuid.NewString()))
	}
	for _, body := range []string{"", "null", `{}`, `[]`, `{"requestid":null}`, `{"requestid":42}`, `{"requestid":""}`, `{"requestid":"not-a-uuid"}`} {
		f.call(400, "POST", path, "owner", body)
	}
	f.call(409, "POST", path, "owner", launchBody(uuid.NewString()))
	f.call(200, "PUT", "/lobby/"+id+"/ready", "owner", `{"ready":true}`)
	f.call(409, "POST", path, "owner", launchBody(uuid.NewString()))
	if f.count("matches") != 0 || f.count("players") != 0 || f.count("launch_receipts") != 0 || f.lobbyView(id, "owner").Status != "open" {
		t.Fatal("invalid launch wrote state")
	}
	f.call(404, "POST", "/lobby/missing/launch", "owner", launchBody(uuid.NewString()))
	f.call(404, "GET", "/match/missing", "owner", "")
	f.call(404, "POST", "/admin/match/missing/end", "admin", "")
	for _, size := range []int{1, 5} {
		id := f.createLobby(size)
		f.ready(id, size)
		w := f.call(409, "POST", "/lobby/"+id+"/launch", "owner", launchBody(uuid.NewString()))
		if !strings.Contains(w.Body.String(), "invalid_roster_size") {
			t.Fatal(w.Body.String())
		}
	}
}

func TestReadinessRosterAndOwnership(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	base := "/lobby/" + id
	for _, body := range []string{"", "null", `{}`, `{"ready":null}`, `{"ready":"true"}`, `{"ready":0}`} {
		f.call(400, "PUT", base+"/ready", "owner", body)
	}
	f.call(403, "GET", base, "outsider", "")
	f.call(403, "PUT", base+"/ready", "outsider", `{"ready":true}`)
	f.call(403, "PUT", base, "member", fmt.Sprintf(`{"name":"forged","owner":%q}`, f.ids["member"]))
	f.call(400, "PUT", base, "owner", fmt.Sprintf(`{"name":"invalid owner","owner":%q}`, f.ids["outsider"]))
	f.call(403, "POST", base+"/user", "member", fmt.Sprintf(`{"userid":%q}`, f.ids["third"]))
	f.call(403, "DELETE", base+"/user/"+f.ids["owner"], "member", "")
	f.call(409, "DELETE", base+"/user/"+f.ids["owner"], "owner", "")
	f.ready(id, 2)
	f.call(200, "PUT", base+"/ready", "owner", `{"ready":false}`)
	view := f.lobbyView(id, "owner")
	if view.Readiness[f.ids["owner"]] || !view.Readiness[f.ids["member"]] {
		t.Fatal("readiness changed another member")
	}
	f.call(200, "PUT", base+"/ready", "owner", `{"Ready":true}`)
	f.call(200, "PUT", base, "owner", fmt.Sprintf(`{"name":"Transferred","owner":%q}`, f.ids["member"]))
	view = f.lobbyView(id, "owner")
	if view.Owner != f.ids["member"] || !view.Readiness[f.ids["owner"]] || !view.Readiness[f.ids["member"]] {
		t.Fatal("transfer lost readiness")
	}
	f.call(403, "PUT", base, "owner", fmt.Sprintf(`{"name":"stolen","owner":%q}`, f.ids["owner"]))
	f.call(200, "DELETE", base+"/user/"+f.ids["owner"], "owner", "")
	f.call(403, "GET", base, "owner", "")
	view = f.lobbyView(id, "member")
	if len(view.Members) != 1 || view.Readiness[f.ids["member"]] {
		t.Fatal("departure did not clear readiness")
	}
	f.call(200, "PUT", base+"/ready", "member", `{"ready":true}`)
	f.call(200, "POST", base+"/user", "member", fmt.Sprintf(`{"userid":%q}`, f.ids["third"]))
	view = f.lobbyView(id, "member")
	if view.Readiness[f.ids["member"]] || view.Readiness[f.ids["third"]] {
		t.Fatal("addition did not clear readiness")
	}
	f.call(200, "PUT", base+"/ready", "member", `{"ready":true}`)
	f.call(200, "PUT", base+"/ready", "third", `{"ready":true}`)
	f.call(400, "POST", base+"/user", "member", fmt.Sprintf(`{"userid":%q}`, f.ids["third"]))
	if !f.lobbyView(id, "member").Readiness[f.ids["member"]] {
		t.Fatal("failed duplicate addition changed readiness")
	}
	f.call(200, "DELETE", base+"/user/"+f.ids["third"], "member", "")
	if f.lobbyView(id, "member").Readiness[f.ids["member"]] {
		t.Fatal("owner removal did not clear readiness")
	}
	f.call(200, "DELETE", base, "member", "")
	f.call(404, "GET", base, "member", "")
}

func TestLaunchRollback(t *testing.T) {
	for _, target := range []struct{ name, trigger string }{
		{"membership", `create trigger fail before insert on players begin select raise(abort, 'private injected failure'); end`},
		{"closure", `create trigger fail before update of status on lobbies begin select raise(abort, 'private injected failure'); end`},
		{"receipt", `create trigger fail before insert on launch_receipts begin select raise(abort, 'private injected failure'); end`},
		{"commit", `create table deferred_failure(id text references users(id) deferrable initially deferred);
create trigger fail after insert on launch_receipts begin insert into deferred_failure values('nonexistent'); end`},
	} {
		t.Run(target.name, func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(2)
			f.ready(id, 2)
			f.exec(target.trigger)
			requestID := uuid.NewString()
			w := f.call(500, "POST", "/lobby/"+id+"/launch", "owner", launchBody(requestID))
			if strings.Contains(w.Body.String(), "private") || strings.Contains(w.Body.String(), "FOREIGN") {
				t.Fatal("failure leaked internal information")
			}
			if f.count("matches") != 0 || f.count("players") != 0 || f.count("launch_receipts") != 0 {
				t.Fatal("failed launch left rows")
			}
			view := f.lobbyView(id, "owner")
			if view.Status != "open" || view.MatchID != nil || !view.Readiness[f.ids["owner"]] || !view.Readiness[f.ids["member"]] {
				t.Fatal("failed launch changed lobby")
			}
			f.exec("drop trigger fail")
			f.call(201, "POST", "/lobby/"+id+"/launch", "owner", launchBody(requestID))
		})
	}
}

func TestConcurrentLaunches(t *testing.T) {
	for _, identical := range []bool{true, false} {
		t.Run(fmt.Sprintf("identical=%t", identical), func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(4)
			f.ready(id, 4)
			requestID := uuid.NewString()
			start := make(chan struct{})
			results := make(chan *httptest.ResponseRecorder, 6)
			var wg sync.WaitGroup
			for i := 0; i < 6; i++ {
				body := launchBody(requestID)
				if !identical {
					body = launchBody(uuid.NewString())
				}
				wg.Add(1)
				go func() { defer wg.Done(); <-start; results <- f.request("POST", "/lobby/"+id+"/launch", "owner", body) }()
			}
			close(start)
			wg.Wait()
			close(results)
			var original string
			successes := 0
			for w := range results {
				if w.Code == 201 {
					successes++
					if original != "" && original != w.Body.String() {
						t.Fatal("retries returned different results")
					}
					original = w.Body.String()
				} else if identical || w.Code != 409 {
					t.Fatalf("unexpected concurrent result: %d %s", w.Code, w.Body.String())
				}
			}
			if (identical && successes != 6) || (!identical && successes != 1) {
				t.Fatalf("unexpected successful launch count %d", successes)
			}
			if f.count("matches") != 1 || f.count("players") != 4 || f.count("launch_receipts") != 1 {
				t.Fatal("concurrent launch duplicated state")
			}
		})
	}
}

func TestLaunchRacesWithRosterChange(t *testing.T) {
	for i := 0; i < 6; i++ {
		f := newLifecycleFixture(t)
		id := f.createLobby(2)
		f.ready(id, 2)
		start := make(chan struct{})
		launches, additions := make(chan *httptest.ResponseRecorder, 1), make(chan *httptest.ResponseRecorder, 1)
		go func() {
			<-start
			launches <- f.request("POST", "/lobby/"+id+"/launch", "owner", launchBody(uuid.NewString()))
		}()
		go func() {
			<-start
			additions <- f.request("POST", "/lobby/"+id+"/user", "owner", fmt.Sprintf(`{"userid":%q}`, f.ids["third"]))
		}()
		close(start)
		launched, added := <-launches, <-additions
		switch {
		case launched.Code == 201 && added.Code == 409:
			view := matchView(t, launched)
			if len(view.Participants) != 2 || f.count("matches") != 1 {
				t.Fatal("launch used mutated roster")
			}
		case launched.Code == 409 && added.Code == 200:
			view := f.lobbyView(id, "owner")
			if view.Status != "open" || len(view.Members) != 3 || f.count("matches") != 0 {
				t.Fatal("addition left invalid lobby")
			}
			for _, ready := range view.Readiness {
				if ready {
					t.Fatal("addition retained readiness")
				}
			}
		default:
			t.Fatalf("unexpected race results: launch %d %s; add %d %s", launched.Code, launched.Body.String(), added.Code, added.Body.String())
		}
	}
}

func TestEndRollbackAndConcurrentRetries(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	f.ready(id, 2)
	view, _, _ := f.launch(id)
	before, _ := f.saved(view.ID)
	f.exec(`create trigger fail before update on matches begin select raise(abort, 'private failure'); end`)
	w := f.call(500, "POST", "/admin/match/"+view.ID+"/end", "admin", "")
	if strings.Contains(w.Body.String(), "private") || f.deletedAt("matches", view.ID).Valid {
		t.Fatal("failed termination leaked details or left a deletion timestamp")
	}
	after, _ := f.saved(view.ID)
	if !slices.Equal(before, after) {
		t.Fatal("failed ending changed state")
	}
	f.exec("drop trigger fail")
	f.exec(`create table deferred_failure(id text references users(id) deferrable initially deferred);
create trigger fail after update on matches begin insert into deferred_failure values('nonexistent'); end`)
	w = f.call(500, "POST", "/admin/match/"+view.ID+"/end", "admin", "")
	if strings.Contains(w.Body.String(), "FOREIGN") || f.deletedAt("matches", view.ID).Valid {
		t.Fatal("failed termination commit leaked details or left a deletion timestamp")
	}
	after, _ = f.saved(view.ID)
	if !slices.Equal(before, after) {
		t.Fatal("failed end commit changed state")
	}
	f.exec("drop trigger fail")
	start := make(chan struct{})
	results := make(chan *httptest.ResponseRecorder, 4)
	for i := 0; i < 4; i++ {
		go func() { <-start; results <- f.request("POST", "/admin/match/"+view.ID+"/end", "admin", "") }()
	}
	close(start)
	for i := 0; i < 4; i++ {
		w := <-results
		if w.Code != 204 || w.Body.Len() != 0 {
			t.Fatalf("unexpected concurrent end: %d %s", w.Code, w.Body.String())
		}
	}
	after, _ = f.saved(view.ID)
	if !slices.Equal(before, after) || !f.deletedAt("matches", view.ID).Valid {
		t.Fatal("termination changed saved state or failed to set timestamp")
	}
}

func TestUnsupportedStateAndExistingOutcome(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	f.ready(id, 2)
	view, requestID, original := f.launch(id)
	_, state := f.saved(view.ID)
	state.SchemaVersion = 99
	unsupported, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	f.exec("update matches set snapshot = ? where id = ?", string(unsupported), view.ID)
	w := f.call(409, "GET", "/match/"+view.ID, "owner", "")
	if !strings.Contains(w.Body.String(), "unsupported_state") {
		t.Fatal(w.Body.String())
	}
	f.call(403, "GET", "/match/"+view.ID, "outsider", "")
	var unchanged []byte
	if err := f.db.QueryRow("select snapshot from matches where id = ?", view.ID).Scan(&unchanged); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(unchanged, unsupported) {
		t.Fatal("unsupported snapshot was rewritten")
	}
	if f.call(201, "POST", "/lobby/"+id+"/launch", "owner", launchBody(requestID)).Body.String() != original {
		t.Fatal("unsupported newer snapshot broke launch receipt")
	}
	f.exec("update launch_receipts set command = 'different' where lobby_id = ?", id)
	w = f.call(409, "POST", "/lobby/"+id+"/launch", "owner", launchBody(requestID))
	if !strings.Contains(w.Body.String(), "request_id_conflict") {
		t.Fatal(w.Body.String())
	}
	f.exec("update launch_receipts set command = 'launch' where lobby_id = ?", id)
	f.call(204, "POST", "/admin/match/"+view.ID+"/end", "admin", "")
	if !f.deletedAt("matches", view.ID).Valid {
		t.Fatal("unsupported match was not deleted")
	}
	if err := f.db.QueryRow("select snapshot from matches where id = ?", view.ID).Scan(&unchanged); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(unchanged, unsupported) {
		t.Fatal("termination rewrote unsupported snapshot")
	}
	f.call(404, "GET", "/match/"+view.ID, "owner", "")
	f.call(404, "POST", "/lobby/"+id+"/launch", "owner", launchBody(requestID))
	// Set up a separate finished match; no restore operation is exposed.
	id = f.createLobby(2)
	f.ready(id, 2)
	view, _, _ = f.launch(id)
	_, state = f.saved(view.ID)
	state.SchemaVersion = game.SchemaVersion
	outcome := "victory"
	state.Status, state.Outcome, state.Revision = "finished", &outcome, 7
	terminal, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	f.exec("update matches set status = 'finished', outcome = ?, revision = 7, snapshot = ? where id = ?", outcome, string(terminal), view.ID)
	f.call(204, "POST", "/admin/match/"+view.ID+"/end", "admin", "")
	after, _ := f.saved(view.ID)
	var status, savedOutcome string
	var revision int64
	if err := f.db.QueryRow("select status, outcome, revision from matches where id = ?", view.ID).Scan(&status, &savedOutcome, &revision); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(after, terminal) || status != "finished" || savedOutcome != outcome || revision != 7 || !f.deletedAt("matches", view.ID).Valid {
		t.Fatal("admin overwrote existing victory")
	}
}

func TestLobbySoftDeletion(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	f.ready(id, 2)
	base := "/lobby/" + id
	for _, actor := range []string{"member", "outsider", "admin"} {
		f.call(403, "DELETE", base, actor, "")
	}
	if f.deletedAt("lobbies", id).Valid {
		t.Fatal("unauthorized deletion changed lobby")
	}
	f.call(200, "DELETE", base, "owner", "")
	deletedAt := f.deletedAt("lobbies", id)
	if !deletedAt.Valid || f.count("lobbies") != 1 || f.count("lobby_users") != 2 {
		t.Fatal("deletion removed lobby data or failed to set timestamp")
	}
	f.restart()
	for _, request := range []struct{ method, suffix, body string }{
		{"GET", "", ""},
		{"DELETE", "", ""},
		{"PUT", "", fmt.Sprintf(`{"name":"changed","owner":%q}`, f.ids["owner"])},
		{"PUT", "/ready", `{"ready":false}`},
		{"POST", "/user", fmt.Sprintf(`{"userid":%q}`, f.ids["third"])},
		{"DELETE", "/user/" + f.ids["member"], ""},
		{"POST", "/launch", launchBody(uuid.NewString())},
	} {
		f.call(404, request.method, base+request.suffix, "owner", request.body)
	}
	if f.deletedAt("lobbies", id) != deletedAt || f.count("matches") != 0 {
		t.Fatal("deleted lobby was modified or launched")
	}
	var name, status string
	var readyCount int
	if err := f.db.QueryRow("select name, status from lobbies where id = ?", id).Scan(&name, &status); err != nil {
		t.Fatal(err)
	}
	if err := f.db.QueryRow("select sum(ready) from lobby_users where lobby_id = ?", id).Scan(&readyCount); err != nil {
		t.Fatal(err)
	}
	if name != "Test match" || status != "open" || readyCount != 2 {
		t.Fatal("deletion changed lobby or membership state")
	}
	repo := lobby.NewLobbyRepo(f.db)
	if found, err := repo.GetLobbyByNameAndOwner("Test match", f.ids["owner"]); err != nil || found != nil {
		t.Fatalf("name lookup exposed deleted lobby: %+v %v", found, err)
	}
	newID := f.createLobby(2)
	if found, err := repo.GetLobbyByNameAndOwner("Test match", f.ids["owner"]); err != nil || found == nil || found.Id != newID {
		t.Fatalf("name lookup failed for replacement lobby: %+v %v", found, err)
	}
}

func TestLobbyDeletionRollback(t *testing.T) {
	for _, failure := range []struct{ name, trigger string }{
		{"statement", `create trigger fail before update of deleted_at on lobbies begin select raise(abort, 'private failure'); end`},
		{"commit", `create table deferred_failure(id text references users(id) deferrable initially deferred);
create trigger fail after update of deleted_at on lobbies begin insert into deferred_failure values('nonexistent'); end`},
	} {
		t.Run(failure.name, func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(2)
			f.ready(id, 2)
			before := f.lobbyView(id, "owner")
			f.exec(failure.trigger)
			w := f.call(500, "DELETE", "/lobby/"+id, "owner", "")
			if strings.Contains(w.Body.String(), "private") || strings.Contains(w.Body.String(), "FOREIGN") || f.deletedAt("lobbies", id).Valid || !reflect.DeepEqual(before, f.lobbyView(id, "owner")) {
				t.Fatal("failed deletion leaked details or changed lobby")
			}
		})
	}
}

func TestAdminFlagNotPubliclyWritable(t *testing.T) {
	f := newLifecycleFixture(t)
	f.call(200, "POST", "/user", "", `{"email":"new@example.com","name":"New","password":"test","admin":true}`)
	var admin bool
	if err := f.db.QueryRow("select admin from users where email = 'new@example.com'").Scan(&admin); err != nil {
		t.Fatal(err)
	}
	if admin {
		t.Fatal("public registration granted admin")
	}
	tx, err := f.db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	if err := user.NewUserRepo(f.db).UpdateUser(tx, &model.User{Id: f.ids["admin"], Name: "renamed", PasswordHash: []byte("unused"), Token: "token-admin"}); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if err := f.db.QueryRow("select admin from users where id = ?", f.ids["admin"]).Scan(&admin); err != nil {
		t.Fatal(err)
	}
	if !admin {
		t.Fatal("ordinary user update lost admin flag")
	}
}

type failedAuthentication struct{}

func (failedAuthentication) VerifyBearerToken(string) (*model.User, error) {
	return nil, fmt.Errorf("private database error")
}

func TestAuthenticationErrors(t *testing.T) {
	f := newLifecycleFixture(t)
	for _, route := range []struct{ method, path, body string }{
		{"POST", "/user/getbyemail", `{"email":"owner@example.com"}`},
		{"POST", "/lobby", `{"name":"Test"}`},
		{"GET", "/lobby/missing", ""},
		{"PUT", "/lobby/missing", `{"name":"Test","owner":"missing"}`},
		{"DELETE", "/lobby/missing", ""},
		{"POST", "/lobby/missing/user", `{"userid":"missing"}`},
		{"DELETE", "/lobby/missing/user/missing", ""},
		{"PUT", "/lobby/missing/ready", `{"ready":true}`},
		{"POST", "/lobby/missing/launch", launchBody(uuid.NewString())},
		{"GET", "/match/missing", ""},
		{"GET", "/match/missing/catalogue", ""},
		{"POST", "/match/missing/commands", commandBody(uuid.NewString(), 0, "choose_chantry", "chantry_fire")},
		{"POST", "/admin/match/missing/end", ""},
	} {
		for _, actor := range []string{"", "invalid"} {
			w := f.call(401, route.method, route.path, actor, route.body)
			var response api.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.Error == nil || response.Error.Code != "unauthorized" || response.Error.Message == "" || !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
				t.Fatalf("invalid authentication response for %s %s: %s", route.method, route.path, w.Body.String())
			}
		}
	}
	router := gin.New()
	router.GET("/match/failure", createAuthedHandler(failedAuthentication{}, func(*gin.Context) { t.Fatal("handler called after auth error") }))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/match/failure", nil))
	var response api.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != 500 || response.Error == nil || response.Error.Code != "internal_error" || strings.Contains(w.Body.String(), "private") || !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("unsafe auth failure: %d %s", w.Code, w.Body.String())
	}
}

func TestFreshSchemaAndFixtures(t *testing.T) {
	f := newLifecycleFixture(t)
	var admin int
	if err := f.db.QueryRow("select admin from users where id = ?", f.ids["owner"]).Scan(&admin); err != nil || admin != 0 {
		t.Fatalf("unexpected admin default: %d %v", admin, err)
	}
	id := f.createLobby(2)
	var createdAt string
	if err := f.db.QueryRow("select created_at from lobbies where id = ?", id).Scan(&createdAt); err != nil || createdAt == "" {
		t.Fatalf("lobby timestamp missing: %q %v", createdAt, err)
	}
	if f.count("matches") != 0 {
		t.Fatal("fresh fixture has matches")
	}
	fixture, err := os.ReadFile("../../sql/populate-test-data.sql")
	if err != nil {
		t.Fatal(err)
	}
	// SQLite's shell readfile() is unavailable through database/sql.
	f.exec(strings.ReplaceAll(string(fixture), "readfile('password_hash.dat')", "x'00'"))
	rows, err := f.db.Query("pragma foreign_key_check")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Fatal("schema/fixtures violate foreign keys")
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
}
