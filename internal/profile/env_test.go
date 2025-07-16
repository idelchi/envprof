package profile_test

import (
	"testing"

	"github.com/goccy/go-yaml"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/env"
)

// TestEnv_UnmarshalYAML tests unmarshaling environment variables from YAML in both map and sequence formats.
//
//nolint:gocognit // Complex table-driven tests with multiple scenarios
func TestEnv_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		yaml    string
		want    profile.Env
		wantErr bool
	}{
		{
			name: "map format with various types",
			yaml: `KEY1: value1
KEY2: 123
KEY3: true
KEY4: 3.14
KEY5: [1, 2, 3]
KEY6:
  nested: value
  count: 42`,
			want: profile.Env{
				"KEY1": "value1",
				"KEY2": 123,
				"KEY3": true,
				"KEY4": 3.14,
				"KEY5": []any{1, 2, 3},
				"KEY6": map[string]any{
					"nested": "value",
					"count":  42,
				},
			},
		},
		{
			name: "sequence format",
			yaml: `- KEY1=value1
- KEY2=123
- KEY3=true
- "KEY4=value with spaces"
- KEY5=`,
			want: profile.Env{
				"KEY1": "value1",
				"KEY2": "123", // Note: sequence format treats all as strings
				"KEY3": "true",
				"KEY4": "value with spaces",
				"KEY5": "",
			},
		},
		{
			name: "empty env",
			yaml: ``,
			want: profile.Env{},
		},
		{
			name: "null values",
			yaml: `KEY1: null
KEY2: ~
KEY3:`,
			want: profile.Env{
				"KEY1": nil,
				"KEY2": nil,
				"KEY3": nil,
			},
		},
		{
			name:    "invalid sequence format",
			yaml:    `- invalid format without equals`,
			wantErr: true,
		},
		{
			name:    "invalid type",
			yaml:    `123`, // Not a map or sequence
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var e profile.Env
			err := yaml.Unmarshal([]byte(tt.yaml), &e)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err == nil {
				if len(e) != len(tt.want) {
					t.Errorf("length mismatch: got %d, want %d", len(e), len(tt.want))

					return
				}

				for k, wantV := range tt.want {
					gotV, ok := e[k]
					if !ok {
						t.Errorf("missing key %q", k)

						continue
					}

					if !equalValues(gotV, wantV) {
						t.Errorf("key %q: got %v (%T), want %v (%T)", k, gotV, gotV, wantV, wantV)
					}
				}
			}
		})
	}
}

// TestEnv_Stringified tests conversion of environment variables to string format for shell export.
func TestEnv_Stringified(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		env     profile.Env
		want    env.Env
		wantErr bool
	}{
		{
			name: "simple strings",
			env: profile.Env{
				"KEY1": "value1",
				"KEY2": "value2",
			},
			want: env.Env{
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
		{
			name: "various types",
			env: profile.Env{
				"STRING": "hello",
				"NUMBER": 42,
				"FLOAT":  3.14,
				"BOOL":   true,
				"NIL":    nil,
				"ARRAY":  []int{1, 2, 3},
				"MAP":    map[string]any{"key": "value"},
			},
			want: env.Env{
				"STRING": "hello",
				"NUMBER": "42",
				"FLOAT":  "3.14",
				"BOOL":   "true",
				"NIL":    "",
				"ARRAY":  "'[1,2,3]'",
				"MAP":    `'{"key":"value"}'`,
			},
		},
		{
			name: "empty env",
			env:  profile.Env{},
			want: env.Env{},
		},
		{
			name: "strings with special characters",
			env: profile.Env{
				"QUOTED":  `"already quoted"`,
				"SPACES":  "has spaces",
				"NEWLINE": "has\nnewline",
				"EQUALS":  "key=value",
			},
			want: env.Env{
				"QUOTED":  `"\"already quoted\""`,
				"SPACES":  `"has spaces"`,
				"NEWLINE": `"has\nnewline"`,
				"EQUALS":  `"key=value"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.env.Stringified()
			if (err != nil) != tt.wantErr {
				t.Errorf("Stringified() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("length mismatch: got %d, want %d", len(got), len(tt.want))

				return
			}

			for k, wantV := range tt.want {
				gotV, ok := got[k]
				if !ok {
					t.Errorf("missing key %q", k)

					continue
				}

				if gotV != wantV {
					t.Errorf("key %q: got %q, want %q", k, gotV, wantV)
				}
			}
		})
	}
}

// equalValues compares two values, handling slices and maps specially.
//
//nolint:gocognit // Complex type comparison requires checking multiple type cases
func equalValues(a, b any) bool {
	// Handle numeric comparisons
	if isNumeric(a) && isNumeric(b) {
		return toFloat64(a) == toFloat64(b)
	}

	switch av := a.(type) {
	case []any:
		bv, ok := b.([]any)
		if !ok || len(av) != len(bv) {
			return false
		}

		for i := range av {
			if !equalValues(av[i], bv[i]) {
				return false
			}
		}

		return true
	case map[string]any:
		bv, ok := b.(map[string]any)
		if !ok || len(av) != len(bv) {
			return false
		}

		for k, v := range av {
			if bVal, exists := bv[k]; !exists || !equalValues(v, bVal) {
				return false
			}
		}

		return true
	default:
		return a == b
	}
}

// isNumeric checks if a value is a numeric type.
func isNumeric(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	}

	return false
}

// toFloat64 converts numeric values to float64 for comparison.
func toFloat64(v any) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	}

	return 0
}
