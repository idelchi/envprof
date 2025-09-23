package envprof

import (
	"errors"
	"fmt"

	"github.com/idelchi/envprof/internal/profiles"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
)

// Type represents the type of the profile file.
type Type string

const (
	// YAML is the YAML file type.
	YAML Type = "yaml"
	// TOML is the TOML file type.
	TOML Type = "toml"
)

// EnvProf represents an environment profile file and its loaded content.
type EnvProf struct {
	file     file.File
	format   Type
	profiles profiles.Profiles // loaded profiles
}

// New creates a new EnvProf instance from the given file.
func New(file file.File) *EnvProf {
	return &EnvProf{file: file}
}

// NewFrom creates a new EnvProf instance from the first found among the given files.
func NewFrom(files files.Files) (*EnvProf, error) {
	files.Expanded()

	file, ok := files.Exists()
	if !ok {
		return nil, fmt.Errorf("profile file not found: searched for %v", files)
	}

	return New(file), nil
}

// File returns the resolved file.
func (e *EnvProf) File() file.File {
	return e.file
}

// Type determines and sets the file format based on the file extension.
func (e *EnvProf) Type() error {
	switch ext := e.file.Extension(); ext {
	case "yaml", "yml":
		e.format = YAML
	case "toml":
		e.format = TOML
	default:
		return fmt.Errorf("unsupported file extension: %q", ext)
	}

	return nil
}

// TryParse attempts to parse the given data into the supported profile formats.
func (e *EnvProf) TryParse(data []byte) error {
	types := []Type{YAML, TOML}

	for _, e.format = range types {
		if _, err := Unmarshal(data, e.format); err == nil {
			return nil
		}
	}

	return errors.New("format cannot be detected from content")
}

// GetOrDefault returns the profile name if it exists, or the default profile if none is specified.
func (e *EnvProf) GetOrDefault(name string) (string, error) {
	if name == "" {
		name = e.profiles.Default()
	}

	if name == "" {
		return "", errors.New("no default profile found and none specified")
	}

	return name, nil
}

// Profiles returns the loaded profiles.
func (e *EnvProf) Profiles() profiles.Profiles {
	return e.profiles
}

// Load reads the file and unmarshals it into the store.
func (e *EnvProf) Load() error {
	data, err := e.file.Read()
	if err != nil {
		return err
	}

	data, err = Template(data, env.FromEnv())
	if err != nil {
		return fmt.Errorf("templating profile file %q: %w", e.file.Path(), err)
	}

	if errType := e.Type(); errType != nil {
		if errParse := e.TryParse(data); errParse != nil {
			return fmt.Errorf(
				"parsing profile file %q: %w: %w",
				e.file.Path(),
				errType, errParse,
			)
		}
	}

	profiles, err := Unmarshal(data, e.format)
	if err != nil {
		return err
	}

	if err = profiles.Validate(); err != nil {
		return err
	}

	e.profiles = profiles

	return nil
}
