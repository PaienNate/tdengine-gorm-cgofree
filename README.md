# TDengine Gorm Dialect

[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/thinkgos/tdengine-gorm?tab=doc)
[![codecov](https://codecov.io/gh/thinkgos/tdengine-gorm/graph/badge.svg?token=aHu5wq1m6i)](https://codecov.io/gh/thinkgos/tdengine-gorm)
[![Tests](https://github.com/thinkgos/tdengine-gorm/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/thinkgos/tdengine-gorm/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thinkgos/tdengine-gorm)](https://goreportcard.com/report/github.com/thinkgos/tdengine-gorm)
[![Licence](https://img.shields.io/github/license/thinkgos/tdengine-gorm)](https://raw.githubusercontent.com/thinkgos/tdengine-gorm/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/thinkgos/tdengine-gorm)](https://github.com/thinkgos/tdengine-gorm/tags)


## Instructions

Not support migrate, update, deletion and transaction

## Driver and DSN

The default driver depends on the build mode:

* `CGO_ENABLED=1`: `taosSql`, use native TCP DSN, for example `user:password@tcp(tdengine-host:6030)/datacenter?loc=Local`
* `CGO_ENABLED=0`: `taosWS`, use WebSocket DSN, for example `user:password@ws(tdengine-host:6041)/datacenter?loc=Local`

`taosWS` requires a `ws(...)` or `wss(...)` DSN. Passing `tcp(...)` to `taosWS` fails during dialect initialization with an actionable error before a network connection is attempted.

The dialect enables TDengine string parameter interpolation by default so GORM queries such as `Where("v = ?", "abc")` work with the official driver. To preserve the upstream driver behavior, disable it explicitly:

```go
db, err := gorm.Open(&tdengine_gorm.Dialect{
    DSN:               dsn,
    InterpolateParams: tdengine_gorm.WithInterpolateParams(false),
})
```

Integration tests are disabled by default. Run them against an isolated test database:

```sh
TDENGINE_INTEGRATION=1 \
TDENGINE_HOST=127.0.0.1 \
TDENGINE_TCP_PORT=6030 \
TDENGINE_WS_PORT=6041 \
TDENGINE_USER=your-user \
TDENGINE_PASSWORD=your-password \
go test ./... -count=1
```

Set `TDENGINE_TEST_DB` to override the temporary test database name. Do not point integration tests at production databases.

Add clauses

* "CREATE TABLE"
* "FILL"
* "SLIMIT"
* "USING"
* "WINDOW"

## EXAMPLE

Check example code [example](./example/example.go)
