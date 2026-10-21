// Package bot makes server-controlled gameplay decisions using only participant
// projections and printed card definitions. It has no HTTP or storage access.
package bot

import (
	"cmp"
	"slices"
	"strings"

	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
)

// Draft picks a Might-oriented opening, with stable catalogue-ID tie breaking.
func Draft(offer *api.DraftResponse, catalogue *api.Catalogue) (kind, id string) {
	cards := slices.Clone(offer.Cards)
	slices.SortFunc(cards, func(a, b api.CardDesign) int {
		score := func(d api.CardDesign) int {
			if offer.Draft.Step == "chantry" {
				best := 0
				for _, minion := range game.OpeningOffers().Minions[d.ID] {
					m, _ := catalogue.Card(minion)
					best = max(best, m.Expertise["might"]*10+m.Affinity["might"])
				}
				return best
			}
			return d.Expertise["might"]*10 + d.Affinity["might"]
		}
		if n := cmp.Compare(score(b), score(a)); n != 0 {
			return n
		}
		return strings.Compare(a.ID, b.ID)
	})
	if len(cards) == 0 {
		return "", ""
	}
	kind = map[string]string{"chantry": "choose_chantry", "victory_path": "choose_victory_path", "minion": "choose_minion"}[offer.Draft.Step]
	return kind, cards[0].ID
}

type planner struct {
	v     *api.MatchView
	c     *api.Catalogue
	actor string
	w     api.WizardView
	hand  []api.CardInstance
}

// Candidates returns preferred moves followed by safe fallbacks. The caller
// applies candidates to immutable engine states and accepts the first legal one.
// actions bounds proactive play within a single automated turn, including loops
// created by card effects. Required choices and response windows always resolve.
func Candidates(view *api.MatchView, catalogue *api.Catalogue, actions int) []game.Command {
	if view.You == nil || view.Status != "playing" {
		return nil
	}
	p := planner{v: view, c: catalogue, actor: view.You.UserID, hand: ordered(view.You.Hand)}
	p.w = view.Wizards[p.actor]
	if view.Choice != nil {
		return p.choices()
	}
	if view.Exchange != nil {
		return p.reactions()
	}
	if actions >= 64 {
		return []game.Command{{Type: "end_turn"}}
	}
	return p.actions()
}

func ordered(cards []api.CardInstance) []api.CardInstance {
	result := slices.Clone(cards)
	slices.SortFunc(result, func(a, b api.CardInstance) int { return strings.Compare(a.ID, b.ID) })
	return result
}
func (p planner) design(id string) api.CardDesign { d, _ := p.c.Card(id); return d }
func (p planner) expertise(card api.CardInstance, domain string) int {
	if card.Expertise != nil {
		return card.Expertise[domain]
	}
	return p.design(card.DesignID).Expertise[domain]
}
func (p planner) value(card api.CardInstance) int {
	d := p.design(card.DesignID)
	switch {
	case card.Ephemeral:
		return 0
	case card.DesignID == "subjugate_rival":
		return 100
	case card.DesignID == "arcane_focus":
		return 90
	case slices.Contains(d.Categories, "minion"):
		return 80 + p.expertise(card, "might")
	case slices.Contains(d.Categories, "location"):
		return 70
	case d.Persistence == "durable":
		return 60
	case d.Timing == "defence_reaction" || d.Timing == "combat_reaction":
		return 55
	case d.Timing == "attack_reaction":
		return 50
	case card.DesignID == "research_spell":
		return 40
	case len(d.Resources) > 0:
		return 20
	default:
		return 30
	}
}
func (p planner) ranked(cards []api.CardInstance, descending bool) []api.CardInstance {
	result := ordered(cards)
	slices.SortStableFunc(result, func(a, b api.CardInstance) int {
		if descending {
			return cmp.Compare(p.value(b), p.value(a))
		}
		return cmp.Compare(p.value(a), p.value(b))
	})
	return result
}
func (p planner) ready(combat bool) []string {
	ids := []string{}
	for _, card := range ordered(p.w.InPlay) {
		if !card.Exhausted && slices.Contains(p.design(card.DesignID).Categories, "minion") && (!combat || p.expertise(card, "might") > 0) && (combat && p.v.Exchange != nil || card.DeployedTurn < p.v.Turn.Number) {
			ids = append(ids, card.ID)
		}
	}
	return ids
}
func (p planner) focus() string {
	for _, card := range p.hand {
		if card.DesignID == "arcane_focus" {
			return card.ID
		}
	}
	return ""
}
func (p planner) payment(cost map[string]int, exclude []string) ([]string, bool) {
	have := map[string]int{}
	for k, v := range p.w.Resources {
		have[k] = v
	}
	ids := []string{}
	for _, card := range p.ranked(p.hand, false) {
		if slices.Contains(exclude, card.ID) {
			continue
		}
		resources := p.design(card.DesignID).Resources
		useful := false
		for _, k := range []string{"wis", "wealth"} {
			useful = useful || have[k] < cost[k] && resources[k] > 0
		}
		if useful {
			ids = append(ids, card.ID)
			for k, v := range resources {
				have[k] += v
			}
		}
	}
	for k, v := range cost {
		if have[k] < v {
			return nil, false
		}
	}
	return ids, true
}

