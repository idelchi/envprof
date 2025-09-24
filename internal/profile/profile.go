package profile

import (
	"github.com/idelchi/envprof/internal/environment"
	"github.com/idelchi/envprof/internal/extends"
)

// Profile represents a configuration profile with environment variables and metadata.
type Profile struct {
	// Env is a collection of environment variables.
	Env Env `toml:"env,omitempty" yaml:"env,omitempty"`
	// Extends is a list of references to other places to extend from.
	Extends extends.Extends `toml:"extends,omitempty" yaml:"extends,omitempty"`
	// Output is the desired output file.
	Output string `toml:"output,omitempty" yaml:"output,omitempty"`
	// Default indicates whether this profile is the default one.
	Default bool `toml:"default,omitempty" yaml:"default,omitempty"`
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
