package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davelawson/grim/api"
	randv2 "math/rand/v2"
	"slices"
)

var (
	ErrTiming         = errors.New("command is unavailable at this point")
	ErrPriority       = errors.New("another wizard has priority")
	ErrStaleChoice    = errors.New("choice is no longer pending")
	ErrFinished       = errors.New("match has finished")
	ErrInvalidCommand = errors.New("invalid command payload")
)

// RandomDraw is an operator diagnostic for one committed random decision.
// It is never persisted as history or included in participant projections.
type RandomDraw struct {
	Bound int `json:"bound"`
	Value int `json:"value"`
}

type PendingEffect struct {
	api.PendingEffectView
	Capacity map[string]int `json:"capacity,omitempty"`
	Focus    bool           `json:"focus"`
	Unbound  bool           `json:"unbound"`
	Stage    string         `json:"stage,omitempty"`
}

type Exchange struct {
	Participants []string         `json:"participants"`
	Priority     string           `json:"priority"`
	Passes       int              `json:"passes"`
	Stack        []*PendingEffect `json:"stack"`
}

type PhaseEffect struct {
	ID       string `json:"id"`
	Actor    string `json:"actor"`
	Source   string `json:"source"`
	Kind     string `json:"kind"`
	Optional bool   `json:"optional"`
}

func clone[T any](v T) T {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	var result T
	if err = json.Unmarshal(b, &result); err != nil {
		panic(err)
	}
	return result
}
func (s *State) nextID(kind string) string {
	s.Sequence++
	return fmt.Sprintf("%s/%s/%d", s.ID, kind, s.Sequence)
}
func (s *State) random(n int) (int, error) {
	if n <= 0 {
		return 0, errors.New("no random candidates")
	}
	if s.Random.Generator != "chacha8-v1" {
		return 0, errors.New("unsupported random generator")
	}
	g := randv2.NewChaCha8([32]byte{})
	if err := g.UnmarshalBinary(s.Random.State); err != nil {
		return 0, err
	}
	v := randv2.New(g).IntN(n)
	s.RandomDraws = append(s.RandomDraws, RandomDraw{Bound: n, Value: v})
	var err error
	s.Random.State, err = g.MarshalBinary()
	return v, err
}
func (s *State) shuffle(cards []CardInstance) error {
	for i := len(cards) - 1; i > 0; i-- {
		j, err := s.random(i + 1)
		if err != nil {
			return err
		}
		cards[i], cards[j] = cards[j], cards[i]
	}
	return nil
}

