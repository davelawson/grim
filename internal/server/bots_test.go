package server

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/bot"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/internal/testcatalogue"
	"github.com/google/uuid"
)

func (f *lifecycleFixture) addBot(id string) api.Lobby {
	f.t.Helper()
	w := f.call(201, "POST", api.AddBotToLobbyPath(id), "owner", "")
	var response api.GetLobbyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		f.t.Fatal(err)
	}
	return response.Lobby
}

func TestBotLobbyRoster(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	for _, actor := range []string{"member", "outsider", "admin"} {
		errorCode(t, f.call(403, "POST", api.AddBotToLobbyPath(id), actor, ""), "forbidden")
	}
	f.call(401, "POST", api.AddBotToLobbyPath(id), "", "")
	f.ready(id, 2)
	lobby := f.addBot(id)
	if len(lobby.Members) != 2 || len(lobby.Bots) != 1 || len(lobby.Readiness) != 2 {
		t.Fatal("bots changed human memberships")
	}
	botID := ""
	for key, b := range lobby.Bots {
		botID = key
		if b.Name != "Bot 1" {
			t.Fatal(b)
		}
	}
	for _, ready := range lobby.Readiness {
		if ready {
			t.Fatal("roster change did not clear readiness")
		}
	}
	errorCode(t, f.call(403, "DELETE", api.RemoveBotFromLobbyPath(id, botID), "member", ""), "forbidden")
	errorCode(t, f.call(404, "DELETE", api.RemoveBotFromLobbyPath(id, "missing"), "owner", ""), "not_found")
	f.call(400, "PUT", api.UpdateLobbyPath(id), "owner", fmt.Sprintf(`{"name":"Test","owner":%q}`, botID))
	lobby = f.addBot(id)
	if len(lobby.Bots) != 2 {
		t.Fatal("missing bots")
	}
	errorCode(t, f.call(409, "POST", api.AddBotToLobbyPath(id), "owner", ""), "lobby_full")
	errorCode(t, f.call(409, "POST", api.JoinLobbyPath(id), "third", ""), "lobby_full")
	if body := f.call(200, "GET", api.ListPublicLobbies.Route, "outsider", "").Body.String(); strings.Contains(body, id) {
		t.Fatal("full bot lobby listed publicly")
	}
	f.ready(id, 2)
	f.call(200, "DELETE", api.RemoveBotFromLobbyPath(id, botID), "owner", "")
	w := f.call(200, "GET", api.ListPublicLobbies.Route, "outsider", "")
	var public api.ListPublicLobbiesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &public); err != nil {
		t.Fatal(err)
	}
	if len(public.Lobbies) != 1 || public.Lobbies[0].MemberCount != 3 {
		t.Fatal(w.Body.String())
	}
	f.call(409, "POST", api.LaunchMatchPath(id), "owner", launchBody(uuid.NewString()))
	f.ready(id, 2)
	view, _, _ := f.launch(id)
	errorCode(t, f.call(409, "POST", api.AddBotToLobbyPath(id), "owner", ""), "lobby_closed")
	for key := range view.Bots {
		errorCode(t, f.call(409, "DELETE", api.RemoveBotFromLobbyPath(id, key), "owner", ""), "lobby_closed")
	}
	if f.count("users") != 6 || f.count("players") != 2 {
		t.Fatal("automated opponents created accounts or human memberships")
	}
}

func TestBotLaunchDraftAndRestart(t *testing.T) {
	for _, example := range []struct{ humans, bots int }{{1, 1}, {1, 2}, {1, 3}, {2, 1}, {2, 2}, {3, 1}} {
		t.Run(fmt.Sprintf("%dhumans/%dbots", example.humans, example.bots), func(t *testing.T) {
			f := newLifecycleFixture(t)
			id := f.createLobby(example.humans)
			for range example.bots {
				f.addBot(id)
			}
			f.ready(id, example.humans)
			view, requestID, body := f.launch(id)
			if len(view.Bots) != example.bots || len(view.Participants) != example.humans+example.bots || view.Revision != int64(3*example.bots) || view.Status != "setup" {
				t.Fatalf("bad automated setup: %+v", view)
			}
			for botID := range view.Bots {
				if !view.Wizards[botID].DraftComplete || view.Wizards[botID].InPlay != nil {
					t.Fatal("draft incomplete or private draft leaked")
				}
			}
			if retry := f.call(201, "POST", api.LaunchMatchPath(id), "owner", launchBody(requestID)); retry.Body.String() != body {
				t.Fatal("launch retry changed result")
			}
			f.restart()
			if !reflect.DeepEqual(view, f.current(view.ID, "owner")) {
				t.Fatal("restart changed setup")
			}
			for _, actor := range roster[:example.humans] {
				view = f.draft(view.ID, actor, "chantry_fire", "subjugate_rival", "bastian_redhand")
			}
			if view.Status != "playing" {
				t.Fatal("automatic draft did not start play")
			}
			assertHumanBoundary(t, view)
			before, _ := f.saved(view.ID)
			f.current(view.ID, "owner")
			f.restart()
			after, _ := f.saved(view.ID)
			if !slices.Equal(before, after) {
				t.Fatal("reads or restart changed saved state")
			}
			f.call(403, "GET", api.GetMatchPath(view.ID), "outsider", "")
		})
	}
}

