package profile

import (
	"fmt"
	"slices"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/env"
)

// Env is a map of environment variable names to their non-stringified values.
type Env map[string]any

// UnmarshalYAML allows Env to be unmarshalled as its regular type or a sequence of strings.
func (e *Env) UnmarshalYAML(node ast.Node) error {
	if seq, ok := node.(*ast.SequenceNode); ok {
		var envs []string

		if err := yaml.NodeToValue(seq, &envs); err != nil {
			return fmt.Errorf("decoding env sequence: %w", err)
		}

		env, err := env.AsEnv(envs...)
		if err != nil {
			return fmt.Errorf("converting to env.Env: %w", err)
		}

		e.FromEnv(env)

		return nil
	}

	type raw Env

	if err := yaml.NodeToValue(node, (*raw)(e)); err != nil {
		return fmt.Errorf("decoding env map: %w", err)
	}

	return nil
}

// FromEnv initializes Env from an env.Env.
func (e *Env) FromEnv(env env.Env) {
	*e = make(Env)
	for k, v := range env {
		(*e)[k] = v
	}
}

// Stringified serializes Env into env.Env using Stringify.
// – Scalars pass through unchanged.
// – Non-scalars are JSON-minified and single-quoted (see Stringify).
func (e *Env) Stringified() (env.Env, error) {
	env := make(env.Env, len(*e))

	keys := make([]string, 0, len(*e))
	for k := range *e {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, key := range keys {
		str, err := Stringify((*e)[key])
		if err != nil {
			return nil, fmt.Errorf("profile %q: %w", key, err)
		}

		if err := env.AddPair(key, str); err != nil {
			return nil, fmt.Errorf("profile %q: %w", key, err)
		}
	}

	return env, nil
}
