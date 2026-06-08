package tdengine_gorm

import (
	"strings"
	"testing"

	legacyusing "github.com/PaienNate/tdengine-gorm-cgofree/legacy/clause/using"
	"gorm.io/gorm"
)

func TestWithInterpolateParamsFalseDisablesCurrentInterpolation(t *testing.T) {
	if interpolateParamsEnabled(Dialect{InterpolateParams: WithInterpolateParams(false)}) {
		t.Fatal("expected false to disable interpolation")
	}
}

func TestLegacyClauseCanBeUsedWithCurrentDialect(t *testing.T) {
	db, err := gorm.Open(&Dialect{
		DriverName:        "taosWS",
		DSN:               "user:secret@ws(example-host:6041)/gorm_test?loc=Local",
		InterpolateParams: WithInterpolateParams(false),
	}, &gorm.Config{
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("open dialect: %v", err)
	}

	stmt := db.Table("tb_legacy").
		Clauses(legacyusing.SetUsing("stb_legacy", map[string]interface{}{"tbn": "tb_legacy"})).
		Create(map[string]any{
			"ts":    "2026-06-08T20:30:01Z",
			"value": "abc",
		}).Statement

	gotSQL := stmt.SQL.String()
	if !strings.Contains(gotSQL, "USING `stb_legacy`") {
		t.Fatalf("expected legacy USING clause in SQL, got %q", gotSQL)
	}
	if !strings.Contains(gotSQL, "VALUES (?,?)") {
		t.Fatalf("expected placeholders to remain when interpolation is disabled, got %q", gotSQL)
	}
}
