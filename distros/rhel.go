package distros

import (
	"fmt"
	"strings"

	c "github.com/mtesauro/commandeer"
)

// Slice of Target structs supported RHEL Install Targets
var rhelReleases = []c.Target{
	{
		ID:      "RHEL:8",
		Distro:  "RHEL",
		Release: "8",
		OS:      "Linux",
		Shell:   "bash",
	},
	{
		ID:      "RHEL:9",
		Distro:  "RHEL",
		Release: "9",
		OS:      "Linux",
		Shell:   "bash",
	},
}

// Commands for RHEL
func GetRHEL(bc *c.CmdPkg, t string) error {
	// Use the label and target to get the correct commands
	switch {
	case bc.Label == "bootstrap":
		err := getRHELBootstrap(bc, t)
		if err != nil {
			// Return error from getRHELBootstrap()
			return err
		}
	case bc.Label == "installerprep":
		err := getRHELInstallerPrep(bc, t)
		if err != nil {
			// Return error from getRHELInstallerPrep()
			return err
		}
	case bc.Label == "prepdjango":
		err := getRHELPrepDjango(bc, t)
		if err != nil {
			// Return error from getRHELInstallerPrep()
			return err
		}
	case bc.Label == "createsettings":
		err := getRHELCreateSettings(bc, t)
		if err != nil {
			// Return error from getRHELCreateSettings()
			return err
		}
	case bc.Label == "setupdojo":
		err := getRHELSetupDojo(bc, t)
		if err != nil {
			// Return error from getRHELCreateSettings()
			return err
		}
	default:
		return fmt.Errorf("Unable to find a set of commands for the label %s\n", bc.Label)
	}

	return nil
}

func GetRHELDB(bc *c.CmdPkg, t string, d string) error {
	// Use the label and target to get the correct commands
	switch {
	case bc.Label == "installdb":
		// Determine target DB
		switch {
		case strings.ToLower(d) == "mysql":
			err := getRHELInstallMySQL(bc, t)
			if err != nil {
				// Return error from getRHELInstallMySQL()
				return err
			}
		case strings.ToLower(d) == "postgresql":
			err := getRHELInstallPostgres(bc, t)
			if err != nil {
				// Return error from getRHELInstallPostgres()
				return err
			}
		default:
			return fmt.Errorf("Unable to find a set of commands for the database %s\n", d)
		}
	case bc.Label == "startdb":
		// Determine target DB
		switch {
		case strings.ToLower(d) == "mysql":
			err := getRHELStartMySQL(bc, t)
			if err != nil {
				// Return error from getRHELInstallMySQL()
				return err
			}
		case strings.ToLower(d) == "postgresql":
			err := getRHELStartPostgres(bc, t)
			if err != nil {
				// Return error from getRHELInstallPostgres()
				return err
			}
		default:
			return fmt.Errorf("Unable to find commands to start the database %s\n", d)
		}
	case bc.Label == "installdbclient":
		// Determine target DB
		switch {
		case strings.ToLower(d) == "mysql":
			err := getRHELInstallMySQLClient(bc, t)
			if err != nil {
				// Return error from getRHELInstallMySQLClient()
				return err
			}
		case strings.ToLower(d) == "postgresql":
			err := getRHELInstallPgClient(bc, t)
			if err != nil {
				// Return error from getRHELInstallPostgres()
				return err
			}
		default:
			return fmt.Errorf("Unable to find commands to start the database %s\n", d)
		}
	default:
		return fmt.Errorf("Unable to find a set of commands for the label %s\n", bc.Label)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
//                           Bootstrap commands                              //
///////////////////////////////////////////////////////////////////////////////

func setRHELBootstrap() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8Bootstrap
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9Bootstrap
		}
	}
}

