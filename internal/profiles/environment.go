package profiles

import (
	"fmt"

	"github.com/idelchi/envprof/internal/environment"
	"github.com/idelchi/envprof/internal/step"
)

// Environment returns a fully resolved environment for a specific profile.
func (p Profiles) Environment(name string, steps step.Steps) (environment.Environment, error) {
	cur, err := p.Get(name)
	if err != nil {
		return environment.Environment{}, err
	}

	out := environment.New(name, cur.Output)

	for _, stp := range steps {
		switch stp.Kind {
		case step.DotEnv:
			if err := out.OverlayDotEnv(stp.Name, stp.Owner); err != nil {
				return out, fmt.Errorf("profile %q: dotenv %q: %w", stp.Owner, stp.Name, err)
			}
		case step.Profile:
			pr, err := p.Get(stp.Name)
			if err != nil {
				return out, err
			}

			pe, err := pr.ToEnv(stp.Name)
			if err != nil {
				return out, fmt.Errorf("stringify %q: %w", stp.Name, err)
			}

			out.OverlayOther(pe)
		case step.Overlay:
			steps, err := p.Plan(stp.Name)
			if err != nil {
				return out, fmt.Errorf("applying overlay %q: %w", stp.Name, err)
			}

			e, err := p.Environment(stp.Name, steps)
			if err != nil {
				return out, fmt.Errorf("applying overlay %q: %w", stp.Name, err)
			}

			out.OverlayOther(e)
		}
	}

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
