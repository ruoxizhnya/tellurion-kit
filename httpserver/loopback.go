package httpserver

// loopback.go —— 回环来源判据（KT-10，裁-18：从 core pkg/httpserver 搬入）。
//
// 只搬判据本身，不搬 gin 中间件：`LoopbackOnly` 的 403 文案与
// 「容器内 RemoteAddr 不是回环」的装配口径属于应用层策略，core 保留那层
// gin 包装（core 二期台账在案）。本仓的依赖白名单允许 gin，但本单按
// 派发只补 core 二期切包所需的导出缺口。
//
// 判据只看 `host:port` 文本，**不看** `X-Forwarded-For` / `ClientIP()`：
// 那些头是客户端可伪造的，在无鉴权模式下信任它们等于把门交给攻击者。

import (
	"net"
	"strings"
)

// IsLoopbackRemoteAddr reports whether a `host:port` remote address is the
// local loopback interface.
//
// Accepted forms: `127.0.0.1:5000`, `[::1]:5000`, `localhost:5000`, and the
// port-less `127.0.0.1` / `::1`. Empty strings, parse failures and non-loopback
// IPs are all false (fail-closed; never guesses).
//
// Migrated verbatim from core pkg/httpserver (KT-10, 裁-18) so core can switch
// to this package without a semantic change.
func IsLoopbackRemoteAddr(remoteAddr string) bool {
	host := strings.TrimSpace(remoteAddr)
	if host == "" {
		return false
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	// IPv6: strip brackets and zone (fe80::1%eth0 is not loopback, but ::1%lo is).
	host = strings.Trim(host, "[]")
	if i := strings.IndexByte(host, '%'); i >= 0 {
		host = host[:i]
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
