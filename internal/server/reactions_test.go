package server

import (
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"testing"
)

func reactionFixture(t *testing.T) (*lifecycleFixture, *game.View) {
	t.Helper()
	f := newLifecycleFixture(t)
	view := f.playable(2)
	_, s := f.saved(view.ID)
	s.Turn = &api.Turn{Number: 10, ActiveWizard: f.ids["owner"], Phase: "actions", Step: "actions"}
	s.Exchange = nil
	s.Choice = nil
	s.Resolving = nil
	s.PhaseEffects = nil
	for _, actor := range []string{"owner", "member"} {
		owner := f.ids[actor]
		w := s.Wizards[owner]
		w.Turns = 1
		w.Resources = map[string]int{"wis": 10}
		w.Rumours = nil
		w.Reserved = nil
		w.Deferred = nil
		w.Hand = nil
		w.Discard = nil
		w.DrawPile = nil
		w.InPlay = w.InPlay[:1]
		w.AttackUsed = false
		w.Kneel = false
		w.ImmuneUntil = 0
		w.Modifiers = []api.Modifier{{Source: "test", Affinity: map[string]int{"earth": 2, "air": 2, "fire": 2, "water": 2}, ExpiresAfterTurn: 99}}
		w.Affinities = map[string]int{"earth": 2, "air": 2, "fire": 3, "water": 2, "might": 1}
		assets := []string{"guard-member"}
		designs := map[string]string{"normal": "study", "counter": "counterspell", "stoneward": "stoneward", "focus": "arcane_focus"}
		if actor == "owner" {
			assets = []string{"guard-one", "guard-two"}
			w.Affinities["might"] = 2
			designs = map[string]string{"study": "study", "counter-owner": "counterspell", "attack": "subjugate_rival"}
		}
		for _, id := range assets {
			w.InPlay = append(w.InPlay, api.CardInstance{ID: "test/" + id, DesignID: "caravan_guard", Owner: owner, Location: "inplay"})
		}
		// Explicit iteration fixes initial test state independently of map order.
		for _, id := range []string{"normal", "counter", "stoneward", "focus", "study", "counter-owner", "attack"} {
			if design, ok := designs[id]; ok {
				w.Hand = append(w.Hand, api.CardInstance{ID: "test/" + id, DesignID: design, Owner: owner, Location: "hand"})
			}
		}
	}
	f.exec("update matches set snapshot = ? where id = ?", string(mustJSON(t, s)), view.ID)
	return f, f.current(view.ID, "owner")
}
