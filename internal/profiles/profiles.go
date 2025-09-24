package profiles

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/idelchi/envprof/internal/profile"
)

// ErrValidation is returned when profiles validation fails.
var ErrValidation = errors.New("validation error")

// Profiles is a map of profile names to their metadata.
type Profiles map[string]profile.Profile

// Exists checks if a profile exists in the profiles.
func (p Profiles) Exists(name string) bool {
	_, ok := p[name]

	return ok
}

// Names returns the names of the profiles in sorted order.
func (p Profiles) Names() []string {
	names := slices.Collect(maps.Keys(p))

	slices.Sort(names)

	return names
}

// Defaults returns the names of the default profiles.
func (p Profiles) Defaults() (defaults []string) {
	for name, profile := range p {
		if profile.Default {
			defaults = append(defaults, name)
		}
	}

	return defaults
}

// Default returns the name of the first default profile found.
func (p Profiles) Default() string {
	defaults := p.Defaults()
	if len(defaults) == 0 {
		return ""
	}

	return defaults[0]
}

// Validate checks that the profiles are valid.
func (p Profiles) Validate() error {
	var errs []error

	defaults := p.Defaults()

	if len(defaults) > 1 {
		errs = append(errs, fmt.Errorf("%w: more than one default profile: %v", ErrValidation, defaults))
	}

	for name, profile := range p {
		if err := profile.Extends.Resolve(); err != nil {
			return fmt.Errorf("%w: profile %q: %w", ErrValidation, name, err)
		}

		p[name] = profile
	}

	return errors.Join(errs...)
}

// Get retrieves a profile by name.
// Returns an error for empty or non-existing profile names.
func (p Profiles) Get(name string) (profile.Profile, error) {
	if name == "" {
		return profile.Profile{}, errors.New("empty profile name")
	}

	if !p.Exists(name) {
		return profile.Profile{}, fmt.Errorf("profile %q not found", name)
	}

	return p[name], nil
}
