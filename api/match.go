package api

// MatchResponse wraps the requesting participant's public and private view.
type MatchResponse struct {
	Match *MatchView `json:"match"`
}

type DraftPickRequest struct {
	RequestID        string `json:"requestid" binding:"required" format:"uuid"`
	ExpectedRevision *int64 `json:"expectedrevision" binding:"required" minimum:"0"`
	Type             string `json:"type" binding:"required" enums:"choose_chantry,choose_victory_path,choose_minion"`
	CardID           string `json:"cardid" binding:"required"`
}

// DraftResponse contains only the caller's progress and current draft options.
type DraftResponse struct {
	Draft    Draft        `json:"draft"`
	Cards    []CardDesign `json:"cards"`
	Revision int64        `json:"revision"`
	Status   string       `json:"status"`
}

// MatchView excludes authoritative private state and random generator data.
type MatchView struct {
	Bots             map[string]Bot        `json:"bots,omitempty"`
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	LobbyID          string                `json:"lobbyid"`
	Status           string                `json:"status"`
	Revision         int64                 `json:"revision"`
	Participants     []string              `json:"participants"`
	Ring             []string              `json:"ring"`
	Winners          []string              `json:"winners"`
	Exchange         *ExchangeView         `json:"exchange"`
	Choice           *RequiredChoice       `json:"choice"`
	Outcome          *string               `json:"outcome"`
	RulesVersion     string                `json:"rulesversion"`
	CatalogueVersion string                `json:"catalogueversion"`
	Wizards          map[string]WizardView `json:"wizards"`
	Turn             *Turn                 `json:"turn"`
	You              *PrivateView          `json:"you,omitempty"`
}

type Draft struct {
	Step        string `json:"step"`
	Chantry     string `json:"chantry,omitempty"`
	VictoryPath string `json:"victorypath,omitempty"`
	Minion      string `json:"minion,omitempty"`
}

type CardInstance struct {
	ID           string         `json:"id"`
	DesignID     string         `json:"designid"`
	Owner        string         `json:"owner"`
	Location     string         `json:"location"`
	Exhausted    bool           `json:"exhausted"`
	DeployedTurn int            `json:"deployedturn,omitempty"`
	SkipRecovery bool           `json:"skiprecovery,omitempty"`
	Expertise    map[string]int `json:"expertise,omitempty"`
	Ephemeral    bool           `json:"ephemeral,omitempty"`
}

type Turn struct {
	Number       int    `json:"number"`
	ActiveWizard string `json:"activewizard"`
	Phase        string `json:"phase"`
	Step         string `json:"step"`
}

type WizardView struct {
	Turns           int            `json:"turns"`
	ImmuneUntil     int            `json:"immuneuntil"`
	KneelEligible   bool           `json:"kneeleligible"`
	Reserved        []string       `json:"reserved"`
	Resources       map[string]int `json:"resources"`
	Eliminated      bool           `json:"eliminated"`
	Immune          bool           `json:"immune"`
	Rumours         []RumourView   `json:"rumours"`
	Modifiers       []Modifier     `json:"modifiers"`
	MaximumHandSize int            `json:"maximumhandsize"`
	AttackUsed      bool           `json:"attackused"`
	DraftComplete   bool           `json:"draftcomplete"`
	Integrity       *int           `json:"integrity,omitempty"`
	Affinities      map[string]int `json:"affinities,omitempty"`
	VictoryProgress map[string]int `json:"victoryprogress,omitempty"`
	InPlay          []CardInstance `json:"inplay,omitempty"`
	Discard         []CardInstance `json:"discard"`
	HandCount       int            `json:"handcount"`
	DrawCount       int            `json:"drawcount"`
}

type PrivateView struct {
	UserID   string         `json:"userid,omitempty"`
	Draft    Draft          `json:"draft"`
	Hand     []CardInstance `json:"hand"`
	DrawPile []CardInstance `json:"drawpile"`
}
