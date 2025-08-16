package profile

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
)

// unmarshal decodes the data into the store's profiles.
func (s *Store) unmarshal(data []byte) error {
	switch s.Type {
	case YAML:
		if err := yaml.UnmarshalWithOptions(data, &s.Profiles, yaml.Strict()); err != nil {
			return err
		}

	case TOML:
		md, err := toml.Decode(string(data), &s.Profiles)
		if err != nil {
			return err
		}

		if undecoded := md.Undecoded(); len(undecoded) > 0 {
			errs := make([]error, len(undecoded))
			for i, key := range undecoded {
				errs[i] = fmt.Errorf("unknown field: %s", key.String())
			}

			return errors.Join(errs...)
		}

	default:
		return fmt.Errorf("%w: %q", ErrUnsupportedFileType, s.Type)
	}

	return nil
}
