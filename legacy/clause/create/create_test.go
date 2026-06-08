package create_test

import (
	"testing"

	legacycreate "github.com/PaienNate/tdengine-gorm-cgofree/legacy/clause/create"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/tests"
	"gorm.io/gorm/clause"
)

func TestLegacyCreateTableClauseCommonTableUsingStable(t *testing.T) {
	clauses := []clause.Interface{
		legacycreate.NewCreateTableClause([]*legacycreate.Table{
			legacycreate.NewTable("t_1", true, nil, "st_1", map[string]interface{}{
				"tag_int":    1,
				"tag_string": "string",
			}),
		}),
	}

	results := []string{
		"CREATE TABLE IF NOT EXISTS `t_1` USING `st_1`(`tag_int`,`tag_string`) TAGS (?,?)",
		"CREATE TABLE IF NOT EXISTS `t_1` USING `st_1`(`tag_string`,`tag_int`) TAGS (?,?)",
	}
	vars := [][][]any{{{1, "string"}}, {{"string", 1}}}

	tests.CheckBuildClauses(t, clauses, results, vars)
}

func TestLegacyCreateTableClauseCommonTableColumns(t *testing.T) {
	clauses := []clause.Interface{
		legacycreate.NewCreateTableClause(nil).AddTables(&legacycreate.Table{
			TableType:   legacycreate.CommonTableType,
			Table:       "t_1",
			IfNotExists: true,
			Column: []*legacycreate.Column{
				{Name: "ts", ColumnType: legacycreate.TimestampType},
				{Name: "c_binary", ColumnType: legacycreate.BinaryType, Length: 128},
				{Name: "c_nchar", ColumnType: legacycreate.NCharType, Length: 128},
			},
		}),
	}

	tests.CheckBuildClauses(t, clauses, []string{
		"CREATE TABLE IF NOT EXISTS `t_1` (`ts` TIMESTAMP,`c_binary` BINARY(128),`c_nchar` NCHAR(128))",
	}, nil)
}

func TestLegacyCreateTableClauseStableAndConstructors(t *testing.T) {
	stable := legacycreate.NewSTable("st_1", true, []*legacycreate.Column{
		{Name: "ts", ColumnType: legacycreate.TimestampType},
		{Name: "value", ColumnType: legacycreate.DoubleType},
	}, []*legacycreate.Column{
		{Name: "t_int", ColumnType: legacycreate.IntType},
	})

	clauses := []clause.Interface{
		legacycreate.NewCreateTableClause([]*legacycreate.Table{stable}),
	}

	tests.CheckBuildClauses(t, clauses, []string{
		"CREATE STABLE IF NOT EXISTS `st_1` (`ts` TIMESTAMP,`value` DOUBLE) TAGS(`t_int` INT)",
	}, nil)
}

