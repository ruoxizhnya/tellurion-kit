# 本仓边界与依赖方向 — tellurion-kit

## 定位

横切工具层：logging · errors · id · decimal · statistics · httpclient · httpserver · auth/authz · observability · clock · sandbox · **db**（仅 connection/migrate/DSN）。**零业务概念**。

状态：**S1 落代码（QT-05）／S2 打 tag（QT-07）——当前为占位仓，尚无实现**

## 依赖方向（编译器强制）

```
github.com/ruoxizhnya/tellurion-kit
        ↑
contracts / core / agent / 各服务
```

- 箭头 = import 方向。**不允许反向**，也不允许成环 —— 反向 import 应当**编译失败**，这是判据不是口号。
- 跨仓只认版本号：`github.com/ruoxizhnya/tellurion-kit` 以需要的版本出现在下游 `go.mod`。

## 边界表

| | 内容 |
|---|---|
| 允许 | 被任意仓 `require`（它是依赖图的底座）；纯技术能力：日志、错误包装、十进制、统计、HTTP 客户端/服务端、鉴权、可观测、时钟、沙箱、DB 连接 |
| 禁止 | 任何业务概念：backtest / strategy / order / factor / market / portfolio …；任何表结构（`kit/db` 只做 connection / migrate / DSN 解析）；反向依赖 contracts 或任何应用仓 |

## 关键约束

1. **单仓可构建**：`go build ./...` / `go vet ./...` / `go test ./...` 不依赖任何兄弟目录。
2. **无 `replace`**：指向仓外的 `replace` 一律不接受（CI 会查）。
3. **无 `go.work` 入库**：本地联调可以建，提交即错。
4. **业务词零命中**：`kit` 内出现 backtest/strategy/order/factor/market 等词即违规（QT-05 验收项）。

## 未决 / 待定

- `db` 包的边界（只连接/迁移/DSN，零表结构）随 QT-05 落地时再核对一次。
