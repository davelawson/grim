package server

import (
	"encoding/json"
	"fmt"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/internal/testcatalogue"
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

func draftResult(t *testing.T, response *httptest.ResponseRecorder) *api.DraftResponse {
	t.Helper()
	var draft api.DraftResponse
	if err := json.Unmarshal(response.Body.Bytes(), &draft); err != nil {
		t.Fatal(err)
	}
	if draft.Draft.Step == "" || draft.Cards == nil {
		t.Fatalf("missing draft progress or cards: %s", response.Body.String())
	}
	return &draft
}

func (f *lifecycleFixture) choose(id, actor, kind, card string) (*game.View, string, string) {
	f.t.Helper()
	requestID := uuid.NewString()
	current := draftResult(f.t, f.call(200, "GET", "/match/"+id+"/draft", actor, ""))
	body := commandBody(requestID, current.Revision, kind, card)
	w := f.call(200, "POST", "/match/"+id+"/draft", actor, body)
	result := draftResult(f.t, w)
	view := f.current(id, actor)
	if result.Draft != view.You.Draft || result.Revision != view.Revision || result.Status != view.Status {
		f.t.Fatal("draft submission differs from current participant progress")
	}
	return view, body, w.Body.String()
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

func assertDraftCards(t *testing.T, result *api.DraftResponse, step string, ids []string) {
	t.Helper()
	if result.Draft.Step != step || len(result.Cards) != len(ids) || result.Cards == nil {
		t.Fatalf("wrong %s options: %+v", step, result)
	}
	catalogue := testcatalogue.Load(t)
	for i, id := range ids {
		card, ok := catalogue.Card(id)
		if !ok || !reflect.DeepEqual(result.Cards[i], card) {
			t.Fatalf("expected full definition for %s, got %+v", id, result.Cards[i])
		}
	}
}

func TestDraftOptionsAreRevealedOneStepAtATime(t *testing.T) {
	for _, example := range []struct {
		chantry string
		minions []string
	}{
		{"chantry_fire", []string{"vera_coalheart", "bastian_redhand", "orla_nine_cinders"}},
		{"chantry_earth", []string{"hedda_stonewake", "leontine_grange", "rurik_ironvale"}},
		{"chantry_air", []string{"mina_cloudscript", "garrick_stormward", "elsie_many_tongues"}},
		{"chantry_water", []string{"nadia_rivermark", "sister_agnes_reed", "tamsin_lockward"}},
	} {
		t.Run(example.chantry, func(t *testing.T) {
			f := newLifecycleFixture(t)
			lobbyID := f.createLobby(2)
			f.ready(lobbyID, 2)
			launch, _, launchResponse := f.launch(lobbyID)
			path := api.GetDraftPath(launch.ID)
			before, _ := f.saved(launch.ID)
			initial := f.call(200, "GET", path, "owner", "")
			assertDraftCards(t, draftResult(t, initial), "chantry", []string{"chantry_fire", "chantry_earth", "chantry_air", "chantry_water"})
			after, _ := f.saved(launch.ID)
			if !slices.Equal(before, after) || f.count("command_receipts") != 0 {
				t.Fatal("opening draft changed saved state")
			}
			for _, response := range []string{
				launchResponse,
				f.call(200, "GET", api.GetMatchPath(launch.ID), "owner", "").Body.String(),
				f.call(200, "GET", api.GetCataloguePath(launch.ID), "owner", "").Body.String(),
			} {
				if strings.Contains(response, `"draftoffers"`) {
					t.Fatal("draft offers exposed outside draft endpoint")
				}
			}
			f.call(404, "POST", "/match/"+launch.ID+"/commands", "owner", commandBody(uuid.NewString(), 0, "choose_chantry", example.chantry))
			_, _, response := f.choose(launch.ID, "owner", "choose_chantry", example.chantry)
			var result api.DraftResponse
			if err := json.Unmarshal([]byte(response), &result); err != nil {
				t.Fatal(err)
			}
			assertDraftCards(t, &result, "victory_path", []string{"subjugate_rival"})
			f.choose(launch.ID, "owner", "choose_victory_path", "subjugate_rival")
			f.restart()
			options := draftResult(t, f.call(200, "GET", path, "owner", ""))
			assertDraftCards(t, options, "minion", example.minions)
			if options.Draft.Chantry != example.chantry || options.Draft.VictoryPath != "subjugate_rival" {
				t.Fatal("resumed draft lost selections")
			}
			// Another player's draft progress does not advance or reveal our picks.
			other := draftResult(t, f.call(200, "GET", path, "member", ""))
			assertDraftCards(t, other, "chantry", []string{"chantry_fire", "chantry_earth", "chantry_air", "chantry_water"})
			if other.Draft.Chantry != "" || other.Draft.VictoryPath != "" || other.Draft.Minion != "" {
				t.Fatal("another player's draft selections leaked")
			}
			completed, _, _ := f.choose(launch.ID, "owner", "choose_minion", example.minions[0])
			assertDraftCards(t, draftResult(t, f.call(200, "GET", path, "owner", "")), "complete", nil)
			if completed.Status != "setup" || !completed.Wizards[f.ids["owner"]].DraftComplete || completed.Wizards[f.ids["member"]].DraftComplete {
				t.Fatal("draft completion did not mark only the caller ready")
			}
			errorCode(t, f.call(409, "POST", path, "owner", commandBody(uuid.NewString(), completed.Revision, "choose_minion", example.minions[1])), "draft_sequence")
		})
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
			if catalogue.Catalogue.CatalogueVersion != launch.CatalogueVersion || len(catalogue.Catalogue.Cards) != 98 {
				t.Fatal("wrong pinned catalogue")
			}
			for _, actor := range roster[:size-1] {
				view := f.draft(launch.ID, actor, "chantry_fire", "subjugate_rival", "orla_nine_cinders")
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
			if view.Status != "playing" || view.Revision != int64(size*3) || view.Turn == nil || view.Turn.ActiveWizard != view.Ring[0] {
				t.Fatalf("bad started view %+v", view)
			}
			ring := slices.Clone(view.Ring)
			slices.Sort(ring)
			if !slices.Equal(ring, view.Participants) {
				t.Fatal("seating lost participants")
			}
			for id, public := range view.Wizards {
				handCount := 5
				if id == view.Turn.ActiveWizard {
					handCount++
				}
				if public.Integrity == nil || *public.Integrity != 3 || public.HandCount != handCount || public.DrawCount != 5 || len(public.InPlay) != 1 || !public.DraftComplete {
					t.Fatal("wrong public start state")
				}
			}
			if len(view.You.Hand) != view.Wizards[f.ids[actor]].HandCount || len(view.You.DrawPile) != 5 {
				t.Fatal("incorrect private start state")
			}
			// Opponent selections remain absent even when drafting completes.
			for _, hidden := range []string{`"random"`, `"seed"`, `"generator"`, `"schemaversion"`} {
				if strings.Contains(finalResponse, hidden) {
					t.Fatalf("leaked %s", hidden)
				}
			}
			if strings.Contains(finalResponse, "orla_nine_cinders") {
				t.Fatal("opponent drafted cards leaked")
			}
			before, state := f.saved(launch.ID)
			if len(state.Wizards) != size || reflect.DeepEqual(state.Random, initial.Random) || f.count("command_receipts") != size*3 {
				t.Fatal("bad persisted start")
			}
			finalMatch := f.call(200, "GET", "/match/"+view.ID, actor, "").Body.String()
			f.restart()
			if f.call(200, "GET", "/match/"+view.ID, actor, "").Body.String() != finalMatch {
				t.Fatal("restart changed participant view")
			}
			if f.call(200, "GET", "/match/"+view.ID+"/draft", actor, "").Body.String() != finalResponse {
				t.Fatal("restart changed completed draft")
			}
			if f.call(200, "POST", "/match/"+view.ID+"/draft", actor, finalBody).Body.String() != finalResponse {
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
			errorCode(t, f.call(409, "POST", "/match/"+view.ID+"/draft", actor, commandBody(uuid.NewString(), view.Revision, "choose_minion", "nadia_rivermark")), "draft_closed")
		})
	}
}

func TestDraftCommandValidationAuthorizationAndRetries(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	f.ready(id, 2)
	view, _, _ := f.launch(id)
	path := "/match/" + view.ID + "/draft"
	valid := commandBody(uuid.NewString(), 0, "choose_chantry", "chantry_fire")
	for _, actor := range []string{"", "invalid"} {
		f.call(401, "POST", path, actor, valid)
		f.call(401, "GET", path, actor, "")
		f.call(401, "GET", "/match/"+view.ID+"/catalogue", actor, "")
	}
	for _, actor := range []string{"outsider", "admin"} {
		f.call(403, "POST", path, actor, valid)
		f.call(403, "GET", path, actor, "")
		f.call(403, "GET", "/match/"+view.ID+"/catalogue", actor, "")
	}
	f.call(404, "POST", "/match/missing/draft", "owner", valid)
	f.call(404, "GET", "/match/missing/draft", "owner", "")
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
	var reordered api.DraftPickRequest
	if err := json.Unmarshal([]byte(valid), &reordered); err != nil {
		t.Fatal(err)
	}
	body := fmt.Sprintf(`{ "type":%q, "cardid":%q, "expectedrevision":0, "requestid":%q }`, reordered.Type, reordered.CardID, strings.ToUpper(reordered.RequestID))
	if f.call(200, "POST", path, "owner", body).Body.String() != accepted {
		t.Fatal("equivalent retry conflicted")
	}
	// Another actor's identical request ID is independent and cannot replay owner data.
	body = commandBody(reordered.RequestID, 3, "choose_victory_path", "subjugate_rival")
	other := draftResult(t, f.call(200, "POST", path, "member", body))
	if other.Draft.Chantry != "chantry_air" || other.Draft.VictoryPath != "subjugate_rival" {
		t.Fatal("receipt leaked another actor's view")
	}
	f.call(204, "POST", "/admin/match/"+view.ID+"/end", "admin", "")
	for _, actor := range []string{"owner", "outsider", "admin"} {
		f.call(404, "POST", path, actor, valid)
		f.call(404, "GET", path, actor, "")
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
			f.choose(view.ID, "member", "choose_victory_path", "subjugate_rival")
			before, state := f.saved(view.ID)
			expected, err := state.Choose(f.ids["member"], "choose_minion", "elsie_many_tongues")
			if err != nil {
				t.Fatal(err)
			}
			f.exec(trigger)
			body := commandBody(uuid.NewString(), 5, "choose_minion", "elsie_many_tongues")
			w := f.call(500, "POST", "/match/"+view.ID+"/draft", "member", body)
			if strings.Contains(w.Body.String(), "private") {
				t.Fatal("exposed internal failure")
			}
			after, _ := f.saved(view.ID)
			if !slices.Equal(before, after) || f.count("command_receipts") != 5 {
				t.Fatal("failed final command partially committed")
			}
			f.exec("drop trigger fail")
			f.call(200, "POST", "/match/"+view.ID+"/draft", "member", body)
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
			f.choose(view.ID, "member", "choose_victory_path", "subjugate_rival")
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
					responses <- f.request("POST", "/match/"+view.ID+"/draft", "member", input)
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
			responses <- f.request("POST", "/match/"+view.ID+"/draft", actor, commandBody(uuid.NewString(), 0, "choose_chantry", "chantry_fire"))
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
		for _, path := range []string{"", "/catalogue", "/draft"} {
			errorCode(t, f.call(409, "GET", "/match/"+view.ID+path, "owner", ""), "unsupported_state")
		}
		errorCode(t, f.call(409, "POST", "/match/"+view.ID+"/draft", "owner", commandBody(uuid.NewString(), 1, "choose_victory_path", "subjugate_rival")), "unsupported_state")
		if f.call(200, "POST", "/match/"+view.ID+"/draft", "owner", body).Body.String() != response {
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
		f.call(404, "POST", "/match/"+view.ID+"/draft", "owner", body)
	}
}

func TestLegacyDraftAndLaunchReceiptsRemainUsable(t *testing.T) {
	f := newLifecycleFixture(t)
	lobbyID := f.createLobby(2)
	f.ready(lobbyID, 2)
	view, launchID, launchResponse := f.launch(lobbyID)
	_, body, accepted := f.choose(view.ID, "owner", "choose_chantry", "chantry_fire")
	var request api.DraftPickRequest
	if err := json.Unmarshal([]byte(body), &request); err != nil {
		t.Fatal(err)
	}
	// Old receipts stored match views containing the complete shared menu.
	addLegacyOffers := func(data string) string {
		t.Helper()
		var fields map[string]json.RawMessage
		if err := json.Unmarshal([]byte(data), &fields); err != nil {
			t.Fatal(err)
		}
		fields["draftoffers"] = json.RawMessage(`{"chantries":["legacy-offer"],"victorypaths":[],"minions":{}}`)
		encoded, err := json.Marshal(fields)
		if err != nil {
			t.Fatal(err)
		}
		return string(encoded)
	}
	var commandReceipt, launchReceipt string
	if err := f.db.QueryRow("select response from command_receipts where request_id = ?", request.RequestID).Scan(&commandReceipt); err != nil {
		t.Fatal(err)
	}
	if err := f.db.QueryRow("select response from launch_receipts where request_id = ?", launchID).Scan(&launchReceipt); err != nil {
		t.Fatal(err)
	}
	commandReceipt, launchReceipt = addLegacyOffers(commandReceipt), addLegacyOffers(launchReceipt)
	f.exec("update command_receipts set response = ? where request_id = ?", commandReceipt, request.RequestID)
	f.exec("update launch_receipts set response = ? where request_id = ?", launchReceipt, launchID)
	f.choose(view.ID, "owner", "choose_victory_path", "subjugate_rival")
	before, _ := f.saved(view.ID)
	f.restart()
	if response := f.call(200, "POST", api.SubmitDraftPath(view.ID), "owner", body).Body.String(); response != accepted {
		t.Fatal("legacy command receipt did not replay original step and options")
	}
	if response := f.call(201, "POST", api.LaunchMatchPath(lobbyID), "owner", launchBody(launchID)).Body.String(); response != launchResponse || strings.Contains(response, "draftoffers") {
		t.Fatal("legacy launch receipt exposed draft offers or changed participant view")
	}
	var savedCommand, savedLaunch string
	if err := f.db.QueryRow("select response from command_receipts where request_id = ?", request.RequestID).Scan(&savedCommand); err != nil {
		t.Fatal(err)
	}
	if err := f.db.QueryRow("select response from launch_receipts where request_id = ?", launchID).Scan(&savedLaunch); err != nil {
		t.Fatal(err)
	}
	after, _ := f.saved(view.ID)
	if savedCommand != commandReceipt || savedLaunch != launchReceipt || !slices.Equal(before, after) || f.count("command_receipts") != 2 {
		t.Fatal("legacy receipt replay rewrote receipts or draft state")
	}
	assertDraftCards(t, draftResult(t, f.call(200, "GET", api.GetDraftPath(view.ID), "owner", "")), "minion", []string{"vera_coalheart", "bastian_redhand", "orla_nine_cinders"})
}
