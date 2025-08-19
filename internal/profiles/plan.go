package profiles

import (
	"fmt"

	"github.com/idelchi/envprof/internal/extends"
	"github.com/idelchi/envprof/internal/step"
)

// Plan creates an execution plan for a profile with optional overlays.
func (p Profiles) Plan(root string, overlays ...string) (step.Steps, error) {
	steps, err := p.plan(root)
	if err != nil {
		return nil, err
	}

	for _, overlay := range overlays {
		steps = append(steps, step.Step{
			Kind:  step.Overlay,
			Owner: root,
			Name:  overlay,
		})
	}

	return steps, nil
}

// plan creates an execution plan for a single profile, handling inheritance.
func (p Profiles) plan(root string) (step.Steps, error) {
	if _, err := p.Get(root); err != nil {
		return nil, err
	}

	type state uint8

	const (
		visiting state = 1
		visited  state = 2
	)

	seen := map[string]state{}
	cache := map[string]step.Steps{}

	var visit func(string) (step.Steps, error)

	visit = func(node string) (step.Steps, error) {
		switch seen[node] {
		case visited:
			out := make(step.Steps, len(cache[node]))
			copy(out, cache[node])

			return out, nil
		case visiting:
			// unreachable if we detect cycles on edges below, but keep as fallback
			return nil, fmt.Errorf("cycle detected: %s", node)
		}

		seen[node] = visiting

		defer func() { seen[node] = visited }()

		profile, err := p.Get(node)
		if err != nil {
			return nil, err
		}

		var plan step.Steps

		for _, extend := range profile.Extends {
			switch extend.Type() {
			case extends.Profile:
				child := extend.Path()

				// Show only the two nodes involved in the back-edge.
				if seen[child] == visiting {
					return nil, fmt.Errorf("cycle detected: %s -> %s -> %s", node, child, node)
				}

				sub, err := visit(child)
				if err != nil {
					return nil, err
				}

				plan = append(plan, sub...) // parent (and its stuff) first

			case extends.DotEnv:
				plan = append(plan, step.Step{Kind: step.DotEnv, Owner: node, Name: extend.Path()}) // interleave

			case extends.Invalid:
				return nil, fmt.Errorf("profile %q: unsupported extends %q", node, extend.Type())
			}
		}

		plan = append(plan, step.Step{Kind: step.Profile, Owner: "", Name: node}) // inline env last

		cache[node] = plan

		out := make(step.Steps, len(plan))

		copy(out, plan)

		return out, nil
	}

	return visit(root)
}
