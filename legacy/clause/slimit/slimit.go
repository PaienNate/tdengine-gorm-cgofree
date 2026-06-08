package slimit

import current "github.com/PaienNate/tdengine-gorm-cgofree/clause/slimit"

type SLimit = current.SLimit

func SetSLimit(limit, offset int) SLimit {
	return current.SetSLimit(limit, offset)
}

