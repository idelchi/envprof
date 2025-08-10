package profile

import (
	"github.com/idelchi/godyl/pkg/path/file"
)

// DotEnv represents a dotenv file for a specific profile.
type DotEnv struct {
	Profile string
	Path    file.File
}
