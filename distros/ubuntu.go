package distros

import (
	"fmt"
	"strings"

	c "github.com/mtesauro/commandeer"
)

// Slice of Target structs supported Ubuntu Install Targets
var ubuntuReleases = []c.Target{
	{
		ID:      "Ubuntu:22.04",
		Distro:  "Ubuntu",
		Release: "22.04",
		OS:      "Linux",
		Shell:   "bash",
	},
	{
		ID:      "Ubuntu:21.04",
		Distro:  "Ubuntu",
		Release: "21.04",
		OS:      "Linux",
		Shell:   "bash",
	},
}

// Commands for Ubuntu
func GetUbuntu(bc *c.CmdPkg, t string) error {
	// Use the label and target to get the correct commands
	switch {
	case bc.Label == "bootstrap":
		err := getUbuntuBootstrap(bc, t)
		if err != nil {
			// Return error from getUbuntuBootstrap()
			return err
		}
	case bc.Label == "installerprep":
		err := getUbuntuInstallerPrep(bc, t)
		if err != nil {
			// Return error from getUbuntuInstallerPrep()
			return err
		}
	case bc.Label == "prepdjango":
		err := getUbuntuPrepDjango(bc, t)
		if err != nil {
			// Return error from getUbuntuInstallerPrep()
			return err
		}
	case bc.Label == "createsettings":
		err := getUbuntuCreateSettings(bc, t)
		if err != nil {
			// Return error from getUbuntuCreateSettings()
			return err
		}
	case bc.Label == "setupdojo":
		err := getUbuntuSetupDojo(bc, t)
		if err != nil {
			// Return error from getUbuntuCreateSettings()
			return err
		}
	default:
		return fmt.Errorf("Unable to find a set of commands for the label %s\n", bc.Label)
	}

	return nil
}

func GetUbuntuDB(bc *c.CmdPkg, t string, d string) error {
	// Use the label and target to get the correct commands
	switch {
	case bc.Label == "installdb":
		// Determine target DB
		switch {
		case strings.ToLower(d) == "mysql":
			err := getUbuntuInstallMySQL(bc, t)
			if err != nil {
				// REturn error from getUbuntuInstallMySQL()
				return err
			}
		case strings.ToLower(d) == "postgresql":
			err := getUbuntuInstallPostgres(bc, t)
			if err != nil {
				// REturn error from getUbuntuInstallPostgres()
				return err
			}
		default:
			return fmt.Errorf("Unable to find a set of commands for the database %s\n", d)
		}
	case bc.Label == "startdb":
		// Determine target DB
		switch {
		case strings.ToLower(d) == "mysql":
			err := getUbuntuStartMySQL(bc, t)
			if err != nil {
				// Return error from getUbuntuInstallMySQL()
				return err
			}
		case strings.ToLower(d) == "postgresql":
			err := getUbuntuStartPostgres(bc, t)
			if err != nil {
				// Return error from getUbuntuInstallPostgres()
				return err
			}
		default:
			return fmt.Errorf("Unable to find commands to start the database %s\n", d)
		}
	case bc.Label == "installdbclient":
		// Determine target DB
		switch {
		case strings.ToLower(d) == "mysql":
			err := getUbuntuInstallMySQLClient(bc, t)
			if err != nil {
				// Return error from getUbuntuInstallMySQLClient()
				return err
			}
		case strings.ToLower(d) == "postgresql":
			err := getUbuntuInstallPgClient(bc, t)
			if err != nil {
				// Return error from getUbuntuInstallPostgres()
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

func setUbuntuBootstrap() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204Bootstrap
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104Bootstrap
		}
	}
}

func getUbuntuBootstrap(bc *c.CmdPkg, t string) error {
	// Set bootstrap as the commands to use
	setUbuntuBootstrap()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Ubuntu 22.04 Bootstrap commands
var u2204Bootstrap = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get update",
		Errmsg:     "Unable to update apt database",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get -y upgrade",
		Errmsg:     "Unable to upgrade OS packages with apt",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::=\"--force-confdef\" -o Dpkg::Options::=\"--force-confold\" install python3 python3-virtualenv ca-certificates curl gnupg git sudo",
		Errmsg:     "Unable to install prerequisites for installer via apt",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Ubuntu 21.04
var u2104Bootstrap = append([]c.SingleCmd{}, u2204Bootstrap...)

///////////////////////////////////////////////////////////////////////////////
//                           Installer Prep commands                         //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuInstallerPrep() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204InstallerPrep
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104InstallerPrep
		}
	}
}

