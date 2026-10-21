package bot_test

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/bot"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/internal/testcatalogue"
)

func playing(t *testing.T) (*game.State, *api.Catalogue) {
	t.Helper()
	c := testcatalogue.Load(t)
	d, err := game.NewDefinitions(*c)
	if err != nil {
		t.Fatal(err)
	}
	s, err := game.NewSetup("match", "Test", "lobby", []string{"one", "two"}, d)
	if err != nil {
		t.Fatal(err)
	}
	for _, actor := range s.Participants {
		for _, pick := range [][2]string{{"choose_chantry", "chantry_fire"}, {"choose_victory_path", "subjugate_rival"}, {"choose_minion", "bastian_redhand"}} {
			s, err = s.Choose(actor, pick[0], pick[1])
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	s.Choice = nil
	s.Exchange = nil
	s.Resolving = nil
	s.PhaseEffects = nil
	s.Turn = &api.Turn{Number: 10, ActiveWizard: "one", Phase: "actions", Step: "actions"}
	s.Wizards["one"].Resources = map[string]int{"wis": 10, "wealth": 10}
	return s, c
}

func accept(t *testing.T, s *game.State, c *api.Catalogue, actions int) (*game.State, game.Command) {
	t.Helper()
	before, _ := json.Marshal(s)
	view := s.View("one")
	viewBefore, _ := json.Marshal(view)
	candidates := bot.Candidates(view, c, actions)
	if !reflect.DeepEqual(candidates, bot.Candidates(view, c, actions)) {
		t.Fatal("decisions are not deterministic")
	}
	for _, command := range candidates {
		next, err := s.Apply("one", command)
		if err == nil {
			after, _ := json.Marshal(s)
			viewAfter, _ := json.Marshal(view)
			if !slices.Equal(before, after) || !slices.Equal(viewBefore, viewAfter) {
				t.Fatal("candidate generation or validation mutated its inputs")
			}
			return next, command
		}
	}
	t.Fatalf("no legal candidates for choice=%+v exchange=%+v", s.Choice, s.Exchange)
	return nil, game.Command{}
}

func TestEveryRequiredChoice(t *testing.T) {
	for _, kind := range []string{"maintenance", "kneel", "effects", "end_discard", "research", "hand_discard", "asset_discard", "attack_loss", "scry"} {
		t.Run(kind, func(t *testing.T) {
			s, c := playing(t)
			w := s.Wizards["one"]
			q := &api.RequiredChoice{ID: "choice", Actor: "one", Type: kind, Minimum: 1, Maximum: 1}
			s.Choice = q
			resolving := func(stage string) {
				s.Resolving = &game.PendingEffect{PendingEffectView: api.PendingEffectView{ID: "effect", Actor: "one", Kind: "research_spell", Source: "source"}, Stage: stage}
				s.Exchange = &game.Exchange{Participants: []string{"one", "two"}, Stack: []*game.PendingEffect{}}
			}
			switch kind {
			case "maintenance":
				s.Turn.Phase = "start"
				s.Turn.Step = "maintenance"
				q.Minimum = 0
				q.Options = []string{w.InPlay[0].ID}
			case "kneel":
				s.Turn.Phase = "start"
				s.Turn.Step = "bend_the_knee"
				q.Options = []string{"kneel", "decline"}
			case "effects":
				s.Turn.Phase = "start"
				s.Turn.Step = "effects"
				s.PhaseEffects = []game.PhaseEffect{{ID: "solar", Actor: "one", Source: "source", Kind: "solar_eclipse"}}
				q.Options = []string{"solar"}
				q.Effects = []api.EffectOption{{ID: "solar", Kind: "solar_eclipse", Source: "source"}}
			case "end_discard":
				s.Turn.Phase = "end"
				s.Turn.Step = "discard"
				q.Minimum = 0
				q.Maximum = len(w.Hand)
				for _, card := range w.Hand {
					q.Options = append(q.Options, card.ID)
				}
			case "research":
				resolving("research")
				q.Options = []string{"spell_healing_spring", "counterspell"}
			case "hand_discard":
				resolving("discard")
				q.Options = []string{w.Hand[0].ID}
				q.Private = true
			case "asset_discard", "attack_loss":
				resolving("loss")
				q.Options = []string{w.InPlay[0].ID}
			case "scry":
				resolving("scry")
				q.Private = true
				q.Cards = slices.Clone(w.DrawPile[:3])
				for _, card := range q.Cards {
					q.Options = append(q.Options, card.ID)
				}
			}
			_, command := accept(t, s, c, 100)
			if command.Type != "resolve_choice" || command.ChoiceID != q.ID {
				t.Fatal("did not answer the current choice")
			}
		})
	}
}

func TestOptionalHuntCanBeDeclinedWithoutChallenger(t *testing.T) {
	s, c := playing(t)
	s.Turn.Phase = "start"
	s.Turn.Step = "effects"
	s.Wizards["one"].Hand = nil
	s.Wizards["one"].Rumours = []api.RumourView{{ID: "moon", DesignID: "rumour_full_moon", ExpiresAfterTurn: 1}}
	s.PhaseEffects = []game.PhaseEffect{{ID: "hunt", Actor: "one", Source: "moon", Kind: "moonlit_hunt", Optional: true}}
	s.Choice = &api.RequiredChoice{ID: "choice", Actor: "one", Type: "effects", Options: []string{"hunt", "skip:hunt"}, Effects: []api.EffectOption{{ID: "hunt", Source: "moon", Kind: "moonlit_hunt", Optional: true}}, Minimum: 1, Maximum: 1}
	_, command := accept(t, s, c, 0)
	if !slices.Equal(command.Selections, []string{"skip:hunt"}) {
		t.Fatal("attempted an impossible hunt")
	}
}

func TestCombatDefenceAndActionFallback(t *testing.T) {
	s, c := playing(t)
	w := s.Wizards["one"]
	w.InPlay = append(w.InPlay, api.CardInstance{ID: "guard", DesignID: "bastian_redhand", Owner: "one", Location: "inplay"})
	w.Hand = []api.CardInstance{{ID: "focus", DesignID: "arcane_focus", Owner: "one", Location: "hand"}}
	s.Exchange = &game.Exchange{Participants: []string{"one", "two"}, Priority: "one", Stack: []*game.PendingEffect{{PendingEffectView: api.PendingEffectView{ID: "attack", Kind: "subjugate_rival", Actor: "two", Target: "one", Minions: []string{"other-guard"}}}}}
	_, command := accept(t, s, c, 0)
	if command.Type != "defend" || !slices.Contains(command.Minions, "guard") {
		t.Fatal("failed to defend an attack")
	}
	s.Exchange = nil
	_, command = accept(t, s, c, 64)
	if command.Type != "end_turn" {
		t.Fatal("action budget did not end proactive play")
	}
	s.Exchange = &game.Exchange{Participants: []string{"one", "two"}, Priority: "one", Stack: []*game.PendingEffect{{PendingEffectView: api.PendingEffectView{ID: "asset", Kind: "study", Actor: "two"}}}}
	_, command = accept(t, s, c, 64)
	if command.Type != "pass" {
		t.Fatal("response fallback did not pass")
	}
}

func TestDraftUsesCurrentOffersAndPinnedDefinitions(t *testing.T) {
	s, c := playing(t)
	d := s.Definitions()
	s, err := game.NewSetup("setup", "Draft", "lobby", []string{"one", "two"}, d)
	if err != nil {
		t.Fatal(err)
	}
	for range 3 {
		offer, err := game.DraftView(s.View("one"), d)
		if err != nil {
			t.Fatal(err)
		}
		kind, id := bot.Draft(offer, c)
		if kind2, id2 := bot.Draft(offer, c); kind != kind2 || id != id2 {
			t.Fatal("draft was not repeatable")
		}
		s, err = s.Choose("one", kind, id)
		if err != nil {
			t.Fatal(err)
		}
	}
	if s.Wizards["one"].Draft.Step != "complete" {
		t.Fatal("draft incomplete")
	}
}
