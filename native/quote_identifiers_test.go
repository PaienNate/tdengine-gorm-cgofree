//go:build cgo

package native

import (
	"strings"
	"testing"

	"gorm.io/gorm"
)

func TestDialectCanDisableIdentifierQuoting(t *testing.T) {
	db, err := gorm.Open(&Dialect{
		DriverName:        DefaultDriverName,
		DSN:               "user:secret@tcp(example-host:6030)/gorm_test?loc=Local",
		InterpolateParams: WithInterpolateParams(false),
		QuoteIdentifiers:  WithQuotedIdentifiers(false),
	}, &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open dialect: %v", err)
	}

	stmt := db.Table("TbName").Create(map[string]any{
		"TsColumn": "2026-07-21T00:00:00Z",
		"Value":    "abc",
	}).Statement

	gotSQL := stmt.SQL.String()
	if !strings.Contains(gotSQL, "INSERT INTO TbName") {
		t.Fatalf("expected SQL without quoted table name, got %q", gotSQL)
	}
	if strings.Contains(gotSQL, "`TbName`") {
		t.Fatalf("expected SQL without backticks, got %q", gotSQL)
	}
}
