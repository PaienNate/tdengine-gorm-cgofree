package tdengine_gorm

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/create"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/fill"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/using"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/window"
	"gorm.io/gorm"
)

const (
	testDb       = "gorm_test"
	dsn          = "root:taosdata@tcp(localhost:6030)/"
	dsnWithoutDb = dsn + "?loc=Local"
	dsnWithDb    = dsn + testDb + "?loc=Local"
)

func Test_Dialect(t *testing.T) {
	testCases := []struct {
		name         string
		dialect      *Dialect
		openSuccess  bool
		query        string
		querySuccess bool
	}{
		{
			name: "Default driver",
			dialect: &Dialect{
				DSN: dsn,
			},
			openSuccess:  true,
			query:        "SELECT 1",
			querySuccess: true,
		},
		{
			name: "create db",
			dialect: &Dialect{
				DriverName: DefaultDriverName,
				DSN:        dsn,
			},
			openSuccess:  true,
			query:        "create database if not exists `gorm_test`",
			querySuccess: true,
		},
		{
			name: "create table",
			dialect: &Dialect{
				DriverName: DefaultDriverName,
				DSN:        dsn,
			},
			openSuccess:  true,
			query:        "create table if not exists `gorm_test`.`test` (`ts` timestamp, `value` double)",
			querySuccess: true,
		},
		{
			name: "insert data",
			dialect: &Dialect{
				DriverName: DefaultDriverName,
				DSN:        dsn,
			},
			openSuccess:  true,
			query:        "insert into `gorm_test`.`test` values (now,12)",
			querySuccess: true,
		},
		{
			name: "query data",
			dialect: &Dialect{
				DriverName: DefaultDriverName,
				DSN:        dsn,
			},
			openSuccess:  true,
			query:        "select * from `gorm_test`.`test` limit 1",
			querySuccess: true,
		},
		{
			name: "syntax error",
			dialect: &Dialect{
				DriverName: DefaultDriverName,
				DSN:        dsn,
			},
			openSuccess:  true,
			query:        "select * rfom `gorm_test`.`test` limit 1",
			querySuccess: false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d/%s", i, tc.name), func(t *testing.T) {
			db, err := gorm.Open(tc.dialect, &gorm.Config{})
			if !tc.openSuccess {
				if err == nil {
					t.Errorf("Expected Open to fail.")
				}
				return
			}
			if err != nil {
				t.Errorf("Expected Open to succeed; got error: %v", err)
				return
			}
			if db == nil {
				t.Errorf("Expected db to be non-nil.")
				return
			}
			if tc.query != "" {
				err = db.Exec(tc.query).Error
				if !tc.querySuccess {
					if err == nil {
						t.Errorf("Expected query to fail.")
					}
					return
				}

				if err != nil {
					t.Errorf("Expected query to succeed; got error: %v", err)
				}
			}
		})
	}
}

const TestStb1 = "stb_1"

type TestTb1 struct {
	TS    time.Time
	Value int64
}

func (*TestTb1) TableName() string {
	return "tb_1"
}

type TestTb2 struct {
	TS    time.Time
	Value int64
}

func (*TestTb2) TableName() string {
	return "tb_2"
}

type TestTbAggregate struct {
	TS    time.Time
	Value int64
}

func (*TestTbAggregate) TableName() string {
	return "tb_aggregate"
}

