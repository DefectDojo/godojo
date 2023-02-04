package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// InstallTargets holds the distro and versions supported for Dojo installations
// TODO: Use this to check install target after determining what it is -
//       currently unused
var InstallTargets = map[string][]string{

	"ubuntu": {"18.04", "20.04", "20.10"},
	"debian": {"stretch", "buster"},
}

// Supported OSes
type targetOS struct {
	id      string
	os      string
	distro  string
	release string
}

func determineOS(tOS *targetOS) {
	// Determine OS first
	tOS.os = runtime.GOOS
	traceMsg(fmt.Sprintf("Determining OS based on GOOS: %+v", tOS.os))

	switch tOS.os {
	case "linux":
		traceMsg("OS determined to be Linux")
		determineLinux(tOS)
	case "darwin":
		traceMsg("OS determined to be Darwin/OS X")
		fmt.Println("OS X/Darwin")
		errorMsg("OS X is not YET a supported installation platform")
		os.Exit(1)
	case "windows":
		traceMsg("OS determined to be Windows")
		errorMsg("Windows is not a supported installation platform")
		os.Exit(1)
	}

	return
}

func determineLinux(tOS *targetOS) {
	// Determine the Linux Distro the installer is running on
	// Based on Based on https://unix.stackexchange.com/questions/6345/how-can-i-get-distribution-name-and-version-number-in-a-simple-shell-script
	traceMsg("Determining what Linux distro is the target OS")

	// freedesktop.org and systemd
	_, err := os.Stat("/etc/os-release")
	if err == nil {
		// That file exists
		traceMsg("Determining Linux distro from /etc/os-release")
		tOS.distro, tOS.release, tOS.id = parseOSRelease("/etc/os-release")
		return
	}

	// lsb_release command is present
	lsbCmd, err := exec.LookPath("lsb_release")
	if err == nil {
		// The command was found
		traceMsg("Determining Linux distro from lsb_release command")
		tOS.distro, tOS.release, tOS.id = parseLsbCmd(lsbCmd)
		return
	}

	// /etc/lsb-release is present
	_, err = os.Stat("/etc/lsb-release")
	if err == nil {
		// The file was found
		traceMsg("Determining Linux distro from /etc/lsb-release")
		tOS.distro, tOS.release, tOS.id = parseEtcLsb("/etc/lsb-release")
		return
	}

	// /etc/issue is present
	_, err = os.Stat("/etc/issue")
	if err == nil {
		// The file was found
		traceMsg("Determining Linux distro from /etc/issue")
		tOS.distro, tOS.release, tOS.id = parseEtcIss("/etc/issue")
		return
	}

	// /etc/debian_version is present
	_, err = os.Stat("/etc/debian_version")
	if err == nil {
		// The file was found
		traceMsg("Determining Linux distro from /etc/debian_version")
		tOS.distro, tOS.release, tOS.id = parseEtcDeb("/etc/debian_version")
		return
	}

	// Older SUSE Linux installation
	_, err = os.Stat("/etc/SuSe-release")
	if err == nil {
		// Distro is too old, not supported
		traceMsg("Older SuSe Linux distro isn't supported by this installer")
		errorMsg("Older versions of SuSe Linux are not suppported, quitting")
		os.Exit(1)
	}
	_, err = os.Stat("/etc/redhat-release")
	if err == nil {
		// Distro is too old, not supported
		traceMsg("Older RedHat Linux distro isn't supported by this installer")
		errorMsg("Older versions of Redhat Linux are not suppported, quitting")
		os.Exit(1)
	}

	traceMsg("Unable to determine the linux distro, assuming unsupported.")
	errorMsg("Unable to determine the Linux install target, quitting")
	os.Exit(1)
}

func parseOSRelease(f string) (string, string, string) {
	// Setup a map of what we need to what /etc/os-release uses
	fields := map[string]string{
		"distro":  "ID",
		"release": "VERSION_ID",
	}
	linMap := parseFile(f, "=", fields)

	return linMap["distro"], linMap["release"], linMap["distro"] + ":" + linMap["release"]

}

