package api

// Effect identifies a typed engine handler.
type Effect string

const (
	EffectChantryIncome    Effect = "chantry_income"
	EffectElement          Effect = "element_support"
	EffectResource         Effect = "resource"
	EffectResearchSpell    Effect = "research_spell"
	EffectResearchArtifact Effect = "research_artifact"
	EffectFocus            Effect = "arcane_focus"
	EffectInvestigate      Effect = "investigate"
	EffectAttack           Effect = "subjugate_rival"
	EffectMinion           Effect = "starting_minion"
)

// CardDesign contains printed properties, separate from individual copies.
type CardDesign struct {
	Abilities        []AbilityDesign   `json:"abilities,omitempty"`
	Upkeep           map[string]int    `json:"upkeep,omitempty"`
	Timing           string            `json:"timing,omitempty"`
	Duration         int               `json:"duration,omitempty"`
	Eligibility      map[string]int    `json:"eligibility,omitempty"`
	Challenges       []ChallengeDesign `json:"challenges,omitempty"`
	Support          map[string]int    `json:"support,omitempty"`
	Pool             string            `json:"pool,omitempty"`
	Text             string            `json:"text,omitempty"`
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Categories       []string          `json:"categories"`
	Cost             map[string]int    `json:"cost"`
	PlayRequirements map[string]int    `json:"playrequirements"`
	Expertise        map[string]int    `json:"expertise"`
	Affinity         map[string]int    `json:"affinity"`
	Income           map[string]int    `json:"income"`
	Resources        map[string]int    `json:"resources"`
	Element          string            `json:"element,omitempty"`
	VictoryPath      string            `json:"victorypath,omitempty"`
	Persistence      string            `json:"persistence"`
	Disposal         string            `json:"disposal"`
	Placeholder      bool              `json:"placeholder"`
	Effects          []Effect          `json:"effects"`
}

type Catalogue struct {
	RulesVersion     string       `json:"rulesversion"`
	CatalogueVersion string       `json:"catalogueversion"`
	Cards            []CardDesign `json:"cards"`
}

type CatalogueResponse struct {
	Catalogue *Catalogue `json:"catalogue"`
}

func (c *Catalogue) Card(id string) (CardDesign, bool) {
	for _, card := range c.Cards {
		if card.ID == id {
			return card, true
		}
	}
	return CardDesign{}, false
}

// ChallengeDesign describes requirements beyond the normal Challenger rule.
type ChallengeDesign struct {
	RewardDuration int            `json:"rewardduration,omitempty"`
	ID             string         `json:"id"`
	Requirements   map[string]int `json:"requirements,omitempty"`
	Cost           map[string]int `json:"cost,omitempty"`
	RequiresFocus  bool           `json:"requiresfocus,omitempty"`
	Rewards        []string       `json:"rewards,omitempty"`
	AffinityReward map[string]int `json:"affinityreward,omitempty"`
}

// AbilityDesign describes a card-specific activation or scheduled effect.
type AbilityDesign struct {
	ID       string         `json:"id"`
	Timing   string         `json:"timing"`
	Cost     map[string]int `json:"cost,omitempty"`
	Exhaust  bool           `json:"exhaust"`
	Optional bool           `json:"optional"`
	Effect   Effect         `json:"effect"`
	Text     string         `json:"text"`
}
