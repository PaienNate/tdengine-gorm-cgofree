package tdengine_gorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestInterpolateParamsFeatureDefault(t *testing.T) {
	if !interpolateParamsEnabled(Dialect{}) {
		t.Fatal("expected parameter interpolation to be enabled by default")
	}
	if !interpolateParamsEnabled(Dialect{InterpolateParams: WithInterpolateParams(true)}) {
		t.Fatal("expected explicit true to enable parameter interpolation")
	}
	if interpolateParamsEnabled(Dialect{InterpolateParams: WithInterpolateParams(false)}) {
		t.Fatal("expected explicit false to disable parameter interpolation")
	}
}

func TestInterpolateTDengineParamsFormatsLiterals(t *testing.T) {
	ts := time.Date(2026, 6, 8, 20, 30, 1, 123456789, time.UTC)
	got, err := interpolateTDengineParams(
		"insert into t values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		[]any{
			"plain",
			"quote ' value",
			[]byte("bytes ' value"),
			json.RawMessage(`{"k":"v'1"}`),
			ts,
			nil,
			true,
			false,
			int64(-7),
			uint(8),
			float64(1.5),
		},
	)
	if err != nil {
		t.Fatalf("unexpected interpolation error: %v", err)
	}

	want := "insert into t values('plain', 'quote '' value', 'bytes '' value', '{\"k\":\"v''1\"}', '2026-06-08T20:30:01.123456789Z', NULL, 1, 0, -7, 8, 1.500000)"
	if got != want {
		t.Fatalf("unexpected interpolation:\nwant: %s\n got: %s", want, got)
	}
}

func TestInterpolateTDengineParamsReturnsErrSkipOnMismatch(t *testing.T) {
	_, err := interpolateTDengineParams("select ?", []any{"a", "b"})
	if !errors.Is(err, driver.ErrSkip) {
		t.Fatalf("expected driver.ErrSkip, got %v", err)
	}
}

func TestInterpolatingConnPoolConsumesArgsWhenEnabled(t *testing.T) {
	base := &recordingConnPool{}
	pool := newInterpolatingConnPool(base, true)

	_, err := pool.ExecContext(context.Background(), "insert into t values(?)", "abc")
	if err != nil {
		t.Fatalf("unexpected exec error: %v", err)
	}

	if base.lastQuery != "insert into t values('abc')" {
		t.Fatalf("expected interpolated query, got %q", base.lastQuery)
	}
	if len(base.lastArgs) != 0 {
		t.Fatalf("expected wrapper to consume args, got %#v", base.lastArgs)
	}
}

func TestInterpolatingConnPoolPassesArgsWhenDisabled(t *testing.T) {
	base := &recordingConnPool{}
	pool := newInterpolatingConnPool(base, false)

	_, err := pool.QueryContext(context.Background(), "select * from t where v = ?", "abc")
	if err != nil {
		t.Fatalf("unexpected query error: %v", err)
	}

	if base.lastQuery != "select * from t where v = ?" {
		t.Fatalf("expected original query, got %q", base.lastQuery)
	}
	if len(base.lastArgs) != 1 || base.lastArgs[0] != "abc" {
		t.Fatalf("expected original args, got %#v", base.lastArgs)
	}
}

func TestInterpolatingConnPoolQueryRowConsumesArgsWhenEnabled(t *testing.T) {
	base := &recordingConnPool{}
	pool := newInterpolatingConnPool(base, true)

	_ = pool.QueryRowContext(context.Background(), "select * from t where v = ?", "abc")

	if base.lastQuery != "select * from t where v = 'abc'" {
		t.Fatalf("expected interpolated query, got %q", base.lastQuery)
	}
	if len(base.lastArgs) != 0 {
		t.Fatalf("expected wrapper to consume args, got %#v", base.lastArgs)
	}
}

type recordingConnPool struct {
	lastQuery string
	lastArgs  []any
}

func (p *recordingConnPool) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, errors.New("prepare not implemented in test double")
}

func (p *recordingConnPool) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	p.record(query, args)
	return driver.RowsAffected(1), nil
}

func (p *recordingConnPool) QueryContext(_ context.Context, query string, args ...any) (*sql.Rows, error) {
	p.record(query, args)
	return nil, nil
}

func (p *recordingConnPool) QueryRowContext(_ context.Context, query string, args ...any) *sql.Row {
	p.record(query, args)
	return nil
}

func (p *recordingConnPool) record(query string, args []any) {
	p.lastQuery = query
	p.lastArgs = append([]any(nil), args...)
}

func TestInterpolateTDengineParamsRejectsUnsupportedValues(t *testing.T) {
	_, err := interpolateTDengineParams("select ?", []any{strings.Builder{}})
	if !errors.Is(err, driver.ErrSkip) {
		t.Fatalf("expected driver.ErrSkip, got %v", err)
	}
}
