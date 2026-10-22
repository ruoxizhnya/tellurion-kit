// Package kit 是 tellurion-kit 的根包，只有文档：本仓的代码一律住子包。
//
// 留这个文件的两个理由（实现落地后仍然成立）：
//  1. 让本仓可被 `go list ./...` 识别，从而能接入 R1/R2 自检；
//  2. 给「module 根上有什么」一个明确答案：什么都没有，别往根包堆代码。
//
// 十三个子包的边界（收什么 / 不收什么 / 允许的包间依赖）在 docs/SCOPE.md §3，
// 对外形状在 docs/api-schema.md §1；`go run ./tools/check_scope` 按前者核对，越界即红。
package kit
