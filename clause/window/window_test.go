package window_test

import (
	"testing"
	"time"

	"github.com/PaienNate/tdengine-gorm-cgofree/clause/tests"
	"github.com/PaienNate/tdengine-gorm-cgofree/clause/window"

	"gorm.io/gorm/clause"
)

func Test_SetInterval(t *testing.T) {
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			Name: "",
			Clauses: []clause.Interface{
				clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
				clause.From{Tables: []clause.Table{{Name: "t_1"}}},
				window.SetInterval(window.Duration{Value: 10, Unit: window.Minute}),
			},
			Result: []string{"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(10m)"},
			Vars:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}

func Test_SetStateWindow(t *testing.T) {
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			Name: "",
			Clauses: []clause.Interface{
				clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
				clause.From{Tables: []clause.Table{{Name: "t_1"}}},
				window.SetStateWindow("state"),
			},
			Result: []string{"SELECT avg(`t_1`.`value`) FROM `t_1` STATE_WINDOW(`state`)"},
			Vars:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}

func Test_SetSessionWindow(t *testing.T) {
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			Name: "",
			Clauses: []clause.Interface{
				clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
				clause.From{Tables: []clause.Table{{Name: "t_1"}}},
				window.SetSessionWindow("ts", window.Duration{
					Value: 10,
					Unit:  window.Minute,
				}),
			},
			Result: []string{"SELECT avg(`t_1`.`value`) FROM `t_1` SESSION(`ts`,10m)"},
			Vars:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}

func Test_SetOffset(t *testing.T) {
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			Clauses: []clause.Interface{
				clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
				clause.From{Tables: []clause.Table{{Name: "t_1"}}},
				window.SetInterval(window.Duration{Value: 10, Unit: window.Minute}).SetOffset(window.Duration{
					Value: 5,
					Unit:  window.Minute,
				}),
			},
			Result: []string{"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(10m,5m)"},
			Vars:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}

func Test_SetSliding(t *testing.T) {
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			Name: "",
			Clauses: []clause.Interface{
				clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
				clause.From{Tables: []clause.Table{{Name: "t_1"}}},
				window.SetInterval(window.Duration{Value: 10, Unit: window.Minute}).SetOffset(window.Duration{
					Value: 5,
					Unit:  window.Minute,
				}).SetSliding(window.Duration{
					Value: 2,
					Unit:  window.Minute,
				}),
			},
			Result: []string{"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(10m,5m) SLIDING(2m)"},
			Vars:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}

func Test_NewDuration(t *testing.T) {
	duration5Min, err := window.NewDuration(time.Minute * 5)
	if err != nil {
		t.Errorf("NewDurationFromTimeDuration error : %s", err.Error())
		return
	}
	_, err = window.NewDuration(-time.Second)
	if err == nil {
		t.Errorf("Need error")
		return
	}
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			Name: "",
			Clauses: []clause.Interface{
				clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
				clause.From{Tables: []clause.Table{{Name: "t_1"}}},
				window.SetInterval(*duration5Min),
			},
			Result: []string{"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(300000000u)"},
			Vars:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}

// ParseDuration
func Test_ParseDuration(t *testing.T) {
	duration5Min, err := window.ParseDuration("5m")
	if err != nil {
		t.Errorf("ParseDuration error : %s", err.Error())
		return
	}
	_, err = window.ParseDuration("1")
	if err == nil {
		t.Errorf("need error")
		return
	}
	_, err = window.ParseDuration("1K")
	if err == nil {
		t.Errorf("need error")
		return
	}
	_, err = window.ParseDuration("mm")
	if err == nil {
		t.Errorf("need error")
		return
	}
	var testCases = []struct {
		Name    string
		Clauses []clause.Interface
		Result  []string
		Vars    [][][]any
	}{
		{
			Clauses: []clause.Interface{
				clause.Select{Columns: []clause.Column{{Name: "avg(`t_1`.`value`)", Raw: true}}},
				clause.From{Tables: []clause.Table{{Name: "t_1"}}},
				window.SetInterval(*duration5Min),
			},
			Result: []string{"SELECT avg(`t_1`.`value`) FROM `t_1` INTERVAL(5m)"},
			Vars:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tests.CheckBuildClauses(t, tc.Clauses, tc.Result, tc.Vars)
		})
	}
}
