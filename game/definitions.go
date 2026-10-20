package game

import (
	"fmt"
	"slices"
	"strings"
)

// Definitions is an immutable, validated catalogue. It is injected into game
// state by the application and never serialized into a match snapshot.
type Definitions struct {
	catalogue Catalogue
	byID      map[string]CardDesign
}

func NewDefinitions(c Catalogue) (*Definitions, error) {
	if c.RulesVersion != OpeningRulesVersion || strings.TrimSpace(c.CatalogueVersion) == "" {
		return nil, ErrUnsupportedVersion
	}
	d := &Definitions{catalogue: clone(c), byID: map[string]CardDesign{}}
	if len(c.Cards) == 0 {
		return nil, fmt.Errorf("empty catalogue")
	}
	for _, card := range d.catalogue.Cards {
		if strings.TrimSpace(card.ID) == "" || strings.TrimSpace(card.Name) == "" {
			return nil, fmt.Errorf("card ID and name are required")
		}
		if _, ok := d.byID[card.ID]; ok {
			return nil, fmt.Errorf("duplicate card %s", card.ID)
		}
		d.byID[card.ID] = card
	}
	for _, card := range d.catalogue.Cards {
		if err := d.validateCard(card); err != nil {
			return nil, fmt.Errorf("card %s: %w", card.ID, err)
		}
	}
	offers := OpeningOffers()
	required := append(slices.Clone(offers.Chantries), offers.VictoryPaths...)
	for _, ids := range offers.Minions {
		required = append(required, ids...)
	}
	required = append(required, "research_spell", "arcane_focus", "investigate", "subjugate_rival", "basic_wis", "basic_wealth", "ephemeral_wis", "ephemeral_wealth", "element_fire", "element_earth", "element_air", "element_water")
	for _, id := range required {
		card, ok := d.byID[id]
		if !ok || card.Placeholder {
			return nil, fmt.Errorf("missing required playable design %s", id)
		}
	}
	return d, nil
}

// Catalogue returns a copy; callers cannot alter another match's definitions.
func (d *Definitions) Catalogue() *Catalogue { c := clone(d.catalogue); return &c }

func (s *State) compatibleDefinitions() bool {
	return s.definitions != nil && s.RulesVersion == OpeningRulesVersion && s.RulesVersion == s.definitions.catalogue.RulesVersion && s.CatalogueVersion == s.definitions.catalogue.CatalogueVersion
}
func (s *State) clone() *State { next := clone(s); next.definitions = s.definitions; return next }

func validateStats(values map[string]int, keys []string) error {
	for key, value := range values {
		if !slices.Contains(keys, key) || value < 0 {
			return fmt.Errorf("invalid stat %s=%d", key, value)
		}
	}
	return nil
}

// Handler names belong to the Go rules engine; printed card values belong to
// the database. Generic Assets and rumours use their own design identity.
func supportedEffect(card CardDesign, effect string) bool {
	if effect == card.ID && (category(card, "minion") || category(card, "location") || category(card, "rumour") || card.Persistence == "durable") {
		return true
	}
	switch effect {
	case "chantry_income", "element_support", "resource", "starting_minion", "arcane_focus", "investigate", "subjugate_rival", "research_spell", "ephemeral_wis", "ephemeral_wealth",
		"flame_lash", "firewall", "flamecraft", "rockslide", "stoneward", "shape_stone", "dislocate", "deflect", "keen_insight", "deep_freeze", "moat", "flowing_form", "counterspell", "study", "unbound_will", "molten_armour", "wildfire", "flash_of_insight", "scalding_mist", "purifying_steam", "sandstorm", "quagmire", "freezing_fog", "scrying_pool",
		"heal_activation", "whispering_wind_trigger", "solar_eclipse", "burn_failure":
		return true
	}
	return false
}

func (d *Definitions) validateCard(card CardDesign) error {
	if len(card.Categories) == 0 || !unique(card.Categories) {
		return fmt.Errorf("invalid categories")
	}
	for _, category := range card.Categories {
		if !slices.Contains([]string{"activity", "placeholder", "location", "element", "minion", "spell", "rumour", "resource"}, category) {
			return fmt.Errorf("unsupported category %s", category)
		}
	}
	if !slices.Contains([]string{"action", "action_or_support", "start", "support", "unavailable", "attack_reaction", "defence_reaction", "combat_reaction", "spell_reaction"}, card.Timing) {
		return fmt.Errorf("unsupported timing %s", card.Timing)
	}
	if !slices.Contains([]string{"instant", "durable", "duration"}, card.Persistence) || !slices.Contains([]string{"discard", "destroy", "mill"}, card.Disposal) || card.Duration < 0 {
		return fmt.Errorf("invalid lifetime")
	}
	if card.Element != "" && !slices.Contains([]string{"fire", "earth", "air", "water"}, card.Element) {
		return fmt.Errorf("invalid element %s", card.Element)
	}
	if card.Pool != "" && !slices.Contains([]string{"unaligned", "fire", "earth", "air", "water", "earth_fire", "air_fire", "fire_water", "air_earth", "earth_water", "air_water"}, card.Pool) {
		return fmt.Errorf("unsupported Spell pool %s", card.Pool)
	}
	if category(card, "spell") && card.Pool == "" {
		return fmt.Errorf("missing Spell pool")
	}
	for _, values := range []map[string]int{card.PlayRequirements, card.Expertise, card.Affinity, card.Eligibility, card.Support} {
		if err := validateStats(values, domains); err != nil {
			return err
		}
	}
	for _, values := range []map[string]int{card.Cost, card.Resources, card.Upkeep} {
		if err := validateStats(values, []string{"wis", "wealth"}); err != nil {
			return err
		}
	}
	if err := validateStats(card.Income, []string{"ephemeral_wis", "ephemeral_wealth"}); err != nil {
		return err
	}
	for _, effect := range card.Effects {
		if !supportedEffect(card, string(effect)) {
			return fmt.Errorf("unsupported effect %s", effect)
		}
	}
	if !card.Placeholder && len(card.Effects) == 0 {
		return fmt.Errorf("missing effect")
	}
	ids := []string{}
	for _, ability := range card.Abilities {
		ids = append(ids, ability.ID)
		if !slices.Contains([]string{"action", "start_effects", "expiry_failure"}, ability.Timing) || !supportedEffect(card, string(ability.Effect)) {
			return fmt.Errorf("unsupported ability %s", ability.ID)
		}
		if err := validateStats(ability.Cost, []string{"wis", "wealth"}); err != nil {
			return err
		}
	}
	if !unique(ids) {
		return fmt.Errorf("invalid ability IDs")
	}
	ids = nil
	for _, challenge := range card.Challenges {
		ids = append(ids, challenge.ID)
		for _, values := range []map[string]int{challenge.Requirements, challenge.AffinityReward} {
			if err := validateStats(values, domains); err != nil {
				return err
			}
		}
		if err := validateStats(challenge.Cost, []string{"wis", "wealth"}); err != nil {
			return err
		}
		if challenge.RewardDuration < 0 {
			return fmt.Errorf("negative reward duration")
		}
		for _, id := range challenge.Rewards {
			reward, ok := d.byID[id]
			if !ok || reward.Placeholder || category(reward, "rumour") {
				return fmt.Errorf("invalid reward reference %s", id)
			}
		}
	}
	if !unique(ids) {
		return fmt.Errorf("invalid Challenge IDs")
	}
	return nil
}
