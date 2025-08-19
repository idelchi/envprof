package environment

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/idelchi/godyl/pkg/env"
)

type Change struct {
	Key string
	Old string
	New string
}

type Diff struct {
	Added   env.Env  // present in B, not in A
	Removed env.Env  // present in A, not in B
	Changed []Change // present in both, different values
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
	for _, k := range first.Keys() {
		av := first.Get(k)
		if !second.Exists(k) {
			out.Removed[k] = av

			continue
		}

		if bv := second.Get(k); av != bv {
			out.Changed = append(out.Changed, Change{Key: k, Old: av, New: bv})
		}
	}

	// Added
	for _, k := range second.Keys() {
		if !first.Exists(k) {
			out.Added[k] = second.Get(k)
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

// RenderUnified prints a git-like unified diff.
// Lines: + added, - removed, ~ changed ("old -> new").
// aName/bName are labels (e.g., "env1", "env2").
func (d Diff) RenderUnified(w io.Writer, aName, bName string) error {
	fmt.Fprintf(w, "--- %s\n+++ %s\n", aName, bName)

	addKeys := d.Added.Keys()
	rmKeys := d.Removed.Keys()

	for _, k := range rmKeys {
		fmt.Fprintf(w, "%s %s=%q\n", "-", k, d.Removed[k])
	}

	for _, k := range addKeys {
		fmt.Fprintf(w, "%s %s=%q\n", "+", k, d.Added[k])
	}

	for _, ch := range d.Changed {
		fmt.Fprintf(w, "%s %s: %q -> %q\n", "~", ch.Key, ch.Old, ch.New)
	}

	return nil
}
