package create

import (
	"strconv"

	"gorm.io/gorm/clause"
)

type CreateTable struct {
	tables []*Table
}

const (
	STableType = iota + 1
	CommonTableType
)

type Table struct {
	TableType   int
	Table       string
	IfNotExists bool
	STable      string
	Tags        map[string]interface{}
	Column      []*Column
	TagColumn   []*Column
}

type Column struct {
	Name       string
	ColumnType string
	Length     uint64
}

const (
	TimestampType = "TIMESTAMP"
	IntType       = "INT"
	BigIntType    = "BIGINT"
	FloatType     = "FLOAT"
	DoubleType    = "DOUBLE"
	BinaryType    = "BINARY"
	SmallIntType  = "SMALLINT"
	TinyIntType   = "TINYINT"
	BoolType      = "BOOL"
	NCharType     = "NCHAR"
)

func NewTable(name string, ifNotExist bool, column []*Column, stable string, tags map[string]interface{}) *Table {
	return &Table{
		TableType:   CommonTableType,
		Table:       name,
		IfNotExists: ifNotExist,
		STable:      stable,
		Tags:        tags,
		Column:      column,
	}
}

func NewSTable(name string, ifNotExists bool, column []*Column, tagColumn []*Column) *Table {
	return &Table{
		TableType:   STableType,
		Table:       name,
		IfNotExists: ifNotExists,
		Column:      column,
		TagColumn:   tagColumn,
	}
}

func NewCreateTableClause(tables []*Table) CreateTable {
	return CreateTable{tables: tables}
}

func (c CreateTable) AddTables(tables ...*Table) CreateTable {
	c.tables = append(c.tables, tables...)
	return c
}

func (CreateTable) Name() string {
	return "CREATE TABLE"
}

func (c CreateTable) Build(builder clause.Builder) {
	for _, table := range c.tables {
		if table == nil {
			continue
		}

		switch table.TableType {
		case CommonTableType:
			_, _ = builder.WriteString("CREATE TABLE ")
		case STableType:
			_, _ = builder.WriteString("CREATE STABLE ")
		default:
			return
		}

		if table.IfNotExists {
			_, _ = builder.WriteString("IF NOT EXISTS ")
		}
		builder.WriteQuoted(table.Table)

		if table.TableType == CommonTableType && table.STable != "" {
			_, _ = builder.WriteString(" USING ")
			builder.WriteQuoted(table.STable)

			tagValueList := make([]interface{}, 0, len(table.Tags))
			index := 0
			_ = builder.WriteByte('(')
			for tagName, tagValue := range table.Tags {
				if index > 0 {
					_ = builder.WriteByte(',')
				}
				builder.WriteQuoted(tagName)
				tagValueList = append(tagValueList, tagValue)
				index++
			}
			_, _ = builder.WriteString(") TAGS ")
			builder.AddVar(builder, tagValueList)
		} else {
			_, _ = builder.WriteString(" (")
			for i, column := range table.Column {
				if i > 0 {
					_ = builder.WriteByte(',')
				}
				column.build(builder)
			}
			_ = builder.WriteByte(')')
		}

		if table.TableType == STableType {
			_, _ = builder.WriteString(" TAGS(")
			for i, tagColumn := range table.TagColumn {
				if i > 0 {
					_ = builder.WriteByte(',')
				}
				tagColumn.build(builder)
			}
			_ = builder.WriteByte(')')
		}
	}
}

func (c CreateTable) MergeClause(clause *clause.Clause) {
	clause.Name = ""
	clause.Expression = c
}

func (c *Column) build(builder clause.Builder) {
	builder.WriteQuoted(c.Name)
	_ = builder.WriteByte(' ')
	_, _ = builder.WriteString(c.ColumnType)
	if c.ColumnType == NCharType || c.ColumnType == BinaryType {
		_ = builder.WriteByte('(')
		_, _ = builder.WriteString(strconv.FormatUint(c.Length, 10))
		_ = builder.WriteByte(')')
	}
}
