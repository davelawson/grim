package server

import (
	"encoding/json"
	"fmt"
	"main/api"
	"main/game"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func commandBody(requestID string, revision int64, kind, card string) string {
	return fmt.Sprintf(`{"requestid":%q,"expectedrevision":%d,"type":%q,"cardid":%q}`, requestID, revision, kind, card)
}

func (f *lifecycleFixture) current(id, actor string) *game.View {
	f.t.Helper()
	return matchView(f.t, f.call(200, "GET", "/match/"+id, actor, ""))
}

func (f *lifecycleFixture) choose(id, actor, kind, card string) (*game.View, string, string) {
	f.t.Helper()
	requestID := uuid.NewString()
	body := commandBody(requestID, f.current(id, actor).Revision, kind, card)
	w := f.call(200, "POST", "/match/"+id+"/commands", actor, body)
	return matchView(f.t, w), body, w.Body.String()
}

func (f *lifecycleFixture) draft(id, actor, chantry, path, minion string) *game.View {
	f.t.Helper()
	f.choose(id, actor, "choose_chantry", chantry)
	f.choose(id, actor, "choose_victory_path", path)
	view, _, _ := f.choose(id, actor, "choose_minion", minion)
	return view
}

func errorCode(t *testing.T, w *httptest.ResponseRecorder, want string) {
	t.Helper()
	var result struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil || result.Error.Code != want {
		t.Fatalf("want %s: %s", want, w.Body.String())
	}
}

func TestDraftingThroughGameStartAndRestart(t *testing.T) {
	for _, size := range []int{2, 3, 4} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			f := newLifecycleFixture(t)
			lobbyID := f.createLobby(size)
			f.ready(lobbyID, size)
			launch, launchID, launchResponse := f.launch(lobbyID)
			_, initial := f.saved(launch.ID)
			catalogueResponse := f.call(200, "GET", "/match/"+launch.ID+"/catalogue", "member", "")
			var catalogue api.CatalogueResponse
			if err := json.Unmarshal(catalogueResponse.Body.Bytes(), &catalogue); err != nil {
				t.Fatal(err)
			}
			if catalogue.Catalogue.CatalogueVersion != launch.CatalogueVersion || len(catalogue.Catalogue.Cards) != 31 || !reflect.DeepEqual(catalogue.Catalogue.DraftOffers, *launch.DraftOffers) {
				t.Fatal("wrong pinned catalogue")
			}
			for _, actor := range roster[:size-1] {
				view := f.draft(launch.ID, actor, "chantry_fire", "victory_fame", "orla_nine_cinders")
				if view.Status != "setup" || len(view.Ring) != 0 || view.Turn != nil || len(view.You.Hand) != 0 {
					t.Fatal("started before everyone drafted")
				}
				other := f.current(launch.ID, roster[size-1])
				if !other.Wizards[f.ids[actor]].DraftComplete || other.You.Draft.Step != "chantry" {
					t.Fatal("wrong private progress")
				}
				_, saved := f.saved(launch.ID)
				if !reflect.DeepEqual(saved.Random, initial.Random) {
					t.Fatal("partial draft consumed randomness")
				}
			}
			actor := roster[size-1]
			f.choose(launch.ID, actor, "choose_chantry", "chantry_water")
			f.choose(launch.ID, actor, "choose_victory_path", "subjugate_rival")
			f.restart()
			view, finalBody, finalResponse := f.choose(launch.ID, actor, "choose_minion", "nadia_rivermark")
			if view.Status != "playing" || view.Revision != int64(size*3) || view.Turn == nil || view.Turn.ActiveWizard != view.Ring[0] || view.Turn.Phase != "start" || view.Turn.Step != "recovery" {
				t.Fatalf("bad started view %+v", view)
			}
			ring := slices.Clone(view.Ring)
			slices.Sort(ring)
			if !slices.Equal(ring, view.Participants) {
				t.Fatal("seating lost participants")
			}
			for _, public := range view.Wizards {
				if public.Integrity == nil || *public.Integrity != 3 || public.HandCount != 5 || public.DrawCount != 5 || len(public.InPlay) != 1 || !public.DraftComplete {
					t.Fatal("wrong public start state")
				}
			}
			if len(view.You.Hand) != 5 || len(view.You.DrawPile) != 5 || view.DraftOffers != nil {
				t.Fatal("incorrect private start state")
			}
			// Opponent selections remain absent even when drafting completes.
			for _, hidden := range []string{`"random"`, `"seed"`, `"generator"`, `"schemaversion"`} {
				if strings.Contains(finalResponse, hidden) {
					t.Fatalf("leaked %s", hidden)
				}
			}
			if strings.Contains(finalResponse, "orla_nine_cinders") || strings.Contains(finalResponse, "victory_fame") {
				t.Fatal("opponent drafted cards leaked")
			}
			before, state := f.saved(launch.ID)
			if len(state.Wizards) != size || reflect.DeepEqual(state.Random, initial.Random) || f.count("command_receipts") != size*3 {
				t.Fatal("bad persisted start")
			}
			f.restart()
			if f.call(200, "GET", "/match/"+view.ID, actor, "").Body.String() != finalResponse {
				t.Fatal("restart changed participant view")
			}
			if f.call(200, "POST", "/match/"+view.ID+"/commands", actor, finalBody).Body.String() != finalResponse {
				t.Fatal("final retry changed response")
			}
			after, _ := f.saved(view.ID)
			if !slices.Equal(before, after) || f.count("command_receipts") != size*3 {
				t.Fatal("retry reshuffled or saved again")
			}
			if f.call(201, "POST", "/lobby/"+lobbyID+"/launch", "owner", launchBody(launchID)).Body.String() != launchResponse {
				t.Fatal("drafting changed launch receipt")
			}
			f.call(200, "GET", "/match/"+view.ID+"/catalogue", actor, "")
			errorCode(t, f.call(409, "POST", "/match/"+view.ID+"/commands", actor, commandBody(uuid.NewString(), view.Revision, "choose_minion", "nadia_rivermark")), "draft_closed")
		})
	}
}

