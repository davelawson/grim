package server

import (
	"encoding/json"
	"fmt"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/internal/testcatalogue"
	"github.com/google/uuid"
	randv2 "math/rand/v2"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func (f *lifecycleFixture) command(id, actor string, c game.Command) *game.View {
	f.t.Helper()
	view := f.current(id, actor)
	revision := view.Revision
	c.ExpectedRevision = &revision
	c.RequestID = uuid.NewString()
	path, input := typedGameplay(f.t, view, c)
	return matchView(f.t, f.call(200, "POST", path, actor, gameplayBody(f.t, input)))
}
func (f *lifecycleFixture) actor(id string) string {
	for actor, userID := range f.ids {
		if userID == id {
			return actor
		}
	}
	f.t.Fatalf("unknown user %s", id)
	return ""
}
func (f *lifecycleFixture) playable(size int) *game.View {
	f.t.Helper()
	lobby := f.createLobby(size)
	f.ready(lobby, size)
	view, _, _ := f.launch(lobby)
	// Pin only the initial setup seed, before any choices or gameplay, so the
	// entire match is reproducible. Everything after this uses production HTTP.
	_, s := f.saved(view.ID)
	seed := [32]byte{5, 8, 13}
	g := randv2.NewChaCha8(seed)
	state, err := g.MarshalBinary()
	if err != nil {
		f.t.Fatal(err)
	}
	s.Random.Seed = seed[:]
	s.Random.State = state
	data, _ := json.Marshal(s)
	f.exec("update matches set snapshot = ? where id = ?", string(data), view.ID)
	for _, actor := range roster[:size] {
		f.draft(view.ID, actor, "chantry_fire", "subjugate_rival", "bastian_redhand")
	}
	return f.current(view.ID, "owner")
}
func containsCategory(d api.CardDesign, category string) bool {
	return slices.Contains(d.Categories, category)
}

// playToVictory is a deliberately cooperative strategy. Opponents pass and
// forgo healing; the winner retains Subjugate Rival and Focus until a legal
// Attack is available. It uses only projected state and catalogue definitions.
func (f *lifecycleFixture) playToVictory(view *game.View, kneel bool) *game.View {
	return f.playToVictoryUsing(view, kneel, f.command, true)
}

func (f *lifecycleFixture) playToVictoryUsing(view *game.View, kneel bool, submit func(string, string, game.Command) *game.View, restart bool) *game.View {
	f.t.Helper()
	id := view.ID
	winner := f.ids["owner"]
	catalogue := testcatalogue.Load(f.t)
	restarted := map[string]bool{}
	for n := 0; n < 3500 && view.Status == "playing"; n++ {
		actor := f.actor(view.Turn.ActiveWizard)
		c := game.Command{}
		if view.Choice != nil {
			actor = f.actor(view.Choice.Actor)
			view = f.current(id, actor)
			q := view.Choice
			c.Type = "resolve_choice"
			c.ChoiceID = q.ID
			switch q.Type {
			case "effects":
				selected := q.Options[0]
				for _, option := range q.Options {
					if strings.HasPrefix(option, "skip:") {
						selected = option
						break
					}
				}
				c.Selections = []string{selected}
			case "kneel":
				selected := "decline"
				if kneel {
					selected = "kneel"
				}
				c.Selections = []string{selected}
			case "end_discard":
				for _, card := range view.You.Hand {
					keep := q.Actor == winner && (card.DesignID == "arcane_focus" || card.DesignID == "subjugate_rival")
					if !keep {
						c.Selections = append(c.Selections, card.ID)
					}
				}
			case "asset_discard":
				selected := q.Options[0]
				for _, card := range view.Wizards[q.Actor].InPlay {
					d, _ := catalogue.Card(card.DesignID)
					if containsCategory(d, "location") {
						selected = card.ID
						break
					}
				}
				c.Selections = []string{selected}
			case "maintenance": // This strategy never creates upkeep assets.
			default:
				c.Selections = q.Options[:q.Minimum]
			}
		} else if view.Exchange != nil {
			actor = f.actor(view.Exchange.Priority)
			c.Type = "pass"
		} else {
			view = f.current(id, actor)
			owner := view.Turn.ActiveWizard
			public := view.Wizards[owner]
			c.Type = "end_turn"
			if owner == winner {
				minions := []string{}
				might := public.Affinities["might"]
				for _, card := range public.InPlay {
					d, _ := catalogue.Card(card.DesignID)
					if containsCategory(d, "minion") && !card.Exhausted && card.DeployedTurn < view.Turn.Number && d.Expertise["might"] > 0 {
						minions = append(minions, card.ID)
						might += d.Expertise["might"]
					}
				}
				attack, focus := "", ""
				for _, card := range view.You.Hand {
					if card.DesignID == "subjugate_rival" {
						attack = card.ID
					}
					if card.DesignID == "arcane_focus" {
						focus = card.ID
					}
				}
				index := slices.Index(view.Ring, owner)
				target := view.Ring[(index+1)%len(view.Ring)]
				if attack != "" && focus != "" && len(minions) > 0 && !public.AttackUsed && !view.Wizards[target].Immune && might+1 >= view.Wizards[target].Affinities["might"]+2 {
					c = game.Command{Type: "play_card", CardID: attack, Minions: minions, Support: []api.CardSupport{{CardID: focus}}}
				} else {
					// Replay a lost Chantry before the Minion that depends on it.
					for _, card := range view.You.Hand {
						d, _ := catalogue.Card(card.DesignID)
						if card.DesignID != "chantry_fire" && card.DesignID != "bastian_redhand" {
							continue
						}
						if !meetsView(public.Affinities, d.PlayRequirements) {
							continue
						}
						payment := []string{}
						wis := public.Resources["wis"]
						for _, resource := range view.You.Hand {
							if wis >= d.Cost["wis"] {
								break
							}
							rd, _ := catalogue.Card(resource.DesignID)
							if rd.Resources["wis"] > 0 {
								wis += rd.Resources["wis"]
								payment = append(payment, resource.ID)
							}
						}
						if wis >= d.Cost["wis"] {
							c = game.Command{Type: "play_card", CardID: card.ID, Payment: payment}
							break
						}
					}
				}
			}
		}
		boundary := ""
		if view.Choice != nil {
			boundary = view.Choice.Type
		} else if view.Exchange != nil {
			boundary = "response"
		}
		if restart && boundary != "" && !restarted[boundary] {
			before := f.call(200, "GET", api.GetMatchPath(id), actor, "").Body.String()
			f.restart()
			if after := f.call(200, "GET", api.GetMatchPath(id), actor, "").Body.String(); after != before {
				f.t.Fatal("restart changed input boundary")
			}
			restarted[boundary] = true
		}
		view = submit(id, actor, c)
	}
	if view.Status != "finished" || view.Outcome == nil || *view.Outcome != "victory" || !slices.Equal(view.Winners, []string{winner}) {
		f.t.Fatalf("no winner after legal HTTP play: turn=%+v progress=%+v", view.Turn, view.Wizards[winner])
	}
	if kneel && view.Wizards[winner].VictoryProgress["domination"] != 5 {
		f.t.Fatal("did not exercise five-step victory")
	}
	if !kneel && len(view.Ring) != 1 {
		f.t.Fatal("did not exercise last-survivor victory")
	}
	return view
}
func meetsView(have, need map[string]int) bool {
	for k, v := range need {
		if have[k] < v {
			return false
		}
	}
	return true
}
func TestCompleteHTTPGames(t *testing.T) {
	for _, size := range []int{2, 3, 4} {
		for _, kneel := range []bool{false, true} {
			if size == 2 && kneel {
				continue
			}
			t.Run(fmt.Sprintf("%d/kneel=%v", size, kneel), func(t *testing.T) {
				f := newLifecycleFixture(t)
				view := f.playToVictory(f.playable(size), kneel)
				f.restart()
				restored := f.current(view.ID, f.actor(view.Turn.ActiveWizard))
				if !reflect.DeepEqual(view.Winners, restored.Winners) || view.Status != restored.Status {
					t.Fatal("winner lost after restart")
				}
				revision := view.Revision
				body := gameplayBody(t, api.EndTurnRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision})
				errorCode(t, f.call(409, "POST", api.EndTurnPath(view.ID), "owner", body), "match_finished")
			})
		}
	}
}

