# tellurion-kit

> module: `github.com/ruoxizhnya/tellurion-kit` ｜ 状态：**S1 落代码（QT-05）／S2 打 tag（QT-07）——当前为占位仓，尚无实现**

横切工具层：logging · errors · id · decimal · statistics · httpclient · httpserver · auth/authz · observability · clock · sandbox · **db**（仅 connection/migrate/DSN）。**零业务概念**。

## 目录

| 路径 | 内容 |
|---|---|
| `docs/` | 本仓文档：台账 / 边界 / 决策 / 归档 —— 从这里进：`docs/README.md` |
| `docs/SCOPE.md` | **13 包逐包准入 / 禁入清单** + 业务词表 + 依赖方向 + 实测偏差登记（KT-01） |
| `tools/check_scope/` | 包边界护栏：业务词 / 越界 import / 未登记包 → 红（`go run ./tools/check_scope`） |
| `AGENTS.md` | 本仓入口（构建、边界、禁止、依赖版本） |
| `doc.go` | 占位包（让空仓 `go build ./...` 从一开始就是绿的） |

## 快速开始

```bash
go build ./...
go vet ./...
go run ./tools/check_scope   # 包边界 ↔ docs/SCOPE.md 逐包清单（CI 同款，KT-01 护栏）
go test ./...
```

本仓**单独 checkout 就能跑通**上面四条。需要别的仓源码才过 → 边界错了，不是加 `replace` 能解决的问题。

## 边界

**可以**：被任意仓 `require`（它是依赖图的底座）；纯技术能力：日志、错误包装、十进制、统计、HTTP 客户端/服务端、鉴权、可观测、时钟、沙箱、DB 连接

**不可以**：任何业务概念：backtest / strategy / order / factor / market / portfolio …；任何表结构（`kit/db` 只做 connection / migrate / DSN 解析）；反向依赖 contracts 或任何应用仓

> 上面两行是**一句话口径**；逐包「收什么 / 不收什么」的完整清单在 [`docs/SCOPE.md`](docs/SCOPE.md)，它是 CI 核对的登记表，不是补充说明。

## 文档

- 文档导航：`docs/README.md`
- **包边界清单**：`docs/SCOPE.md`（13 包逐包准入/禁入 + 业务词表；`go run ./tools/check_scope` 核对）
- 本仓台账：`docs/TASKS.md`（唯一；跨仓任务只留指针）
- 边界与依赖方向：`docs/ARCHITECTURE.md`
- 决策记录：`docs/adr/`

## 许可

MIT（见 `LICENSE`）。
