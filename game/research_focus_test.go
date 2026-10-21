package game

import (
	"reflect"
	"testing"

	"github.com/davelawson/grim/api"
)

func TestResearchFocusPaysWithoutChangingPool(t *testing.T) {
	for _, tc := range []struct {
		name, mode, element, pool string
		inputs                    []string
	}{
		{"learn", "element", "earth", "", nil},
		{"unaligned", "discover", "", "unaligned", nil},
		{"elemental", "discover", "", "fire", []string{"fire"}},
		{"hybrid", "discover", "", "air_fire", []string{"fire", "air"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := table(t, 2)
			actor := s.Ring[0]
			s.Wizards[actor].Resources["wis"] = 0
			affinities := s.affinities(actor)
			source, focus := hand(s, actor, "research_spell"), hand(s, actor, "arcane_focus")
			support := []api.CardSupport{{CardID: focus, Domains: []string{"water"}}}
			for _, input := range tc.inputs {
				support = append(support, api.CardSupport{CardID: hand(s, actor, "element_"+input)})
			}
			s = act(t, s, actor, Command{Type: "play_card", CardID: source, Mode: tc.mode, Element: tc.element, Support: support})
			s = settle(t, s)
			if tc.mode == "discover" {
				if s.Choice == nil || s.Choice.Type != "research" || len(s.Choice.Options) == 0 {
					t.Fatal("missing Research offer")
				}
				for _, id := range s.Choice.Options {
					if s.definition(id).Pool != tc.pool {
						t.Fatal("Focus changed discovery pool", id)
					}
				}
				s = answer(t, s, s.Choice.Options[0])
			}
			w := s.Wizards[actor]
			if w.Resources["wis"] != 0 || !reflect.DeepEqual(affinities, s.affinities(actor)) {
				t.Fatal("Focus leaked into standing resources or Affinities")
			}
			for _, id := range []string{source, focus} {
				if _, card := find(w.Discard, id); card == nil {
					t.Fatal("Research or Focus was not discarded", id)
				}
			}
			if tc.mode == "element" {
				found := false
				for _, card := range w.Discard {
					found = found || card.DesignID == "element_earth"
				}
				if !found {
					t.Fatal("chosen Element was not gained")
				}
			}
		})
	}
}
