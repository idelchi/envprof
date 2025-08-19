package environment

import (
	"fmt"
	"strings"
)

// Origin tracks environment variable keys to their sources.
type Origin map[string]Heritage

// Add adds a new origin for the given keys.
func (o *Origin) Add(origin string, keys ...string) {
	for _, k := range keys {
		(*o)[k] = append((*o)[k], origin)
	}
}

// Clear removes all origins for the given keys.
func (o *Origin) Clear(keys ...string) {
	for _, k := range keys {
		delete(*o, k)
	}
}

// Heritage represents the inheritance chain of an environment variable.
type Heritage []string

// String returns the string representation of the heritage chain.
func (h Heritage) String() string {
	result := make([]string, 0, len(h))

	for _, part := range h {
		result = append(result, fmt.Sprintf("%q", part))
	}

	return strings.Join(result, " -> ")
}
