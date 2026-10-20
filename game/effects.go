package game

import (
	"github.com/davelawson/grim/api"
	"slices"
)

func (s *State) validateTargets(e *PendingEffect) error {
	n := len(e.Targets)
	a := s.effect(e.Parent)
	opposing := func(id string, low bool) bool {
		if a == nil {
			return false
		}
		ids := a.Minions
		owner := a.Actor
		if e.Actor == a.Actor {
			ids = a.Defenders
			owner = a.Target
		}
		c := s.inPlay(id)
		return c != nil && c.Owner == owner && slices.Contains(ids, id) && (!low || expertise(*c)["might"] <= 1)
	}
	switch e.Kind {
	case "counterspell":
		if n != 1 || s.effect(e.Targets[0]) == nil || !s.effect(e.Targets[0]).Spell {
			return ErrInvalidChoice
		}
	case "deep_freeze":
		if n != 1 {
			return ErrInvalidChoice
		}
		c := s.inPlay(e.Targets[0])
		if c == nil || c.Exhausted || !category(definition(c.DesignID), "minion") || (c.Owner != s.neighbour(e.Actor, 1) && c.Owner != s.neighbour(e.Actor, -1)) {
			return ErrInvalidChoice
		}
	case "firewall", "dislocate", "deflect", "wildfire", "molten_armour":
		if n != 1 || !opposing(e.Targets[0], e.Kind == "molten_armour") {
			return ErrInvalidChoice
		}
	case "scalding_mist", "freezing_fog":
		if n < 1 || n > 2 || !opposing(e.Targets[0], true) {
			return ErrInvalidChoice
		}
		if n == 2 && !opposing(e.Targets[1], false) {
			return ErrInvalidChoice
		}
	case "moat", "quagmire":
		if n > 2 {
			return ErrInvalidChoice
		}
		for _, id := range e.Targets {
			if !opposing(id, false) {
				return ErrInvalidChoice
			}
		}
	default:
		if n != 0 {
			return ErrInvalidChoice
		}
	}
	return nil
}
func (s *State) resolve(e *PendingEffect) error {
	if e.Countered {
		if e.Kind == "subjugate_rival" {
			return s.attackFailure(e)
		}
		return nil
	}
	w := s.Wizards[e.Actor]
	switch e.Kind {
	case "defend":
		return nil
	case "challenge":
		return s.resolveChallenge(e)
	case "subjugate_rival":
		return s.resolveAttack(e)
	case "basic_wis", "basic_wealth", "ephemeral_wis", "ephemeral_wealth":
		addStats(w.Resources, definition(e.Kind).Resources)
	case "arcane_focus", "heal_activation", "purifying_steam":
		if s.alive(e.Actor) {
			w.Integrity = min(3, w.Integrity+1)
		}
	case "research_spell":
		if e.Mode == "element" {
			return s.gain(e.Actor, "element_"+e.Element)
		}
		pool := []string{}
		for _, d := range OpeningCatalogue().Cards {
			if category(d, "spell") && d.Pool == e.Element {
				pool = append(pool, d.ID)
			}
		}
		offers := []string{}
		for len(pool) > 0 && len(offers) < 3 {
			n, err := s.random(len(pool))
			if err != nil {
				return err
			}
			offers = append(offers, pool[n])
			pool = slices.Delete(pool, n, n+1)
		}
		e.Stage = "research"
		s.choice(e.Actor, "research", e.Source, offers, 1, 1, false)
	case "solar_eclipse":
		w.Resources["wis"] += 2
		w.Modifiers = append(w.Modifiers, api.Modifier{Source: e.Source, Affinity: map[string]int{"fire": 1, "infernal": 1}, ExpiresAfterTurn: w.Turns})
	case "burn_failure":
		assets := []string{}
		for _, c := range w.InPlay {
			if !category(definition(c.DesignID), "spell") {
				assets = append(assets, c.ID)
			}
		}
		if len(assets) > 0 {
			e.Stage = "burn_failure"
			s.choice(e.Actor, "asset_discard", e.Source, assets, 1, 1, false)
		}
	case "study", "flash_of_insight", "whispering_wind_trigger":
		n := 2
		if e.Kind == "whispering_wind_trigger" {
			n = 1
		}
		if err := s.draw(e.Actor, n); err != nil {
			return err
		}
		if len(w.Hand) > 0 {
			e.Stage = "discard"
			s.choice(e.Actor, "hand_discard", e.Source, cardIDs(w.Hand), 1, 1, true)
		}
	case "scrying_pool":
		n := min(3, len(w.DrawPile))
		if n > 0 {
			e.Stage = "scry"
			s.choice(e.Actor, "scry", e.Source, cardIDs(w.DrawPile[:n]), 1, 1, true)
			s.Choice.Cards = clone(w.DrawPile[:n])
		}
	case "deep_freeze":
		if c := s.inPlay(e.Targets[0]); c != nil && !c.Exhausted && (c.Owner == s.neighbour(e.Actor, 1) || c.Owner == s.neighbour(e.Actor, -1)) {
			c.Exhausted = true
			c.SkipRecovery = true
		}
	case "counterspell":
		if target := s.effect(e.Targets[0]); target != nil && target.Spell {
			target.Countered = true
		}
	case "flame_lash", "stoneward", "rockslide", "firewall", "dislocate", "deflect", "moat", "molten_armour", "wildfire", "scalding_mist", "sandstorm", "quagmire", "freezing_fog":
		s.resolveCombatSpell(e)
	default:
		// Durable ordinary plays are the common asset/lasting-Spell handler.
		d := definition(e.Kind)
		if d.Persistence != "durable" {
			return ErrInvalidChoice
		}
		for i, c := range e.Cards {
			if c.ID == e.Source {
				c.Location = "inplay"
				c.DeployedTurn = s.Turn.Number
				w.InPlay = append(w.InPlay, c)
				e.Cards = slices.Delete(e.Cards, i, i+1)
				break
			}
		}
	}
	return nil
}
func (s *State) resolveChallenge(e *PendingEffect) error {
	if e.ChallengeID != "investigate" && s.rumour(e.Actor, e.Source) == nil {
		return nil
	}
	ch := s.challengeDesign(e)
	if ch.ID == "" || !s.hasChallenger(e) || !meets(s.capacity(e), ch.Requirements) {
		return nil
	}
	if ch.ID == "investigate" {
		r, err := s.drawRumour(e.Actor, s.capacity(e))
		if err != nil {
			return err
		}
		s.Wizards[e.Actor].Reserved = append(s.Wizards[e.Actor].Reserved, r)
		return nil
	}
	for _, reward := range ch.Rewards {
		if err := s.gain(e.Actor, reward); err != nil {
			return err
		}
	}
	w := s.Wizards[e.Actor]
	if len(ch.AffinityReward) > 0 {
		w.Modifiers = append(w.Modifiers, api.Modifier{Source: e.Source, Affinity: cloneStats(ch.AffinityReward), ExpiresAfterTurn: w.Turns + 1})
	}
	w.Rumours = slices.DeleteFunc(w.Rumours, func(r api.RumourView) bool { return r.ID == e.Source })
	return nil
}
func (s *State) resolveAttack(e *PendingEffect) error {
	attack := s.affinities(e.Actor)["might"] + e.BonusAttack
	defence := s.affinities(e.Target)["might"] + e.BonusDefence
	survivor := false
	for _, id := range e.Minions {
		if c := s.inPlay(id); c != nil && c.Owner == e.Actor && expertise(*c)["might"] > 0 {
			survivor = true
			attack += max(0, expertise(*c)["might"]+e.MinionBonuses[id])
		}
	}
	for _, id := range e.Defenders {
		if c := s.inPlay(id); c != nil && c.Owner == e.Target {
			defence += max(0, expertise(*c)["might"]+e.MinionBonuses[id])
		}
	}
	for _, c := range s.Wizards[e.Target].InPlay {
		if c.DesignID == "rampart" {
			defence += 2
		}
	}
	if !survivor || !s.alive(e.Actor) || !s.alive(e.Target) || attack-defence < 2 {
		return s.attackFailure(e)
	}
	target := s.Wizards[e.Target]
	if target.Integrity > 0 {
		target.Integrity--
		target.ImmuneUntil = target.Turns + 1
		target.Kneel = true
		return nil
	}
	target.Eliminated = true
	target.Kneel = false
	s.Wizards[e.Actor].VictoryProgress["domination"]++
	s.transferChantry(e.Target, e.Actor)
	return nil
}
func (s *State) attackFailure(e *PendingEffect) error {
	options := []string{}
	for _, id := range e.Minions {
		if c := s.inPlay(id); c != nil && c.Owner == e.Actor {
			options = append(options, id)
		}
	}
	if len(options) > 0 {
		e.Stage = "attack_loss"
		s.choice(e.Actor, "attack_loss", e.Source, options, 1, 1, false)
	}
	return nil
}
func (s *State) transferChantry(from, to string) {
	id := s.Wizards[from].StartingChantry
	for _, owner := range s.Participants {
		w := s.Wizards[owner]
		for _, zone := range []*[]CardInstance{&w.InPlay, &w.Hand, &w.DrawPile, &w.Discard} {
			if c, ok := take(zone, id); ok {
				c.Owner = to
				s.dispose(c, "discard")
				return
			}
		}
	}
	if s.Exchange != nil {
		for _, e := range s.Exchange.Stack {
			if c, ok := take(&e.Cards, id); ok {
				c.Owner = to
				s.dispose(c, "discard")
				return
			}
		}
	}
}
func (s *State) resolveCombatSpell(e *PendingEffect) {
	a := s.effect(e.Parent)
	if a == nil {
		return
	}
	if a.MinionBonuses == nil {
		a.MinionBonuses = map[string]int{}
	}
	opposing := a.Minions
	if e.Actor == a.Actor {
		opposing = a.Defenders
	}
	survives := func(id string) bool { return slices.Contains(opposing, id) && s.inPlay(id) != nil }
	delay := func(id string) {
		if survives(id) {
			s.inPlay(id).SkipRecovery = true
		}
	}
	remove := func(id, zone string, low bool) {
		if survives(id) && (!low || expertise(*s.inPlay(id))["might"] <= 1) {
			s.removeAsset(id, zone)
		}
	}
	switch e.Kind {
	case "flame_lash":
		a.BonusAttack += 2
	case "stoneward":
		a.BonusDefence += 2
	case "rockslide":
		for _, id := range a.Minions {
			if s.inPlay(id) != nil {
				a.MinionBonuses[id]++
			}
		}
	case "firewall", "wildfire":
		remove(e.Targets[0], "discard", false)
	case "dislocate", "deflect":
		remove(e.Targets[0], "hand", false)
	case "moat":
		for _, id := range e.Targets {
			delay(id)
		}
	case "molten_armour":
		a.BonusDefence += 2
		remove(e.Targets[0], "discard", true)
	case "scalding_mist":
		remove(e.Targets[0], "discard", true)
		if len(e.Targets) > 1 {
			delay(e.Targets[1])
		}
	case "freezing_fog":
		remove(e.Targets[0], "hand", true)
		if len(e.Targets) > 1 {
			delay(e.Targets[1])
		}
	case "sandstorm":
		for _, id := range opposing {
			if survives(id) {
				a.MinionBonuses[id]--
			}
		}
	case "quagmire":
		for _, id := range e.Targets {
			if survives(id) {
				a.MinionBonuses[id]--
				delay(id)
			}
		}
	}
}
