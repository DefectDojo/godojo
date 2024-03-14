package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	c "github.com/mtesauro/commandeer"
)

// godojo default value struct
type DDConfig struct {
	ver         string           // Holds the version of godojo
	cf          string           // Name of the config file
	conf        dojoConfig       // Global config struct
	sensStr     []string         // Holds sensitive strings to redact
	logLocation string           // Where the logs are written, relative to the directory godojo is called in
	Trace       *log.Logger      // Logger for trace logs
	Info        *log.Logger      // Logger for info logs
	Warning     *log.Logger      // Logger for warning logs
	Error       *log.Logger      // Logger for error logs
	cmdLogger   *log.Logger      // File pointer to the file in logLocation where command output is written
	helpURL     string           // Location of the godojo help URL
	releaseURL  string           // Location to download DefectDojo releases
	cloneURL    string           // URL to git clone DefectDojo
	yarnGPG     string           // URL to the yarn GPG key
	yarnRepo    string           // URL for the yarn repo
	nodeURL     string           // URL for the node repo
	quiet       bool             // Runtime flag to suppress output
	traceOn     bool             // Runtime flag to turn on trace logging
	redact      bool             // Runtime flag to redact sensitive info (defaults to on)
	spin        *spinner.Spinner // Progress spinner
	defInstall  bool             // Holds command-line bool asking for a default install
	emdir       string
	otdir       string
	bdir        string
	modf        string
	tgzf        string
}

// Set the godojo defaults in the DDConfig struct
func (d *DDConfig) setGodojoDefaults() {
	d.ver = "1.2.2"
	d.cf = "dojoConfig.yml"

	// Setup default logging
	d.logLocation = "logs"
	logHandler := d.prepLogging()
	d.Trace = log.New(logHandler, "TRACE:   ", log.Ldate|log.Ltime)
	d.Info = log.New(logHandler, "INFO:    ", log.Ldate|log.Ltime)
	d.Warning = log.New(logHandler, "WARNING: ", log.Ldate|log.Ltime)
	d.Error = log.New(logHandler, "ERROR:   ", log.Ldate|log.Ltime)

	// Set some installer defaults - .deb specific
	d.helpURL = "https://github.com/DefectDojo/godojo"
	d.releaseURL = "https://github.com/DefectDojo/django-DefectDojo/archive/"
	d.cloneURL = "https://github.com/DefectDojo/django-DefectDojo.git"
	d.yarnGPG = "https://dl.yarnpkg.com/debian/pubkey.gpg"
	d.yarnRepo = "deb https://dl.yarnpkg.com/debian/ stable main"
	d.nodeURL = "https://deb.nodesource.com/setup_18.x"
	d.quiet = false
	d.traceOn = true
	d.redact = true
	d.defInstall = false
	d.emdir = "embd/"
	d.otdir = "/tmp/.dojo-temp/"
	d.bdir = "/opt/"
	d.modf = ".dd.mod"
	d.tgzf = "gdj.tar.gz"

}

func (gd *DDConfig) prepLogging() io.Writer {
	// Setup logging for the installer
	n := time.Now()
	when := strconv.Itoa(int(n.UnixNano()))
	logName := "dojo-install_" + when + ".log"
	logPath := path.Join(gd.logLocation, logName)
	// Create the logs directory if it does not exist
	_, err := os.Stat(logPath)
	if err != nil {
		// logs directory doesn't exist
		err = os.MkdirAll(gd.logLocation, 0755)
		if err != nil {
			// Can't create logs directory for some reason, exit after showing error
			fmt.Println("")
			fmt.Println("##############################################################################")
			fmt.Printf("  Error creating godojo installer logging directory was %+v\n", err)
			fmt.Println("    Installation requires a logging directory.  Either create one in the same")
			fmt.Println("    directory as the godojo installer or correct the error above.")
			fmt.Println("##############################################################################")
			fmt.Println("")
			fmt.Println("Exiting install")
			os.Exit(1)
		}
	}

	// Create log file for the install
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("")
		fmt.Println("##############################################################################")
		fmt.Printf("  ERROR: Failed to open log file %s.  Error was:\n    %+v\n", logPath, err)
		fmt.Println("##############################################################################")
		fmt.Println("")
		fmt.Println("Log files are required for the install, exiting install")
		os.Exit(1)
	}

	// Return the logfile
	return logFile

}

// Output a section message and log the same string
func (gd *DDConfig) sectionMsg(s string) {
	// Pring status message if quiet isn't set
	if !gd.quiet {
		fmt.Println("")
		fmt.Println("==============================================================================")
		fmt.Printf("  %s\n", s)
		fmt.Println("==============================================================================")
		fmt.Println("")
	}
	gd.Info.Println("SECTION: " + s)
}

