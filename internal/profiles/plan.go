package profiles

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/idelchi/envprof/internal/extends"
)

// replay steps
type StepKind string

const (
	StepDotenv  StepKind = "dotenv"
	StepProfile StepKind = "env"
	StepOverlay StepKind = "overlay"
)

type Step struct {
	Kind  StepKind
	Owner string // for dotenv origin
	Name  string // profile name or dotenv path
}

type Steps []Step

func (s Steps) Table() string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 4, 2, ' ', 0)

	fmt.Fprintln(w, "STEP\tPROFILE\tKIND\tNAME")
	for i, st := range s {
		if st.Owner == "" {
			st.Owner = st.Name // for env steps, owner == profile
			st.Name = ""       // no name for env steps
		}
		fmt.Fprintf(w, "%02d\t%s\t%s\t%s\n", i+1, st.Owner, st.Kind, st.Name)
	}
	_ = w.Flush()
	return b.String()
}

func (p Profiles) Plan(root string, overlays ...string) (Steps, error) {
	steps, err := p.plan(root)
	if err != nil {
		return nil, err
	}

	for _, overlay := range overlays {
		steps = append(steps, Step{
			Kind:  StepOverlay,
			Owner: root,
			Name:  overlay,
		})
	}

	return steps, nil
}

func (p Profiles) plan(root string) (Steps, error) {
	if _, err := p.Get(root); err != nil {
		return nil, err
	}

	type state uint8
	const (
		visiting state = 1
		visited  state = 2
	)

	seen := map[string]state{}
	cache := map[string]Steps{}

	var visit func(string) (Steps, error)
	visit = func(n string) (Steps, error) {
		switch seen[n] {
		case visited:
			out := make(Steps, len(cache[n]))
			copy(out, cache[n])
			return out, nil
		case visiting:
			// unreachable if we detect cycles on edges below, but keep as fallback
			return nil, fmt.Errorf("cycle detected: %s", n)
		}

		seen[n] = visiting
		defer func() { seen[n] = visited }()

		pr, err := p.Get(n)
		if err != nil {
			return nil, err
		}

		var plan Steps
		for _, e := range pr.Extends {
			switch e.Type() {
			case extends.Profile:
				child := e.Path()

				// Show only the two nodes involved in the back-edge.
				if seen[child] == visiting {
					return nil, fmt.Errorf("cycle detected: %s -> %s -> %s", n, child, n)
				}

				sub, err := visit(child)
				if err != nil {
					return nil, err
				}
				plan = append(plan, sub...) // parent (and its stuff) first

			case extends.DotEnv:
				plan = append(plan, Step{StepDotenv, n, e.Path()}) // interleave

			default:
				return nil, fmt.Errorf("profile %q: unsupported extends %q", n, e.Type())
			}
		}

		plan = append(plan, Step{StepProfile, "", n}) // inline env last

		cache[n] = plan
		out := make(Steps, len(plan))
		copy(out, plan)
		return out, nil
	}

	return visit(root)
}