func getRHELBootstrap(bc *c.CmdPkg, t string) error {
	// Set bootstrap as the commands to use
	setRHELBootstrap()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// RHEL 8 Bootstrap commands
var rhel8Bootstrap = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "dnf check-update || [ $? -eq 100 ]", // WTF, dnf returns a 100 exit code if this command is successful!!
		Errmsg:     "Unable to update RHEL package database",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "dnf update -y",
		Errmsg:     "Unable to upgrade OS packages with dnf",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "dnf install -y python39 python3-virtualenv ca-certificates curl gnupg git sudo",
		Errmsg:     "Unable to install prerequisites for installer via dnf",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9Bootstrap = append([]c.SingleCmd{}, rhel8Bootstrap...)

///////////////////////////////////////////////////////////////////////////////
//                           Installer Prep commands                         //
///////////////////////////////////////////////////////////////////////////////

func setRHELInstallerPrep() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8InstallerPrep
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9InstallerPrep
		}
	}
}

func getRHELInstallerPrep(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setRHELInstallerPrep()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// RHEL 8 installer prep Commands
var rhel8InstallerPrep = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "curl --silent --location https://dl.yarnpkg.com/rpm/yarn.repo | sudo tee /etc/yum.repos.d/yarn.repo",
		Errmsg:     "Unable to add the repo for Yarn",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "curl --silent --location https://rpm.nodesource.com/setup_18.x | sudo bash -",
		Errmsg:     "Unable to add yard repo as an apt source",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "dnf check-update || [ $? -eq 100 ]", // WTF, dnf returns a 100 exit code if this command is successful!!
		Errmsg:     "Unable to update RHEL package database",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "dnf install -y sudo mysql yarn expect gcc python39-devel python39-pip initscripts mariadb-connector-c-devel libcurl-devel",
		Errmsg:     "Unable to install RHEL packages needed to prep the installer",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9InstallerPrep = append([]c.SingleCmd{}, rhel8InstallerPrep...)

///////////////////////////////////////////////////////////////////////////////
//                           Install MySQL commands                          //
///////////////////////////////////////////////////////////////////////////////

func setRHELInstallMySQL() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8NoDBMySQL
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9NoDBMySQL
		}
	}
}

func getRHELInstallMySQL(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setRHELInstallMySQL()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands to install MySQL for target %s\n", t)
}

// RHEL 8 install MySQL Commands
// TODO: https://computingforgeeks.com/install-mysql-5-7-on-centos-rhel-linux/
var rhel8NoDBMySQL = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "echo 'CURRENTLY UNSUPPORTED' && false",
		Errmsg:     "Unable to install MySQL",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9NoDBMySQL = append([]c.SingleCmd{}, rhel8NoDBMySQL...)

///////////////////////////////////////////////////////////////////////////////
//                           Install Postgres commands                       //
///////////////////////////////////////////////////////////////////////////////

func setRHELInstallPostgres() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8NoDBPostgres
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9NoDBPostgres
		}
	}
}

func getRHELInstallPostgres(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setRHELInstallPostgres()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands to install PostgreSQL for target %s\n", t)
}

// RHEL 8 install Postgres Commands
var rhel8NoDBPostgres = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "dnf module enable -y postgresql:13",
		Errmsg:     "Unable to enable install of PostgreSQL 13",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "dnf install -y postgresql-server",
		Errmsg:     "Unable to install PostgreSQL 13",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "postgresql-setup --initdb",
		Errmsg:     "Unable to initialize PostgreSQL 13",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9NoDBPostgres = append([]c.SingleCmd{}, rhel8NoDBPostgres...)

///////////////////////////////////////////////////////////////////////////////
//                           Install MySQL client commands                //
///////////////////////////////////////////////////////////////////////////////

func setRHELInstallMySQLClient() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			//rhelReleases[k].PkgCmds = rhel8InstMySQLClient
		case rhelReleases[k].Release == "9":
			//rhelReleases[k].PkgCmds = rhel9InstMySQLClient
		}
	}
}

func getRHELInstallMySQLClient(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setRHELInstallMySQLClient()

	// No match for the target provided
	//return fmt.Errorf("Unable to find commands for target %s\n", t)
	return fmt.Errorf("Commands for target %s have not been implemented\n", t)
}

///////////////////////////////////////////////////////////////////////////////
//                           Install Postgres client commands                //
///////////////////////////////////////////////////////////////////////////////

