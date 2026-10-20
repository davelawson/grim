// Package catalogue loads shared card definitions from SQLite. It never seeds or edits data.
package catalogue

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davelawson/grim/api"
	"strings"
)

var ErrNotFound = errors.New("catalogue version not found")

type Querier interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type statGroup struct {
	prefix string
	keys   []string
	values *map[string]int
}

func stats(d *api.CardDesign) []statGroup {
	domains := []string{"fire", "earth", "air", "water", "nobility", "church", "infernal", "might", "popularity"}
	return []statGroup{
		{"cost", []string{"wis", "wealth"}, &d.Cost},
		{"play_requirement", domains, &d.PlayRequirements},
		{"expertise", domains, &d.Expertise},
		{"affinity", domains, &d.Affinity},
		{"income", []string{"ephemeral_wis", "ephemeral_wealth"}, &d.Income},
		{"resource", []string{"wis", "wealth"}, &d.Resources},
		{"upkeep", []string{"wis", "wealth"}, &d.Upkeep},
		{"eligibility", domains, &d.Eligibility},
		{"support", domains, &d.Support},
	}
}

// Load reads a complete version in stable catalogue order using the caller's transaction.
func Load(q Querier, version string) (*api.Catalogue, error) {
	var shape api.CardDesign
	columns := []string{"design_id", "rules_version", "sort_order", "pool", "name", "text", "timing", "duration", "element", "victory_path", "persistence", "disposal", "placeholder"}
	for _, g := range stats(&shape) {
		for _, key := range g.keys {
			columns = append(columns, g.prefix+"_"+key)
		}
	}
	columns = append(columns, "categories", "effects", "abilities", "challenges")
	rows, err := q.Query("select "+strings.Join(columns, ",")+" from card_designs where catalogue_version = ? order by sort_order", version)
	if err != nil {
		return nil, fmt.Errorf("load catalogue %q: %w", version, err)
	}
	defer rows.Close()
	result := &api.Catalogue{CatalogueVersion: version, Cards: []api.CardDesign{}}
	for rows.Next() {
		d := api.CardDesign{Cost: map[string]int{}, PlayRequirements: map[string]int{}, Expertise: map[string]int{}, Affinity: map[string]int{}, Income: map[string]int{}, Resources: map[string]int{}}
		var rules string
		var ordinal int
		target := []any{&d.ID, &rules, &ordinal, &d.Pool, &d.Name, &d.Text, &d.Timing, &d.Duration, &d.Element, &d.VictoryPath, &d.Persistence, &d.Disposal, &d.Placeholder}
		groups := stats(&d)
		count := 0
		for _, g := range groups {
			count += len(g.keys)
		}
		values := make([]sql.NullInt64, count)
		for i := range values {
			target = append(target, &values[i])
		}
		var lists [4]string
		for i := range lists {
			target = append(target, &lists[i])
		}
		if err := rows.Scan(target...); err != nil {
			return nil, fmt.Errorf("decode catalogue %q row: %w", version, err)
		}
		if ordinal != len(result.Cards) {
			return nil, fmt.Errorf("catalogue %q has missing or invalid order at %s", version, d.ID)
		}
		if result.RulesVersion != "" && result.RulesVersion != rules {
			return nil, fmt.Errorf("catalogue %q mixes rules versions", version)
		}
		result.RulesVersion = rules
		offset := 0
		for _, g := range groups {
			for _, key := range g.keys {
				v := values[offset]
				offset++
				if v.Valid {
					if v.Int64 < 0 || int64(int(v.Int64)) != v.Int64 {
						return nil, fmt.Errorf("invalid %s_%s on %s", g.prefix, key, d.ID)
					}
					if *g.values == nil {
						*g.values = map[string]int{}
					}
					(*g.values)[key] = int(v.Int64)
				}
			}
		}
		for i, value := range []any{&d.Categories, &d.Effects, &d.Abilities, &d.Challenges} {
			decoder := json.NewDecoder(strings.NewReader(lists[i]))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(value); err != nil {
				return nil, fmt.Errorf("invalid catalogue list on %s: %w", d.ID, err)
			}
		}
		if len(d.Abilities) == 0 {
			d.Abilities = nil
		}
		if len(d.Challenges) == 0 {
			d.Challenges = nil
		}
		result.Cards = append(result.Cards, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(result.Cards) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, version)
	}
	return result, nil
}
