package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// OS commands to perform an action e.g. install DB from OS packages
type osCmds struct {
	id     string   // Holds distro + release e.g. ubuntu:18.04
	cmds   []string // Holds the os commands
	errmsg []string // Holds the error messages if the matching command fails
	hard   []bool   // Flag to know if an error on the matching command is fatal
}

// TODO: Document this and/or move it to a separate package
func sendCmd(d *DDConfig, o *log.Logger, cmd string, lerr string, hard bool) {
	// Setup command
	runCmd := exec.Command("bash", "-c", cmd)
	d.cmdLogger.Printf("[godojo] # %s\n", d.redactatron(cmd, d.redact))

	// Run and gather its output
	cmdOut, err := runCmd.CombinedOutput()
	if err != nil {
		d.errorMsg(fmt.Sprintf("%s - Failed to run OS command %+v, error was: %+v",
			timeStamp(), d.redactatron(cmd, d.redact), err))
		if hard {
			// Exit on hard aka fatal errors
			os.Exit(1)
		}
	}
	d.cmdLogger.Printf("%s\n", string(cmdOut))
	if err != nil {
		d.errorMsg(fmt.Sprintf("Failed to write to OS command log file, error was: %+v", err))
	}
}

// TODO: Document this and/or move it to a separate package
func tryCmd(d *DDConfig, cmd string, lerr string, hard bool) error {
	d.traceMsg("Entering tryCmd")
	// Setup command
	runCmd := exec.Command("bash", "-c", cmd)
	d.cmdLogger.Printf("[godojo] # " + d.redactatron(cmd, d.redact) + "\n")

	// Hook up stdout and strerr
	runCmd.Stdout = d.cmdLogger.Writer()
	runCmd.Stderr = d.cmdLogger.Writer()

	// Start the command
	err := runCmd.Start()
	if err != nil {
		d.traceMsg(fmt.Sprintf("Failed to start command, error was: %+v", err))
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
				d.traceMsg(fmt.Sprintf("%s - %s errored with exit status: %d", timeStamp(), cmd, status.ExitStatus()))
				return err
			}
		} else {
			d.traceMsg(fmt.Sprintf("%s - %s errored from Wait(): %v", timeStamp(), cmd, err))
			return err
		}
	}

	d.traceMsg("Non-error return from tryCmd")
	return nil
}

func tryCmds(d *DDConfig, c osCmds) error {
	// Cycle through the provided commands, trying them one at at time
	for i := range c.cmds {
		err := tryCmd(d,
			c.cmds[i],
			c.errmsg[i],
			c.hard[i])

		if err != nil {
			d.traceMsg(fmt.Sprintf("%s - Command %s errored with %s. Underlying error is %+v",
				timeStamp(), c.cmds[i], c.errmsg[i], err))
			return errors.New(c.errmsg[i])
		}
	}

	return nil
}

// TODO: Document this and/or move it to a separate package
func inspectCmd(d *DDConfig, cmd string, lerr string, hard bool) (string, error) {
	d.traceMsg("Inside inspectCmd")
	// Setup command
	runCmd := exec.Command("bash", "-c", cmd)
	d.cmdLogger.Printf("[godojo] # " + d.redactatron(cmd, d.redact) + "\n")
	//}

	// Hook up stdout and strerr
	var tmpBuf bytes.Buffer
	multi := io.MultiWriter(d.cmdLogger.Writer(), &tmpBuf)
	runCmd.Stdout = multi
	runCmd.Stderr = d.cmdLogger.Writer()

	// Start the command
	err := runCmd.Start()
	if err != nil {
		d.traceMsg(fmt.Sprintf("%s - Failed to start command %+v, error was: %+v",
			timeStamp(), d.redactatron(cmd, d.redact), err))
		return "", err
	}

	d.traceMsg("Before runCmd.Wait()")
	// Wait for command to exit, then check the exit code
	err = runCmd.Wait()
	if err != nil {
		// Check if the error is a ExitError
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				// Above casts the exiterr to syscll.WaitStatus aka unint32
				d.traceMsg(fmt.Sprintf("%s - %s errored with exit status: %d", timeStamp(), cmd, status.ExitStatus()))
				return "", err
			}
		} else {
			d.traceMsg(fmt.Sprintf("%s - %s errored from Wait(): %v", timeStamp(), cmd, err))
			return "", err
		}
	}
	d.traceMsg("After runCmd.Wait()")

	d.traceMsg("Non-error return from inspectCmd")
	return tmpBuf.String(), nil
}

func inspectCmds(d *DDConfig, c osCmds) ([]string, error) {
	d.traceMsg("Inside inspectCmds")
	ret := make([]string, 1)
	// Cycle through the provided commands, trying them one at at time
	for i := range c.cmds {
		d.traceMsg(fmt.Sprintf("Current cmd: %+v", c.cmds[i]))
		out, err := inspectCmd(d,
			c.cmds[i],
			c.errmsg[i],
			c.hard[i])

		if err != nil {
			d.traceMsg(fmt.Sprintf("%s - Command %s errored with %s. Underlying error is %+v",
				timeStamp(), c.cmds[i], c.errmsg[i], err))
			return ret, errors.New(c.errmsg[i])
		}
		ret = append(ret, out)
	}

	return ret, nil
}

func timeStamp() string {
	return time.Now().Format("2006/01/02 15:04:05")
}
