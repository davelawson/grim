package game

import (
	"encoding/json"
	"errors"
	randv2 "math/rand/v2"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func setup(t *testing.T, size int) *State {
	t.Helper()
	s, err := NewSetup("match", "Test", "lobby", []string{"d", "c", "b", "a"}[4-size:])
	if err != nil {
		t.Fatal(err)
	}
	seed := [32]byte{1, 2, 3}
	g := randv2.NewChaCha8(seed)
	state, err := g.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s.Random.Seed, s.Random.State = slices.Clone(seed[:]), state
	return s
}

func choose(t *testing.T, s *State, actor, kind, card string) *State {
	t.Helper()
	next, err := s.Choose(actor, kind, card)
	if err != nil {
		t.Fatal(err)
	}
	if next.Revision != s.Revision+1 {
		t.Fatal("choice did not advance revision once")
	}
	return next
}

func jsonBytes(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestAllDraftOffersAndStartingComposition(t *testing.T) {
	offers := OpeningOffers()
	for _, chantry := range offers.Chantries {
		for _, victory := range offers.VictoryPaths {
			for _, minion := range offers.Minions[chantry] {
				t.Run(chantry+"/"+victory+"/"+minion, func(t *testing.T) {
					s := setup(t, 2)
					// Identical designs remain available to both players.
					for _, actor := range s.Participants {
						s = choose(t, s, actor, "choose_chantry", chantry)
						s = choose(t, s, actor, "choose_victory_path", victory)
						s = choose(t, s, actor, "choose_minion", minion)
					}
					if s.Status != "playing" || s.Revision != 6 || s.Turn == nil || s.Turn.ActiveWizard != s.Ring[0] || s.Turn.Phase != "start" || s.Turn.Step != "recovery" || s.Turn.Number != 1 {
						t.Fatalf("incorrect game start: %+v", s)
					}
					seen := map[string]bool{}
					for _, owner := range s.Participants {
						w := s.Wizards[owner]
						if w.Integrity != 3 || len(w.InPlay) != 1 || len(w.Hand) != 5 || len(w.DrawPile) != 5 || len(w.Discard) != 0 {
							t.Fatal("wrong opening zones")
						}
						if w.InPlay[0].DesignID != chantry || w.InPlay[0].Location != "inplay" {
							t.Fatal("wrong starting Chantry")
						}
						design, _ := OpeningCatalogue().Card(chantry)
						if !reflect.DeepEqual(w.Affinities, map[string]int{design.Element: 1}) || len(w.VictoryProgress) != 5 {
							t.Fatal("incorrect public stats")
						}
						for _, progress := range w.VictoryProgress {
							if progress != 0 {
								t.Fatal("nonzero initial progress")
							}
						}
						counts := map[string]int{}
						for _, card := range append(slices.Clone(w.Hand), w.DrawPile...) {
							counts[card.DesignID]++
						}
						want := map[string]int{"element_" + design.Element: 1, "research_spell": 1, "research_artifact": 1, "basic_wis": 2, "basic_wealth": 1, "arcane_focus": 1, "investigate": 1, victory: 1, minion: 1}
						if !reflect.DeepEqual(counts, want) {
							t.Fatalf("wrong deck: %v", counts)
						}
						for _, zone := range []struct {
							location string
							cards    []CardInstance
						}{{"inplay", w.InPlay}, {"hand", w.Hand}, {"drawpile", w.DrawPile}} {
							for _, card := range zone.cards {
								if seen[card.ID] || card.Owner != owner || card.Location != zone.location || card.Exhausted {
									t.Fatalf("bad instance: %+v", card)
								}
								seen[card.ID] = true
							}
						}
					}
					if len(seen) != 22 {
						t.Fatal("incorrect instance count")
					}
				})
			}
		}
	}
}

func TestDraftRejectionsAreImmutable(t *testing.T) {
	s := setup(t, 2)
	check := func(actor, kind, card string, expected error) {
		t.Helper()
		before := jsonBytes(t, s)
		next, err := s.Choose(actor, kind, card)
		if next != nil || !errors.Is(err, expected) || !slices.Equal(before, jsonBytes(t, s)) {
			t.Fatalf("bad rejection %s: %v", kind, err)
		}
	}
	check("outsider", "choose_chantry", "chantry_fire", ErrNotParticipant)
	check("a", "choose_minion", "bastian_redhand", ErrDraftSequence)
	check("a", "choose_chantry", "fire", ErrInvalidChoice)
	s = choose(t, s, "a", "choose_chantry", "chantry_fire")
	check("a", "choose_chantry", "chantry_air", ErrDraftSequence)
	s = choose(t, s, "a", "choose_victory_path", "subjugate_rival")
	check("a", "choose_minion", "rurik_ironvale", ErrInvalidChoice)
	s = choose(t, s, "a", "choose_minion", "orla_nine_cinders")
	check("a", "choose_minion", "bastian_redhand", ErrDraftSequence)
	s = choose(t, s, "b", "choose_chantry", "chantry_water")
	s = choose(t, s, "b", "choose_victory_path", "victory_arcane")
	s.Random.State = []byte("bad generator state")
	before := jsonBytes(t, s)
	if next, err := s.Choose("b", "choose_minion", "nadia_rivermark"); next != nil || err == nil || !slices.Equal(before, jsonBytes(t, s)) {
		t.Fatal("failed start mutated state")
	}
}

func TestDeterministicStartAndRestart(t *testing.T) {
	for _, size := range []int{2, 3, 4} {
		s := setup(t, size)
		for _, actor := range s.Participants {
			s = choose(t, s, actor, "choose_chantry", "chantry_fire")
			s = choose(t, s, actor, "choose_victory_path", "subjugate_rival")
			if actor != s.Participants[size-1] {
				s = choose(t, s, actor, "choose_minion", "bastian_redhand")
			}
		}
		before := jsonBytes(t, s)
		restored, err := Decode(before)
		if err != nil {
			t.Fatal(err)
		}
		a := choose(t, s, s.Participants[size-1], "choose_minion", "bastian_redhand")
		b := choose(t, restored, s.Participants[size-1], "choose_minion", "bastian_redhand")
		if !slices.Equal(jsonBytes(t, a), jsonBytes(t, b)) || !slices.Equal(before, jsonBytes(t, s)) {
			t.Fatal("start is not reproducible or mutated prior state")
		}
		generator := randv2.NewChaCha8([32]byte(s.Random.Seed))
		random := randv2.New(generator)
		ring := slices.Clone(s.Participants)
		random.Shuffle(size, func(i, j int) { ring[i], ring[j] = ring[j], ring[i] })
		if !slices.Equal(ring, a.Ring) {
			t.Fatal("seed does not reproduce seating")
		}
		for _, actor := range s.Participants {
			deck := []string{"element_fire", "research_spell", "research_artifact", "basic_wis", "basic_wis", "subjugate_rival", "basic_wealth", "bastian_redhand", "arcane_focus", "investigate"}
			random.Shuffle(10, func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
			cards := append(slices.Clone(a.Wizards[actor].Hand), a.Wizards[actor].DrawPile...)
			for i, card := range cards {
				if card.DesignID != deck[i] {
					t.Fatal("seed does not reproduce deal")
				}
			}
		}
		resumed := randv2.NewChaCha8([32]byte{})
		if err := resumed.UnmarshalBinary(a.Random.State); err != nil {
			t.Fatal(err)
		}
		if resumed.Uint64() != generator.Uint64() {
			t.Fatal("wrong saved random state")
		}
	}
}

func TestParticipantPrivacyAndProjectionIsolation(t *testing.T) {
	s := setup(t, 2)
	before := jsonBytes(t, s.View("b"))
	s = choose(t, s, "a", "choose_chantry", "chantry_fire")
	s = choose(t, s, "a", "choose_victory_path", "victory_fame")
	otherView := s.View("b")
	otherView.Revision = 0 // The global revision is intentionally public.
	if !slices.Equal(before, jsonBytes(t, otherView)) {
		t.Fatal("opponent choices or step leaked")
	}
	owner := s.View("a")
	if owner.You.Draft.Chantry != "chantry_fire" || owner.You.Draft.VictoryPath != "victory_fame" {
		t.Fatal("own choices missing")
	}
	s = choose(t, s, "a", "choose_minion", "bastian_redhand")
	if !s.View("b").Wizards["a"].DraftComplete || s.View("b").Wizards["a"].Integrity != nil {
		t.Fatal("bad completion visibility")
	}
	for _, selection := range []struct{ kind, card string }{{"choose_chantry", "chantry_water"}, {"choose_victory_path", "victory_arcane"}, {"choose_minion", "nadia_rivermark"}} {
		s = choose(t, s, "b", selection.kind, selection.card)
	}
	view := s.View("b")
	data := string(jsonBytes(t, view))
	for _, hidden := range []string{"bastian_redhand", "victory_fame", "seed", "random", "generator"} {
		if strings.Contains(data, hidden) {
			t.Fatalf("leaked %s", hidden)
		}
	}
	if len(view.You.Hand) != 5 || len(view.You.DrawPile) != 5 || view.Wizards["a"].HandCount != 5 {
		t.Fatal("wrong pile view")
	}
	slices.Reverse(s.Wizards["b"].DrawPile)
	if !slices.Equal(jsonBytes(t, view), jsonBytes(t, s.View("b"))) {
		t.Fatal("draw order affects view")
	}
	view.You.Hand[0].DesignID = "changed"
	view.You.DrawPile[0].DesignID = "changed"
	view.Wizards["b"].Affinities["water"] = 99
	view.Turn.Step = "changed"
	if strings.Contains(string(jsonBytes(t, s)), "changed") || s.Wizards["b"].Affinities["water"] != 1 {
		t.Fatal("view aliases state")
	}
}

func TestOpeningCataloguePropertiesAndCompatibility(t *testing.T) {
	c := OpeningCatalogue()
	if len(c.Cards) != 31 {
		t.Fatalf("unexpected catalogue size %d", len(c.Cards))
	}
	seen := map[string]bool{}
	for _, card := range c.Cards {
		if seen[card.ID] {
			t.Fatal("duplicate design ID")
		}
		seen[card.ID] = true
	}
	for _, chantry := range c.DraftOffers.Chantries {
		card, ok := c.Card(chantry)
		if !ok || card.Affinity[card.Element] != 1 || card.Income["ephemeral_wis"] != 1 || card.Cost["wis"] != 1 {
			t.Fatal("wrong Chantry stats")
		}
		for _, id := range c.DraftOffers.Minions[chantry] {
			minion, ok := c.Card(id)
			if !ok || minion.Cost["wis"] != 1 || minion.PlayRequirements[card.Element] != 1 || len(minion.Expertise) != 1 || len(minion.Affinity) != 0 || len(minion.Income) != 0 {
				t.Fatal("wrong Minion stats")
			}
			for _, rank := range minion.Expertise {
				if rank != 1 {
					t.Fatal("wrong Expertise rank")
				}
			}
		}
	}
	for _, path := range c.DraftOffers.VictoryPaths {
		card, ok := c.Card(path)
		if !ok || card.Placeholder != (path != "subjugate_rival") {
			t.Fatal("wrong placeholder flag")
		}
		if card.Placeholder && len(card.Effects) != 0 {
			t.Fatal("placeholder has behavior")
		}
	}
	c.Cards[0].Cost["wis"] = 99
	if OpeningCatalogue().Cards[0].Cost["wis"] != 1 {
		t.Fatal("catalogue aliases shared data")
	}
	s := setup(t, 2)
	s.SchemaVersion = 1
	if _, err := Decode(jsonBytes(t, s)); !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatal("old schema accepted")
	}
	s.SchemaVersion = SchemaVersion
	s.RulesVersion = "unknown"
	if _, err := Decode(jsonBytes(t, s)); !errors.Is(err, ErrUnsupportedVersion) {
		t.Fatal("unknown rules accepted")
	}
	s.RulesVersion, s.CatalogueVersion = OpeningRulesVersion, "unknown"
	if _, err := Decode(jsonBytes(t, s)); !errors.Is(err, ErrUnsupportedVersion) {
		t.Fatal("unknown catalogue accepted")
	}
}
