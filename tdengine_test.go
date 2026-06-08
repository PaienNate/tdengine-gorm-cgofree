package tdengine_gorm

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/create"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/fill"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/using"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/window"
	"gorm.io/gorm"
)

type integrationConfig struct {
	Host      string
	TCPPort   int
	WSPort    int
	User      string
	Password  string
	Database  string
	NativeDSN string
	WSDSN     string
}

type driverDSN struct {
	name   string
	driver string
	dsn    string
}

func tdengineIntegrationConfig(t *testing.T) integrationConfig {
	t.Helper()

	if os.Getenv("TDENGINE_INTEGRATION") != "1" {
		t.Skip("set TDENGINE_INTEGRATION=1 to run TDengine integration tests")
	}

	host := getenvDefault("TDENGINE_HOST", "127.0.0.1")
	tcpPort := getenvIntDefault(t, "TDENGINE_TCP_PORT", 6030)
	wsPort := getenvIntDefault(t, "TDENGINE_WS_PORT", 6041)
	user := getenvDefault("TDENGINE_USER", "user")
	password := getenvDefault("TDENGINE_PASSWORD", "password")
	dbName := getenvDefault("TDENGINE_TEST_DB", fmt.Sprintf("tdengine_gorm_cgofree_test_%d", os.Getpid()))

	return integrationConfig{
		Host:      host,
		TCPPort:   tcpPort,
		WSPort:    wsPort,
		User:      user,
		Password:  password,
		Database:  dbName,
		NativeDSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, tcpPort),
		WSDSN:     fmt.Sprintf("%s:%s@ws(%s:%d)/", user, password, host, wsPort),
	}
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getenvIntDefault(t *testing.T, key string, fallback int) int {
	t.Helper()

	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		t.Fatalf("invalid %s %q: %v", key, value, err)
	}
	return parsed
}

func dsnWithDB(baseDSN, dbName string) string {
	return strings.TrimRight(baseDSN, "/") + "/" + dbName + "?loc=Local"
}

func dsnWithoutDB(baseDSN string) string {
	return strings.TrimRight(baseDSN, "/") + "/?loc=Local"
}

