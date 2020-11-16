package main

// TODO: Consider Cobra for command-line args - https://github.com/spf13/cobra
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mtesauro/godojo/config"
	"github.com/spf13/viper"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Global vars
var (
	// Installer version
	ver = "1.1.1"
	// Configuration file name
	cf = "dojoConfig.yml"
	// Global config struct
	conf    config.DojoConfig
	sensStr [12]string // Hold sensitive strings to redact
	emdir   = "embd/"
	otdir   = "/tmp/.dojo-temp/"
	bdir    = "/opt/"
	modf    = ".dd.mod"
	tgzf    = "gdj.tar.gz"
	// For logging
	logLocation = "logs"
	Trace       *log.Logger
	Info        *log.Logger
	Warning     *log.Logger
	Error       *log.Logger
	// For Global config flags
	Quiet   bool
	TraceOn bool
	Redact  bool
	// Spinner FTW
	Spin spinner.Spinner
)

// Global Constants
const (
	// URLs needed by the installer
	// TODO: Move most of these into dojoConfig.yml optinal section
	HelpURL    = "https://github.com/mtesauro/godojo"
	ReleaseURL = "https://github.com/DefectDojo/django-DefectDojo/archive/"
	CloneURL   = "https://github.com/DefectDojo/django-DefectDojo.git"
	YarnGPG    = "https://dl.yarnpkg.com/debian/pubkey.gpg"
	YarnRepo   = "deb https://dl.yarnpkg.com/debian/ stable main"
	NodeURL    = "https://deb.nodesource.com/setup_12.x"
)

// Setup logging with type appended to the log lines - this logs all types to a single file
func logSetup(logHandler io.Writer) {
	// Setup logging 'levels' which can be called globally like Info.Println("Example info log")
	Trace = log.New(logHandler, "TRACE:   ", log.Ldate|log.Ltime)
	Info = log.New(logHandler, "INFO:    ", log.Ldate|log.Ltime)
	Warning = log.New(logHandler, "WARNING: ", log.Ldate|log.Ltime)
	Error = log.New(logHandler, "ERROR:   ", log.Ldate|log.Ltime)
}

// Output the installer banner
func dojoBanner() {
	fmt.Println("        ____       ____          __     ____          _      ")
	fmt.Println("       / __ \\___  / __/__  _____/ /_   / __ \\____    (_)___  ")
	fmt.Println("      / / / / _ \\/ /_/ _ \\/ ___/ __/  / / / / __ \\  / / __ \\ ")
	fmt.Println("     / /_/ /  __/ __/  __/ /__/ /_   / /_/ / /_/ / / / /_/ / ")
	fmt.Println("    /_____/\\___/_/  \\___/\\___/\\__/  /_____/\\____/_/ /\\____/  ")
	fmt.Println("                                               /___/         ")
	fmt.Println("    version ", ver)
	fmt.Println("")
	fmt.Println("  Welcome to goDojo, the official way to install DefectDojo.")
	fmt.Println("  For more information on how goDojo does an install, see:")
	fmt.Printf("  %s", HelpURL)
	fmt.Println("")
}

// Output a section message and log the same string
func sectionMsg(s string) {
	// Pring status message if quiet isn't set
	if !Quiet {
		fmt.Println("")
		fmt.Println("==============================================================================")
		fmt.Printf("  %s\n", s)
		fmt.Println("==============================================================================")
		fmt.Println("")
	}
	Info.Println("SECTION: " + s)
}

// Output a status message and log the same string
func statusMsg(s string) {
	// Redact sensitive info in redact is true
	Redactatron(s, Redact)
	// Pring status message if quiet isn't set
	if !Quiet {
		fmt.Printf("%s\n", s)
	}
	Info.Println(s)
}

// Output a blatant error message and log the string as an error
func errorMsg(s string) {
	// Pring status message if quiet isn't set
	if !Quiet {
		fmt.Println("")
		fmt.Println("##############################################################################")
		fmt.Printf("  ERROR: %s\n", s)
		fmt.Println("##############################################################################")
		fmt.Println("")
	}
	Error.Println(s)
}

// Output a blatant error message and log the string as an error
func traceMsg(s string) {
	// Pring status message if quiet isn't set
	if TraceOn {
		Trace.Println(s)
	}
}

