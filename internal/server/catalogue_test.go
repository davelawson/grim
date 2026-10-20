package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/internal/testcatalogue"
	"github.com/davelawson/grim/match"
	"github.com/davelawson/grim/user"
	"github.com/davelawson/grim/util"
	"github.com/google/uuid"
)

func (f *lifecycleFixture) catalogue(id string) *api.Catalogue {
	f.t.Helper()
	var result api.CatalogueResponse
	if err := json.Unmarshal(f.call(200, "GET", api.GetCataloguePath(id), "owner", "").Body.Bytes(), &result); err != nil || result.Catalogue == nil {
		f.t.Fatal("missing catalogue", err)
	}
	return result.Catalogue
}

func TestCatalogueVersionsDriveMatchesIndependently(t *testing.T) {
	f, old := reactionFixture(t)
	testcatalogue.CopyVersion(t, f.db, "playable-v1", "edition-two")
	f.exec(`update card_designs set name='Database Study',cost_wis=4 where catalogue_version='edition-two' and design_id='study'`)
	f.exec(`update card_designs set name='Database Chantry' where catalogue_version='edition-two' and design_id='chantry_fire'`)
	f.exec(`update card_designs set name='Database Attack' where catalogue_version='edition-two' and design_id='subjugate_rival'`)
	f.exec(`update card_designs set abilities=json_set(abilities,'$[0].cost.wis',5) where catalogue_version='edition-two' and design_id='spell_healing_spring'`)
	// Draft receipts from the original version must not adopt the newly configured edition.
	lobby := f.createLobby(2)
	f.ready(lobby, 2)
	draft, _, _ := f.launch(lobby)
	revision := draft.Revision
	pick := gameplayBody(t, api.DraftPickRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, Type: "choose_chantry", CardID: "chantry_fire"})
	recorded := f.call(200, "POST", api.SubmitDraftPath(draft.ID), "owner", pick).Body.String()
	f.catalogueVersion = "edition-two"
	f.restart()
	if f.call(200, "POST", api.SubmitDraftPath(draft.ID), "owner", pick).Body.String() != recorded {
		t.Fatal("draft replay used configured current version")
	}
	oldCard, _ := f.catalogue(old.ID).Card("study")
	if oldCard.Cost["wis"] != 1 || oldCard.Name != "Study" {
		t.Fatal("existing match changed catalogue")
	}
	// A newly launched match pins the configured edition and exposes its database values.
	fresh := f.playable(2)
	if fresh.CatalogueVersion != "edition-two" {
		t.Fatal("new match pinned wrong edition")
	}
	freshCard, _ := f.catalogue(fresh.ID).Card("study")
	if freshCard.Cost["wis"] != 4 || freshCard.Name != "Database Study" {
		t.Fatal("catalogue did not come from database")
	}
	for _, entry := range []struct {
		id, card string
		cost     int
	}{{old.ID, "test/study", 1}, {fresh.ID, "new/study", 4}} {
		_, s := f.saved(entry.id)
		owner := f.ids["owner"]
		w := s.Wizards[owner]
		s.Turn = &api.Turn{Number: 10, ActiveWizard: owner, Phase: "actions", Step: "actions"}
		s.Choice = nil
		s.Exchange = nil
		s.Resolving = nil
		s.PhaseEffects = nil
		w.Hand = []api.CardInstance{{ID: entry.card, DesignID: "study", Owner: owner, Location: "hand"}}
		w.Resources = map[string]int{"wis": 10}
		f.exec("update matches set snapshot=? where id=?", string(mustJSON(t, s)), entry.id)
		result := f.command(entry.id, "owner", game.Command{Type: "play_card", CardID: entry.card})
		if result.Wizards[owner].Resources["wis"] != 10-entry.cost {
			t.Fatal("play cost ignored pinned database definition")
		}
	}
	// Activated-ability costs are also definition data, rather than duplicated Go values.
	_, s := f.saved(fresh.ID)
	owner := f.ids["owner"]
	w := s.Wizards[owner]
	s.Choice = nil
	s.Exchange = nil
	s.Resolving = nil
	s.PhaseEffects = nil
	w.Resources = map[string]int{"wis": 10}
	w.InPlay = append(w.InPlay, api.CardInstance{ID: "db/spring", DesignID: "spell_healing_spring", Owner: owner, Location: "inplay"})
	f.exec("update matches set snapshot=? where id=?", string(mustJSON(t, s)), fresh.ID)
	activated := f.command(fresh.ID, "owner", game.Command{Type: "activate_ability", CardID: "db/spring", AbilityID: "heal"})
	if activated.Wizards[owner].Resources["wis"] != 5 {
		t.Fatal("activation ignored database cost")
	}
	if _, err := f.db.Exec(`update card_designs set cost_wis=99 where catalogue_version='edition-two' and design_id='study'`); err == nil {
		t.Fatal("active edition remained editable")
	}
	f.restart()
	if f.current(old.ID, "owner").CatalogueVersion != "playable-v1" || f.current(fresh.ID, "owner").CatalogueVersion != "edition-two" {
		t.Fatal("restart lost edition pins")
	}
}

