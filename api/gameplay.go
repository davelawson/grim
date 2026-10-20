package api

// PlayCardRequest declares an ordinary card play. Attack, Research, and Focus healing have dedicated requests.
type PlayCardRequest struct {
	RequestID        string        `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64        `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string        `json:"cardid" binding:"required"`
	Targets          []string      `json:"targets,omitempty"`
	Support          []CardSupport `json:"support,omitempty"`
	Payment          []string      `json:"payment,omitempty"`
	Discard          []string      `json:"discard,omitempty"`
}

// AttackRequest commits an Attack card and ready Minions.
type AttackRequest struct {
	RequestID        string        `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64        `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string        `json:"cardid" binding:"required"`
	Minions          []string      `json:"minions" binding:"required"`
	Support          []CardSupport `json:"support,omitempty"`
	Payment          []string      `json:"payment,omitempty"`
}

// AttemptChallengeRequest attempts a printed Challenge on Investigate or a rumour.
type AttemptChallengeRequest struct {
	RequestID        string        `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64        `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string        `json:"cardid" binding:"required"`
	ChallengeID      string        `json:"challengeid" binding:"required"`
	Minions          []string      `json:"minions,omitempty"`
	Support          []CardSupport `json:"support,omitempty"`
	Payment          []string      `json:"payment,omitempty"`
}

// ActivateAbilityRequest activates an ability on an in-play card.
type ActivateAbilityRequest struct {
	RequestID        string        `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64        `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string        `json:"cardid" binding:"required"`
	AbilityID        string        `json:"abilityid" binding:"required"`
	Support          []CardSupport `json:"support,omitempty"`
	Payment          []string      `json:"payment,omitempty"`
}

// ResearchElementRequest uses Research to gain the selected Element.
type ResearchElementRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string `json:"cardid" binding:"required"`
	Element          string `json:"element" binding:"required" enums:"fire,earth,air,water"`
}

// ResearchSpellRequest uses Research with supporting Element-card instance IDs.
type ResearchSpellRequest struct {
	RequestID        string   `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64   `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string   `json:"cardid" binding:"required"`
	Elements         []string `json:"elements,omitempty"`
}

// HealWithFocusRequest uses Focus for healing; another Focus can contribute through support.
type HealWithFocusRequest struct {
	RequestID        string        `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64        `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string        `json:"cardid" binding:"required"`
	Support          []CardSupport `json:"support,omitempty"`
	Payment          []string      `json:"payment,omitempty"`
}

// EndTurnRequest enters the End Phase.
type EndTurnRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
}

// MaintenanceRequest selects Assets to keep and pays their upkeep. Keep can be empty.
type MaintenanceRequest struct {
	RequestID        string   `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64   `json:"expectedrevision" binding:"required" minimum:"0"`
	Keep             []string `json:"keep" binding:"required"`
	Payment          []string `json:"payment,omitempty"`
}

// KneelRequest requires an explicit boolean, including false.
type KneelRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
	Kneel            *bool  `json:"kneel" binding:"required"`
}

// ResolveEffectRequest resolves an offered effect other than a Full Moon hunt.
type ResolveEffectRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
	EffectID         string `json:"effectid" binding:"required"`
}

// SkipEffectRequest declines an optional effect using its unprefixed ID.
type SkipEffectRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
	EffectID         string `json:"effectid" binding:"required"`
}

// UndertakeHuntRequest commits to Full Moon’s hunt.
type UndertakeHuntRequest struct {
	RequestID        string        `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64        `json:"expectedrevision" binding:"required" minimum:"0"`
	EffectID         string        `json:"effectid" binding:"required"`
	Minions          []string      `json:"minions,omitempty"`
	Support          []CardSupport `json:"support,omitempty"`
	Payment          []string      `json:"payment,omitempty"`
}

// ChooseResearchRequest selects an offered card design ID.
type ChooseResearchRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
	DesignID         string `json:"designid" binding:"required"`
}

