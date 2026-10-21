package game

import "strings"

// ValidateOperation binds a typed route to the selected card or pending choice.
// It is deliberately evaluated after receipt replay by the match service.
func (s *State) ValidateOperation(actor string, c Command, operation string) error {
	if strings.HasPrefix(operation, "choices/") {
		q := s.Choice
		if q == nil || q.ID != c.ChoiceID {
			return ErrStaleChoice
		}
		if q.Actor != actor {
			return ErrPriority
		}
		kind := map[string]string{"maintenance": "maintenance", "kneel": "kneel", "resolve-effect": "effects", "skip-effect": "effects", "undertake-hunt": "effects", "research": "research", "end-discard": "end_discard", "hand-discard": "hand_discard", "scry": "scry"}[strings.TrimPrefix(operation, "choices/")]
		if operation == "choices/asset-loss" {
			if q.Type != "asset_discard" && q.Type != "attack_loss" {
				return ErrInvalidChoice
			}
		} else if kind == "" || kind != q.Type {
			return ErrInvalidChoice
		}
		if q.Type == "hand_discard" && c.Mode != "discard" && c.Mode != "destroy" {
			return ErrInvalidChoice
		}
		if q.Type == "effects" {
			if len(c.Selections) != 1 {
				return ErrInvalidChoice
			}
			selected := c.Selections[0]
			if operation == "choices/skip-effect" {
				selected = strings.TrimPrefix(selected, "skip:")
			}
			for _, effect := range s.PhaseEffects {
				if effect.ID == selected {
					if operation == "choices/undertake-hunt" && effect.Kind != "moonlit_hunt" {
						return ErrInvalidChoice
					}
					if operation == "choices/resolve-effect" && effect.Kind == "moonlit_hunt" {
						return ErrInvalidChoice
					}
					return nil
				}
			}
			return ErrInvalidChoice
		}
		return nil
	}
	switch operation {
	case "actions/play-card", "actions/attack", "actions/research-element", "actions/research-spell", "actions/heal-with-focus":
		w := s.Wizards[actor]
		if w == nil {
			return ErrNotParticipant
		}
		_, card := find(w.Hand, c.CardID)
		if card == nil {
			return ErrInvalidChoice
		}
		expected := map[string]string{"actions/attack": "subjugate_rival", "actions/research-element": "research_spell", "actions/research-spell": "research_spell", "actions/heal-with-focus": "arcane_focus"}[operation]
		if expected != "" && card.DesignID != expected {
			return ErrInvalidChoice
		}
		if operation == "actions/play-card" && (card.DesignID == "subjugate_rival" || card.DesignID == "research_spell" || card.DesignID == "arcane_focus") {
			return ErrInvalidChoice
		}
		if operation == "actions/research-spell" || operation == "actions/research-element" {
			for _, support := range c.Support {
				_, paired := find(w.Hand, support.CardID)
				if paired == nil {
					return ErrInvalidChoice
				}
				if paired.DesignID == "arcane_focus" && len(support.Domains) == 1 {
					continue // The engine validates the chosen Affinity and single-Focus limit.
				}
				if operation != "actions/research-spell" || !category(s.definition(paired.DesignID), "element") || len(support.Domains) != 0 {
					return ErrInvalidChoice
				}
			}
		}
	case "actions/attempt-challenge", "actions/activate-ability", "actions/end-turn", "reactions/play-card", "reactions/defend", "reactions/pass":
	default:
		return ErrInvalidCommand
	}
	return nil
}