func Test_Dialect(t *testing.T) {
	cfg := tdengineIntegrationConfig(t)

	testCases := []struct {
		name         string
		query        string
		querySuccess bool
	}{
		{
			name:         "select",
			query:        "SELECT 1",
			querySuccess: true,
		},
		{
			name:         "create db",
			query:        fmt.Sprintf("create database if not exists `%s`", cfg.Database),
			querySuccess: true,
		},
		{
			name:         "create table",
			query:        fmt.Sprintf("create table if not exists `%s`.`test` (`ts` timestamp, `value` double)", cfg.Database),
			querySuccess: true,
		},
		{
			name:         "insert data",
			query:        fmt.Sprintf("insert into `%s`.`test` values (now,12)", cfg.Database),
			querySuccess: true,
		},
		{
			name:         "query data",
			query:        fmt.Sprintf("select * from `%s`.`test` limit 1", cfg.Database),
			querySuccess: true,
		},
		{
			name:         "syntax error",
			query:        fmt.Sprintf("select * rfom `%s`.`test` limit 1", cfg.Database),
			querySuccess: false,
		},
	}

	driverDSNs := make([]driverDSN, 0, 2)
	if DefaultDriverName == "taosSql" {
		driverDSNs = append(driverDSNs, driverDSN{
			name:   "native",
			driver: "taosSql",
			dsn:    dsnWithoutDB(cfg.NativeDSN),
		})
	}
	driverDSNs = append(driverDSNs, driverDSN{
		name:   "websocket",
		driver: "taosWS",
		dsn:    dsnWithoutDB(cfg.WSDSN),
	})

	for _, driverDSN := range driverDSNs {
		t.Run(driverDSN.name, func(t *testing.T) {
			for i, tc := range testCases {
				t.Run(fmt.Sprintf("%d/%s", i, tc.name), func(t *testing.T) {
					db, err := gorm.Open(&Dialect{DriverName: driverDSN.driver, DSN: driverDSN.dsn}, &gorm.Config{})
					if err != nil {
						t.Errorf("Expected Open to succeed; got error: %v", err)
						return
					}
					if db == nil {
						t.Errorf("Expected db to be non-nil.")
						return
					}
					defer closeGormDB(t, db)

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
				})
			}

			db, err := sql.Open(driverDSN.driver, driverDSN.dsn)
			if err != nil {
				t.Fatalf("open cleanup connection: %v", err)
			}
			defer db.Close()
			_, _ = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", cfg.Database))
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
	cfg := tdengineIntegrationConfig(t)
	baseDSN := cfg.NativeDSN
	if DefaultDriverName == "taosWS" {
		baseDSN = cfg.WSDSN
	}

	nativeDb, err := sql.Open(DefaultDriverName, dsnWithoutDB(baseDSN))
	if err != nil {
		t.Errorf("connect db error: %v", err)
		return
	}
	defer nativeDb.Close()

	_, err = nativeDb.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", cfg.Database))
	if err != nil {
		t.Errorf("drop database error: %v", err)
		return
	}
	_, err = nativeDb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.Database))
	if err != nil {
		t.Errorf("create database error: %v", err)
		return
	}
	defer func() {
		_, _ = nativeDb.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", cfg.Database))
	}()

	db, err := gorm.Open(&Dialect{DSN: dsnWithDB(baseDSN, cfg.Database)})
	if err != nil {
		t.Errorf("unexpected error:%v", err)
		return
	}
	defer closeGormDB(t, db)
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
		t.Run("tb_aggregate: create table using stable when insert data", func(t *testing.T) {
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

func Test_StringParameters(t *testing.T) {
	cfg := tdengineIntegrationConfig(t)

	driverDSNs := make([]driverDSN, 0, 2)
	if DefaultDriverName == "taosSql" {
		driverDSNs = append(driverDSNs, driverDSN{
			name:   "native",
			driver: "taosSql",
			dsn:    dsnWithDB(cfg.NativeDSN, cfg.Database),
		})
	}
	driverDSNs = append(driverDSNs, driverDSN{
		name:   "websocket",
		driver: "taosWS",
		dsn:    dsnWithDB(cfg.WSDSN, cfg.Database),
	})

	for _, driverDSN := range driverDSNs {
		t.Run(driverDSN.name, func(t *testing.T) {
			adminDSN := cfg.NativeDSN
			if driverDSN.driver == "taosWS" {
				adminDSN = cfg.WSDSN
			}
			admin, err := sql.Open(driverDSN.driver, dsnWithoutDB(adminDSN))
			if err != nil {
				t.Fatalf("open admin connection: %v", err)
			}
			defer admin.Close()

			_, _ = admin.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", cfg.Database))
			if _, err := admin.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.Database)); err != nil {
				t.Fatalf("create database: %v", err)
			}
			defer func() {
				_, _ = admin.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", cfg.Database))
			}()

			db, err := gorm.Open(&Dialect{
				DriverName:        driverDSN.driver,
				DSN:               driverDSN.dsn,
				InterpolateParams: WithInterpolateParams(true),
			})
			if err != nil {
				t.Fatalf("open interpolating DB: %v", err)
			}
			defer closeGormDB(t, db)

			if err := db.Exec(fmt.Sprintf("create table if not exists `%s`.`param_test` (`ts` timestamp, `v` binary(64))", cfg.Database)).Error; err != nil {
				t.Fatalf("create table: %v", err)
			}

			value := "abc'xyz"
			if err := db.Exec(fmt.Sprintf("insert into `%s`.`param_test` values(now, ?)", cfg.Database), value).Error; err != nil {
				t.Fatalf("insert string parameter with interpolation enabled: %v", err)
			}

			type row struct {
				V string
			}
			var got row
			if err := db.Raw(fmt.Sprintf("select `v` from `%s`.`param_test` where `v` = ?", cfg.Database), value).Scan(&got).Error; err != nil {
				t.Fatalf("query string parameter with interpolation enabled: %v", err)
			}
			if got.V != value {
				t.Fatalf("expected queried value %q, got %q", value, got.V)
			}

			rawDB, err := gorm.Open(&Dialect{
				DriverName:        driverDSN.driver,
				DSN:               driverDSN.dsn,
				InterpolateParams: WithInterpolateParams(false),
			})
			if err != nil {
				t.Fatalf("open raw DB: %v", err)
			}
			defer closeGormDB(t, rawDB)

			err = rawDB.Exec(fmt.Sprintf("insert into `%s`.`param_test` values(now, ?)", cfg.Database), "plain").Error
			if err == nil {
				t.Fatal("expected official driver raw interpolation to fail for string parameters")
			}
		})
	}
}

func closeGormDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql DB handle: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("close sql DB handle: %v", err)
	}
}

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
