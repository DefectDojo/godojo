package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

// Supported OSes
type targetOS struct {
	id      string
	os      string
	distro  string
	release string
}

func checkOS(d *gdjDefault) targetOS {
	// Check install OS
	d.sectionMsg("Determining OS for installation")

	// TODO: write OS determination code for OS X
	// TODO: test OS detection on Alpine Linux docker
	target := targetOS{}
	determineOS(d, &target)

	// TODO: Need to write a function that takes target and validates it's supporeted by godojo
	d.statusMsg(fmt.Sprintf("OS was determined to be %+v, %+v", strings.Title(target.os), strings.Title(target.id)))
	d.statusMsg("DefectDojo installation on this OS is supported, continuing")

	return target
}

func determineOS(d *gdjDefault, tOS *targetOS) {
	// Determine OS first
	tOS.os = runtime.GOOS
	d.traceMsg(fmt.Sprintf("Determining OS based on GOOS: %+v", tOS.os))

	switch tOS.os {
	case "linux":
		d.traceMsg("OS determined to be Linux")
		determineLinux(d, tOS)
	case "darwin":
		d.traceMsg("OS determined to be Darwin/OS X")
		fmt.Println("OS X/Darwin")
		d.errorMsg("OS X is not YET a supported installation platform")
		os.Exit(1)
	case "windows":
		d.traceMsg("OS determined to be Windows")
		d.errorMsg("Windows is not a supported installation platform")
		os.Exit(1)
	}

	return
}

func determineLinux(d *gdjDefault, tOS *targetOS) {
	// Determine the Linux Distro the installer is running on
	// Based on Based on https://unix.stackexchange.com/questions/6345/how-can-i-get-distribution-name-and-version-number-in-a-simple-shell-script
	d.traceMsg("Determining what Linux distro is the target OS")

	// freedesktop.org and systemd
	_, err := os.Stat("/etc/os-release")
	if err == nil {
		// That file exists
		d.traceMsg("Determining Linux distro from /etc/os-release")
		tOS.distro, tOS.release, tOS.id = parseOSRelease(d, "/etc/os-release")
		return
	}

	// lsb_release command is present
	lsbCmd, err := exec.LookPath("lsb_release")
	if err == nil {
		// The command was found
		d.traceMsg("Determining Linux distro from lsb_release command")
		tOS.distro, tOS.release, tOS.id = parseLsbCmd(d, lsbCmd)
		return
	}

	// /etc/lsb-release is present
	_, err = os.Stat("/etc/lsb-release")
	if err == nil {
		// The file was found
		d.traceMsg("Determining Linux distro from /etc/lsb-release")
		tOS.distro, tOS.release, tOS.id = parseEtcLsb(d, "/etc/lsb-release")
		return
	}

	// /etc/issue is present
	_, err = os.Stat("/etc/issue")
	if err == nil {
		// The file was found
		d.traceMsg("Determining Linux distro from /etc/issue")
		tOS.distro, tOS.release, tOS.id = parseEtcIss(d, "/etc/issue")
		return
	}

	// /etc/debian_version is present
	_, err = os.Stat("/etc/debian_version")
	if err == nil {
		// The file was found
		d.traceMsg("Determining Linux distro from /etc/debian_version")
		tOS.distro, tOS.release, tOS.id = parseEtcDeb(d, "/etc/debian_version")
		return
	}

	// Older SUSE Linux installation
	_, err = os.Stat("/etc/SuSe-release")
	if err == nil {
		// Distro is too old, not supported
		d.traceMsg("Older SuSe Linux distro isn't supported by this installer")
		d.errorMsg("Older versions of SuSe Linux are not suppported, quitting")
		os.Exit(1)
	}
	_, err = os.Stat("/etc/redhat-release")
	if err == nil {
		// Distro is too old, not supported
		d.traceMsg("Older RedHat Linux distro isn't supported by this installer")
		d.errorMsg("Older versions of Redhat Linux are not suppported, quitting")
		os.Exit(1)
	}

	d.traceMsg("Unable to determine the linux distro, assuming unsupported.")
	d.errorMsg("Unable to determine the Linux install target, quitting")
	os.Exit(1)
}

func parseOSRelease(d *gdjDefault, f string) (string, string, string) {
	// Setup a map of what we need to what /etc/os-release uses
	fields := map[string]string{
		"distro":  "ID",
		"release": "VERSION_ID",
	}
	linMap := parseFile(d, f, "=", fields)

	return linMap["distro"], linMap["release"], linMap["distro"] + ":" + linMap["release"]

}

