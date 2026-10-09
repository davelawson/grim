package game

import (
	"errors"
	"fmt"
	"main/api"
	randv2 "math/rand/v2"
	"slices"
)

var (
	ErrDraftClosed    = errors.New("drafting is complete")
	ErrDraftSequence  = errors.New("choice is not the current draft step")
	ErrInvalidChoice  = errors.New("card is not offered for this choice")
	ErrNotParticipant = errors.New("actor is not a participant")
)

type Draft = api.Draft

type CardInstance = api.CardInstance

type Wizard struct {
	Draft           Draft          `json:"draft"`
	Integrity       int            `json:"integrity"`
	Affinities      map[string]int `json:"affinities"`
	VictoryProgress map[string]int `json:"victoryprogress"`
	InPlay          []CardInstance `json:"inplay"`
	Hand            []CardInstance `json:"hand"`
	DrawPile        []CardInstance `json:"drawpile"`
	Discard         []CardInstance `json:"discard"`
}

type Turn = api.Turn

type WizardView = api.WizardView

type PrivateView = api.PrivateView

// Choose returns a new state; even a failed automatic transition leaves the
// input and its random state intact. HTTP envelopes belong to the service.
func (s *State) Choose(actorID, kind, cardID string) (*State, error) {
	if !slices.Contains(s.Participants, actorID) {
		return nil, ErrNotParticipant
	}
	if s.Status != "setup" {
		return nil, ErrDraftClosed
	}
	wizard := s.Wizards[actorID]
	if wizard == nil {
		return nil, errors.New("missing wizard state")
	}
	offers := OpeningOffers()
	var step, next string
	var offered []string
	switch kind {
	case "choose_chantry":
		step, next, offered = "chantry", "victory_path", offers.Chantries
	case "choose_victory_path":
		step, next, offered = "victory_path", "minion", offers.VictoryPaths
	case "choose_minion":
		step, next, offered = "minion", "complete", offers.Minions[wizard.Draft.Chantry]
	default:
		return nil, ErrInvalidChoice
	}
	if wizard.Draft.Step != step {
		return nil, ErrDraftSequence
	}
	if !slices.Contains(offered, cardID) {
		return nil, ErrInvalidChoice
	}
	if s.RulesVersion != OpeningRulesVersion || s.CatalogueVersion != OpeningCatalogueVersion {
		return nil, ErrUnsupportedVersion
	}

	// Deep copy mutable wizard zones; immutable instance values are sufficient.
	state := *s
	state.Wizards = make(map[string]*Wizard, len(s.Wizards))
	for id, original := range s.Wizards {
		if original == nil {
			return nil, errors.New("missing wizard state")
		}
		copy := *original
		copy.Hand, copy.DrawPile = slices.Clone(original.Hand), slices.Clone(original.DrawPile)
		copy.InPlay, copy.Discard = slices.Clone(original.InPlay), slices.Clone(original.Discard)
		copy.Affinities, copy.VictoryProgress = cloneStats(original.Affinities), cloneStats(original.VictoryProgress)
		state.Wizards[id] = &copy
	}
	draft := &state.Wizards[actorID].Draft
	switch step {
	case "chantry":
		draft.Chantry = cardID
	case "victory_path":
		draft.VictoryPath = cardID
	case "minion":
		draft.Minion = cardID
	}
	draft.Step = next
	complete := true
	for _, id := range state.Participants {
		if state.Wizards[id] == nil {
			return nil, errors.New("missing wizard state")
		}
		complete = complete && state.Wizards[id].Draft.Step == "complete"
	}
	if complete {
		if err := state.start(); err != nil {
			return nil, err
		}
	}
	state.Revision++
	return &state, nil
}

func (s *State) start() error {
	if s.Random.Generator != "chacha8-v1" {
		return errors.New("unsupported random generator")
	}
	generator := randv2.NewChaCha8([32]byte{})
	if err := generator.UnmarshalBinary(s.Random.State); err != nil {
		return err
	}
	random := randv2.New(generator)
	s.Ring = slices.Clone(s.Participants)
	random.Shuffle(len(s.Ring), func(i, j int) { s.Ring[i], s.Ring[j] = s.Ring[j], s.Ring[i] })
	catalogue := OpeningCatalogue()
	// Sorted participant order fixes which random draws belong to each deck.
	for _, owner := range s.Participants {
		wizard := s.Wizards[owner]
		chantry, ok := catalogue.Card(wizard.Draft.Chantry)
		if !ok {
			return errors.New("missing selected Chantry")
		}
		wizard.Integrity = 3
		wizard.Affinities = cloneStats(chantry.Affinity)
		wizard.VictoryProgress = map[string]int{"domination": 0, "arcane": 0, "influence": 0, "council": 0, "fame": 0}
		newCard := func(design string, index int, location string) CardInstance {
			// IDs are assigned before shuffle, never derived from draw order.
			return CardInstance{ID: fmt.Sprintf("%s/%s/%02d", s.ID, owner, index), DesignID: design, Owner: owner, Location: location}
		}
		wizard.InPlay = []CardInstance{newCard(chantry.ID, 0, "inplay")}
		deckIDs := []string{"element_" + chantry.Element, "research_spell", "research_artifact", "basic_wis", "basic_wis", wizard.Draft.VictoryPath, "basic_wealth", wizard.Draft.Minion, "arcane_focus", "investigate"}
		deck := make([]CardInstance, len(deckIDs))
		for i, id := range deckIDs {
			deck[i] = newCard(id, i+1, "drawpile")
		}
		random.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
		wizard.Hand, wizard.DrawPile, wizard.Discard = deck[:5:5], deck[5:10:10], []CardInstance{}
		for i := range wizard.Hand {
			wizard.Hand[i].Location = "hand"
		}
	}
	state, err := generator.MarshalBinary()
	if err != nil {
		return err
	}
	s.Random.State = state
	s.Status = "playing"
	s.Turn = &Turn{Number: 1, ActiveWizard: s.Ring[0], Phase: "start", Step: "recovery"}
	return nil
}
