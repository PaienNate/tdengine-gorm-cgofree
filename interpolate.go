package tdengine_gorm

import (
	"github.com/PaienNate/tdengine-gorm-cgofree/internal/base"
	"gorm.io/gorm"
)

func interpolateParamsEnabled(dialect Dialect) bool {
	return base.InterpolateParamsEnabled(dialect.baseDialect())
}

func newInterpolatingConnPool(conn gorm.ConnPool, enabled bool) gorm.ConnPool {
	return base.NewInterpolatingConnPool(conn, enabled)
}

func interpolateTDengineParams(query string, args []any) (string, error) {
	return base.InterpolateTDengineParams(query, args)
}