func TestHTTPResearchAndPrivateResolutionRecovery(t *testing.T) {
	f := newLifecycleFixture(t)
	view := f.playable(2)
	id := view.ID
	_, state := f.saved(id)
	owner := f.ids["owner"]
	state.Turn = &api.Turn{Number: 10, ActiveWizard: owner, Phase: "actions", Step: "actions"}
	state.Exchange = nil
	state.Choice = nil
	state.Resolving = nil
	state.PhaseEffects = nil
	w := state.Wizards[owner]
	w.Resources = map[string]int{"wis": 10}
	w.Hand = []api.CardInstance{{ID: "test/research", DesignID: "research_spell", Owner: owner, Location: "hand"}}
	w.Rumours = nil
	data, _ := json.Marshal(state)
	f.exec("update matches set snapshot = ? where id = ?", string(data), id)
	view = f.command(id, "owner", game.Command{Type: "play_card", CardID: "test/research", Mode: "discover"})
	// Stop just before the final pass creates its random Research offer.
	for view.Exchange.Passes < len(view.Exchange.Participants)-1 {
		view = f.command(id, f.actor(view.Exchange.Priority), game.Command{Type: "pass"})
	}
	actor := f.actor(view.Exchange.Priority)
	revision := view.Revision
	command := game.Command{Type: "pass", RequestID: uuid.NewString(), ExpectedRevision: &revision}
	body := gameplayBody(t, api.PassRequest{RequestID: command.RequestID, ExpectedRevision: command.ExpectedRevision})
	before, prior := f.saved(id)
	expected, err := prior.Apply(f.ids[actor], command)
	if err != nil {
		t.Fatal(err)
	}
	f.exec(`create trigger fail before insert on command_receipts begin select raise(abort, 'private random failure'); end`)
	f.call(500, "POST", api.PassPath(id), actor, body)
	after, _ := f.saved(id)
	if !slices.Equal(before, after) {
		t.Fatal("rollback consumed random offer")
	}
	f.exec("drop trigger fail")
	response := f.call(200, "POST", api.PassPath(id), actor, body)
	view = matchView(t, response)
	_, saved := f.saved(id)
	if !slices.Equal(mustJSON(t, expected), mustJSON(t, saved)) {
		t.Fatal("retried random transition changed outcomes")
	}
	if view.Choice == nil || view.Choice.Type != "research" || len(view.Choice.Options) != 3 {
		t.Fatal("Research offer not persisted")
	}
	offer := slices.Clone(view.Choice.Options)
	f.restart()
	view = f.current(id, "owner")
	if !slices.Equal(view.Choice.Options, offer) {
		t.Fatal("restart rerolled offer")
	}
	if f.call(200, "POST", api.PassPath(id), actor, body).Body.String() != response.Body.String() {
		t.Fatal("random command retry changed response")
	}
	q := view.Choice
	view = f.command(id, "owner", game.Command{Type: "resolve_choice", ChoiceID: q.ID, Selections: offer[:1]})
	// Set up the later Scrying Pool play with a known two-card draw pile. The
	// continuation and every response/choice still use production persistence.
	_, state = f.saved(id)
	w = state.Wizards[owner]
	w.Resources = map[string]int{"wis": 10}
	w.Modifiers = append(w.Modifiers, api.Modifier{Source: "test", Affinity: map[string]int{"air": 1, "water": 1}, ExpiresAfterTurn: 99})
	w.Affinities["air"] = 1
	w.Affinities["water"] = 1
	w.Hand = append(w.Hand, api.CardInstance{ID: "test/scry", DesignID: "scrying_pool", Owner: owner, Location: "hand"})
	w.DrawPile = []api.CardInstance{{ID: "test/top", DesignID: "basic_wis", Owner: owner, Location: "drawpile"}, {ID: "test/second", DesignID: "basic_wealth", Owner: owner, Location: "drawpile"}}
	data, _ = json.Marshal(state)
	f.exec("update matches set snapshot = ? where id = ?", string(data), id)
	view = f.command(id, "owner", game.Command{Type: "play_card", CardID: "test/scry"})
	for view.Choice == nil {
		view = f.command(id, f.actor(view.Exchange.Priority), game.Command{Type: "pass"})
	}
	own := f.current(id, "owner")
	other := f.current(id, "member")
	if own.Choice.Type != "scry" || len(own.Choice.Cards) != 2 || len(other.Choice.Cards) != 0 || len(other.Choice.Options) != 0 {
		t.Fatal("private HTTP inspection leaked")
	}
	serialized := f.call(200, "GET", api.GetMatchPath(id), "member", "").Body.String()
	if strings.Contains(serialized, "test/top") || strings.Contains(serialized, "test/second") {
		t.Fatal("private draw identity leaked through exchange")
	}
	f.restart()
	own = f.current(id, "owner")
	q = own.Choice
	revision = own.Revision
	decision := game.Command{Type: "resolve_choice", ChoiceID: q.ID, Selections: q.Options[:1], Order: q.Options[1:], ExpectedRevision: &revision, RequestID: uuid.NewString()}
	body = gameplayBody(t, api.ScryRequest{RequestID: decision.RequestID, ExpectedRevision: decision.ExpectedRevision, Take: decision.Selections[0], Order: decision.Order})
	accepted := f.call(200, "POST", api.ScryPath(id, q.ID), "owner", body).Body.String()
	f.restart()
	if f.call(200, "POST", api.ScryPath(id, q.ID), "owner", body).Body.String() != accepted {
		t.Fatal("scry decision was repeated")
	}
	view = f.current(id, "owner")
	if len(view.You.DrawPile) != 1 || view.You.DrawPile[0].ID != "test/second" {
		t.Fatal("scry continuation lost remaining card")
	}
}
func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
