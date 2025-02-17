package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnion(t *testing.T) {
	type args[T comparable] struct {
		a []T
		b []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[string]{
		{
			name: "nil slices",
			args: args[string]{
				a: nil,
				b: nil,
			},
			want: []string{},
		},
		{
			name: "nil slices",
			args: args[string]{
				a: nil,
				b: nil,
			},
			want: []string{},
		},
		{
			name: "a nil",
			args: args[string]{
				a: nil,
				b: []string{"abc"},
			},
			want: []string{},
		},
		{
			name: "b nil",
			args: args[string]{
				b: nil,
				a: []string{"abc"},
			},
			want: []string{},
		},
		{
			name: "empty",
			args: args[string]{
				a: []string{},
				b: []string{},
			},
			want: []string{},
		},
		{
			name: "both with",
			args: args[string]{
				a: []string{"abc", "def", "ghi"},
				b: []string{"ghi", "abc", "def"},
			},
			want: []string{"abc", "def", "ghi"},
		},
		{
			name: "differences",
			args: args[string]{
				a: []string{"abc", "def", "ghi"},
				b: []string{"ghi", "asdf", "def"},
			},
			want: []string{"def", "ghi"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Union(tt.args.a, tt.args.b), "Union(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}

func TestRemove(t *testing.T) {
	type args[T comparable] struct {
		a []T
		b T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[string]{
		{
			name: "nil slices",
			args: args[string]{
				a: nil,
				b: "",
			},
			want: nil,
		},
		{
			name: "single in array no match",
			args: args[string]{
				a: []string{"abc"},
				b: "cef",
			},
			want: []string{"abc"},
		},
		{
			name: "single in array match",
			args: args[string]{
				a: []string{"abc"},
				b: "abc",
			},
			want: []string{},
		},
		{
			name: "multiple in array no match",
			args: args[string]{
				a: []string{"abc", "def"},
				b: "ghi",
			},
			want: []string{"abc", "def"},
		},
		{
			name: "multiple in array single match",
			args: args[string]{
				a: []string{"abc", "def"},
				b: "def",
			},
			want: []string{"abc"},
		},
		{
			name: "multiple in array multiple match",
			args: args[string]{
				a: []string{"abc", "def", "ghi", "jkl", "def"},
				b: "def",
			},
			want: []string{"abc", "ghi", "jkl"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Remove(tt.args.a, tt.args.b), "Remove(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}