func assertHumanBoundary(t *testing.T, view *game.View) {
	t.Helper()
	if view.Status == "finished" {
		return
	}
	actor := view.Turn.ActiveWizard
	if view.Choice != nil {
		actor = view.Choice.Actor
	} else if view.Exchange != nil {
		actor = view.Exchange.Priority
	}
	if _, ok := view.Bots[actor]; ok {
		t.Fatalf("request returned an automated input boundary: %+v", view)
	}
}

func TestBotConcurrentCapacity(t *testing.T) {
	f := newLifecycleFixture(t)
	id := f.createLobby(2)
	f.addBot(id)
	results := make(chan int, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); results <- f.request("POST", api.AddBotToLobbyPath(id), "owner", "").Code }()
	go func() { defer wg.Done(); results <- f.request("POST", api.JoinLobbyPath(id), "third", "").Code }()
	wg.Wait()
	close(results)
	success, conflict := 0, 0
	for status := range results {
		if status == 200 || status == 201 {
			success++
		} else if status == 409 {
			conflict++
		} else {
			t.Fatal(status)
		}
	}
	if success != 1 || conflict != 1 {
		t.Fatal("concurrent requests overfilled lobby")
	}
}

func TestBotCommandRollbackReplayAndStaleRevision(t *testing.T) {
	f := newLifecycleFixture(t)
	lobby := f.createLobby(1)
	f.addBot(lobby)
	f.ready(lobby, 1)
	view, _, _ := f.launch(lobby)
	f.choose(view.ID, "owner", "choose_chantry", "chantry_fire")
	view, _, _ = f.choose(view.ID, "owner", "choose_victory_path", "subjugate_rival")
	body := commandBody(uuid.NewString(), view.Revision, "choose_minion", "bastian_redhand")
	before, _ := f.saved(view.ID)
	receipts := f.count("command_receipts")
	f.exec(`create trigger fail_bot_save before update on matches begin select raise(abort,'private failure'); end`)
	f.call(500, "POST", api.SubmitDraftPath(view.ID), "owner", body)
	after, _ := f.saved(view.ID)
	if !slices.Equal(before, after) || f.count("command_receipts") != receipts {
		t.Fatal("automated draft failure did not roll back")
	}
	f.exec("drop trigger fail_bot_save")
	accepted := f.call(200, "POST", api.SubmitDraftPath(view.ID), "owner", body).Body.String()
	f.restart()
	if replay := f.call(200, "POST", api.SubmitDraftPath(view.ID), "owner", body); replay.Body.String() != accepted {
		t.Fatal("draft retry reran automation")
	}
	view = f.current(view.ID, "owner")
	stale := view.Revision - 1
	errorCode(t, f.call(409, "POST", api.EndTurnPath(view.ID), "owner", gameplayBody(t, api.EndTurnRequest{RequestID: uuid.NewString(), ExpectedRevision: &stale})), "stale_revision")
	// A human asset play opens Priority for the bot, ensuring automated inputs.
	_, state := f.saved(view.ID)
	human := f.ids["owner"]
	state.Turn.ActiveWizard, state.Turn.Phase, state.Turn.Step = human, "actions", "actions"
	state.Exchange, state.Resolving, state.PhaseEffects = nil, nil, nil
	state.Choice = nil
	state.Wizards[human].Resources = map[string]int{"wis": 10, "wealth": 10}
	state.Wizards[human].Hand = append(state.Wizards[human].Hand, api.CardInstance{ID: "bot-test-study", DesignID: "study", Owner: human, Location: "hand"})
	f.exec("update matches set snapshot = ? where id = ?", string(mustJSON(t, state)), state.ID)
	view = f.current(view.ID, "owner")
	c := game.Command{Type: "play_card", CardID: "bot-test-study"}
	c.RequestID = uuid.NewString()
	c.ExpectedRevision = &view.Revision
	path, input := typedGameplay(t, view, c)
	command := gameplayBody(t, input)
	before, _ = f.saved(view.ID)
	f.exec(`create trigger fail_bot_receipt before insert on command_receipts begin select raise(abort,'private failure'); end`)
	f.call(500, "POST", path, "owner", command)
	after, _ = f.saved(view.ID)
	if !slices.Equal(before, after) {
		t.Fatal("receipt failure saved automated moves")
	}
	f.exec("drop trigger fail_bot_receipt")
	responses := make(chan string, 4)
	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := f.request("POST", path, "owner", command)
			responses <- fmt.Sprintf("%d:%s", w.Code, w.Body.String())
		}()
	}
	wg.Wait()
	close(responses)
	first := ""
	for response := range responses {
		if !strings.HasPrefix(response, "200:") {
			t.Fatal(response)
		}
		if first == "" {
			first = response
		}
		if response != first {
			t.Fatal("concurrent retries differed")
		}
	}
	var completed api.MatchResponse
	if err := json.Unmarshal([]byte(strings.TrimPrefix(first, "200:")), &completed); err != nil {
		t.Fatal(err)
	}
	if completed.Match.Revision <= view.Revision+1 {
		t.Fatal("asset play did not include automated transitions")
	}
	f.restart()
	if replay := f.call(200, "POST", path, "owner", command); "200:"+replay.Body.String() != first {
		t.Fatal("restart lost gameplay receipt")
	}
}

