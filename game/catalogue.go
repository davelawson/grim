package game

import (
	"github.com/davelawson/grim/api"
	"maps"
)

const (
	OpeningRulesVersion     = "playable-v1"
	DefaultCatalogueVersion = "playable-v1"
)

// Public card definitions are shared with API consumers.
type Effect = api.Effect
type CardDesign = api.CardDesign
type Catalogue = api.Catalogue

// Offers is the engine's opening menu, separate from public card definitions.
type Offers struct {
	Chantries    []string
	VictoryPaths []string
	Minions      map[string][]string
}

const (
	EffectChantryIncome    = api.EffectChantryIncome
	EffectElement          = api.EffectElement
	EffectResource         = api.EffectResource
	EffectResearchSpell    = api.EffectResearchSpell
	EffectResearchArtifact = api.EffectResearchArtifact
	EffectFocus            = api.EffectFocus
	EffectInvestigate      = api.EffectInvestigate
	EffectAttack           = api.EffectAttack
	EffectMinion           = api.EffectMinion
)

func OpeningOffers() Offers {
	return Offers{
		Chantries:    []string{"chantry_fire", "chantry_earth", "chantry_air", "chantry_water"},
		VictoryPaths: []string{"subjugate_rival"},
		Minions: map[string][]string{
			"chantry_fire":  {"vera_coalheart", "bastian_redhand", "orla_nine_cinders"},
			"chantry_earth": {"hedda_stonewake", "leontine_grange", "rurik_ironvale"},
			"chantry_air":   {"mina_cloudscript", "garrick_stormward", "elsie_many_tongues"},
			"chantry_water": {"nadia_rivermark", "sister_agnes_reed", "tamsin_lockward"},
		},
	}
}

func cloneStats(stats map[string]int) map[string]int { return maps.Clone(stats) }

var domains = []string{"fire", "earth", "air", "water", "nobility", "church", "infernal", "might", "popularity"}

func stats(domain string, n int) map[string]int { return map[string]int{domain: n} }
