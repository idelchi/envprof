package profiles

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/idelchi/envprof/internal/environment"
	"github.com/idelchi/envprof/internal/profile"
)

// Profiles is a map of profile names to their metadata.
type Profiles map[string]profile.Profile

// GetOrDefault returns the profile name if it exists, or the default profile if none is specified.
func (p Profiles) GetOrDefault(name string) (string, error) {
	if name == "" {
		name = p.Default()
	}

	if name == "" {
		return "", errors.New("no default profile found and none specified")
	}

	return name, nil
}

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
		errs = append(errs, fmt.Errorf("more than one default profile: %v", defaults))
	}

	return errors.Join(errs...)
}

// Environments returns a fully resolved list of environments for all profiles.
func (p Profiles) Environments() (environments []environment.Environment, err error) {
	for name := range p {
		env, err := p.Environment(name)
		if err != nil {
			return nil, err
		}

		environments = append(environments, env)
	}

	return environments, nil
}

// Environment returns a fully resolved environment for a specific profile.
func (p Profiles) Environment(name string, overlays ...string) (environment.Environment, error) {
	cur, err := p.Get(name)
	out := environment.New(name, cur.Output)

	if err != nil {
		return out, err
	}

	plan, err := p.Plan(name)
	if err != nil {
		return out, err
	}

	for _, s := range plan {
		switch s.Kind {
		case StepDotenv:
			if err := out.OverlayDotEnv(s.Name, s.Owner); err != nil {
				return out, fmt.Errorf("profile %q: dotenv %q: %w", s.Owner, s.Name, err)
			}
		case StepProfile:
			pr, err := p.Get(s.Name)
			if err != nil {
				return out, err
			}

			pe, err := pr.ToEnv(s.Name)
			if err != nil {
				return out, fmt.Errorf("stringify %q: %w", s.Name, err)
			}

			out.OverlayOther(pe)
		}
	}

	for _, ov := range overlays {
		e, err := p.Environment(ov)
		if err != nil {
			return out, fmt.Errorf("applying overlay %q: %w", ov, err)
		}

		out.OverlayOther(e)
	}

	return out, nil
}

// get attempts to retrieve a profile by name.
// Errors for empty or non-existing profile names.
func (p Profiles) Get(name string) (profile.Profile, error) {
	if name == "" {
		return profile.Profile{}, errors.New("empty profile name")
	}

	if !p.Exists(name) {
		return profile.Profile{}, fmt.Errorf("profile %q not found", name)
	}

	return p[name], nil
}
