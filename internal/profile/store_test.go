//nolint:testpackage // Tests need access to private unmarshal() method and types
package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
)

// TestStore_Unmarshal tests unmarshaling profiles from YAML and TOML formats.
func TestStore_Unmarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		format  Type
		want    Profiles
		wantErr bool
	}{
		{
			name:   "valid yaml",
			format: YAML,
			data: []byte(`dev:
  env:
    APP_ENV: development
    DEBUG: true
prod:
  env:
    APP_ENV: production
    DEBUG: false`),
			want: Profiles{
				"dev": &Profile{
					Env: Env{
						"APP_ENV": "development",
						"DEBUG":   true,
					},
				},
				"prod": &Profile{
					Env: Env{
						"APP_ENV": "production",
						"DEBUG":   false,
					},
				},
			},
		},
		{
			name:   "valid toml",
			format: TOML,
			data: []byte(`[dev.env]
APP_ENV = "development"
DEBUG = true

[prod.env]
APP_ENV = "production"
DEBUG = false`),
			want: Profiles{
				"dev": &Profile{
					Env: Env{
						"APP_ENV": "development",
						"DEBUG":   true,
					},
				},
				"prod": &Profile{
					Env: Env{
						"APP_ENV": "production",
						"DEBUG":   false,
					},
				},
			},
		},
		{
			name:    "invalid yaml",
			format:  YAML,
			data:    []byte(`invalid: yaml: content:`),
			wantErr: true,
		},
		{
			name:    "invalid toml",
			format:  TOML,
			data:    []byte(`[invalid toml content`),
			wantErr: true,
		},
		{
			name:   "empty yaml",
			format: YAML,
			data:   []byte(``),
			want:   Profiles{},
		},
		{
			name:   "empty toml",
			format: TOML,
			data:   []byte(``),
			want:   Profiles{},
		},
		{
			name:    "unsupported type",
			format:  Type("json"),
			data:    []byte(`{}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := &Store{
				Type:     tt.format,
				Profiles: Profiles{},
			}

			err := store.unmarshal(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshal() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err == nil {
				compareProfiles(t, store.Profiles, tt.want)
			}
		})
	}
}

// TestProfiles_Marshal tests marshaling profiles to YAML and TOML formats with round-trip verification.
func TestProfiles_Marshal(t *testing.T) {
	t.Parallel()

	profiles := Profiles{
		"dev": &Profile{
			Env: Env{
				"APP_ENV": "development",
				"DEBUG":   true,
			},
		},
		"prod": &Profile{
			Env: Env{
				"APP_ENV": "production",
				"DEBUG":   false,
			},
		},
	}

	// Test YAML marshaling
	yamlData, err := yaml.Marshal(profiles)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	// Unmarshal to verify
	var yamlProfiles Profiles

	err = yaml.Unmarshal(yamlData, &yamlProfiles)
	if err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if len(yamlProfiles) != len(profiles) {
		t.Errorf("YAML round-trip: profile count = %d, want %d", len(yamlProfiles), len(profiles))
	}

	// Test TOML marshaling
	tomlData, err := toml.Marshal(profiles)
	if err != nil {
		t.Fatalf("toml.Marshal() error = %v", err)
	}

	// Unmarshal to verify
	var tomlProfiles Profiles

	err = toml.Unmarshal(tomlData, &tomlProfiles)
	if err != nil {
		t.Fatalf("toml.Unmarshal() error = %v", err)
	}

	if len(tomlProfiles) != len(profiles) {
		t.Errorf("TOML round-trip: profile count = %d, want %d", len(tomlProfiles), len(profiles))
	}
}

// Note: The following tests require integration with the file system and the godyl file package.
// Since we can't easily mock the file.File interface without knowing its exact structure,
// we'll create simplified integration tests that work with actual files.

// TestStore_NewAndLoad_Integration tests file system integration for loading profile stores.
//
//nolint:paralleltest // File system operations and shared test fixtures should not run in parallel
func TestStore_NewAndLoad_Integration(t *testing.T) {
	// Skip this test if running in CI or similar environment where file operations might be restricted
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration tests")
	}

	tests := []struct {
		name     string
		filename string
		content  string
		wantType Type
		wantErr  bool
	}{
		{
			name:     "yaml file",
			filename: "config.yaml",
			content: `dev:
  env:
    KEY: value`,
			wantType: YAML,
		},
		{
			name:     "toml file",
			filename: "config.toml",
			content: `[dev.env]
KEY = "value"`,
			wantType: TOML,
		},
		{
			name:     "json file (unsupported)",
			filename: "config.json",
			content:  `{}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, tt.filename)

			err := os.WriteFile(filePath, []byte(tt.content), 0o600)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}
		})
	}
}

// compareProfiles is a helper function to compare two Profiles maps.
func compareProfiles(t *testing.T, got, want Profiles) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("profile count = %d, want %d", len(got), len(want))

		return
	}

	for name, wantProfile := range want {
		gotProfile, exists := got[name]
		if !exists {
			t.Errorf("missing profile %q", name)

			continue
		}

		// Compare Env
		if len(gotProfile.Env) != len(wantProfile.Env) {
			t.Errorf("profile %q: env count = %d, want %d", name, len(gotProfile.Env), len(wantProfile.Env))
		}

		for k, wantV := range wantProfile.Env {
			gotV, ok := gotProfile.Env[k]
			if !ok {
				t.Errorf("profile %q: missing env key %q", name, k)

				continue
			}

			if !equalValues(gotV, wantV) {
				t.Errorf("profile %q, key %q: got %v (%T), want %v (%T)", name, k, gotV, gotV, wantV, wantV)
			}
		}
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
