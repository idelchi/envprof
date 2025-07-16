//nolint:testpackage // Tests need access to private types and functions
package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/idelchi/godyl/pkg/env"
)

// TestProfiles_BasicOperations tests basic profile operations including existence checks, name retrieval, and profile
// creation.
//
//nolint:paralleltest // File system operations and shared test fixtures should not run in parallel
func TestProfiles_BasicOperations(t *testing.T) {
	profiles := Profiles{
		"dev": &Profile{
			Env: Env{"APP_ENV": "development"},
		},
		"prod": &Profile{
			Env: Env{"APP_ENV": "production"},
		},
	}

	// Test Exists
	if !profiles.Exists("dev") {
		t.Error("Exists('dev') should return true")
	}

	if profiles.Exists("staging") {
		t.Error("Exists('staging') should return false")
	}

	// Test Names
	names := profiles.Names()
	if len(names) != 2 {
		t.Errorf("Names() returned %d names, want 2", len(names))
	}

	// Test Create
	profiles.Create("staging")
	profiles.Create("dev") // Should not error even if exists

	if !profiles.Exists("staging") {
		t.Error("staging profile should exist after creation")
	}
}

// TestProfiles_Environment tests environment variable resolution with various inheritance scenarios including single,
// multiple, and circular dependencies.
//
//nolint:gocognit,paralleltest // Complex inheritance tests; File operations should not run in parallel
func TestProfiles_Environment(t *testing.T) {
	tests := []struct {
		name     string
		profiles Profiles
		target   string
		want     env.Env
		wantErr  bool
	}{
		{
			name: "simple profile",
			profiles: Profiles{
				"dev": &Profile{
					Env: Env{
						"APP_ENV": "development",
						"DEBUG":   "true",
					},
				},
			},
			target: "dev",
			want: env.Env{
				"APP_ENV": "development",
				"DEBUG":   "true",
			},
		},
		{
			name: "single inheritance",
			profiles: Profiles{
				"base": &Profile{
					Env: Env{
						"APP_NAME": "myapp",
						"VERSION":  "1.0",
						"DEBUG":    "false",
					},
				},
				"dev": &Profile{
					Extends: []string{"base"},
					Env: Env{
						"APP_ENV": "development",
						"DEBUG":   "true", // Override base
					},
				},
			},
			target: "dev",
			want: env.Env{
				"APP_NAME": "myapp",       // Inherited
				"VERSION":  "1.0",         // Inherited
				"DEBUG":    "true",        // Overridden
				"APP_ENV":  "development", // New
			},
		},
		{
			name: "multiple inheritance",
			profiles: Profiles{
				"base": &Profile{
					Env: Env{
						"APP_NAME": "myapp",
						"VERSION":  "1.0",
					},
				},
				"common": &Profile{
					Env: Env{
						"LOG_LEVEL": "info",
						"TIMEOUT":   "30",
					},
				},
				"dev": &Profile{
					Extends: []string{"base", "common"},
					Env: Env{
						"APP_ENV":   "development",
						"LOG_LEVEL": "debug", // Override common
					},
				},
			},
			target: "dev",
			want: env.Env{
				"APP_NAME":  "myapp",       // From base
				"VERSION":   "1.0",         // From base
				"TIMEOUT":   "30",          // From common
				"LOG_LEVEL": "debug",       // Overridden
				"APP_ENV":   "development", // New
			},
		},
		{
			name: "deep inheritance chain",
			profiles: Profiles{
				"base": &Profile{
					Env: Env{"LEVEL": "1", "BASE": "true"},
				},
				"middle": &Profile{
					Extends: []string{"base"},
					Env:     Env{"LEVEL": "2", "MIDDLE": "true"},
				},
				"top": &Profile{
					Extends: []string{"middle"},
					Env:     Env{"LEVEL": "3", "TOP": "true"},
				},
			},
			target: "top",
			want: env.Env{
				"BASE":   "true", // From base
				"MIDDLE": "true", // From middle
				"TOP":    "true", // From top
				"LEVEL":  "3",    // Overridden twice
			},
		},
		{
			name: "circular dependency",
			profiles: Profiles{
				"a": &Profile{
					Extends: []string{"b"},
					Env:     Env{"VAR_A": "a"},
				},
				"b": &Profile{
					Extends: []string{"c"},
					Env:     Env{"VAR_B": "b"},
				},
				"c": &Profile{
					Extends: []string{"a"}, // Creates cycle
					Env:     Env{"VAR_C": "c"},
				},
			},
			target:  "a",
			wantErr: true,
		},
		{
			name: "non-existent profile",
			profiles: Profiles{
				"dev": &Profile{
					Env: Env{"APP_ENV": "development"},
				},
			},
			target:  "staging",
			wantErr: true,
		},
		{
			name: "extends non-existent profile",
			profiles: Profiles{
				"dev": &Profile{
					Extends: []string{"missing"},
					Env:     Env{"APP_ENV": "development"},
				},
			},
			target:  "dev",
			wantErr: true,
		},
		{
			name:     "empty profiles",
			profiles: Profiles{},
			target:   "any",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.profiles.Environment(tt.target)

			if (err != nil) != tt.wantErr {
				t.Errorf("Environment() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err == nil {
				if len(got.Env) != len(tt.want) {
					t.Errorf("env count = %d, want %d", len(got.Env), len(tt.want))

					return
				}

				for k, wantV := range tt.want {
					gotV := got.Env.Get(k)
					if gotV == "" && wantV != "" {
						t.Errorf("missing key %q", k)

						continue
					}

					if gotV != wantV {
						t.Errorf("key %q = %q, want %q", k, gotV, wantV)
					}
				}
			}
		})
	}
}

