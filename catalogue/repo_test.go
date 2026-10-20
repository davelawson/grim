package catalogue_test

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/davelawson/grim/catalogue"
	"github.com/davelawson/grim/internal/testcatalogue"
)

func database(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	if _, err := db.Exec("create table matches (snapshot text not null)"); err != nil {
		t.Fatal(err)
	}
	testcatalogue.Install(t, db)
	return db
}
func exec(t *testing.T, db *sql.DB, statement string, args ...any) {
	t.Helper()
	if _, err := db.Exec(statement, args...); err != nil {
		t.Fatal(err)
	}
}

func TestOriginalCatalogueParity(t *testing.T) {
	c := testcatalogue.Load(t)
	if len(c.Cards) != 98 {
		t.Fatal("original catalogue size changed", len(c.Cards))
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	// Captured from OpeningCatalogue before the Go definition builders were removed.
	const original = "ee78d152b47a008306cd625e0ce85d7a9df88d504b4d0c76710d6ebadee88ffa"
	if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != original {
		t.Fatalf("original catalogue properties/order changed: %s", got)
	}
	zero, _ := c.Card("flamecraft")
	if value, exists := zero.Cost["wis"]; !exists || value != 0 {
		t.Fatal("explicit zero cost lost")
	}
	if _, exists := zero.Cost["wealth"]; exists {
		t.Fatal("absent cost became an explicit zero")
	}
	rumour, _ := c.Card("rumour_restless_kiln")
	if !slices.Equal(rumour.Challenges[0].Rewards, []string{"basic_wealth", "basic_wealth"}) {
		t.Fatal("repeated rewards lost")
	}
}

func TestLoadColumnsAndVersionIsolation(t *testing.T) {
	db := database(t)
	testcatalogue.CopyVersion(t, db, "playable-v1", "edition-two")
	exec(t, db, `update card_designs set name='A new Study',text='Changed in SQL',cost_wis=4,play_requirement_fire=0,expertise_earth=2,affinity_air=1,income_ephemeral_wis=3,resource_wealth=2,upkeep_wealth=0,eligibility_church=1,support_water=2 where catalogue_version='edition-two' and design_id='study'`)
	updated, err := catalogue.Load(db, "edition-two")
	if err != nil {
		t.Fatal(err)
	}
	card, _ := updated.Card("study")
	if card.Name != "A new Study" || card.Text != "Changed in SQL" || card.Cost["wis"] != 4 || card.Expertise["earth"] != 2 || card.Affinity["air"] != 1 || card.Income["ephemeral_wis"] != 3 || card.Resources["wealth"] != 2 || card.Eligibility["church"] != 1 || card.Support["water"] != 2 {
		t.Fatal("stat columns not mapped", card)
	}
	if v, ok := card.PlayRequirements["fire"]; !ok || v != 0 {
		t.Fatal("explicit zero requirement lost")
	}
	if v, ok := card.Upkeep["wealth"]; !ok || v != 0 {
		t.Fatal("explicit zero upkeep lost")
	}
	original, err := catalogue.Load(db, "playable-v1")
	if err != nil {
		t.Fatal(err)
	}
	old, _ := original.Card("study")
	if old.Name == card.Name || old.Cost["wis"] != 1 {
		t.Fatal("version leaked into original")
	}
	card.Cost["wis"] = 99
	again, err := catalogue.Load(db, "edition-two")
	if err != nil {
		t.Fatal(err)
	}
	fresh, _ := again.Card("study")
	if fresh.Cost["wis"] != 4 {
		t.Fatal("load shares mutable state")
	}
	if _, err := catalogue.Load(db, "missing"); !errors.Is(err, catalogue.ErrNotFound) {
		t.Fatal("missing version fell back", err)
	}
}

func TestSQLConstraintsAndInvalidCatalogueData(t *testing.T) {
	for _, assignment := range []string{"cost_wis=-1", "cost_wis=0.5", "duration=-1", "duration=0.5", "placeholder=2", "name=''", "categories='{}'", "abilities='not json'", "effects='null'", "sort_order=-1"} {
		t.Run(assignment, func(t *testing.T) {
			db := database(t)
			if _, err := db.Exec("update card_designs set " + assignment + " where design_id='study'"); err == nil {
				t.Fatal("invalid column accepted")
			}
		})
	}
	for _, assignment := range []string{`abilities='[{"id":"heal","timing":"action","effect":"heal_activation","unexpected":true}]'`, `rules_version='future'`, `sort_order=1000`} {
		t.Run(assignment, func(t *testing.T) {
			db := database(t)
			exec(t, db, "update card_designs set "+assignment+" where design_id='study'")
			if _, err := catalogue.Load(db, "playable-v1"); err == nil {
				t.Fatal("invalid catalogue accepted")
			}
		})
	}
}

func TestUsedVersionIsFrozenAndInstallationIsIdempotent(t *testing.T) {
	db := database(t)
	testcatalogue.CopyVersion(t, db, "playable-v1", "edition-two")
	exec(t, db, `insert into matches(snapshot) values('not-json'),('{"catalogueversion":"playable-v1"}')`)
	for _, statement := range []string{
		`update card_designs set cost_wis=99 where catalogue_version='playable-v1' and design_id='study'`,
		`delete from card_designs where catalogue_version='playable-v1' and design_id='study'`,
		`insert into card_designs(catalogue_version,design_id,rules_version,sort_order,name) values('playable-v1','new-card','playable-v1',100,'New card')`,
		`update card_designs set catalogue_version='playable-v1' where catalogue_version='edition-two' and design_id='study'`,
		`update card_designs set catalogue_version='other' where catalogue_version='playable-v1' and design_id='study'`,
	} {
		if _, err := db.Exec(statement); err == nil {
			t.Fatal("frozen version changed", statement)
		}
	}
	exec(t, db, `update card_designs set cost_wis=7 where catalogue_version='edition-two' and design_id='study'`)
	testcatalogue.Install(t, db)
	c, err := catalogue.Load(db, "edition-two")
	if err != nil {
		t.Fatal(err)
	}
	card, _ := c.Card("study")
	if card.Cost["wis"] != 7 {
		t.Fatal("installation overwrote edited data")
	}
	var count int
	if err := db.QueryRow("select count(*) from card_designs").Scan(&count); err != nil || count != 196 {
		t.Fatal("installation changed row count", count, err)
	}
}
