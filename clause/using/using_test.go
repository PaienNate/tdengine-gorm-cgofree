package using_test

import (
	"testing"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/tests"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/using"
	"gorm.io/gorm/clause"
)

func Test_Using(t *testing.T) {
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			Name: "USING",
			Clauses: []clause.Interface{
				clause.Insert{
					Table: clause.Table{Name: "tb"},
				},
				using.SetUsing("stb", map[string]any{"tag1": 1}).
					AddTag("tag2", "string"),
			},
			Result: []string{
				"INSERT INTO `tb` USING `stb`(?,?) TAGS(?,?)",
			},
			Vars: [][][]any{{{"tag1", "tag2", 1, "string"}, {"tag2", "tag1", "string", 1}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}
