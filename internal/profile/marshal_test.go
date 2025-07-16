//nolint:testpackage // Tests need access to private unmarshal() method and types
package profile

import (
	"errors"
	"testing"
)

// TestStore_UnsupportedType tests that unmarshaling returns an error for unsupported file types.
func TestStore_UnsupportedType(t *testing.T) {
	t.Parallel()

	store := &Store{
		Type:     Type("json"),
		Profiles: Profiles{},
	}

	err := store.unmarshal([]byte(`{}`))
	if err == nil {
		t.Error("expected error for unsupported file type")
	}

	if !errors.Is(err, ErrUnsupportedFileType) {
		t.Errorf("expected ErrUnsupportedFileType, got %v", err)
	}
}
