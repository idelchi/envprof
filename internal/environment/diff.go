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

// DiffEnvs computes a structured diff of two environments.
func DiffEnvs(a, b env.Env) Diff {
	// Normalize for platform semantics (e.g., Windows case-insensitivity).
	a = a.Normalized()
	b = b.Normalized()

	out := Diff{
		Added:   make(env.Env),
		Removed: make(env.Env),
	}

	// Removed / Changed
	for _, k := range a.Keys() {
		av := a.Get(k)
		if !b.Exists(k) {
			out.Removed[k] = av
			continue
		}
		if bv := b.Get(k); av != bv {
			out.Changed = append(out.Changed, Change{Key: k, Old: av, New: bv})
		}
	}

	// Added
	for _, k := range b.Keys() {
		if !a.Exists(k) {
			out.Added[k] = b.Get(k)
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

type RenderOptions struct {
	Color bool // ANSI colors
}

// RenderUnified prints a git-like unified diff.
// Lines: + added, - removed, ~ changed ("old -> new").
// aName/bName are labels (e.g., "env1", "env2").
func (d Diff) RenderUnified(w io.Writer, aName, bName string, opts RenderOptions) error {
	color := func(code, s string) string {
		if !opts.Color {
			return s
		}
		return code + s + "\x1b[0m"
	}
	green := func(s string) string { return color("\x1b[32m", s) }
	red := func(s string) string { return color("\x1b[31m", s) }
	yellow := func(s string) string { return color("\x1b[33m", s) }

	fmt.Fprintf(w, "--- %s\n+++ %s\n", aName, bName)

	addKeys := d.Added.Keys()
	rmKeys := d.Removed.Keys()

	for _, k := range rmKeys {
		fmt.Fprintf(w, "%s %s=%q\n", red("-"), k, d.Removed[k])
	}
	for _, k := range addKeys {
		fmt.Fprintf(w, "%s %s=%q\n", green("+"), k, d.Added[k])
	}
	for _, ch := range d.Changed {
		fmt.Fprintf(w, "%s %s: %q -> %q\n", yellow("~"), ch.Key, ch.Old, ch.New)
	}
	return nil
}
