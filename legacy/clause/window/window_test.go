package window_test

import (
	"testing"
	"time"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/tests"
	legacywindow "github.com/PaienNate/tdengine-gorm-cgofree/legacy/clause/window"
	"gorm.io/gorm/clause"
)

func TestLegacyWindowClauses(t *testing.T) {
	tests.CheckBuildClauses(t, []clause.Interface{
		clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
		clause.From{Tables: []clause.Table{{Name: "t_1"}}},
		legacywindow.SetInterval(legacywindow.Duration{Value: 10, Unit: legacywindow.Minute}).
			SetOffset(legacywindow.Duration{Value: 5, Unit: legacywindow.Minute}).
			SetSliding(legacywindow.Duration{Value: 2, Unit: legacywindow.Minute}),
	}, []string{
		"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(10m,5m) SLIDING(2m)",
	}, nil)

	tests.CheckBuildClauses(t, []clause.Interface{
		clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
		clause.From{Tables: []clause.Table{{Name: "t_1"}}},
		legacywindow.SetStateWindow("state"),
	}, []string{
		"SELECT avg(`t_1`.`value`) FROM `t_1` STATE_WINDOW(`state`)",
	}, nil)

	tests.CheckBuildClauses(t, []clause.Interface{
		clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
		clause.From{Tables: []clause.Table{{Name: "t_1"}}},
		legacywindow.SetSessionWindow("ts", legacywindow.Duration{Value: 10, Unit: legacywindow.Minute}),
	}, []string{
		"SELECT avg(`t_1`.`value`) FROM `t_1` SESSION(`ts`,10m)",
	}, nil)
}

func TestLegacyWindowDurationHelpers(t *testing.T) {
	duration5Min, err := legacywindow.NewDurationFromTimeDuration(time.Minute * 5)
	if err != nil {
		t.Fatalf("NewDurationFromTimeDuration error: %v", err)
	}
	if _, err := legacywindow.NewDurationFromTimeDuration(-time.Second); err == nil {
		t.Fatal("expected negative duration error")
	}

	tests.CheckBuildClauses(t, []clause.Interface{
		clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
		clause.From{Tables: []clause.Table{{Name: "t_1"}}},
		legacywindow.SetInterval(*duration5Min),
	}, []string{
		"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(300000000u)",
	}, nil)

	parsed, err := legacywindow.ParseDuration("5m")
	if err != nil {
		t.Fatalf("ParseDuration error: %v", err)
	}
	if _, err := legacywindow.ParseDuration("1"); err == nil {
		t.Fatal("expected short duration error")
	}
	if _, err := legacywindow.ParseDuration("1K"); err == nil {
		t.Fatal("expected invalid unit error")
	}
	if _, err := legacywindow.ParseDuration("mm"); err == nil {
		t.Fatal("expected invalid value error")
	}

	tests.CheckBuildClauses(t, []clause.Interface{
		clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
		clause.From{Tables: []clause.Table{{Name: "t_1"}}},
		legacywindow.SetInterval(*parsed),
	}, []string{
		"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(5m)",
	}, nil)
}

