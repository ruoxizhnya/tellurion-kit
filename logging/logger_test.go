package logging

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// restoreGlobals snapshots the package/global state Init mutates and puts it
// back when the test finishes. Tests in this package must not run in parallel:
// Init writes both Logger and the zerolog global level (see docs/api-schema.md
// §1 for the "call Init once at start-up" constraint).
func restoreGlobals(t *testing.T) {
	t.Helper()
	prevLogger := Logger
	prevLevel := zerolog.GlobalLevel()
	prevStackMarshaler := zerolog.ErrorStackMarshaler
	prevTimeFormat := zerolog.TimeFieldFormat

	t.Cleanup(func() {
		Logger = prevLogger
		zerolog.SetGlobalLevel(prevLevel)
		zerolog.ErrorStackMarshaler = prevStackMarshaler
		zerolog.TimeFieldFormat = prevTimeFormat
	})
}

// ---- parseLevel ----------------------------------------------------------

func TestParseLevel(t *testing.T) {
	tests := []struct {
		in   string
		want zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"", zerolog.InfoLevel},
		{"bogus", zerolog.InfoLevel},
		// Case-sensitive by design: an upper-case value is not recognised and
		// falls back to info. Pinned here so a future change is visible.
		{"DEBUG", zerolog.InfoLevel},
	}
	for _, tt := range tests {
		t.Run("level="+tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, parseLevel(tt.in))
		})
	}
}

// ---- newWriter -----------------------------------------------------------

func TestNewWriter_ConsoleReturnsConsoleWriter(t *testing.T) {
	var buf bytes.Buffer
	w := newWriter("console", &buf)
	require.IsType(t, zerolog.ConsoleWriter{}, w)
}

func TestNewWriter_ConsoleIsHumanReadable(t *testing.T) {
	var buf bytes.Buffer
	w := newWriter("console", &buf)
	_, err := w.Write([]byte(`{"level":"info","message":"hello console"}`))
	require.NoError(t, err)
	// The console writer re-renders JSON as text, so the raw JSON keys are gone.
	assert.Contains(t, buf.String(), "hello console")
	assert.NotContains(t, buf.String(), `"message":`)
}

func TestNewWriter_NonConsolePassesThrough(t *testing.T) {
	for _, format := range []string{"json", "", "text", "CONSOLE"} {
		t.Run("format="+format, func(t *testing.T) {
			var buf bytes.Buffer
			w := newWriter(format, &buf)
			assert.Equal(t, &buf, w, "non-console formats must write the given writer untouched")

			_, err := w.Write([]byte(`{"a":1}`))
			require.NoError(t, err)
			assert.Equal(t, `{"a":1}`, buf.String())
		})
	}
}

// ---- Init ----------------------------------------------------------------

// TestInit_LevelIsAlwaysForcedToDebug pins today's behaviour: the level
// argument is parsed but then unconditionally overwritten with DebugLevel.
// This is a known defect tracked as KT-07 — the test documents the current
// contract, not the desired one. When KT-07 is fixed this test must change.
func TestInit_LevelIsAlwaysForcedToDebug(t *testing.T) {
	restoreGlobals(t)

	for _, level := range []string{"debug", "info", "warn", "error", "bogus"} {
		t.Run("level="+level, func(t *testing.T) {
			Init(level, "json")
			assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
		})
	}
}

func TestInit_SetsStackMarshalerAndTimeFormat(t *testing.T) {
	restoreGlobals(t)

	Init("info", "json")
	assert.NotNil(t, zerolog.ErrorStackMarshaler)
	assert.Equal(t, time.RFC3339Nano, zerolog.TimeFieldFormat)
}

func TestInit_ConsoleFormatStillInitialisesLogger(t *testing.T) {
	restoreGlobals(t)

	// Smoke test: the console branch must not panic and must leave a usable
	// logger behind.
	Init("info", "console")
	assert.NotPanics(t, func() {
		Logger.Info().Msg("console path smoke test")
	})
}

// ---- WithContext ---------------------------------------------------------

func TestWithContext_InjectsFields(t *testing.T) {
	restoreGlobals(t)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	var buf bytes.Buffer
	Logger = zerolog.New(&buf)

	l := WithContext(map[string]any{
		"request_id": "req-1",
		"attempt":    2,
	})
	l.Info().Msg("hello")

	out := buf.String()
	assert.Contains(t, out, `"request_id":"req-1"`)
	assert.Contains(t, out, `"attempt":2`)
	assert.Contains(t, out, `"message":"hello"`)
}

func TestWithContext_DoesNotMutateBaseLogger(t *testing.T) {
	restoreGlobals(t)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	var buf bytes.Buffer
	base := zerolog.New(&buf).With().Str("keep", "yes").Logger()
	Logger = base

	child := WithContext(map[string]any{"extra": "x"})
	child.Info().Msg("child")
	base.Info().Msg("base")

	// The child inherits the same writer, so tell the two records apart by
	// their message rather than by separate buffers.
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, lines, 2)

	childLine, baseLine := lines[0], lines[1]
	assert.Contains(t, childLine, `"message":"child"`)
	assert.Contains(t, childLine, `"extra":"x"`, "child logger must carry the injected field")
	assert.Contains(t, baseLine, `"message":"base"`)
	assert.NotContains(t, baseLine, `"extra":"x"`, "base logger must stay untouched")
	assert.Contains(t, baseLine, `"keep":"yes"`, "base logger must keep its own context")
}

func TestWithContext_EmptyMapReturnsBaseLogger(t *testing.T) {
	restoreGlobals(t)

	var buf bytes.Buffer
	Logger = zerolog.New(&buf).With().Str("keep", "yes").Logger()

	l := WithContext(map[string]any{})
	l.Info().Msg("hello")

	assert.Contains(t, buf.String(), `"keep":"yes"`)
}
