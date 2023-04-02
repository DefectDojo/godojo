package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

func run(d *DDConfig) {
	// Print the install banner
	if !(d.quiet || d.conf.Options.Embd) {
		d.dojoBanner()
	}

	// Setup command logging
	d.cmdLogger = setCmdLogging(d)

	// Check embedded
	embdCk(d)

	// Check install OS
	osTarget := checkOS(d)

	// Bootstrap install
	bootstrapInstall(d, &osTarget)

	// Validate Python version
	validPython(d)

	// Download DefectDojo release or source
	downloadDojo(d)

	// Install OS packges need by DefectDojo
	prepOSForDojo(d, &osTarget)

	// Install DB if needed
	installDBForDojo(d, &osTarget)

	// Prepare the DB for DefectDojo
	prepDBForDojo(d, &osTarget)

	// Prepare for Django - virtenv, etc
	// TODO Convert to Commandeer
	prepDjango(d, &osTarget)

	// Create settings.py
	createSettings(d, &osTarget)

	// Setup DefectDojo
	setupDefectDojo(d, &osTarget)

	d.statusMsg(fmt.Sprintf("\nSuccessfully installed DefectDojo using godojo version %+v", d.ver))
}

func setCmdLogging(d *DDConfig) *log.Logger {
	// Setup OS command logging
	d.traceMsg("Creating log file for OS command output for debugging reasons")
	n := time.Now()
	when := strconv.Itoa(int(n.UnixNano()))
	cmdLog := "cmd-output_" + when + ".log"
	cmdPath := path.Join(d.logLocation, cmdLog)
	// Create command output log file in the existing logging directory
	cmdLogger, err := os.OpenFile(cmdPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("")
		fmt.Println("##############################################################################")
		fmt.Printf("  ERROR: Failed to open OS Command log file %s.  Error was:\n    %+v\n", cmdPath, err)
		fmt.Println("##############################################################################")
		fmt.Println("")
		fmt.Println("Log files are required for the install, exiting install")
		os.Exit(1)
	}
	//cmdLogger = cmdFile
	d.traceMsg(fmt.Sprintf("Successfully created OS Command log file at %+v", cmdPath))

	return log.New(cmdLogger, "[godojo] # ", log.Ldate|log.Ltime)
}
