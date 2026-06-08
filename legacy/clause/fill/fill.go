package fill

import current "github.com/PaienNate/tdengine-gorm-cgofree/clause/fill"

type Type = current.FillType
type Fill = current.Fill

const (
	FillNone   = current.FillNone
	FillValue  = current.FillValue
	FillPrev   = current.FillPrev
	FillNull   = current.FillNull
	FillLinear = current.FillLinear
	FillNext   = current.FillNext
)

func SetFill(fillType Type) Fill {
	return current.SetFill(fillType)
}

