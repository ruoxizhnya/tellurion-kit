// Package observability provides a Prometheus metrics registry plus
// request_id / trace propagation middleware for any service.
//
// The registry is deliberately **private** (not the global default
// registerer): tests can run in parallel without colliding on duplicate
// metric registration, and a service can expose exactly the metric
// families it chose to register — nothing leaks in from imported
// libraries.
//
// This package owns only the mechanism and a couple of generic metric
// families (outbound HTTP calls, cache hit ratio). Service-specific
// metric families are registered by their owning package via
// RegisterCollectors — the tool layer does not name them.
package observability

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// RequestIDHeader is the canonical header used to propagate the
// per-request correlation token. It is intentionally distinct from
// the W3C `traceparent` header so that OTel can be layered on top
// later without breaking existing tooling that scrapes X-Request-ID.
const RequestIDHeader = "X-Request-ID"

// Metrics is the central handle to a private Prometheus registry plus
// the generic metric families. Construct it once at process startup via
// NewMetrics and share it via dependency injection.
type Metrics struct {
	// HTTPRequests counts outbound HTTP calls, labeled by target
	// service and HTTP status.
	HTTPRequests *prometheus.CounterVec
	// CacheHitRatio is a rolling cache hit ratio in [0,1], labeled by
	// cache kind (the kind strings are the caller's vocabulary).
	CacheHitRatio *prometheus.GaugeVec

	registry *prometheus.Registry
}

// NewMetrics creates a Metrics handle with a private registry.
// Register() must be called before exposing /metrics so the built-in
// families are attached exactly once.
func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_client_requests_total",
			Help: "Count of outbound HTTP calls made by the service, labeled by target service and HTTP status.",
		}, []string{"service", "status"}),

		CacheHitRatio: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cache_hit_ratio",
			Help: "Rolling cache hit ratio in [0,1], labeled by cache kind.",
		}, []string{"kind"}),

		registry: prometheus.NewRegistry(),
	}
}

// Register attaches the built-in metric families to the underlying
// registry. The constructor does not auto-register so tests can inspect
// or re-create metrics; production code MUST call Register before
// exposing /metrics.
func (m *Metrics) Register() {
	m.registry.MustRegister(
		m.HTTPRequests,
		m.CacheHitRatio,
	)
}

// RegisterCollectors attaches additional collectors — Go process /
// runtime collectors, or metric families owned by a service package.
//
// This is the extension point that keeps service-specific metrics out
// of the tool layer: the owning package builds the prometheus collector
// and registers it here, so /metrics exposes it without this package
// naming it.
func (m *Metrics) RegisterCollectors(collectors ...prometheus.Collector) {
	for _, c := range collectors {
		m.registry.MustRegister(c)
	}
}

// Registry returns the underlying registry. Useful for tests that need
// to scrape /metrics directly via testutil.GatherAndCount.
func (m *Metrics) Registry() *prometheus.Registry { return m.registry }

// ObserveHTTP records one outbound HTTP call. A zero status (transport
// error before a response) is emitted as "0".
func (m *Metrics) ObserveHTTP(service string, status int) {
	if m == nil || m.HTTPRequests == nil {
		return
	}
	m.HTTPRequests.WithLabelValues(service, strconv.Itoa(status)).Inc()
}

// SetCacheHitRatio sets the cache hit ratio for a named cache kind.
// The ratio is clamped to [0, 1] to avoid emitting invalid gauge values
// that would confuse alerting rules.
func (m *Metrics) SetCacheHitRatio(kind string, ratio float64) {
	if m == nil || m.CacheHitRatio == nil {
		return
	}
	if ratio < 0 {
		ratio = 0
	} else if ratio > 1 {
		ratio = 1
	}
	m.CacheHitRatio.WithLabelValues(kind).Set(ratio)
}
