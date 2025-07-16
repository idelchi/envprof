package profile_test

import (
	"testing"

	"github.com/idelchi/envprof/internal/profile"
)

// TestStringify tests string conversion of various data types including primitives, slices, maps, and special
// characters.
func TestStringify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		want  string
	}{
		// Simple types
		{
			name:  "simple string",
			value: "hello",
			want:  "hello",
		},
		{
			name:  "string with spaces",
			value: "hello world",
			want:  `"hello world"`,
		},
		{
			name:  "string with quotes",
			value: `he said "hello"`,
			want:  `"he said \"hello\""`,
		},
		{
			name:  "string with newline",
			value: "line1\nline2",
			want:  `"line1\nline2"`,
		},
		{
			name:  "string with equals",
			value: "key=value",
			want:  `"key=value"`,
		},
		{
			name:  "empty string",
			value: "",
			want:  "",
		},
		{
			name:  "already quoted string",
			value: `"already quoted"`,
			want:  `"\"already quoted\""`,
		},
		{
			name:  "boolean true",
			value: true,
			want:  "true",
		},
		{
			name:  "boolean false",
			value: false,
			want:  "false",
		},
		{
			name:  "integer",
			value: 42,
			want:  "42",
		},
		{
			name:  "negative integer",
			value: -123,
			want:  "-123",
		},
		{
			name:  "float",
			value: 3.14,
			want:  "3.14",
		},
		{
			name:  "float with many decimals",
			value: 3.14159265359,
			want:  "3.14159265359",
		},
		{
			name:  "nil value",
			value: nil,
			want:  "",
		},

		// Complex types
		{
			name:  "integer slice",
			value: []int{1, 2, 3},
			want:  `'[1,2,3]'`,
		},
		{
			name:  "string slice",
			value: []string{"a", "b", "c"},
			want:  `'["a","b","c"]'`,
		},
		{
			name:  "mixed slice",
			value: []any{1, "two", true, 3.14},
			want:  `'[1,"two",true,3.14]'`,
		},
		{
			name:  "empty slice",
			value: []any{},
			want:  `'[]'`,
		},
		{
			name:  "simple map",
			value: map[string]int{"a": 1, "b": 2},
			want:  `'{"a":1,"b":2}'`,
		},
		{
			name: "complex map",
			value: map[string]any{
				"string": "value",
				"number": 123,
				"bool":   true,
				"null":   nil,
			},
			want: `'{"bool":true,"null":null,"number":123,"string":"value"}'`,
		},
		{
			name: "nested structure",
			value: map[string]any{
				"array": []int{1, 2, 3},
				"map": map[string]string{
					"key": "value",
				},
			},
			want: `'{"array":[1,2,3],"map":{"key":"value"}}'`,
		},
		{
			name:  "empty map",
			value: map[string]any{},
			want:  `'{}'`,
		},

		// Other types
		{
			name:  "uint",
			value: uint(100),
			want:  "100",
		},
		{
			name:  "int64",
			value: int64(9223372036854775807),
			want:  "9223372036854775807",
		},
		{
			name:  "float32",
			value: float32(3.14),
			want:  "3.14",
		},
		{
			name:  "byte slice as string",
			value: []byte("hello"),
			want:  `'"aGVsbG8="'`, // JSON marshals byte slices as base64
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := profile.Stringify(tt.value)
			if err != nil {
				t.Fatalf("Stringify() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("stringify(%v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}
