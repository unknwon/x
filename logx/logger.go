package logx

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Logger wraps *slog.Logger with scoping and builder methods that preserve the
// Logger type.
type Logger struct {
	*slog.Logger
}

// New creates a Logger backed by the given slog.Handler.
// All standard logging methods (Info, Error, etc.) are promoted from the
// embedded *slog.Logger.
func New(h slog.Handler) *Logger {
	return &Logger{Logger: slog.New(h)}
}

// NewNoopLogger returns a Logger whose output is discarded.
// Useful in tests or when logging must be suppressed entirely.
func NewNoopLogger() *Logger {
	return New(slog.NewTextHandler(io.Discard, nil))
}

// Scoped returns a child Logger whose every record has the ComponentKey
// attribute set to name. Scoping a scoped Logger joins the names with dots, e.g.,
// Scoped("db").Scoped("gorm") logs "component":"db.gorm", so log explorers can
// filter a component and everything under it by prefix. The attribute stays at
// the top level even when the parent Logger has groups.
func (l *Logger) Scoped(name string) *Logger {
	if h, ok := l.Handler().(*componentHandler); ok {
		return New(newComponentHandler(h.component+"."+name, h.base, h.ops))
	}
	return New(newComponentHandler(name, l.Handler(), nil))
}

// With returns a child Logger that includes the given key-value pairs or Attrs
// in every log record. The arguments follow the same rules as slog.Logger.With.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{Logger: l.Logger.With(args...)}
}

// WithGroup returns a child Logger that nests all subsequent attributes under
// the provided group name. Equivalent to slog.Logger.WithGroup.
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{Logger: l.Logger.WithGroup(name)}
}

// FatalContext logs msg at ERROR level with the given context and arguments,
// then terminates the process with exit code 1.
// The args are handled identically to slog.Logger.ErrorContext.
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}
