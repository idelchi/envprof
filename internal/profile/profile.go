package profile

// profile holds metadata plus an env-var map.
type profile struct {
	Env     Env      `toml:"env,omitempty"     yaml:"env,omitempty"`
	DotEnv  []string `toml:"dotenv,omitempty"  yaml:"dotenv,omitempty"`
	Extends []string `toml:"extends,omitempty" yaml:"extends,omitempty"`
}

// newProfile creates a new profile with an empty env-var map.
func newProfile() *profile {
	return &profile{
		DotEnv:  []string{},
		Env:     make(Env),
		Extends: []string{},
	}
}
