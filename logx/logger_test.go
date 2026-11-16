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
	decode := func(t *testing.T, buf *bytes.Buffer) map[string]any {
		t.Helper()
		var rec map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
		return rec
	}

	t.Run("records the component as an attribute", func(t *testing.T) {
		var buf bytes.Buffer
		l := New(slog.NewJSONHandler(&buf, nil)).Scoped("svc")

		l.Info("event", slog.String("trace", "abc123"))

		rec := decode(t, &buf)
		require.Equal(t, "event", rec["msg"])
		require.Equal(t, "svc", rec[ComponentKey])
		require.Equal(t, "abc123", rec["trace"])
	})

	t.Run("joins nested components with dots in one attribute", func(t *testing.T) {
		var buf bytes.Buffer
		l := New(slog.NewTextHandler(&buf, nil)).Scoped("a").With("id", 1).Scoped("b").Scoped("c")

		l.Info("msg")

		out := buf.String()
		require.Equal(t, 1, strings.Count(out, ComponentKey+"="), out)
		require.Contains(t, out, "component=a.b.c")
		require.Contains(t, out, "id=1")
		require.Contains(t, out, "msg=msg")
	})

	t.Run("keeps the component out of groups", func(t *testing.T) {
		var buf bytes.Buffer
		l := New(slog.NewJSONHandler(&buf, nil)).Scoped("worker").WithGroup("task").Scoped("cleanup")

		l.Info("done", "name", "tmp")

		rec := decode(t, &buf)
		require.Equal(t, "worker.cleanup", rec[ComponentKey])
		task, ok := rec["task"].(map[string]any)
		require.True(t, ok, "expected task group")
		require.Equal(t, "tmp", task["name"])
		require.NotContains(t, task, ComponentKey)
	})

	t.Run("does not change the parent logger", func(t *testing.T) {
		var buf bytes.Buffer
		parent := New(slog.NewJSONHandler(&buf, nil)).Scoped("a")
		_ = parent.Scoped("b")

		parent.Info("msg")

		require.Equal(t, "a", decode(t, &buf)[ComponentKey])
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
		require.Equal(t, "processing", rec["msg"])
		require.Equal(t, "worker", rec[ComponentKey])
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
