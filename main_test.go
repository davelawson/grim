package main

import (
	"github.com/davelawson/grim/game"
	"testing"
)

func TestCatalogueConfiguration(t *testing.T) {
	t.Setenv("GRIM_DB", "/tmp/config-test.db")
	t.Setenv("GRIM_SSL", "/tmp/config-test-ssl")
	for _, tc := range []struct{ configured, want string }{{"", game.DefaultCatalogueVersion}, {"edition-two", "edition-two"}} {
		t.Setenv("GRIM_CATALOGUE_VERSION", tc.configured)
		cfg, err := getConfig()
		if err != nil || cfg.CatalogueVersion != tc.want {
			t.Fatal(cfg, err)
		}
	}
}