func setRHELInstallPgClient() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8InstPgClient
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9InstPgClient
		}
	}
}

func getRHELInstallPgClient(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setRHELInstallPgClient()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// RHEL 8 install Postgres client Commands
var rhel8InstPgClient = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "dnf module enable -y postgresql:13 && dnf install -y postgresql",
		Errmsg:     "Unable to install PostgreSQL client",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "/usr/sbin/groupadd -f postgres",
		Errmsg:     "Unable to add postgres group",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "id postgres &>/dev/null; if [ $? -ne 0 ]; then useradd -s /bin/bash -m -g postgres postgres; fi",
		Errmsg:     "Unable to add postgres user",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "mkdir -p /var/lib/pgsql",
		Errmsg:     "Unable to create postgres user directory",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9InstPgClient = append([]c.SingleCmd{}, rhel8InstPgClient...)

///////////////////////////////////////////////////////////////////////////////
//                           Start MySQL commands                            //
///////////////////////////////////////////////////////////////////////////////

func setRHELStartMySQL() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8StartMySQL
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9StartMySQL
		}
	}
}

func getRHELStartMySQL(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setRHELStartMySQL()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// RHEL 8 Start MySQL Commands
var rhel8StartMySQL = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "service mysql start && false",
		Errmsg:     "Unable to start MySQL server",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9StartMySQL = append([]c.SingleCmd{}, rhel8StartMySQL...)

///////////////////////////////////////////////////////////////////////////////
//                           Start Postgres commands                         //
///////////////////////////////////////////////////////////////////////////////

func setRHELStartPostgres() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8StartPostgres
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9StartPostgres
		}
	}
}

func getRHELStartPostgres(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setRHELStartPostgres()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// RHEL 8 Start Postgres Commands
var rhel8StartPostgres = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "systemctl start postgresql",
		Errmsg:     "Unable to start PostgreSQL",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9StartPostgres = append([]c.SingleCmd{}, rhel8StartPostgres...)

///////////////////////////////////////////////////////////////////////////////
//                           Prep Django commands                            //
///////////////////////////////////////////////////////////////////////////////

func setRHELPrepDjango() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8PrepDjango
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9PrepDjango
		}
	}
}

func getRHELPrepDjango(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setRHELPrepDjango()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// RHEL 8 Prep Django Commands
var rhel8PrepDjango = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "python3.9 -m pip install virtualenv",
		Errmsg:     "Unable to install virtualenv module for DefectDojo",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "python3.9 -m virtualenv --python=/usr/bin/python3.9 {conf.Install.Root}",
		Errmsg:     "Unable to create virtualenv for DefectDojo",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "{conf.Install.Root}/bin/python3 -m pip install --upgrade pip",
		Errmsg:     "Upgrade of Python pip failed",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "{conf.Install.Root}/bin/pip3 install -r {conf.Install.Root}/django-DefectDojo/requirements.txt",
		Errmsg:     "Unable to install Python3 modules for DefectDojo",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "mkdir {conf.Install.Root}/logs",
		Errmsg:     "Unable to create a directory for logs",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "/usr/sbin/groupadd -f {conf.Install.OS.Group}",
		Errmsg:     "Unable to create a group for DefectDojo OS user",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd: "id {conf.Install.OS.User} &>/dev/null; if [ $? -ne 0 ]; then useradd -s /bin/bash -m -g " +
			"{conf.Install.OS.Group} {conf.Install.OS.User}; fi",
		Errmsg:     "Unable to create an OS user for DefectDojo",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "chown -R {conf.Install.OS.User}.{conf.Install.OS.Group} {conf.Install.Root}",
		Errmsg:     "",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9PrepDjango = append([]c.SingleCmd{}, rhel8PrepDjango...)

///////////////////////////////////////////////////////////////////////////////
//                          Create Settings commands                         //
///////////////////////////////////////////////////////////////////////////////

func setRHELCreateSettings() {
	// Connect bootstrap commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8CreateSettings
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9CreateSettings
		}
	}
}