func TestBotCompleteGames(t *testing.T) {
	for count := 1; count <= 3; count++ {
		t.Run(fmt.Sprintf("%dbots", count), func(t *testing.T) {
			f := newLifecycleFixture(t)
			lobby := f.createLobby(1)
			for range count {
				f.addBot(lobby)
			}
			f.ready(lobby, 1)
			view, _, _ := f.launch(lobby)
			f.pinCLISeed(view.ID)
			view = f.draft(view.ID, "owner", "chantry_water", "subjugate_rival", "tamsin_lockward")
			catalogue := testcatalogue.Load(t)
			for n := 0; n < 2000 && view.Status != "finished"; n++ {
				assertHumanBoundary(t, view)
				view = f.current(view.ID, "owner")
				c := game.Command{Type: "end_turn"}
				if view.Choice != nil {
					c = bot.Candidates(view, catalogue, 0)[0]
					if view.Choice.Type == "effects" {
						for _, option := range view.Choice.Options {
							if strings.HasPrefix(option, "skip:") {
								c = game.Command{Type: "resolve_choice", ChoiceID: view.Choice.ID, Selections: []string{option}}
								break
							}
						}
					}
				} else if view.Exchange != nil {
					c.Type = "pass"
				}
				view = f.command(view.ID, "owner", c)
			}
			if view.Status != "finished" || view.Outcome == nil || *view.Outcome != "victory" || len(view.Winners) != 1 || view.Winners[0] == f.ids["owner"] {
				t.Fatalf("automated game did not finish: turn=%+v", view.Turn)
			}
			f.restart()
			if !reflect.DeepEqual(view, f.current(view.ID, "owner")) {
				t.Fatal("completed result lost after restart")
			}
		})
	}
}

func TestInstallBotStorageIsAdditive(t *testing.T) {
	f := newLifecycleFixture(t)
	oldLobby := f.createLobby(2)
	f.ready(oldLobby, 2)
	oldMatch, _, _ := f.launch(oldLobby)
	before, _ := f.saved(oldMatch.ID)
	id := f.createLobby(1)
	f.exec("drop table lobby_bots")
	script, err := os.ReadFile("../../sql/install-lobby-bots.sql")
	if err != nil {
		t.Fatal(err)
	}
	f.exec(string(script))
	lobby := f.addBot(id)
	f.exec(string(script))
	if f.count("users") != 6 || f.count("lobbies") != 2 || f.count("lobby_bots") != 1 || len(lobby.Bots) != 1 || f.count("players") != 2 || f.count("launch_receipts") != 1 {
		t.Fatal("installation destroyed existing data")
	}
	after, _ := f.saved(oldMatch.ID)
	if !slices.Equal(before, after) || !reflect.DeepEqual(oldMatch, f.current(oldMatch.ID, "owner")) {
		t.Fatal("installation changed an existing match")
	}
}

