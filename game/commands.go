package game

import (
	"encoding/json"
	"github.com/davelawson/grim/api"
	"slices"
	"strings"
)

// ValidateCommand rejects fields inappropriate to a command, instead of silently
// accepting unused selections or payments. HTTP envelope checks live in match.
func ValidateCommand(c Command) error {
	allowed := map[string]string{
		"play_card":         "cardid mode element targets minions support payment discard",
		"attempt_challenge": "cardid challengeid minions support payment",
		"activate_ability":  "cardid abilityid support payment",
		"defend":            "minions support",
		"pass":              "",
		"end_turn":          "",
		"resolve_choice":    "choiceid selections order mode payment minions support",
	}
	fields, ok := allowed[c.Type]
	if !ok {
		return ErrInvalidCommand
	}
	data, _ := json.Marshal(c)
	var input map[string]any
	_ = json.Unmarshal(data, &input)
	for k := range input {
		if k == "type" || k == "requestid" || k == "expectedrevision" {
			continue
		}
		if !slices.Contains(strings.Fields(fields), k) {
			return ErrInvalidCommand
		}
	}
	if (c.Type == "play_card" || c.Type == "attempt_challenge" || c.Type == "activate_ability") && c.CardID == "" {
		return ErrInvalidCommand
	}
	if c.Type == "resolve_choice" && c.ChoiceID == "" {
		return ErrInvalidCommand
	}
	return nil
}
func (s *State) Apply(actor string, c Command) (*State, error) {
	if err := ValidateCommand(c); err != nil {
		return nil, err
	}
	if !slices.Contains(s.Participants, actor) {
		return nil, ErrNotParticipant
	}
	if !s.compatibleDefinitions() {
		return nil, ErrUnsupportedVersion
	}
	if s.Status == "finished" {
		return nil, ErrFinished
	}
	if s.Status != "playing" || (!s.alive(actor) && c.Type != "resolve_choice") {
		return nil, ErrTiming
	}
	next := s.clone()
	if err := next.apply(actor, c); err != nil {
		return nil, err
	}
	if err := next.advance(); err != nil {
		return nil, err
	}
	next.refresh()
	next.Revision++
	return next, nil
}
func (s *State) apply(actor string, c Command) error {
	if c.Type == "resolve_choice" {
		return s.resolveChoice(actor, c)
	}
	if s.Choice != nil || s.Resolving != nil {
		return ErrTiming
	}
	if s.Exchange != nil {
		if s.Exchange.Priority != actor {
			return ErrPriority
		}
		switch c.Type {
		case "pass":
			return s.pass()
		case "play_card":
			return s.declare(actor, c, false)
		case "defend":
			return s.defend(actor, c)
		default:
			return ErrTiming
		}
	}
	if actor != s.Turn.ActiveWizard || s.Turn.Phase != "actions" {
		return ErrTiming
	}
	switch c.Type {
	case "end_turn":
		s.Turn.Phase = "end"
		s.Turn.Step = "prepare_effects"
		return nil
	case "play_card", "attempt_challenge", "activate_ability":
		return s.declare(actor, c, false)
	default:
		return ErrTiming
	}
}
func (s *State) attack() *PendingEffect {
	if s.Exchange != nil {
		for _, e := range s.Exchange.Stack {
			if e.Kind == "subjugate_rival" {
				return e
			}
		}
	}
	return nil
}
func (s *State) effect(id string) *PendingEffect {
	if s.Exchange != nil {
		for _, e := range s.Exchange.Stack {
			if e.ID == id {
				return e
			}
		}
	}
	return nil
}
func unique(ids []string) bool {
	seen := map[string]bool{}
	for _, id := range ids {
		if id == "" || seen[id] {
			return false
		}
		seen[id] = true
	}
	return true
}

