// Package game owns match state without dependencies on HTTP or storage.
package game

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	randv2 "math/rand/v2"
	"slices"
)

const SchemaVersion = 1

var ErrUnsupportedSchema = errors.New("unsupported match schema")

// State is limited to pre-draft lifecycle state. Rules and catalogue binding
// will be introduced with drafting, not inferred on load.
type State struct {
	SchemaVersion    int         `json:"schemaversion"`
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	LobbyID          string      `json:"lobbyid"`
	Status           string      `json:"status"`
	Revision         int64       `json:"revision"`
	Participants     []string    `json:"participants"`
	Ring             []string    `json:"ring"`
	Outcome          *string     `json:"outcome"`
	RulesVersion     *string     `json:"rulesversion"`
	CatalogueVersion *string     `json:"catalogueversion"`
	Random           RandomState `json:"random"`
}

type RandomState struct {
	Generator string `json:"generator"`
	Seed      []byte `json:"seed"`
	State     []byte `json:"state"`
}

// View is separate from State so new private fields cannot leak by default.
type View struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	LobbyID      string   `json:"lobbyid"`
	Status       string   `json:"status"`
	Revision     int64    `json:"revision"`
	Participants []string `json:"participants"`
	Ring         []string `json:"ring"`
	Outcome      *string  `json:"outcome"`
}

func NewSetup(id, name, lobbyID string, participants []string) (*State, error) {
	var seed [32]byte
	if _, err := rand.Read(seed[:]); err != nil {
		return nil, err
	}
	players := slices.Clone(participants)
	slices.Sort(players)
	ring := slices.Clone(players)
	generator := randv2.NewChaCha8(seed)
	randv2.New(generator).Shuffle(len(ring), func(i, j int) { ring[i], ring[j] = ring[j], ring[i] })
	randomState, err := generator.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &State{
		SchemaVersion: SchemaVersion, ID: id, Name: name, LobbyID: lobbyID,
		Status: "setup", Participants: players, Ring: ring,
		Random: RandomState{Generator: "chacha8-v1", Seed: seed[:], State: randomState},
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
	return &state, nil
}

func (s *State) View() *View {
	var outcome *string
	if s.Outcome != nil {
		value := *s.Outcome
		outcome = &value
	}
	return &View{ID: s.ID, Name: s.Name, LobbyID: s.LobbyID, Status: s.Status,
		Revision: s.Revision, Participants: slices.Clone(s.Participants),
		Ring: slices.Clone(s.Ring), Outcome: outcome}
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