func parseLsbCmd(d *gdjDefault, cmd string) (string, string, string) {
	// Setup map to hold parsed values
	vals := make(map[string]string)

	// Execute the lsb_release command with -a (all) and parse the output
	runCmd := exec.Command(cmd, "-a")

	// Run command and gather its output
	cmdOut, err := runCmd.CombinedOutput()
	if err != nil {
		d.errorMsg(fmt.Sprintf("Failed to run OS command, error was: %+v", err))
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
		d.errorMsg("Unable to determine distro from lsb_release command, quitting.")
		os.Exit(1)
	}
	if _, ok := vals["release"]; !ok {
		// The distro key hasn't been set above
		d.errorMsg("Unable to determine release from lsb_release command, quitting.")
		os.Exit(1)
	}

	return vals["distro"], vals["release"], vals["distro"] + ":" + vals["release"]
}

func parseEtcLsb(d *gdjDefault, f string) (string, string, string) {
	// Setup a map of what we need to what /etc/lsb-release uses
	fields := map[string]string{
		"distro":  "DISTRIB_ID",
		"release": "DISTRIB_RELEASE",
	}
	linMap := parseFile(d, f, "=", fields)

	return linMap["distro"], linMap["release"], linMap["distro"] + ":" + linMap["release"]
}

func parseEtcIss(d *gdjDefault, f string) (string, string, string) {
	// Setup return map
	vals := make(map[string]string)

	// Open the file for parsing
	file, err := os.Open(f)
	if err != nil {
		d.errorMsg(fmt.Sprintf("Unable to open file: %+v\nError was: %v", f, err))
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			d.traceMsg(fmt.Sprintf("Erro closing file\nError was: %v", err))
			os.Exit(1)
		}
	}()

	// Read the file in, pull off the first line and split it
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		d.errorMsg(fmt.Sprintf("Unable to read file: %+v\nError was: %v", f, err))
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

func parseEtcDeb(d *gdjDefault, f string) (string, string, string) {
	// Setup map to hold parsed values
	vals := make(map[string]string)
	vals["distro"] = "debian"

	// Open the file for parsing
	file, err := os.Open(f)
	if err != nil {
		d.errorMsg(fmt.Sprintf("Unable to open file: %+v\nError was: %v", f, err))
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			d.errorMsg(fmt.Sprintf("Unable to close file\nError was: %v", err))
			os.Exit(1)
		}
	}()

	// Read the file in, pull off the first line
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		d.errorMsg(fmt.Sprintf("Unable to read file: %+v\nError was: %v", f, err))
		os.Exit(1)
	}
	// TODO: Test this with a Debian docker
	vals["release"] = strings.ToLower(strings.Trim(line, "\n\t "))

	return vals["distro"], vals["release"], vals["distro"] + ":" + vals["release"]
}

