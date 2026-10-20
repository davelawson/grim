package server

import (
	"encoding/json"
	"fmt"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"github.com/google/uuid"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestTypedGameplaySchemas(t *testing.T) {
	f, view := reactionFixture(t)
	revision := view.Revision
	no := false
	cases := []struct {
		name     string
		path     string
		input    any
		required []string
	}{
		{"PlayCard", api.PlayCardPath(view.ID), api.PlayCardRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance"}, []string{"requestid", "expectedrevision", "cardid"}},
		{"Attack", api.AttackPath(view.ID), api.AttackRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance", Minions: []string{}}, []string{"requestid", "expectedrevision", "cardid", "minions"}},
		{"AttemptChallenge", api.AttemptChallengePath(view.ID), api.AttemptChallengeRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance", ChallengeID: "instance"}, []string{"requestid", "expectedrevision", "cardid", "challengeid"}},
		{"ActivateAbility", api.ActivateAbilityPath(view.ID), api.ActivateAbilityRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance", AbilityID: "instance"}, []string{"requestid", "expectedrevision", "cardid", "abilityid"}},
		{"ResearchElement", api.ResearchElementPath(view.ID), api.ResearchElementRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance", Element: "fire"}, []string{"requestid", "expectedrevision", "cardid", "element"}},
		{"ResearchSpell", api.ResearchSpellPath(view.ID), api.ResearchSpellRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance"}, []string{"requestid", "expectedrevision", "cardid"}},
		{"HealWithFocus", api.HealWithFocusPath(view.ID), api.HealWithFocusRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance"}, []string{"requestid", "expectedrevision", "cardid"}},
		{"EndTurn", api.EndTurnPath(view.ID), api.EndTurnRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision}, []string{"requestid", "expectedrevision"}},
		{"Maintenance", api.MaintenancePath(view.ID, "choice"), api.MaintenanceRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, Keep: []string{}}, []string{"requestid", "expectedrevision", "keep"}},
		{"Kneel", api.KneelPath(view.ID, "choice"), api.KneelRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, Kneel: &no}, []string{"requestid", "expectedrevision", "kneel"}},
		{"ResolveEffect", api.ResolveEffectPath(view.ID, "choice"), api.ResolveEffectRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, EffectID: "instance"}, []string{"requestid", "expectedrevision", "effectid"}},
		{"SkipEffect", api.SkipEffectPath(view.ID, "choice"), api.SkipEffectRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, EffectID: "instance"}, []string{"requestid", "expectedrevision", "effectid"}},
		{"UndertakeHunt", api.UndertakeHuntPath(view.ID, "choice"), api.UndertakeHuntRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, EffectID: "instance"}, []string{"requestid", "expectedrevision", "effectid"}},
		{"ChooseResearch", api.ChooseResearchPath(view.ID, "choice"), api.ChooseResearchRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, DesignID: "instance"}, []string{"requestid", "expectedrevision", "designid"}},
		{"EndDiscard", api.EndDiscardPath(view.ID, "choice"), api.EndDiscardRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, Cards: []string{}}, []string{"requestid", "expectedrevision", "cards"}},
		{"HandDiscard", api.HandDiscardPath(view.ID, "choice"), api.HandDiscardRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance", Disposition: "discard"}, []string{"requestid", "expectedrevision", "cardid", "disposition"}},
		{"AssetLoss", api.AssetLossPath(view.ID, "choice"), api.AssetLossRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance"}, []string{"requestid", "expectedrevision", "cardid"}},
		{"Scry", api.ScryPath(view.ID, "choice"), api.ScryRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, Take: "instance", Order: []string{}}, []string{"requestid", "expectedrevision", "take", "order"}},
		{"React", api.ReactPath(view.ID), api.ReactRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "instance"}, []string{"requestid", "expectedrevision", "cardid"}},
		{"Defend", api.DefendPath(view.ID), api.DefendRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision}, []string{"requestid", "expectedrevision"}},
		{"Pass", api.PassPath(view.ID), api.PassRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision}, []string{"requestid", "expectedrevision"}},
	}
	before, _ := f.saved(view.ID)
	count := f.count("command_receipts")
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := gameplayBody(t, tc.input)
			f.call(404, "POST", strings.Replace(tc.path, view.ID, "missing", 1), "owner", body)
			for _, actor := range []string{"", "invalid"} {
				f.call(401, "POST", tc.path, actor, body)
			}
			for _, actor := range []string{"outsider", "admin"} {
				f.call(403, "POST", tc.path, actor, body)
			}
			for _, bad := range []string{"null", "{}", "[]", body + body, body + " trailing"} {
				errorCode(t, f.call(400, "POST", tc.path, "owner", bad), "invalid_request")
			}
			for _, field := range []string{"type", "mode", "choiceid", "unknown"} {
				errorCode(t, f.call(400, "POST", tc.path, "owner", strings.TrimSuffix(body, "}")+fmt.Sprintf(",%q:null}", field)), "invalid_request")
			}
			for _, field := range tc.required {
				for _, null := range []bool{false, true} {
					var input map[string]any
					_ = json.Unmarshal([]byte(body), &input)
					if null {
						input[field] = nil
					} else {
						delete(input, field)
					}
					errorCode(t, f.call(400, "POST", tc.path, "owner", gameplayBody(t, input)), "invalid_request")
				}
			}
			var input map[string]any
			_ = json.Unmarshal([]byte(body), &input)
			input["expectedrevision"] = -1
			errorCode(t, f.call(400, "POST", tc.path, "owner", gameplayBody(t, input)), "invalid_request")
		})
	}
	after, _ := f.saved(view.ID)
	if !slices.Equal(before, after) || count != f.count("command_receipts") {
		t.Fatal("schema rejection mutated match")
	}
	for _, path := range []string{"/match/" + view.ID + "/commands", "/match/" + view.ID + "/reactions"} {
		f.call(404, "POST", path, "owner", "{}")
	}
}

