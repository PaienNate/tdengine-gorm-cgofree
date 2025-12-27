//go:build !cgo

package tdengine_gorm

import _ "github.com/taosdata/driver-go/v3/taosWS"

func init() {
    DefaultDriverName = "taosWS"
}