func parseLsbCmd(cmd string) (string, string, string) {
	// Setup map to hold parsed values
	vals := make(map[string]string)

	// Execute the lsb_release command with -a (all) and parse the output
	runCmd := exec.Command(cmd, "-a")

	// Run command and gather its output
	cmdOut, err := runCmd.CombinedOutput()
	if err != nil {
		errorMsg(fmt.Sprintf("Failed to run OS command, error was: %+v", err))
		os.Exit(1)
	}

	// Parse command output for the strings we need
	lines := bytes.Split(cmdOut, []byte("\n"))
	for _, line := range lines {
		l := string(line)

		// Look for the distro
		if strings.HasPrefix(l, "Distributor ID") {
			dis := strings.SplitN(l, ":", 2)
			vals["distro"] = strings.ToLower(strings.Trim(dis[1], "\n\t\" "))
		}

		// Look for the release
		if strings.HasPrefix(l, "Release") {
			rel := strings.SplitN(l, ":", 2)
			vals["release"] = strings.ToLower(strings.Trim(rel[1], "\n\t\" "))
		}
	}

	if _, ok := vals["distro"]; !ok {
		// The distro key hasn't been set above
		errorMsg("Unable to determine distro from lsb_release command, quitting.")
		os.Exit(1)
	}
	if _, ok := vals["release"]; !ok {
		// The distro key hasn't been set above
		errorMsg("Unable to determine release from lsb_release command, quitting.")
		os.Exit(1)
	}

	return vals["distro"], vals["release"], vals["distro"] + ":" + vals["release"]
}

func parseEtcLsb(f string) (string, string, string) {
	// Setup a map of what we need to what /etc/lsb-release uses
	fields := map[string]string{
		"distro":  "DISTRIB_ID",
		"release": "DISTRIB_RELEASE",
	}
	linMap := parseFile(f, "=", fields)

	return linMap["distro"], linMap["release"], linMap["distro"] + ":" + linMap["release"]
}

func parseEtcIss(f string) (string, string, string) {
	// Setup return map
	vals := make(map[string]string)

	// Open the file for parsing
	file, err := os.Open(f)
	if err != nil {
		errorMsg(fmt.Sprintf("Unable to open file: %+v\nError was: %v", f, err))
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			traceMsg(fmt.Sprintf("Erro closing file\nError was: %v", err))
			os.Exit(1)
		}
	}()

	// Read the file in, pull off the first line and split it
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		errorMsg(fmt.Sprintf("Unable to read file: %+v\nError was: %v", f, err))
		os.Exit(1)
	}
	fields := strings.Split(line, " ")
	vals["distro"] = strings.ToLower(fields[0])
	vals["release"] = fields[1]

	// Correct for Ubuntu 'minor' releases aka 18.04.2
	if vals["distro"] == "ubuntu" {
		tmp := strings.Split(vals["release"], ".")
		vals["release"] = tmp[0] + "." + tmp[1]
	}

	return vals["distro"], vals["release"], vals["distro"] + ":" + vals["release"]
}

func parseEtcDeb(f string) (string, string, string) {
	// Setup map to hold parsed values
	vals := make(map[string]string)
	vals["distro"] = "debian"

	// Open the file for parsing
	file, err := os.Open(f)
	if err != nil {
		errorMsg(fmt.Sprintf("Unable to open file: %+v\nError was: %v", f, err))
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			errorMsg(fmt.Sprintf("Unable to close file\nError was: %v", err))
			os.Exit(1)
		}
	}()

	// Read the file in, pull off the first line
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		errorMsg(fmt.Sprintf("Unable to read file: %+v\nError was: %v", f, err))
		os.Exit(1)
	}
	// TODO: Test this with a Debian docker
	vals["release"] = strings.ToLower(strings.Trim(line, "\n\t "))

	return vals["distro"], vals["release"], vals["distro"] + ":" + vals["release"]
}

func parseFile(f string, sep string, flds map[string]string) map[string]string {
	// Setup return map
	vals := make(map[string]string)

	// Open the file for parsing
	file, err := os.Open(f)
	if err != nil {
		errorMsg(fmt.Sprintf("Unable to open file: %+v\nError was: %v", f, err))
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			errorMsg(fmt.Sprintf("Unable to close file\nError was: %v", err))
			os.Exit(1)
		}
	}()

	// Read the file one line at a time till done
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// Look for the distro
		if strings.HasPrefix(line, flds["distro"]+sep) {
			dis := strings.SplitN(line, sep, 2)
			vals["distro"] = strings.ToLower(strings.Trim(dis[1], "\n\""))
		}

		// Look for the release
		if strings.HasPrefix(line, flds["release"]+sep) {
			rel := strings.SplitN(line, sep, 2)
			vals["release"] = strings.ToLower(strings.Trim(rel[1], "\n\""))
		}
	}

	return vals
}

// TODO: Make this fail for unsupported OSes
// TODO: Consider making a bootstrap per Linux distro switching to those here
func initBootstrap(id string, b *osCmds) {
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

func checkPythonVersion() bool {
	// DefectDojo is now Python 3+, lets make sure that's installed
	_, err := exec.LookPath("python3")
	if err != nil {
		errorMsg(fmt.Sprintf("Unable to find python binary. Error was: %+v", err))
		os.Exit(1)
	}

	// Execute the python3 command with --version to get the version
	runCmd := exec.Command("python3", "--version")

	// Run command and gather its output
	cmdOut, err := runCmd.CombinedOutput()
	if err != nil {
		errorMsg(fmt.Sprintf("Failed to run python3 command, error was: %+v", err))
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
