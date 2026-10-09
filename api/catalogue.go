package api

// Effect identifies documented behavior for future typed engine handlers.
// Drafting does not execute any of these effects.
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
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Categories       []string       `json:"categories"`
	Cost             map[string]int `json:"cost"`
	PlayRequirements map[string]int `json:"playrequirements"`
	Expertise        map[string]int `json:"expertise"`
	Affinity         map[string]int `json:"affinity"`
	Income           map[string]int `json:"income"`
	Resources        map[string]int `json:"resources"`
	Element          string         `json:"element,omitempty"`
	VictoryPath      string         `json:"victorypath,omitempty"`
	Persistence      string         `json:"persistence"`
	Disposal         string         `json:"disposal"`
	Placeholder      bool           `json:"placeholder"`
	Effects          []Effect       `json:"effects"`
}

type Catalogue struct {
	RulesVersion     string       `json:"rulesversion"`
	CatalogueVersion string       `json:"catalogueversion"`
	Cards            []CardDesign `json:"cards"`
	DraftOffers      Offers       `json:"draftoffers"`
}

type Offers struct {
	Chantries    []string            `json:"chantries"`
	VictoryPaths []string            `json:"victorypaths"`
	Minions      map[string][]string `json:"minions"`
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
