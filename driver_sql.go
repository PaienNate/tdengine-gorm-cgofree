//go:build cgo

package tdengine_gorm

import (
	"github.com/PaienNate/tdengine-gorm-cgofree/internal/base"
	_ "github.com/taosdata/driver-go/v3/taosSql"
	_ "github.com/taosdata/driver-go/v3/taosWS"
)

func init() {
	DefaultDriverName = "taosSql"
	base.DefaultDriverName = DefaultDriverName
}
