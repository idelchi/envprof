// Package dotenv provides a simple helper to write profiles to dotenv files.
package dotenv

import (
	"fmt"
	"strings"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/path/file"
)

// DotEnv represents a dotenv file for a specific profile.
type DotEnv struct {
	Profile string
	Path    file.File
}

// WriteFrom writes the environment variables for the profile to the dotenv file.
func (d DotEnv) WriteFrom(profiles *profile.Profiles) error {
	vars, err := profiles.Environment(d.Profile)
	if err != nil {
		return err //nolint:wrapcheck // Error does not need additional wrapping.
	}

	envs := vars.Env.AsSlice()

	envs = append([]string{fmt.Sprintf("# Active profile: %q", d.Profile)}, envs...)
	envs = append(envs, "")

	if err := d.Path.Write([]byte(strings.Join(envs, "\n"))); err != nil {
		return fmt.Errorf("writing to dotenv file %q: %w", d.Path, err)
	}

	return nil
}