// TestProfiles_EnvironmentWithDotEnv tests loading environment variables from .env files and profile override
// precedence.
//
//nolint:paralleltest // File system operations and shared test fixtures should not run in parallel
func TestProfiles_EnvironmentWithDotEnv(t *testing.T) {
	// Create test .env file
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "test.env")
	envContent := `FROM_ENV=dotenv_value
DEBUG=false
OVERRIDE_ME=from_env`

	err := os.WriteFile(envFile, []byte(envContent), 0o600)
	if err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	profiles := Profiles{
		"dev": &Profile{
			DotEnv: []string{envFile},
			Env: Env{
				"APP_ENV":     "development",
				"DEBUG":       "true",      // Override .env
				"OVERRIDE_ME": "from_prof", // Override .env
			},
		},
	}

	got, err := profiles.Environment("dev")
	if err != nil {
		t.Fatalf("Environment() error = %v", err)
	}

	want := env.Env{
		"FROM_ENV":    "dotenv_value", // From .env only
		"DEBUG":       "true",         // Profile overrides .env
		"OVERRIDE_ME": "from_prof",    // Profile overrides .env
		"APP_ENV":     "development",  // From profile only
	}

	if len(got.Env) != len(want) {
		t.Errorf("env count = %d, want %d", len(got.Env), len(want))
	}

	for k, wantV := range want {
		gotV := got.Env.Get(k)
		if gotV == "" && wantV != "" {
			t.Errorf("missing key %q", k)

			continue
		}

		if gotV != wantV {
			t.Errorf("key %q = %q, want %q", k, gotV, wantV)
		}
	}
}

// TestProfiles_ComplexInheritanceWithDotEnv tests complex inheritance scenarios with multiple .env files and profile
// overrides.
//
//nolint:paralleltest // File system operations and shared test fixtures should not run in parallel
func TestProfiles_ComplexInheritanceWithDotEnv(t *testing.T) {
	// Create test .env files
	tmpDir := t.TempDir()

	baseEnv := filepath.Join(tmpDir, "base.env")

	err := os.WriteFile(baseEnv, []byte("BASE_ENV=from_base_env\nSHARED=base"), 0o600)
	if err != nil {
		t.Fatalf("failed to write base.env: %v", err)
	}

	devEnv := filepath.Join(tmpDir, "dev.env")

	err = os.WriteFile(devEnv, []byte("DEV_ENV=from_dev_env\nSHARED=dev"), 0o600)
	if err != nil {
		t.Fatalf("failed to write dev.env: %v", err)
	}

	profiles := Profiles{
		"base": &Profile{
			DotEnv: []string{baseEnv},
			Env: Env{
				"APP_NAME": "myapp",
				"VERSION":  "1.0",
				"SHARED":   "base_profile", // Override base.env
			},
		},
		"dev": &Profile{
			Extends: []string{"base"},
			DotEnv:  []string{devEnv},
			Env: Env{
				"APP_ENV": "development",
				"VERSION": "dev",         // Override base profile
				"SHARED":  "dev_profile", // Override everything
			},
		},
	}

	got, err := profiles.Environment("dev")
	if err != nil {
		t.Fatalf("Environment() error = %v", err)
	}

	want := env.Env{
		"DEV_ENV":  "from_dev_env", // From dev.env
		"APP_NAME": "myapp",        // From base profile
		"APP_ENV":  "development",  // From dev profile
		"VERSION":  "dev",          // Dev overrides base
		"SHARED":   "dev_profile",  // Dev profile overrides all
	}

	for k, wantV := range want {
		gotV := got.Env.Get(k)
		if gotV == "" && wantV != "" {
			t.Errorf("missing key %q", k)

			continue
		}

		if gotV != wantV {
			t.Errorf("key %q = %q, want %q", k, gotV, wantV)
		}
	}
}

// TestProfiles_EnvironmentPrecedence tests the precedence order: .env < inherited profiles < current profile.
//
//nolint:paralleltest // File system operations and shared test fixtures should not run in parallel
func TestProfiles_EnvironmentPrecedence(t *testing.T) {
	// Test precedence: .env < inherited profiles < current profile
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, "test.env")

	err := os.WriteFile(envFile, []byte("VAR=from_env"), 0o600)
	if err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	profiles := Profiles{
		"base1": &Profile{
			Env: Env{"VAR": "from_base1"},
		},
		"base2": &Profile{
			Env: Env{"VAR": "from_base2"},
		},
		"dev": &Profile{
			Extends: []string{"base1", "base2"}, // base2 should override base1
			DotEnv:  []string{envFile},
			Env:     Env{"VAR": "from_dev"}, // Should override everything
		},
	}

	got, err := profiles.Environment("dev")
	if err != nil {
		t.Fatalf("Environment() error = %v", err)
	}

	if got.Env.Get("VAR") != "from_dev" {
		t.Errorf("VAR = %q, want 'from_dev' (should have highest precedence)", got.Env.Get("VAR"))
	}
}
