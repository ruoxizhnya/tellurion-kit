---
status: active
last-verified: 2026-11-12
verified-by: 2026-11-12 库存三单全部收口——KIT-DEP-SECURITY-BUMP（govulncheck EXIT=0 / 破坏验证真做真还原 / 发版 `v0.4.0`）、KIT-SURFACE-FREEZE（导出面基线 272 行 / 两向破坏验证 / CI 接线）、KIT-REQUESTID-STDLIB（stdlib 中间件 4 条单测 / 删透传即红 / 发版 `v0.5.0`）；同日台账卫生批次将三条原文逐字移入归档并补检索索引行；此前 2026-11-11 全面审查 W7 台账立单（A5 授权）——新立 3 条六字段单，证据坐标引 ODR-069/072（`tellurion-docs:docs/archive/odr/`）；`KIT-DEP-SECURITY-BUMP` 是 `tellurion:X-DEP-BUMP-CONSUMERS` 的前置
---

# TASKS — tellurion-kit 未完成项

> **归档说明（2026-11-11）**：本文件已瘦身为「只装未做的事」。2026-11-10 及以前已收口（✅）的任务原文、施工记录、验收输出与破坏验证，已逐字移入 [`docs/archive/tasks-archive-2026-11.md`](archive/tasks-archive-2026-11.md)；该文件顶部有「检索索引」表（ID / 标题 / 收口日期 / commit / 证据链接）。检索：`grep -n "<任务ID或commit>" docs/archive/tasks-archive-2026-11.md`。
>
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `⏸ 挂起（外部依赖或已裁决暂不实施；**不计入开放项**，重开须另立单并满足该条「重开条件」）` / `✅ 完成（附日期，随即移入归档）`

## 一、未做 / 挂起 / 外部依赖

**当前无未做 / 挂起 / 外部依赖项**（2026-11-12 台账卫生批次：库存三单 `KIT-DEP-SECURITY-BUMP` / `KIT-SURFACE-FREEZE` / `KIT-REQUESTID-STDLIB` 已全部 ✅，原文与收口记录逐字移入归档「2026-11-12 收口（台账卫生批次）」节并补检索索引行）。已收口单的正文与施工记录历史部分在归档文件里；新增工作一律另立单。

## 二、跨仓指针（不在本仓，勿双写）

| 事项 | 去哪 |
|---|---|
| 契约形状与版本纪律 | `tellurion-contracts:docs/TASKS.md` |
| core 侧装配与鉴权消费 | `tellurion:docs/TASKS.md` |

## 维护本文件

- **新增单**：写全六字段；六字段里必须显式引用**冻结契约的路径与节号**（契约仓的包路径、设计文档的节号），不许写"参考相关实现"。
- **派工前先证伪前提**：读码核对工单写的前提；不符就停下回报 + 给拆法，不硬做、不扩范围。
- **跨仓引用写 `repo:path`**（如 `tellurion:docs/TASKS.md`），不做可点击相对链接。
- **顺序纪律**：同一波次内撞同一批文件的单按单号串行，一单一 commit。
