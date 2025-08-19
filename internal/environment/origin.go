package environment

import (
	"fmt"
	"strings"
)

// Origin tracks a key to it's source.
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

// Heritage represents the origin parts of a key.
type Heritage []string

// String returns the string representation of the origin parts.
func (h Heritage) String() string {
	result := make([]string, 0, len(h))

	for _, part := range h {
		result = append(result, fmt.Sprintf("%q", part))
	}

	return strings.Join(result, " -> ")
}