func TestTypedActionIntentAndResearchElements(t *testing.T) {
	f, view := reactionFixture(t)
	id := view.ID
	_, state := f.saved(id)
	owner := f.ids["owner"]
	w := state.Wizards[owner]
	for _, pair := range [][2]string{{"research", "research_spell"}, {"heal-focus", "arcane_focus"}, {"focus-two", "arcane_focus"}, {"element", "element_fire"}} {
		w.Hand = append(w.Hand, api.CardInstance{ID: "test/" + pair[0], DesignID: pair[1], Owner: owner, Location: "hand"})
	}
	f.exec("update matches set snapshot = ? where id = ?", string(mustJSON(t, state)), id)
	revision := view.Revision
	before, _ := f.saved(id)
	invalid := []struct {
		path  string
		input any
	}{
		{api.PlayCardPath(id), api.PlayCardRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/attack"}},
		{api.PlayCardPath(id), api.PlayCardRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/research"}},
		{api.PlayCardPath(id), api.PlayCardRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/heal-focus"}},
		{api.AttackPath(id), api.AttackRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/study", Minions: []string{"test/guard-one"}}},
		{api.ResearchSpellPath(id), api.ResearchSpellRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/research", Elements: []string{"test/heal-focus"}}},
		{api.ResearchElementPath(id), api.ResearchElementRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/study", Element: "fire"}},
		{api.HealWithFocusPath(id), api.HealWithFocusRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/study"}},
	}
	for _, tc := range invalid {
		errorCode(t, f.call(400, "POST", tc.path, "owner", gameplayBody(t, tc.input)), "invalid_choice")
	}
	after, _ := f.saved(id)
	if !slices.Equal(before, after) {
		t.Fatal("intent rejection mutated match")
	}
	// A second Focus explicitly contributes toward the healing cost.
	heal := api.HealWithFocusRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/heal-focus", Support: []api.CardSupport{{CardID: "test/focus-two", Domains: []string{"fire"}}}}
	healed := matchView(t, f.call(200, "POST", api.HealWithFocusPath(id), "owner", gameplayBody(t, heal)))
	if healed.Exchange.Stack[0].Mode != "heal" || healed.Wizards[owner].Resources["wis"] != 9 {
		t.Fatal("support-assisted healing lost cost contribution")
	}
	// Reset the fixture to exercise supporting Element instance IDs.
	f.exec("update matches set snapshot = ?, revision = ? where id = ?", string(before), revision, id)
	research := api.ResearchSpellRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/research", Elements: []string{"test/element"}}
	researched := matchView(t, f.call(200, "POST", api.ResearchSpellPath(id), "owner", gameplayBody(t, research)))
	effect := researched.Exchange.Stack[0]
	if effect.Mode != "discover" || effect.Element != "fire" || len(effect.Support) != 1 || effect.Support[0].CardID != "test/element" {
		t.Fatal("Research Elements not translated")
	}
}

func TestTypedReceiptsReplayRollbackConcurrencyAndTerminal(t *testing.T) {
	f, view := reactionFixture(t)
	id := view.ID
	rev := view.Revision
	action := api.PlayCardRequest{RequestID: uuid.NewString(), ExpectedRevision: &rev, CardID: "test/study"}
	actionBody := gameplayBody(t, action)
	accepted := f.call(200, "POST", api.PlayCardPath(id), "owner", actionBody)
	view = matchView(t, accepted)
	rev = view.Revision
	reaction := api.ReactRequest{RequestID: uuid.NewString(), ExpectedRevision: &rev, CardID: "test/counter", Targets: []string{view.Exchange.Stack[0].ID}}
	body := gameplayBody(t, reaction)
	path := api.ReactPath(id)
	before, _ := f.saved(id)
	count := f.count("command_receipts")
	f.exec(`create trigger fail before insert on command_receipts begin select raise(abort, 'private reaction failure'); end`)
	failure := f.call(500, "POST", path, "member", body)
	if strings.Contains(failure.Body.String(), "private reaction failure") {
		t.Fatal("internal error leaked")
	}
	after, _ := f.saved(id)
	if !slices.Equal(before, after) || count != f.count("command_receipts") {
		t.Fatal("failed transition committed")
	}
	f.exec("drop trigger fail")
	responses := make(chan string, 6)
	var wg sync.WaitGroup
	for i := range 6 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			input := body
			if i%2 == 1 {
				input = strings.ReplaceAll(input, reaction.RequestID, strings.ToUpper(reaction.RequestID))
			}
			w := f.request("POST", path, "member", input)
			responses <- fmt.Sprintf("%d:%s", w.Code, w.Body.String())
		}(i)
	}
	wg.Wait()
	close(responses)
	recorded := ""
	for response := range responses {
		if !strings.HasPrefix(response, "200:") {
			t.Fatal(response)
		}
		if recorded != "" && recorded != response {
			t.Fatal("concurrent replay diverged")
		}
		recorded = response
	}
	if count+1 != f.count("command_receipts") {
		t.Fatal("replay created receipts")
	}
	errorCode(t, f.call(409, "POST", api.PlayCardPath(id), "member", body), "request_id_conflict")
	errorCode(t, f.call(409, "POST", api.ReactPath(id), "owner", actionBody), "request_id_conflict")
	changed := reaction
	changed.Targets = []string{"changed"}
	errorCode(t, f.call(409, "POST", path, "member", gameplayBody(t, changed)), "request_id_conflict")
	changed = reaction
	newRev := rev + 1
	changed.ExpectedRevision = &newRev
	errorCode(t, f.call(409, "POST", path, "member", gameplayBody(t, changed)), "request_id_conflict")
	changed = reaction
	changed.RequestID = uuid.NewString()
	errorCode(t, f.call(409, "POST", path, "member", gameplayBody(t, changed)), "stale_revision")
	// Counter the Counterspell, then let Study pause for a required discard.
	view = f.current(id, "owner")
	view = f.command(id, "owner", game.Command{Type: "play_card", CardID: "test/counter-owner", Targets: []string{view.Exchange.Stack[1].ID}})
	for view.Choice == nil {
		view = f.command(id, f.actor(view.Exchange.Priority), game.Command{Type: "pass"})
	}
	if view.Choice.Type != "hand_discard" {
		t.Fatal("missing nested resolution")
	}
	f.restart()
	if got := f.call(200, "POST", path, "member", body).Body.String(); got != strings.TrimPrefix(recorded, "200:") {
		t.Fatal("replay checked live choice/Priority")
	}
	// A terminal state must not affect an accepted receipt.
	_, state := f.saved(id)
	state.Status = "finished"
	state.Exchange = nil
	state.Resolving = nil
	state.Choice = nil
	outcome := "victory"
	state.Outcome = &outcome
	state.Winners = []string{f.ids["owner"]}
	f.exec("update matches set status = ?, outcome = ?, snapshot = ? where id = ?", state.Status, outcome, string(mustJSON(t, state)), id)
	f.restart()
	if got := f.call(200, "POST", path, "member", body).Body.String(); got != strings.TrimPrefix(recorded, "200:") {
		t.Fatal("terminal replay changed")
	}
	currentRev := state.Revision
	errorCode(t, f.call(409, "POST", api.EndTurnPath(id), "owner", gameplayBody(t, api.EndTurnRequest{RequestID: uuid.NewString(), ExpectedRevision: &currentRev})), "match_finished")
	// Historical broad receipts keep their UUID reserved and are never reinterpreted.
	f.exec("update command_receipts set command = ? where match_id = ? and actor_id = ? and request_id = ?", `{"operation":"reaction","command":{}}`, id, f.ids["member"], reaction.RequestID)
	errorCode(t, f.call(409, "POST", path, "member", body), "request_id_conflict")
	var saved string
	_ = f.db.QueryRow("select command from command_receipts where match_id=? and actor_id=? and request_id=?", id, f.ids["member"], reaction.RequestID).Scan(&saved)
	if saved != `{"operation":"reaction","command":{}}` {
		t.Fatal("historical receipt rewritten")
	}
	// Receipt recovery does not require decoding the latest snapshot.
	f.exec("update matches set snapshot = ? where id = ?", `{"schemaversion":999}`, id)
	if f.call(200, "POST", api.PlayCardPath(id), "owner", actionBody).Body.String() != accepted.Body.String() {
		t.Fatal("receipt read current snapshot")
	}
}

func TestTypedChoiceIdentityIntentAndLiteralPaths(t *testing.T) {
	f, view := reactionFixture(t)
	id := view.ID
	_, state := f.saved(id)
	owner := f.ids["owner"]
	choiceID := "choice/:id + %雪?#/end"
	state.Choice = &api.RequiredChoice{ID: choiceID, Actor: owner, Type: "kneel", Options: []string{"kneel", "decline"}, Minimum: 1, Maximum: 1}
	state.Turn.Phase = "start"
	state.Turn.Step = "bend_the_knee"
	f.exec("update matches set snapshot = ? where id = ?", string(mustJSON(t, state)), id)
	rev := state.Revision
	no := false
	input := api.KneelRequest{RequestID: uuid.NewString(), ExpectedRevision: &rev, Kneel: &no}
	path := api.KneelPath(id, choiceID)
	body := gameplayBody(t, input)
	errorCode(t, f.call(400, "POST", api.MaintenancePath(id, choiceID), "owner", gameplayBody(t, api.MaintenanceRequest{RequestID: uuid.NewString(), ExpectedRevision: &rev, Keep: []string{}})), "invalid_choice")
	errorCode(t, f.call(409, "POST", path, "member", body), "wrong_priority")
	errorCode(t, f.call(409, "POST", api.KneelPath(id, choiceID+"stale"), "owner", body), "stale_choice")
	accepted := f.call(200, "POST", path, "owner", body)
	f.restart()
	if f.call(200, "POST", path, "owner", body).Body.String() != accepted.Body.String() {
		t.Fatal("choice replay changed")
	}
	errorCode(t, f.call(409, "POST", api.KneelPath(id, choiceID+"other"), "owner", body), "request_id_conflict")
	yes := true
	changed := input
	changed.Kneel = &yes
	errorCode(t, f.call(409, "POST", path, "owner", gameplayBody(t, changed)), "request_id_conflict")
	errorCode(t, f.call(409, "POST", api.EndDiscardPath(id, choiceID), "owner", gameplayBody(t, api.EndDiscardRequest{RequestID: input.RequestID, ExpectedRevision: &rev, Cards: []string{}})), "request_id_conflict")
}

func TestTypedEffectRoutes(t *testing.T) {
	f, view := reactionFixture(t)
	id := view.ID
	_, state := f.saved(id)
	owner := f.ids["owner"]
	state.PhaseEffects = []game.PhaseEffect{{ID: "hunt", Actor: owner, Source: "moon", Kind: "moonlit_hunt", Optional: true}, {ID: "ordinary", Actor: owner, Source: "source", Kind: "income", Optional: true}}
	state.Choice = &api.RequiredChoice{ID: "effects/one", Actor: owner, Type: "effects", Options: []string{"hunt", "skip:hunt", "ordinary", "skip:ordinary"}, Minimum: 1, Maximum: 1}
	state.Turn.Phase = "start"
	state.Turn.Step = "effects"
	f.exec("update matches set snapshot=? where id=?", string(mustJSON(t, state)), id)
	rev := state.Revision
	before, _ := f.saved(id)
	for _, tc := range []struct {
		path  string
		input any
	}{
		{api.ResolveEffectPath(id, state.Choice.ID), api.ResolveEffectRequest{RequestID: uuid.NewString(), ExpectedRevision: &rev, EffectID: "hunt"}},
		{api.UndertakeHuntPath(id, state.Choice.ID), api.UndertakeHuntRequest{RequestID: uuid.NewString(), ExpectedRevision: &rev, EffectID: "ordinary"}},
	} {
		errorCode(t, f.call(400, "POST", tc.path, "owner", gameplayBody(t, tc.input)), "invalid_choice")
	}
	after, _ := f.saved(id)
	if !slices.Equal(before, after) {
		t.Fatal("effect intent rejection mutated state")
	}
	// A skip option is not a legal effect ID on either resolution route.
	for _, path := range []string{api.ResolveEffectPath(id, state.Choice.ID), api.UndertakeHuntPath(id, state.Choice.ID)} {
		errorCode(t, f.call(400, "POST", path, "owner", gameplayBody(t, map[string]any{"requestid": uuid.NewString(), "expectedrevision": rev, "effectid": "skip:hunt"})), "invalid_choice")
	}
	// Skipping a hunt never accepts commitments and uses the unprefixed effect ID.
	request := api.SkipEffectRequest{RequestID: uuid.NewString(), ExpectedRevision: &rev, EffectID: "hunt"}
	body := gameplayBody(t, request)
	errorCode(t, f.call(400, "POST", api.SkipEffectPath(id, state.Choice.ID), "owner", strings.TrimSuffix(body, "}")+`,"minions":[]}`), "invalid_request")
	f.call(200, "POST", api.SkipEffectPath(id, state.Choice.ID), "owner", body)
}

// Each controller must preserve the existing engine transition, including its
// randomness and continuation, and replay after that input boundary has moved.
func TestTypedRoutesApplyEngineTransitions(t *testing.T) {
	for _, operation := range []string{
		"actions/play-card", "actions/attack", "actions/attempt-challenge", "actions/activate-ability", "actions/research-element", "actions/research-spell", "actions/heal-with-focus", "actions/end-turn",
		"choices/maintenance", "choices/kneel", "choices/resolve-effect", "choices/skip-effect", "choices/undertake-hunt", "choices/research", "choices/end-discard", "choices/hand-discard", "choices/asset-loss", "choices/attack-loss", "choices/scry",
		"reactions/play-card", "reactions/defend", "reactions/pass",
	} {
		t.Run(operation, func(t *testing.T) {
			f, view := reactionFixture(t)
			id := view.ID
			_, state := f.saved(id)
			owner := f.ids["owner"]
			actor := "owner"
			w := state.Wizards[owner]
			for _, pair := range [][2]string{{"investigate", "investigate"}, {"focus", "arcane_focus"}, {"research", "research_spell"}, {"element", "element_fire"}} {
				w.Hand = append(w.Hand, api.CardInstance{ID: "typed/" + pair[0], DesignID: pair[1], Owner: owner, Location: "hand"})
			}
			w.InPlay = append(w.InPlay, api.CardInstance{ID: "typed/spring", DesignID: "spell_healing_spring", Owner: owner, Location: "inplay"})
			c := game.Command{}
			if strings.HasPrefix(operation, "choices/") {
				kind := strings.TrimPrefix(operation, "choices/")
				kind = mapChoiceKind(kind)
				state.Choice = &api.RequiredChoice{ID: "typed/choice + %?#雪", Actor: owner, Type: kind, Minimum: 1, Maximum: 1}
				q := state.Choice
				c.Type = "resolve_choice"
				c.ChoiceID = q.ID
				switch kind {
				case "maintenance":
					q.Minimum = 0
					q.Maximum = 2
					q.Options = []string{"test/guard-one", "test/guard-two"}
					c.Selections = []string{}
					state.Turn.Phase = "start"
					state.Turn.Step = "maintenance"
				case "kneel":
					q.Options = []string{"kneel", "decline"}
					c.Selections = []string{"decline"}
					state.Turn.Phase = "start"
					state.Turn.Step = "bend_the_knee"
				case "effects":
					effect := game.PhaseEffect{ID: "typed/effect", Actor: owner, Source: "typed/spring", Kind: "whispering_wind_trigger", Optional: true}
					if operation == "choices/undertake-hunt" {
						effect.Source = "typed/moon"
						effect.Kind = "moonlit_hunt"
						w.Rumours = []api.RumourView{{ID: "typed/moon", DesignID: "rumour_full_moon", ExpiresAfterTurn: 10}}
						c.Minions = []string{"test/guard-one"}
					}
					state.PhaseEffects = []game.PhaseEffect{effect}
					q.Effects = []api.EffectOption{{ID: effect.ID, Kind: effect.Kind, Source: effect.Source, Optional: effect.Optional}}
					q.Options = []string{effect.ID, "skip:" + effect.ID}
					c.Selections = []string{effect.ID}
					if operation == "choices/skip-effect" {
						c.Selections[0] = "skip:" + effect.ID
					}
					state.Turn.Phase = "start"
					state.Turn.Step = "effects"
				case "end_discard":
					w.Hand = w.Hand[:1]
					q.Minimum = 0
					q.Options = []string{w.Hand[0].ID}
					c.Selections = []string{}
					state.Turn.Phase = "end"
					state.Turn.Step = "discard"
				default:
					stage, kind := "research", "research_spell"
					switch q.Type {
					case "hand_discard":
						stage, kind = "discard", "flash_of_insight"
					case "asset_discard", "attack_loss":
						stage, kind = "attack_loss", "subjugate_rival"
					case "scry":
						stage, kind = "scry", "scrying_pool"
					}
					state.Resolving = &game.PendingEffect{PendingEffectView: api.PendingEffectView{ID: "typed/resolving", Actor: owner, Kind: kind, Cards: []api.CardInstance{}}, Stage: stage}
					state.Exchange = &game.Exchange{Participants: slices.Clone(state.Ring), Stack: []*game.PendingEffect{}}
					switch q.Type {
					case "research":
						q.Options = []string{"flame_lash"}
						c.Selections = q.Options
					case "hand_discard":
						q.Options = []string{"test/study"}
						c.Selections = q.Options
						c.Mode = "destroy"
					case "asset_discard", "attack_loss":
						q.Options = []string{"test/guard-one"}
						c.Selections = q.Options
					case "scry":
						w.DrawPile = []api.CardInstance{{ID: "typed/top", DesignID: "basic_wis", Owner: owner, Location: "drawpile"}}
						q.Options = []string{"typed/top"}
						c.Selections = q.Options
						c.Order = []string{}
					}
				}
			} else if strings.HasPrefix(operation, "reactions/") {
				declaration := game.Command{Type: "play_card", CardID: "test/study"}
				if operation == "reactions/defend" {
					declaration.CardID = "test/attack"
					declaration.Minions = []string{"test/guard-one"}
				}
				var err error
				state, err = state.Apply(owner, declaration)
				if err != nil {
					t.Fatal(err)
				}
				actor = "member"
				switch operation {
				case "reactions/play-card":
					c = game.Command{Type: "play_card", CardID: "test/counter", Targets: []string{state.Exchange.Stack[0].ID}}
				case "reactions/defend":
					c = game.Command{Type: "defend", Minions: []string{"test/guard-member"}, Support: []api.CardSupport{{CardID: "test/focus"}}}
				case "reactions/pass":
					c.Type = "pass"
				}
			} else {
				switch operation {
				case "actions/play-card":
					c = game.Command{Type: "play_card", CardID: "test/study"}
				case "actions/attack":
					c = game.Command{Type: "play_card", CardID: "test/attack", Minions: []string{"test/guard-one"}}
				case "actions/attempt-challenge":
					c = game.Command{Type: "attempt_challenge", CardID: "typed/investigate", ChallengeID: "investigate", Minions: []string{"test/guard-one"}}
				case "actions/activate-ability":
					c = game.Command{Type: "activate_ability", CardID: "typed/spring", AbilityID: "heal"}
				case "actions/research-element":
					c = game.Command{Type: "play_card", CardID: "typed/research", Mode: "element", Element: "water"}
				case "actions/research-spell":
					c = game.Command{Type: "play_card", CardID: "typed/research", Mode: "discover", Support: []api.CardSupport{{CardID: "typed/element"}}}
				case "actions/heal-with-focus":
					c = game.Command{Type: "play_card", CardID: "typed/focus", Mode: "heal"}
				case "actions/end-turn":
					c.Type = "end_turn"
				}
			}
			revision := state.Revision
			c.RequestID = uuid.NewString()
			c.ExpectedRevision = &revision
			expected, err := state.Apply(f.ids[actor], c)
			if err != nil {
				t.Fatal(err)
			}
			f.exec("update matches set snapshot=?, revision=? where id=?", string(mustJSON(t, state)), state.Revision, id)
			current := f.current(id, actor)
			path, input := typedGameplay(t, current, c)
			suffix := strings.Split(operation, "/")[1]
			if suffix == "attack-loss" {
				suffix = "asset-loss"
			}
			if !strings.HasSuffix(path, "/"+suffix) {
				t.Fatal("strategy selected wrong typed route", path)
			}
			body := gameplayBody(t, input)
			accepted := f.call(200, "POST", path, actor, body)
			got := matchView(t, accepted)
			if !slices.Equal(mustJSON(t, got), mustJSON(t, expected.View(f.ids[actor]))) {
				t.Fatal("typed request changed projected transition")
			}
			_, saved := f.saved(id)
			if !slices.Equal(mustJSON(t, saved), mustJSON(t, expected)) {
				t.Fatal("typed request changed persisted transition")
			}
			count := f.count("command_receipts")
			f.restart()
			if f.call(200, "POST", path, actor, body).Body.String() != accepted.Body.String() || f.count("command_receipts") != count {
				t.Fatal("typed input replay changed response or receipts")
			}
		})
	}
}

func mapChoiceKind(kind string) string {
	switch kind {
	case "resolve-effect", "skip-effect", "undertake-hunt":
		return "effects"
	case "end-discard":
		return "end_discard"
	case "hand-discard":
		return "hand_discard"
	case "asset-loss":
		return "asset_discard"
	case "attack-loss":
		return "attack_loss"
	default:
		return kind
	}
}

func TestTypedCompetingReactions(t *testing.T) {
	f, view := reactionFixture(t)
	id := view.ID
	view = f.command(id, "owner", game.Command{Type: "play_card", CardID: "test/study"})
	revision := view.Revision
	count := f.count("command_receipts")
	bodies := []string{}
	for range 2 {
		bodies = append(bodies, gameplayBody(t, api.ReactRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/counter", Targets: []string{view.Exchange.Stack[0].ID}}))
	}
	results := make(chan int, 2)
	var wg sync.WaitGroup
	for _, body := range bodies {
		wg.Add(1)
		go func(body string) {
			defer wg.Done()
			results <- f.request("POST", api.ReactPath(id), "member", body).Code
		}(body)
	}
	wg.Wait()
	close(results)
	statuses := []int{}
	for status := range results {
		statuses = append(statuses, status)
	}
	slices.Sort(statuses)
	if !slices.Equal(statuses, []int{200, 409}) || f.count("command_receipts") != count+1 {
		t.Fatal("competing submissions committed twice", statuses)
	}
	current := f.current(id, "member")
	if current.Revision != revision+1 || current.Wizards[f.ids["member"]].Resources["wis"] != 9 {
		t.Fatal("competing reaction spent twice")
	}
}
