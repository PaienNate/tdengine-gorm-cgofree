package tdengine_gorm

import (
	"strings"
	"testing"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/using"
	"gorm.io/gorm"
)

func TestDialectQuotesIdentifiersByDefault(t *testing.T) {
	db, err := gorm.Open(&Dialect{
		DriverName:        "taosWS",
		DSN:               "user:secret@ws(example-host:6041)/gorm_test?loc=Local",
		InterpolateParams: WithInterpolateParams(false),
	}, &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open dialect: %v", err)
	}

	stmt := db.Table("TbName").
		Clauses(using.SetUsing("StbName", map[string]any{"TagName": "TagVal"})).
		Create(map[string]any{
			"TsColumn": "2026-07-21T00:00:00Z",
			"Value":    "abc",
		}).Statement

	gotSQL := stmt.SQL.String()
	for _, want := range []string{
		"INSERT INTO `TbName`",
		"USING `StbName`",
		"(`TsColumn`,`Value`)",
	} {
		if !strings.Contains(gotSQL, want) {
			t.Fatalf("expected SQL to contain %q, got %q", want, gotSQL)
		}
	}
}

func TestDialectCanDisableIdentifierQuoting(t *testing.T) {
	db, err := gorm.Open(&Dialect{
		DriverName:        "taosWS",
		DSN:               "user:secret@ws(example-host:6041)/gorm_test?loc=Local",
		InterpolateParams: WithInterpolateParams(false),
		QuoteIdentifiers:  WithQuotedIdentifiers(false),
	}, &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open dialect: %v", err)
	}

	stmt := db.Table("TbName").
		Clauses(using.SetUsing("StbName", map[string]any{"TagName": "TagVal"})).
		Create(map[string]any{
			"TsColumn": "2026-07-21T00:00:00Z",
			"Value":    "abc",
		}).Statement

	gotSQL := stmt.SQL.String()
	for _, want := range []string{
		"INSERT INTO TbName",
		"USING StbName",
		"(TsColumn,Value)",
	} {
		if !strings.Contains(gotSQL, want) {
			t.Fatalf("expected SQL to contain %q, got %q", want, gotSQL)
		}
	}
	for _, unwanted := range []string{"`TbName`", "`StbName`", "`TsColumn`", "`Value`"} {
		if strings.Contains(gotSQL, unwanted) {
			t.Fatalf("expected SQL not to contain %q, got %q", unwanted, gotSQL)
		}
	}
}

func TestLegacyDryRunComparisonDocumentsRemainingDifferences(t *testing.T) {
	db, err := gorm.Open(&Dialect{
		DriverName:        "taosWS",
		DSN:               "user:secret@ws(example-host:6041)/gorm_test?loc=Local",
		InterpolateParams: WithInterpolateParams(false),
		QuoteIdentifiers:  WithQuotedIdentifiers(false),
	}, &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open dialect: %v", err)
	}

	stmt := db.Clauses(
		using.SetUsing("StbName", map[string]any{"TagName": "TagVal"}),
	).Table("TbName").Create(map[string]any{
		"ts":    "2020-01-01T00:00:00Z",
		"value": "abc",
	}).Statement

	gotSQL := stmt.SQL.String()
	if strings.Contains(gotSQL, "`") {
		t.Fatalf("expected legacy-compat SQL to avoid backticks, got %q", gotSQL)
	}
	if !strings.Contains(gotSQL, "VALUES (?,?)") {
		t.Fatalf("expected placeholders to remain with interpolation disabled, got %q", gotSQL)
	}
}