func getRHELCreateSettings(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setRHELCreateSettings()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// RHEL 8 Create Settings Commands
var rhel8CreateSettings = []c.SingleCmd{
	c.SingleCmd{
		Cmd: "ln -s {conf.Install.Root}/django-DefectDojo/dojo/settings/ " +
			"{conf.Install.Root}/customizations",
		Errmsg:     "Unable to create customization directory",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd: "echo '# Add customizations here\n# For more details see:" +
			" https://documentation.defectdojo.com/getting_started/configuration/' > {conf.Install.Root}/customizations/local_settings.py",
		Errmsg:     "Unable to change ownership of .env.prod file",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd: "chown {conf.Install.OS.User}.{conf.Install.OS.Group} {conf.Install.Root}" +
			"/django-DefectDojo/dojo/settings/.env.prod",
		Errmsg:     "Unable to change ownership of .env.prod file",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9CreateSettings = append([]c.SingleCmd{}, rhel8CreateSettings...)

///////////////////////////////////////////////////////////////////////////////
//                           Setup DefectDojo commands                       //
///////////////////////////////////////////////////////////////////////////////

func setRHELSetupDojo() {
	// Connect setup DefectDojo commands to the supported RHEL releases
	for k := range rhelReleases {
		switch {
		case rhelReleases[k].Release == "8":
			rhelReleases[k].PkgCmds = rhel8SetupDojo
		case rhelReleases[k].Release == "9":
			rhelReleases[k].PkgCmds = rhel9SetupDojo
		}
	}
}

func getRHELSetupDojo(bc *c.CmdPkg, t string) error {
	// Set setup DefectDojo as the commands to use
	setRHELSetupDojo()

	// Cycle through RHEL install targets
	for k, v := range rhelReleases {
		// Find a match for the target ID and the existing list of commands in rhelReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, rhelReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// RHEL 8 setup DefectDojo Commands
var rhel8SetupDojo = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && python3 manage.py makemigrations dojo",
		Errmsg:     "Failed during makemgration dojo",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && python3 manage.py migrate",
		Errmsg:     "Failed during database migrate",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd: "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && python3 manage.py createsuperuser" +
			" --noinput --username=\"{conf.Install.Admin.User}\" --email=\"{conf.Install.Admin.Email}\"",
		Errmsg:     "Failed while creating DefectDojo superuser",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd: "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && " +
			"{conf.Install.Root}/django-DefectDojo/setup-superuser.expect {conf.Install.Admin.User} \"{conf.Install.Admin.Pass}\"",
		Errmsg:     "Failed while setting the password for the DefectDojo superuser",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd: "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && python3 manage.py loaddata " +
			"system_settings initial_banner_conf product_type test_type development_environment benchmark_type " +
			"benchmark_category benchmark_requirement language_type objects_review regulation initial_surveys role",
		Errmsg:     "Failed while the loading data for a default install",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && python3 manage.py migrate_textquestions",
		Errmsg:     "Failed while the loading data for a default survey questions",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && python3 manage.py buildwatson",
		Errmsg:     "Failed while the running buildwatson",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && python3 manage.py installwatson",
		Errmsg:     "Failed while the running installwatson",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && python3 manage.py initialize_test_types",
		Errmsg:     "Failed to initialize test_types",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "cd {conf.Install.Root}/django-DefectDojo && source ../bin/activate && python3 manage.py initialize_permissions",
		Errmsg:     "Failed to initialize permissions",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "cd {conf.Install.Root}/django-DefectDojo/components && yarn",
		Errmsg:     "Failed while the running yarn",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "cd {conf.Install.Root}/django-DefectDojo/ && source ../bin/activate && python3 manage.py collectstatic --noinput",
		Errmsg:     "Failed while the running collectstatic",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "chown -R {conf.Install.OS.User}.{conf.Install.OS.Group} {conf.Install.Root}",
		Errmsg:     "Unable to change ownership of the DefectDojo directory",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for RHEL 9
var rhel9SetupDojo = append([]c.SingleCmd{}, rhel8SetupDojo...)