func (s *State) definition(id string) CardDesign { return s.definitions.byID[id] }
func category(d CardDesign, k string) bool       { return slices.Contains(d.Categories, k) }
func addStats(to, from map[string]int) {
	for k, v := range from {
		to[k] += v
	}
}
func meets(have, need map[string]int) bool {
	for k, v := range need {
		if have[k] < v {
			return false
		}
	}
	return true
}
func find(cards []CardInstance, id string) (int, *CardInstance) {
	for i := range cards {
		if cards[i].ID == id {
			return i, &cards[i]
		}
	}
	return -1, nil
}
func take(cards *[]CardInstance, id string) (CardInstance, bool) {
	i, c := find(*cards, id)
	if c == nil {
		return CardInstance{}, false
	}
	v := *c
	*cards = slices.Delete(*cards, i, i+1)
	return v, true
}
func (s *State) inPlay(id string) *CardInstance {
	for _, owner := range s.Participants {
		_, c := find(s.Wizards[owner].InPlay, id)
		if c != nil {
			return c
		}
	}
	return nil
}
func (s *State) expertise(c CardInstance) map[string]int {
	if c.Expertise != nil {
		return c.Expertise
	}
	return s.definition(c.DesignID).Expertise
}
func (s *State) affinities(owner string) map[string]int {
	result := map[string]int{}
	w := s.Wizards[owner]
	for _, c := range w.InPlay {
		addStats(result, s.definition(c.DesignID).Affinity)
	}
	for _, m := range w.Modifiers {
		addStats(result, m.Affinity)
	}
	for k, v := range result {
		result[k] = max(0, v)
	}
	return result
}
func (s *State) refresh() {
	for _, id := range s.Participants {
		s.Wizards[id].Affinities = s.affinities(id)
	}
}
func (s *State) maxHand(owner string) int {
	n := 5
	for _, c := range s.Wizards[owner].InPlay {
		if c.DesignID == "mnemonic_seal" {
			n++
		}
	}
	return n
}
func (s *State) neighbour(owner string, delta int) string {
	i := slices.Index(s.Ring, owner)
	if i < 0 || len(s.Ring) < 2 {
		return ""
	}
	return s.Ring[(i+delta+len(s.Ring))%len(s.Ring)]
}
func (s *State) alive(owner string) bool { w := s.Wizards[owner]; return w != nil && !w.Eliminated }
func (s *State) immune(owner string) bool {
	w := s.Wizards[owner]
	return w.ImmuneUntil > 0 && w.Turns <= w.ImmuneUntil
}
func (s *State) newCard(owner, design, zone string) CardInstance {
	return CardInstance{ID: s.nextID("card"), DesignID: design, Owner: owner, Location: zone, Ephemeral: design == "ephemeral_wis" || design == "ephemeral_wealth"}
}
func (s *State) gain(owner, design string) error {
	c := s.newCard(owner, design, "discard")
	if design == "pupil" {
		eligible := []string{}
		el := s.definition(s.Wizards[owner].Draft.Chantry).Element
		for _, d := range domains {
			if d != el {
				eligible = append(eligible, d)
			}
		}
		n, err := s.random(len(eligible))
		if err != nil {
			return err
		}
		c.Expertise = stats(eligible[n], 1)
	}
	s.Wizards[owner].Discard = append(s.Wizards[owner].Discard, c)
	return nil
}
func (s *State) dispose(c CardInstance, zone string) {
	if c.Ephemeral && zone == "discard" {
		return
	}
	c.Exhausted = false
	c.SkipRecovery = false
	c.DeployedTurn = 0
	c.Location = zone
	w := s.Wizards[c.Owner]
	if zone == "hand" {
		w.Hand = append(w.Hand, c)
	} else {
		w.Discard = append(w.Discard, c)
	}
}
func (s *State) removeAsset(id, zone string) {
	for _, owner := range s.Participants {
		if c, ok := take(&s.Wizards[owner].InPlay, id); ok {
			s.dispose(c, zone)
			s.refresh()
			return
		}
	}
}
func (s *State) draw(owner string, n int) error {
	w := s.Wizards[owner]
	for range n {
		if len(w.DrawPile) == 0 {
			w.DrawPile = w.Discard
			w.Discard = []CardInstance{}
			if err := s.shuffle(w.DrawPile); err != nil {
				return err
			}
			for i := range w.DrawPile {
				w.DrawPile[i].Location = "drawpile"
			}
		}
		if len(w.DrawPile) == 0 {
			break
		}
		c := w.DrawPile[0]
		w.DrawPile = w.DrawPile[1:]
		c.Location = "hand"
		w.Hand = append(w.Hand, c)
	}
	return nil
}
func (s *State) choice(actor, kind, source string, options []string, minCount, maxCount int, private bool) {
	s.Choice = &api.RequiredChoice{ID: s.nextID("choice"), Actor: actor, Type: kind, Source: source, Options: options, Minimum: minCount, Maximum: maxCount, Private: private}
}
func cardIDs(cards []CardInstance) []string {
	ids := []string{}
	for _, c := range cards {
		ids = append(ids, c.ID)
	}
	return ids
}
func (s *State) rumour(owner, id string) *api.RumourView {
	for i := range s.Wizards[owner].Rumours {
		r := &s.Wizards[owner].Rumours[i]
		if r.ID == id {
			return r
		}
	}
	return nil
}
func (s *State) drawRumour(owner string, capacity map[string]int) (string, error) {
	w := s.Wizards[owner]
	available := []string{}
	for _, d := range s.definitions.catalogue.Cards {
		if !category(d, "rumour") || !meets(capacity, d.Eligibility) {
			continue
		}
		busy := slices.Contains(w.Reserved, d.ID)
		for _, r := range w.Rumours {
			busy = busy || r.DesignID == d.ID
		}
		if !busy {
			available = append(available, d.ID)
		}
	}
	n, err := s.random(len(available))
	if err != nil {
		return "", err
	}
	return available[n], nil
}
func (s *State) playRumour(owner, design string) {
	w := s.Wizards[owner]
	d := s.definition(design)
	r := api.RumourView{ID: s.nextID("rumour"), DesignID: design, ExpiresAfterTurn: w.Turns + d.Duration - 1}
	w.Rumours = append(w.Rumours, r)
	switch design {
	case "rumour_solar_eclipse":
		s.schedule(owner, r.ID, "solar_eclipse", false)
	case "rumour_full_moon":
		s.schedule(owner, r.ID, "moonlit_hunt", true)
	}
}
func (s *State) schedule(actor, source, kind string, optional bool) {
	s.PhaseEffects = append(s.PhaseEffects, PhaseEffect{ID: s.nextID("phaseeffect"), Actor: actor, Source: source, Kind: kind, Optional: optional})
}
