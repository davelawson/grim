package game

import (
	"errors"
	"github.com/davelawson/grim/api"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func table(t *testing.T, size int) *State {
	t.Helper()
	s := setup(t, size)
	for _, actor := range s.Participants {
		s = choose(t, s, actor, "choose_chantry", "chantry_fire")
		s = choose(t, s, actor, "choose_victory_path", "subjugate_rival")
		s = choose(t, s, actor, "choose_minion", "bastian_redhand")
	}
	s.Ring = slices.Clone(s.Participants)
	s.Turn = &Turn{Number: 10, ActiveWizard: s.Ring[0], Phase: "actions", Step: "actions"}
	s.Exchange = nil
	s.Choice = nil
	s.Resolving = nil
	s.PhaseEffects = nil
	for _, id := range s.Participants {
		w := s.Wizards[id]
		w.Turns = 1
		w.Hand = nil
		w.DrawPile = nil
		w.Discard = nil
		w.Rumours = nil
		w.Resources = map[string]int{"wis": 100, "wealth": 100}
		w.Modifiers = nil
		w.InPlay = w.InPlay[:1]
	}
	return s
}
func hand(s *State, owner, design string) string {
	c := s.newCard(owner, design, "hand")
	s.Wizards[owner].Hand = append(s.Wizards[owner].Hand, c)
	return c.ID
}
func asset(s *State, owner, design string) string {
	c := s.newCard(owner, design, "inplay")
	s.Wizards[owner].InPlay = append(s.Wizards[owner].InPlay, c)
	s.refresh()
	return c.ID
}
func act(t *testing.T, s *State, actor string, c Command) *State {
	t.Helper()
	next, err := s.Apply(actor, c)
	if err != nil {
		t.Fatalf("%s by %s (%s): %v; choice=%+v exchange=%+v", c.Type, actor, c.CardID, err, s.Choice, s.Exchange)
	}
	if next.Revision != s.Revision+1 {
		t.Fatal("revision not advanced once")
	}
	return next
}
func settle(t *testing.T, s *State) *State {
	t.Helper()
	for n := 0; s.Exchange != nil && s.Choice == nil && s.Status == "playing"; n++ {
		if n > 100 {
			t.Fatal("exchange did not finish")
		}
		s = act(t, s, s.Exchange.Priority, Command{Type: "pass"})
	}
	return s
}
func answer(t *testing.T, s *State, selections ...string) *State {
	t.Helper()
	if s.Choice == nil {
		t.Fatal("no required choice")
	}
	return act(t, s, s.Choice.Actor, Command{Type: "resolve_choice", ChoiceID: s.Choice.ID, Selections: selections})
}
func allAffinity(s *State, owner string) {
	m := map[string]int{}
	for _, d := range domains {
		m[d] = 10
	}
	s.Wizards[owner].Modifiers = append(s.Wizards[owner].Modifiers, api.Modifier{Source: "test", Affinity: m, ExpiresAfterTurn: 99})
	s.refresh()
}
func resumed(t *testing.T, s *State) *State {
	t.Helper()
	r, err := Decode(jsonBytes(t, s))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(jsonBytes(t, s), jsonBytes(t, r)) {
		t.Fatal("restart changed state")
	}
	return r
}
func reject(t *testing.T, s *State, actor string, c Command, want error) {
	t.Helper()
	before := jsonBytes(t, s)
	next, err := s.Apply(actor, c)
	if next != nil || !errors.Is(err, want) || !slices.Equal(before, jsonBytes(t, s)) {
		t.Fatalf("rejection changed state or returned wrong error: %v", err)
	}
}

func TestEveryCompletedSpell(t *testing.T) {
	for _, d := range OpeningCatalogue().Cards {
		if !category(d, "spell") {
			continue
		}
		t.Run(d.ID, func(t *testing.T) {
			s := table(t, 2)
			owner, opponent := s.Ring[0], s.Ring[1]
			allAffinity(s, owner)
			if d.Timing == "support" {
				source := hand(s, owner, "investigate")
				support := api.CardSupport{CardID: hand(s, owner, d.ID)}
				if d.ID == "keen_insight" {
					support.Domains = []string{"church"}
				}
				if d.ID == "flowing_form" {
					support.Domains = []string{"might", "church"}
				}
				minions := []string{}
				if d.ID != "unbound_will" {
					minions = append(minions, asset(s, owner, "bastian_redhand"))
				}
				s = act(t, s, owner, Command{Type: "attempt_challenge", CardID: source, ChallengeID: "investigate", Minions: minions, Support: []api.CardSupport{support}})
				s = settle(t, resumed(t, s))
				if len(s.Wizards[owner].Reserved) != 1 {
					t.Fatal("support failed to resolve")
				}
				return
			}
			reaction := strings.HasSuffix(d.Timing, "reaction")
			targets := []string{}
			if reaction {
				attacker, defender := owner, opponent
				if d.Timing != "attack_reaction" {
					attacker, defender = opponent, owner
					s.Turn.ActiveWizard = attacker
				}
				attackers := []string{asset(s, attacker, "caravan_guard"), asset(s, attacker, "caravan_guard")}
				source := hand(s, attacker, "subjugate_rival")
				if d.Timing == "spell_reaction" {
					source = hand(s, attacker, "study")
					s = act(t, s, attacker, Command{Type: "play_card", CardID: source})
					targets = []string{s.Exchange.Stack[0].ID}
				} else {
					s = act(t, s, attacker, Command{Type: "play_card", CardID: source, Minions: attackers})
					opposing := attackers
					if d.Timing == "attack_reaction" {
						defenders := []string{asset(s, defender, "caravan_guard"), asset(s, defender, "caravan_guard")}
						s = act(t, s, defender, Command{Type: "defend", Minions: defenders})
						opposing = defenders
					}
					switch d.ID {
					case "firewall", "dislocate", "deflect", "wildfire", "molten_armour":
						targets = opposing[:1]
					case "scalding_mist", "freezing_fog", "moat", "quagmire":
						targets = opposing
					}
				}
			} else if d.ID == "deep_freeze" {
				targets = []string{asset(s, opponent, "bastian_redhand")}
			}
			source := hand(s, owner, d.ID)
			command := Command{Type: "play_card", CardID: source, Targets: targets}
			if d.ID == "purifying_steam" {
				command.Discard = []string{hand(s, owner, "basic_wealth")}
			}
			if d.ID == "scrying_pool" {
				for _, id := range []string{"basic_wis", "basic_wealth", "arcane_focus"} {
					c := s.newCard(owner, id, "drawpile")
					s.Wizards[owner].DrawPile = append(s.Wizards[owner].DrawPile, c)
				}
			}
			s = act(t, s, owner, command)
			s = settle(t, resumed(t, s))
			if s.Choice != nil {
				q := s.Choice
				if q.Type == "scry" {
					s = act(t, s, owner, Command{Type: "resolve_choice", ChoiceID: q.ID, Selections: q.Options[:1], Order: []string{q.Options[2], q.Options[1]}})
				} else {
					s = answer(t, s, q.Options[0])
				}
				s = settle(t, s)
			}
			if s.Choice != nil && s.Choice.Type == "attack_loss" {
				s = answer(t, s, s.Choice.Options[0])
				s = settle(t, s)
			}
			if s.Exchange != nil || s.Choice != nil {
				t.Fatalf("unfinished spell: %+v", s.Choice)
			}
			if d.Persistence == "durable" {
				if s.inPlay(source) == nil {
					t.Fatal("lasting Spell did not enter play")
				}
			} else {
				_, card := find(s.Wizards[owner].Discard, source)
				if card == nil {
					t.Fatal("Spell did not discard")
				}
			}
		})
	}
}

func TestResearchOffersAndPrivateScryResume(t *testing.T) {
	s := table(t, 2)
	a, b := s.Ring[0], s.Ring[1]
	research := hand(s, a, "research_spell")
	fire := hand(s, a, "element_fire")
	water := hand(s, a, "element_water")
	s = act(t, s, a, Command{Type: "play_card", CardID: research, Mode: "discover", Support: []api.CardSupport{{CardID: fire}, {CardID: water}}})
	s = settle(t, s)
	if s.Choice == nil || s.Choice.Type != "research" || len(s.Choice.Options) != 2 {
		t.Fatal("wrong hybrid offer")
	}
	if !slices.Contains(s.Choice.Options, "scalding_mist") || !slices.Contains(s.Choice.Options, "purifying_steam") {
		t.Fatal("wrong pool")
	}
	reject(t, s, b, Command{Type: "resolve_choice", ChoiceID: s.Choice.ID, Selections: s.Choice.Options[:1]}, ErrPriority)
	reject(t, s, a, Command{Type: "end_turn"}, ErrTiming)
	one := answer(t, s, s.Choice.Options[0])
	two := answer(t, resumed(t, s), s.Choice.Options[0])
	if !reflect.DeepEqual(one, two) {
		t.Fatal("offer continuation changed")
	}
	s = one
	allAffinity(s, a)
	pool := hand(s, a, "scrying_pool")
	for _, id := range []string{"basic_wis", "basic_wealth"} {
		s.Wizards[a].DrawPile = append(s.Wizards[a].DrawPile, s.newCard(a, id, "drawpile"))
	}
	hand(s, a, "mnemonic_seal")
	s = act(t, s, a, Command{Type: "play_card", CardID: pool})
	s = settle(t, s)
	q := s.Choice
	if q == nil || len(q.Cards) != 2 {
		t.Fatal("scry reshuffled or missed cards")
	}
	view := s.View(b)
	if len(view.Choice.Cards) != 0 || len(view.Choice.Options) != 0 || view.Exchange.Resolving == nil {
		t.Fatal("private inspection leaked or resolution missing")
	}
	s = act(t, resumed(t, s), a, Command{Type: "resolve_choice", ChoiceID: q.ID, Selections: q.Options[:1], Order: q.Options[1:]})
	if len(s.Wizards[a].DrawPile) != 1 || s.Wizards[a].DrawPile[0].ID != q.Options[1] {
		t.Fatal("wrong returned order")
	}
}

func TestCounteredCostsAndCounterspellStack(t *testing.T) {
	s := table(t, 2)
	a, b := s.Ring[0], s.Ring[1]
	allAffinity(s, a)
	spring := asset(s, a, "spell_healing_spring")
	s.Wizards[a].Integrity = 1
	s = act(t, s, a, Command{Type: "activate_ability", CardID: spring, AbilityID: "heal"})
	activation := s.Exchange.Stack[0].ID
	counter := hand(s, b, "counterspell")
	s = act(t, s, b, Command{Type: "play_card", CardID: counter, Targets: []string{activation}})
	s = settle(t, resumed(t, s))
	if s.Wizards[a].Integrity != 1 || s.Wizards[a].Resources["wis"] != 97 || !s.inPlay(spring).Exhausted {
		t.Fatal("counter refunded activation")
	}
	steam := hand(s, a, "purifying_steam")
	discard := hand(s, a, "basic_wealth")
	s = act(t, s, a, Command{Type: "play_card", CardID: steam, Discard: []string{discard}})
	counter = hand(s, b, "counterspell")
	s = act(t, s, b, Command{Type: "play_card", CardID: counter, Targets: []string{s.Exchange.Stack[0].ID}})
	own := hand(s, a, "counterspell")
	s = act(t, s, a, Command{Type: "play_card", CardID: own, Targets: []string{s.Exchange.Stack[1].ID}})
	s = settle(t, s)
	if s.Wizards[a].Integrity != 2 {
		t.Fatal("countered Counterspell still resolved")
	}
	if _, c := find(s.Wizards[a].Discard, discard); c == nil {
		t.Fatal("additional cost refunded")
	}
}

func TestSupportCapacityAndDeployment(t *testing.T) {
	s := table(t, 2)
	a := s.Ring[0]
	s.Wizards[a].Modifiers = []api.Modifier{{Source: "test", Affinity: stats("earth", 1), ExpiresAfterTurn: 99}}
	r := api.RumourView{ID: s.nextID("rumour"), DesignID: "rumour_boundary_stones", ExpiresAfterTurn: 3}
	s.Wizards[a].Rumours = append(s.Wizards[a].Rumours, r)
	support := []api.CardSupport{{CardID: hand(s, a, "shape_stone")}, {CardID: hand(s, a, "shape_stone")}, {CardID: hand(s, a, "unbound_will")}}
	s = act(t, s, a, Command{Type: "attempt_challenge", CardID: r.ID, ChallengeID: "raise_stones", Support: support})
	if s.capacity(s.Exchange.Stack[0])["earth"] != 5 {
		t.Fatal("support did not stack")
	}
	s = settle(t, s)
	if len(s.Wizards[a].Rumours) != 0 {
		t.Fatal("Challenge did not succeed")
	}
	s = table(t, 2)
	a = s.Ring[0]
	minion := hand(s, a, "bastian_redhand")
	s = act(t, s, a, Command{Type: "play_card", CardID: minion})
	s = settle(t, s)
	attack := hand(s, a, "subjugate_rival")
	reject(t, s, a, Command{Type: "play_card", CardID: attack, Minions: []string{minion}}, ErrInvalidChoice)
	// A Support Spell's bonus cannot satisfy another Support Spell's requirement.
	investigate := hand(s, a, "investigate")
	keen := hand(s, a, "keen_insight")
	shape := hand(s, a, "shape_stone")
	s.inPlay(minion).DeployedTurn = 0
	s.Wizards[a].Modifiers = []api.Modifier{{Affinity: stats("air", 1), ExpiresAfterTurn: 99}}
	reject(t, s, a, Command{Type: "attempt_challenge", CardID: investigate, ChallengeID: "investigate", Minions: []string{minion}, Support: []api.CardSupport{{CardID: keen, Domains: []string{"earth"}}, {CardID: shape}}}, ErrInvalidChoice)
}

func TestEveryRumourChallengeAndReward(t *testing.T) {
	for _, d := range OpeningCatalogue().Cards {
		if !category(d, "rumour") {
			continue
		}
		for _, ch := range d.Challenges {
			t.Run(d.ID+"/"+ch.ID, func(t *testing.T) {
				s := table(t, 2)
				a := s.Ring[0]
				allAffinity(s, a)
				r := api.RumourView{ID: s.nextID("rumour"), DesignID: d.ID, ExpiresAfterTurn: s.Wizards[a].Turns + d.Duration - 1}
				s.Wizards[a].Rumours = []api.RumourView{r}
				minion := asset(s, a, "bastian_redhand")
				support := []api.CardSupport{}
				if ch.RequiresFocus {
					support = append(support, api.CardSupport{CardID: hand(s, a, "arcane_focus"), Domains: []string{"might"}})
				}
				if ch.ID == "moonlit_hunt" {
					s.Turn.Phase = "start"
					s.Turn.Step = "effects"
					s.schedule(a, r.ID, "moonlit_hunt", true)
					if err := s.advance(); err != nil {
						t.Fatal(err)
					}
					s = act(t, s, a, Command{Type: "resolve_choice", ChoiceID: s.Choice.ID, Selections: s.Choice.Options[:1], Minions: []string{minion}, Support: support})
				} else {
					s = act(t, s, a, Command{Type: "attempt_challenge", CardID: r.ID, ChallengeID: ch.ID, Minions: []string{minion}, Support: support})
				}
				s = settle(t, resumed(t, s))
				if len(s.Wizards[a].Rumours) != 0 {
					t.Fatal("successful rumour stayed in play")
				}
				counts := map[string]int{}
				for _, c := range s.Wizards[a].Discard {
					counts[c.DesignID]++
					if c.DesignID == "pupil" {
						if len(c.Expertise) != 1 || c.Expertise["fire"] != 0 {
							t.Fatal("invalid permanent Pupil Expertise")
						}
						again := resumed(t, s)
						_, copy := find(again.Wizards[a].Discard, c.ID)
						if !reflect.DeepEqual(copy.Expertise, c.Expertise) {
							t.Fatal("Pupil changed Expertise")
						}
					}
				}
				want := map[string]int{}
				for _, reward := range ch.Rewards {
					want[reward]++
				}
				for id, n := range want {
					if counts[id] != n {
						t.Fatalf("wrong reward %s: %v", id, counts)
					}
				}
				if len(ch.AffinityReward) > 0 {
					found := false
					for _, m := range s.Wizards[a].Modifiers {
						found = found || reflect.DeepEqual(m.Affinity, ch.AffinityReward) && m.ExpiresAfterTurn == 2
					}
					if !found {
						t.Fatal("missing two-turn reward")
					}
				}
			})
		}
	}
}

func TestFailedAttackAndRemovalPenalty(t *testing.T) {
	s := table(t, 2)
	a, b := s.Ring[0], s.Ring[1]
	s.Wizards[a].Modifiers = []api.Modifier{{Affinity: stats("might", 2), ExpiresAfterTurn: 99}}
	s.Wizards[b].Modifiers = []api.Modifier{{Affinity: map[string]int{"might": 2, "fire": 2}, ExpiresAfterTurn: 99}}
	first, second := asset(s, a, "bastian_redhand"), asset(s, a, "bastian_redhand")
	attack := hand(s, a, "subjugate_rival")
	s = act(t, s, a, Command{Type: "play_card", CardID: attack, Minions: []string{first, second}})
	firewall := hand(s, b, "firewall")
	s = act(t, s, b, Command{Type: "play_card", CardID: firewall, Targets: []string{first}})
	s = settle(t, s)
	if s.Choice == nil || s.Choice.Type != "attack_loss" || !slices.Equal(s.Choice.Options, []string{second}) {
		t.Fatal("missing failure loss")
	}
	reject(t, s, b, Command{Type: "pass"}, ErrTiming)
	s = answer(t, resumed(t, s), second)
	if s.Wizards[b].Integrity != 3 || s.inPlay(first) != nil || s.inPlay(second) != nil {
		t.Fatal("wrong failed Attack result")
	}
	again := hand(s, a, "subjugate_rival")
	minion := asset(s, a, "caravan_guard")
	reject(t, s, a, Command{Type: "play_card", CardID: again, Minions: []string{minion}}, ErrInvalidChoice)
}

func TestEliminationChantryTransferAndVictoryCheckpoints(t *testing.T) {
	for _, five := range []bool{false, true} {
		t.Run(map[bool]string{false: "last_survivor", true: "five_steps"}[five], func(t *testing.T) {
			s := table(t, 2)
			a, b := s.Ring[0], s.Ring[1]
			s.Wizards[b].Integrity = 0
			starting := s.Wizards[b].StartingChantry
			s.removeAsset(starting, "hand")
			asset(s, b, "chantry_earth")
			if five {
				s.Wizards[a].VictoryProgress["domination"] = 4
			}
			attack := hand(s, a, "subjugate_rival")
			guard := asset(s, a, "caravan_guard")
			s = act(t, s, a, Command{Type: "play_card", CardID: attack, Minions: []string{guard}})
			s = settle(t, s)
			if s.Status != "finished" || *s.Outcome != "victory" || !slices.Equal(s.Winners, []string{a}) || !s.Wizards[b].Eliminated {
				t.Fatal("missing terminal winner")
			}
			_, chantry := find(s.Wizards[a].Discard, starting)
			if chantry == nil || chantry.Owner != a {
				t.Fatal("starting Chantry was not transferred")
			}
			if len(s.Wizards[b].InPlay) != 1 {
				t.Fatal("additional Chantry transferred")
			}
			if five && len(s.Ring) != 2 {
				t.Fatal("ran Ring bookkeeping after five-step victory")
			}
			if !five && len(s.Ring) != 1 {
				t.Fatal("last-survivor check preceded Ring removal")
			}
			reject(t, s, a, Command{Type: "end_turn"}, ErrFinished)
			resumed(t, s)
		})
	}
	s := table(t, 2)
	for _, id := range s.Ring {
		s.Wizards[id].Eliminated = true
	}
	s.beginExchange()
	if err := s.finishExchange(); err != nil {
		t.Fatal(err)
	}
	if s.Status != "finished" || *s.Outcome != "draw" || len(s.Winners) != 0 {
		t.Fatal("zero survivors did not draw")
	}
}

func automatic(t *testing.T, s *State) *State {
	t.Helper()
	for i := 0; i < 100; i++ {
		if s.Status != "playing" {
			return s
		}
		if s.Exchange != nil && s.Choice == nil {
			s = settle(t, s)
			continue
		}
		if s.Choice == nil {
			return s
		}
		q := s.Choice
		switch q.Type {
		case "effects":
			selection := q.Options[0]
			for _, opt := range q.Options {
				if strings.HasPrefix(opt, "skip:") {
					selection = opt
					break
				}
			}
			s = answer(t, s, selection)
		case "kneel":
			s = answer(t, s, "decline")
		case "asset_discard", "attack_loss", "hand_discard", "research":
			s = answer(t, s, q.Options[0])
		case "end_discard":
			s = answer(t, s, q.Options...)
		case "maintenance":
			s = answer(t, s)
		default:
			t.Fatalf("unexpected automatic choice %s", q.Type)
		}
	}
	t.Fatal("automatic progress did not stop")
	return nil
}

func TestImmunityThroughNextTurnAndKneelingExpiry(t *testing.T) {
	s := table(t, 3)
	a, b, c := s.Ring[0], s.Ring[1], s.Ring[2]
	guard := asset(s, a, "caravan_guard")
	attack := hand(s, a, "subjugate_rival")
	s = act(t, s, a, Command{Type: "play_card", CardID: attack, Minions: []string{guard}})
	s = settle(t, s)
	if !s.immune(b) || s.Wizards[b].Integrity != 2 {
		t.Fatal("successful Attack did not grant immunity")
	}
	// Expiry-only Failure is due on B's upcoming forfeited turn.
	s.Wizards[b].Rumours = []api.RumourView{{ID: s.nextID("rumour"), DesignID: "rumour_burn_the_witch", ExpiresAfterTurn: 2}}
	delay := asset(s, b, "bastian_redhand")
	s.inPlay(delay).Exhausted = true
	s.inPlay(delay).SkipRecovery = true
	s = act(t, s, a, Command{Type: "end_turn"})
	s = answer(t, s, s.Choice.Options...)
	if s.Turn.ActiveWizard != b || s.Choice == nil || s.Choice.Type != "kneel" || !s.immune(b) {
		t.Fatal("missing kneeling checkpoint")
	}
	if !s.inPlay(delay).Exhausted || s.inPlay(delay).SkipRecovery {
		t.Fatal("kneeling did not consume recovery delay")
	}
	beforeHand := len(ordinary(s.Wizards[b].Hand))
	s = answer(t, resumed(t, s), "kneel")
	if s.Turn.ActiveWizard != c || !slices.Equal(s.Ring, []string{b, a, c}) || s.Wizards[a].VictoryProgress["domination"] != 1 || s.immune(b) || len(ordinary(s.Wizards[b].Hand)) != beforeHand || len(s.Wizards[b].Deferred) != 1 {
		t.Fatal("incorrect kneeling handoff or expiry")
	}
	if len(s.Wizards[b].Rumours) != 0 {
		t.Fatal("expired rumour stayed in play")
	}
	s = automatic(t, s)
	for s.Turn.ActiveWizard != b {
		s = act(t, s, s.Turn.ActiveWizard, Command{Type: "end_turn"})
		s = automatic(t, s)
	}
	if len(s.Wizards[b].Deferred) != 0 {
		t.Fatal("deferred failure was not consumed")
	}
	if s.inPlay(delay) == nil {
		t.Fatal("discarded wrong asset in deferred Failure")
	}
}

func TestUpkeepIncomeStacksAndResourceExpiry(t *testing.T) {
	s := table(t, 2)
	a, b := s.Ring[0], s.Ring[1]
	one, two := asset(s, a, "rampart"), asset(s, a, "rampart")
	asset(s, a, "verdant_spring")
	asset(s, a, "mnemonic_seal")
	asset(s, a, "mnemonic_seal")
	asset(s, a, "wind_road")
	asset(s, a, "wind_road")
	if s.maxHand(a) != 7 {
		t.Fatal("hand size did not stack")
	}
	s.Wizards[a].Resources = map[string]int{}
	s.Turn.Phase = "start"
	s.Turn.Step = "recovery"
	if err := s.advance(); err != nil {
		t.Fatal(err)
	}
	q := s.Choice
	if q == nil || q.Type != "maintenance" || len(q.Options) != 2 {
		t.Fatal("missing upkeep")
	}
	wis := ""
	wealth := ""
	for _, c := range s.Wizards[a].Hand {
		if c.DesignID == "ephemeral_wis" {
			wis = c.ID
		}
		if c.DesignID == "ephemeral_wealth" {
			wealth = c.ID
		}
	}
	if wis == "" || wealth == "" {
		t.Fatal("income not generated before upkeep")
	}
	reject(t, s, a, Command{Type: "resolve_choice", ChoiceID: q.ID, Selections: []string{one, two}, Payment: []string{wis}}, ErrInvalidChoice)
	s = act(t, s, a, Command{Type: "resolve_choice", ChoiceID: q.ID, Selections: []string{one}, Payment: []string{wis}})
	s = automatic(t, s)
	if s.inPlay(one) == nil || s.inPlay(two) != nil {
		t.Fatal("unpaid upkeep not discarded independently")
	}
	minion := hand(s, a, "bastian_redhand")
	before := s.Wizards[a].Resources["wis"]
	s = act(t, s, a, Command{Type: "play_card", CardID: minion})
	s = settle(t, s)
	if s.Wizards[a].Resources["wis"] != before {
		t.Fatal("stacked Wind Road underflowed")
	}
	s.Wizards[b].Resources["wis"] = 7
	s.Wizards[b].Hand = append(s.Wizards[b].Hand, s.newCard(b, "ephemeral_wis", "hand"))
	s = act(t, s, a, Command{Type: "end_turn"})
	if len(s.Wizards[b].Resources) != 0 {
		t.Fatal("responder resources survived cleanup")
	}
	for _, c := range s.Wizards[b].Hand {
		if c.Ephemeral {
			t.Fatal("responder ephemeral survived cleanup")
		}
	}
	if s.Choice == nil || s.Choice.Type != "end_discard" {
		t.Fatal("missing retention choice")
	}
}

func TestFixedResponseEligibilityAndPriorityReset(t *testing.T) {
	s := table(t, 4)
	a, b, c, d := s.Ring[0], s.Ring[1], s.Ring[2], s.Ring[3]
	source := hand(s, a, "study")
	s = act(t, s, a, Command{Type: "play_card", CardID: source})
	if !slices.Equal(s.Exchange.Participants, []string{a, b, d}) || s.Exchange.Priority != b {
		t.Fatal("wrong fixed eligibility")
	}
	counter := hand(s, c, "counterspell")
	reject(t, s, c, Command{Type: "play_card", CardID: counter, Targets: []string{s.Exchange.Stack[0].ID}}, ErrPriority)
	counter = hand(s, b, "counterspell")
	s = act(t, s, b, Command{Type: "play_card", CardID: counter, Targets: []string{s.Exchange.Stack[0].ID}})
	if s.Exchange.Priority != d || s.Exchange.Passes != 0 || slices.Contains(s.Exchange.Participants, c) {
		t.Fatal("Reaction expanded eligibility or failed to reset passes")
	}
	s = act(t, s, d, Command{Type: "pass"})
	s = act(t, s, a, Command{Type: "pass"})
	s = act(t, s, b, Command{Type: "pass"})
	if s.Exchange.Priority != a || s.Exchange.Passes != 0 || !s.Exchange.Stack[0].Countered {
		t.Fatal("resolution did not restore active Priority")
	}
	s = settle(t, s)
	if len(s.Wizards[a].Hand) != 0 {
		t.Fatal("countered Study drew cards")
	}
}

func TestRecoveryDelaysOverlapAndReturnClearsTemporaryState(t *testing.T) {
	s := table(t, 2)
	a, b := s.Ring[0], s.Ring[1]
	s.Wizards[a].Modifiers = []api.Modifier{{Affinity: stats("might", 1), ExpiresAfterTurn: 99}}
	s.Wizards[b].Modifiers = []api.Modifier{{Affinity: map[string]int{"earth": 1, "water": 1}, ExpiresAfterTurn: 99}}
	minion := asset(s, a, "caravan_guard")
	attack := hand(s, a, "subjugate_rival")
	s = act(t, s, a, Command{Type: "play_card", CardID: attack, Minions: []string{minion}})
	moat := hand(s, b, "moat")
	s = act(t, s, b, Command{Type: "play_card", CardID: moat, Targets: []string{minion}})
	s = act(t, s, a, Command{Type: "pass"})
	quagmire := hand(s, b, "quagmire")
	s = act(t, s, b, Command{Type: "play_card", CardID: quagmire, Targets: []string{minion}})
	s = settle(t, s)
	if s.Choice != nil {
		s = answer(t, s, s.Choice.Options[0])
	}
	// The surviving guard skips exactly one Recovery for both delays.
	if s.inPlay(minion) == nil {
		t.Fatal("unexpected Attack failure")
	}
	if !s.inPlay(minion).SkipRecovery {
		t.Fatal("missing recovery delay")
	}
	s.Turn.Phase = "start"
	s.Turn.Step = "recovery"
	if err := s.advance(); err != nil {
		t.Fatal(err)
	}
	if !s.inPlay(minion).Exhausted || s.inPlay(minion).SkipRecovery {
		t.Fatal("delays accumulated or did not skip Recovery")
	}
	s = automatic(t, s)
	s.Turn.Phase = "start"
	s.Turn.Step = "recovery"
	if err := s.advance(); err != nil {
		t.Fatal(err)
	}
	if s.inPlay(minion).Exhausted {
		t.Fatal("second Recovery was incorrectly skipped")
	}
	// Permanent instance traits survive return while temporary state does not.
	c := s.newCard(a, "pupil", "inplay")
	c.Expertise = stats("might", 1)
	c.Exhausted = true
	c.SkipRecovery = true
	s.Wizards[a].InPlay = append(s.Wizards[a].InPlay, c)
	s.removeAsset(c.ID, "hand")
	_, returned := find(s.Wizards[a].Hand, c.ID)
	if returned == nil || returned.Exhausted || returned.SkipRecovery || returned.Expertise["might"] != 1 {
		t.Fatal("return lost traits or retained temporary state")
	}
}

func TestResearchRecipesAndUnfilteredOffers(t *testing.T) {
	recipes := []struct {
		elements []string
		pool     string
		count    int
	}{{nil, "unaligned", 3}, {[]string{"fire"}, "fire", 3}, {[]string{"earth", "air"}, "air_earth", 2}}
	for _, recipe := range recipes {
		t.Run(recipe.pool, func(t *testing.T) {
			s := table(t, 2)
			a := s.Ring[0]
			source := hand(s, a, "research_spell")
			support := []api.CardSupport{}
			for _, element := range recipe.elements {
				support = append(support, api.CardSupport{CardID: hand(s, a, "element_"+element)})
			}
			s = act(t, s, a, Command{Type: "play_card", CardID: source, Mode: "discover", Support: support})
			s = settle(t, s)
			q := s.Choice
			if q == nil || len(q.Options) != recipe.count || !unique(q.Options) {
				t.Fatal("wrong distinct offer")
			}
			for _, id := range q.Options {
				if definition(id).Pool != recipe.pool {
					t.Fatal("wrong Research pool")
				}
			}
			s = answer(t, s, q.Options[0])
			if len(s.Wizards[a].Discard) < 2 {
				t.Fatal("Research or reward missing")
			}
		})
	}
	s := table(t, 2)
	a := s.Ring[0]
	source := hand(s, a, "research_spell")
	fire1, fire2 := hand(s, a, "element_fire"), hand(s, a, "element_fire")
	reject(t, s, a, Command{Type: "play_card", CardID: source, Mode: "discover", Support: []api.CardSupport{{CardID: fire1}, {CardID: fire2}}}, ErrInvalidChoice)
	s = act(t, s, a, Command{Type: "play_card", CardID: source, Mode: "element", Element: "water"})
	s = settle(t, s)
	found := false
	for _, c := range s.Wizards[a].Discard {
		found = found || c.DesignID == "element_water"
	}
	if !found {
		t.Fatal("could not learn an Element without standing Affinity")
	}
}

func TestRetentionRefillAndIndependentActivationRequirements(t *testing.T) {
	s := table(t, 2)
	a := s.Ring[0]
	keep := hand(s, a, "arcane_focus")
	drop := hand(s, a, "basic_wis")
	for range 6 {
		s.Wizards[a].Discard = append(s.Wizards[a].Discard, s.newCard(a, "basic_wealth", "discard"))
	}
	s = act(t, s, a, Command{Type: "end_turn"})
	q := s.Choice
	s = act(t, s, a, Command{Type: "resolve_choice", ChoiceID: q.ID, Selections: []string{drop}})
	if len(s.Wizards[a].Hand) != 5 {
		t.Fatal("did not refill retained hand")
	}
	if _, card := find(s.Wizards[a].Hand, keep); card == nil {
		t.Fatal("lost retained card")
	}
	s = table(t, 2)
	a = s.Ring[0]
	spring := asset(s, a, "spell_healing_spring")
	s.Wizards[a].Integrity = 2
	// Water was required to play the lasting Spell, not to activate its printed
	// pay-and-exhaust ability later after losing the source of that Affinity.
	s = act(t, s, a, Command{Type: "activate_ability", CardID: spring, AbilityID: "heal"})
	s = settle(t, s)
	if s.Wizards[a].Integrity != 3 {
		t.Fatal("activation reapplied the card's play requirements")
	}
}

func TestFullMoonOrderingAndFailureOnlyAtExpiry(t *testing.T) {
	s := table(t, 2)
	a := s.Ring[0]
	moon := api.RumourView{ID: s.nextID("rumour"), DesignID: "rumour_full_moon", ExpiresAfterTurn: 1}
	s.Wizards[a].Rumours = []api.RumourView{moon}
	s.Turn.Phase = "start"
	s.Turn.Step = "effects"
	wind := asset(s, a, "whispering_wind")
	s.schedule(a, wind, "whispering_wind_trigger", true)
	s.schedule(a, moon.ID, "moonlit_hunt", true)
	if err := s.advance(); err != nil {
		t.Fatal(err)
	}
	for _, e := range s.Choice.Effects {
		if e.Kind != "moonlit_hunt" {
			t.Fatal("Whispering Wind became available before the hunt decision")
		}
	}
	s = answer(t, s, s.Choice.Options[1])
	if len(s.Choice.Effects) != 1 || s.Choice.Effects[0].Kind != "whispering_wind_trigger" {
		t.Fatal("did not unlock later optional effects")
	}
	s = answer(t, s, s.Choice.Options[1])
	burn := api.RumourView{ID: s.nextID("rumour"), DesignID: "rumour_burn_the_witch", ExpiresAfterTurn: 3}
	s.Wizards[a].Rumours = []api.RumourView{burn}
	minion := asset(s, a, "bastian_redhand")
	s = act(t, s, a, Command{Type: "attempt_challenge", CardID: burn.ID, ChallengeID: "arrest_witch", Minions: []string{minion}})
	s.removeAsset(minion, "discard")
	s = settle(t, s)
	if s.Choice != nil || s.rumour(a, burn.ID) == nil {
		t.Fatal("failed attempt applied expiry penalty early")
	}
}

func TestRequiredChoiceRejectionsAndSnapshotValidation(t *testing.T) {
	s := table(t, 2)
	a := s.Ring[0]
	hand(s, a, "basic_wis")
	s = act(t, s, a, Command{Type: "end_turn"})
	q := s.Choice
	reject(t, s, a, Command{Type: "resolve_choice", ChoiceID: "old-choice"}, ErrStaleChoice)
	reject(t, s, a, Command{Type: "resolve_choice", ChoiceID: q.ID, Selections: []string{"unknown"}}, ErrInvalidChoice)
	reject(t, s, a, Command{Type: "resolve_choice", ChoiceID: q.ID, Selections: []string{q.Options[0], q.Options[0]}}, ErrInvalidChoice)
	reject(t, s, a, Command{Type: "resolve_choice", ChoiceID: q.ID, Mode: "destroy"}, ErrInvalidChoice)
	reject(t, s, a, Command{Type: "end_turn", Payment: []string{q.Options[0]}}, ErrInvalidCommand)
	s = answer(t, s)
	resumed(t, s)
	for _, corrupt := range []func(*State){
		func(s *State) { s.Turn = nil },
		func(s *State) { s.Ring = append(s.Ring, s.Ring[0]) },
		func(s *State) { s.Wizards[s.Ring[0]].Resources = nil },
		func(s *State) { s.Random.State = []byte("invalid") },
		func(s *State) {
			s.Exchange = &Exchange{Participants: []string{s.Ring[0]}, Stack: []*PendingEffect{nil}, Priority: s.Ring[0]}
		},
		func(s *State) {
			s.Choice = &api.RequiredChoice{ID: "test", Actor: s.Ring[0], Options: []string{"a"}, Minimum: 2, Maximum: 1}
		},
	} {
		bad := clone(s)
		corrupt(bad)
		if _, err := Decode(jsonBytes(t, bad)); err == nil {
			t.Fatal("corrupt continuation accepted")
		}
	}
	setup := setup(t, 2)
	setup.Abandon()
	resumed(t, setup)
}

func TestEliminatedActivePriorityAndSurvivorHandoff(t *testing.T) {
	s := table(t, 4)
	a, b, c, d := s.Ring[0], s.Ring[1], s.Ring[2], s.Ring[3]
	s.Turn.ActiveWizard = c
	source := hand(s, c, "study")
	s = act(t, s, c, Command{Type: "play_card", CardID: source})
	counter := hand(s, d, "counterspell")
	s = act(t, s, d, Command{Type: "play_card", CardID: counter, Targets: []string{s.Exchange.Stack[0].ID}})
	s.Wizards[c].Eliminated = true
	s = act(t, s, b, Command{Type: "pass"})
	s = act(t, s, d, Command{Type: "pass"})
	if s.Exchange.Priority != d || !slices.Contains(s.Ring, c) {
		t.Fatal("removed Ring position too early or started Priority before the former active position")
	}
	s = act(t, s, d, Command{Type: "pass"})
	s = act(t, s, b, Command{Type: "pass"})
	if s.Turn.ActiveWizard != d || !slices.Equal(s.Ring, []string{a, b, d}) {
		t.Fatal("wrong active-elimination handoff")
	}
	// With no eligible responders left, remaining effects finish automatically
	// instead of leaving an impossible Priority wait before the draw checkpoint.
	s = table(t, 2)
	a, b = s.Ring[0], s.Ring[1]
	source = hand(s, a, "study")
	s = act(t, s, a, Command{Type: "play_card", CardID: source})
	counter = hand(s, b, "counterspell")
	s = act(t, s, b, Command{Type: "play_card", CardID: counter, Targets: []string{s.Exchange.Stack[0].ID}})
	s.Wizards[a].Eliminated = true
	s.Wizards[b].Eliminated = true
	if err := s.pass(); err != nil {
		t.Fatal(err)
	}
	if s.Status != "finished" || *s.Outcome != "draw" || s.Exchange != nil {
		t.Fatal("zero-survivor exchange did not finish")
	}
}

func TestProgressTieBreakAndSharedVictory(t *testing.T) {
	for _, shared := range []bool{false, true} {
		s := table(t, 3)
		a, b := s.Ring[0], s.Ring[1]
		s.Wizards[a].VictoryProgress["domination"] = 5
		s.Wizards[b].VictoryProgress["domination"] = 5
		s.Wizards[a].VictoryProgress["arcane"] = 2
		s.Wizards[b].VictoryProgress["arcane"] = 1
		if shared {
			s.Wizards[b].VictoryProgress["arcane"] = 2
		}
		if !s.checkProgress() {
			t.Fatal("missing progress victory")
		}
		want := []string{a}
		if shared {
			want = append(want, b)
		}
		if !slices.Equal(s.Winners, want) {
			t.Fatal("wrong progress tie-break")
		}
	}
}
