package logx

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger_New(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := New(h)

	l.Info("hello", "key", "value")

	out := buf.String()
	require.Contains(t, out, "hello")
	require.Contains(t, out, "key=value")
}

func TestLogger_NewNoopLogger(t *testing.T) {
	l := NewNoopLogger()
	require.NotNil(t, l)

	l.Info("should be discarded")
	l.Error("also discarded")
}

func TestLogger_Scoped(t *testing.T) {
	t.Run("adds prefix to messages", func(t *testing.T) {
		var buf bytes.Buffer
		h := slog.NewTextHandler(&buf, nil)
		l := New(h).Scoped("component")

		l.Info("doing work")

		out := buf.String()
		require.Contains(t, out, "[component] doing work")
	})

	t.Run("chains prefixes when called repeatedly", func(t *testing.T) {
		var buf bytes.Buffer
		h := slog.NewTextHandler(&buf, nil)
		l := New(h).Scoped("a").Scoped("b")

		l.Info("msg")

		out := buf.String()
		require.Contains(t, out, "[a] [b] msg")
	})

	t.Run("preserves attributes on the record", func(t *testing.T) {
		var buf bytes.Buffer
		h := slog.NewJSONHandler(&buf, nil)
		l := New(h).Scoped("svc")

		l.Info("event", slog.String("trace", "abc123"))

		var rec map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
		require.Equal(t, "[svc] event", rec["msg"])
		require.Equal(t, "abc123", rec["trace"])
	})
}

func TestLogger_With(t *testing.T) {
	t.Run("attaches key value pairs", func(t *testing.T) {
		var buf bytes.Buffer
		h := slog.NewJSONHandler(&buf, nil)
		l := New(h).With("service", "api", slog.Int("port", 8080))

		l.Info("started")

		var rec map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
		require.Equal(t, "api", rec["service"])
		require.Equal(t, float64(8080), rec["port"])
	})

	t.Run("nests attributes under group", func(t *testing.T) {
		var buf bytes.Buffer
		h := slog.NewJSONHandler(&buf, nil)
		l := New(h).WithGroup("http")

		l.Info("request", "method", "GET")

		var rec map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
		http, ok := rec["http"].(map[string]any)
		require.True(t, ok, "expected http group")
		require.Equal(t, "GET", http["method"])
	})

	t.Run("composes with Scoped and WithGroup", func(t *testing.T) {
		var buf bytes.Buffer
		h := slog.NewJSONHandler(&buf, nil)
		l := New(h).Scoped("worker").With("id", 42).WithGroup("task")

		l.Info("processing", "name", "cleanup")

		var rec map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
		require.Equal(t, "[worker] processing", rec["msg"])
		require.Equal(t, float64(42), rec["id"])

		task, ok := rec["task"].(map[string]any)
		require.True(t, ok)
		require.Equal(t, "cleanup", task["name"])
	})
}

func TestLogger_StandardLogMethods(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := New(h)

	ctx := context.Background()

	l.Debug("debug msg")
	l.DebugContext(ctx, "debug ctx")
	l.Info("info msg")
	l.InfoContext(ctx, "info ctx")
	l.Warn("warn msg")
	l.WarnContext(ctx, "warn ctx")
	l.Error("error msg")
	l.ErrorContext(ctx, "error ctx")

	out := buf.String()
	require.True(t, strings.Contains(out, "debug msg"))
	require.True(t, strings.Contains(out, "info msg"))
	require.True(t, strings.Contains(out, "warn msg"))
	require.True(t, strings.Contains(out, "error msg"))
}
