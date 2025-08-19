package envprof

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"

	"github.com/idelchi/envprof/internal/profiles"
)

// Unmarshal decodes the data into profiles.
func Unmarshal(data []byte, format Type) (profiles profiles.Profiles, err error) {
	switch format {
	case YAML:
		if err := yaml.UnmarshalWithOptions(data, &profiles, yaml.Strict()); err != nil {
			return profiles, err
		}

	case TOML:
		md, err := toml.Decode(string(data), &profiles)
		if err != nil {
			return profiles, err
		}

		if undecoded := md.Undecoded(); len(undecoded) > 0 {
			errs := make([]error, len(undecoded))
			for i, key := range undecoded {
				errs[i] = fmt.Errorf("unknown field: %s", key.String())
			}

			return nil, errors.Join(errs...)
		}

	default:
		return nil, fmt.Errorf("unsupported file format: %q", format)
	}

	return profiles, nil
}
