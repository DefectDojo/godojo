package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// bootstrapInstall takes a pointer to a gdjDefault struct and a targetOS struct
// to run the commands necessary to bootstrap the installation
func bootstrapInstall(d *gdjDefault, t *targetOS) {
	d.sectionMsg("Bootstrapping the godojo installer")
	bs := osCmds{}
	initBootstrap(t.id, &bs)

	// Start the spinner
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Bootstrapping..."
	d.spin.Start()
	// Run the boostrapping commands for the target OS
	for i := range bs.cmds {
		sendCmd(d,
			d.cmdLogger,
			bs.cmds[i],
			bs.errmsg[i],
			bs.hard[i])
	}
	d.spin.Stop()
	d.statusMsg("Boostraping godojo installer complete")

}

// initBootstrap takes an id from targetOS and a pointer to a osCmds struct to
// run the commands needd to bootstrap the OS for installation
func initBootstrap(id string, b *osCmds) {
	// TODO: Make this fail for unsupported OSes
	// TODO: Consider making a bootstrap per Linux distro switching to those here
	switch strings.ToLower(distOnly(id)) {
	case "debian":
		fallthrough
	case "ubuntu":
		b.id = id
		b.cmds = []string{
			"DEBIAN_FRONTEND=noninteractive apt-get update",
			"DEBIAN_FRONTEND=noninteractive apt-get -y upgrade",
			"DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::=\"--force-confdef\" -o Dpkg::Options::=\"--force-confold\" install python3 python3-virtualenv ca-certificates curl gnupg git sudo",
		}
		b.errmsg = []string{
			"Unable to update apt database",
			"Unable to upgrade OS packages with apt",
			"Unable to install prerequisites for installer via apt",
		}
		b.hard = []bool{
			true,
			true,
			false,
		}

		return
	default:
		fmt.Println("Unsupported OS to bootstrap, quitting.")
		os.Exit(1)

	}
}

// validPython checks to ensure the correct version of Python is available
func validPython(d *gdjDefault) {
	d.sectionMsg("Checking for Python 3")
	if checkPythonVersion(d) {
		d.statusMsg("Python 3 found, install can continue")
	} else {
		d.errorMsg("Python 3 wasn't found, quitting installer")
		os.Exit(1)
	}
}

// checkPythonVersion verifies that python3 is availble on the install target
func checkPythonVersion(d *gdjDefault) bool {
	// DefectDojo is now Python 3+, lets make sure that's installed
	_, err := exec.LookPath("python3")
	if err != nil {
		d.errorMsg(fmt.Sprintf("Unable to find python binary. Error was: %+v", err))
		os.Exit(1)
	}

	// Execute the python3 command with --version to get the version
	runCmd := exec.Command("python3", "--version")

	// Run command and gather its output
	cmdOut, err := runCmd.CombinedOutput()
	if err != nil {
		d.errorMsg(fmt.Sprintf("Failed to run python3 command, error was: %+v", err))
		os.Exit(1)
	}

	// Parse command output for the strings we need
	lines := bytes.Split(cmdOut, []byte("\n"))
	line := strings.Split(string(lines[0]), " ")
	pyVer := line[1]

	// TODO: Consider checking the minor version of Python3 as well - probably not needed (yet)
	if strings.HasPrefix(pyVer, "3.") {
		return true
	}
	// DefectDojo requires Python 3.x
	return false
}

// downloadDojo takes a ponter to gdjDefault and downloads a release or source
// code depending on the configuration of dojoConfig.yml
func downloadDojo(d *gdjDefault) {
	d.sectionMsg("Downloading the source for DefectDojo")

	// Determine if a release or Dojo source will be installed
	d.traceMsg(fmt.Sprintf("Determining if this is a source or release install: SourceInstall is %+v", d.conf.Install.SourceInstall))
	if d.conf.Install.PullSource {
		// TODO: Move this to a separate function
		if d.conf.Install.SourceInstall {
			// Checkout the Dojo source directly from Github
			d.traceMsg("Dojo will be installed from source")

			err := getDojoSource(d)
			if err != nil {
				d.errorMsg(fmt.Sprintf("Error attempting to install Dojo source was:\n    %+v", err))
				os.Exit(1)
			}
		} else {
			// Download Dojo source as a Github release tarball
			d.traceMsg("Dojo will be installed from a release tarball")

			err := getDojoRelease(d)
			if err != nil {
				d.errorMsg(fmt.Sprintf("Error attempting to install Dojo from a release tarball was:\n    %+v", err))
				os.Exit(1)
			}
		}
	} else {
		d.statusMsg("No source for DefectDojo downloaded per configuration")
		d.traceMsg("Source NOT downloaded as PullSource is false")
	}
}