func (p planner) funded(c game.Command, cost map[string]int) (game.Command, bool) {
	paid := map[string]int{}
	for k, v := range cost {
		paid[k] = v
	}
	exclude := []string{c.CardID}
	exclude = append(exclude, c.Discard...)
	mainDesign := ""
	for _, card := range p.hand {
		if card.ID == c.CardID {
			mainDesign = card.DesignID
			break
		}
	}
	focusWis := 0
	for _, s := range c.Support {
		exclude = append(exclude, s.CardID)
		for _, card := range p.hand {
			if card.ID != s.CardID {
				continue
			}
			if card.DesignID == "arcane_focus" {
				if mainDesign != "subjugate_rival" && c.Type != "defend" {
					focusWis++
				}
			} else {
				for k, v := range p.design(card.DesignID).Cost {
					paid[k] += v
				}
			}
			break
		}
	}
	paid["wis"] = max(0, paid["wis"]-focusWis)
	var ok bool
	c.Payment, ok = p.payment(paid, exclude)
	return c, ok
}

func (p planner) actions() []game.Command {
	result := []game.Command{}
	add := func(c game.Command, cost map[string]int) {
		if c, ok := p.funded(c, cost); ok {
			result = append(result, c)
		}
	}
	minions := p.ready(true)
	index := slices.Index(p.v.Ring, p.actor)
	if index < 0 {
		return []game.Command{{Type: "end_turn"}}
	}
	target := p.v.Wizards[p.v.Ring[(index+1)%len(p.v.Ring)]]
	attack := p.w.Affinities["might"]
	for _, card := range p.w.InPlay {
		if slices.Contains(minions, card.ID) {
			attack += p.expertise(card, "might")
		}
	}
	defence := target.Affinities["might"]
	for _, card := range target.InPlay {
		if !card.Exhausted {
			defence += p.expertise(card, "might")
		}
		if card.DesignID == "rampart" {
			defence += 2
		}
	}
	if p.focus() != "" {
		attack++
	}
	for _, card := range p.hand {
		if card.DesignID == "subjugate_rival" && len(minions) > 0 && !p.w.AttackUsed && !target.Immune && attack >= defence+2 {
			c := game.Command{Type: "play_card", CardID: card.ID, Minions: minions}
			if focus := p.focus(); focus != "" {
				c.Support = []api.CardSupport{{CardID: focus}}
			}
			add(c, p.design(card.DesignID).Cost)
		}
	}
	// Build lasting assets before spending ready Minions on Challenges.
	for _, card := range p.ranked(p.hand, true) {
		d := p.design(card.DesignID)
		if d.Persistence == "durable" && !d.Placeholder && !strings.HasSuffix(d.Timing, "reaction") {
			cost := map[string]int{}
			for k, v := range d.Cost {
				cost[k] = v
			}
			if slices.Contains(d.Categories, "minion") {
				for _, asset := range p.w.InPlay {
					if asset.DesignID == "wind_road" {
						cost["wis"] = max(0, cost["wis"]-1)
					}
				}
			}
			add(game.Command{Type: "play_card", CardID: card.ID}, cost)
			for _, el := range p.hand {
				if slices.Contains(p.design(el.DesignID).Categories, "element") {
					add(game.Command{Type: "play_card", CardID: card.ID, Support: []api.CardSupport{{CardID: el.ID}}}, cost)
				}
			}
		}
	}
	if p.w.Integrity != nil && *p.w.Integrity < 2 {
		if focus := p.focus(); focus != "" {
			add(game.Command{Type: "play_card", CardID: focus, Mode: "heal"}, map[string]int{"wis": 2})
		}
		for _, asset := range ordered(p.w.InPlay) {
			for _, a := range p.design(asset.DesignID).Abilities {
				if a.Effect == "heal_activation" && !asset.Exhausted {
					add(game.Command{Type: "activate_ability", CardID: asset.ID, AbilityID: a.ID}, a.Cost)
				}
			}
		}
	}
	for _, rumour := range p.w.Rumours {
		if rumour.DesignID == "rumour_full_moon" {
			continue
		}
		for _, ch := range p.design(rumour.DesignID).Challenges {
			result = append(result, p.challenge(rumour.ID, ch)...)
		}
	}
	for _, card := range p.hand {
		if card.DesignID == "research_spell" {
			for _, el := range p.hand {
				if slices.Contains(p.design(el.DesignID).Categories, "element") {
					add(game.Command{Type: "play_card", CardID: card.ID, Mode: "discover", Support: []api.CardSupport{{CardID: el.ID}}}, p.design(card.DesignID).Cost)
				}
			}
			add(game.Command{Type: "play_card", CardID: card.ID, Mode: "discover"}, p.design(card.DesignID).Cost)
		}
	}
	return append(result, game.Command{Type: "end_turn"})
}

