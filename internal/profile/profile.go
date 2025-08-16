package profile

// profile holds metadata plus an env-var map.
type profile struct {
	// Env is a collection of environment variables.
	Env Env
	// DotEnv is a list of dotenv files to import.
	DotEnv []string
	// Extends is a list of references to other profiles to extend from.
	Extends []string
	// Output is the desired output file.
	Output string
	// Default indicates whether this profile is the default one.
	Default bool
}

// newProfile creates a new profile with an empty env-var map.
func newProfile() *profile {
	return &profile{
		DotEnv:  []string{},
		Env:     make(Env),
		Extends: []string{},
	}
}
