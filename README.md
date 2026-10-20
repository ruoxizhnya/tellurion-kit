# tellurion-kit

> module: `github.com/ruoxizhnya/tellurion-kit` ｜ 状态：**待建（S2 / QT-07）**——当前为占位仓，**尚无实现**。抽取动作见 QT-05（S1）与 QT-07（S2）。

横切工具，**零业务概念**：logging · db（仅 connection/migrate/DSN）· errors · decimal · statistics · httpclient · httpserver · authz · observability · clock · sandbox。

`kit/db` 是 K 裁决的落地（ADR-D29）：只提供连接、迁移与 DSN 解析，**不持有任何表结构**——各服务用自己的 schema。

## 边界

**可以**
- 被 contracts / core / agent / 各服务 `require`
- 放与业务无关的通用设施

**不可以**
- 出现业务词（backtest / strategy / order / factor / market / portfolio …）
- 依赖 contracts（方向相反：contracts 依赖 kit 的 decimal/errors）
- 持有任何表结构或 schema 迁移

## 依赖方向

```
github.com/ruoxizhnya/tellurion-kit
```

- 契约与工具一律以 **tag** 依赖（`github.com/ruoxizhnya/tellurion-contracts` / `-kit`），
  **禁止 `replace` 指向仓外路径**（单仓不可构建）。
- `go.work` 一律不入库。

## 施工入口

本仓所有任务登记在 **`tellurion` 仓的 `docs/TASKS.md`**（单一台账，避免双写）。
在台账出现本仓的单号之前，**不要在此写实现**。

## 权威文档

| 主题 | 位置 |
|---|---|
| 仓清单与命名 | `tellurion:docs/design/qt-decoupling/repo-inventory.md` |
| 分仓拓扑与阶梯 | `tellurion:docs/design/qt-decoupling/repo-topology-and-naming.md` |
| P0 修复基线 | `tellurion:docs/design/qt-decoupling/p0-round2-review.md` |
