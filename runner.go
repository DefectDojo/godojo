package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// TODO: Document this and/or move it to a separate package
func sendCmd(o io.Writer, cmd string, lerr string, hard bool) {
	// Setup command
	runCmd := exec.Command("bash", "-c", cmd)
	_, err := o.Write([]byte("[godojo] # " + Redactatron(cmd, Redact) + "\n"))
	if err != nil {
		errorMsg(fmt.Sprintf("Failed to setup command, error was: %+v", err))
	}

	// Run and gather its output
	cmdOut, err := runCmd.CombinedOutput()
	if err != nil {
		errorMsg(fmt.Sprintf("Failed to run OS command, error was: %+v", err))
		if hard {
			// Exit on hard aka fatal errors
			os.Exit(1)
		}
	}
	_, err = o.Write(cmdOut)
	if err != nil {
		errorMsg(fmt.Sprintf("Failed to write to OS command log file, error was: %+v", err))
	}
}

func runCmds(o io.Writer, c osCmds) {
	// Cycle through the provided commands, running them one at at time
	for i := range c.cmds {
		sendCmd(o,
			c.cmds[i],
			c.errmsg[i],
			c.hard[i])
	}
}

// TODO: Document this and/or move it to a separate package
func tryCmd(o io.Writer, cmd string, lerr string, hard bool) error {
	// Setup command
	runCmd := exec.Command("bash", "-c", cmd)
	_, err := o.Write([]byte("[godojo] # " + Redactatron(cmd, Redact) + "\n"))
	if err != nil {
		traceMsg(fmt.Sprintf("Failed to setup command, error was: %+v", err))
		return err
	}

	// Hook up stdout and strerr
	runCmd.Stdout = o
	runCmd.Stderr = o

	// Start the command
	err = runCmd.Start()
	if err != nil {
		traceMsg(fmt.Sprintf("Failed to start command, error was: %+v", err))
		return err
	}

	// Wait for command to exit, then check the exit code
	err = runCmd.Wait()
	if err != nil {
		// Check if the error is a ExitError
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				// Above casts the exiterr to syscll.WaitStatus aka unint32
				traceMsg(fmt.Sprintf("%s errored with exit status: %d", cmd, status.ExitStatus()))
				return err
			}
		} else {
			traceMsg(fmt.Sprintf("%s errored from Wait(): %v", cmd, err))
			return err
		}
	}

	return nil
}

func tryCmds(o io.Writer, c osCmds) error {
	// Cycle through the provided commands, trying them one at at time
	for i := range c.cmds {
		err := tryCmd(o,
			c.cmds[i],
			c.errmsg[i],
			c.hard[i])

		if err != nil {
			traceMsg(fmt.Sprintf("Command %s errored with %s. Underlying error is %+v", c.cmds[i], c.errmsg[i], err))
			return errors.New(c.errmsg[i])
		}
	}

	return nil
}
