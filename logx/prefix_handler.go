package logx

import (
	"context"
	"log/slog"
)

// prefixHandler wraps a slog.Handler and prepends a fixed prefix to every
// record's message. Used by Logger.Scoped to add component context without
// relying on attribute grouping.
type prefixHandler struct {
	prefix  string
	handler slog.Handler
}

func newPrefixHandler(prefix string, handler slog.Handler) *prefixHandler {
	return &prefixHandler{prefix: prefix, handler: handler}
}

// Enabled reports whether the handler handles records at the given level.
// It delegates to the underlying handler.
func (h *prefixHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle handles the Record by prepending the configured prefix to the message
// and delegating to the underlying handler.
func (h *prefixHandler) Handle(ctx context.Context, r slog.Record) error {
	r2 := r.Clone()
	r2.Message = h.prefix + r2.Message
	return h.handler.Handle(ctx, r2)
}

// WithAttrs returns a new Handler whose attributes consist of the receiver's
// attributes followed by the arguments. The underlying handler is wrapped with
// the same prefix.
func (h *prefixHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newPrefixHandler(h.prefix, h.handler.WithAttrs(attrs))
}

// WithGroup returns a new Handler with the given group appended to the
// receiver's existing groups. The underlying handler is wrapped with the same
// prefix.
func (h *prefixHandler) WithGroup(name string) slog.Handler {
	return newPrefixHandler(h.prefix, h.handler.WithGroup(name))
}