func tableContents(t *testing.T, db *sql.DB, name string) []byte {
	t.Helper()
	rows, err := db.Query("select * from " + name + " order by rowid")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}
	records := [][]any{}
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			t.Fatal(err)
		}
		records = append(records, values)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return mustJSON(t, records)
}

func TestAdditiveCatalogueInstallPreservesMatchesAndReceipts(t *testing.T) {
	f, view := reactionFixture(t)
	id := view.ID
	revision := view.Revision
	request := api.PlayCardRequest{RequestID: uuid.NewString(), ExpectedRevision: &revision, CardID: "test/study"}
	body := gameplayBody(t, request)
	accepted := f.call(200, "POST", api.PlayCardPath(id), "owner", body).Body.String()
	before := map[string][]byte{}
	for _, name := range []string{"users", "lobbies", "lobby_users", "matches", "players", "command_receipts", "launch_receipts"} {
		before[name] = tableContents(t, f.db, name)
	}
	// Represent an existing pre-catalogue database without altering its saved games.
	f.exec("drop table card_designs")
	f.restart()
	f.call(500, "GET", api.GetCataloguePath(id), "owner", "")
	if f.call(200, "POST", api.PlayCardPath(id), "owner", body).Body.String() != accepted {
		t.Fatal("gameplay replay required current definitions")
	}
	testcatalogue.Install(t, f.db)
	testcatalogue.Install(t, f.db)
	for name, content := range before {
		if !bytes.Equal(content, tableContents(t, f.db, name)) {
			t.Fatal("additive installation changed", name)
		}
	}
	f.restart()
	if f.call(200, "POST", api.PlayCardPath(id), "owner", body).Body.String() != accepted {
		t.Fatal("installation changed accepted receipt")
	}
	if f.catalogue(id).CatalogueVersion != "playable-v1" {
		t.Fatal("installation lost original version")
	}
	if f.current(id, "owner").ID != id {
		t.Fatal("existing match cannot resume")
	}
}

func TestConfiguredCatalogueValidationAndMissingPinnedVersion(t *testing.T) {
	f, view := reactionFixture(t)
	check := func(version string) error {
		service := match.NewService(match.NewRepo(f.db), user.NewUserRepo(f.db), version)
		return util.InTx(f.db, func(tx *sql.Tx) error { return service.ValidateCatalogue(tx) })()
	}
	if err := check("playable-v1"); err != nil {
		t.Fatal(err)
	}
	if err := check("missing"); err == nil {
		t.Fatal("missing current version accepted")
	}
	testcatalogue.CopyVersion(t, f.db, "playable-v1", "broken")
	f.exec(`update card_designs set effects='["not-implemented"]' where catalogue_version='broken' and design_id='study'`)
	if err := check("broken"); err == nil {
		t.Fatal("invalid current version accepted")
	}
	// Explicitly missing historical versions are rejected rather than falling back.
	_, state := f.saved(view.ID)
	state.CatalogueVersion = "missing"
	f.exec("update matches set snapshot=? where id=?", string(mustJSON(t, state)), view.ID)
	errorCode(t, f.call(409, "GET", api.GetMatchPath(view.ID), "owner", ""), "unsupported_state")
	// Validation never silently repairs or rewrites the installed definitions.
	var effects string
	if err := f.db.QueryRow(`select effects from card_designs where catalogue_version='broken' and design_id='study'`).Scan(&effects); err != nil {
		t.Fatal(err)
	}
	if effects != `["not-implemented"]` {
		t.Fatal(fmt.Sprintf("invalid definition was rewritten: %s", effects))
	}
}
