package execx

import (
	"os/exec"
)

// Replace replaces the current process with command+args using env and dir.
// On success it never returns.
func Replace(command string, args, env []string) error {
	path, err := exec.LookPath(command)
	if err != nil {
		return err //nolint:wrapcheck	// Error does not need additional wrapping.
	}

	return replace(path, args, env)
}
