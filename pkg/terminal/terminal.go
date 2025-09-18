// Package terminal provides functionality to spawn a terminal with a specific environment.
package terminal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

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

// Shell represents a terminal shell.
type Shell string

// Type returns the type of shell.
func (s Shell) Type() Type {
	switch s {
	case "":
		return None
	case "cmd":
		return Cmd
	case "powershell", "pwsh":
		return Powershell
	default:
		return Unix
	}
}

// Interactive returns true if the shell is interactive.
func (s Shell) Interactive() bool {
	return s.Type() != None
}

// Type represents the type of shell.
type Type int

const (
	// Unix represents a Unix-like shell.
	Unix Type = iota
	// Powershell represents the Powershell or pwsh shell.
	Powershell
	// Cmd represents the Windows Command Prompt shell.
	Cmd
	// None represents no shell.
	None
)

// Current tries to determine the current terminal being used.
func Current() Shell {
	env := env.FromEnv()

	if shell := env.GetAny("SHELL", "STARSHIP_SHELL"); shell != "" {
		return Shell(shell)
	}

	switch runtime.GOOS {
	case "windows":
		switch {
		case env.Exists("PROMPT"):
			return Shell("cmd")
		case env.Exists("PSMODULEPATH"):
			PSModulePath := env.Get("PSMODULEPATH")
			if strings.Contains(PSModulePath, "microsoft.powershell") {
				return Shell("pwsh")
			}

			return Shell("powershell")
		default:
			return Shell("cmd")
		}

	default:
		return Shell("sh")
	}
}
