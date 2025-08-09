// Package terminal provides functionality to spawn a terminal with a specific environment.
package terminal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/idelchi/godyl/pkg/env"
)

// Spawn launches a new shell with the specified environment variables.
func Spawn(shell string, env []string) error {
	cmd := exec.CommandContext(context.Background(), shell)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("spawning terminal %q: %w", shell, err)
	}

	return nil
}

// Current tries to determine the current terminal being used.
func Current() string {
	env := env.FromEnv()

	if shell := env.GetAny("SHELL", "STARSHIP_SHELL"); shell != "" {
		return shell
	}

	switch runtime.GOOS {
	case "windows":
		switch {
		case env.Exists("PROMPT"):
			return "cmd.exe"
		case env.Exists("PSMODULEPATH"):
			return "powershell.exe"
		default:
			return "cmd.exe"
		}
	default:
		return "sh"
	}
}
