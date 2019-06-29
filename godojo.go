package main

// TODO:
// Add Cobra for command-line args - https://github.com/spf13/cobra
// Add go-git to pull Dojo source - https://github.com/src-d/go-git
// Add code to create a log subdirectory and log the install process there
// Add redactatron function like prior installer
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mtesauro/godojo/config"
	"github.com/mtesauro/godojo/util"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Global vars
var (
	logLocation = "logs"
	Trace       *log.Logger
	Info        *log.Logger
	Warning     *log.Logger
	Error       *log.Logger
)

func logSetup(logHandler io.Writer) {
	// Setup logging 'levels' which can be called globally like Info.Println("Example info log")
	Trace = log.New(logHandler, "TRACE:   ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(logHandler, "INFO:    ", log.Ldate|log.Ltime)
	Warning = log.New(logHandler, "WARNING: ", log.Ldate|log.Ltime)
	Error = log.New(logHandler, "ERROR:   ", log.Ldate|log.Ltime|log.Lshortfile)
}

func dojoBanner() {
	fmt.Println("        ____       ____          __     ____          _      ")
	fmt.Println("       / __ \\___  / __/__  _____/ /_   / __ \\____    (_)___  ")
	fmt.Println("      / / / / _ \\/ /_/ _ \\/ ___/ __/  / / / / __ \\  / / __ \\ ")
	fmt.Println("     / /_/ /  __/ __/  __/ /__/ /_   / /_/ / /_/ / / / /_/ / ")
	fmt.Println("    /_____/\\___/_/  \\___/\\___/\\__/  /_____/\\____/_/ /\\____/  ")
	fmt.Println("                                               /___/         ")
	fmt.Println("")
	fmt.Println("  Welcome to goDojo, the official way to install Defect Dojo.")
	fmt.Println("  For more information on how goDojo does an install, see:")
	fmt.Println("  https://github.com/mtesauro/godojo")
	fmt.Println("")
}

func statusMsg(s string) {
	fmt.Println("==============================================================================")
	fmt.Printf("  %s\n", s)
	fmt.Println("==============================================================================")
	fmt.Println("")
	Info.Println(s)
}

// getDojo retrives the supplied version of DefectDojo from the Git repo
// and places it in the specified dojoSource directory (default is /opt/dojo)
func getDojoRelease(i *config.InstallConfig) error {
	// Setup needed info
	dwnURL := i.ReleaseURL + i.Version + ".tar.gz"
	tarball := i.Root + "/dojo-v" + i.Version + ".tar.gz"
	fmt.Printf("Download link is %+v\n", dwnURL)
	fmt.Printf("File path is %+v\n", tarball)

	// Setup a custom http client for downloading the Dojo release
	var ddClient = &http.Client{
		// Set time to a max of 20 seconds
		Timeout: time.Second * 20,
	}

	// Download requested release from Dojo's Github repo
	resp, err := ddClient.Get(dwnURL)
	if err != nil {
		fmt.Println("Error downloading from %+v\n", dwnURL)
		fmt.Println("Error downloading was: %+v\n", err)
		return err
	}
	defer resp.Body.Close()

	// Create the file handle
	out, err := os.Create(tarball)
	if err != nil {
		fmt.Println("Error creating tartball was: %+v\n", err)
		return err
	}

	// Write the content downloaded into the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error writing file contents was = %+v\n", err)
		return err
	}

	fmt.Printf("resp.Status = %+v\n", resp.Status)
	fmt.Printf("out = %+v\n", out)

	// Extract the tarball to create the Dojo source directory
	tb, err := os.Open(tarball)
	if err != nil {
		fmt.Printf("Error opening tarball was = %+v\n", err)
		return err
	}
	err = util.Untar(i.Root, tb)
	if err != nil {
		fmt.Printf("Error extracting tarball was = %+v\n", err)
		return err
	}

	// Remane source directory to the non-versioned name
	oldPath := filepath.Join(i.Root, "django-DefectDojo-"+i.Version)
	newPath := filepath.Join(i.Root, i.Source)
	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Printf("Error renaming Dojo source path was %+v\n", err)
		return err
	}

	// Successfully extracted the file, return nil
	return nil
}