func getUbuntuInstallerPrep(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setUbuntuInstallerPrep()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Ubuntu 22.04 installer prep Commands
// TODO Check if the yarn command needs updating
var u2204InstallerPrep = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "curl -sS {yarnGPG} | apt-key add -",
		Errmsg:     "Unable to obtain the gpg key for Yarn",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "echo -n {yarnRepo} > /etc/apt/sources.list.d/yarn.list",
		Errmsg:     "Unable to add yard repo as an apt source",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get update",
		Errmsg:     "Unable to update apt database",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get -y install sudo libmysqlclient-dev",
		Errmsg:     "Unable to install sudo and MySQL client library",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "curl -sL {nodeURL} | bash - ",
		Errmsg:     "Unable to install nodejs",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get install -y apt-transport-https libjpeg-dev gcc libssl-dev python3-dev python3-pip python3-virtualenv yarn build-essential expect libcurl4-openssl-dev",
		Errmsg:     "Installing OS packages with apt failed",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Ubuntu 21.04
var u2104InstallerPrep = append([]c.SingleCmd{}, u2204InstallerPrep...)

///////////////////////////////////////////////////////////////////////////////
//                           Install MySQL commands                          //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuInstallMySQL() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204NoDBMySQL
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104NoDBMySQL
		}
	}
}

func getUbuntuInstallMySQL(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setUbuntuInstallMySQL()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands to install MySQL for target %s\n", t)
}

// Ubuntu 22.04 install MySQL Commands
var u2204NoDBMySQL = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get install -y mysql-server libmysqlclient-dev",
		Errmsg:     "Unable to install MySQL",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Ubuntu 21.04
var u2104NoDBMySQL = append([]c.SingleCmd{}, u2204NoDBMySQL...)

///////////////////////////////////////////////////////////////////////////////
//                           Install Postgres commands                       //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuInstallPostgres() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204NoDBPostgres
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104NoDBPostgres
		}
	}
}

func getUbuntuInstallPostgres(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setUbuntuInstallPostgres()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands to install PostgreSQL for target %s\n", t)
}

// Ubuntu 22.04 install Postgres Commands
var u2204NoDBPostgres = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get install -y libpq-dev postgresql postgresql-contrib postgresql-client-common",
		Errmsg:     "Unable to install PostgreSQL",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Ubuntu 21.04
var u2104NoDBPostgres = append([]c.SingleCmd{}, u2204NoDBPostgres...)

///////////////////////////////////////////////////////////////////////////////
//                           Install MySQL client commands                //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuInstallMySQLClient() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			//ubuntuReleases[k].PkgCmds = u2204InstMySQLClient
		case ubuntuReleases[k].Release == "21.04":
			//ubuntuReleases[k].PkgCmds = u2104InstMySQLClient
		}
	}
}

func getUbuntuInstallMySQLClient(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setUbuntuInstallMySQLClient()

	// No match for the target provided
	//return fmt.Errorf("Unable to find commands for target %s\n", t)
	return fmt.Errorf("Commands for target %s have not been implemented\n", t)
}

///////////////////////////////////////////////////////////////////////////////
//                           Install Postgres client commands                //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuInstallPgClient() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204InstPgClient
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104InstPgClient
		}
	}
}

func getUbuntuInstallPgClient(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setUbuntuInstallPgClient()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Ubuntu 22.04 install Postgres client Commands
var u2204InstPgClient = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get install -y postgresql-client-14",
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
		Cmd:        "/usr/sbin/useradd -s /bin/bash -m -g postgres postgres",
		Errmsg:     "Unable to add postgres user",
		Hard:       false, // incase there is an existing postgres user, useradd returns a 9 exit code
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Ubuntu 21.04
var u2104InstPgClient = append([]c.SingleCmd{}, u2204InstPgClient...)

///////////////////////////////////////////////////////////////////////////////
//                           Start MySQL commands                            //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuStartMySQL() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204StartMySQL
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104StartMySQL
		}
	}
}

