package profiles

import (
	"fmt"

	"github.com/idelchi/envprof/internal/environment"
)

// Environment returns a fully resolved environment for a specific profile.
func (p Profiles) Environment(name string, steps Steps) (environment.Environment, error) {
	cur, err := p.Get(name)
	if err != nil {
		return environment.Environment{}, err
	}

	out := environment.New(name, cur.Output)

	for _, step := range steps {
		switch step.Kind {
		case StepDotenv:
			if err := out.OverlayDotEnv(step.Name, step.Owner); err != nil {
				return out, fmt.Errorf("profile %q: dotenv %q: %w", step.Owner, step.Name, err)
			}
		case StepProfile:
			pr, err := p.Get(step.Name)
			if err != nil {
				return out, err
			}

			pe, err := pr.ToEnv(step.Name)
			if err != nil {
				return out, fmt.Errorf("stringify %q: %w", step.Name, err)
			}

			out.OverlayOther(pe)
		case StepOverlay:
			steps, err := p.Plan(step.Name)
			if err != nil {
				return out, fmt.Errorf("applying overlay %q: %w", step.Name, err)
			}
			e, err := p.Environment(step.Name, steps)
			if err != nil {
				return out, fmt.Errorf("applying overlay %q: %w", step.Name, err)
			}

			out.OverlayOther(e)
		}
	}

	// for _, ov := range overlays {
	// 	e, err := p.Environment(ov)
	// 	if err != nil {
	// 		return out, fmt.Errorf("applying overlay %q: %w", ov, err)
	// 	}

	// 	out.OverlayOther(e)
	// }

	return out, nil
}

// Environments returns a fully resolved list of environments for all profiles.
func (p Profiles) Environments() (environments []environment.Environment, err error) {
	for name := range p {
		steps, err := p.Plan(name)
		if err != nil {
			return nil, err
		}
		env, err := p.Environment(name, steps)
		if err != nil {
			return nil, err
		}

		environments = append(environments, env)
	}

	return environments, nil
}
