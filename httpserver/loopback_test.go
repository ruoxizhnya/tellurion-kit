package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsLoopbackRemoteAddr 钉住判据本身（KT-10 从 core 逐字搬入）：
// 回环放行、非回环拒绝、不可解析一律拒绝（fail-closed）。
func TestIsLoopbackRemoteAddr(t *testing.T) {
	t.Parallel()

	cases := []struct {
		addr string
		want bool
	}{
		{"127.0.0.1:5000", true},
		{"127.1.2.3:1", true},
		{"[::1]:8080", true},
		{"localhost:8080", true},
		{"127.0.0.1", true},
		{"::1", true},
		{"[::1]", true},
		{"192.168.1.10:5000", false},
		{"10.0.0.1:5000", false},
		{"0.0.0.0:5000", false},
		{"[fe80::1%eth0]:5000", false},
		{"example.com:80", false},
		{"", false},
		{"not-an-addr", false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.addr, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, IsLoopbackRemoteAddr(tc.addr))
		})
	}
}