// getDojoRelease retrives the supplied version of DefectDojo from the Git repo
// and places it in the specified dojoSource directory (default is /opt/dojo)
func getDojoRelease(i *config.InstallConfig) error {
	statusMsg(fmt.Sprintf("Downloading the configured release of DefectDojo => version %+v", i.Version))
	s := spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	s.Prefix = "Downloading release..."
	s.Start()

	// Create the directory to clone the source into if it doesn't exist already
	traceMsg("Creating the Dojo root directory if it doesn't exist already")
	_, err := os.Stat(i.Root)
	if err != nil {
		// Source directory doesn't exist
		err = os.MkdirAll(i.Root, 0755)
		if err != nil {
			traceMsg(fmt.Sprintf("Error creating Dojo root directory was: %+v", err))
			// TODO: Better handle the case when the repo already exists at that path - maybe?
			return err
		}
	}

	// Setup needed info
	dwnURL := ReleaseURL + i.Version + ".tar.gz"
	tarball := i.Root + "/dojo-v" + i.Version + ".tar.gz"
	traceMsg(fmt.Sprintf("Relese download list is %+v", dwnURL))
	traceMsg(fmt.Sprintf("File path to write tarball is %+v", tarball))

	// Setup a custom http client for downloading the Dojo release
	var ddClient = &http.Client{
		// Set time to a max of 20 seconds
		Timeout: time.Second * 20,
	}
	traceMsg("http.Client timeout set to 20 seconds for release download")

	// Download requested release from Dojo's Github repo
	traceMsg(fmt.Sprintf("Downloading release from %+v", dwnURL))
	resp, err := ddClient.Get(dwnURL)
	if resp != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				traceMsg(fmt.Sprintf("Error closing response.\nError was: %v", err))
				os.Exit(1)
			}
		}()
	}
	if err != nil {
		traceMsg(fmt.Sprintf("Error downloading from %+v", dwnURL))
		traceMsg(fmt.Sprintf("Error downloading was: %+v", err))
		return err
	}

	// TODO: Check for 200 status before moving on
	traceMsg(fmt.Sprintf("Status of http.Client response was %+v", resp.Status))

	// Create the file handle
	traceMsg("Creating file for downloaded tarball")
	out, err := os.Create(tarball)
	if err != nil {
		traceMsg(fmt.Sprintf("Error creating tarball was: %+v", err))
		return err
	}

	// Write the content downloaded into the file
	traceMsg("Writing downloaded content to tarball file")
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		traceMsg(fmt.Sprintf("Error writing file contents was: %+v", err))
		return err
	}

	// Extract the tarball to create the Dojo source directory
	traceMsg("Extracting tarball into the Dojo source directory")
	tb, err := os.Open(tarball)
	if err != nil {
		traceMsg(fmt.Sprintf("Error openging tarball was: %+v", err))
		return err
	}
	err = Untar(i.Root, tb)
	if err != nil {
		traceMsg(fmt.Sprintf("Error extracting tarball was: %+v", err))
		return err
	}

	// Remane source directory to the non-versioned name
	traceMsg("Renaming source directory to the non-versioned name")
	oldPath := filepath.Join(i.Root, "django-DefectDojo-"+i.Version)
	newPath := filepath.Join(i.Root, i.Source)
	err = os.Rename(oldPath, newPath)
	if err != nil {
		traceMsg(fmt.Sprintf("Error renaming Dojo source directory was: %+v", err))
		return err
	}

	// Successfully extracted the file, return nil
	s.Stop()
	statusMsg("Successfully downloaded and extracted the DefectDojo release file")
	return nil
}

