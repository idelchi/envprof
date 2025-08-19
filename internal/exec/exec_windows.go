//go:build windows

package execx

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/idelchi/godyl/pkg/path/file"
)

// replace simulates process replacement on Windows by running the command and exiting with its code.
func replace(path string, args, env []string) error {
	var command *exec.Cmd

	ext := strings.ToLower(file.New(path).Extension())
	if ext == "bat" || ext == "cmd" {
		//nolint:gosec	// The user can execute whatever they'd like.
		command = exec.CommandContext(
			context.Background(),
			"cmd.exe",
			append([]string{"/c", path}, args...)...)
	} else {
		command = exec.CommandContext(context.Background(), path, args...)
	}

	command.Env = env

	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := command.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			os.Exit(ee.ExitCode()) //nolint:forbidigo // Allowing exit code propagation.
		}

		return err
	}

	if ps := command.ProcessState; ps != nil {
		os.Exit(ps.ExitCode()) //nolint:forbidigo	// Allowing exit code propagation.
	}

	os.Exit(0) //nolint:forbidigo	// Allowing exit code propagation.

	return nil
}
