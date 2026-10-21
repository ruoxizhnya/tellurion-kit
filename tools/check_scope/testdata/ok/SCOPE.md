# SCOPE（夹具：十三包齐、内容中性的最小清单）

<!-- scope:begin -->
| 包 | 准入 | 禁入 | kit 内允许依赖 |
|---|---|---|---|
| `logging` | 日志形状与脱敏 | 业务埋点 | — |
| `errors` | 包装与分类 | 业务错误码语义 | — |
| `id` | 与域无关的 ID | 交易对象 ID | — |
| `decimal` | 定点运算 | 费率模型 | — |
| `statistics` | 纯数学 | 金融口径 | — |
| `httpclient` | 客户端形状 | endpoint 表 | logging |
| `httpserver` | gin 骨架 | 路由清单 | errors |
| `auth` | token 机制 | 表结构 SQL | — |
| `authz` | 授权裁决机制 | 路由到角色的映射表 | auth |
| `observability` | registry 与埋点 | 业务指标名 | — |
| `clock` | 时间原语 | 回测语义 | — |
| `sandbox` | 静态检查与进程隔离 | WASM 策略宿主 | — |
| `db` | 连接 / DSN / 迁移机制 | 任何表结构 | — |
<!-- scope:end -->