func parseFile(d *gdjDefault, f string, sep string, flds map[string]string) map[string]string {
	// Setup return map
	vals := make(map[string]string)

	// Open the file for parsing
	file, err := os.Open(f)
	if err != nil {
		d.errorMsg(fmt.Sprintf("Unable to open file: %+v\nError was: %v", f, err))
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			d.errorMsg(fmt.Sprintf("Unable to close file\nError was: %v", err))
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

// distOnly takes an string such as the targetOS.id and returns only the distro
// name portion of that id e.g. if sent ubuntu:20.04, ubuntu would be returned
func distOnly(d string) string {
	if strings.Contains(d, "ubuntu") {
		dist := strings.Split(d, ":")
		return dist[0]
	}
	if strings.Contains(d, "debian") {
		dist := strings.Split(d, ":")
		return dist[0]
	}
	// TODO Add more DISTROS here

	return "Unable-to-parse-distro"
}

// prepOSForDojo takes a pointer to a gdjDefault struct and a string representing
// the id for the target OS and installs the necessary OS software required by
// DefectDojo
func prepOSForDojo(d *gdjDefault, o *targetOS) {
	// Gather OS commands to bootstrap the install
	d.sectionMsg("Installing OS packages needed for DefectDojo")
	osInst := osCmds{}
	initOSInst(d, o.id, &osInst)

	// Install the OS packages
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Installing OS packages..."
	d.spin.Start()
	// Send commands to install OS packages for Target OS
	for i := range osInst.cmds {
		sendCmd(d,
			d.cmdLogger,
			osInst.cmds[i],
			osInst.errmsg[i],
			osInst.hard[i])
	}
	d.spin.Stop()
	d.statusMsg("Installing OS packages complete")
}

// initOSInst takes an id from osTarget and a pointer to osCmds struct to add
// the commands needed for the id provided to prepare the OS for installing
// DefectDojo
func initOSInst(d *gdjDefault, id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		ubuntuInitOSInst(d, id, b)

	}
	return
}

// Commands to bootstrap Ubuntu for the installer
func ubuntuInitOSInst(d *gdjDefault, id string, b *osCmds) {
	switch strings.ToLower(id) {
	case "debian:10":
		fallthrough
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		b.id = id
		b.cmds = []string{
			fmt.Sprintf("curl -sS %s | apt-key add -", d.yarnGPG),
			fmt.Sprintf("echo -n %s > /etc/apt/sources.list.d/yarn.list", d.yarnRepo),
			"DEBIAN_FRONTEND=noninteractive apt-get update",
			"DEBIAN_FRONTEND=noninteractive apt-get install sudo",
			fmt.Sprintf("curl -sL %s | bash - ", d.nodeURL),
			"DEBIAN_FRONTEND=noninteractive apt-get install -y apt-transport-https libjpeg-dev gcc libssl-dev python3-dev python3-pip python3-virtualenv yarn build-essential expect libcurl4-openssl-dev",
		}
		b.errmsg = []string{
			"Unable to obtain the gpg key for Yarn",
			"Unable to add yard repo as an apt source",
			"Unable to update apt database",
			"Unable to install sudo",
			"Unable to install nodejs",
			"Installing OS packages with apt failed",
		}
		b.hard = []bool{
			true,
			true,
			true,
			true,
			true,
			true,
		}
	}
	return
}

// prepDjango(d, &osTarget)
func prepDjango(d *gdjDefault, o *targetOS) {
	// Prep OS (user, virtualenv, chownership)
	d.sectionMsg("Preparing the OS for DefectDojo installation")
	prepCmds := osCmds{}
	osPrep(d, o.id, &prepCmds)
	// Run the OS Prep commands
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Preparing the OS for DefectDojo..."
	d.spin.Start()
	for i := range prepCmds.cmds {
		sendCmd(d,
			d.cmdLogger,
			prepCmds.cmds[i],
			prepCmds.errmsg[i],
			prepCmds.hard[i])
	}
	d.spin.Stop()
	d.statusMsg("Preparing the OS complete")
}

// osPrep
func osPrep(d *gdjDefault, id string, cmds *osCmds) {
	switch id {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		ubuntuOSPrep(d, id, cmds)
	}
	return
}

// ubuntuOSPrep
func ubuntuOSPrep(d *gdjDefault, id string, b *osCmds) {
	// Setup virutalenv, setup OS User, and chown DefectDojo app root to the dojo user
	switch id {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		b.id = id
		b.cmds = []string{
			"python3 -m virtualenv --python=/usr/bin/python3 " + d.conf.Install.Root,
			d.conf.Install.Root + "/bin/python3 -m pip install --upgrade pip",
			d.conf.Install.Root + "/bin/pip3 install -r " + d.conf.Install.Root + "/django-DefectDojo/requirements.txt",
			"mkdir " + d.conf.Install.Root + "/logs",
			"/usr/sbin/groupadd -f " + d.conf.Install.OS.Group, // TODO: check with os.Group.Lookup
			"id " + d.conf.Install.OS.User + " &>/dev/null; if [ $? -ne 0 ]; then useradd -s /bin/bash -m -g " +
				d.conf.Install.OS.Group + " " + d.conf.Install.OS.User + "; fi", // TODO: check with os.User.Lookup
			"chown -R " + d.conf.Install.OS.User + "." + d.conf.Install.OS.Group + " " + d.conf.Install.Root,
		}
		b.errmsg = []string{
			"Unable to setup virtualenv for DefectDojo",
			"Unable to update pip to latest",
			"Unable to install Python3 modules for DefectDojo",
			"Unable to create a directory for logs",
			"Unable to create a group for DefectDojo OS user",
			"Unable to create an OS user for DefectDojo",
			"Unable to change ownership of the DefectDojo app root directory",
		}
		b.hard = []bool{
			true,
			true,
			true,
			true,
			true,
			true,
			true,
		}
	}

	return
}

// createSettings
func createSettings(d *gdjDefault, o *targetOS) {
	// Create settings.py for DefectDojo
	d.sectionMsg("Creating settings.py for DefectDojo")
	settCmds := osCmds{}
	createSettingsPy(d, o.id, &settCmds)
	// Run the commands to create settings.py
	// TODO: Write values to .env.prod file
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Creating settings.py for DefectDojo..."
	d.spin.Start()
	for i := range settCmds.cmds {
		sendCmd(d,
			d.cmdLogger,
			settCmds.cmds[i],
			settCmds.errmsg[i],
			settCmds.hard[i])
	}
	d.spin.Stop()
	d.statusMsg("Creating settings.py for DefectDojo complete")

}

// createSettingsPy
func createSettingsPy(d *gdjDefault, id string, cmds *osCmds) {
	// Setup the env.prod file used by settings.py

	// Create the database URL for the env file - https://github.com/kennethreitz/dj-database-url
	dbURL := ""
	switch d.conf.Install.DB.Engine {
	case "SQLite":
		// sqlite:///PATH
		dbURL = "sqlite:///defectdojo.db"
	case "MySQL":
		// mysql://USER:PASSWORD@HOST:PORT/NAME
		dbURL = "mysql://" + d.conf.Install.DB.User + ":" + d.conf.Install.DB.Pass + "@" + d.conf.Install.DB.Host + ":" +
			strconv.Itoa(d.conf.Install.DB.Port) + "/" + d.conf.Install.DB.Name
	case "PostgreSQL":
		// postgres://USER:PASSWORD@HOST:PORT/NAME
		dbURL = "postgres://" + d.conf.Install.DB.User + ":" + d.conf.Install.DB.Pass + "@" + d.conf.Install.DB.Host + ":" +
			strconv.Itoa(d.conf.Install.DB.Port) + "/" + d.conf.Install.DB.Name
	}

	// Setup env file for production
	genAndWriteEnv(d, dbURL)

	// Create a settings.py for Dojo to use
	cmds.id = id
	cmds.cmds = []string{
		"ln -s " + d.conf.Install.Root + "/django-DefectDojo/dojo/settings " +
			d.conf.Install.Root + "/customizations",
		"chown " + d.conf.Install.OS.User + "." + d.conf.Install.OS.Group + " " + d.conf.Install.Root +
			"/django-DefectDojo/dojo/settings/settings.py",
	}
	cmds.errmsg = []string{
		"Unable to create settings.py file",
		"Unable to change ownership of settings.py file",
	}
	cmds.hard = []bool{
		true,
		true,
	}

	return
}

// setupDefectDojo
func setupDefectDojo(d *gdjDefault, o *targetOS) {
	// Django/Python installs
	d.sectionMsg("Setting up Django for DefectDojo")
	setupDj := osCmds{}
	setupDjango(d, o.id, &setupDj)
	// Run the Django commands
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Setting up Django for DefectDojo..."
	d.spin.Start()
	for i := range setupDj.cmds {
		sendCmd(d,
			d.cmdLogger,
			setupDj.cmds[i],
			setupDj.errmsg[i],
			setupDj.hard[i])
	}
	d.spin.Stop()
	d.statusMsg("Setting up Django complete")
}

// setupDjango
func setupDjango(d *gdjDefault, id string, cmds *osCmds) {
	// Generate the commands to do the Django install
	switch id {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		ubuntuSetupDDjango(d, id, cmds)
	}
	return
}

// ubuntuSetupDDjango
func ubuntuSetupDDjango(d *gdjDefault, id string, b *osCmds) {
	// Setup expect script needed to set initial admin password
	d.traceMsg(fmt.Sprintf("Injecting file %s at %s", "setup-superuser.expect", d.conf.Install.Root+"/django-DefectDojo"))
	// Inject expect script to change admin password
	terr := injectFile(d, suExpect, d.conf.Install.Root+"/django-DefectDojo", 0755)
	if terr != nil {
		fmt.Println("SOMETHING BAD HAPPENED HERE")
		fmt.Printf("Error was: %+v\n", terr)
		os.Exit(1)
	}

	err := patchOMatic(d)
	if err != nil {
		d.traceMsg(fmt.Sprintf("patchOMatic failed with non-blocking error: %+v", err))
		d.traceMsg("A failure of patchOMatic may lead to a corrupt install - be warned")
	}

	// Django installs - migrations, create Django superuser
	// TODO: Remove this switch to simplify
	switch id {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		// Add commands to setup DefectDojo - migrations, super user,
		// removed - "cd " + inst.Root + "/django-DefectDojo && source ../bin/activate && python3 manage.py makemigrations --merge --noinput", "Initial makemgrations failed",
		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo && source ../bin/activate && python3 manage.py makemigrations dojo",
			"Failed during makemgration dojo", true)

		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo && source ../bin/activate && python3 manage.py migrate",
			"Failed during database migrate", true)

		// Ensure there's a value for email as the call will fail without one
		adminEmail := "default.user@defectdojo.org"
		if len(d.conf.Install.Admin.Email) > 0 {
			// If user configures an incorrect email, this will still fail but that's on them, not godojo
			adminEmail = d.conf.Install.Admin.Email
		}
		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo && source ../bin/activate && python3 manage.py createsuperuser --noinput --username=\""+
			d.conf.Install.Admin.User+"\" --email=\""+adminEmail+"\"",
			"Failed while creating DefectDojo superuser", true)

		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo && source ../bin/activate && "+
			d.conf.Install.Root+"/django-DefectDojo/setup-superuser.expect "+d.conf.Install.Admin.User+" \""+escSpCar(d.conf.Install.Admin.Pass)+"\"",
			"Failed while setting the password for the DefectDojo superuser", true)

		// Roles showed up in 2.x.x
		if onlyAfter(d, 2, 0, 0) {
			addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo && source ../bin/activate && python3 manage.py loaddata role",
				"Failed while the loading data for role", true)
		}

		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo && source ../bin/activate && python3 manage.py loaddata "+
			"system_settings initial_banner_conf product_type test_type development_environment benchmark_type "+
			"benchmark_category benchmark_requirement language_type objects_review regulation initial_surveys role",
			"Failed while the loading data for a default install", true)

		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo && source ../bin/activate && python3 manage.py migrate_textquestions",
			"Failed while the loading data for a default survey questions", true)

		// removed - "cd " + inst.Root + "/django-DefectDojo && source ../bin/activate && python3 manage.py import_surveys", "Failed while the running import_surveys",
		// removed - "cd " + inst.Root + "/django-DefectDojo && source ../bin/activate && python3 manage.py loaddata initial_surveys", "Failed while the loading data for initial_surveys",

		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo && source ../bin/activate && python3 manage.py buildwatson",
			"Failed while the running buildwatson", true)

		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo && source ../bin/activate && python3 manage.py installwatson",
			"Failed while the running installwatson", true)

		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo/components && yarn",
			"Failed while the running yarn", true)

		addCmd(b, "cd "+d.conf.Install.Root+"/django-DefectDojo/ && source ../bin/activate && python3 manage.py collectstatic --noinput",
			"Failed while the running collectstatic", true)

		addCmd(b, "chown -R "+d.conf.Install.OS.User+"."+d.conf.Install.OS.Group+" "+d.conf.Install.Root,
			"Unable to change ownership of the DefectDojo directory", true)
	}

	return
}

