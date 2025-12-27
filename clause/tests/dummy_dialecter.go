package tests

import (
	"github.com/PaienNate/tdengine-gorm-cgofree/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DummyDialector struct{}

func (DummyDialector) Name() string {
	return "dummy"
}

func (DummyDialector) Initialize(*gorm.DB) error {
	return nil
}

func (DummyDialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return clause.Expr{SQL: "NULL"}
}

func (DummyDialector) Migrator(*gorm.DB) gorm.Migrator {
	return nil
}

func (DummyDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v any) {
	writer.WriteByte('?')
}

func (DummyDialector) QuoteTo(writer clause.Writer, str string) {
	utils.QuoteTo(writer, str)
}

func (DummyDialector) Explain(sql string, vars ...any) string {
	return logger.ExplainSQL(sql, nil, `'`, vars...)
}

func (DummyDialector) DataTypeOf(*schema.Field) string {
	return ""
}
