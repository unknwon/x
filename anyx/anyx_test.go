package anyx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoalesce(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			name string
			in   []string
			want string
		}{
			{name: "returns first non-zero", in: []string{"a", "b"}, want: "a"},
			{name: "skips leading zero values", in: []string{"", "b", "c"}, want: "b"},
			{name: "returns zero value when all are zero", in: []string{"", "", ""}, want: ""},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Coalesce(tt.in...)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("int", func(t *testing.T) {
		tests := []struct {
			name string
			in   []int
			want int
		}{
			{name: "returns first non-zero", in: []int{1, 2}, want: 1},
			{name: "skips leading zeros", in: []int{0, 42, 0}, want: 42},
			{name: "returns zero value when all are zero", in: []int{0, 0, 0}, want: 0},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Coalesce(tt.in...)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("bool", func(t *testing.T) {
		tests := []struct {
			name string
			in   []bool
			want bool
		}{
			{name: "returns first true", in: []bool{true, false}, want: true},
			{name: "skips leading falses", in: []bool{false, true}, want: true},
			{name: "returns zero value when all are false", in: []bool{false, false}, want: false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Coalesce(tt.in...)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("pointer", func(t *testing.T) {
		a, b := new(int), new(int)
		tests := []struct {
			name string
			in   []*int
			want *int
		}{
			{name: "returns first non-nil", in: []*int{a, b}, want: a},
			{name: "skips leading nils", in: []*int{nil, b, a}, want: b},
			{name: "returns zero value when all are nil", in: []*int{nil, nil}, want: nil},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Coalesce(tt.in...)
				assert.Same(t, tt.want, got)
			})
		}
	})

	t.Run("mixed types via any", func(t *testing.T) {
		tests := []struct {
			name string
			in   []any
			want any
		}{
			{name: "returns first non-nil regardless of concrete type", in: []any{"a", 1, true}, want: "a"},
			{name: "skips leading nils across types", in: []any{nil, 42, "b"}, want: 42},
			{name: "skips leading nils to reach a bool", in: []any{nil, nil, true}, want: true},
			{name: "returns nil when all are nil", in: []any{nil, nil}, want: nil},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Coalesce(tt.in...)
				assert.Equal(t, tt.want, got)
			})
		}
	})
}
