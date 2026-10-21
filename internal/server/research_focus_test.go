package server

import (
	"bytes"
	"testing"

	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"github.com/google/uuid"
)

func researchFocusFixture(t *testing.T) (*lifecycleFixture, *game.View) {
	t.Helper()
	f, view := reactionFixture(t)
	_, state := f.saved(view.ID)
	owner := f.ids["owner"]
	w := state.Wizards[owner]
	w.Resources["wis"] = 0
	for _, pair := range [][2]string{{"research", "research_spell"}, {"focus", "arcane_focus"}, {"focus-two", "arcane_focus"}, {"fire", "element_fire"}, {"fire-two", "element_fire"}, {"air", "element_air"}, {"earth", "element_earth"}, {"support-spell", "flamecraft"}} {
		w.Hand = append(w.Hand, api.CardInstance{ID: "research/" + pair[0], DesignID: pair[1], Owner: owner, Location: "hand"})
	}
	f.exec("update matches set snapshot = ? where id = ?", string(mustJSON(t, state)), view.ID)
	return f, f.current(view.ID, "owner")
}

func TestTypedResearchFocusAndReceiptReplay(t *testing.T) {
	for _, mode := range []string{"discover", "element"} {
		t.Run(mode, func(t *testing.T) {
			f, view := researchFocusFixture(t)
			support := []api.CardSupport{{CardID: "research/focus", Domains: []string{"water"}}}
			var input any
			path := api.ResearchSpellPath(view.ID)
			if mode == "element" {
				path = api.ResearchElementPath(view.ID)
				input = api.ResearchElementRequest{RequestID: uuid.NewString(), ExpectedRevision: &view.Revision, CardID: "research/research", Element: "earth", Support: support}
			} else {
				input = api.ResearchSpellRequest{RequestID: uuid.NewString(), ExpectedRevision: &view.Revision, CardID: "research/research", Elements: []string{"research/fire"}, Support: support}
			}
			body := gameplayBody(t, input)
			accepted := f.call(200, "POST", path, "owner", body)
			result := matchView(t, accepted)
			effect := result.Exchange.Stack[0]
			if effect.Mode != mode || len(effect.Support) == 0 || effect.Support[len(effect.Support)-1].CardID != "research/focus" || result.Wizards[f.ids["owner"]].Resources["wis"] != 0 {
				t.Fatal("Research lost Focus payment support")
			}
			if mode == "discover" && effect.Element != "fire" {
				t.Fatal("Focus changed discovery pool")
			}
			before, _ := f.saved(view.ID)
			replayed := f.call(200, "POST", path, "owner", body)
			after, _ := f.saved(view.ID)
			if accepted.Body.String() != replayed.Body.String() || !bytes.Equal(before, after) {
				t.Fatal("accepted retry changed its response or mutated the match")
			}
		})
	}
}

func TestTypedResearchRejectsInvalidSupport(t *testing.T) {
	for _, mode := range []string{"discover", "element"} {
		t.Run(mode, func(t *testing.T) {
			f, view := researchFocusFixture(t)
			before, _ := f.saved(view.ID)
			focus := api.CardSupport{CardID: "research/focus", Domains: []string{"water"}}
			for _, support := range [][]api.CardSupport{
				nil, // No available Wis and no Focus.
				{{CardID: "research/focus"}},
				{{CardID: "research/focus", Domains: []string{"unknown"}}},
				{{CardID: "research/focus", Domains: []string{"water", "air"}}},
				{{CardID: "research/support-spell", Domains: []string{"water"}}},
				{{CardID: "research/missing", Domains: []string{"water"}}},
				{focus, focus},
				{focus, {CardID: "research/focus-two", Domains: []string{"fire"}}},
				{{CardID: "research/research", Domains: []string{"water"}}},
				{{CardID: "research/fire", Domains: []string{"water"}}},
			} {
				var input any
				path := api.ResearchSpellPath(view.ID)
				if mode == "element" {
					path = api.ResearchElementPath(view.ID)
					input = api.ResearchElementRequest{RequestID: uuid.NewString(), ExpectedRevision: &view.Revision, CardID: "research/research", Element: "earth", Support: support}
				} else {
					input = api.ResearchSpellRequest{RequestID: uuid.NewString(), ExpectedRevision: &view.Revision, CardID: "research/research", Support: support}
				}
				errorCode(t, f.call(400, "POST", path, "owner", gameplayBody(t, input)), "invalid_choice")
			}
			if mode == "discover" {
				for _, elements := range [][]string{{"research/focus"}, {"research/fire", "research/fire-two"}, {"research/fire", "research/air", "research/earth"}, {"research/fire", "research/fire"}} {
					input := api.ResearchSpellRequest{RequestID: uuid.NewString(), ExpectedRevision: &view.Revision, CardID: "research/research", Elements: elements, Support: []api.CardSupport{focus}}
					errorCode(t, f.call(400, "POST", api.ResearchSpellPath(view.ID), "owner", gameplayBody(t, input)), "invalid_choice")
				}
			} else {
				input := api.ResearchElementRequest{RequestID: uuid.NewString(), ExpectedRevision: &view.Revision, CardID: "research/research", Element: "earth", Support: []api.CardSupport{focus, {CardID: "research/fire"}}}
				errorCode(t, f.call(400, "POST", api.ResearchElementPath(view.ID), "owner", gameplayBody(t, input)), "invalid_choice")
			}
			after, _ := f.saved(view.ID)
			if !bytes.Equal(before, after) {
				t.Fatal("rejected support mutated the match")
			}
		})
	}
}