// Use go-git to checkout latest source - either from a specific commit or HEAD on a branch
// and places it in the specified dojoSource directory (default is /opt/dojo)
func getDojoSource(i *config.InstallConfig) error {
	statusMsg("Downloading DefectDojo source as a branch or commit from the repo directly")
	s := spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	s.Prefix = "Downloading DefectDojo source..."

	// Create the directory to clone the source into if it doesn't exist already
	traceMsg("Creating source directory if it doesn't exist already")
	srcPath := filepath.Join(i.Root, i.Source)
	_, err := os.Stat(srcPath)
	if err != nil {
		// Source directory doesn't exist
		err = os.MkdirAll(srcPath, 0755)
		if err != nil {
			traceMsg(fmt.Sprintf("Error creating Dojo source directory was: %+v", err))
			// TODO: Better handle the case when the repo already exists at that path - maybe?
			return err
		}
	}

	// Check out a specific branch or commit - but only one of those
	// In the case that both commit and branch are set to non-empty strings,
	// the configured commit will win (aka only the commit alone will be done)
	traceMsg("Determining if a commit or branch will be checked out of the repo")
	if len(i.SourceCommit) > 0 {
		// Commit is set, so it will be used and branch ignored
		statusMsg(fmt.Sprintf("Dojo will be installed from commit %+v", i.SourceCommit))
		s.Start()

		// Do the initial clone of DefectDojo from Github
		traceMsg(fmt.Sprintf("Initial clone of %+v", CloneURL))
		repo, err := git.PlainClone(srcPath, false, &git.CloneOptions{URL: CloneURL})
		if err != nil {
			traceMsg(fmt.Sprintf("Error cloning the DefectDojo repo was: %+v", err))
			return err
		}

		// Setup the working tree for checking out a particular commit
		traceMsg("Setting up the working tree to checkout the commit")
		wk, _ := repo.Worktree()
		// TODO: consider checking the err above that is removed with _
		err = wk.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(i.SourceCommit)})
		if err != nil {
			fmt.Printf("Error checking out was %+v\n", err)
			traceMsg(fmt.Sprintf("Error checking out was: %+v", err))
			return err
		}

	} else {
		if len(i.SourceBranch) == 0 {
			// Handle the case that both source commit and branch are wonky
			err = fmt.Errorf("Both source commit and branch have empty or nonsensical values configured.\n"+
				"  Source commit was configured as %s and branch was configured as %s", i.SourceCommit, i.SourceBranch)
			traceMsg(fmt.Sprintf("Error checking out Dojo source was: %+v", err))
			return err
		}
		statusMsg(fmt.Sprintf("DefectDojo will be installed from %+v branch", i.SourceBranch))
		s.Start()

		// Check out a specific branch
		// Note: Branch and tag references are a bit odd, see https://github.com/src-d/go-git/blob/master/_examples/branch/main.go#L33
		//       However, the installer appends the necessary string to the 'normal' branch name
		traceMsg(fmt.Sprintf("Checking out branch %+v", i.SourceBranch))
		_, err = git.PlainClone(srcPath, false, &git.CloneOptions{
			URL:           CloneURL,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + i.SourceBranch),
			SingleBranch:  true,
		})
		if err != nil {
			traceMsg(fmt.Sprintf("Error checking out branch was: %+v", err))
			return err
		}

	}

	// Successfully checked out the configured source, return nil
	s.Stop()
	statusMsg("Successfully checked out the configured DefectDojo source")
	return nil
}