func TestDraftCommandValidationAuthorizationAndRetries(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	f.ready(id, 2)
	view, _, _ := f.launch(id)
	path := "/match/" + view.ID + "/commands"
	valid := commandBody(uuid.NewString(), 0, "choose_chantry", "chantry_fire")
	for _, actor := range []string{"", "invalid"} {
		f.call(401, "POST", path, actor, valid)
		f.call(401, "GET", "/match/"+view.ID+"/catalogue", actor, "")
	}
	for _, actor := range []string{"outsider", "admin"} {
		f.call(403, "POST", path, actor, valid)
		f.call(403, "GET", "/match/"+view.ID+"/catalogue", actor, "")
	}
	f.call(404, "POST", "/match/missing/commands", "owner", valid)
	f.call(404, "GET", "/match/missing/catalogue", "owner", "")
	for _, body := range []string{"", "null", `{}`, `[]`, valid + valid,
		strings.Replace(valid, `"expectedrevision":0,`, "", 1),
		strings.Replace(valid, `"expectedrevision":0`, `"expectedrevision":null`, 1),
		strings.Replace(valid, `"expectedrevision":0`, `"expectedrevision":-1`, 1),
		strings.Replace(valid, `"expectedrevision":0`, `"expectedrevision":"0"`, 1),
		strings.Replace(valid, `"expectedrevision":0`, `"expectedrevision":0.5`, 1),
		strings.Replace(valid, `"requestid":"`, `"requestid":"bad`, 1),
		strings.Replace(valid, `"choose_chantry"`, `"play_card"`, 1),
		strings.Replace(valid, `"chantry_fire"`, `""`, 1),
		strings.Replace(valid, `"cardid"`, `"card"`, 1),
		strings.TrimSuffix(valid, "}") + `,"actorid":"member"}`,
	} {
		errorCode(t, f.call(400, "POST", path, "owner", body), "invalid_request")
	}
	before, _ := f.saved(view.ID)
	errorCode(t, f.call(409, "POST", path, "owner", commandBody(uuid.NewString(), 0, "choose_minion", "bastian_redhand")), "draft_sequence")
	errorCode(t, f.call(400, "POST", path, "owner", commandBody(uuid.NewString(), 0, "choose_chantry", "nonexistent")), "invalid_choice")
	after, _ := f.saved(view.ID)
	if !slices.Equal(before, after) || f.count("command_receipts") != 0 {
		t.Fatal("invalid choice persisted state")
	}
	accepted := f.call(200, "POST", path, "owner", valid).Body.String()
	f.choose(view.ID, "member", "choose_chantry", "chantry_air")
	if f.call(200, "POST", path, "owner", valid).Body.String() != accepted {
		t.Fatal("old retry did not replay original response")
	}
	errorCode(t, f.call(409, "POST", path, "owner", strings.Replace(valid, "chantry_fire", "chantry_earth", 1)), "request_id_conflict")
	errorCode(t, f.call(409, "POST", path, "owner", strings.Replace(valid, `"expectedrevision":0`, `"expectedrevision":2`, 1)), "request_id_conflict")
	errorCode(t, f.call(409, "POST", path, "owner", commandBody(uuid.NewString(), 0, "choose_victory_path", "subjugate_rival")), "stale_revision")
	errorCode(t, f.call(409, "POST", path, "owner", commandBody(uuid.NewString(), 2, "choose_chantry", "chantry_water")), "draft_sequence")
	f.choose(view.ID, "owner", "choose_victory_path", "subjugate_rival")
	errorCode(t, f.call(400, "POST", path, "owner", commandBody(uuid.NewString(), 3, "choose_minion", "nadia_rivermark")), "invalid_choice")
	// Whitespace, key ordering and UUID casing do not change canonical input.
	var reordered api.CommandRequest
	if err := json.Unmarshal([]byte(valid), &reordered); err != nil {
		t.Fatal(err)
	}
	body := fmt.Sprintf(`{ "type":%q, "cardid":%q, "expectedrevision":0, "requestid":%q }`, reordered.Type, reordered.CardID, strings.ToUpper(reordered.RequestID))
	if f.call(200, "POST", path, "owner", body).Body.String() != accepted {
		t.Fatal("equivalent retry conflicted")
	}
	// Another actor's identical request ID is independent and cannot replay owner data.
	body = commandBody(reordered.RequestID, 3, "choose_victory_path", "victory_council")
	other := matchView(t, f.call(200, "POST", path, "member", body))
	if other.You.Draft.Chantry != "chantry_air" || other.You.Draft.VictoryPath != "victory_council" {
		t.Fatal("receipt leaked another actor's view")
	}
	f.call(204, "POST", "/admin/match/"+view.ID+"/end", "admin", "")
	for _, actor := range []string{"owner", "outsider", "admin"} {
		f.call(404, "POST", path, actor, valid)
		f.call(404, "GET", "/match/"+view.ID+"/catalogue", actor, "")
	}
}

