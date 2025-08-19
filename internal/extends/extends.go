package extends

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// Extend represents a extend entry.
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
		return Invalid
	default:
		return Profile
	}
}

// Path strips the type prefix from the extend entry.
func (e Extend) Path() string {
	return strings.TrimPrefix(strings.TrimPrefix(string(e), string(e.Type())), ":")
}

// ToType converts the extend entry to its type.
func ToType(e []string, t Extend) []Extend {
	var extends []Extend

	for _, s := range e {
		extends = append(extends, Extend(fmt.Sprintf("%s:%s", t, s)))
	}

	return extends
}

// Extends represents a list of extend entries.
type Extends []Extend

// Valid checks if all extend entries are valid.
func (es Extends) Valid() error {
	for _, e := range es {
		if e.Type() == Invalid {
			return errors.New("invalid extend: " + string(e))
		}
	}

	return nil
}

func (es *Extends) Resolve() error {
	var extends Extends

	for _, extend := range *es {
		if extend.Type() == DotEnv {
			path := extend.Path()

			matches, err := filepath.Glob(path)
			if err != nil {
				return fmt.Errorf("dotenv %q: %w", path, err)
			}

			extends = append(extends, ToType(matches, extend.Type())...)
		} else {
			extends = append(extends, extend)
		}
	}

	*es = extends

	return nil
}
