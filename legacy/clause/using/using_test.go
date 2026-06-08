package using_test

import (
	"testing"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/tests"
	legacyusing "github.com/PaienNate/tdengine-gorm-cgofree/legacy/clause/using"
	"gorm.io/gorm/clause"
)

func TestLegacyUsingAddTagPair(t *testing.T) {
	tests.CheckBuildClauses(t, []clause.Interface{
		clause.Insert{Table: clause.Table{Name: "tb"}},
		legacyusing.SetUsing("stb", map[string]interface{}{
			"tag1": 1,
		}).ADDTagPair("tag2", "string"),
	}, []string{
		"INSERT INTO `tb` USING `stb`(?,?) TAGS(?,?)",
	}, [][][]any{{{"tag1", "tag2", 1, "string"}, {"tag2", "tag1", "string", 1}}})
}

