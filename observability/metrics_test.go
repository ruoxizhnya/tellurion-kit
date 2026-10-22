package observability

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetrics_RegisterAndObserve(t *testing.T) {
	m := NewMetrics()
	m.Register()

	m.ObserveHTTP("data", http.StatusOK)
	m.ObserveHTTP("data", http.StatusOK)
	m.ObserveHTTP("risk", http.StatusServiceUnavailable)
	assert.InDelta(t, 2, testutil.ToFloat64(m.HTTPRequests.WithLabelValues("data", "200")), 1e-9)
	assert.InDelta(t, 1, testutil.ToFloat64(m.HTTPRequests.WithLabelValues("risk", "503")), 1e-9)

	// 比例钳制到 [0,1]，负值与超界都不允许被写进 gauge。
	m.SetCacheHitRatio("factor", -0.5)
	assert.InDelta(t, 0, testutil.ToFloat64(m.CacheHitRatio.WithLabelValues("factor")), 1e-9)
	m.SetCacheHitRatio("factor", 1.5)
	assert.InDelta(t, 1, testutil.ToFloat64(m.CacheHitRatio.WithLabelValues("factor")), 1e-9)
	m.SetCacheHitRatio("sync_job", 0.25)
	assert.InDelta(t, 0.25, testutil.ToFloat64(m.CacheHitRatio.WithLabelValues("sync_job")), 1e-9)

	// 私有 registry 只暴露显式注册的两族：不含 Go runtime 默认采集器。
	assert.Equal(t, []string{"cache_hit_ratio", "http_client_requests_total"}, metricFamilyNames(t, m))
}

func TestRegisterCollectors_AcceptsCallerOwnedFamilies(t *testing.T) {
	m := NewMetrics()
	m.Register()

	custom := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "example_jobs_total",
		Help: "caller-owned family",
	})
	m.RegisterCollectors(custom)
	custom.Inc()

	assert.Contains(t, metricFamilyNames(t, m), "example_jobs_total",
		"RegisterCollectors 是服务自定义指标的扩展点，必须真的进 registry")
	assert.InDelta(t, 1, testutil.ToFloat64(custom), 1e-9)
}

func metricFamilyNames(t *testing.T, m *Metrics) []string {
	t.Helper()
	families, err := m.Registry().Gather()
	require.NoError(t, err)
	names := make([]string, 0, len(families))
	for _, f := range families {
		names = append(names, f.GetName())
	}
	sort.Strings(names)
	return names
}

func TestNilMetricsIsNoop(t *testing.T) {
	var m *Metrics
	assert.NotPanics(t, func() {
		m.ObserveHTTP("data", 200)
		m.SetCacheHitRatio("factor", 0.5)
	})
}

func TestHandler_NilMetricsIs503(t *testing.T) {
	r := gin.New()
	r.GET("/metrics", Handler(nil))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestHandler_ServesRegisteredRegistry(t *testing.T) {
	m := NewMetrics()
	m.Register()
	m.ObserveHTTP("data", 200)

	r := gin.New()
	r.GET("/metrics", Handler(m))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http_client_requests_total")
}
