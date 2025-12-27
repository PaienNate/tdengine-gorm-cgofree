package create_test

import (
	"testing"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/create"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/tests"

	"gorm.io/gorm/clause"
)

func Test_CreateTable(t *testing.T) {
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			"create table use stable",
			[]clause.Interface{
				create.NewCreateTable(
					create.NewCTableBuilder("t_1").
						IfNotExists().
						BuildWithSTable(
							"st_1",
							map[string]any{
								"tag_int":    1,
								"tag_string": "string",
							}),
				)},
			[]string{
				"CREATE TABLE IF NOT EXISTS `t_1` USING `st_1`(`tag_int`,`tag_string`) TAGS (?,?)",
				"CREATE TABLE IF NOT EXISTS `t_1` USING `st_1`(`tag_string`,`tag_int`) TAGS (?,?)",
			},
			[][][]any{{{1, "string"}}, {{"string", 1}}},
		},
		{
			"create table without stable",
			[]clause.Interface{
				create.NewCreateTable().
					AddTables(
						create.NewCTableBuilder("t_1").
							IfNotExists().
							Columns(
								[]*create.Column{
									{
										Name:   "ts",
										Type:   create.Timestamp,
										Length: 0,
									},
									{
										Name:   "c_int",
										Type:   create.Int,
										Length: 0,
									},
									{
										Name:   "c_bigint",
										Type:   create.BigInt,
										Length: 0,
									},
									{
										Name:   "c_float",
										Type:   create.Float,
										Length: 0,
									},
									{
										Name:   "c_double",
										Type:   create.Double,
										Length: 0,
									},
									{
										Name:   "c_binary",
										Type:   create.Binary,
										Length: 128,
									},
									{
										Name:   "c_smallint",
										Type:   create.SmallInt,
										Length: 0,
									},
									{
										Name:   "c_tinyint",
										Type:   create.TinyInt,
										Length: 0,
									},
									{
										Name:   "c_bool",
										Type:   create.Bool,
										Length: 0,
									},
									{
										Name:   "c_nchar",
										Type:   create.NChar,
										Length: 128,
									},
								}...,
							).
							Build(),
					),
			},
			[]string{
				"CREATE TABLE IF NOT EXISTS `t_1` (`ts` TIMESTAMP,`c_int` INT,`c_bigint` BIGINT,`c_float` FLOAT,`c_double` DOUBLE,`c_binary` BINARY(128),`c_smallint` SMALLINT,`c_tinyint` TINYINT,`c_bool` BOOL,`c_nchar` NCHAR(128))",
			},
			nil,
		},
		{
			"create stable",
			[]clause.Interface{
				create.NewCreateTable(
					create.NewSTableBuilder("st_1").
						IfNotExists().
						Columns(
							[]*create.Column{
								{
									Name:   "ts",
									Type:   create.Timestamp,
									Length: 0,
								},
								{
									Name:   "c_int",
									Type:   create.Int,
									Length: 0,
								},
								{
									Name:   "c_bigint",
									Type:   create.BigInt,
									Length: 0,
								},
								{
									Name:   "c_float",
									Type:   create.Float,
									Length: 0,
								},
								{
									Name:   "c_double",
									Type:   create.Double,
									Length: 0,
								},
								{
									Name:   "c_binary",
									Type:   create.Binary,
									Length: 128,
								},
								{
									Name:   "c_smallint",
									Type:   create.SmallInt,
									Length: 0,
								},
								{
									Name:   "c_tinyint",
									Type:   create.TinyInt,
									Length: 0,
								},
								{
									Name:   "c_bool",
									Type:   create.Bool,
									Length: 0,
								},
								{
									Name:   "c_nchar",
									Type:   create.NChar,
									Length: 128,
								},
							}...,
						).
						TagColumns(
							[]*create.Column{
								{
									Name:   "t_int",
									Type:   create.Int,
									Length: 0,
								},
								{
									Name:   "c_nchar",
									Type:   create.NChar,
									Length: 128,
								},
							}...,
						).
						Build(),
				)},
			[]string{
				"CREATE STABLE IF NOT EXISTS `st_1` (`ts` TIMESTAMP,`c_int` INT,`c_bigint` BIGINT,`c_float` FLOAT,`c_double` DOUBLE,`c_binary` BINARY(128),`c_smallint` SMALLINT,`c_tinyint` TINYINT,`c_bool` BOOL,`c_nchar` NCHAR(128)) TAGS(`t_int` INT,`c_nchar` NCHAR(128))",
				"CREATE STABLE IF NOT EXISTS `st_1` (`ts` TIMESTAMP,`c_int` INT,`c_bigint` BIGINT,`c_float` FLOAT,`c_double` DOUBLE,`c_binary` BINARY(128),`c_smallint` SMALLINT,`c_tinyint` TINYINT,`c_bool` BOOL,`c_nchar` NCHAR(128)) TAGS(`c_nchar` NCHAR(128),`t_int` INT)",
			},
			nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}