func TestBotsFinishAfterLastHumanEliminated(t *testing.T) {
	f := newLifecycleFixture(t)
	lobby := f.createLobby(1)
	f.addBot(lobby)
	f.addBot(lobby)
	f.ready(lobby, 1)
	view, _, _ := f.launch(lobby)
	view = f.draft(view.ID, "owner", "chantry_fire", "subjugate_rival", "bastian_redhand")
	_, s := f.saved(view.ID)
	human := f.ids["owner"]
	ids := []string{}
	for id := range s.Bots {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	s.Ring = []string{human, ids[0], ids[1]}
	s.Turn = &api.Turn{Number: 10, ActiveWizard: human, Phase: "actions", Step: "actions"}
	s.Choice = nil
	s.Exchange = nil
	s.Resolving = nil
	s.PhaseEffects = nil
	for _, id := range s.Participants {
		w := s.Wizards[id]
		w.InPlay = w.InPlay[:1]
		w.Hand = []api.CardInstance{}
		w.DrawPile = []api.CardInstance{}
		w.Discard = []api.CardInstance{}
		w.Rumours = nil
		w.Reserved = nil
		w.Modifiers = nil
		w.Deferred = nil
		w.ImmuneUntil = 0
		w.Kneel = false
		w.AttackUsed = false
		w.Resources = map[string]int{}
	}
	s.Wizards[human].Integrity = 0
	winner := s.Wizards[ids[1]]
	winner.InPlay = append(winner.InPlay, api.CardInstance{ID: "strong-minion", DesignID: "bastian_redhand", Expertise: map[string]int{"might": 10}, Owner: ids[1], Location: "inplay"})
	winner.Hand = []api.CardInstance{{ID: "winning-attack", DesignID: "subjugate_rival", Owner: ids[1], Location: "hand"}}
	f.exec("update matches set snapshot = ? where id = ?", string(mustJSON(t, s)), s.ID)
	view = f.current(s.ID, "owner")
	for n := 0; n < 100 && view.Status != "finished"; n++ {
		c := game.Command{Type: "end_turn"}
		if view.Choice != nil {
			c = bot.Candidates(view, testcatalogue.Load(t), 0)[0]
		} else if view.Exchange != nil {
			c.Type = "pass"
		}
		view = f.command(view.ID, "owner", c)
	}
	if view.Status != "finished" || !view.Wizards[human].Eliminated || !slices.Equal(view.Winners, []string{ids[1]}) {
		t.Fatalf("bots did not finish after human elimination: %+v", view.Turn)
	}
}

func TestHumanVictoryStopsAutomation(t *testing.T) {
	f := newLifecycleFixture(t)
	lobby := f.createLobby(1)
	f.addBot(lobby)
	f.ready(lobby, 1)
	view, _, _ := f.launch(lobby)
	view = f.draft(view.ID, "owner", "chantry_fire", "subjugate_rival", "bastian_redhand")
	_, s := f.saved(view.ID)
	human := f.ids["owner"]
	s.Turn = &api.Turn{Number: 10, ActiveWizard: human, Phase: "actions", Step: "actions"}
	s.Choice = nil
	s.Exchange = nil
	s.Resolving = nil
	s.PhaseEffects = nil
	for id := range s.Bots {
		s.Wizards[id].Integrity = 0
		s.Wizards[id].ImmuneUntil = 0
		s.Wizards[id].InPlay = s.Wizards[id].InPlay[:1]
		s.Wizards[id].Hand = nil
	}
	w := s.Wizards[human]
	w.InPlay = append(w.InPlay, api.CardInstance{ID: "human-minion", DesignID: "bastian_redhand", Owner: human, Location: "inplay"})
	w.Hand = []api.CardInstance{{ID: "human-attack", DesignID: "subjugate_rival", Owner: human, Location: "hand"}, {ID: "human-focus", DesignID: "arcane_focus", Owner: human, Location: "hand"}}
	w.AttackUsed = false
	f.exec("update matches set snapshot = ? where id = ?", string(mustJSON(t, s)), s.ID)
	view = f.command(s.ID, "owner", game.Command{Type: "play_card", CardID: "human-attack", Minions: []string{"human-minion"}, Support: []api.CardSupport{{CardID: "human-focus"}}})
	for n := 0; n < 20 && view.Status != "finished"; n++ {
		view = f.command(view.ID, "owner", game.Command{Type: "pass"})
	}
	if view.Status != "finished" || !slices.Equal(view.Winners, []string{human}) {
		t.Fatal("human did not win")
	}
	before, _ := f.saved(view.ID)
	errorCode(t, f.call(409, "POST", api.EndTurnPath(view.ID), "owner", gameplayBody(t, api.EndTurnRequest{RequestID: uuid.NewString(), ExpectedRevision: &view.Revision})), "match_finished")
	after, _ := f.saved(view.ID)
	if !slices.Equal(before, after) {
		t.Fatal("automation changed a finished match")
	}
}
