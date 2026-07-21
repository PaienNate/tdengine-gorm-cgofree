package tdengine_gorm

import "testing"

func TestOpenEnablesQuotedIdentifiers(t *testing.T) {
	dialect, ok := Open("dsn").(*Dialect)
	if !ok {
		t.Fatalf("Open should return *Dialect")
	}
	if dialect.QuoteIdentifiers == nil || !*dialect.QuoteIdentifiers {
		t.Fatalf("Open should enable quoted identifiers by default")
	}
}

func TestOpenLegacyDisablesQuotedIdentifiers(t *testing.T) {
	dialect, ok := OpenLegacy("dsn").(*Dialect)
	if !ok {
		t.Fatalf("OpenLegacy should return *Dialect")
	}
	if dialect.QuoteIdentifiers == nil || *dialect.QuoteIdentifiers {
		t.Fatalf("OpenLegacy should disable quoted identifiers")
	}
}