// getDojoRelease retrives the supplied version of DefectDojo from the Git repo
// and places it in the specified dojoSource directory (default is /opt/dojo)
func getDojoRelease(d *gdjDefault) error {
	d.statusMsg(fmt.Sprintf("Downloading the configured release of DefectDojo => version %+v", d.conf.Install.Version))
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Downloading release..."
	d.spin.Start()

	// Create the directory to clone the source into if it doesn't exist already
	d.traceMsg("Creating the Dojo root directory if it doesn't exist already")
	_, err := os.Stat(d.conf.Install.Root)
	if err != nil {
		// Source directory doesn't exist
		err = os.MkdirAll(d.conf.Install.Root, 0755)
		if err != nil {
			d.traceMsg(fmt.Sprintf("Error creating Dojo root directory was: %+v", err))
			// TODO: Better handle the case when the repo already exists at that path - maybe?
			return err
		}
	}

	// Setup needed info
	dwnURL := d.releaseURL + d.conf.Install.Version + ".tar.gz"
	tarball := d.conf.Install.Root + "/dojo-v" + d.conf.Install.Version + ".tar.gz"
	d.traceMsg(fmt.Sprintf("Relese download list is %+v", dwnURL))
	d.traceMsg(fmt.Sprintf("File path to write tarball is %+v", tarball))

	// Check for existing tarball before downloading, might be a re-run of godojo
	_, err = os.Stat(tarball)
	if err == nil {
		// File already downloaded so return early
		err = extractRelease(d, tarball)
		if err != nil {
			return err
		}
		d.spin.Stop()
		d.statusMsg("Tarball already downloaded and extracted the DefectDojo release file")
		return nil
	}

	// Setup a custom http client for downloading the Dojo release
	var ddClient = &http.Client{
		// Set time to a max of 60 seconds
		Timeout: time.Second * 120,
	}
	d.traceMsg("http.Client timeout set to 120 seconds for release download")

	// Download requested release from Dojo's Github repo
	d.traceMsg(fmt.Sprintf("Downloading release from %+v", dwnURL))
	resp, err := ddClient.Get(dwnURL)
	if resp != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				d.traceMsg(fmt.Sprintf("Error closing response.\nError was: %v", err))
				os.Exit(1)
			}
		}()
	}
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error downloading from %+v", dwnURL))
		d.traceMsg(fmt.Sprintf("Error downloading was: %+v", err))
		return err
	}

	// TODO: Check for 200 status before moving on
	d.traceMsg(fmt.Sprintf("Status of http.Client response was %+v", resp.Status))

	// Create the file handle
	d.traceMsg("Creating file for downloaded tarball")
	out, err := os.Create(tarball)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error creating tarball was: %+v", err))
		return err
	}

	// Write the content downloaded into the file
	d.traceMsg("Writing downloaded content to tarball file")
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error writing file contents was: %+v", err))
		return err
	}

	// Extract the tarball to create the Dojo source directory
	err = extractRelease(d, tarball)
	if err != nil {
		return err
	}

	// Successfully extracted the file, return nil
	d.spin.Stop()
	d.statusMsg("Successfully downloaded and extracted the DefectDojo release file")
	return nil
}

