package slimit_test

import (
	"testing"

	legacyslimit "github.com/PaienNate/tdengine-gorm-cgofree/legacy/clause/slimit"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/tests"
	"gorm.io/gorm/clause"
)

func TestLegacySLimit(t *testing.T) {
	tests.CheckBuildClauses(t, []clause.Interface{
		clause.Select{},
		clause.From{},
		legacyslimit.SetSLimit(10, 20),
		legacyslimit.SetSLimit(0, 30),
		legacyslimit.SetSLimit(50, 0),
	}, []string{
		"SELECT * FROM `users` SLIMIT 50 SOFFSET 30",
	}, [][][]any{{nil}})
}

