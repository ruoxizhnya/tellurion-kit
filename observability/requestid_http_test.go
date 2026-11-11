package observability

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// echoHandler 把从 context 取回的 request_id 写进响应体，方便断言
// 「handler 看到的 token」与「响应头」是不是同一个。
func echoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, RequestIDFromContext(r.Context()))
	})
}

// 无入站 token ⇒ 生成一个，且 handler 看到的与响应头一致。
func TestRequestIDHandler_GeneratesWhenAbsent(t *testing.T) {
	rec := httptest.NewRecorder()
	RequestIDHandler(echoHandler()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", nil))

	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.NotEmpty(t, body, "缺失时应生成一个 request_id")
	assert.Equal(t, body, rec.Header().Get(RequestIDHeader), "响应头必须回镜像")
}

// 入站已有 token ⇒ 原样透传，不重新生成。
func TestRequestIDHandler_PreservesInboundToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(RequestIDHeader, "upstream-token")

	rec := httptest.NewRecorder()
	RequestIDHandler(echoHandler()).ServeHTTP(rec, req)

	assert.Equal(t, "upstream-token", rec.Body.String())
	assert.Equal(t, "upstream-token", rec.Header().Get(RequestIDHeader))
}

// 出站传播：handler 用**同一个 context** 发下游请求，HTTPTransport 把
// token 写进下游的 X-Request-ID，下游看到的就是入站那个 token。
func TestRequestIDHandler_PropagatesToOutboundRequest(t *testing.T) {
	var downstream string
	base := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		downstream = req.Header.Get(RequestIDHeader)
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
	})
	client := &http.Client{Transport: &HTTPTransport{Base: base}}

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "http://downstream.test/x", nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		_ = resp.Body.Close()
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(RequestIDHeader, "trace-9")
	RequestIDHandler(inner).ServeHTTP(httptest.NewRecorder(), req)

	assert.Equal(t, "trace-9", downstream, "出站请求必须用同一个 token")
}

// 并发下每个请求各自拿自己的 token：既不能串号，也不能互相覆盖。
func TestRequestIDHandler_ConcurrentNoCrosstalk(t *testing.T) {
	const n = 64
	handler := RequestIDHandler(echoHandler())

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			token := fmt.Sprintf("rid-%03d", i)

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			req.Header.Set(RequestIDHeader, token)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if got := rec.Body.String(); got != token {
				t.Errorf("第 %d 个请求串号：context 里是 %q，期望 %q", i, got, token)
			}
			if got := rec.Header().Get(RequestIDHeader); got != token {
				t.Errorf("第 %d 个请求响应头串号：%q，期望 %q", i, got, token)
			}
		}(i)
	}
	wg.Wait()
}
