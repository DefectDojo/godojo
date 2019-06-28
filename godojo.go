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
	"strconv"
	"strings"
	"time"

	"github.com/mtesauro/godojo/config"
	"github.com/spf13/viper"
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
		Timeout: time.Second * 20
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

	return nil
}

func getDojoSource(c string) error {
	// Use go-git to checkout latest source
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
		//panic("TODO: Better error handling 03")
	}

	// Checkout the Dojo source - either a versioned release or using git
	fmt.Println("Getting dojo from repository")
	err = getDojoRelease(&conf.Install)
	if err != nil {
		panic("TODO: Better error handling 03")
	}

	// Start stub'ing out stuff
	// Look at setup.bash's high-level workflow
}
