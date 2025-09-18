package execx

import (
	"os/exec"
	"strings"

	"github.com/idelchi/envprof/pkg/terminal"
)

// Replace replaces the current process with command+args using env and dir.
// On success it never returns.
func Replace(command string, args, env []string, shell terminal.Shell) error {
	if !shell.Interactive() {
		path, err := exec.LookPath(command)
		if err != nil {
			return err
		}

		return replace(path, args, env)
	}

	cmd := strings.Join(append([]string{command}, args...), " ")

	switch shell.Type() {
	case terminal.Unix:
		args = []string{"-i", "-c", cmd}

	case terminal.Powershell:
		args = []string{"-NoExit", "-Command", cmd}
	case terminal.Cmd:
		args = []string{"/K", cmd}
	default:
		args = []string{"-c", cmd}
	}

	return replace(string(shell), args, env)
}
