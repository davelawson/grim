package server

import (
	"encoding/json"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"strings"
	"testing"
)

func typedGameplay(t *testing.T, view *game.View, c game.Command) (string, any) {
	t.Helper()
	operation := ""
	switch c.Type {
	case "play_card":
		operation = "actions/play-card"
		if view.Exchange != nil && view.Choice == nil {
			operation = "reactions/play-card"
		} else {
			for _, card := range view.You.Hand {
				if card.ID == c.CardID {
					switch card.DesignID {
					case "subjugate_rival":
						operation = "actions/attack"
					case "research_spell":
						operation = "actions/research-spell"
						if c.Mode == "element" {
							operation = "actions/research-element"
						}
					case "arcane_focus":
						operation = "actions/heal-with-focus"
					}
				}
			}
		}
	case "attempt_challenge":
		operation = "actions/attempt-challenge"
	case "activate_ability":
		operation = "actions/activate-ability"
	case "end_turn":
		operation = "actions/end-turn"
	case "pass":
		operation = "reactions/pass"
	case "defend":
		operation = "reactions/defend"
	case "resolve_choice":
		q := view.Choice
		if q == nil {
			t.Fatal("strategy missing choice")
		}
		switch q.Type {
		case "asset_discard", "attack_loss":
			operation = "choices/asset-loss"
		case "end_discard":
			operation = "choices/end-discard"
		case "hand_discard":
			operation = "choices/hand-discard"
		case "effects":
			operation = "choices/resolve-effect"
			if strings.HasPrefix(c.Selections[0], "skip:") {
				operation = "choices/skip-effect"
			} else {
				for _, effect := range q.Effects {
					if effect.ID == c.Selections[0] && effect.Kind == "moonlit_hunt" {
						operation = "choices/undertake-hunt"
					}
				}
			}
		default:
			operation = "choices/" + q.Type
		}
	}
	if c.Selections == nil {
		c.Selections = []string{}
	}
	if c.Order == nil {
		c.Order = []string{}
	}
	first := ""
	if len(c.Selections) > 0 {
		first = c.Selections[0]
	}
	disposition := c.Mode
	if disposition == "" {
		disposition = "discard"
	}
	kneel := first == "kneel"
	elements := []string{}
	for _, e := range c.Support {
		elements = append(elements, e.CardID)
	}
	switch operation {
	case "actions/play-card":
		return api.PlayCardPath(view.ID), api.PlayCardRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: c.CardID, Targets: c.Targets, Support: c.Support, Payment: c.Payment, Discard: c.Discard}
	case "actions/attack":
		return api.AttackPath(view.ID), api.AttackRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: c.CardID, Minions: c.Minions, Support: c.Support, Payment: c.Payment}
	case "actions/attempt-challenge":
		return api.AttemptChallengePath(view.ID), api.AttemptChallengeRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: c.CardID, ChallengeID: c.ChallengeID, Minions: c.Minions, Support: c.Support, Payment: c.Payment}
	case "actions/activate-ability":
		return api.ActivateAbilityPath(view.ID), api.ActivateAbilityRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: c.CardID, AbilityID: c.AbilityID, Support: c.Support, Payment: c.Payment}
	case "actions/research-element":
		return api.ResearchElementPath(view.ID), api.ResearchElementRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: c.CardID, Element: c.Element}
	case "actions/research-spell":
		return api.ResearchSpellPath(view.ID), api.ResearchSpellRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: c.CardID, Elements: elements}
	case "actions/heal-with-focus":
		return api.HealWithFocusPath(view.ID), api.HealWithFocusRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: c.CardID, Support: c.Support, Payment: c.Payment}
	case "actions/end-turn":
		return api.EndTurnPath(view.ID), api.EndTurnRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision}
	case "choices/maintenance":
		return api.MaintenancePath(view.ID, c.ChoiceID), api.MaintenanceRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, Keep: c.Selections, Payment: c.Payment}
	case "choices/kneel":
		return api.KneelPath(view.ID, c.ChoiceID), api.KneelRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, Kneel: &kneel}
	case "choices/resolve-effect":
		return api.ResolveEffectPath(view.ID, c.ChoiceID), api.ResolveEffectRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, EffectID: strings.TrimPrefix(first, "skip:")}
	case "choices/skip-effect":
		return api.SkipEffectPath(view.ID, c.ChoiceID), api.SkipEffectRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, EffectID: strings.TrimPrefix(first, "skip:")}
	case "choices/undertake-hunt":
		return api.UndertakeHuntPath(view.ID, c.ChoiceID), api.UndertakeHuntRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, EffectID: strings.TrimPrefix(first, "skip:"), Minions: c.Minions, Support: c.Support, Payment: c.Payment}
	case "choices/research":
		return api.ChooseResearchPath(view.ID, c.ChoiceID), api.ChooseResearchRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, DesignID: first}
	case "choices/end-discard":
		return api.EndDiscardPath(view.ID, c.ChoiceID), api.EndDiscardRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, Cards: c.Selections}
	case "choices/hand-discard":
		return api.HandDiscardPath(view.ID, c.ChoiceID), api.HandDiscardRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: first, Disposition: disposition}
	case "choices/asset-loss":
		return api.AssetLossPath(view.ID, c.ChoiceID), api.AssetLossRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: first}
	case "choices/scry":
		return api.ScryPath(view.ID, c.ChoiceID), api.ScryRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, Take: first, Order: c.Order}
	case "reactions/play-card":
		return api.ReactPath(view.ID), api.ReactRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, CardID: c.CardID, Targets: c.Targets, Support: c.Support, Payment: c.Payment}
	case "reactions/defend":
		return api.DefendPath(view.ID), api.DefendRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision, Minions: c.Minions, Support: c.Support}
	case "reactions/pass":
		return api.PassPath(view.ID), api.PassRequest{RequestID: c.RequestID, ExpectedRevision: c.ExpectedRevision}
	}
	t.Fatalf("unsupported test intent %s", operation)
	return "", nil
}
func gameplayBody(t *testing.T, input any) string {
	t.Helper()
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