// EndDiscardRequest selects hand cards to discard. Cards can be empty.
type EndDiscardRequest struct {
	RequestID        string   `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64   `json:"expectedrevision" binding:"required" minimum:"0"`
	Cards            []string `json:"cards" binding:"required"`
}

// HandDiscardRequest selects one hand card and its explicit discard or destroy disposition.
type HandDiscardRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string `json:"cardid" binding:"required"`
	Disposition      string `json:"disposition" binding:"required" enums:"discard,destroy"`
}

// AssetLossRequest answers either an Asset-discard or Attack-loss prompt.
type AssetLossRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string `json:"cardid" binding:"required"`
}

// ScryRequest takes one inspected card and orders the remaining cards top first.
type ScryRequest struct {
	RequestID        string   `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64   `json:"expectedrevision" binding:"required" minimum:"0"`
	Take             string   `json:"take" binding:"required"`
	Order            []string `json:"order" binding:"required"`
}

// ReactRequest plays a reaction card while holding Priority.
type ReactRequest struct {
	RequestID        string        `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64        `json:"expectedrevision" binding:"required" minimum:"0"`
	CardID           string        `json:"cardid" binding:"required"`
	Targets          []string      `json:"targets,omitempty"`
	Support          []CardSupport `json:"support,omitempty"`
	Payment          []string      `json:"payment,omitempty"`
}

// DefendRequest commits ready Minions and/or Focus while holding Priority.
type DefendRequest struct {
	RequestID        string        `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64        `json:"expectedrevision" binding:"required" minimum:"0"`
	Minions          []string      `json:"minions,omitempty"`
	Support          []CardSupport `json:"support,omitempty"`
}

// PassRequest declines the current Priority opportunity.
type PassRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
}

type CardSupport struct {
	CardID  string   `json:"cardid"`
	Domains []string `json:"domains,omitempty"`
}

type Modifier struct {
	Source           string         `json:"source"`
	Affinity         map[string]int `json:"affinity"`
	ExpiresAfterTurn int            `json:"expiresafterturn"`
}

type RumourView struct {
	ID               string `json:"id"`
	DesignID         string `json:"designid"`
	ExpiresAfterTurn int    `json:"expiresafterturn"`
}

// RequiredChoice has public identity and timing. Private choices omit options
// and inspected cards from everyone except the deciding wizard.
type EffectOption struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Source   string `json:"source"`
	Optional bool   `json:"optional"`
}

type RequiredChoice struct {
	Effects []EffectOption `json:"effects,omitempty"`
	ID      string         `json:"id"`
	Actor   string         `json:"actor"`
	Type    string         `json:"type"`
	Source  string         `json:"source"`
	Options []string       `json:"options,omitempty"`
	Cards   []CardInstance `json:"cards,omitempty"`
	Minimum int            `json:"minimum"`
	Maximum int            `json:"maximum"`
	Private bool           `json:"private"`
}

type PendingEffectView struct {
	Spell         bool           `json:"spell"`
	Defenders     []string       `json:"defenders,omitempty"`
	BonusAttack   int            `json:"bonusattack"`
	BonusDefence  int            `json:"bonusdefence"`
	MinionBonuses map[string]int `json:"minionbonuses,omitempty"`
	DefenceFocus  bool           `json:"defencefocus"`
	Parent        string         `json:"parent,omitempty"`
	ID            string         `json:"id"`
	Actor         string         `json:"actor"`
	Kind          string         `json:"kind"`
	Source        string         `json:"source"`
	Target        string         `json:"target,omitempty"`
	Targets       []string       `json:"targets,omitempty"`
	Minions       []string       `json:"minions,omitempty"`
	Cards         []CardInstance `json:"cards"`
	Support       []CardSupport  `json:"support,omitempty"`
	Payment       []CardInstance `json:"payment,omitempty"`
	Mode          string         `json:"mode,omitempty"`
	Element       string         `json:"element,omitempty"`
	ChallengeID   string         `json:"challengeid,omitempty"`
	Countered     bool           `json:"countered"`
}

type ExchangeView struct {
	Resolving    *PendingEffectView  `json:"resolving"`
	Participants []string            `json:"participants"`
	Priority     string              `json:"priority"`
	Passes       int                 `json:"passes"`
	Stack        []PendingEffectView `json:"stack"`
}
