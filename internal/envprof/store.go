package envprof

import (
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

// File returns the path and type of the profile file selected.
func File(paths files.Files) (file.File, Type, error) {
	file, ok := paths.Exists()
	if !ok {
		return file, "", fmt.Errorf("profile file not found: searched for %v", paths)
	}

	var format Type

	switch ext := file.Extension(); ext {
	case "yaml", "yml":
		format = YAML
	case "toml":
		format = TOML
	default:
		return file, "", fmt.Errorf("unsupported file type: %q: %q", file.Path(), ext)
	}

	return file, format, nil
}

// Load reads the file and unmarshals it into the store.
func Load(paths files.Files) (profiles.Profiles, file.File, error) {
	file, format, err := File(paths)
	if err != nil {
		return nil, file, err
	}

	data, err := file.Read()
	if err != nil {
		return nil, file, err
	}

	data, err = Template(data, env.FromEnv())
	if err != nil {
		return nil, file, fmt.Errorf("templating profile file %q: %w", file.Path(), err)
	}

	profiles, err := Unmarshal(data, format)
	if err != nil {
		return nil, file, err
	}

	if err = profiles.Validate(); err != nil {
		return nil, file, err
	}

	for name, profile := range profiles {
		if err := profile.Extends.Resolve(); err != nil {
			return nil, file, fmt.Errorf("profile %q: %w", name, err)
		}

		profiles[name] = profile
	}

	return profiles, file, nil
}
