package profile

import (
	"github.com/idelchi/envprof/internal/environment"
	"github.com/idelchi/envprof/internal/extends"
)

// Profile represents a configuration profile with environment variables and metadata.
type Profile struct {
	// Env is a collection of environment variables.
	Env Env
	// Extends is a list of references to other places to extend from.
	Extends extends.Extends
	// Output is the desired output file.
	Output string
	// Default indicates whether this profile is the default one.
	Default bool
}

// ToEnv converts the profile to an environment representation,
// stringifying the environment variables.
func (p *Profile) ToEnv(name string) (environment.Environment, error) {
	stringified, err := p.Env.Stringified()
	if err != nil {
		return environment.Environment{}, err
	}

	return environment.Environment{
		Name: name,
		Env:  stringified,
	}, nil
}
