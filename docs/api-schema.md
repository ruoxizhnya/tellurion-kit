---
status: active
last-verified: 2026-10-20
verified-by: 功能边界固化（三份契约：db / api / dataflow）
---

# API Schema — tellurion-kit

> **契约地位**：本仓的"API"是 **Go 包级 API**（导出函数 / 类型）。消费方式是 `require` + 版本号。
> 本文件是**本仓对外接口的唯一真相源**。代码与它不一致时，先判断谁错：契约错 → 改本文 + 立 ADR；实现错 → 改代码。
> 跨仓消费方（core / agent / ui / 其他服务）按本文件编码，**不读本仓实现**。

## 1. 接口面

| 包 | 提供什么 | 典型用途 | 实现状态 |
|---|---|---|---|
| `logging` | 结构化日志、级别、字段约定 | 全仓统一日志形状 | ⬜ 未落地（仅 doc.go 占位） |
| `errors` | 包装 / 分类 / 堆栈 | 与 contracts 的三类哨兵配合 | ⬜ 未落地（仅 doc.go 占位） |
| `id` | 有序 ID 生成 | run_id / job_id | ⬜ 未落地（仅 doc.go 占位） |
| `decimal` | 定点小数 | 价格与金额（禁止 float 做钱） | ⬜ 未落地（仅 doc.go 占位） |
| `statistics` | 常用统计量 | 指标计算 | ⬜ 未落地（仅 doc.go 占位） |
| `httpclient` | 带超时 / 重试 / 熔断的客户端 | 服务间调用 | ⬜ 未落地（仅 doc.go 占位） |
| `httpserver` | 统一的错误响应与中间件 | 各服务 HTTP 面 | ⬜ 未落地（仅 doc.go 占位） |
| `auth` / `authz` | 认证与鉴权原语 | API 保护 | ⬜ 未落地（仅 doc.go 占位） |
| `observability` | 指标 / 追踪 / 健康检查 | `/metrics`、`/health` | ⬜ 未落地（仅 doc.go 占位） |
| `clock` | 可注入时钟 | 让"纯函数"可测试 | ⬜ 未落地（仅 doc.go 占位） |
| `sandbox` | staticcheck / runner | 策略沙箱 | ⬜ 未落地（仅 doc.go 占位） |
| `db` | 连接 / DSN / 迁移执行（**零表结构**） | 各仓自己的 store | ⬜ 未落地（仅 doc.go 占位） |

## 2. 错误码与 fail-closed 语义

- 统一用 `errors` 包包装，保留 `%w` 链（可 `errors.Is` / `As`）。
- 对外错误必须能映射到 HTTP 码（`httpserver` 提供映射表）；**不许把内部错误字符串直接吐给客户端**。
- 重试只由 `httpclient` 决定，调用方不自己写重试循环。

## 3. 幂等性与重试

不适用（工具层不定义业务幂等；幂等键属于各服务的 API 契约）。

## 4. 版本与兼容

- 加函数 / 加包 → MINOR
- 改签名 / 改语义 → MAJOR（0.x 期间 bump MINOR 并写 CHANGELOG）
- 删导出符号 → MAJOR（1.0 后先 deprecated 一个 MINOR 周期）

## 5. 不提供什么（边界）

- 不提供业务概念（哪怕是"通用"的 Strategy 接口）
- 不提供表结构 / 迁移 SQL
- 不提供 DI 容器 / 应用框架（那是 core 的装配层，QT-50 后外移）
- 不提供测试 mock 框架