func TestDraftAtomicRollbackIncludingFinalStart(t *testing.T) {
	for _, trigger := range []string{
		`create trigger fail before update on matches begin select raise(abort, 'private snapshot failure'); end`,
		`create trigger fail before insert on command_receipts begin select raise(abort, 'private receipt failure'); end`,
	} {
		t.Run(trigger, func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(2)
			f.ready(id, 2)
			view, _, _ := f.launch(id)
			f.draft(view.ID, "owner", "chantry_fire", "subjugate_rival", "bastian_redhand")
			f.choose(view.ID, "member", "choose_chantry", "chantry_air")
			f.choose(view.ID, "member", "choose_victory_path", "victory_fame")
			before, state := f.saved(view.ID)
			expected, err := state.Choose(f.ids["member"], "choose_minion", "elsie_many_tongues")
			if err != nil {
				t.Fatal(err)
			}
			f.exec(trigger)
			body := commandBody(uuid.NewString(), 5, "choose_minion", "elsie_many_tongues")
			w := f.call(500, "POST", "/match/"+view.ID+"/commands", "member", body)
			if strings.Contains(w.Body.String(), "private") {
				t.Fatal("exposed internal failure")
			}
			after, _ := f.saved(view.ID)
			if !slices.Equal(before, after) || f.count("command_receipts") != 5 {
				t.Fatal("failed final command partially committed")
			}
			f.exec("drop trigger fail")
			f.call(200, "POST", "/match/"+view.ID+"/commands", "member", body)
			_, started := f.saved(view.ID)
			if !reflect.DeepEqual(expected, started) {
				t.Fatal("rollback changed random outcome or cards")
			}
		})
	}
}