// injectFile
func injectFile(d *gdjDefault, n string, p string, mask fs.FileMode) error {
	// Extract embedded file
	f, err := embd.ReadFile(n)
	if err != nil {
		// Embeded file was not found.
		fmt.Println("Unable to extract embedded patch file")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Strip off embedded directory from filename
	name := strings.Replace(n, "embd/", "", 1)

	// Write the file to disk
	err = ioutil.WriteFile(p+"/"+name, f, mask)
	if err != nil {
		// File can't be written
		return err
	}

	d.traceMsg(fmt.Sprintf("Wrote file %s at %s", name, p))

	return nil
}

func patchOMatic(d *gdjDefault) error {
	// If a source or commit install, do no patching
	if d.conf.Install.SourceInstall {
		return nil
	}

	// Check the install version for any needed patches
	switch d.conf.Install.Version {
	case "1.15.1":
		// Replace dojo/tools/factory to work around bug in Python 3.8 - https://bugs.python.org/issue44061
		_ = injectFile(d, factory2, d.conf.Install.Root+"/django-DefectDojo/dojo/tools", 755)
		_ = tryCmd(d,
			"mv -f "+d.conf.Install.Root+"/django-DefectDojo/dojo/tools/factory.py "+
				d.conf.Install.Root+"/django-DefectDojo/dojo/tools/factory_py.buggy",
			"Error renaming factory.py to factory_py.buggy", false)
		_ = tryCmd(d,
			"mv -f "+d.conf.Install.Root+"/django-DefectDojo/dojo/tools/factory_2.0.3 "+
				d.conf.Install.Root+"/django-DefectDojo/dojo/tools/factory.py",
			"Error replacing factory.py with updated one from version 2.0.3", false)
	}

	return nil
}

//onlyAfter
func onlyAfter(d *gdjDefault, major int, minor int, patch int) bool {
	// Split up version
	vBits := strings.Split(d.ver, ".")
	if len(vBits) != 3 {
		d.traceMsg(fmt.Sprintf("Bad version string: %s sent to onlyAfter()", d.ver))
		return false
	}

	// Convert version bits
	vMaj, _ := strconv.Atoi(vBits[0])
	vMin, _ := strconv.Atoi(vBits[1])
	vPat, _ := strconv.Atoi(vBits[2])

	//
	if vMaj < major {
		return false
	}
	if vMin < minor {
		return false
	}
	if vPat < patch {
		return false
	}

	return true
}
