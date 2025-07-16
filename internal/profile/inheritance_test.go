package profile_test

import (
	"strings"
	"testing"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/env"
)

// TestInheritanceTracker_Inheritance tests basic inheritance tracking for environment variables.
func TestInheritanceTracker_Inheritance(t *testing.T) {
	t.Parallel()

	tracker := profile.InheritanceTracker{
		Inheritance: profile.Inheritance{
			"APP_NAME": "base",
			"DEBUG":    "dev",
			"PORT":     "local.env",
			"VERSION":  "prod",
		},
	}

	// Test inheritance tracking
	if source := tracker.Inheritance["APP_NAME"]; source != "base" {
		t.Errorf("Inheritance['APP_NAME'] = %q, want 'base'", source)
	}

	if source := tracker.Inheritance["DEBUG"]; source != "dev" {
		t.Errorf("Inheritance['DEBUG'] = %q, want 'dev'", source)
	}

	if source := tracker.Inheritance["UNKNOWN"]; source != "" {
		t.Errorf("Inheritance['UNKNOWN'] = %q, want empty string", source)
	}
}

// TestInheritanceTracker_Format tests formatting of environment variables with optional inheritance information.
func TestInheritanceTracker_Format(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tracker profile.InheritanceTracker
		key     string
		verbose bool
		withKey bool
		want    string
	}{
		{
			name: "non-verbose with key",
			tracker: profile.InheritanceTracker{
				Inheritance: profile.Inheritance{"KEY1": "dev"},
				Env:         env.Env{"KEY1": "value1"},
				Name:        "dev",
			},
			key:     "KEY1",
			verbose: false,
			withKey: true,
			want:    "KEY1=value1",
		},
		{
			name: "non-verbose without key",
			tracker: profile.InheritanceTracker{
				Inheritance: profile.Inheritance{"KEY1": "dev"},
				Env:         env.Env{"KEY1": "value1"},
				Name:        "dev",
			},
			key:     "KEY1",
			verbose: false,
			withKey: false,
			want:    "value1",
		},
		{
			name: "verbose with inheritance from different source",
			tracker: profile.InheritanceTracker{
				Inheritance: profile.Inheritance{"KEY1": "base"},
				Env:         env.Env{"KEY1": "value1"},
				Name:        "dev",
			},
			key:     "KEY1",
			verbose: true,
			withKey: true,
			want:    "KEY1=value1                                                  (inherited from \"base\")",
		},
		{
			name: "verbose with inheritance from same source",
			tracker: profile.InheritanceTracker{
				Inheritance: profile.Inheritance{"KEY1": "dev"},
				Env:         env.Env{"KEY1": "value1"},
				Name:        "dev",
			},
			key:     "KEY1",
			verbose: true,
			withKey: true,
			want:    "KEY1=value1",
		},
		{
			name: "value with special characters",
			tracker: profile.InheritanceTracker{
				Inheritance: profile.Inheritance{"PATH": "system.env"},
				Env:         env.Env{"PATH": "/usr/bin:/usr/local/bin"},
				Name:        "dev",
			},
			key:     "PATH",
			verbose: true,
			withKey: true,
			want:    "PATH=/usr/bin:/usr/local/bin                                 (inherited from \"system.env\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.tracker.Format(tt.key, tt.verbose, tt.withKey)

			if got != tt.want {
				t.Errorf("Format() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestInheritanceTracker_FormatAll tests formatting all environment variables with verbose inheritance annotations.
func TestInheritanceTracker_FormatAll(t *testing.T) {
	t.Parallel()

	tracker := profile.InheritanceTracker{
		Inheritance: profile.Inheritance{
			"APP_NAME": "base",
			"DEBUG":    "dev.env",
			"PORT":     "dev",
		},
		Env: env.Env{
			"APP_NAME":  "myapp",
			"DEBUG":     "true",
			"PORT":      "3000",
			"UNTRACKED": "value",
		},
		Name: "current",
	}

	// Test non-verbose
	result := tracker.FormatAll("", false)
	lines := strings.Split(strings.TrimSpace(result), "\n")

	if len(lines) != 4 {
		t.Errorf("FormatAll (non-verbose) returned %d lines, want 4", len(lines))
	}

	// Check that lines don't contain inheritance info
	for _, line := range lines {
		if strings.Contains(line, "inherited from") {
			t.Errorf("Non-verbose output should not contain inheritance info: %q", line)
		}
	}

	// Test verbose
	result = tracker.FormatAll("", true)
	lines = strings.Split(strings.TrimSpace(result), "\n")

	if len(lines) != 4 {
		t.Errorf("FormatAll (verbose) returned %d lines, want 4", len(lines))
	}

	// Check specific content and inheritance annotations
	if !strings.Contains(result, "APP_NAME=myapp") {
		t.Error("Missing APP_NAME in output")
	}

	if !strings.Contains(result, "inherited from \"base\"") {
		t.Error("Missing inheritance info for APP_NAME")
	}

	if !strings.Contains(result, "inherited from \"dev.env\"") {
		t.Error("Missing inheritance info for DEBUG")
	}

	if !strings.Contains(result, "UNTRACKED=value") {
		t.Error("Missing UNTRACKED in output")
	}

	// Test with prefix
	result = tracker.FormatAll("export ", false)
	lines = strings.Split(strings.TrimSpace(result), "\n")

	for _, line := range lines {
		if !strings.HasPrefix(line, "export ") {
			t.Errorf("Line should start with prefix 'export ': %q", line)
		}
	}
}

// TestInheritanceTracker_EmptyTracker tests behavior with empty inheritance tracking.
func TestInheritanceTracker_EmptyTracker(t *testing.T) {
	t.Parallel()

	tracker := profile.InheritanceTracker{
		Inheritance: profile.Inheritance{},
		Env:         env.Env{"KEY": "value"},
		Name:        "current",
	}

	// Test empty inheritance
	if source := tracker.Inheritance["ANY_KEY"]; source != "" {
		t.Errorf("Empty inheritance[ANY_KEY] = %q, want empty string", source)
	}

	// Format with empty inheritance should work
	formatted := tracker.Format("KEY", true, true)
	if formatted != "KEY=value" {
		t.Errorf("Format with empty inheritance = %q, want 'KEY=value'", formatted)
	}

	// FormatAll with empty tracker
	tracker.Env = env.Env{"KEY1": "value1", "KEY2": "value2"}
	result := tracker.FormatAll("", true)

	// Should not contain any inheritance info
	if strings.Contains(result, "inherited from") {
		t.Error("Empty inheritance should not add inheritance comments")
	}
}
