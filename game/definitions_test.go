package game

import (
	"bytes"
	"strings"
	"testing"
)

func TestDefinitionsAreImmutableAndSurviveCloning(t *testing.T) {
	input := testCatalogue(t)
	defs, err := NewDefinitions(*input)
	if err != nil {
		t.Fatal(err)
	}
	input.Cards[0].Cost["wis"] = 99
	exported := defs.Catalogue()
	exported.Cards[0].Cost["wis"] = 88
	s, err := NewSetup("m", "name", "l", []string{"a", "b"}, defs)
	if err != nil {
		t.Fatal(err)
	}
	next, err := s.Choose("a", "choose_chantry", "chantry_fire")
	if err != nil {
		t.Fatal(err)
	}
	if next.definitions != defs || s.definition("chantry_fire").Cost["wis"] != 1 {
		t.Fatal("definitions lost or mutated")
	}
	saved := jsonBytes(t, next)
	if bytes.Contains(saved, []byte(`"definitions"`)) || bytes.Contains(saved, []byte(`"cards"`)) {
		t.Fatal("definitions entered match snapshot")
	}
	restored, err := Decode(saved, defs)
	if err != nil {
		t.Fatal(err)
	}
	if restored.definitions != defs {
		t.Fatal("definitions not rebound on restart")
	}
	if _, err := Decode(saved, nil); err == nil {
		t.Fatal("missing definitions accepted")
	}
	different := testCatalogue(t)
	different.CatalogueVersion = "another-version"
	other, err := NewDefinitions(*different)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Decode(saved, other); err == nil {
		t.Fatal("wrong version accepted")
	}
}

func TestInvalidDefinitionsFailClosed(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*Catalogue)
	}{
		{"unknown rules", func(c *Catalogue) { c.RulesVersion = "future" }},
		{"missing draft card", func(c *Catalogue) { c.Cards = c.Cards[1:] }},
		{"duplicate card", func(c *Catalogue) { c.Cards = append(c.Cards, c.Cards[0]) }},
		{"negative stat", func(c *Catalogue) { c.Cards[0].Cost["wis"] = -1 }},
		{"unknown stat", func(c *Catalogue) { c.Cards[0].Affinity["unknown"] = 1 }},
		{"unknown handler", func(c *Catalogue) { c.Cards[0].Effects = []Effect{"not-implemented"} }},
		{"unknown ability", func(c *Catalogue) {
			for i := range c.Cards {
				if len(c.Cards[i].Abilities) > 0 {
					c.Cards[i].Abilities[0].Effect = "not-implemented"
					break
				}
			}
		}},
		{"missing reward", func(c *Catalogue) {
			for i := range c.Cards {
				if len(c.Cards[i].Challenges) > 0 {
					c.Cards[i].Challenges[0].Rewards = []string{"missing"}
					break
				}
			}
		}},
		{"duplicate Challenge", func(c *Catalogue) {
			for i := range c.Cards {
				if len(c.Cards[i].Challenges) > 0 {
					c.Cards[i].Challenges = append(c.Cards[i].Challenges, c.Cards[i].Challenges[0])
					break
				}
			}
		}},
		{"unknown pool", func(c *Catalogue) {
			for i := range c.Cards {
				if strings.Contains(c.Cards[i].Pool, "fire") {
					c.Cards[i].Pool = "unknown"
					break
				}
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := testCatalogue(t)
			tc.edit(c)
			if _, err := NewDefinitions(*c); err == nil {
				t.Fatal("invalid definitions accepted")
			}
		})
	}
}
