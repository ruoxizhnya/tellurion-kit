# tellurion-kit（夹具：尚未发版的基线，对应本仓今天的真实形状）

## 版本规则

<!-- rules:begin -->
| 位 | 何时 bump | 例子 |
|---|---|---|
| `MAJOR` | 删符号、改签名、改语义 | 删掉 `Decimal` 的一个方法 |
| `MINOR` | 新增向后兼容的能力 | 新增一个纯数学函数 |
| `PATCH` | 文档 / 注释 / 无行为修正 | 措辞订正 |
<!-- rules:end -->

## 兼容矩阵

<!-- matrix:begin -->
| kit 版本 | 状态 | 支持的 contracts | 支持的 core | 支持的 agent |
|---|---|---|---|---|
| `main`（未发版） | 开发基线，不做兼容承诺 | 尚无 require | 伪版本 pin `main` | 伪版本 pin `main` |
<!-- matrix:end -->
