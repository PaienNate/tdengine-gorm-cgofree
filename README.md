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

### Root package compatibility mode

The root package keeps the original behavior and chooses the default driver from the build mode:

* `CGO_ENABLED=1`: `taosSql`, use native TCP DSN, for example `user:password@tcp(tdengine-host:6030)/datacenter?loc=Local`
* `CGO_ENABLED=0`: `taosWS`, use WebSocket DSN, for example `user:password@ws(tdengine-host:6041)/datacenter?loc=Local`

### Driver-isolated packages

If you want to guarantee that WebSocket-only builds never import the native CGO TDengine driver, use the dedicated subpackages instead of the root package:

* `github.com/PaienNate/tdengine-gorm-cgofree/ws`: always registers only `taosWS`, even when `CGO_ENABLED=1`
* `github.com/PaienNate/tdengine-gorm-cgofree/native`: registers only `taosSql`, available only when `CGO_ENABLED=1`

Example:

```go
import (
    tdws "github.com/PaienNate/tdengine-gorm-cgofree/ws"
    "gorm.io/gorm"
)

db, err := gorm.Open(tdws.Open("user:password@ws(tdengine-host:6041)/datacenter?loc=Local"), &gorm.Config{})
```

`taosWS` requires a `ws(...)` or `wss(...)` DSN. Passing `tcp(...)` to `taosWS` fails during dialect initialization with an actionable error before a network connection is attempted.

The dialect enables TDengine string parameter interpolation by default so GORM queries such as `Where("v = ?", "abc")` work with the official driver. To preserve the upstream driver behavior, disable it explicitly:

```go
import (
    tdengine_gorm "github.com/PaienNate/tdengine-gorm-cgofree"
    "gorm.io/gorm"
)

db, err := gorm.Open(&tdengine_gorm.Dialect{
    DSN:               dsn,
    InterpolateParams: tdengine_gorm.WithInterpolateParams(false),
})
```

### Identifier quoting compatibility

By default, the dialect quotes identifiers with backticks. For example, table names, super table names, tag names, and column names are emitted as `` `name` ``.

If you are migrating from `wild-River2016/tdengine_gorm-master` and need the old TDengine behavior during the transition, disable identifier quoting explicitly:

```go
import (
    tdengine_gorm "github.com/PaienNate/tdengine-gorm-cgofree"
    "gorm.io/gorm"
)

db, err := gorm.Open(&tdengine_gorm.Dialect{
    DSN:               dsn,
    InterpolateParams: tdengine_gorm.WithInterpolateParams(false),
    QuoteIdentifiers:  tdengine_gorm.WithQuotedIdentifiers(false),
}, &gorm.Config{})
```

This changes SQL generation from quoted identifiers such as:

```sql
INSERT INTO `TbName` USING `StbName` (`ts`,`value`) VALUES (?,?)
```

to unquoted identifiers:

```sql
INSERT INTO TbName USING StbName (ts,value) VALUES (?,?)
```

The same option is available in the driver-isolated packages:

```go
import (
    tdnative "github.com/PaienNate/tdengine-gorm-cgofree/native"
    tdws "github.com/PaienNate/tdengine-gorm-cgofree/ws"
    "gorm.io/gorm"
)

nativeDB, err := gorm.Open(&tdnative.Dialect{
    DriverName:        tdnative.DefaultDriverName,
    DSN:               nativeDSN,
    InterpolateParams: tdnative.WithInterpolateParams(false),
    QuoteIdentifiers:  tdnative.WithQuotedIdentifiers(false),
}, &gorm.Config{})

wsDB, err := gorm.Open(&tdws.Dialect{
    DriverName:        tdws.DefaultDriverName,
    DSN:               wsDSN,
    InterpolateParams: tdws.WithInterpolateParams(false),
    QuoteIdentifiers:  tdws.WithQuotedIdentifiers(false),
}, &gorm.Config{})
```

Notes:

* `QuoteIdentifiers` only controls whether identifiers are wrapped in backticks.
* Keep the default `true` for new code unless you specifically need legacy SQL naming compatibility.
* Disabling identifier quoting addresses TDengine name case-sensitivity behavior during migration.
* `DryRun` SQL can still differ from the old repository in placeholder rendering for string values; this does not affect identifier case behavior.

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
