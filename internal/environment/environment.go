package environment

import (
	"errors"
	"fmt"
	"strings"

	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Environment represents the resolved environment variables for a named profile.
type Environment struct {
	// Name is the name of the profile.
	Name string
	// Output is the output file for the profile.
	Output file.File

	// Env is the environment variables for the profile.
	Env env.Env
	// Origin tracks the source of each environment variable.
	Origin Origin
}

// New returns a new environment for the given profile,
// with default values for the output file.
func New(name, output string) Environment {
	if output == "" {
		output = name + ".env"
	}

	return Environment{
		Name:   name,
		Output: file.New(output),
		Env:    make(env.Env),
		Origin: make(Origin),
	}
}

// OverlayDotEnv overlays the environment variables from a .env file.
func (e *Environment) OverlayDotEnv(path, profile string) error {
	file := file.New(path)

	if !file.Exists() {
		return fmt.Errorf("dotenv file %q does not exist", path)
	}

	if file.IsDir() {
		return fmt.Errorf("dotenv file %q is a directory", path)
	}

	env, err := env.FromDotEnv(file.Path())
	if err != nil {
		return err
	}

	e.UpdateOrigin(profile, env)

	e.Origin.Add(path, env.Keys()...)

	e.Env = env.MergedWith(e.Env)

	return nil
}

// UpdateOrigin updates the origin of the environment variables.
func (e *Environment) UpdateOrigin(profile string, env env.Env) {
	if e.Name == profile {
		e.Origin.Clear(env.Keys()...)
	} else {
		e.Origin.Clear(env.Keys()...)
		e.Origin.Add(profile, env.Keys()...)
	}
}

// OverlayOther overlays the environment variables from another environment.
func (e *Environment) OverlayOther(other Environment) {
	env := other.Env
	profile := other.Name

	e.UpdateOrigin(profile, env)

	e.Env = other.Env.MergedWith(e.Env)
}

// Write saves the environment variables to a dotenv file.
func (e *Environment) Write() error {
	if e.Output == "" {
		return errors.New("no output file specified")
	}

	envs := e.Env.AsSlice()

	envs = append([]string{fmt.Sprintf("# Active profile: %q", e.Name)}, envs...)
	envs = append(envs, "")

	if err := e.Output.Write([]byte(strings.Join(envs, "\n"))); err != nil {
		return fmt.Errorf("writing to dotenv file %q: %w", e.Output, err)
	}

	return nil
}
