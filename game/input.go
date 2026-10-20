package game

import "github.com/davelawson/grim/api"

// Command is the engine input; HTTP clients use operation-specific api requests.
type Command struct {
	RequestID        string            `json:"requestid"`
	ExpectedRevision *int64            `json:"expectedrevision"`
	Type             string            `json:"type"`
	CardID           string            `json:"cardid,omitempty"`
	AbilityID        string            `json:"abilityid,omitempty"`
	ChallengeID      string            `json:"challengeid,omitempty"`
	Targets          []string          `json:"targets,omitempty"`
	Minions          []string          `json:"minions,omitempty"`
	Support          []api.CardSupport `json:"support,omitempty"`
	Payment          []string          `json:"payment,omitempty"`
	Discard          []string          `json:"discard,omitempty"`
	Mode             string            `json:"mode,omitempty"`
	Element          string            `json:"element,omitempty"`
	ChoiceID         string            `json:"choiceid,omitempty"`
	Selections       []string          `json:"selections,omitempty"`
	Order            []string          `json:"order,omitempty"`
}
