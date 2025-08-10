package profile

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Type represents the type of the profile file.
type Type string

const (
	// YAML is the YAML file type.
	YAML Type = "yaml"
	// TOML is the TOML file type.
	TOML Type = "toml"
)

// Store represents the entire file and codec information.
type Store struct {
	// File is the file containing the profiles.
	File file.File
	// Profiles is a map of profile names to their metadata.
	Profiles Profiles
	// Type is the type of the file.
	Type Type
}

var (
	// ErrUnsupportedFileType is returned when the file type is not supported.
	ErrUnsupportedFileType = errors.New("unsupported file type")
	// ErrProfileNotFound is returned when a profile is not found in the store.
	ErrProfileNotFound = errors.New("profile not found")
)

// New creates a new store from the given path.
func New(file file.File) (*Store, error) {
	store := &Store{
		File:     file,
		Profiles: map[string]*profile{},
	}

	switch ext := file.Extension(); ext {
	case "yaml", "yml":
		store.Type = YAML
	case "toml":
		store.Type = TOML
	default:
		return nil, fmt.Errorf("%w: %q: %q", ErrUnsupportedFileType, file.Path(), ext)
	}

	return store, nil
}

// Load reads the file and unmarshals it into the store.
func (s *Store) Load() (*Store, error) {
	data, err := s.File.Read()
	if err != nil {
		return nil, err //nolint:wrapcheck	// Error does not need additional wrapping.
	}

	if err = s.unmarshal(data); err != nil {
		return nil, err
	}

	// Make sure no nil entries exist in the profiles map.
	for name, profile := range s.Profiles {
		if profile == nil {
			s.Profiles.Create(name)
		} else if profile.Env == nil {
			profile.Env = make(Env)
			s.Profiles[name] = profile
		}
	}

	return s, err
}
