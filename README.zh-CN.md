# TDengine Gorm Dialect

[English README](./README.md)

## 说明

当前方言不支持 migrate、update、delete 和 transaction。

## 驱动与 DSN

### 根包兼容模式

根包会根据构建模式自动选择默认驱动：

* `CGO_ENABLED=1`：默认使用 `taosSql`，DSN 形如 `user:password@tcp(tdengine-host:6030)/datacenter?loc=Local`
* `CGO_ENABLED=0`：默认使用 `taosWS`，DSN 形如 `user:password@ws(tdengine-host:6041)/datacenter?loc=Local`

### 驱动隔离子包

如果你希望显式区分 WebSocket 与 Native 驱动，可以使用子包：

* `github.com/PaienNate/tdengine-gorm-cgofree/ws`：始终只注册 `taosWS`
* `github.com/PaienNate/tdengine-gorm-cgofree/native`：始终只注册 `taosSql`，仅在 `CGO_ENABLED=1` 时可用

示例：

```go
import (
    tdws "github.com/PaienNate/tdengine-gorm-cgofree/ws"
    "gorm.io/gorm"
)

db, err := gorm.Open(tdws.Open("user:password@ws(tdengine-host:6041)/datacenter?loc=Local"), &gorm.Config{})
```

`taosWS` 只接受 `ws(...)` 或 `wss(...)` DSN。如果把 `tcp(...)` DSN 传给 `taosWS`，会在方言初始化阶段直接报错。

## 参数插值

方言默认开启 TDengine 字符串参数插值，这样 `Where("v = ?", "abc")` 之类的查询可以直接工作。

如果你希望保留上游 driver 的原始行为，可以显式关闭：

```go
import (
    tdengine_gorm "github.com/PaienNate/tdengine-gorm-cgofree"
    "gorm.io/gorm"
)

db, err := gorm.Open(&tdengine_gorm.Dialect{
    DSN:               dsn,
    InterpolateParams: tdengine_gorm.WithInterpolateParams(false),
}, &gorm.Config{})
```

## legacy 标识符兼容

当前版本默认会给标识符加反引号，例如表名、超级表名、tag 名、列名都会输出成 `` `name` ``。

如果你正在从 `wild-River2016/tdengine_gorm-master` 迁移，并且希望保留旧仓库“不加反引号”的行为，直接使用 `OpenLegacy`：

```go
import (
    tdengine_gorm "github.com/PaienNate/tdengine-gorm-cgofree"
    "gorm.io/gorm"
)

db, err := gorm.Open(tdengine_gorm.OpenLegacy(dsn), &gorm.Config{})
```

它等价于：

```go
db, err := gorm.Open(&tdengine_gorm.Dialect{
    DSN:              dsn,
    QuoteIdentifiers: tdengine_gorm.WithQuotedIdentifiers(false),
}, &gorm.Config{})
```

也就是把：

```sql
INSERT INTO `TbName` USING `StbName` (`ts`,`value`) VALUES (?,?)
```

变成：

```sql
INSERT INTO TbName USING StbName (ts,value) VALUES (?,?)
```

`native` / `ws` 子包也提供同名 helper：

```go
import (
    tdnative "github.com/PaienNate/tdengine-gorm-cgofree/native"
    tdws "github.com/PaienNate/tdengine-gorm-cgofree/ws"
    "gorm.io/gorm"
)

nativeDB, err := gorm.Open(tdnative.OpenLegacy(nativeDSN), &gorm.Config{})
wsDB, err := gorm.Open(tdws.OpenLegacy(wsDSN), &gorm.Config{})
```

说明：

* `QuoteIdentifiers` 只控制标识符是否加反引号
* 新代码建议保留默认值，也就是继续使用 `Open(...)`
* 迁移期如果主要目标是避免 TDengine 名称大小写变化，优先使用 `OpenLegacy(...)`
* `DryRun` SQL 在字符串占位符展示上仍可能和旧仓库不同，但这不影响标识符大小写行为

## 集成测试

集成测试默认关闭。运行前请准备独立测试库：

```sh
TDENGINE_INTEGRATION=1 \
TDENGINE_HOST=127.0.0.1 \
TDENGINE_TCP_PORT=6030 \
TDENGINE_WS_PORT=6041 \
TDENGINE_USER=your-user \
TDENGINE_PASSWORD=your-password \
go test ./... -count=1
```

可以通过 `TDENGINE_TEST_DB` 覆盖临时测试库名。不要把集成测试指向生产库。

## 已扩展子句

* `CREATE TABLE`
* `FILL`
* `SLIMIT`
* `USING`
* `WINDOW`

## 示例

示例代码见 [example](./example/example.go)。
