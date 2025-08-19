// Package debug provides a basic debug printing utility.
package debug

import (
	"fmt"
	"os"

	"github.com/idelchi/godyl/pkg/pretty"
)

// DebugEnvVar is the environment variable that controls debug output.
const DebugEnvVar = "ENVPROF_DEBUG"

// Debug prints a debug message if the DebugEnvVar environment variable is set.
func Debug(format string, args ...any) {
	if os.Getenv(DebugEnvVar) != "" {
		fmt.Printf( //nolint:forbidigo // Debug package is meant for printing
			"DEBUG: "+format+"\n",
			args...)
	}
}

// Print prints a debug message if the DebugEnvVar environment variable is set.
func Print(a ...any) {
	if os.Getenv(DebugEnvVar) != "" {
		pretty.PrintYAML(a)
	}
}
