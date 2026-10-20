package game

import (
	"github.com/davelawson/grim/internal/testcatalogue"
	"testing"
)

func testCatalogue(t testing.TB) *Catalogue { t.Helper(); return testcatalogue.Load(t) }
func testDefinitions(t testing.TB) *Definitions {
	t.Helper()
	d, err := NewDefinitions(*testCatalogue(t))
	if err != nil {
		t.Fatal(err)
	}
	return d
}
