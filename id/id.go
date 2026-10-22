// Package id 是 UUIDv7 的集中生成器（可按时间排序，适合做数据库主键）。
//
// 为什么是 v7 而不是 v4：前 48 bit 是 unix 毫秒时间戳，主键按时间有序，
// 高写入量下 B-tree 插入落在最右页，操作员也能按主键还原时间顺序。
//
// 本包只放**与业务无关**的 ID 生成原语。交易对象 ID（OrderID 这类）属于
// 交易域的领域对象，不进本仓（SCOPE.md §5 偏差 #1 的裁决）。
package id

import (
	"github.com/google/uuid"
)

// JobID 是异步任务的 ID（任务队列 / 同步作业 / 构建任务）。
// 与其它 ID 共用 v7 原语：DB 列形状相同（VARCHAR(64) PRIMARY KEY），
// 排序收益也相同。
func JobID() string {
	return newV7()
}

// SubscriptionID 用于内存 pub/sub 注册表的键。它不落库，v7 只是防御性
// 默认值 —— 无论哪个版本，碰撞概率都低到可以忽略。
func SubscriptionID() string {
	return newV7()
}

// newV7 把 uuid.NewV7() 收敛到一处：将来换 v8 或 Snowflake 只改这一行。
//
// 返回规范 36 字符形式（xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx）。
func newV7() string {
	u, err := uuid.NewV7()
	if err != nil {
		// uuid.NewV7 只会在系统时间出错（时钟倒退到 1970 之前）时失败。
		// 这是进程级故障：回落到 v4 会悄悄丢掉排序保证，所以必须响亮地
		// 失败（panic），让配置错误的 NTP 当场暴露。
		panic("id: UUID v7 generation failed (system clock?): " + err.Error())
	}
	return u.String()
}