func Test_Clause(t *testing.T) {
	nativeDb, err := sql.Open(DefaultDriverName, dsnWithoutDb)
	if err != nil {
		t.Errorf("connect db error: %v", err)
		return
	}
	_, err = nativeDb.Exec("CREATE DATABASE IF NOT EXISTS `gorm_test`")
	nativeDb.Close()
	if err != nil {
		t.Errorf("create database error: %v", err)
		return
	}
	db, err := gorm.Open(&Dialect{DSN: dsnWithDb})
	if err != nil {
		t.Errorf("unexpected error:%v", err)
		return
	}
	db = db.Debug()

	t.Run("create stable", func(t *testing.T) {
		stable := create.NewSTableBuilder(TestStb1).
			IfNotExists().
			Columns(
				&create.Column{Type: create.Timestamp, Name: "ts"},
				&create.Column{Type: create.BigInt, Name: "value"},
			).
			TagColumns(
				&create.Column{Type: create.Binary, Name: "tbn", Length: 64},
			).
			Build()

		err = db.Table(stable.TableName()).
			Clauses(create.NewCreateTable(stable)).
			Create(map[string]any{}).Error
		if err != nil {
			t.Errorf("create stable error %v", err)
			return
		}
	})

	now := time.Now()
	vx1 := rand.Int63()

	t.Run("tb_1", func(t *testing.T) {
		t.Run("tb_1: create table using stable", func(t *testing.T) {
			table := create.NewCTableBuilder("tb_1").
				IfNotExists().
				BuildWithSTable(TestStb1, map[string]any{"tbn": "tb_1"})
			err = db.Clauses(create.NewCreateTable(table)).
				Create(&TestTb1{}).Error
			if err != nil {
				t.Errorf("create table error, %v", err)
				return
			}
		})
		t.Run("tb_1: insert data", func(t *testing.T) {
			err = db.Create(&TestTb1{
				TS:    now,
				Value: vx1,
			}).Error
			if err != nil {
				t.Errorf("tb_1: insert data error, %v", err)
				return
			}
		})
		t.Run("tb_1: find data", func(t *testing.T) {
			var got TestTb1

			err = db.Model(&TestTb1{}).Where("`ts` = ?", now).Find(&got).Error
			if err != nil {
				t.Errorf("find data error, %v", err)
				return
			}
			if got.Value != vx1 {
				t.Errorf("expect value: %v, got: %v", vx1, got.Value)
				return
			}
		})
		t.Run("tb_1: find via stable", func(t *testing.T) {
			var got TestTb1

			err = db.Table(TestStb1).Where("`ts` = ?", now).Find(&got).Error
			if err != nil {
				t.Errorf("find data by stable error %v", err)
				return
			}
			if got.Value != vx1 {
				t.Errorf("expect value %v got %v", vx1, got.Value)
				return
			}
		})
	})

	t1 := now.Add(time.Second)
	t2 := now.Add(time.Second * 2)
	t3 := now.Add(time.Second * 3)
	v1 := 11
	v2 := 12
	v3 := 13
	vx2 := rand.Int63()
	t.Run("tb_2", func(t *testing.T) {
		t.Run("tb_2: create table using stable when insert data", func(t *testing.T) {
			err = db.Clauses(using.SetUsing(TestStb1, map[string]any{"tbn": "tb_2"})).
				Create(&TestTb2{
					TS:    t1,
					Value: vx2,
				}).Error
			if err != nil {
				t.Errorf("tb_2: create table when insert data error, %v", err)
				return
			}
		})
		t.Run("tb_2: find data", func(t *testing.T) {
			var got TestTb2

			err = db.Model(&TestTb2{}).Where("`ts` = ?", t1).Find(&got).Error
			if err != nil {
				t.Errorf("find data error %v", err)
				return
			}
			if got.Value != vx2 {
				t.Errorf("expect value: %v, got: %v", vx2, got.Value)
				return
			}
		})
	})

	t.Run("tb_aggregate", func(t *testing.T) {
		t.Run("tb_aggregate: create table using stable when insert dat", func(t *testing.T) {
			err = db.Clauses(using.SetUsing(TestStb1, map[string]any{"tbn": "tb_aggregate"})).
				Create([]*TestTbAggregate{
					{t1, int64(v1)},
					{t2, int64(v2)},
					{t3, int64(v3)},
				}).Error
			if err != nil {
				t.Errorf("tb_aggregate: create table using stable when insert data error, %v", err)
				return
			}
		})
		t.Run("tb_aggregate: query avg", func(t *testing.T) {
			var result []map[string]any

			err = db.Table("tb_aggregate").
				Select("avg(`value`) as v").
				Where("`ts` >= ?", now.Add(time.Second)).
				Where("`ts` <= ?", now.Add(time.Second*3)).
				Find(&result).Error
			if err != nil {
				t.Errorf("aggregate query error %v", err)
				return
			}
			expectR1 := []map[string]any{
				{
					"v": float64(12),
				},
			}
			if !resultMapEqual(expectR1, result) {
				t.Errorf("expect %v got %v", expectR1, result)
				return
			}
		})

		t.Run("tb_aggregate: query time window", func(t *testing.T) {
			var result2 []map[string]any
			wd, err := window.NewDuration(time.Second)
			if err != nil {
				t.Fatal(err)
			}
			err = db.Table("tb_aggregate").
				Select("`ts`, max(`value`) as v").
				Where("`ts` >= ?", now.Add(time.Second)).
				Where("`ts` <= ?", now.Add(time.Second*4)).
				Clauses(
					window.SetInterval(*wd),
					fill.Fill{Type: fill.FillNull},
				).
				Find(&result2).Error
			if err != nil {
				t.Errorf("aggregate query error %v", err)
				return
			}
			expectR2 := []map[string]any{
				{
					"ts": now.Add(time.Second),
					"v":  int64(11),
				},
				{
					"ts": now.Add(time.Second * 2),
					"v":  int64(12),
				},
				{
					"ts": now.Add(time.Second * 3),
					"v":  int64(13),
				},
				{
					"ts": now.Add(time.Second * 4),
					"v":  nil,
				},
			}
			if !resultMapEqual(result2, expectR2) {
				t.Errorf("aggregate query expect %v got %v", result2, expectR2)
				return
			}
		})
	})

	t.Run("stb_1: delete data", func(t *testing.T) {
		err = db.Table(TestStb1).Where("`ts` <= ?", now).Delete(map[string]any{}).Error
		if err != nil {
			t.Errorf("stb_1: delete data, %v", err)
			return
		}
	})
} // nolint: gocyclo

func resultMapEqual(m1, m2 []map[string]any) bool {
	if len(m1) != len(m2) {
		return false
	}
	for i := range m1 {
		if len(m1[i]) != len(m2[i]) {
			return false
		}
	}
	for i, m := range m1 {
		for s, v := range m {
			_, ok := m2[i][s].(time.Time)
			if ok {
				continue
			}
			if m2[i][s] != v {
				return false
			}
		}
	}
	return true
}
