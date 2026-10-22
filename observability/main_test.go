package observability

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestMain sets gin to test mode once for the entire package: gin's mode is
// a process-wide global, and per-test SetMode calls race with route
// registration. See httpserver's AUD-28 note — same pattern here.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
