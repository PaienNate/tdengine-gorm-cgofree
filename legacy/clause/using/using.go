package using

import current "github.com/PaienNate/tdengine-gorm-cgofree/clause/using"

type Using = current.Using

func SetUsing(sTable string, tags map[string]interface{}) Using {
	return current.SetUsing(sTable, tags)
}

