package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Logger zerolog.Logger

// Init initializes the global logger with the specified level and format.
//
// Init writes the package-level Logger and the zerolog global level, so it is
// only safe to call once during process start-up (see docs/api-schema.md §1).
func Init(level, format string) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	zerolog.SetGlobalLevel(parseLevel(level))
	// Force debug level for testing
	// TODO(KT-07): this unconditional overwrite is a known defect — the level
	// argument has no effect. Kept as-is here (KT-05 fixes no behaviour) and
	// tracked separately.
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	Logger = zerolog.New(newWriter(format, os.Stdout)).With().Timestamp().Caller().Logger()
}

// WithContext returns a logger with additional context fields.
func WithContext(ctx map[string]any) zerolog.Logger {
	l := Logger
	for k, v := range ctx {
		l = l.With().Interface(k, v).Logger()
	}
	return l
}

// parseLevel maps the configuration string to a zerolog level. Unknown values
// fall back to info.
func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

// newWriter selects the log output format. Only "console" opts into the
// human-readable writer; every other value keeps the raw JSON stream.
func newWriter(format string, out io.Writer) io.Writer {
	if format == "console" {
		return zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}
	}
	return out
}
