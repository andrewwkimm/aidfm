package desktop

import (
	"fmt"
	"strings"
)

// ExecLine represents a parsed Exec= value with optional env vars and a binary path.
type ExecLine struct {
	Env    map[string]string
	Binary string
}

// ParseExec parses an Exec= value into its env vars and binary path.
// Handles both plain binary paths and env-prefixed forms:
//
//	Exec=/path/to/binary
//	Exec=env KEY=VALUE /path/to/binary
func ParseExec(exec string) ExecLine {
	parts := strings.Fields(exec)
	if len(parts) == 0 {
		return ExecLine{Env: make(map[string]string)}
	}

	result := ExecLine{Env: make(map[string]string)}
	start := 0

	if parts[0] == "env" {
		start = 1
	}

	for i := start; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "/") {
			result.Binary = parts[i]
			break
		}
		kv := strings.SplitN(parts[i], "=", 2)
		if len(kv) == 2 {
			result.Env[kv[0]] = kv[1]
		}
	}

	return result
}

// String reconstructs the Exec= value from the ExecLine.
// Produces "env KEY=VALUE /path/to/binary" if env vars are present,
// or "/path/to/binary" if not.
func (e ExecLine) String() string {
	if len(e.Env) == 0 {
		return e.Binary
	}

	parts := []string{"env"}
	for k, v := range e.Env {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	parts = append(parts, e.Binary)
	return strings.Join(parts, " ")
}

// SetEnv adds or updates an env var on the ExecLine.
func (e *ExecLine) SetEnv(key, value string) {
	e.Env[key] = value
}

// UnsetEnv removes an env var from the ExecLine.
func (e *ExecLine) UnsetEnv(key string) {
	delete(e.Env, key)
}
