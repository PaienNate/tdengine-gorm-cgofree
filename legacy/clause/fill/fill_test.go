package fill_test

import (
	"testing"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/tests"
	legacyfill "github.com/PaienNate/tdengine-gorm-cgofree/legacy/clause/fill"
	legacywindow "github.com/PaienNate/tdengine-gorm-cgofree/legacy/clause/window"
	"gorm.io/gorm/clause"
)

func TestLegacyFillSetValue(t *testing.T) {
	tests.CheckBuildClauses(t, []clause.Interface{
		clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
		clause.From{Tables: []clause.Table{{Name: "t_1"}}},
		legacywindow.SetInterval(legacywindow.Duration{Value: 10, Unit: legacywindow.Minute}),
		legacyfill.SetFill(legacyfill.FillValue).SetValue(12),
	}, []string{
		"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(10m) FILL (VALUE,12)",
	}, nil)
}

