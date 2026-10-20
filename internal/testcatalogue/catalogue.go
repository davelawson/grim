// Package testcatalogue installs the production schema and catalogue for tests.
package testcatalogue

import (
	"database/sql"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/catalogue"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Install(t testing.TB, db *sql.DB) {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	data, err := os.ReadFile(filepath.Join(filepath.Dir(file), "../../sql/install-card-designs.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(data)); err != nil {
		t.Fatal(err)
	}
}
func Load(t testing.TB) *api.Catalogue {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	// Catalogue freeze guards reference match snapshots, which the pure engine tests do not otherwise need.
	if _, err := db.Exec("create table matches (snapshot text not null)"); err != nil {
		t.Fatal(err)
	}
	Install(t, db)
	c, err := catalogue.Load(db, "playable-v1")
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// CopyVersion creates a new editable catalogue version for integration tests.
func CopyVersion(t testing.TB, db *sql.DB, from, to string) {
	t.Helper()
	rows, err := db.Query("pragma table_info(card_designs)")
	if err != nil {
		t.Fatal(err)
	}
	columns := []string{}
	for rows.Next() {
		var index, notnull, key int
		var name, typ string
		var defaultValue any
		if err := rows.Scan(&index, &name, &typ, &notnull, &defaultValue, &key); err != nil {
			t.Fatal(err)
		}
		columns = append(columns, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	rows.Close()
	values := append([]string{}, columns...)
	values[0] = "?"
	if _, err := db.Exec("insert into card_designs ("+strings.Join(columns, ",")+") select "+strings.Join(values, ",")+" from card_designs where catalogue_version=?", to, from); err != nil {
		t.Fatal(err)
	}
}
