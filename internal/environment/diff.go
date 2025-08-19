package environment

import (
	"fmt"
	"slices"
	"strings"

	"github.com/idelchi/godyl/pkg/env"
)

// Change represents a single environment variable that has changed between two environments.
type Change struct {
	// Key is the environment variable name.
	Key string
	// Old is the previous value of the variable.
	Old string
	// New is the current value of the variable.
	New string
}

// Diff represents the differences between two environments.
type Diff struct {
	// Added contains variables present in B but not in A.
	Added env.Env
	// Removed contains variables present in A but not in B.
	Removed env.Env
	// Changed contains variables present in both with different values.
	Changed []Change
}

// Diffs computes a structured diff of two environments.
func Diffs(first, second env.Env) Diff {
	// Normalize for platform semantics (e.g., Windows case-insensitivity).
	first = first.Normalized()
	second = second.Normalized()

	out := Diff{
		Added:   make(env.Env),
		Removed: make(env.Env),
	}

	// Removed / Changed
	for _, key := range first.Keys() {
		oldValue := first.Get(key)
		if !second.Exists(key) {
			out.Removed[key] = oldValue

			continue
		}

		if bv := second.Get(key); oldValue != bv {
			out.Changed = append(out.Changed, Change{Key: key, Old: oldValue, New: bv})
		}
	}

	// Added
	for _, key := range second.Keys() {
		if !first.Exists(key) {
			out.Added[key] = second.Get(key)
		}
	}

	// Deterministic order for Changed.
	slices.SortFunc(out.Changed, func(x, y Change) int { return strings.Compare(x.Key, y.Key) })

	return out
}

// Equal reports whether there are no differences.
func (d Diff) Equal() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// Render prints a git-like unified diff.
// Lines: + added, - removed, ~ changed ("old -> new").
// aName/bName are labels (e.g., "env1", "env2").
//
//nolint:forbidigo	// Function prints out to the console.
func (d Diff) Render(aName, bName string) error {
	fmt.Printf("--- %s\n+++ %s\n", aName, bName)

	addKeys := d.Added.Keys()
	rmKeys := d.Removed.Keys()

	for _, k := range rmKeys {
		fmt.Printf("%s %s=%q\n", "-", k, d.Removed[k])
	}

	for _, k := range addKeys {
		fmt.Printf("%s %s=%q\n", "+", k, d.Added[k])
	}

	for _, ch := range d.Changed {
		fmt.Printf("%s %s: %q -> %q\n", "~", ch.Key, ch.Old, ch.New)
	}

	return nil
}
