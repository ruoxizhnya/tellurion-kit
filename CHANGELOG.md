# CHANGELOG — tellurion-kit

> 本仓的版本账本。**每个发版一节**，标题格式固定 `## [vX.Y.Z] - YYYY-MM-DD`（`tools/check_release_matrix` 按它和 README 兼容矩阵、`git tag` 三方对拍）。
> 未发版的改动记在 `## [Unreleased]` 下，发版时把那一节的内容搬进新版本节。
> 规则口径见 [`README.md`](README.md) 的「版本规则」（源头：`tellurion:docs/TASKS.md` QT-19 ③）。
> breaking 的条目必须写成「破坏性：<旧 → 新>」，光写「改了签名」不算留账。
> ⚠️ 本仓在最底层：一条 breaking 会同时打到 contracts / core / agent，所以「谁被迫改代码」要写在条目里。

## [Unreleased]

尚未发版（本仓零符号、零 tag）。首个版本将是 `v0.1.0`，由 `tellurion:docs/TASKS.md` QT-07 ③-S2-a 在 S2 打；在此之前使用方只能 pin `main`，**无兼容承诺**。

### 变更（PATCH 级：文档与仓内自查工具，不改导出面）

- **KT-02**（2026-10-22）：README 补「版本规则表 + 兼容矩阵 + 发版流程 + 拉取前提（`GOPRIVATE`）」；新建本 `CHANGELOG.md`；新建护栏 `tools/check_release_matrix`（矩阵 ↔ CHANGELOG ↔ tag 三件套不一致即红）并接进 CI。导出面零变化 —— 本仓今天没有符号。
- **KT-01**（2026-10-22）：新建 `docs/SCOPE.md`（13 包逐包准入/禁入 + 业务词表 + 依赖方向 + 九条实测偏差登记）与 `tools/check_scope`（业务词 / 越界 import / 未登记包即红），接进 CI 步骤「包边界核对」。导出面零变化。
