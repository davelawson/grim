package game

import (
	"errors"
	randv2 "math/rand/v2"
	"slices"
)

// validateSavedState checks structural invariants before a persisted continuation
// can reach handlers. It does not change state or regenerate random outcomes.
func (s *State) validateSavedState() error {
	invalid := errors.New("invalid saved gameplay state")
	if s.Random.Generator != "chacha8-v1" {
		return ErrUnsupportedVersion
	}
	g := randv2.NewChaCha8([32]byte{})
	if len(s.Random.Seed) != 32 {
		return invalid
	}
	if err := g.UnmarshalBinary(s.Random.State); err != nil {
		return invalid
	}
	if s.Revision < 0 || s.Sequence < 0 || !unique(s.Ring) {
		return invalid
	}
	for _, id := range s.Ring {
		if s.Wizards[id] == nil {
			return invalid
		}
	}
	if s.Status == "setup" || (s.Status == "finished" && s.Outcome != nil && *s.Outcome == "abandoned" && s.Turn == nil) {
		if s.Turn != nil || len(s.Ring) > 0 || s.Exchange != nil || s.Choice != nil || s.Resolving != nil {
			return invalid
		}
		return nil
	}
	if s.Status != "playing" && s.Status != "finished" {
		return invalid
	}
	if s.Turn == nil || s.Turn.Number < 1 || s.Wizards[s.Turn.ActiveWizard] == nil {
		return invalid
	}
	if s.Status == "playing" && (len(s.Ring) < 2 || !slices.Contains(s.Ring, s.Turn.ActiveWizard) || s.Outcome != nil) {
		return invalid
	}
	if s.Status == "finished" && (s.Outcome == nil || len(s.Winners) > len(s.Participants)) {
		return invalid
	}
	if s.Status == "playing" {
		validSteps := map[string][]string{"start": {"recovery", "income", "maintenance", "bend_the_knee", "rumour", "effects"}, "actions": {"actions"}, "end": {"prepare_effects", "effects", "cleanup", "discard", "draw"}}
		if !slices.Contains(validSteps[s.Turn.Phase], s.Turn.Step) {
			return invalid
		}
	}
	seen := map[string]bool{}
	checkCards := func(owner, zone string, cards []CardInstance) bool {
		for _, c := range cards {
			d, ok := cardDefinitions[c.DesignID]
			if !ok || seen[c.ID] || c.ID == "" || c.Owner != owner || c.Location != zone || category(d, "rumour") {
				return false
			}
			seen[c.ID] = true
		}
		return true
	}
	for _, id := range s.Participants {
		w := s.Wizards[id]
		if w.Resources == nil || w.VictoryProgress == nil || w.Integrity < 0 || w.Integrity > 3 || w.Turns < 0 || w.ImmuneUntil < 0 {
			return invalid
		}
		for _, n := range w.Resources {
			if n < 0 {
				return invalid
			}
		}
		if !checkCards(id, "hand", w.Hand) || !checkCards(id, "drawpile", w.DrawPile) || !checkCards(id, "discard", w.Discard) || !checkCards(id, "inplay", w.InPlay) {
			return invalid
		}
		busy := map[string]bool{}
		for _, r := range w.Rumours {
			if r.ID == "" || busy[r.DesignID] || !category(definition(r.DesignID), "rumour") {
				return invalid
			}
			busy[r.DesignID] = true
		}
		for _, design := range w.Reserved {
			if busy[design] || !category(definition(design), "rumour") {
				return invalid
			}
			busy[design] = true
		}
	}
	if s.Choice != nil {
		q := s.Choice
		if q.ID == "" || s.Wizards[q.Actor] == nil || q.Minimum < 0 || q.Maximum < q.Minimum || q.Maximum > len(q.Options) || !unique(q.Options) {
			return invalid
		}
	}
	if s.Resolving != nil && (s.Exchange == nil || s.Choice == nil) {
		return invalid
	}
	if s.Exchange != nil {
		x := s.Exchange
		if !unique(x.Participants) || len(x.Participants) == 0 || x.Passes < 0 {
			return invalid
		}
		for _, id := range x.Participants {
			if s.Wizards[id] == nil {
				return invalid
			}
		}
		if s.Status == "playing" && s.Resolving == nil && (!slices.Contains(x.Participants, x.Priority) || !s.alive(x.Priority) || len(x.Stack) == 0) {
			return invalid
		}
		effects := slices.Clone(x.Stack)
		if s.Resolving != nil {
			effects = append(effects, s.Resolving)
		}
		ids := map[string]bool{}
		for _, e := range effects {
			if e == nil || e.ID == "" || ids[e.ID] || s.Wizards[e.Actor] == nil || !checkCards(e.Actor, "stack", e.Cards) {
				return invalid
			}
			ids[e.ID] = true
		}
	}
	return nil
}