func (s *State) declare(actor string, c Command, startHunt bool) error {
	w := s.Wizards[actor]
	e := &PendingEffect{PendingEffectView: api.PendingEffectView{ID: s.nextID("effect"), Actor: actor, Source: c.CardID, Targets: slices.Clone(c.Targets), Minions: slices.Clone(c.Minions), Support: clone(c.Support), Mode: c.Mode, Element: c.Element, ChallengeID: c.ChallengeID, Cards: []CardInstance{}, MinionBonuses: map[string]int{}}, Capacity: map[string]int{}}
	d := CardDesign{}
	cost := map[string]int{}
	mainHand := false
	exhaustAbility := false
	challenge := false
	if c.Type == "attempt_challenge" {
		challenge = true
		e.Kind = "challenge"
		if r := s.rumour(actor, c.CardID); r != nil {
			d = s.definition(r.DesignID)
			if d.ID == "rumour_full_moon" && !startHunt {
				return ErrTiming
			}
		} else {
			_, card := find(w.Hand, c.CardID)
			if card == nil || card.DesignID != "investigate" {
				return ErrInvalidChoice
			}
			d = s.definition(card.DesignID)
			mainHand = true
		}
		var found bool
		for _, ch := range d.Challenges {
			if ch.ID == c.ChallengeID {
				e.ChallengeID = ch.ID
				addStats(cost, ch.Cost)
				found = true
			}
		}
		if !found {
			return ErrInvalidChoice
		}
	} else if c.Type == "activate_ability" {
		_, card := find(w.InPlay, c.CardID)
		if card == nil || card.DesignID != "spell_healing_spring" {
			return ErrInvalidChoice
		}
		d = s.definition(card.DesignID)
		found := false
		for _, ability := range d.Abilities {
			if ability.ID == c.AbilityID && ability.Timing == "action" && ability.Effect == "heal_activation" {
				if ability.Exhaust && card.Exhausted {
					return ErrInvalidChoice
				}
				e.Kind = string(ability.Effect)
				e.Spell = category(d, "spell")
				addStats(cost, ability.Cost)
				exhaustAbility = ability.Exhaust
				found = true
				break
			}
		}
		if !found {
			return ErrInvalidChoice
		}
	} else {
		_, card := find(w.Hand, c.CardID)
		if card == nil {
			return ErrInvalidChoice
		}
		d = s.definition(card.DesignID)
		if d.Placeholder || category(d, "element") || d.Timing == "support" || d.ID == "investigate" {
			return ErrInvalidChoice
		}
		mainHand = true
		e.Kind = d.ID
		e.Spell = category(d, "spell")
		addStats(cost, d.Cost)
	}
	reacting := s.Exchange != nil
	if reacting {
		if c.Type != "play_card" {
			return ErrTiming
		}
		a := s.attack()
		switch d.Timing {
		case "attack_reaction":
			if a == nil || a.Actor != actor {
				return ErrTiming
			}
			e.Parent = a.ID
		case "defence_reaction":
			if a == nil || a.Target != actor {
				return ErrTiming
			}
			e.Parent = a.ID
		case "combat_reaction":
			if a == nil || (a.Actor != actor && a.Target != actor) {
				return ErrTiming
			}
			e.Parent = a.ID
		case "spell_reaction":
			if len(c.Targets) != 1 {
				return ErrInvalidChoice
			}
			target := s.effect(c.Targets[0])
			if target == nil || !target.Spell {
				return ErrInvalidChoice
			}
		default:
			return ErrTiming
		}
	} else if strings.HasSuffix(d.Timing, "reaction") {
		return ErrTiming
	}
	if c.Type != "activate_ability" && c.AbilityID != "" {
		return ErrInvalidCommand
	}
	if !unique(c.Minions) || !unique(c.Payment) || !unique(c.Discard) || !unique(c.Targets) {
		return ErrInvalidChoice
	}
	used := []string{}
	if mainHand {
		used = append(used, c.CardID)
	}
	for _, x := range c.Support {
		used = append(used, x.CardID)
	}
	used = append(used, c.Payment...)
	used = append(used, c.Discard...)
	used = append(used, c.Minions...)
	if !unique(used) {
		return ErrInvalidChoice
	}
	supportBase := s.affinities(actor)
	supportCost := map[string]int{}
	focus := false
	elements := []string{}
	// Elements and Focus must be applied before checking any Support Spell's
	// own play requirements. Expertise and Spell bonuses never enable support.
	for _, x := range c.Support {
		_, card := find(w.Hand, x.CardID)
		if card == nil {
			return ErrInvalidChoice
		}
		sd := s.definition(card.DesignID)
		switch {
		case category(sd, "element"):
			if len(x.Domains) > 0 {
				return ErrInvalidChoice
			}
			supportBase[sd.Element]++
			e.Capacity[sd.Element]++
			elements = append(elements, sd.Element)
		case sd.ID == "arcane_focus":
			if focus {
				return ErrInvalidChoice
			}
			focus = true
			e.Focus = true
			if e.Kind == "subjugate_rival" {
				if len(x.Domains) > 0 {
					return ErrInvalidChoice
				}
				e.BonusAttack++
			} else {
				if len(x.Domains) != 1 || !slices.Contains(domains, x.Domains[0]) {
					return ErrInvalidChoice
				}
				supportBase[x.Domains[0]]++
				e.Capacity[x.Domains[0]]++
				supportCost["wis"]++
			}
		case sd.Timing == "support":
			if !challenge {
				return ErrInvalidChoice
			}
		default:
			return ErrInvalidChoice
		}
	}
	for _, x := range c.Support {
		_, card := find(w.Hand, x.CardID)
		sd := s.definition(card.DesignID)
		if sd.Timing != "support" {
			continue
		}
		if !meets(supportBase, sd.PlayRequirements) {
			return ErrInvalidChoice
		}
		addStats(cost, sd.Cost)
		addStats(e.Capacity, sd.Support)
		switch sd.ID {
		case "keen_insight":
			if len(x.Domains) != 1 || !slices.Contains(domains, x.Domains[0]) {
				return ErrInvalidChoice
			}
			e.Capacity[x.Domains[0]]++
		case "flowing_form":
			if len(x.Domains) != 2 || !unique(x.Domains) {
				return ErrInvalidChoice
			}
			for _, domain := range x.Domains {
				if !slices.Contains(domains, domain) {
					return ErrInvalidChoice
				}
				e.Capacity[domain]++
			}
		case "unbound_will":
			if len(x.Domains) > 0 {
				return ErrInvalidChoice
			}
			e.Unbound = true
		default:
			if len(x.Domains) > 0 {
				return ErrInvalidChoice
			}
		}
	}
	if c.Type != "activate_ability" && !meets(supportBase, d.PlayRequirements) {
		return ErrInvalidChoice
	}
	if e.Kind == "subjugate_rival" {
		if w.AttackUsed || len(c.Minions) == 0 {
			return ErrInvalidChoice
		}
		e.Target = s.neighbour(actor, 1)
		if e.Target == "" || !s.alive(e.Target) || s.immune(e.Target) {
			return ErrInvalidChoice
		}
	} else if !challenge && len(c.Minions) > 0 {
		return ErrInvalidChoice
	}
	for _, id := range c.Minions {
		_, card := find(w.InPlay, id)
		if card == nil || !category(s.definition(card.DesignID), "minion") || card.Exhausted || card.DeployedTurn >= s.Turn.Number {
			return ErrInvalidChoice
		}
		if e.Kind == "subjugate_rival" && s.expertise(*card)["might"] < 1 {
			return ErrInvalidChoice
		}
	}
	if challenge {
		if !focus && !e.Unbound && len(c.Minions) == 0 {
			return ErrInvalidChoice
		}
		ch := s.challengeDesign(e)
		if ch.RequiresFocus && !focus {
			return ErrInvalidChoice
		}
		if !meets(s.capacity(e), ch.Requirements) {
			return ErrInvalidChoice
		}
	}
	switch e.Kind {
	case "research_spell":
		if c.Mode == "element" {
			if len(elements) > 0 || !slices.Contains([]string{"fire", "earth", "air", "water"}, c.Element) {
				return ErrInvalidChoice
			}
		} else if c.Mode == "discover" {
			if c.Element != "" || len(elements) > 2 || !unique(elements) {
				return ErrInvalidChoice
			}
			slices.Sort(elements)
			e.Element = strings.Join(elements, "_")
			if len(elements) == 0 {
				e.Element = "unaligned"
			}
		} else {
			return ErrInvalidChoice
		}
	case "arcane_focus":
		if c.Mode != "heal" {
			return ErrInvalidChoice
		}
		cost["wis"] = 2
	default:
		if c.Mode != "" || c.Element != "" {
			return ErrInvalidChoice
		}
	}
	if e.Kind == "purifying_steam" {
		if len(c.Discard) != 1 {
			return ErrInvalidChoice
		}
	} else if len(c.Discard) > 0 {
		return ErrInvalidChoice
	}
	if err := s.validateTargets(e); err != nil {
		return err
	}
	// Wind Road reduces only the Wis play cost of Minions.
	if mainHand && category(d, "minion") {
		for _, asset := range w.InPlay {
			if asset.DesignID == "wind_road" {
				cost["wis"] = max(0, cost["wis"]-1)
			}
		}
	}
	for k, v := range supportCost {
		cost[k] = max(0, cost[k]-v)
	}
	payments, err := s.pay(actor, cost, c.Payment)
	if err != nil {
		return err
	}
	e.Payment = payments
	for _, id := range c.Discard {
		card, ok := take(&w.Hand, id)
		if !ok {
			return ErrInvalidChoice
		}
		e.Payment = append(e.Payment, card)
		s.dispose(card, "discard")
	}
	if mainHand {
		card, ok := take(&w.Hand, c.CardID)
		if !ok {
			return ErrInvalidChoice
		}
		card.Location = "stack"
		e.Cards = append(e.Cards, card)
	}
	for _, x := range c.Support {
		card, ok := take(&w.Hand, x.CardID)
		if !ok {
			return ErrInvalidChoice
		}
		card.Location = "stack"
		e.Cards = append(e.Cards, card)
	}
	for _, id := range c.Minions {
		_, card := find(w.InPlay, id)
		card.Exhausted = true
	}
	if exhaustAbility {
		_, card := find(w.InPlay, c.CardID)
		card.Exhausted = true
	}
	if e.Kind == "subjugate_rival" {
		w.AttackUsed = true
	}
	s.open(e)
	return nil
}
func (s *State) pay(actor string, cost map[string]int, ids []string) ([]CardInstance, error) {
	w := s.Wizards[actor]
	if w.Resources == nil {
		w.Resources = map[string]int{}
	}
	if !unique(ids) {
		return nil, ErrInvalidChoice
	}
	available := cloneStats(w.Resources)
	payments := []CardInstance{}
	for _, id := range ids {
		_, card := find(w.Hand, id)
		if card == nil || !category(s.definition(card.DesignID), "resource") {
			return nil, ErrInvalidChoice
		}
		addStats(available, s.definition(card.DesignID).Resources)
	}
	if !meets(available, cost) {
		return nil, ErrInvalidChoice
	}
	for _, id := range ids {
		card, _ := take(&w.Hand, id)
		payments = append(payments, card)
		s.dispose(card, "discard")
	}
	for k, v := range cost {
		available[k] -= v
	}
	w.Resources = available
	return payments, nil
}
func (s *State) challengeDesign(e *PendingEffect) api.ChallengeDesign {
	design := "investigate"
	if r := s.rumour(e.Actor, e.Source); r != nil {
		design = r.DesignID
	}
	for _, ch := range s.definition(design).Challenges {
		if ch.ID == e.ChallengeID {
			return ch
		}
	}
	return api.ChallengeDesign{}
}
func (s *State) capacity(e *PendingEffect) map[string]int {
	result := s.affinities(e.Actor)
	addStats(result, e.Capacity)
	for _, id := range e.Minions {
		if c := s.inPlay(id); c != nil {
			addStats(result, s.expertise(*c))
		}
	}
	return result
}
func (s *State) hasChallenger(e *PendingEffect) bool {
	if e.Focus || e.Unbound {
		return true
	}
	for _, id := range e.Minions {
		if c := s.inPlay(id); c != nil && c.Owner == e.Actor {
			return true
		}
	}
	return false
}
func (s *State) defend(actor string, c Command) error {
	a := s.attack()
	if a == nil || a.Target != actor {
		return ErrTiming
	}
	if len(c.Minions) == 0 && len(c.Support) == 0 {
		return ErrInvalidChoice
	}
	if !unique(c.Minions) {
		return ErrInvalidChoice
	}
	w := s.Wizards[actor]
	for _, id := range c.Minions {
		_, card := find(w.InPlay, id)
		if card == nil || card.Exhausted || s.expertise(*card)["might"] < 1 {
			return ErrInvalidChoice
		}
	}
	if len(c.Support) > 1 || a.DefenceFocus && len(c.Support) > 0 {
		return ErrInvalidChoice
	}
	cards := []CardInstance{}
	if len(c.Support) == 1 {
		x := c.Support[0]
		_, card := find(w.Hand, x.CardID)
		if card == nil || card.DesignID != "arcane_focus" || len(x.Domains) != 0 {
			return ErrInvalidChoice
		}
		v, _ := take(&w.Hand, x.CardID)
		v.Location = "stack"
		cards = append(cards, v)
		a.DefenceFocus = true
		a.BonusDefence++
	}
	for _, id := range c.Minions {
		_, card := find(w.InPlay, id)
		card.Exhausted = true
		a.Defenders = append(a.Defenders, id)
	}
	s.open(&PendingEffect{PendingEffectView: api.PendingEffectView{ID: s.nextID("effect"), Actor: actor, Kind: "defend", Source: a.ID, Minions: slices.Clone(c.Minions), Cards: cards, Support: clone(c.Support), Parent: a.ID}})
	return nil
}