func sendCmd(o io.Writer, cmd string, lerr string, hard bool) {
	// Setup command
	runCmd := exec.Command("bash", "-c", cmd)
	_, err := o.Write([]byte("[godojo] # " + Redactatron(cmd, Redact) + "\n"))
	if err != nil {
		errorMsg(fmt.Sprintf("Failed to setup command, error was: %+v", err))
	}
	// TODO: Remove DEBUG below
	//fmt.Println("\nRunning ", cmd)

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

func main() {
	// Read command-line args, if any
	arg := readArgs()

	// Handle default and dev installs
	if arg.Default || arg.Dev {
		// Set config options based on embedded default config
		defaultConfig()
		if arg.Dev {
			// TODO: Write this bit
			setDevDefaults()
		}
	} else {
		// Read dojoConfig.yml file
		readConfigFile()
	}

	// Read in any environmental variables
	readEnvVars()

	// TODO Consider a "show vars" command-line option to print supported envivonmental variables

	// Setup output and logging levels and print the DefectDojo banner if needed
	Quiet = conf.Install.Quiet
	TraceOn = conf.Install.Trace
	Redact = conf.Install.Redact
	if !(Quiet || conf.Options.Embd) {
		dojoBanner()
	}
	// Setup strings to be redacted
	InitRedact(&conf)

	// Check that user is root for the installer or run with "sudo godojo"
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if usr.Uid != "0" && !conf.Options.UsrInst {
		fmt.Println("")
		fmt.Println("##############################################################################")
		fmt.Println("  ERROR: This program must be run as root or with sudo\n  Please correct and run installer again")
		fmt.Println("##############################################################################")
		fmt.Println("")
		os.Exit(1)
	}

	// Setup logging for the installer
	n := time.Now()
	when := strconv.Itoa(int(n.UnixNano()))
	logName := "dojo-install_" + when + ".log"
	logPath := path.Join(logLocation, logName)
	// Create the logs directory if it does not exist
	_, err = os.Stat(logPath)
	if err != nil {
		// logs directory doesn't exist
		err = os.MkdirAll(logLocation, 0755)
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
	// Log everything to the specificied log file location
	logSetup(logFile)

	// Logging is setup, start using statusMsg and errorMsg functions for output
	traceMsg("Logging established, trace log begins here")
	sectionMsg("Starting the dojo install at " + n.Format("Mon Jan 2, 2006 15:04:05 MST"))

	// Setup OS command logging
	traceMsg("Creating log file for OS command output for debugging reasons")
	cmdLog := "cmd-output_" + when + ".log"
	cmdPath := path.Join(logLocation, cmdLog)
	// Create command output log file in the existing logging directory
	cmdFile, err := os.OpenFile(cmdPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("")
		fmt.Println("##############################################################################")
		fmt.Printf("  ERROR: Failed to open OS Command log file %s.  Error was:\n    %+v\n", cmdPath, err)
		fmt.Println("##############################################################################")
		fmt.Println("")
		fmt.Println("Log files are required for the install, exiting install")
		os.Exit(1)
	}
	traceMsg(fmt.Sprintf("Successfully created OS Command log file at %+v", cmdPath))

	// Write out the runtime config based on the net of the config file + ENV variables
	// TODO: Consider moving this closer to the end of main - if it's not getting changed, here is fine...
	traceMsg("Writing out the runtime install configuration file")
	err = viper.WriteConfigAs("runtime-install-config.yml")
	if err != nil {
		errorMsg(fmt.Sprintf("Error from writing the runtime config was: %+v", err))
		os.Exit(1)
	}

	// Check options after logging is turned on
	if conf.Options.Embd {
		Quiet = true
		err = extr()
		if err != nil {
			fmt.Printf("Configuration has Embd = %v but no embedded assets available\n", conf.Options.Embd)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check install OS
	sectionMsg("Determining OS for installation")

	// TODO: write OS determination code for OS X
	// TODO: test OS detection on Alpine Linux docker
	target := targetOS{}
	determineOS(&target)

	// TODO: Need to write a function that takes target and validates it's supporeted by godojo
	statusMsg(fmt.Sprintf("OS was determined to be %+v, %+v", strings.Title(target.os), strings.Title(target.id)))
	statusMsg("DefectDojo installation on this OS is supported, continuing")

	// Bootstrap installer
	sectionMsg("Bootstrapping the godojo installer")
	bs := osCmds{}
	initBootstrap(target.id, &bs)

	Spin := spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	Spin.Prefix = "Bootstrapping..."
	Spin.Start()
	// TODO REMOVE COMMENTS BELOW
	for i := range bs.cmds {
		sendCmd(cmdFile,
			bs.cmds[i],
			bs.errmsg[i],
			bs.hard[i])
	}
	Spin.Stop()
	statusMsg("Boostraping godojo installer complete")

	sectionMsg("Checking for Python 3")
	if checkPythonVersion() {
		statusMsg("Python 3 found, install can continue")
	} else {
		errorMsg("Python 3 wasn't found, quitting installer")
		os.Exit(1)
	}

	sectionMsg("Downloading the source for DefectDojo")

	// Determine if a release or Dojo source will be installed
	traceMsg(fmt.Sprintf("Determining if this is a source or release install: SourceInstall is %+v", conf.Install.SourceInstall))
	if conf.Install.PullSource {
		// TODO: Move this to a separate function
		if conf.Install.SourceInstall {
			// Checkout the Dojo source directly from Github
			traceMsg("Dojo will be installed from source")

			err = getDojoSource(&conf.Install)
			if err != nil {
				errorMsg(fmt.Sprintf("Error attempting to install Dojo source was:\n    %+v", err))
				os.Exit(1)
			}
		} else {
			// Download Dojo source as a Github release tarball
			traceMsg("Dojo will be installed from a release tarball")

			err = getDojoRelease(&conf.Install)
			if err != nil {
				errorMsg(fmt.Sprintf("Error attempting to install Dojo from a release tarball was:\n    %+v", err))
				os.Exit(1)
			}
		}
	} else {
		statusMsg("No source for DefectDojo downloaded per configuration")
		traceMsg("Source NOT downloaded sa PullSource is false")
	}

	// Stup for prompting for install-time items
	if conf.Install.Prompt {
		sectionMsg("Prompt set to true, interactive installation beginning")
		fmt.Println("TODO: Write prompting code")
		os.Exit(1)
	} else {
		sectionMsg("Prompt set to false, non-interactive installation")
	}

	// Gather OS commands to bootstrap the install
	sectionMsg("Installing OS packages needed for DefectDojo")
	osInst := osCmds{}
	initOSInst(target.id, &osInst)

	// Install the OS packages
	Spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	Spin.Prefix = "Installing OS packages..."
	Spin.Start()
	// TODO REMOVE COMMENTS BELOW
	for i := range osInst.cmds {
		sendCmd(cmdFile,
			osInst.cmds[i],
			osInst.errmsg[i],
			osInst.hard[i])
	}
	Spin.Stop()
	statusMsg("Installing OS packages complete")

	// InstallDB (if needed)
	if !conf.Install.DB.Local && !conf.Install.DB.Exists {
		// Remote database that doesn't exist - godojo can't help you here
		errorMsg("Remote database which doens't exist confgiured - unsupported option")
		statusMsg("Correct configuration or install remote DB before continuing")
		fmt.Printf("Exiting...\n\n")
		os.Exit(1)
	} else if !conf.Install.DB.Exists {
		// Handle the case that the DB is local and doesn't exist
		sectionMsg("Installing database needed for DefectDojo")

		// Gather OS commands to install the DB
		dbInst := osCmds{}
		dbConf := &conf.Install.DB
		installDB(target.id, dbConf, &dbInst)

		// Run the commands to install the chosen DB
		Spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
		Spin.Prefix = "Installing " + conf.Install.DB.Engine + " database for DefectDojo..."
		Spin.Start()
		for i := range dbInst.cmds {
			sendCmd(cmdFile,
				dbInst.cmds[i],
				dbInst.errmsg[i],
				dbInst.hard[i])
		}
		Spin.Stop()
		statusMsg("Installing Database complete")

	}

	// Start the database if local and didn't already exist
	dbConf := &conf.Install.DB
	if conf.Install.DB.Local && !conf.Install.DB.Exists {
		// Handle the case that the DB is local and doesn't exist
		sectionMsg("Starting the database needed for DefectDojo")

		// Gather OS commands to install the DB
		dbStart := osCmds{}
		startDB(target.id, dbConf, &dbStart)

		// Run the commands to install the chosen DB
		Spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
		Spin.Prefix = "Starting " + conf.Install.DB.Engine + " database for DefectDojo..."
		Spin.Start()
		for i := range dbStart.cmds {
			sendCmd(cmdFile,
				dbStart.cmds[i],
				dbStart.errmsg[i],
				dbStart.hard[i])
		}
		Spin.Stop()
		statusMsg("Installing Database complete")
	}

	// Preapare the database for DefectDojo by:
	// (1) Checking connectivity to the DB, (2) checking that the configured Dojo database name doesn't exit already
	// (3) Droping the existing database if Drop = true is configured (4) Create the DefectDojo database
	// (5) Add the DB user for DefectDojo to use
	sectionMsg("Preparing the database needed for DefectDojo")
	err = dbPrep(target.id, dbConf)
	if err != nil {
		errorMsg(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	// Prep OS (user, virtualenv, chownership)
	sectionMsg("Preparing the OS for DefectDojo installation")
	prepCmds := osCmds{}
	osPrep(target.id, &conf.Install, &prepCmds)
	// Run the OS Prep commands
	Spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	Spin.Prefix = "Preparing the OS for DefectDojo..."
	Spin.Start()
	for i := range prepCmds.cmds {
		sendCmd(cmdFile,
			prepCmds.cmds[i],
			prepCmds.errmsg[i],
			prepCmds.hard[i])
	}
	Spin.Stop()
	statusMsg("Preparing the OS complete")

	// Create settings.py for DefectDojo
	sectionMsg("Creating settings.py for DefectDojo")
	settCmds := osCmds{}
	createSettingsPy(target.id, &conf, &settCmds)
	// Run the commands to create settings.py
	// TODO: Write values to .env.prod file
	Spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	Spin.Prefix = "Creating settings.py for DefectDojo..."
	Spin.Start()
	for i := range settCmds.cmds {
		sendCmd(cmdFile,
			settCmds.cmds[i],
			settCmds.errmsg[i],
			settCmds.hard[i])
	}
	Spin.Stop()
	statusMsg("Creating settings.py for DefectDojo complete")

	// Django/Python installs
	sectionMsg("Setting up Django for DefectDojo")
	setupDj := osCmds{}
	setupDjango(target.id, &conf, &setupDj)
	// Run the Django commands
	Spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	Spin.Prefix = "Setting up Django for DefectDojo..."
	Spin.Start()
	for i := range setupDj.cmds {
		sendCmd(cmdFile,
			setupDj.cmds[i],
			setupDj.errmsg[i],
			setupDj.hard[i])
	}
	Spin.Stop()
	statusMsg("Setting up Django complete")

	// Static items

	// Celery / TODO: RabitMQ

	// Optional Installs

	// Look at setup.bash's high-level workflow
	statusMsg(fmt.Sprintf("\n\nSuccessfully reached the end of main in godojo version %+v", ver))
}
