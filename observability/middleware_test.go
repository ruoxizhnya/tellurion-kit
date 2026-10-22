package observability

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestIDMiddleware_GeneratesAndPropagates(t *testing.T) {
	r := gin.New()
	r.Use(RequestIDMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%s|%s", c.GetString(ContextKey), RequestIDFromContext(c.Request.Context()))
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.NotEmpty(t, body)
	assert.Contains(t, body, "|")
	parts := splitOnce(body, '|')
	assert.NotEmpty(t, parts[0], "缺失时应生成一个 request_id")
	assert.Equal(t, parts[0], parts[1], "gin context 与 request context 必须是同一个 token")
	assert.Equal(t, parts[0], w.Header().Get(RequestIDHeader), "响应头必须回镜像")
}

func TestRequestIDMiddleware_PreservesInboundToken(t *testing.T) {
	r := gin.New()
	r.Use(RequestIDMiddleware())
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, c.GetString(ContextKey)) })

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(RequestIDHeader, "upstream-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, "upstream-token", w.Body.String())
	assert.Equal(t, "upstream-token", w.Header().Get(RequestIDHeader))
}

func TestHTTPTransport_PropagatesFromContext(t *testing.T) {
	var gotHeader string
	base := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		gotHeader = req.Header.Get(RequestIDHeader)
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
	})

	m := NewMetrics()
	m.Register()
	transport := &HTTPTransport{Base: base, Metrics: m, Service: "data"}

	ctx := withRequestID(context.Background(), "rid-42")
	req := httptest.NewRequest(http.MethodGet, "http://example.test/x", nil).WithContext(ctx)
	_, err := transport.RoundTrip(req)
	require.NoError(t, err)

	assert.Equal(t, "rid-42", gotHeader)
	// transport 的 metrics 开关真的记了一次出站调用。
	assert.InDelta(t, 1, testutil.ToFloat64(m.HTTPRequests.WithLabelValues("data", "200")), 1e-9)
}

func TestHTTPTransport_DoesNotClobberExplicitHeader(t *testing.T) {
	var gotHeader string
	base := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		gotHeader = req.Header.Get(RequestIDHeader)
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
	})
	transport := &HTTPTransport{Base: base}

	ctx := withRequestID(context.Background(), "ctx-token")
	req := httptest.NewRequest(http.MethodGet, "http://example.test/x", nil).WithContext(ctx)
	req.Header.Set(RequestIDHeader, "explicit-token")
	_, err := transport.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, "explicit-token", gotHeader)
}

func TestWithRequestIDHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.test/x", nil)
	out := WithRequestIDHeader(req, "rid-1")
	assert.Equal(t, "rid-1", out.Header.Get(RequestIDHeader))

	// 已显式设置过的不覆盖。
	req2 := httptest.NewRequest(http.MethodGet, "http://example.test/x", nil)
	req2.Header.Set(RequestIDHeader, "keep")
	assert.Equal(t, "keep", WithRequestIDHeader(req2, "rid-2").Header.Get(RequestIDHeader))

	assert.Nil(t, WithRequestIDHeader(nil, "rid"))
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func splitOnce(s string, sep byte) [2]string {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return [2]string{s[:i], s[i+1:]}
		}
	}
	return [2]string{s, ""}
}