func extractRelease(d *gdjDefault, t string) error {
	// Extract the tarball to create the Dojo source directory
	d.traceMsg("Extracting tarball into the Dojo source directory")
	tb, err := os.Open(t)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error openging tarball was: %+v", err))
		return err
	}
	err = untar(d, d.conf.Install.Root, tb)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error extracting tarball was: %+v", err))
		return err
	}

	// Remane source directory to the non-versioned name
	d.traceMsg("Renaming source directory to the non-versioned name")
	oldPath := filepath.Join(d.conf.Install.Root, "django-DefectDojo-"+d.conf.Install.Version)
	newPath := filepath.Join(d.conf.Install.Root, d.conf.Install.Source)
	err = os.Rename(oldPath, newPath)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error renaming Dojo source directory was: %+v", err))
		return err
	}
	return nil
}

// Use go-git to checkout latest source - either from a specific commit or HEAD
// on a branch and places it in the specified dojoSource directory
// (default is /opt/dojo)
func getDojoSource(d *gdjDefault) error {
	d.statusMsg("Downloading DefectDojo source as a branch or commit from the repo directly")
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Downloading DefectDojo source..."

	// Create the directory to clone the source into if it doesn't exist already
	d.traceMsg("Creating source directory if it doesn't exist already")
	srcPath := filepath.Join(d.conf.Install.Root, d.conf.Install.Source)
	_, err := os.Stat(srcPath)
	if err != nil {
		// Source directory doesn't exist
		err = os.MkdirAll(srcPath, 0755)
		if err != nil {
			d.traceMsg(fmt.Sprintf("Error creating Dojo source directory was: %+v", err))
			// TODO: Better handle the case when the repo already exists at that path - maybe?
			return err
		}
	}

	// Check out a specific branch or commit - but only one of those
	// In the case that both commit and branch are set to non-empty strings,
	// the configured commit will win (aka only the commit alone will be done)
	d.traceMsg("Determining if a commit or branch will be checked out of the repo")
	if len(d.conf.Install.SourceCommit) > 0 {
		// Commit is set, so it will be used and branch ignored
		d.statusMsg(fmt.Sprintf("Dojo will be installed from commit %+v", d.conf.Install.SourceCommit))
		d.spin.Start()

		// Do the initial clone of DefectDojo from Github
		d.traceMsg(fmt.Sprintf("Initial clone of %+v", d.cloneURL))
		repo, err := git.PlainClone(srcPath, false, &git.CloneOptions{URL: d.cloneURL})
		if err != nil {
			d.traceMsg(fmt.Sprintf("Error cloning the DefectDojo repo was: %+v", err))
			return err
		}

		// Setup the working tree for checking out a particular commit
		d.traceMsg("Setting up the working tree to checkout the commit")
		wk, _ := repo.Worktree()
		// TODO: consider checking the err above that is removed with _
		err = wk.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(d.conf.Install.SourceCommit)})
		if err != nil {
			fmt.Printf("Error checking out was %+v\n", err)
			d.traceMsg(fmt.Sprintf("Error checking out was: %+v", err))
			return err
		}

	} else {
		if len(d.conf.Install.SourceBranch) == 0 {
			// Handle the case that both source commit and branch are wonky
			err = fmt.Errorf("Both source commit and branch have empty or nonsensical values configured.\n"+
				"  Source commit was configured as %s and branch was configured as %s", d.conf.Install.SourceCommit, d.conf.Install.SourceBranch)
			d.traceMsg(fmt.Sprintf("Error checking out Dojo source was: %+v", err))
			return err
		}
		d.statusMsg(fmt.Sprintf("DefectDojo will be installed from %+v branch", d.conf.Install.SourceBranch))
		d.spin.Start()

		// Check out a specific branch
		// Note: Branch and tag references are a bit odd, see https://github.com/src-d/go-git/blob/master/_examples/branch/main.go#L33
		//       However, the installer appends the necessary string to the 'normal' branch name
		d.traceMsg(fmt.Sprintf("Checking out branch %+v", d.conf.Install.SourceBranch))
		_, err = git.PlainClone(srcPath, false, &git.CloneOptions{
			URL:           d.cloneURL,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + d.conf.Install.SourceBranch),
			SingleBranch:  true,
		})
		if err != nil {
			d.traceMsg(fmt.Sprintf("Error checking out branch was: %+v", err))
			return err
		}

	}

	// Successfully checked out the configured source, return nil
	d.spin.Stop()
	d.statusMsg("Successfully checked out the configured DefectDojo source")
	return nil
}
