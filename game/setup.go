// Package game owns match state without dependencies on HTTP or storage.
package game

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/davelawson/grim/api"
	randv2 "math/rand/v2"
	"slices"
	"strings"
)

const SchemaVersion = 3

var ErrUnsupportedSchema = errors.New("unsupported match schema")
var ErrUnsupportedVersion = errors.New("unsupported rules or catalogue version")

// State is authoritative and private; only participant projections are served.
type State struct {
	RandomDraws      []RandomDraw        `json:"-"`
	Winners          []string            `json:"winners"`
	Sequence         int                 `json:"sequence"`
	Exchange         *Exchange           `json:"exchange"`
	Choice           *api.RequiredChoice `json:"choice"`
	Resolving        *PendingEffect      `json:"resolving"`
	PhaseEffects     []PhaseEffect       `json:"phaseeffects"`
	SchemaVersion    int                 `json:"schemaversion"`
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	LobbyID          string              `json:"lobbyid"`
	Status           string              `json:"status"`
	Revision         int64               `json:"revision"`
	Participants     []string            `json:"participants"`
	Ring             []string            `json:"ring"`
	Outcome          *string             `json:"outcome"`
	RulesVersion     string              `json:"rulesversion"`
	CatalogueVersion string              `json:"catalogueversion"`
	Wizards          map[string]*Wizard  `json:"wizards"`
	Turn             *Turn               `json:"turn"`
	Random           RandomState         `json:"random"`
}

type RandomState struct {
	Generator string `json:"generator"`
	Seed      []byte `json:"seed"`
	State     []byte `json:"state"`
}

// View is separate from State so new private fields cannot leak by default.
type View = api.MatchView

func NewSetup(id, name, lobbyID string, participants []string) (*State, error) {
	var seed [32]byte
	if _, err := rand.Read(seed[:]); err != nil {
		return nil, err
	}
	players := slices.Clone(participants)
	slices.Sort(players)
	generator := randv2.NewChaCha8(seed)
	randomState, err := generator.MarshalBinary()
	if err != nil {
		return nil, err
	}
	wizards := make(map[string]*Wizard, len(players))
	for _, player := range players {
		wizards[player] = &Wizard{Draft: Draft{Step: "chantry"}}
	}
	return &State{
		SchemaVersion: SchemaVersion, ID: id, Name: name, LobbyID: lobbyID,
		Status: "setup", Participants: players, Ring: []string{}, Wizards: wizards,
		RulesVersion: OpeningRulesVersion, CatalogueVersion: OpeningCatalogueVersion,
		Winners: []string{}, Random: RandomState{Generator: "chacha8-v1", Seed: seed[:], State: randomState},
	}, nil
}

func Decode(data []byte) (*State, error) {
	var header struct {
		SchemaVersion int `json:"schemaversion"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, err
	}
	if header.SchemaVersion != SchemaVersion {
		return nil, ErrUnsupportedSchema
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	if state.RulesVersion != OpeningRulesVersion || state.CatalogueVersion != OpeningCatalogueVersion {
		return nil, ErrUnsupportedVersion
	}
	if len(state.Participants) < 2 || len(state.Participants) > 4 || len(state.Wizards) != len(state.Participants) || !slices.IsSorted(state.Participants) {
		return nil, errors.New("invalid saved wizard roster")
	}
	for i, id := range state.Participants {
		if id == "" || state.Wizards[id] == nil || (i > 0 && id == state.Participants[i-1]) {
			return nil, errors.New("invalid saved wizard roster")
		}
	}
	if err := state.validateSavedState(); err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *State) View(actorID string) *View {
	var outcome *string
	if s.Outcome != nil {
		value := *s.Outcome
		outcome = &value
	}
	view := &View{ID: s.ID, Name: s.Name, LobbyID: s.LobbyID, Status: s.Status,
		Revision: s.Revision, Participants: slices.Clone(s.Participants),
		Ring: slices.Clone(s.Ring), Outcome: outcome, Winners: slices.Clone(s.Winners),
		RulesVersion: s.RulesVersion, CatalogueVersion: s.CatalogueVersion,
		Wizards: make(map[string]WizardView, len(s.Wizards))}
	if s.Turn != nil {
		turn := *s.Turn
		view.Turn = &turn
	}
	for id, wizard := range s.Wizards {
		public := WizardView{Turns: wizard.Turns, ImmuneUntil: wizard.ImmuneUntil, KneelEligible: wizard.Kneel && len(s.Ring) > 2, Reserved: slices.Clone(wizard.Reserved), Resources: cloneStats(wizard.Resources), Eliminated: wizard.Eliminated, Immune: s.immune(id), Rumours: slices.Clone(wizard.Rumours), Modifiers: clone(wizard.Modifiers), MaximumHandSize: s.maxHand(id), AttackUsed: wizard.AttackUsed, DraftComplete: wizard.Draft.Step == "complete"}
		if s.Turn != nil {
			integrity := wizard.Integrity
			public.Integrity = &integrity
			public.InPlay = slices.Clone(wizard.InPlay)
			public.Discard = slices.Clone(wizard.Discard)
			public.HandCount, public.DrawCount = len(wizard.Hand), len(wizard.DrawPile)
			public.Affinities = cloneStats(wizard.Affinities)
			public.VictoryProgress = cloneStats(wizard.VictoryProgress)
		}
		view.Wizards[id] = public
		if id == actorID {
			// Instance identity order is unrelated to shuffled draw order.
			draw := slices.Clone(wizard.DrawPile)
			slices.SortFunc(draw, func(a, b CardInstance) int { return strings.Compare(a.ID, b.ID) })
			view.You = &PrivateView{UserID: actorID, Draft: wizard.Draft, Hand: slices.Clone(wizard.Hand), DrawPile: draw}
		}
	}
	if s.Exchange != nil {
		view.Exchange = &api.ExchangeView{Participants: slices.Clone(s.Exchange.Participants), Priority: s.Exchange.Priority, Passes: s.Exchange.Passes, Stack: []api.PendingEffectView{}}
		for _, e := range s.Exchange.Stack {
			view.Exchange.Stack = append(view.Exchange.Stack, clone(e.PendingEffectView))
		}
	}
	if s.Resolving != nil && view.Exchange != nil {
		e := clone(s.Resolving.PendingEffectView)
		view.Exchange.Resolving = &e
	}
	if s.Choice != nil {
		view.Choice = clone(s.Choice)
		if view.Choice.Private && actorID != view.Choice.Actor {
			view.Choice.Options = nil
			view.Choice.Cards = nil
		}
	}
	// Deep-copy nested instance traits as well as zones and public continuations.
	return clone(view)
}

// Abandon preserves every existing terminal result.
func (s *State) Abandon() bool {
	if s.Status == "finished" {
		return false
	}
	outcome := "abandoned"
	s.Status, s.Outcome = "finished", &outcome
	s.Revision++
	return true
}
