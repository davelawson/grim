package game

import (
	"github.com/davelawson/grim/api"
	"slices"
	"strings"
)

func (s *State) advance() error {
	for s.Status == "playing" && s.Choice == nil && s.Exchange == nil {
		w := s.Wizards[s.Turn.ActiveWizard]
		owner := s.Turn.ActiveWizard
		switch s.Turn.Phase + "/" + s.Turn.Step {
		case "start/recovery":
			w.Turns++
			w.AttackUsed = false
			for i := range w.InPlay {
				c := &w.InPlay[i]
				if c.SkipRecovery {
					c.SkipRecovery = false
				} else {
					c.Exhausted = false
				}
			}
			s.Turn.Step = "income"
		case "start/income":
			for _, c := range w.InPlay {
				if c.DeployedTurn >= s.Turn.Number {
					continue
				}
				for _, resource := range []string{"ephemeral_wis", "ephemeral_wealth"} {
					for range definition(c.DesignID).Income[resource] {
						w.Hand = append(w.Hand, s.newCard(owner, resource, "hand"))
					}
				}
			}
			s.Turn.Step = "maintenance"
		case "start/maintenance":
			options := []string{}
			for _, c := range w.InPlay {
				if len(definition(c.DesignID).Upkeep) > 0 {
					options = append(options, c.ID)
				}
			}
			if len(options) > 0 {
				s.choice(owner, "maintenance", "", options, 0, len(options), false)
				return nil
			}
			s.Turn.Step = "bend_the_knee"
		case "start/bend_the_knee":
			if w.Kneel && len(s.Ring) > 2 {
				s.choice(owner, "kneel", "", []string{"kneel", "decline"}, 1, 1, false)
				return nil
			}
			w.Kneel = false
			s.Turn.Step = "rumour"
		case "start/rumour":
			design, err := s.drawRumour(owner, s.affinities(owner))
			if err != nil {
				return err
			}
			s.playRumour(owner, design)
			for _, r := range w.Reserved {
				s.playRumour(owner, r)
			}
			w.Reserved = nil
			s.PhaseEffects = append(s.PhaseEffects, w.Deferred...)
			w.Deferred = nil
			for _, c := range w.InPlay {
				if c.DesignID == "whispering_wind" {
					s.schedule(owner, c.ID, "whispering_wind_trigger", true)
				}
			}
			s.Turn.Step = "effects"
		case "start/effects", "end/effects":
			if len(s.PhaseEffects) > 0 {
				return s.offerEffects()
			}
			if s.Turn.Phase == "start" {
				s.Turn.Phase = "actions"
				s.Turn.Step = "actions"
			} else {
				s.Turn.Step = "cleanup"
			}
		case "actions/actions":
			return nil
		case "end/prepare_effects":
			for _, r := range w.Rumours {
				if r.ExpiresAfterTurn <= w.Turns && r.DesignID == "rumour_burn_the_witch" {
					s.schedule(owner, r.ID, "burn_failure", false)
				}
			}
			s.Turn.Step = "effects"
		case "end/cleanup":
			s.expire(false)
			s.Turn.Step = "discard"
		case "end/discard":
			minimum := max(0, len(w.Hand)-s.maxHand(owner))
			s.choice(owner, "end_discard", "", cardIDs(w.Hand), minimum, len(w.Hand), true)
			return nil
		case "end/draw":
			if err := s.draw(owner, max(0, s.maxHand(owner)-len(w.Hand))); err != nil {
				return err
			}
			s.startTurn(s.neighbour(owner, 1))
		default:
			return ErrTiming
		}
	}
	s.refresh()
	return nil
}
func (s *State) startTurn(owner string) {
	s.Turn = &Turn{Number: s.Turn.Number + 1, ActiveWizard: owner, Phase: "start", Step: "recovery"}
	s.PhaseEffects = nil
}
func (s *State) expire(skipped bool) {
	owner := s.Turn.ActiveWizard
	w := s.Wizards[owner]
	for _, id := range s.Participants {
		other := s.Wizards[id]
		other.Resources = map[string]int{}
		for _, zone := range []*[]CardInstance{&other.Hand, &other.DrawPile, &other.Discard, &other.InPlay} {
			*zone = slices.DeleteFunc(*zone, func(c CardInstance) bool { return c.Ephemeral })
		}
	}
	kept := []api.RumourView{}
	for _, r := range w.Rumours {
		if r.ExpiresAfterTurn > w.Turns {
			kept = append(kept, r)
			continue
		}
		if skipped && r.DesignID == "rumour_burn_the_witch" && !w.Eliminated {
			w.Deferred = append(w.Deferred, PhaseEffect{ID: s.nextID("phaseeffect"), Actor: owner, Source: r.ID, Kind: "burn_failure"})
		}
	}
	w.Rumours = kept
	w.Modifiers = slices.DeleteFunc(w.Modifiers, func(m api.Modifier) bool { return m.ExpiresAfterTurn <= w.Turns })
	if w.ImmuneUntil <= w.Turns {
		w.ImmuneUntil = 0
	}
	s.refresh()
}
func (s *State) offerEffects() error {
	// Mandatory effects can resolve in any chosen order. Full Moon precedes
	// every other optional effect, but does not constrain mandatory effects.
	if len(s.PhaseEffects) == 1 && !s.PhaseEffects[0].Optional {
		e := s.PhaseEffects[0]
		s.PhaseEffects = nil
		s.openPhaseEffect(e)
		return nil
	}
	moon := false
	for _, e := range s.PhaseEffects {
		moon = moon || e.Kind == "moonlit_hunt"
	}
	options := []string{}
	for _, e := range s.PhaseEffects {
		if e.Optional && moon && e.Kind != "moonlit_hunt" {
			continue
		}
		options = append(options, e.ID)
		if e.Optional {
			options = append(options, "skip:"+e.ID)
		}
	}
	s.choice(s.Turn.ActiveWizard, "effects", "", options, 1, 1, false)
	for _, e := range s.PhaseEffects {
		if slices.Contains(options, e.ID) {
			s.Choice.Effects = append(s.Choice.Effects, api.EffectOption{ID: e.ID, Kind: e.Kind, Source: e.Source, Optional: e.Optional})
		}
	}
	return nil
}
func (s *State) openPhaseEffect(e PhaseEffect) {
	s.open(&PendingEffect{PendingEffectView: api.PendingEffectView{ID: s.nextID("effect"), Actor: e.Actor, Kind: e.Kind, Source: e.Source, Spell: e.Kind == "whispering_wind_trigger"}})
}
func (s *State) beginExchange() {
	active := s.Turn.ActiveWizard
	participants := []string{}
	for _, id := range s.Ring {
		if s.alive(id) && (id == active || id == s.neighbour(active, 1) || id == s.neighbour(active, -1)) {
			participants = append(participants, id)
		}
	}
	s.Exchange = &Exchange{Participants: participants, Stack: []*PendingEffect{}}
}
func (s *State) eligible() []string {
	ids := []string{}
	if s.Exchange != nil {
		for _, id := range s.Exchange.Participants {
			if s.alive(id) {
				ids = append(ids, id)
			}
		}
	}
	return ids
}
func (s *State) nextPriority(after string) string {
	if s.Exchange == nil {
		return ""
	}
	ids := s.Exchange.Participants
	i := slices.Index(ids, after)
	for j := 1; j <= len(ids); j++ {
		id := ids[(i+j)%len(ids)]
		if s.alive(id) {
			return id
		}
	}
	return ""
}

