//nolint:testpackage // Tests need access to private newProfile() function
package profile

import (
	"testing"
)

// TestNewProfile tests that newProfile properly initializes all fields with empty slices and maps.
func TestNewProfile(t *testing.T) {
	t.Parallel()

	p := newProfile()

	if p.Env == nil {
		t.Error("Env should be initialized")
	}

	if p.DotEnv == nil {
		t.Error("DotEnv should be initialized")
	}

	if p.Extends == nil {
		t.Error("Extends should be initialized")
	}

	if len(p.Env) != 0 {
		t.Error("Env should be empty")
	}

	if len(p.DotEnv) != 0 {
		t.Error("DotEnv should be empty")
	}

	if len(p.Extends) != 0 {
		t.Error("Extends should be empty")
	}
}

// TestProfile_BasicFields tests basic profile field operations with various data types.
func TestProfile_BasicFields(t *testing.T) {
	t.Parallel()

	p := &Profile{
		Env: Env{
			"KEY1": "value1",
			"KEY2": 123,
			"KEY3": true,
		},
		DotEnv:  []string{".env", "app.env"},
		Extends: []string{"base", "common"},
	}

	if len(p.Env) != 3 {
		t.Errorf("Expected 3 env vars, got %d", len(p.Env))
	}

	if p.Env["KEY1"] != "value1" {
		t.Errorf("Expected KEY1=value1, got %v", p.Env["KEY1"])
	}

	if len(p.DotEnv) != 2 {
		t.Errorf("Expected 2 dotenv files, got %d", len(p.DotEnv))
	}

	if len(p.Extends) != 2 {
		t.Errorf("Expected 2 extends, got %d", len(p.Extends))
	}
}
