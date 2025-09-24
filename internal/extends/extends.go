package extends

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// Extend represents an extension entry that can reference profiles or dotenv files.
type Extend string

const (
	// Profile is an extend entry for a profile.
	Profile Extend = "profile"
	// DotEnv is an extend entry for a dotenv file.
	DotEnv Extend = "dotenv"
	// Invalid is an extend entry for an invalid type.
	Invalid Extend = "invalid"
)

// Type returns the type of the extend entry.
func (e Extend) Type() Extend {
	switch {
	case strings.HasPrefix(string(e), "profile:"):
		return Profile
	case strings.HasPrefix(string(e), "dotenv:"):
		return DotEnv
	case strings.Contains(string(e), ":"):
		//nolint:mnd // Selects the first segment before the colon.
		return Extend(strings.SplitN(string(e), ":", 2)[0])
	default:
		return Profile
	}
}

// Path strips the type prefix from the extend entry.
func (e Extend) Path() string {
	return strings.TrimPrefix(strings.TrimPrefix(string(e), string(e.Type())), ":")
}

// ToType converts a slice of strings to Extend entries with the specified type prefix.
func ToType(e []string, t Extend) []Extend {
	extends := make([]Extend, 0, len(e))

	for _, s := range e {
		extends = append(extends, Extend(fmt.Sprintf("%s:%s", t, s)))
	}

	return extends
}

// Extends represents a slice of extension entries.
type Extends []Extend

// Valid checks if all extend entries are valid.
func (es *Extends) Valid() error {
	for _, e := range *es {
		if e.Type() == Invalid {
			return errors.New("invalid extend: " + string(e))
		}
	}

	return nil
}

// Resolve expands glob patterns in dotenv extends and resolves all entries.
func (es *Extends) Resolve() error {
	var extends Extends

	for _, extend := range *es {
		if extend.Type() == DotEnv {
			path := extend.Path()

			matches, err := filepath.Glob(path)
			if err != nil {
				matches = []string{path}
			}

			if len(matches) == 0 {
				return fmt.Errorf("dotenv %q: no matches found", path)
			}

			extends = append(extends, ToType(matches, extend.Type())...)
		} else {
			extends = append(extends, extend)
		}
	}

	*es = extends

	return nil
}
