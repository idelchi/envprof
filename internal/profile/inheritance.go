package profile

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
)

// InheritanceTracker holds the environment variables and their inheritance sources.
type InheritanceTracker struct {
	// Inheritance is a mapping from environment variable names to where they are inherited from.
	Inheritance Inheritance
	// Env is the environment variables.
	Env env.Env
	// Name is the name of the profile.
	Name string
}

// Inheritance carries a mapping from environment variable names to where they are inherited from.
type Inheritance map[string]string

// Format returns the formatted value of a single variable.
func (i InheritanceTracker) Format(key string, verbose, withKey bool) string {
	val := i.Env.Get(key)
	if withKey {
		val = fmt.Sprintf("%v=%v", key, val)
	}

	if verbose {
		if src := i.Inheritance[key]; src != "" && src != i.Name {
			return fmt.Sprintf("%-60v (inherited from %q)", val, src)
		}
	}

	return val
}

// FormatAll returns all variables, one per line.
func (i InheritanceTracker) FormatAll(prefix string, verbose bool) string {
	out := []string{}

	for _, k := range i.Env.Keys() {
		s := i.Format(k, verbose, true)
		if prefix != "" {
			s = fmt.Sprintf("%s%s", prefix, s)
		}

		out = append(out, s)
	}

	return strings.Join(out, "\n")
}

// ToDotEnv writes the environment variables to a dotenv file.
func (i InheritanceTracker) ToDotEnv(dotenv file.File) error {
	envs := i.Env.AsSlice()

	envs = append([]string{fmt.Sprintf("# Active profile: %q", i.Name)}, envs...)
	envs = append(envs, "")

	if err := dotenv.Write([]byte(strings.Join(envs, "\n"))); err != nil {
		return fmt.Errorf("writing to dotenv file %q: %w", dotenv, err)
	}

	return nil
}
