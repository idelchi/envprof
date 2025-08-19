//go:build !windows

package execx

import (
	"golang.org/x/sys/unix"
)

// replace replaces the current process with the specified command using unix.Exec.
// On success it never returns.
func replace(path string, args, env []string) error {
	argv := append([]string{path}, args...) // argv[0] required

	return unix.Exec(path, argv, env)
}
