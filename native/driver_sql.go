//go:build cgo

package native

import (
	_ "github.com/taosdata/driver-go/v3/taosSql"

	"github.com/PaienNate/tdengine-gorm-cgofree/internal/base"
)

func init() {
	base.DefaultDriverName = DefaultDriverName
}
