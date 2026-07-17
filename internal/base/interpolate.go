package base

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func InterpolateParamsEnabled(dialect Dialect) bool {
	return interpolateParamsEnabled(dialect)
}

func NewInterpolatingConnPool(conn gorm.ConnPool, enabled bool) gorm.ConnPool {
	return newInterpolatingConnPool(conn, enabled)
}

func InterpolateTDengineParams(query string, args []any) (string, error) {
	return interpolateTDengineParams(query, args)
}

func interpolateParamsEnabled(dialect Dialect) bool {
	if dialect.InterpolateParams == nil {
		return true
	}
	return *dialect.InterpolateParams
}

type interpolatingConnPool struct {
	conn    gorm.ConnPool
	enabled bool
}

func newInterpolatingConnPool(conn gorm.ConnPool, enabled bool) gorm.ConnPool {
	if !enabled {
		return conn
	}
	return &interpolatingConnPool{conn: conn, enabled: enabled}
}

func (p *interpolatingConnPool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return p.conn.PrepareContext(ctx, query)
}

func (p *interpolatingConnPool) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	query, args, err := p.interpolate(query, args)
	if err != nil {
		return nil, err
	}
	return p.conn.ExecContext(ctx, query, args...)
}

func (p *interpolatingConnPool) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	query, args, err := p.interpolate(query, args)
	if err != nil {
		return nil, err
	}
	return p.conn.QueryContext(ctx, query, args...)
}

func (p *interpolatingConnPool) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	query, args, err := p.interpolate(query, args)
	if err != nil {
		return p.conn.QueryRowContext(ctx, query, args...)
	}
	return p.conn.QueryRowContext(ctx, query, args...)
}

func (p *interpolatingConnPool) GetDBConn() (*sql.DB, error) {
	if dbConnector, ok := p.conn.(gorm.GetDBConnector); ok {
		return dbConnector.GetDBConn()
	}
	if db, ok := p.conn.(*sql.DB); ok {
		return db, nil
	}
	return nil, gorm.ErrInvalidDB
}

func (p *interpolatingConnPool) interpolate(query string, args []any) (string, []any, error) {
	if !p.enabled || len(args) == 0 {
		return query, args, nil
	}

	interpolated, err := interpolateTDengineParams(query, args)
	if err != nil {
		return "", nil, err
	}
	return interpolated, nil, nil
}

func interpolateTDengineParams(query string, args []any) (string, error) {
	if strings.Count(query, "?") != len(args) {
		return "", driver.ErrSkip
	}

	var buf strings.Builder
	argPos := 0
	for i := 0; i < len(query); i++ {
		q := strings.IndexByte(query[i:], '?')
		if q == -1 {
			buf.WriteString(query[i:])
			break
		}
		buf.WriteString(query[i : i+q])
		i += q

		if err := writeTDengineLiteral(&buf, args[argPos]); err != nil {
			return "", err
		}
		argPos++
	}
	if argPos != len(args) {
		return "", driver.ErrSkip
	}
	return buf.String(), nil
}

func writeTDengineLiteral(buf *strings.Builder, arg any) error {
	if arg == nil {
		buf.WriteString("NULL")
		return nil
	}

	switch v := arg.(type) {
	case int8:
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	case int16:
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	case int32:
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	case int64:
		buf.WriteString(strconv.FormatInt(v, 10))
	case int:
		buf.WriteString(strconv.Itoa(v))
	case uint8:
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	case uint16:
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	case uint32:
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	case uint64:
		buf.WriteString(strconv.FormatUint(v, 10))
	case uint:
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	case float32:
		fmt.Fprintf(buf, "%f", v)
	case float64:
		fmt.Fprintf(buf, "%f", v)
	case bool:
		if v {
			buf.WriteByte('1')
		} else {
			buf.WriteByte('0')
		}
	case time.Time:
		writeQuotedString(buf, v.Format(time.RFC3339Nano))
	case json.RawMessage:
		writeQuotedBytes(buf, v)
	case []byte:
		writeQuotedBytes(buf, v)
	case string:
		writeQuotedString(buf, v)
	default:
		return driver.ErrSkip
	}
	return nil
}

func writeQuotedBytes(buf *strings.Builder, value []byte) {
	buf.WriteByte('\'')
	for _, c := range value {
		if c == '\'' {
			buf.WriteByte('\'')
			buf.WriteByte('\'')
		} else {
			buf.WriteByte(c)
		}
	}
	buf.WriteByte('\'')
}

func writeQuotedString(buf *strings.Builder, value string) {
	buf.WriteByte('\'')
	for i := 0; i < len(value); i++ {
		if value[i] == '\'' {
			buf.WriteByte('\'')
			buf.WriteByte('\'')
		} else {
			buf.WriteByte(value[i])
		}
	}
	buf.WriteByte('\'')
}
