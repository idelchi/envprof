package environment

import (
	"fmt"
	"strings"
)

// Formatter holds the configuration for formatting environment variables.
type Formatter struct {
	// WithOrigin indicates whether to include the origin of the variable.
	WithOrigin bool
	// WithKey indicates whether to include the key in the output.
	WithKey bool
	// Prefix is a string to prepend to the variable value.
	Prefix string
	// Padding is the width of the output field.
	Padding int
}

// Key formats an environment variable for output.
func (f Formatter) Key(key string, environment Environment) string {
	if f.Padding == 0 {
		f.Padding = 60
	}

	val := environment.Env.Get(key)

	if f.WithKey {
		val = fmt.Sprintf("%v=%v", key, val)
	}

	if f.WithOrigin {
		if src := environment.Origin[key]; len(src) > 0 {
			return fmt.Sprintf("%-*v (inherited from %s)", f.Padding, val, src)
		}
	}

	if f.Prefix != "" {
		val = f.Prefix + val
	}

	return val
}

// All formats all environment variables for output.
func (f Formatter) All(environment Environment) string {
	out := []string{}

	for _, k := range environment.Env.Keys() {
		out = append(out, f.Key(k, environment))
	}

	return strings.Join(out, "\n")
}