// Output a status message and log the same string
func (gd *DDConfig) statusMsg(s string) {
	// Pring status message if quiet isn't set & redact sensitive info in redact is true
	if !gd.quiet {
		fmt.Printf("%s\n", gd.redactatron(s, gd.redact))
	}
	gd.Info.Println(gd.redactatron(s, gd.redact))
}

// Output a blatant error message and log the string to the error log
func (gd *DDConfig) warnMsg(s string) {
	// Pring status message if quiet isn't set & redact sensitive info in redact is true
	if !gd.quiet {
		fmt.Println("")
		fmt.Println("##############################################################################")
		fmt.Printf("  WARNING: %s\n", gd.redactatron(s, gd.redact))
		fmt.Println("##############################################################################")
		fmt.Println("")
	}
	gd.Warning.Println(gd.redactatron(s, gd.redact))
}

// Output a blatant error message and log the string to the error log
func (gd *DDConfig) errorMsg(s string) {
	// Pring status message if quiet isn't set & redact sensitive info in redact is true
	if !gd.quiet {
		fmt.Println("")
		fmt.Println("##############################################################################")
		fmt.Printf("  ERROR: %s\n", gd.redactatron(s, gd.redact))
		fmt.Println("##############################################################################")
		fmt.Println("")
	}
	gd.Error.Println(gd.redactatron(s, gd.redact))
}

// Log the string as an trace log
func (gd *DDConfig) traceMsg(s string) {
	// Pring status message if quiet isn't set & redact sensitive info in redact is true
	if gd.traceOn {
		gd.Trace.Println(gd.redactatron(s, gd.redact))
	}
}

// Output the installer banner
func (gd *DDConfig) dojoBanner() {
	fmt.Println("        ____       ____          __     ____          _      ")
	fmt.Println("       / __ \\___  / __/__  _____/ /_   / __ \\____    (_)___  ")
	fmt.Println("      / / / / _ \\/ /_/ _ \\/ ___/ __/  / / / / __ \\  / / __ \\ ")
	fmt.Println("     / /_/ /  __/ __/  __/ /__/ /_   / /_/ / /_/ / / / /_/ / ")
	fmt.Println("    /_____/\\___/_/  \\___/\\___/\\__/  /_____/\\____/_/ /\\____/  ")
	fmt.Println("                                               /___/         ")
	fmt.Println("    version ", gd.ver)
	fmt.Println("")
	fmt.Println("  Welcome to godojo, the official way to install DefectDojo on iron.")
	fmt.Println("  For more information on how goDojo does an install, see:")
	fmt.Printf("  %s", gd.helpURL)
	fmt.Println("")
}

func (gd *DDConfig) getReplacements() map[string]string {
	// Setup values to replace
	iv := make(map[string]string)

	// Setup map to inject values for placholders
	iv["{yarnGPG}"] = gd.conf.Options.YarnGPG                      // Yarn's GPG key URL
	iv["{yarnRepo}"] = gd.conf.Options.YarnRepo                    // Yarn's package URL
	iv["{nodeURL}"] = gd.conf.Options.NodeURL                      // Node's URL
	iv["{conf.Install.Root}"] = gd.conf.Install.Root               // Path where DefectDojo is installed defaults to /opt/dojo
	iv["{conf.Install.OS.Group}"] = gd.conf.Install.OS.Group       // OS group used by DefectDojo application
	iv["{conf.Install.OS.User}"] = gd.conf.Install.OS.User         // OS user used by DefectDojo application
	iv["{conf.Install.Admin.User}"] = gd.conf.Install.Admin.User   // Admin user used by DefectDojo web UI
	iv["{conf.Install.Admin.Email}"] = gd.conf.Install.Admin.Email // Admin user's email address used by DefectDojo web UI
	iv["{conf.Install.Admin.Pass}"] = gd.conf.Install.Admin.Pass   // Admin user's password for DefectDojo web UI

	return iv
}

func (gd *DDConfig) injectConfigVals(cmds []c.SingleCmd) {
	// Get replacement values
	confVal := gd.getReplacements()

	// Cycle through commands, making replacements as needed
	for k := range cmds {
		for i, v := range confVal {
			// Check if a replacement is needed
			if strings.Contains(cmds[k].Cmd, i) {
				cmds[k].Cmd = strings.ReplaceAll(cmds[k].Cmd, i, v)
			}
		}
	}

}