func (p planner) challenge(source string, ch api.ChallengeDesign) []game.Command {
	result := []game.Command{}
	base := game.Command{Type: "attempt_challenge", CardID: source, ChallengeID: ch.ID, Minions: p.ready(false)}
	if c, ok := p.funded(base, ch.Cost); ok {
		result = append(result, c)
	}
	if focus := p.focus(); focus != "" {
		for _, domain := range []string{"might", "fire", "earth", "air", "water", "nobility", "church", "infernal", "popularity"} {
			c := base
			c.Support = []api.CardSupport{{CardID: focus, Domains: []string{domain}}}
			if c, ok := p.funded(c, ch.Cost); ok {
				result = append(result, c)
			}
		}
	}
	return result
}

func (p planner) reactions() []game.Command {
	result := []game.Command{}
	var attack *api.PendingEffectView
	for i := range p.v.Exchange.Stack {
		e := &p.v.Exchange.Stack[i]
		if e.Kind == "subjugate_rival" {
			attack = e
			break
		}
	}
	if attack != nil && attack.Target == p.actor {
		c := game.Command{Type: "defend", Minions: p.ready(true)}
		if !attack.DefenceFocus && p.focus() != "" {
			c.Support = []api.CardSupport{{CardID: p.focus()}}
		}
		if len(c.Minions)+len(c.Support) > 0 {
			result = append(result, c)
		}
	}
	// Untargeted combat spells and counters are enough for the basic strategy;
	// every other reaction can safely be declined.
	for _, card := range p.hand {
		d := p.design(card.DesignID)
		if !strings.HasSuffix(d.Timing, "reaction") {
			continue
		}
		c := game.Command{Type: "play_card", CardID: card.ID}
		if d.Timing == "spell_reaction" {
			for i := len(p.v.Exchange.Stack) - 1; i >= 0; i-- {
				e := p.v.Exchange.Stack[i]
				if e.Spell && e.Actor != p.actor {
					c.Targets = []string{e.ID}
					break
				}
			}
			if len(c.Targets) == 0 {
				continue
			}
		}
		if c, ok := p.funded(c, d.Cost); ok {
			result = append(result, c)
		}
	}
	return append(result, game.Command{Type: "pass"})
}

