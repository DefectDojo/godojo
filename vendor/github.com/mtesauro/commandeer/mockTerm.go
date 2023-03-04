package commandeer

import (
	"context"
)

// Create a mock terminal to use while testing to control both sides of an execution

type MockTerm struct {
	Out         []byte
	Err         error
	LastCommand string
}

func (m *MockTerm) ExecCombined(cxt context.Context, cmd string, shell string) ([]byte, error) {
	// Store the last command executed
	m.LastCommand = cmd

	return m.Out, m.Err
}

func (m *MockTerm) ExecError(ctx context.Context, cmd string, shell string) error {
	// Store the last command executed
	m.LastCommand = cmd

	return m.Err
}

func (m *MockTerm) ExecOnly(ctx context.Context, cmd string, shell string) {
	// Store the last command executed
	m.LastCommand = cmd
}

func (m *MockTerm) ExecStdout(cxt context.Context, cmd string, shell string) ([]byte, error) {
	// Store the last command executed
	m.LastCommand = cmd

	return m.Out, m.Err
}

func (m *MockTerm) ExecStderr(cxt context.Context, cmd string, shell string) ([]byte, error) {
	// Store the last command executed
	m.LastCommand = cmd

	return m.Out, m.Err
}
