package game

import (
	"encoding/json"
	"testing"

	"github.com/davelawson/grim/api"
)

func TestAutomatedMetadataCompatibilityAndPrivacy(t *testing.T) {
	definitions := testDefinitions(t)
	legacy := setup(t, 2)
	if _, err := Decode(jsonBytes(t, legacy), definitions); err != nil {
		t.Fatalf("legacy snapshot rejected: %v", err)
	}
	s, err := NewSetup("match", "Bots", "lobby", []string{"human", "bot-example"}, definitions)
	if err != nil {
		t.Fatal(err)
	}
	s.Bots = map[string]api.Bot{"bot-example": {Name: "Bot 1"}}
	data := jsonBytes(t, s)
	restored, err := Decode(data, definitions)
	if err != nil {
		t.Fatal(err)
	}
	view := restored.View("human")
	if view.You.UserID != "human" || view.Bots["bot-example"].Name != "Bot 1" {
		t.Fatal("incorrect participant projection")
	}
	view.Bots["bot-example"] = api.Bot{Name: "changed"}
	if restored.Bots["bot-example"].Name != "Bot 1" {
		t.Fatal("view aliases saved opponent metadata")
	}
	for _, bots := range []map[string]api.Bot{
		nil,
		{"missing": {Name: "Bot 1"}},
		{"human": {Name: "Bot 1"}},
		{"bot-example": {Name: ""}},
	} {
		var bad State
		if err := json.Unmarshal(data, &bad); err != nil {
			t.Fatal(err)
		}
		bad.Bots = bots
		if _, err := Decode(jsonBytes(t, &bad), definitions); err == nil {
			t.Fatal("invalid automated metadata was accepted")
		}
	}
}
