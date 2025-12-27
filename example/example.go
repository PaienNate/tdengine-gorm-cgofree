package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	tdengine_gorm "github.com/PaienNate/tdengine-gorm-cgofree"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/create"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/fill"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/using"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/window"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const testDsn = "root:taosdata@/tcp(127.0.0.1:6030)"

type Data struct {
	TS    time.Time
	Value int64
}

func main() {
	// create database
	createDatabase()
	// connect to the database
	db := connect()
	// create a sTable
	// CREATE STABLE IF NOT EXISTS `stb_1` (`ts` TIMESTAMP,`value` BIGINT) TAGS(`tbn` BINARY(64))
	createSTable(db)

	// CREATE TABLE IF NOT EXISTS `tb_1` USING `stb_1`(`tbn`) TAGS ('tb_1')
	createTableUsingStable(db)
	now := time.Now()
	randValue := rand.Int63()

	// INSERT INTO `tb_1` (`ts`,`value`) VALUES ('2021-08-11 09:43:00.041',0.604660)
	insertData(db, "tb_1", now, randValue)
	t1 := now.Add(time.Second)
	randValue2 := rand.Int63()

	// INSERT INTO tb_2 USING stb_1('tbn') TAGS('tb_2') (`ts`,`value`) VALUES ('2021-08-11 09:43:01.041',0.940509)
	automaticTableCreationWhenInsertingData(db, "tb_2", t1, randValue2)
	// SELECT * FROM `tb_1` WHERE `ts` = '2021-08-11 09:43:00.041'
	tb1Data := queryData(db, "tb_1", now)
	if tb1Data.Value != randValue {
		log.Fatalf("expect value %v got %v", randValue, tb1Data.Value)
	}
	// SELECT * FROM `tb_2` WHERE `ts` = '2021-08-11 09:43:01.041'
	tb2Data := queryData(db, "tb_2", t1)
	if tb2Data.Value != randValue2 {
		log.Fatalf("expect value %v got %v", randValue, tb2Data.Value)
	}
	// SELECT * FROM `stb_1` WHERE `ts` = '2021-08-11 09:43:00.041'
	stbData := queryData(db, "stb_1", now)
	if stbData.Value != randValue {
		log.Fatalf("expect value %v got %v", randValue, stbData.Value)
	}
	t2 := now.Add(time.Second * 2)
	t3 := now.Add(time.Second * 3)
	t4 := now.Add(time.Second * 4)
	v1 := 11
	v2 := 12
	v3 := 13

	// INSERT INTO tb_aggregate USING stb_1('tbn') TAGS('tb_aggregate') (ts,value) VALUES ('2021-08-11 09:43:01.041',11),('2021-08-11 09:43:02.041',12),('2021-08-11 09:43:03.041',13)
	automaticTableCreationWhenInsertingMultiData(db, "tb_aggregate", []map[string]any{
		{
			"ts":    t1,
			"value": v1,
		}, {
			"ts":    t2,
			"value": v2,
		}, {
			"ts":    t3,
			"value": v3,
		},
	})
	// aggregate query
	// SELECT avg(`value`) as v FROM tb_aggregate WHERE ts >= '2021-08-11 09:43:01.041' and ts <= '2021-08-11 09:43:03.041'
	resultAvg := aggregateQuery(db, "tb_aggregate", "avg(`value`) as v", t1, t3, nil)
	expectAvg := []map[string]any{
		{
			"v": int64(12),
		},
	}
	if !resultMapEqual(expectAvg, resultAvg) {
		log.Fatalf("expect %v got %v", expectAvg, resultAvg)
	}
	windowD, err := window.NewDuration(time.Second)
	if err != nil {
		log.Fatal(err)
	}
	// SELECT max(`value`) as v FROM tb_aggregate WHERE ts >= '2021-08-11 09:43:01.041' and ts <= '2021-08-11 09:43:04.041' INTERVAL(1000000u) FILL (NULL)
	resultWindowMax := aggregateQuery(db, "tb_aggregate", "max(`value`) as v", t1, t4, []clause.Expression{
		window.SetInterval(*windowD),
		fill.Fill{Type: fill.FillNull},
	})
	expectWindowMax := []map[string]any{
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
	if !resultMapEqual(expectWindowMax, resultWindowMax) {
		log.Fatalf("expect %v got %v", expectWindowMax, resultWindowMax)
	}
}

func createDatabase() {
	dsnWithoutDB := testDsn + "/?loc=Local"
	nativeDB, err := sql.Open(tdengine_gorm.DefaultDriverName, dsnWithoutDB)
	if err != nil {
		log.Fatalf("connect db error:%v", err)
		return
	}
	_, err = nativeDB.Exec("create database if not exists gorm_test")
	if err != nil {
		log.Fatalf("create database error %v", err)
		return
	}
	_ = nativeDB.Close()
}

func connect() *gorm.DB {
	dsn := testDsn + "/gorm_test?loc=Local"
	db, err := gorm.Open(&tdengine_gorm.Dialect{DSN: dsn})
	if err != nil {
		log.Fatalf("unexpected error:%v", err)
	}
	db = db.Debug()
	return db
}

func createSTable(db *gorm.DB) {
	// create stable
	stable := create.NewSTableBuilder("stb_1").
		IfNotExists().
		Columns(
			[]*create.Column{
				{
					Name: "ts",
					Type: create.Timestamp,
				},
				{
					Name: "value",
					Type: create.Double,
				},
			}...,
		).
		TagColumns(&create.Column{
			Name:   "tbn",
			Type:   create.Binary,
			Length: 64,
		}).
		Build()

	err := db.Table("stb_1").Clauses(create.NewCreateTable(stable)).Create(map[string]any{}).Error
	if err != nil {
		log.Fatalf("create sTable error %v", err)
	}
}

func createTableUsingStable(db *gorm.DB) {
	// create table using sTable
	table := create.NewCTableBuilder("tb_1").
		IfNotExists().
		BuildWithSTable("stb_1", map[string]any{"tbn": "tb_1"})
	err := db.Table("tb_1").Clauses(create.NewCreateTable(table)).Create(map[string]any{}).Error
	if err != nil {
		log.Fatalf("create table error %v", err)
	}
}

func insertData(db *gorm.DB, tableName string, ts time.Time, value any) {
	// insert data
	err := db.Table(tableName).Create(map[string]any{
		"ts":    ts,
		"value": value,
	}).Error
	if err != nil {
		log.Fatalf("insert data error %v", err)
	}
}

func automaticTableCreationWhenInsertingData(db *gorm.DB, tableName string, ts time.Time, value any) {
	// automatic table creation when inserting data
	err := db.Table(tableName).Clauses(using.SetUsing("stb_1", map[string]any{
		"tbn": tableName,
	})).Create(map[string]any{
		"ts":    ts,
		"value": value,
	}).Error
	if err != nil {
		log.Fatalf("create table when insert data error %v", err)
	}
}

func queryData(db *gorm.DB, tableName string, ts time.Time) *Data {
	var d Data
	err := db.Table(tableName).Where("`ts` = ?", ts).Find(&d).Error
	if err != nil {
		log.Fatalf("find data error %v", err)
	}
	return &d
}

func automaticTableCreationWhenInsertingMultiData(db *gorm.DB, tableName string, data []map[string]any) {
	// automatic table creation when inserting data
	err := db.Table(tableName).Clauses(using.SetUsing("stb_1", map[string]any{
		"tbn": tableName,
	})).Create(data).Error
	if err != nil {
		log.Fatalf("create table when insert data error %v", err)
	}
}

func aggregateQuery(db *gorm.DB, tableName string, query string, start, end time.Time, conds []clause.Expression) []map[string]any {
	var result []map[string]any
	err := db.Table(tableName).Select(query).Where("`ts` >= ? and `ts` <= ?", start, end).Clauses(conds...).Find(&result).Error
	if err != nil {
		log.Fatalf("aggregate query error %v", err)
	}
	return result
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
