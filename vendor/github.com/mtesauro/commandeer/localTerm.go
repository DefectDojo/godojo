package commandeer

import (
	"context"
	"io"
	"os/exec"
)

type LocalTerm struct {
}

// Exec the provided command and return the combined output (stdout and stderr) along with an error
func (l *LocalTerm) ExecCombined(ctx context.Context, cmd string, shell string) ([]byte, error) {
	// Wrap the command in the configured shell e.g. "sh -c 'command sent to method'"
	wrappedCmd := exec.CommandContext(ctx, shell, "-c", cmd)
	output, err := wrappedCmd.CombinedOutput()

	return output, err
}

// Exec the provided command and return only a Go error if any occurs. Command output will be silently dropped
func (l *LocalTerm) ExecError(ctx context.Context, cmd string, shell string) error {
	// Wrap the command in the configured shell e.g. "sh -c 'command sent to method'"
	wrappedCmd := exec.CommandContext(ctx, shell, "-c", cmd)
	_, err := wrappedCmd.CombinedOutput()

	return err
}

// Exec the provided command returning nothing. Command output and errors will be silently dropped
func (i *LocalTerm) ExecOnly(ctx context.Context, cmd string, shell string) {
	// Wrap the command in the configured shell e.g. "sh -c 'command sent to method'"
	wrappedCmd := exec.CommandContext(ctx, shell, "-c", cmd)
	// Silently drop any errors
	_ = wrappedCmd.Run()
}

// Exec the provided command and return only the contents of stdout. Any output to stderr will be silently dropped
func (l *LocalTerm) ExecStdout(ctx context.Context, cmd string, shell string) ([]byte, error) {
	// Wrap the command in the configured shell e.g. "sh -c 'command sent to method'"
	wrappedCmd := exec.CommandContext(ctx, shell, "-c", cmd)
	out, err := wrappedCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = wrappedCmd.Start()
	if err != nil {
		return nil, err
	}

	stdout, err := io.ReadAll(out)
	if err != nil {
		return stdout, err
	}

	err = wrappedCmd.Wait()
	if err != nil {
		return stdout, err
	}

	return stdout, err
}

// Exec the provided command and return only the contents of stderr. Any output to stdout will be silently dropped
func (l *LocalTerm) ExecStderr(ctx context.Context, cmd string, shell string) ([]byte, error) {
	// Wrap the command in the configured shell e.g. "sh -c 'command sent to method'"
	wrappedCmd := exec.CommandContext(ctx, shell, "-c", cmd)
	out, err := wrappedCmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = wrappedCmd.Start()
	if err != nil {
		return nil, err
	}

	stderr, err := io.ReadAll(out)
	if err != nil {
		return stderr, err
	}

	err = wrappedCmd.Wait()
	if err != nil {
		return stderr, err
	}

	return stderr, err
}