func (p planner) choices() []game.Command {
	q := p.v.Choice
	base := game.Command{Type: "resolve_choice", ChoiceID: q.ID}
	result := []game.Command{}
	switch q.Type {
	case "maintenance":
		cost := map[string]int{}
		keep := []string{}
		for _, card := range p.ranked(p.w.InPlay, true) {
			if !slices.Contains(q.Options, card.ID) {
				continue
			}
			next := map[string]int{}
			for k, v := range cost {
				next[k] = v
			}
			for k, v := range p.design(card.DesignID).Upkeep {
				next[k] += v
			}
			if _, ok := p.payment(next, nil); ok {
				cost = next
				keep = append(keep, card.ID)
			}
		}
		base.Selections = keep
		base.Payment, _ = p.payment(cost, nil)
	case "kneel":
		base.Selections = []string{"decline"}
		if p.w.Integrity != nil && *p.w.Integrity == 0 {
			base.Selections = []string{"kneel"}
		}
	case "effects":
		for _, option := range q.Options {
			if strings.HasPrefix(option, "skip:") {
				continue
			}
			c := base
			c.Selections = []string{option}
			hunt := false
			for _, effect := range q.Effects {
				if effect.ID == option && effect.Kind == "moonlit_hunt" {
					hunt = true
					for _, rumour := range p.w.Rumours {
						if rumour.ID == effect.Source {
							for _, ch := range p.design(rumour.DesignID).Challenges {
								for _, attempt := range p.challenge(rumour.ID, ch) {
									choice := c
									choice.Minions, choice.Support, choice.Payment = attempt.Minions, attempt.Support, attempt.Payment
									result = append(result, choice)
								}
							}
						}
					}
				}
			}
			if !hunt {
				result = append(result, c)
			}
		}
		for _, option := range q.Options {
			if strings.HasPrefix(option, "skip:") {
				c := base
				c.Selections = []string{option}
				result = append(result, c)
			}
		}
		return result
	case "end_discard":
		cards := p.ranked(p.hand, true)
		keep := map[string]bool{}
		// Leave room for fresh draws rather than retaining an entire inert hand.
		for _, card := range cards[:min(2, p.w.MaximumHandSize, len(cards))] {
			if p.value(card) >= 50 {
				keep[card.ID] = true
			}
		}
		base.Selections = []string{}
		for _, id := range q.Options {
			if !keep[id] {
				base.Selections = append(base.Selections, id)
			}
		}
	case "hand_discard", "asset_discard", "attack_loss":
		cards := p.hand
		if q.Type != "hand_discard" {
			cards = p.w.InPlay
		} else {
			base.Mode = "discard"
		}
		for _, card := range p.ranked(cards, false) {
			if slices.Contains(q.Options, card.ID) {
				base.Selections = []string{card.ID}
				break
			}
		}
	case "scry":
		cards := p.ranked(q.Cards, true)
		if len(cards) > 0 {
			base.Selections = []string{cards[0].ID}
			base.Order = []string{}
			for _, card := range cards[1:] {
				base.Order = append(base.Order, card.ID)
			}
		}
	case "research":
		cards := []api.CardInstance{}
		for _, id := range q.Options {
			cards = append(cards, api.CardInstance{ID: id, DesignID: id})
		}
		if cards = p.ranked(cards, true); len(cards) > 0 {
			base.Selections = []string{cards[0].ID}
		}
	default:
		return nil
	}
	return append(result, base)
}