func TestConcurrentDraftCommandsAndFinalRetries(t *testing.T) {
	for _, identical := range []bool{false, true} {
		t.Run(fmt.Sprint(identical), func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(2)
			f.ready(id, 2)
			view, _, _ := f.launch(id)
			f.draft(view.ID, "owner", "chantry_earth", "subjugate_rival", "rurik_ironvale")
			f.choose(view.ID, "member", "choose_chantry", "chantry_water")
			f.choose(view.ID, "member", "choose_victory_path", "victory_influence")
			body := commandBody(uuid.NewString(), 5, "choose_minion", "sister_agnes_reed")
			responses := make(chan *httptest.ResponseRecorder, 8)
			var wg sync.WaitGroup
			for i := 0; i < 8; i++ {
				input := body
				if !identical {
					input = commandBody(uuid.NewString(), 5, "choose_minion", "tamsin_lockward")
				}
				wg.Add(1)
				go func(input string) {
					defer wg.Done()
					responses <- f.request("POST", "/match/"+view.ID+"/commands", "member", input)
				}(input)
			}
			wg.Wait()
			close(responses)
			accepted := 0
			original := ""
			for w := range responses {
				if w.Code == 200 {
					accepted++
					if original != "" && original != w.Body.String() {
						t.Fatal("concurrent retries returned different outcomes")
					}
					original = w.Body.String()
				} else if !identical && w.Code == 409 {
					errorCode(t, w, "stale_revision")
				} else {
					t.Fatalf("unexpected concurrent response: %d %s", w.Code, w.Body.String())
				}
			}
			want := 1
			if identical {
				want = 8
			}
			_, saved := f.saved(view.ID)
			if accepted != want || saved.Revision != 6 || saved.Status != "playing" || f.count("command_receipts") != 6 {
				t.Fatal("concurrent commands overwrote or duplicated start")
			}
		})
	}
	// Different players compete on the same global revision; the loser refreshes.
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	f.ready(id, 2)
	view, _, _ := f.launch(id)
	responses := make(chan *httptest.ResponseRecorder, 2)
	var wg sync.WaitGroup
	for _, actor := range []string{"owner", "member"} {
		wg.Add(1)
		go func(actor string) {
			defer wg.Done()
			responses <- f.request("POST", "/match/"+view.ID+"/commands", actor, commandBody(uuid.NewString(), 0, "choose_chantry", "chantry_fire"))
		}(actor)
	}
	wg.Wait()
	close(responses)
	accepted, stale := 0, 0
	for w := range responses {
		if w.Code == 200 {
			accepted++
		} else if w.Code == 409 {
			stale++
			errorCode(t, w, "stale_revision")
		} else {
			t.Fatal(w.Code, w.Body.String())
		}
	}
	if accepted != 1 || stale != 1 || f.count("command_receipts") != 1 {
		t.Fatal("cross-player revisions did not serialize")
	}
	for _, actor := range []string{"owner", "member"} {
		if f.current(view.ID, actor).You.Draft.Step == "chantry" {
			f.choose(view.ID, actor, "choose_chantry", "chantry_fire")
		}
	}
	if f.current(view.ID, "owner").Revision != 2 {
		t.Fatal("stale player could not continue")
	}
}

func TestDraftUnsupportedSnapshotsAndReceiptRecovery(t *testing.T) {
	for _, change := range []func(*game.State){
		func(s *game.State) { s.SchemaVersion = 1 },
		func(s *game.State) { s.RulesVersion = "unknown" },
		func(s *game.State) { s.CatalogueVersion = "unknown" },
	} {
		f := newLifecycleFixture(t)
		id := f.createLobby(2)
		f.ready(id, 2)
		view, _, _ := f.launch(id)
		_, body, response := f.choose(view.ID, "owner", "choose_chantry", "chantry_fire")
		_, state := f.saved(view.ID)
		change(state)
		data, err := json.Marshal(state)
		if err != nil {
			t.Fatal(err)
		}
		f.exec("update matches set snapshot = ? where id = ?", string(data), view.ID)
		for _, path := range []string{"", "/catalogue"} {
			errorCode(t, f.call(409, "GET", "/match/"+view.ID+path, "owner", ""), "unsupported_state")
		}
		errorCode(t, f.call(409, "POST", "/match/"+view.ID+"/commands", "owner", commandBody(uuid.NewString(), 1, "choose_victory_path", "subjugate_rival")), "unsupported_state")
		if f.call(200, "POST", "/match/"+view.ID+"/commands", "owner", body).Body.String() != response {
			t.Fatal("unsupported state broke existing receipt")
		}
		var after []byte
		if err := f.db.QueryRow("select snapshot from matches where id = ?", view.ID).Scan(&after); err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(data, after) {
			t.Fatal("unsupported snapshot rewritten")
		}
		f.call(204, "POST", "/admin/match/"+view.ID+"/end", "admin", "")
		f.call(404, "POST", "/match/"+view.ID+"/commands", "owner", body)
	}
}
