package profile

import (
	"errors"
	"fmt"
	"slices"

	"github.com/idelchi/godyl/pkg/dag"
	"github.com/idelchi/godyl/pkg/env"
)

// Profiles is a map of profile names to their metadata.
type Profiles map[string]*profile

// Exists checks if a profile exists in the store.
func (p Profiles) Exists(name string) bool {
	_, ok := p[name]

	return ok
}

// Create creates a new profile in the store.
func (p Profiles) Create(name string) {
	p[name] = newProfile()
}

// Names returns the names of the profiles in sorted order.
func (p Profiles) Names() []string {
	names := make([]string, 0, len(p))

	for k := range p {
		names = append(names, k)
	}

	slices.Sort(names)

	return names
}

// Environment returns the merged environment variables for a profile, resolving dependencies.
//
//nolint:gocognit	// TODO(Idelchi): Refactor this function to reduce cognitive complexity.
func (p Profiles) Environment(name string) (*InheritanceTracker, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: empty profile name", ErrProfileNotFound)
	}

	if !p.Exists(name) {
		return nil, fmt.Errorf("%w: %q", ErrProfileNotFound, name)
	}

	// Build dependency DAG.
	nodes := make([]string, 0, len(p))
	for n := range p {
		nodes = append(nodes, n)
	}

	g, err := dag.Build(nodes, func(n string) []string { return p[n].Extends })
	if err != nil {
		return nil, fmt.Errorf("dag: %w", err)
	}

	chain, err := g.Chain(name)
	if err != nil {
		return nil, fmt.Errorf("chain: %w", err)
	}

	final := &InheritanceTracker{
		Name:        name,
		Env:         make(env.Env),
		Inheritance: make(Inheritance),
	}

	// Check if there's some dotenv files to load and load them first.
	profile := p[name]
	if len(profile.DotEnv) > 0 {
		for _, file := range profile.DotEnv {
			dotenv, err := env.FromDotEnv(file)
			if err != nil {
				return nil, fmt.Errorf("profile %q: dotenv %q: %w", name, file, err)
			}

			for key, value := range dotenv {
				if err := final.Env.AddPair(key, value); err != nil {
					return nil, fmt.Errorf("profile %q: dotenv %q: %w", name, file, err)
				}

				final.Inheritance[key] = file
			}
		}
	}

	errs := []error{}

	for _, profile := range chain {
		m := p[profile].Env

		stringified, err := m.Stringified()
		if err != nil {
			return nil, fmt.Errorf("stringify profile %s: %w", profile, err)
		}

		for k, v := range stringified {
			final.Inheritance[k] = profile
			errs = append(errs, final.Env.AddPair(k, v))
		}
	}

	return final, errors.Join(errs...)
}
