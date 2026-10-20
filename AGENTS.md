# AGENTS — tellurion-kit

## 这是什么仓

横切工具，**零业务概念**：logging · db（仅 connection/migrate/DSN）· errors · decimal · statistics · httpclient · httpserver · authz · observability · clock · sandbox。

状态：**待建（S2 / QT-07）**——当前为占位仓，**尚无实现**。抽取动作见 QT-05（S1）与 QT-07（S2）。

## 构建 / 测试

```bash
go build ./...
go vet ./...
go test ./...
```

本仓必须**单独 checkout 就能跑通**这三条（R1/R2）。若某一步需要别的仓的源码目录，
那是**契约画错了**，不是正常状态——停下来先修边界，不要加 `replace ../xxx`。

## 边界（硬约束）

允许：
- 被 contracts / core / agent / 各服务 `require`
- 放与业务无关的通用设施

禁止：
- 出现业务词（backtest / strategy / order / factor / market / portfolio …）
- 依赖 contracts（方向相反：contracts 依赖 kit 的 decimal/errors）
- 持有任何表结构或 schema 迁移
- `replace` 指向仓外路径；`go.work*` 提交入库。
- 引入未在本文件声明的第三方依赖（先改本文件，再改代码）。

## 依赖的契约版本

| 依赖 | 版本 | 说明 |
|---|---|---|
| `github.com/ruoxizhnya/tellurion-contracts` | 待定（S2 后 `v0.1.0`） | 契约变更只走 tag，不走源码路径 |
| `github.com/ruoxizhnya/tellurion-kit` | 待定（S2 后 `v0.1.0`） | 同上 |

## 该做什么

**唯一台账是 `tellurion` 仓的 `docs/TASKS.md`。** 本仓没有独立 TASKS，
跨仓只留指针、不镜像正文（避免双写漂移）。

开工前先读：
1. `tellurion:docs/design/qt-decoupling/repo-inventory.md` —— 仓名与 module 规则
2. `tellurion:docs/TASKS.md` —— 找本仓的单号；没有单号就别开工
3. `tellurion:docs/design/qt-decoupling/p0-round2-review.md` §4 —— wave 顺序
