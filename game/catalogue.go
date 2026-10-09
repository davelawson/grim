package game

import "maps"

const (
	OpeningRulesVersion     = "opening-v1"
	OpeningCatalogueVersion = "opening-v1"
)

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

func OpeningOffers() Offers {
	return Offers{
		Chantries:    []string{"chantry_fire", "chantry_earth", "chantry_air", "chantry_water"},
		VictoryPaths: []string{"subjugate_rival", "victory_arcane", "victory_influence", "victory_council", "victory_fame"},
		Minions: map[string][]string{
			"chantry_fire":  {"vera_coalheart", "bastian_redhand", "orla_nine_cinders"},
			"chantry_earth": {"hedda_stonewake", "leontine_grange", "rurik_ironvale"},
			"chantry_air":   {"mina_cloudscript", "garrick_stormward", "elsie_many_tongues"},
			"chantry_water": {"nadia_rivermark", "sister_agnes_reed", "tamsin_lockward"},
		},
	}
}

func cloneStats(stats map[string]int) map[string]int { return maps.Clone(stats) }

// OpeningCatalogue returns fresh definitions so callers cannot modify shared rules.
func OpeningCatalogue() *Catalogue {
	c := &Catalogue{RulesVersion: OpeningRulesVersion, CatalogueVersion: OpeningCatalogueVersion, DraftOffers: OpeningOffers()}
	add := func(id, name, category string, effects ...Effect) *CardDesign {
		c.Cards = append(c.Cards, CardDesign{ID: id, Name: name, Categories: []string{category},
			Cost: map[string]int{}, PlayRequirements: map[string]int{}, Expertise: map[string]int{},
			Affinity: map[string]int{}, Income: map[string]int{}, Resources: map[string]int{},
			Persistence: "instant", Disposal: "discard", Effects: append([]Effect{}, effects...)})
		return &c.Cards[len(c.Cards)-1]
	}
	for _, element := range []struct{ id, name string }{{"fire", "Fire"}, {"earth", "Earth"}, {"air", "Air"}, {"water", "Water"}} {
		chantry := add("chantry_"+element.id, element.name+" Chantry", "location", EffectChantryIncome)
		chantry.Element, chantry.Persistence = element.id, "durable"
		chantry.Cost["wis"], chantry.Affinity[element.id], chantry.Income["ephemeral_wis"] = 1, 1, 1
		card := add("element_"+element.id, element.name+" Element", "element", EffectElement)
		card.Element = element.id
	}
	add("research_spell", "Research Spell", "activity", EffectResearchSpell)
	add("research_artifact", "Research Artifact", "activity", EffectResearchArtifact)
	add("basic_wis", "Basic Wis", "resource", EffectResource).Resources["wis"] = 1
	add("basic_wealth", "Basic Wealth", "resource", EffectResource).Resources["wealth"] = 1
	add("arcane_focus", "Arcane Focus", "activity", EffectFocus)
	add("investigate", "Investigate", "activity", EffectInvestigate)
	add("subjugate_rival", "Subjugate Rival", "activity", EffectAttack).VictoryPath = "domination"
	for _, path := range []struct{ id, name string }{{"arcane", "Arcane"}, {"influence", "Influence"}, {"council", "Council"}, {"fame", "Fame"}} {
		card := add("victory_"+path.id, path.name+" (placeholder)", "placeholder")
		card.VictoryPath, card.Placeholder = path.id, true
	}
	for _, minion := range []struct{ id, name, element, expertise string }{
		{"vera_coalheart", "Vera Coalheart", "fire", "fire"},
		{"bastian_redhand", "Bastian Redhand", "fire", "might"},
		{"orla_nine_cinders", "Orla Nine-Cinders", "fire", "infernal"},
		{"hedda_stonewake", "Hedda Stonewake", "earth", "earth"},
		{"leontine_grange", "Leontine Grange", "earth", "nobility"},
		{"rurik_ironvale", "Rurik Ironvale", "earth", "might"},
		{"mina_cloudscript", "Mina Cloudscript", "air", "air"},
		{"garrick_stormward", "Garrick Stormward", "air", "might"},
		{"elsie_many_tongues", "Elsie Many-Tongues", "air", "popularity"},
		{"nadia_rivermark", "Nadia Rivermark", "water", "water"},
		{"sister_agnes_reed", "Sister Agnes Reed", "water", "church"},
		{"tamsin_lockward", "Tamsin Lockward", "water", "might"},
	} {
		card := add(minion.id, minion.name, "minion", EffectMinion)
		card.Persistence = "durable"
		card.Cost["wis"], card.PlayRequirements[minion.element], card.Expertise[minion.expertise] = 1, 1, 1
	}
	return c
}

func (c *Catalogue) Card(id string) (CardDesign, bool) {
	for _, card := range c.Cards {
		if card.ID == id {
			return card, true
		}
	}
	return CardDesign{}, false
}
