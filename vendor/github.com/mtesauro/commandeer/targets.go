package commandeer

import (
	"time"
)

// Holds a single command including the max exec time aka Timeout
type SingleCmd struct {
	Cmd        string        // Holds the command to be executed
	Errmsg     string        // Holds a custom error message to return on error
	Hard       bool          // Exit running if an error occurs during execution e.g. os.Exit(1)
	Timeout    time.Duration // Holds the max time a command can run before being cancelled
	BeforeText string        // Holds the optional text to print before starting the command
	AfterText  string        // Holds the optional text to print after starting the command
	Stdout     string        // Holds stdout from running the command
	Stderr     string        // Holds stderr from running the command
	StdBoth    string        // Holds both stdout and stderr from running the command
}

//////////////////////////////////
// Target for a command package //
//////////////////////////////////

// Targets are the targets a program using commandeer supports
type Target struct {
	ID      string      // Holds the supported Distro ID e.g. Ubuntu:21.10
	Distro  string      // Holds the supported Distro name e.g. Ubuntu
	Release string      // Holds the supported Release e.g. 21.10
	OS      string      // Holds the supported operating system e.g. Linux
	Shell   string      // Holds the supported shell to run the command in e.g. bash or sh
	PkgCmds []SingleCmd // Holds the commands in this command package
}
