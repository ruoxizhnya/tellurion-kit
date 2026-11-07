package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Logger is the package-level logger populated by Init.
//
// Deprecated: build an independent logger with New and pass it to the
// components you construct instead of reading this global. Logger and Init are
// a compatibility thin shell for existing callers and will be removed in the
// next MAJOR release (v1.0.0); see CHANGELOG.md and docs/SCOPE.md.
var Logger zerolog.Logger

// New builds an independent logger for the given level and format that writes
// to out. It neither reads nor writes the package-level Logger, so two loggers
// created in the same process can write to different sinks and filter at
// different levels.
//
// The level mapping matches Init (unknown values fall back to info) and the
// record shape is the same (timestamp + caller fields; JSON unless format is
// "console"). The zerolog process-wide formatting settings that make up that
// shape (TimeFieldFormat / ErrorStackMarshaler) are format constants rather
// than a logger selection, and are applied here for parity with Init.
func New(level, format string, out io.Writer) zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	return zerolog.New(newWriter(format, out)).With().
		Timestamp().Caller().Logger().
		Level(parseLevel(level))
}

// Init initializes the global logger with the specified level and format.
//
// Init writes the package-level Logger and the zerolog global level, so it is
// only safe to call once during process start-up (see docs/api-schema.md §1).
//
// Deprecated: use New to build a logger and pass it to the components you
// construct. Init and Logger will be removed in the next MAJOR release
// (v1.0.0); see CHANGELOG.md and docs/SCOPE.md.
func Init(level, format string) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	zerolog.SetGlobalLevel(parseLevel(level))

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
