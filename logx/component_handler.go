package logx

import (
	"context"
	"log/slog"
	"slices"
)

// ComponentKey is the attribute key under which Logger.Scoped records the
// component of every record, e.g., "component":"gorm".
const ComponentKey = "component"

// componentHandler adds the component attribute to every record. It keeps the
// handler from before the attribute was added, and replays later With and
// WithGroup calls on it, so a nested Logger.Scoped replaces the attribute with a
// longer dotted path instead of repeating the key, and the attribute is never
// nested inside a group.
type componentHandler struct {
	component string
	base      slog.Handler
	ops       []func(slog.Handler) slog.Handler
	handler   slog.Handler
}

func newComponentHandler(component string, base slog.Handler, ops []func(slog.Handler) slog.Handler) *componentHandler {
	handler := base.WithAttrs([]slog.Attr{slog.String(ComponentKey, component)})
	for _, op := range ops {
		handler = op(handler)
	}
	return &componentHandler{
		component: component,
		base:      base,
		ops:       ops,
		handler:   handler,
	}
}

func (h *componentHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *componentHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *componentHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.with(func(handler slog.Handler) slog.Handler { return handler.WithAttrs(attrs) })
}

func (h *componentHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.with(func(handler slog.Handler) slog.Handler { return handler.WithGroup(name) })
}

func (h *componentHandler) with(op func(slog.Handler) slog.Handler) *componentHandler {
	return &componentHandler{
		component: h.component,
		base:      h.base,
		ops:       append(slices.Clip(h.ops), op),
		handler:   op(h.handler),
	}
}