func (s *State) open(e *PendingEffect) {
	if s.Exchange == nil {
		s.beginExchange()
	}
	s.Exchange.Stack = append(s.Exchange.Stack, e)
	s.Exchange.Passes = 0
	s.Exchange.Priority = s.nextPriority(e.Actor)
}
func (s *State) pass() error {
	x := s.Exchange
	x.Passes++
	if x.Passes < len(s.eligible()) {
		x.Priority = s.nextPriority(x.Priority)
		return nil
	}
	n := len(x.Stack) - 1
	e := x.Stack[n]
	x.Stack = x.Stack[:n]
	s.Resolving = e
	if err := s.resolve(e); err != nil {
		return err
	}
	if s.Choice == nil {
		return s.finishResolution()
	}
	x.Priority = ""
	return nil
}
func (s *State) finishResolution() error {
	e := s.Resolving
	for _, c := range e.Cards {
		s.dispose(c, "discard")
	}
	s.Resolving = nil
	s.refresh()
	if s.checkProgress() {
		return nil
	}
	if len(s.Exchange.Stack) == 0 {
		return s.finishExchange()
	}
	if len(s.eligible()) == 0 {
		return s.pass()
	}
	s.Exchange.Passes = 0
	s.Exchange.Priority = s.Turn.ActiveWizard
	if !s.alive(s.Exchange.Priority) {
		s.Exchange.Priority = s.nextPriority(s.Exchange.Priority)
	}
	return nil
}
func (s *State) finishExchange() error {
	old := slices.Clone(s.Ring)
	active := s.Turn.ActiveWizard
	s.Exchange = nil
	s.Ring = slices.DeleteFunc(s.Ring, func(id string) bool { return !s.alive(id) })
	if len(s.Ring) <= 1 {
		s.finish(slices.Clone(s.Ring))
		return nil
	}
	if !s.alive(active) {
		i := slices.Index(old, active)
		next := ""
		for j := 1; j <= len(old); j++ {
			id := old[(i+j)%len(old)]
			if s.alive(id) {
				next = id
				break
			}
		}
		s.expire(true)
		s.startTurn(next)
	}
	return s.advance()
}
func (s *State) finish(winners []string) {
	outcome := "victory"
	if len(winners) == 0 {
		outcome = "draw"
	}
	s.Status = "finished"
	s.Outcome = &outcome
	s.Winners = append([]string{}, winners...)
}
func (s *State) checkProgress() bool {
	contenders := []string{}
	vectors := map[string][]int{}
	for _, id := range s.Participants {
		if !s.alive(id) {
			continue
		}
		progress := s.Wizards[id].VictoryProgress
		winning := ""
		keys := []string{"domination", "arcane", "influence", "council", "fame"}
		for _, k := range keys {
			if progress[k] >= 5 {
				winning = k
				break
			}
		}
		if winning == "" {
			continue
		}
		v := []int{}
		for _, k := range keys {
			if k != winning {
				v = append(v, progress[k])
			}
		}
		slices.Sort(v)
		slices.Reverse(v)
		contenders = append(contenders, id)
		vectors[id] = v
	}
	if len(contenders) == 0 {
		return false
	}
	best := contenders[0]
	for _, id := range contenders {
		if slices.Compare(vectors[id], vectors[best]) > 0 {
			best = id
		}
	}
	winners := []string{}
	for _, id := range contenders {
		if slices.Equal(vectors[id], vectors[best]) {
			winners = append(winners, id)
		}
	}
	s.finish(winners)
	return true
}
func (s *State) kneel() error {
	owner := s.Turn.ActiveWizard
	aggressor := s.neighbour(owner, -1)
	next := s.neighbour(owner, 1)
	s.Wizards[owner].Kneel = false
	s.Wizards[aggressor].VictoryProgress["domination"]++
	a, b := slices.Index(s.Ring, owner), slices.Index(s.Ring, aggressor)
	s.Ring[a], s.Ring[b] = s.Ring[b], s.Ring[a]
	if s.checkProgress() {
		return nil
	}
	s.expire(true)
	s.startTurn(next)
	return nil
}
func skipEffect(id string) (string, bool) {
	return strings.TrimPrefix(id, "skip:"), strings.HasPrefix(id, "skip:")
}
