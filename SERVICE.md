---
service-card: 1
repo: tellurion-kit
kind: library
lifecycle: production
owner: ruoxizhnya
impact: compile-time

requires: []
required-by:
  - tellurion
  - tellurion-agent
  - tellurion-analytics
  - tellurion-data
  - tellurion-orchestration
  - tellurion-risk
  - tellurion-strategy
calls: []
called-by: []

truth:
  readme: README.md
  agents: AGENTS.md
  architecture: docs/ARCHITECTURE.md
  api: docs/api-schema.md
  db: docs/db-schema.md
  dataflow: docs/dataflow.md
  runbook: null
  tasks: docs/TASKS.md

status: evergreen
last-verified: 2026-11-11
verified-by: 2026-11-11 实读 README / AGENTS / SCOPE / api-schema / db-schema / dataflow / docs/ARCHITECTURE.md，并按工作区各仓 go.mod 反向扫描消费方；未写版本号与包数量
---

# 服务卡：tellurion-kit

> **一句话**：依赖图底座上的横切工具层——日志、错误、十进制、统计、HTTP、鉴权、时钟、沙箱、DB 连接；零业务概念。
>
> `library` · `production` · `compile-time` · `ruoxizhnya`

## 1. 这是什么

它是所有服务共享的纯技术能力层：把每仓都要写一遍的日志、错误包装、十进制、统计、HTTP 客户端 / 服务端、鉴权、可观测、时钟、沙箱与 DB 连接收在一处。它只依赖标准库与第三方，不 import 任何本项目的业务仓。

**它拥有**：`logging` / `errors` / `id` / `decimal` / `statistics` / `httpclient` / `httpserver` / `auth`（含 `auth/httpadapter`）/ `authz` / `observability` / `clock` / `sandbox` / `db` 等横切包；`docs/SCOPE.md` 的逐包准入 / 禁入清单。
**它不拥有**：任何业务概念、任何数据库表、任何服务入口。

## 2. 它不做什么

- 不收业务概念：包名 / 类型名 / 函数名 / 字段名不得出现 `backtest` / `strategy` / `order` / `factor` / `market` / `portfolio` 等词表里的词。
- 不收只被单个仓使用的能力——那属于那个仓。
- 不反向依赖 `contracts` / `core` / `agent` 或任何服务仓。
- 不做业务规则：`fees`、`marketimpact` 这类零依赖但属业务的东西留在 core。
- 不提供常驻服务、端口或探针。

## 3. 与谁相连

**运行期没有边**：本仓只是编译期被各仓 import，不监听端口，也不调用任何服务。

### 谁 require 它（编译期）

| 消费方 | 用法 | 它坏了会怎样 |
|---|---|---|
| `tellurion` | 日志、错误、十进制、统计、HTTP、鉴权、DB 连接 | 编译失败；core 无法构建 / 发版 |
| `tellurion-agent` | 日志、DB 迁移机制、鉴权适配 | 编译失败 |
| `tellurion-data` | 日志、HTTP、鉴权、DB 连接 | 编译失败 |
| `tellurion-strategy` | 日志、HTTP、鉴权、DB 连接 | 编译失败 |
| `tellurion-risk` | 日志、HTTP、鉴权、DB 连接 | 编译失败 |
| `tellurion-orchestration` | 日志、HTTP、DB 连接 | 编译失败 |
| `tellurion-analytics` | 日志、HTTP 客户端、DB 连接 | 编译失败 |

> 当前未直接 require：`tellurion-contracts`、`tellurion-broker`、`tellurion-docs`、`tellurion-ui`。消费方清单以各仓 `go.mod` 为准，不复制 pin 版本。

### 它调用（出边）

| 目标 | 协议 | 用来做什么 | 目标挂了会怎样 |
|---|---|---|---|
| 无 | — | — | — |

> `db` 包在运行期会连 PostgreSQL，但那是**消费方进程**的连接，不是本仓自己的出边。

## 4. 运行时身份

- **常驻进程**：无，纯库。
- **存活/就绪语义**：无端口、无探针。
- **构建入口**：`README.md` 快速开始；单独 checkout 必须能跑通 `go build` / `go vet` / `go test` 与两条护栏（`check_scope` / `check_release_matrix`）。
- **端口**：不适用。

## 5. 数据与契约归属

| 维度 | 归属 | 真相源 |
|---|---|---|
| DB schema | 无表、不迁移；`db` 包只提供连接 / DSN / 迁移机制 | `docs/db-schema.md` |
| 跨仓 DTO / 契约 | 无业务 DTO；只提供工具层类型与接口 | `docs/api-schema.md` |
| 对外 HTTP | 无服务；`httpserver` 只是被消费方使用的包装 | `docs/api-schema.md` §1 |

> 逐包「收什么 / 不收什么」见 `docs/SCOPE.md` §3；本卡不复述包清单与导出符号。

## 6. 坏了怎么办

- **影响面**：消费方**编译失败**，不是运行期 503；已构建二进制不受影响。破坏性变更在消费方升级 pin 时才暴露。
- **先确认（按顺序）**：
  1. 消费方 `go build ./...` 是否因工具包签名变化失败；对照 `CHANGELOG.md` 的 breaking 条目。
  2. 本仓 `go run ./tools/check_scope` 与 `go run ./tools/check_release_matrix` 是否红。
  3. 下游 pin 是否落在兼容矩阵声明范围内（版本矩阵见 core `docs/REPO-MAP.md`）。
- **升级路径**：`owner` → 跨仓争议 → core `docs/TASKS.md`
- **runbook**：无（纯库不需要起停手册）；构建与发版见 `README.md`「快速开始」「版本规则」。

## 7. 5 分钟后去哪

| 我想…… | 读这个 |
|---|---|
| 知道每个包收什么、不收什么 | `docs/SCOPE.md` |
| 知道对外工具形状 | `docs/api-schema.md` |
| 知道为什么本仓没有业务类型 | `docs/db-schema.md`、`docs/SCOPE.md` §1 |
| 知道版本规则与兼容矩阵 | `README.md` / `CHANGELOG.md` |
| 当前在做什么 | `docs/TASKS.md` |
| 架构与决策 | `docs/ARCHITECTURE.md` / `docs/adr/` |
