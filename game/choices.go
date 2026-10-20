package game

import (
	"slices"
)

func (s *State) resolveChoice(actor string, c Command) error {
	q := s.Choice
	if q == nil || q.ID != c.ChoiceID {
		return ErrStaleChoice
	}
	if q.Actor != actor {
		return ErrPriority
	}
	if !unique(c.Selections) || len(c.Selections) < q.Minimum || len(c.Selections) > q.Maximum {
		return ErrInvalidChoice
	}
	for _, id := range c.Selections {
		if !slices.Contains(q.Options, id) {
			return ErrInvalidChoice
		}
	}
	if q.Type != "scry" && len(c.Order) > 0 {
		return ErrInvalidChoice
	}
	if q.Type != "hand_discard" && c.Mode != "" {
		return ErrInvalidChoice
	}
	if q.Type != "maintenance" && q.Type != "effects" && (len(c.Payment) > 0 || len(c.Minions) > 0 || len(c.Support) > 0) {
		return ErrInvalidChoice
	}
	if q.Type == "maintenance" && (len(c.Minions) > 0 || len(c.Support) > 0) {
		return ErrInvalidChoice
	}
	w := s.Wizards[actor]
	s.Choice = nil
	switch q.Type {
	case "maintenance":
		cost := map[string]int{}
		for _, id := range c.Selections {
			card := s.inPlay(id)
			if card == nil || card.Owner != actor {
				return ErrInvalidChoice
			}
			addStats(cost, s.definition(card.DesignID).Upkeep)
		}
		if _, err := s.pay(actor, cost, c.Payment); err != nil {
			return err
		}
		for _, id := range q.Options {
			if !slices.Contains(c.Selections, id) {
				s.removeAsset(id, "discard")
			}
		}
		s.Turn.Step = "bend_the_knee"
	case "kneel":
		if c.Selections[0] == "kneel" {
			return s.kneel()
		}
		w.Kneel = false
		s.Turn.Step = "rumour"
	case "effects":
		selected, skip := skipEffect(c.Selections[0])
		index := -1
		var effect PhaseEffect
		for i, e := range s.PhaseEffects {
			if e.ID == selected {
				index = i
				effect = e
				break
			}
		}
		if index < 0 || skip && !effect.Optional {
			return ErrInvalidChoice
		}
		s.PhaseEffects = slices.Delete(s.PhaseEffects, index, index+1)
		if !skip && effect.Kind == "moonlit_hunt" {
			return s.declare(actor, Command{Type: "attempt_challenge", CardID: effect.Source, ChallengeID: "moonlit_hunt", Minions: c.Minions, Support: c.Support, Payment: c.Payment}, true)
		}
		if len(c.Minions) > 0 || len(c.Support) > 0 || len(c.Payment) > 0 {
			return ErrInvalidChoice
		}
		if !skip {
			s.openPhaseEffect(effect)
		}
	case "end_discard":
		for _, id := range c.Selections {
			card, ok := take(&w.Hand, id)
			if !ok {
				return ErrInvalidChoice
			}
			s.dispose(card, "discard")
		}
		if len(w.Hand) > s.maxHand(actor) {
			return ErrInvalidChoice
		}
		s.Turn.Step = "draw"
	case "research":
		if s.Resolving == nil || s.Resolving.Stage != "research" {
			return ErrTiming
		}
		if err := s.gain(actor, c.Selections[0]); err != nil {
			return err
		}
		return s.finishResolution()
	case "hand_discard":
		if s.Resolving == nil || s.Resolving.Stage != "discard" {
			return ErrTiming
		}
		if c.Mode != "" && c.Mode != "discard" && c.Mode != "destroy" {
			return ErrInvalidChoice
		}
		if c.Mode == "destroy" && s.Resolving.Kind != "flash_of_insight" {
			return ErrInvalidChoice
		}
		card, ok := take(&w.Hand, c.Selections[0])
		if !ok {
			return ErrInvalidChoice
		}
		if c.Mode != "destroy" {
			s.dispose(card, "discard")
		}
		return s.finishResolution()
	case "asset_discard", "attack_loss":
		if s.Resolving == nil {
			return ErrTiming
		}
		card := s.inPlay(c.Selections[0])
		if card == nil || card.Owner != actor {
			return ErrInvalidChoice
		}
		s.removeAsset(card.ID, "discard")
		return s.finishResolution()
	case "scry":
		if s.Resolving == nil || s.Resolving.Stage != "scry" {
			return ErrTiming
		}
		remaining := []string{}
		for _, id := range q.Options {
			if id != c.Selections[0] {
				remaining = append(remaining, id)
			}
		}
		order := slices.Clone(c.Order)
		slices.Sort(order)
		slices.Sort(remaining)
		if !slices.Equal(order, remaining) || !unique(c.Order) {
			return ErrInvalidChoice
		}
		selected, ok := take(&w.DrawPile, c.Selections[0])
		if !ok {
			return ErrInvalidChoice
		}
		selected.Location = "hand"
		w.Hand = append(w.Hand, selected)
		top := []CardInstance{}
		for _, id := range c.Order {
			card, ok := take(&w.DrawPile, id)
			if !ok {
				return ErrInvalidChoice
			}
			top = append(top, card)
		}
		w.DrawPile = append(top, w.DrawPile...)
		return s.finishResolution()
	default:
		return ErrInvalidChoice
	}
	return nil
}
