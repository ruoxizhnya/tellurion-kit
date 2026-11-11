// Package observability — stdlib net/http 版的 request_id 中间件。
//
// 为什么单开一份而不是改 gin 版：6 个服务仓用的是 stdlib `net/http` mux
// （`http.ServeMux` / `http.Handler` 链），拿不到 `gin.HandlerFunc`。
// 语义与 [RequestIDMiddleware] 逐条对齐，只是形状换成 stdlib 的
// `func(http.Handler) http.Handler`：
//
//   - 入站有 `X-Request-ID` 则透传（上游 token 优先），无则生成新的 UUIDv4
//   - token 写进 request 的 context.Context（[RequestIDKey]），用
//     [RequestIDFromContext] 取回
//   - 回镜像到响应头 `X-Request-ID`
//   - 出站传播不用另学：下游请求带上**同一个 context**，交给
//     [HTTPTransport]（或 [WithRequestIDHeader]）即可
//
// 幂等性同 gin 版：链上套两次时，第二次读到的是第一次生成的 token，
// 结果不变。
package observability

import (
	"net/http"

	"github.com/google/uuid"
)

// RequestIDHandler 用一般的 http.Handler 链上 request_id 语义。
//
// 典型装配：
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/metrics", metricsHandler)
//	srv := &http.Server{Handler: observability.RequestIDHandler(mux)}
//
// 取 ID：observability.RequestIDFromContext(r.Context())。
// 往下游传播：用同一个 r.Context() 发请求，Transport 换成
// &observability.HTTPTransport{...}。
func RequestIDHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 上游已经带了 token 就必须透传：链路追踪的父子关系靠它串起来，
		// 重生成会把一次调用切成两段互不相干的 trace。
		rid := r.Header.Get(RequestIDHeader)
		if rid == "" {
			rid = uuid.NewString()
		}
		// 先写响应头再进 handler：下游 handler 里 panic 时，外层 recovery
		// 仍能凭这个头把日志和响应关联起来。
		w.Header().Set(RequestIDHeader, rid)
		next.ServeHTTP(w, r.WithContext(withRequestID(r.Context(), rid)))
	})
}
