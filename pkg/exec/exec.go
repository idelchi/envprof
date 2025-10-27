package execx

import (
	"os/exec"
	"strings"

	"github.com/idelchi/envprof/pkg/terminal"

	"mvdan.cc/sh/v3/syntax"
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

	switch shell.Type() {
	case terminal.Unix:
		parts := make([]string, 0, len(args)+1)

		parts = append(parts, quoteArg(command))

		for _, arg := range args {
			parts = append(parts, quoteArg(arg))
		}

		args = []string{"-i", "-c", strings.Join(parts, " ")}
	case terminal.Powershell:
		cmd := strings.Join(append([]string{command}, args...), " ")

		args = []string{"-NoExit", "-Command", cmd}
	case terminal.Cmd:
		cmd := strings.Join(append([]string{command}, args...), " ")

		args = []string{"/K", cmd}
	default:
		cmd := strings.Join(append([]string{command}, args...), " ")

		args = []string{"-c", cmd}
	}

	return replace(string(shell), args, env)
}

// quoteArg quotes a string for safe use in shell commands.
func quoteArg(s string) string {
	quoted, err := syntax.Quote(s, syntax.LangBash)
	if err != nil {
		// Fallback for edge cases (e.g., null bytes)
		return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
	}

	return quoted
}
