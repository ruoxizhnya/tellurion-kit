package auth

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestMain sets gin to test mode once for the entire package. Individual
// tests must NOT call gin.SetMode themselves: gin's mode is a process-wide
// global and the middleware tests run in parallel (see the guard in
// httpserver/ginmode_test.go).
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
