package profile

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Stringify turns any Go value into the dotenv-ready string described earlier.
//   - Scalars → plain strings (`true`, `5432`, `foo`…).
//   - Slices / maps / structs → compact JSON wrapped in single quotes
//     (`'["a","b"]'`, `'{"k":"v"}'`).
func Stringify(v any) (string, error) {
	switch val := v.(type) {
	case nil:
		return "", nil

	case string:
		if needsQuotes(val) {
			return strconv.Quote(val), nil // adds double-quotes and escapes
		}

		return val, nil

	case bool:
		return strconv.FormatBool(val), nil

	// All signed ints
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val), nil

	// Unsigned ints
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val), nil

	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil

	default:
		// Anything not matched above gets JSON-encoded.
		raw, err := json.Marshal(val)
		if err != nil {
			return "", fmt.Errorf("json: %w", err)
		}

		return "'" + string(raw) + "'", nil
	}
}

// needsQuotes returns true if the string contains characters that
// require quoting in a dotenv file (space, =, #, quotes, etc.).
func needsQuotes(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		switch {
		case unicode.IsSpace(r):
			return true
		case strings.ContainsRune(`"'#=`, r):
			return true
		}
	}

	return false
}
