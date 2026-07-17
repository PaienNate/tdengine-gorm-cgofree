package ws

import (
	_ "github.com/taosdata/driver-go/v3/taosWS"

	"github.com/PaienNate/tdengine-gorm-cgofree/internal/base"
)

func init() {
	base.DefaultDriverName = DefaultDriverName
}