func getUbuntuStartMySQL(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setUbuntuStartMySQL()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Ubuntu 22.04 Start MySQL Commands
var u2204StartMySQL = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "service mysql start",
		Errmsg:     "Unable to start MariaDB",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Ubuntu 21.04
var u2104StartMySQL = append([]c.SingleCmd{}, u2204StartMySQL...)

///////////////////////////////////////////////////////////////////////////////
//                           Start Postgres commands                         //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuStartPostgres() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204StartPostgres
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104StartPostgres
		}
	}
}

func getUbuntuStartPostgres(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setUbuntuStartPostgres()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Ubuntu 22.04 Start Postgres Commands
var u2204StartPostgres = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "/usr/sbin/service postgresql start",
		Errmsg:     "Unable to start PostgreSQL",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Ubuntu 21.04
var u2104StartPostgres = append([]c.SingleCmd{}, u2204StartPostgres...)

///////////////////////////////////////////////////////////////////////////////
//                           Prep Django commands                            //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuPrepDjango() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204PrepDjango
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104PrepDjango
		}
	}
}

func getUbuntuPrepDjango(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setUbuntuPrepDjango()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Ubuntu 22.04 Prep Django Commands
var u2204PrepDjango = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "python3 -m virtualenv --python=/usr/bin/python3 {conf.Install.Root}",
		Errmsg:     "Unable to setup virtualenv for DefectDojo",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "{conf.Install.Root}/bin/python3 -m pip install --upgrade pip",
		Errmsg:     "",
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

// No command changes needed for Ubuntu 21.04
var u2104PrepDjango = append([]c.SingleCmd{}, u2204PrepDjango...)

///////////////////////////////////////////////////////////////////////////////
//                          Create Settings commands                         //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuCreateSettings() {
	// Connect bootstrap commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204CreateSettings
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104CreateSettings
		}
	}
}

func getUbuntuCreateSettings(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setUbuntuCreateSettings()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Ubuntu 22.04 Create Settings Commands
var u2204CreateSettings = []c.SingleCmd{
	c.SingleCmd{
		Cmd: "ln -s {conf.Install.Root}/django-DefectDojo/dojo/settings/ " +
			"{conf.Install.Root}/customizations",
		Errmsg:     "Unable to create settings.py file",
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
			"/django-DefectDojo/dojo/settings/settings.py",
		Errmsg:     "Unable to change ownership of settings.py file",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Ubuntu 21.04
var u2104CreateSettings = append([]c.SingleCmd{}, u2204CreateSettings...)

///////////////////////////////////////////////////////////////////////////////
//                           Setup DefectDojo commands                       //
///////////////////////////////////////////////////////////////////////////////

func setUbuntuSetupDojo() {
	// Connect setup DefectDojo commands to the supported Ubuntu releases
	for k := range ubuntuReleases {
		switch {
		case ubuntuReleases[k].Release == "22.04":
			ubuntuReleases[k].PkgCmds = u2204SetupDojo
		case ubuntuReleases[k].Release == "21.04":
			ubuntuReleases[k].PkgCmds = u2104SetupDojo
		}
	}
}

func getUbuntuSetupDojo(bc *c.CmdPkg, t string) error {
	// Set setup DefectDojo as the commands to use
	setUbuntuSetupDojo()

	// Cycle through Ubuntu install targets
	for k, v := range ubuntuReleases {
		// Find a match for the target ID and the existing list of commands in ubuntuReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, ubuntuReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Ubuntu 22.04 setup DefectDojo Commands
var u2204SetupDojo = []c.SingleCmd{
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

// No command changes needed for Ubuntu 21.04
var u2104SetupDojo = append([]c.SingleCmd{}, u2204SetupDojo...)
