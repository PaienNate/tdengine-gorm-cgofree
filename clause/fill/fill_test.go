package fill_test

import (
	"testing"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/fill"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/tests"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/window"

	"gorm.io/gorm/clause"
)

func Test_Fill(t *testing.T) {
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			Name: "",
			Clauses: []clause.Interface{
				clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
				clause.From{Tables: []clause.Table{{Name: "t_1"}}},
				window.SetInterval(window.Duration{Value: 10, Unit: window.Minute}),
				fill.Fill{fill.FillValue, 12},
			},
			Result: []string{"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(10m) FILL (VALUE,12)"},
			Vars:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}
