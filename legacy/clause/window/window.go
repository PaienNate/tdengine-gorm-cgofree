package window

import (
	"time"

	current "github.com/PaienNate/tdengine-gorm-cgofree/clause/window"
)

type UnitType = current.UnitType
type Duration = current.Duration
type Window = current.Window

const (
	Microsecond = current.Microsecond
	Millisecond = current.Millisecond
	Second      = current.Second
	Minute      = current.Minute
	Hour        = current.Hour
	Day         = current.Day
	Week        = current.Week
	Month       = current.Month
	Year        = current.Year
)

func NewDurationFromTimeDuration(duration time.Duration) (*Duration, error) {
	return current.NewDurationFromTimeDuration(duration)
}

func ParseDuration(durationString string) (*Duration, error) {
	return current.ParseDuration(durationString)
}

func SetSessionWindow(tsColumn string, duration Duration) Window {
	return current.SetSessionWindow(tsColumn, duration)
}

func SetStateWindow(column string) Window {
	return current.SetStateWindow(column)
}

func SetInterval(duration Duration) Window {
	return current.SetInterval(duration)
}

