//go:build cgo

package native

import (
	"github.com/PaienNate/tdengine-gorm-cgofree/internal/base"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Migrator = base.Migrator
type Column = base.Column

var DefaultDriverName = "taosSql"

type Dialect struct {
	DriverName        string
	DSN               string
	Conn              gorm.ConnPool
	InterpolateParams *bool
	QuoteIdentifiers  *bool
}

func Open(dsn string) gorm.Dialector {
	return &Dialect{DriverName: DefaultDriverName, DSN: dsn}
}

func WithInterpolateParams(enabled bool) *bool {
	return base.WithInterpolateParams(enabled)
}

func WithQuotedIdentifiers(enabled bool) *bool {
	return base.WithQuotedIdentifiers(enabled)
}

func (Dialect) Name() string {
	return base.Dialect{}.Name()
}

func (dialect Dialect) Initialize(db *gorm.DB) error {
	baseDialect := dialect.baseDialect()
	if baseDialect.DriverName == "" {
		baseDialect.DriverName = DefaultDriverName
	}
	return baseDialect.Initialize(db)
}

func (Dialect) ClauseBuilders() map[string]clause.ClauseBuilder {
	return base.Dialect{}.ClauseBuilders()
}

func (Dialect) DefaultValueOf(field *schema.Field) clause.Expression {
	return base.Dialect{}.DefaultValueOf(field)
}

func (dialect Dialect) Migrator(db *gorm.DB) gorm.Migrator {
	baseDialect := dialect.baseDialect()
	if baseDialect.DriverName == "" {
		baseDialect.DriverName = DefaultDriverName
	}
	return baseDialect.Migrator(db)
}

func (Dialect) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v any) {
	base.Dialect{}.BindVarTo(writer, stmt, v)
}

func (dialect Dialect) QuoteTo(writer clause.Writer, str string) {
	dialect.baseDialect().QuoteTo(writer, str)
}

func (Dialect) Explain(sql string, vars ...any) string {
	return base.Dialect{}.Explain(sql, vars...)
}

func (Dialect) DataTypeOf(field *schema.Field) string {
	return base.Dialect{}.DataTypeOf(field)
}

func (Dialect) SavePoint(tx *gorm.DB, name string) error {
	return base.Dialect{}.SavePoint(tx, name)
}

func (Dialect) RollbackTo(tx *gorm.DB, name string) error {
	return base.Dialect{}.RollbackTo(tx, name)
}

func (dialect Dialect) baseDialect() base.Dialect {
	return base.Dialect{
		DriverName:        dialect.DriverName,
		DSN:               dialect.DSN,
		Conn:              dialect.Conn,
		InterpolateParams: dialect.InterpolateParams,
		QuoteIdentifiers:  dialect.QuoteIdentifiers,
	}
}
