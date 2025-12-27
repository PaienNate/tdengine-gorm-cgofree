module github.com/PaienNate/tdengine-gorm-cgofree

go 1.22

// 暂时保持不变
replace github.com/taosdata/driver-go/v3 v3.7.0 => github.com/thinkgos/driver-go/v3 v3.6.1-0.20250514070202-16d60a2a4d79

require (
	github.com/taosdata/driver-go/v3 v3.7.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/text v0.20.0 // indirect
)
