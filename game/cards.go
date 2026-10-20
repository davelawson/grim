package game

import "github.com/davelawson/grim/api"

var domains = []string{"fire", "earth", "air", "water", "nobility", "church", "infernal", "might", "popularity"}

func stats(domain string, n int) map[string]int { return map[string]int{domain: n} }
func completeCatalogue(c *Catalogue) {
	add := func(id, name, category string, cost int, timing string) *CardDesign {
		c.Cards = append(c.Cards, CardDesign{ID: id, Name: name, Categories: []string{category}, Cost: stats("wis", cost), PlayRequirements: map[string]int{}, Expertise: map[string]int{}, Affinity: map[string]int{}, Income: map[string]int{}, Resources: map[string]int{}, Persistence: "instant", Disposal: "discard", Timing: timing, Effects: []Effect{Effect(id)}})
		return &c.Cards[len(c.Cards)-1]
	}
	for i := range c.Cards {
		d := &c.Cards[i]
		d.Timing = "action"
		if d.Placeholder {
			d.Timing = "unavailable"
		}
		switch d.ID {
		case "research_spell":
			d.Cost = stats("wis", 1)
			d.Text = "Learn an Element, or discover a Spell from the pool selected by zero, one, or two distinct Element cards."
		case "arcane_focus":
			d.Timing = "action_or_support"
			d.Text = "Support: 1 Wis and one chosen Affinity; serve as Challenger. Alternatively pay 2 Wis to heal 1, or support an Attack or Defence with 1 Might."
		case "investigate":
			d.Challenges = []api.ChallengeDesign{{ID: "investigate"}}
			d.Text = "Challenge: reserve one eligible rumour using this attempt's final domain totals."
		case "subjugate_rival":
			d.Text = "Initiate your one Attack this turn with ready, eligible Might Minions against your Target."
		}
	}
	for _, resource := range []string{"wis", "wealth"} {
		d := add("ephemeral_"+resource, "Ephemeral "+resource, "resource", 0, "action")
		d.Resources = stats(resource, 1)
		d.Disposal = "destroy"
	}
	type spell struct {
		id, name           string
		cost               int
		pool, timing, text string
		requirement        map[string]int
	}
	spells := []spell{
		{"flame_lash", "Flame Lash", 1, "fire", "attack_reaction", "During your Attack, add 2 attacking Might.", stats("fire", 1)},
		{"firewall", "Firewall", 2, "fire", "defence_reaction", "Discard one committed attacking Minion.", stats("fire", 2)},
		{"flamecraft", "Flamecraft", 0, "fire", "support", "Challenge support: 2 Fire.", stats("fire", 1)},
		{"inner_flame", "Inner Flame", 2, "fire", "lasting", "Provide 1 Fire Affinity.", stats("fire", 1)},
		{"rockslide", "Rockslide", 1, "earth", "attack_reaction", "Each committed attacker contributes 1 additional Might.", stats("earth", 1)},
		{"stoneward", "Stoneward", 0, "earth", "defence_reaction", "Add 2 defensive Might.", stats("earth", 1)},
		{"shape_stone", "Shape Stone", 0, "earth", "support", "Challenge support: 2 Earth.", stats("earth", 1)},
		{"rampart", "Rampart", 2, "earth", "lasting", "Add 2 defensive Might. Upkeep: 1 Wis.", stats("earth", 2)},
		{"dislocate", "Dislocate", 1, "air", "attack_reaction", "Return one committed defender to its owner's hand.", stats("air", 2)},
		{"deflect", "Deflect", 1, "air", "defence_reaction", "Return one committed attacker to its owner's hand.", stats("air", 1)},
		{"keen_insight", "Keen Insight", 0, "air", "support", "Challenge support: 1 in one chosen Affinity.", stats("air", 1)},
		{"whispering_wind", "Whispering Wind", 2, "air", "lasting", "At Start.Effects, optionally draw one card, then discard one.", stats("air", 2)},
		{"deep_freeze", "Deep Freeze", 2, "water", "action", "Exhaust one ready neighbouring Minion; it skips its next Recovery.", stats("water", 2)},
		{"moat", "Moat", 0, "water", "defence_reaction", "Up to two committed attackers skip their next Recovery.", stats("water", 1)},
		{"flowing_form", "Flowing Form", 0, "water", "support", "Challenge support: 1 in each of two distinct chosen Affinities.", stats("water", 1)},
		{"spell_healing_spring", "Healing Spring", 2, "water", "lasting", "Action: pay 3 Wis and exhaust this to restore 1 Integrity, maximum 3.", stats("water", 2)},
		{"counterspell", "Counterspell", 1, "unaligned", "spell_reaction", "Counter one pending Spell effect.", nil},
		{"study", "Study", 1, "unaligned", "action", "Draw two, then discard one.", nil},
		{"unbound_will", "Unbound Will", 0, "unaligned", "support", "Permit a Challenge without a Minion or Arcane Focus; supply no capacity.", nil},
		{"mnemonic_seal", "Mnemonic Seal", 2, "unaligned", "lasting", "Increase maximum hand size by one.", nil},
		{"molten_armour", "Molten Armour", 2, "earth_fire", "defence_reaction", "Add 2 defensive Might; discard one committed attacker with Might Expertise at most 1.", nil},
		{"forgeheart", "Forgeheart", 2, "earth_fire", "lasting", "Provide Fire 1 and Earth 1.", nil},
		{"wildfire", "Wildfire", 2, "air_fire", "attack_reaction", "Discard one committed defender.", nil},
		{"flash_of_insight", "Flash of Insight", 1, "air_fire", "action", "Draw two, then discard or destroy one card from hand.", nil},
		{"scalding_mist", "Scalding Mist", 2, "fire_water", "combat_reaction", "Discard one opposing committed Minion with Might Expertise at most 1; up to one other opponent skips its next Recovery.", nil},
		{"purifying_steam", "Purifying Steam", 1, "fire_water", "action", "Additional cost: discard one hand card. Restore 1 Integrity, maximum 3.", nil},
		{"sandstorm", "Sandstorm", 1, "air_earth", "combat_reaction", "Opposing committed Minions each contribute 1 less Might, minimum zero.", nil},
		{"wind_road", "Wind Road", 2, "air_earth", "lasting", "Reduce Minion Wis play costs by 1, minimum zero.", nil},
		{"quagmire", "Quagmire", 1, "earth_water", "defence_reaction", "Up to two attackers each contribute 1 less Might and skip their next Recovery.", nil},
		{"verdant_spring", "Verdant Spring", 2, "earth_water", "lasting", "Generate one ephemeral Wealth at Start.Income.", nil},
		{"freezing_fog", "Freezing Fog", 1, "air_water", "defence_reaction", "Return one attacker with Might Expertise at most 1 to hand; up to one other attacker skips its next Recovery.", nil},
		{"scrying_pool", "Scrying Pool", 1, "air_water", "action", "Privately inspect up to the top three draw cards, take one, and reorder the rest; do not reshuffle.", nil},
	}
	for _, x := range spells {
		d := add(x.id, x.name, "spell", x.cost, x.timing)
		d.Pool = x.pool
		d.Text = x.text
		if x.requirement != nil {
			d.PlayRequirements = x.requirement
		} else if x.pool != "unaligned" {
			for _, el := range []string{"fire", "earth", "air", "water"} {
				if containsElement(x.pool, el) {
					d.PlayRequirements[el] = 1
				}
			}
		}
		if x.timing == "lasting" {
			d.Persistence = "durable"
			d.Timing = "action"
		}
		switch x.id {
		case "inner_flame":
			d.Affinity = stats("fire", 1)
		case "forgeheart":
			d.Affinity = map[string]int{"fire": 1, "earth": 1}
		case "rampart":
			d.Upkeep = stats("wis", 1)
		case "verdant_spring":
			d.Income = stats("ephemeral_wealth", 1)
		case "flamecraft":
			d.Support = stats("fire", 2)
		case "shape_stone":
			d.Support = stats("earth", 2)
		}
	}
	rewards := []struct {
		id, name, category          string
		cost                        int
		affinity, expertise, income map[string]int
	}{
		{"pupil", "Pupil", "minion", 1, nil, nil, nil},
		{"witch", "Witch", "minion", 1, nil, stats("infernal", 1), nil},
		{"leyline", "Leyline", "location", 1, nil, nil, stats("ephemeral_wis", 1)},
		{"court_office", "Court Office", "location", 1, stats("nobility", 1), nil, nil},
		{"village_hall", "Village Hall", "location", 1, stats("popularity", 1), nil, nil},
		{"stone_circle", "Stone Circle", "location", 1, stats("earth", 1), nil, nil},
		{"healing_spring", "Healing Spring", "location", 1, stats("water", 1), nil, nil},
		{"abbey_liaison", "Abbey Liaison", "minion", 2, stats("church", 1), stats("church", 1), nil},
		{"crossroads_shrine", "Crossroads Shrine", "location", 1, stats("infernal", 1), nil, nil},
		{"caravan_guard", "Caravan Guard", "minion", 2, stats("might", 1), stats("might", 1), nil},
		{"ember_shrine", "Ember Shrine", "location", 2, map[string]int{"fire": 1, "infernal": 1}, nil, nil},
		{"chartered_quarry", "Chartered Quarry", "location", 3, map[string]int{"earth": 1, "nobility": 1}, nil, stats("ephemeral_wealth", 1)},
		{"reedbound_healer", "Reedbound Healer", "minion", 2, stats("church", 1), map[string]int{"water": 1, "church": 1}, nil},
	}
	for _, x := range rewards {
		d := add(x.id, x.name, x.category, x.cost, "action")
		d.Persistence = "durable"
		if x.affinity != nil {
			d.Affinity = x.affinity
		}
		if x.expertise != nil {
			d.Expertise = x.expertise
		}
		if x.income != nil {
			d.Income = x.income
		}
	}
	challenge := func(id string, req, cost map[string]int, rewards ...string) api.ChallengeDesign {
		return api.ChallengeDesign{ID: id, Requirements: req, Cost: cost, Rewards: rewards}
	}
	rumour := func(id, name string, duration int, eligibility map[string]int, challenges ...api.ChallengeDesign) {
		d := add("rumour_"+id, name, "rumour", 0, "start")
		d.Duration = duration
		d.Eligibility = eligibility
		d.Challenges = challenges
		d.Persistence = "duration"
		d.Disposal = "mill"
	}
	pupil := challenge("accept_pupil", nil, stats("wis", 3), "pupil")
	pupil.RequiresFocus = true
	rumour("promising_pupil", "Promising Pupil", 2, nil, pupil)
	rumour("solar_eclipse", "Solar Eclipse", 1, nil)
	rumour("burn_the_witch", "Burn the Witch", 3, nil, challenge("arrest_witch", stats("might", 1), nil), challenge("recruit_witch", nil, stats("wis", 2), "witch"))
	rumour("new_leyline", "New Leyline", 2, nil, challenge("claim_leyline", nil, stats("wis", 2), "leyline"))
	moon := challenge("moonlit_hunt", stats("might", 1), nil, "basic_wealth")
	moon.AffinityReward = map[string]int{"water": 1, "infernal": 1}
	rumour("full_moon", "Full Moon", 1, nil, moon)
	rumour("wandering_merchant", "Wandering Merchant", 2, nil, challenge("strike_bargain", nil, stats("wealth", 1), "basic_wis", "basic_wealth"))
	courier := challenge("escort_courier", stats("might", 1), nil, "basic_wealth")
	courier.AffinityReward = stats("nobility", 1)
	rumour("lost_courier", "Lost Courier", 2, nil, challenge("follow_messages", stats("air", 1), nil, "court_office"), courier)
	rumour("village_petition", "Village Petition", 2, nil, challenge("rally_villagers", stats("popularity", 1), stats("wis", 1), "village_hall"), challenge("repair_well", stats("earth", 1), stats("wis", 1), "village_hall"))
	rumour("restless_kiln", "Restless Kiln", 2, stats("fire", 1), challenge("tame_furnace", stats("fire", 2), nil, "basic_wealth", "basic_wealth"))
	rumour("boundary_stones", "Buried Boundary Stones", 2, stats("earth", 1), challenge("raise_stones", stats("earth", 2), nil, "stone_circle"))
	whispers := challenge("gather_whispers", stats("air", 2), nil, "basic_wis")
	whispers.AffinityReward = stats("popularity", 1)
	rumour("whispers", "Whispers on the Wind", 2, stats("air", 1), whispers)
	rumour("healing_spring", "The Healing Spring", 2, stats("water", 1), challenge("restore_spring", stats("water", 2), nil, "healing_spring"))
	rumour("inheritance", "Disputed Inheritance", 2, stats("nobility", 1), challenge("witness_settlement", stats("nobility", 2), stats("wis", 1), "court_office"))
	rumour("infirmary", "Empty Abbey Infirmary", 2, stats("church", 1), challenge("organise_relief", stats("church", 1), stats("wis", 1), "abbey_liaison"))
	rumour("crossroads", "Voices at the Crossroads", 2, stats("infernal", 1), challenge("make_offering", stats("infernal", 1), stats("wealth", 1), "crossroads_shrine"), challenge("learn_names", stats("infernal", 2), nil, "crossroads_shrine"))
	rumour("caravan", "Caravan under Escort", 2, stats("might", 1), challenge("guard_road", stats("might", 2), nil, "basic_wealth", "basic_wealth"), challenge("hire_guard", stats("might", 1), stats("wealth", 1), "caravan_guard"))
	rumour("village_fair", "Village Fair", 2, stats("popularity", 1), challenge("win_crowd", stats("popularity", 2), nil, "village_hall", "basic_wealth"))
	rumour("gallows", "Embers beneath the Gallows", 3, map[string]int{"fire": 1, "infernal": 1}, challenge("bind_embers", map[string]int{"fire": 2, "infernal": 2}, stats("wis", 1), "ember_shrine"))
	rumour("quarry", "The Baron's Forgotten Quarry", 3, map[string]int{"earth": 1, "nobility": 1}, challenge("reopen_quarry", map[string]int{"earth": 2, "nobility": 1}, stats("wealth", 1), "chartered_quarry"))
	rumour("crypt", "The Abbey's Flooded Crypt", 3, map[string]int{"water": 1, "church": 1}, challenge("recover_records", map[string]int{"water": 2, "church": 1}, stats("wis", 1), "reedbound_healer"))
	for i := range c.Cards {
		d := &c.Cards[i]
		for j := range d.Challenges {
			if len(d.Challenges[j].AffinityReward) > 0 {
				d.Challenges[j].RewardDuration = 2
			}
		}
		switch d.ID {
		case "spell_healing_spring":
			d.Abilities = []api.AbilityDesign{{ID: "heal", Timing: "action", Cost: stats("wis", 3), Exhaust: true, Effect: "heal_activation", Text: "Restore 1 Integrity, maximum 3."}}
		case "whispering_wind":
			d.Abilities = []api.AbilityDesign{{ID: "whispering_wind_trigger", Timing: "start_effects", Optional: true, Effect: "whispering_wind_trigger", Text: "Draw one card, then discard one."}}
		case "rumour_solar_eclipse":
			d.Abilities = []api.AbilityDesign{{ID: "solar_eclipse", Timing: "start_effects", Effect: "solar_eclipse", Text: "Gain 2 Wis, Fire 1 and Infernal 1 until this turn ends."}}
		case "rumour_burn_the_witch":
			d.Abilities = []api.AbilityDesign{{ID: "burn_failure", Timing: "expiry_failure", Effect: "burn_failure", Text: "If no Challenge succeeded before expiry, discard one chosen Asset; defer to next Start.Effects when this turn's End Phase is skipped."}}
		case "pupil":
			d.Text = "On recruitment choose Expertise 1 uniformly from the eight domains other than the drafted Chantry's Element. This trait is permanent."
		}
	}

}

func containsElement(pool, element string) bool {
	// Pool IDs use canonical alphabetic pairs.
	for i := 0; i+len(element) <= len(pool); i++ {
		if pool[i:i+len(element)] == element {
			return true
		}
	}
	return false
}
