package commandeer

import (
	"context"
)

// CmdLine is an interface for running OS commands either locally or remote
type Terminal interface {
	ExecCombined(ctx context.Context, cmd string, shell string) ([]byte, error)
	ExecError(ctx context.Context, cmd string, shell string) error
	ExecOnly(ctx context.Context, cmd string, shell string)
	ExecStdout(ctx context.Context, cmd string, shell string) ([]byte, error)
	ExecStderr(ctx context.Context, cmd string, shell string) ([]byte, error)
}