func getDojoSource(i *config.InstallConfig) error {
	// Use go-git to checkout latest source
	// TODO clean me up
	CloneURL := "https://github.com/DefectDojo/django-DefectDojo.git"

	// Create the directory to clone the source into if it doesn't exist already
	srcPath := filepath.Join(i.Root, i.Source)
	_, err := os.Stat(srcPath)
	if err != nil {
		// Source directory doesn't exist
		err = os.MkdirAll(srcPath, 0755)
		if err != nil {
			fmt.Printf("Error creating Dojo source directory was %+v\n", err)
			// TODO: Better handle the case when the repo already exists at that path - maybe?
			return err
		}
	}

	// Check out a specific branch or commit - but only one of those
	// In the case that both commit and branch are set to non-empty strings,
	// the configured commit will win (aka only the commit alone will be done)
	if len(i.SourceCommit) > 0 {
		// Commit is set, so it will be used and branch ignored
		fmt.Println("COMMIT WINS")
		// Do the initial clone of Defect Dojo from Github
		repo, err := git.PlainClone(srcPath, false, &git.CloneOptions{URL: CloneURL})
		if err != nil {
			fmt.Printf("Error cloning the DefectDojo repo was %+v\n", err)
			return err
		}
		fmt.Printf("Repo cloned at HEAD\n\n")

		// Setup the working tree for checking out a particular branch or commit
		wk, err := repo.Worktree()
		fmt.Printf("Checkouting out commit %+v\n", i.SourceCommit)
		err = wk.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(i.SourceCommit)})
		if err != nil {
			fmt.Printf("Error checking out was %+v\n", err)
			return err
		}

	} else {
		if len(i.SourceBranch) == 0 {
			// Handle the case that both source commit and branch are wonky
			fmt.Printf("i.SourceBranch = %+v\n", i.SourceBranch)
			fmt.Println("WONKY commit and source values!")
			err = fmt.Errorf("Both source commit and branch have empty or nonsensical values configured.\n"+
				"  Source commit was configured as %s and branch was configured as %s", i.SourceCommit, i.SourceBranch)
			return err
		}
		fmt.Println("BRANCH WINS")
		// Check out a specfic branch
		// Note, Branch and tag references are a bit odd, see https://github.com/src-d/go-git/blob/master/_examples/branch/main.go#L33
		fmt.Printf("Checking out branch %+v\n", i.SourceBranch)
		_, err = git.PlainClone(srcPath, false, &git.CloneOptions{
			URL:           CloneURL,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + i.SourceBranch),
			SingleBranch:  true,
		})
		if err != nil {
			fmt.Printf("Error checking brach out was %+v\n", err)
			return err
		}

	}

	return nil
}

func main() {
	dojoBanner()
	// Check that user is root for the installer or run with "sudo godojo"
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if usr.Uid != "0" {
		fmt.Println("==============================================================================")
		fmt.Println(" This program must be run as root or with sudo\n  Please correct this and try again")
		fmt.Println("==============================================================================")
		fmt.Println("")
		// TODO: Remove DEBUG below
		// DEBUG os.Exit(1)
	}

	// Setup logging for the installer
	n := time.Now()
	when := strconv.Itoa(int(n.UnixNano()))
	logName := "dojo-install_" + when + ".log"
	// TODO: Consider creating the logs directory instead of expecting it
	logPath := path.Join(logLocation, logName)
	_, err = os.Stat(logPath)
	if err != nil {
		// Source directory doesn't exist
		err = os.MkdirAll(logLocation, 0755)
		if err != nil {
			// TODO: Change the below to log and exit more gracefully
			fmt.Printf("Error creating godojo installer logging directory was %+v\n", err)
			os.Exit(1)
		}
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("\nPlease create any directories needed to write logs to %v\n\n", logPath)
		log.Fatalf("Failed to open log file %s.  Error was:\n  %+v\n", logPath, err)
	}
	// Log everthing to the specificied log file location
	logSetup(logFile)
	statusMsg("Starting the dojo install at " + n.Format("Mon Jan 2, 2006 15:04:05 MST"))

	// Read the config and set default values
	err = config.ReadInstallConfig()
	if err != nil {
		panic("Can't read config nor setup default values")
	}

	// Setup viper config
	viper.AddConfigPath(".")
	viper.SetConfigName("dojoConfig")
	var conf config.DojoConfig

	// Setup ENV variables
	viper.SetEnvPrefix("DD")
	replace := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replace)
	viper.AutomaticEnv()

	// DEBUG
	//os.Setenv("DD_DOJO_VER", "6.6.6")

	err = viper.ReadInConfig()
	if err != nil {
		panic("TODO: Better error handling 01")
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic("TODO: Better error handling 02")
	}
	//fmt.Printf("conf.Install.Version = %+v\n\n", conf.Install.Version)
	//fmt.Printf("conf = %+v\n\n", conf)
	//fmt.Printf("conf.Install = %+v\n\n", conf.Install)
	//fmt.Printf("viper.Get(\"DojoVer\") = %+v\n", viper.Get("install.version"))
	//fmt.Printf("viper.AllKeys = %+v\n\n", viper.AllKeys())
	//fmt.Printf("viper.AllSettings = %+v\n\n", viper.AllSettings())
	//fmt.Printf("FROM VIPER viper.Get(\"install.db.engine\") = %+v\n\n", viper.Get("install.db.engine"))
	//fmt.Printf("FROM STRUCT config.Install.db.engine= %+v\n", conf.Install.DB.Engine)
	// Stopped here
	//fmt.Println(config.TestConfig())

	// DEBUG - Testing writing config
	err = viper.WriteConfigAs("runtime-install-config.yml")
	if err != nil {
		fmt.Printf("err from write was: %+v\n", err)
		panic("TODO: Better error handling 03")
	}

	// Download Dojo source as a Github release tarball
	//fmt.Println("Getting dojo from repository")
	//err = getDojoRelease(&conf.Install)
	//if err != nil {
	//	panic("TODO: Better error handling 04")
	//}

	// Checkout the Dojo source directly from Github
	err = getDojoSource(&conf.Install)
	if err != nil {
		fmt.Printf("err from getDojoSource = %+v\n", err)
		panic("TODO: Better error handling 05")
	}

	// Start stub'ing out stuff
	// Look at setup.bash's high-level workflow
}
