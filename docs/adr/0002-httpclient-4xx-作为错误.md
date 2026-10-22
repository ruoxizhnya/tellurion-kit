# 0002 · httpclient 把 4xx 作为错误返回

- **状态**：已裁决（来源：KT-04，2026-10-23）
- **场景**：`httpclient.Client.Do` 原先只把 5xx / 429 当失败（重试），其余状态码一律 `err = nil` 返回。于是调用方写 `resp, err := c.Do(...); if err != nil` 这条最自然的门时，一个 400 / 404 会被当成成功继续往下走 —— 实测（KT-04 ②）400 返回 `err=<nil>` + `StatusCode:400`，调用方必须额外记住「还要看 StatusCode」。同一次取证还确认带 body 的请求在第二次尝试发出空 body。
- **决策**：
  1. 4xx（429 除外，429 走重试）返回 `nil` + `*errors.AppError`，错误码按状态映射：400 / 422 / 未映射 4xx → `INVALID_INPUT`，401 / 403 → `PERMISSION`，404 → `NOT_FOUND`，408 → `TIMEOUT`，409 → `CONFLICT`。上游响应体截断 200 字后写进错误 `message`，`op` 为 `"<METHOD> <Path>"`。
  2. 请求体在每次尝试前重新构造 reader，重试完整重发。
  3. 由此新增 `kit` 内部边 `httpclient → errors`（`errors` 是 stdlib-only 的叶子包，无回边风险）。
- **理由**：本仓与整个项目的失败处理是 fail-closed 口径 —— `err != nil` 必须是调用方唯一的「没成功」门，不能让每个使用方各自记得再看 `StatusCode`。`Response` 的形状不改（不新增字段），错误里带上游 body 已能保留诊断信息；422 与未知 4xx 都归 `INVALID_INPUT` 是因为它们同样表示「请求本身不成立」。
- **代价 / 后果**：
  - 这是**导出面语义 breaking**：`v0.1.0` 之前 pin `main` 的使用方若依赖「4xx 时 `err == nil`」会红；本仓零 tag，按 README 的 0.x 规则记 breaking、不 bump 版本号（CHANGELOG `[Unreleased]` 在册）。
  - 4xx 时不再返回 `Response`，调用方拿不到原始 body 字节；错误 `message` 内的截断引用是有意保留的诊断面。
  - 重试仍**不看方法**（KT-04 ④ 选择「文档钉死」而不是加开关）：有副作用的调用须显式 `New(baseURL, timeout, 0)` 关重试。以后要加「只对幂等方法重试」的开关，属新导出面，另立单。
- **替代关系**：无。同批的 `decimal` 裁决见 0001 之后的 KT-03 记录（不涉及包间边）。
