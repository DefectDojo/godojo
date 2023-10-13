package distros

import (
	"fmt"
	"strings"

	c "github.com/mtesauro/commandeer"
)

///////////////////////////////////////////////////////////////////////////////
//              Template file to add commads for other distros               //
///////////////////////////////////////////////////////////////////////////////

// Slice of Target structs supported Template Install Targets
var templateReleases = []c.Target{
	{
		ID:      "Template:22.04",
		Distro:  "Template",
		Release: "22.04",
		OS:      "Linux",
		Shell:   "bash",
	},
	{
		ID:      "Template:21.04",
		Distro:  "Template",
		Release: "21.04",
		OS:      "Linux",
		Shell:   "bash",
	},
}

// Commands for Tempate
func GetTemplate(bc *c.CmdPkg, t string) error {
	// Use the label and target to get the correct commands
	switch {
	case bc.Label == "bootstrap":
		err := getTemplateBootstrap(bc, t)
		if err != nil {
			// Return error from getTemplateBootstrap()
			return err
		}
	case bc.Label == "installerprep":
		err := getTemplateInstallerPrep(bc, t)
		if err != nil {
			// Return error from getTemplateInstallerPrep()
			return err
		}
	case bc.Label == "prepdjango":
		err := getTemplatePrepDjango(bc, t)
		if err != nil {
			// Return error from getTemplateInstallerPrep()
			return err
		}
	case bc.Label == "createsettings":
		err := getTemplateCreateSettings(bc, t)
		if err != nil {
			// Return error from getTemplateCreateSettings()
			return err
		}
	case bc.Label == "setupdojo":
		err := getTemplateSetupDojo(bc, t)
		if err != nil {
			// Return error from getTemplateCreateSettings()
			return err
		}
	default:
		return fmt.Errorf("Unable to find a set of commands for the label %s\n", bc.Label)
	}

	return nil
}

func GetTemplateDB(bc *c.CmdPkg, t string, d string) error {
	// Use the label and target to get the correct commands
	switch {
	case bc.Label == "installdb":
		// Determine target DB
		switch {
		case strings.ToLower(d) == "mysql":
			err := getTemplateInstallMySQL(bc, t)
			if err != nil {
				// Return error from getTemplateInstallMySQL()
				return err
			}
		case strings.ToLower(d) == "postgresql":
			err := getTemplateInstallPostgres(bc, t)
			if err != nil {
				// Return error from getTemplateInstallPostgres()
				return err
			}
		default:
			return fmt.Errorf("Unable to find a set of commands for the database %s\n", d)
		}
	case bc.Label == "startdb":
		// Determine target DB
		switch {
		case strings.ToLower(d) == "mysql":
			err := getTemplateStartMySQL(bc, t)
			if err != nil {
				// Return error from getTemplateInstallMySQL()
				return err
			}
		case strings.ToLower(d) == "postgresql":
			err := getTemplateStartPostgres(bc, t)
			if err != nil {
				// Return error from getTemplateInstallPostgres()
				return err
			}
		default:
			return fmt.Errorf("Unable to find commands to start the database %s\n", d)
		}
	case bc.Label == "installdbclient":
		// Determine target DB
		switch {
		case strings.ToLower(d) == "mysql":
			err := getTemplateInstallMySQLClient(bc, t)
			if err != nil {
				// Return error from getTemplateInstallMySQLClient()
				return err
			}
		case strings.ToLower(d) == "postgresql":
			err := getTemplateInstallPgClient(bc, t)
			if err != nil {
				// Return error from getTemplateInstallPostgres()
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

func setTemplateBootstrap() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204Bootstrap
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104Bootstrap
		}
	}
}

func getTemplateBootstrap(bc *c.CmdPkg, t string) error {
	// Set bootstrap as the commands to use
	setTemplateBootstrap()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Template 22.04 Bootstrap commands
var t2204Bootstrap = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive tpl-get update",
		Errmsg:     "Unable to update tpl database",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive tpl-get -y upgrade",
		Errmsg:     "Unable to upgrade OS packages with tpl",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive tpl-get -y -o Dpkg::Options::=\"--force-confdef\" -o Dpkg::Options::=\"--force-confold\" install python3 python3-virtualenv ca-certificates curl gnupg git sudo",
		Errmsg:     "Unable to install prerequisites for installer via tpl",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Template 21.04
var t2104Bootstrap = append([]c.SingleCmd{}, t2204Bootstrap...)

///////////////////////////////////////////////////////////////////////////////
//                           Installer Prep commands                         //
///////////////////////////////////////////////////////////////////////////////

func setTemplateInstallerPrep() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204InstallerPrep
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104InstallerPrep
		}
	}
}

func getTemplateInstallerPrep(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setTemplateInstallerPrep()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Template 22.04 installer prep Commands
// TODO Check if the yarn command needs updating
var t2204InstallerPrep = []c.SingleCmd{
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

// No command changes needed for Template 21.04
var t2104InstallerPrep = append([]c.SingleCmd{}, t2204InstallerPrep...)

///////////////////////////////////////////////////////////////////////////////
//                           Install MySQL commands                          //
///////////////////////////////////////////////////////////////////////////////

func setTemplateInstallMySQL() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204NoDBMySQL
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104NoDBMySQL
		}
	}
}

func getTemplateInstallMySQL(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setTemplateInstallMySQL()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands to install MySQL for target %s\n", t)
}

// Template 22.04 install MySQL Commands
var t2204NoDBMySQL = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get install -y mysql-server libmysqlclient-dev",
		Errmsg:     "Unable to install MySQL",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Template 21.04
var t2104NoDBMySQL = append([]c.SingleCmd{}, t2204NoDBMySQL...)

///////////////////////////////////////////////////////////////////////////////
//                           Install Postgres commands                       //
///////////////////////////////////////////////////////////////////////////////

func setTemplateInstallPostgres() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204NoDBPostgres
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104NoDBPostgres
		}
	}
}

func getTemplateInstallPostgres(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setTemplateInstallPostgres()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands to install PostgreSQL for target %s\n", t)
}

// Template 22.04 install Postgres Commands
var t2204NoDBPostgres = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get install -y libpq-dev postgresql postgresql-contrib postgresql-client-common",
		Errmsg:     "Unable to install PostgreSQL",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Template 21.04
var t2104NoDBPostgres = append([]c.SingleCmd{}, t2204NoDBPostgres...)

///////////////////////////////////////////////////////////////////////////////
//                           Install MySQL client commands                //
///////////////////////////////////////////////////////////////////////////////

func setTemplateInstallMySQLClient() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			//templateReleases[k].PkgCmds = t2204InstMySQLClient
		case templateReleases[k].Release == "21.04":
			//templateReleases[k].PkgCmds = t2104InstMySQLClient
		}
	}
}

func getTemplateInstallMySQLClient(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setTemplateInstallMySQLClient()

	// No match for the target provided
	//return fmt.Errorf("Unable to find commands for target %s\n", t)
	return fmt.Errorf("Commands for target %s have not been implemented\n", t)
}

///////////////////////////////////////////////////////////////////////////////
//                           Install Postgres client commands                //
///////////////////////////////////////////////////////////////////////////////

func setTemplateInstallPgClient() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204InstPgClient
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104InstPgClient
		}
	}
}

func getTemplateInstallPgClient(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setTemplateInstallPgClient()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Template 22.04 install Postgres client Commands
var t2204InstPgClient = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "DEBIAN_FRONTEND=noninteractive apt-get install -y postgresql-client-12",
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

// No command changes needed for Template 21.04
var t2104InstPgClient = append([]c.SingleCmd{}, t2204InstPgClient...)

///////////////////////////////////////////////////////////////////////////////
//                           Start MySQL commands                            //
///////////////////////////////////////////////////////////////////////////////

func setTemplateStartMySQL() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204StartMySQL
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104StartMySQL
		}
	}
}

func getTemplateStartMySQL(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setTemplateStartMySQL()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Template 22.04 Start MySQL Commands
var t2204StartMySQL = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "service mysql start",
		Errmsg:     "Unable to start MariaDB",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Template 21.04
var t2104StartMySQL = append([]c.SingleCmd{}, t2204StartMySQL...)

///////////////////////////////////////////////////////////////////////////////
//                           Start Postgres commands                         //
///////////////////////////////////////////////////////////////////////////////

func setTemplateStartPostgres() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204StartPostgres
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104StartPostgres
		}
	}
}

func getTemplateStartPostgres(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setTemplateStartPostgres()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Template 22.04 Start Postgres Commands
var t2204StartPostgres = []c.SingleCmd{
	c.SingleCmd{
		Cmd:        "/usr/sbin/service postgresql start",
		Errmsg:     "Unable to start PostgreSQL",
		Hard:       true,
		Timeout:    0,
		BeforeText: "",
		AfterText:  "",
	},
}

// No command changes needed for Template 21.04
var t2104StartPostgres = append([]c.SingleCmd{}, t2204StartPostgres...)

///////////////////////////////////////////////////////////////////////////////
//                           Prep Django commands                            //
///////////////////////////////////////////////////////////////////////////////

func setTemplatePrepDjango() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204PrepDjango
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104PrepDjango
		}
	}
}

func getTemplatePrepDjango(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setTemplatePrepDjango()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Template 22.04 Prep Django Commands
var t2204PrepDjango = []c.SingleCmd{
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

// No command changes needed for Template 21.04
var t2104PrepDjango = append([]c.SingleCmd{}, t2204PrepDjango...)

///////////////////////////////////////////////////////////////////////////////
//                          Create Settings commands                         //
///////////////////////////////////////////////////////////////////////////////

func setTemplateCreateSettings() {
	// Connect bootstrap commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204CreateSettings
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104CreateSettings
		}
	}
}

func getTemplateCreateSettings(bc *c.CmdPkg, t string) error {
	// Set Installer Prep as the commands to use
	setTemplateCreateSettings()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Template 22.04 Create Settings Commands
var t2204CreateSettings = []c.SingleCmd{
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

// No command changes needed for Template 21.04
var t2104CreateSettings = append([]c.SingleCmd{}, t2204CreateSettings...)

///////////////////////////////////////////////////////////////////////////////
//                           Setup DefectDojo commands                       //
///////////////////////////////////////////////////////////////////////////////

func setTemplateSetupDojo() {
	// Connect setup DefectDojo commands to the supported Template releases
	for k := range templateReleases {
		switch {
		case templateReleases[k].Release == "22.04":
			templateReleases[k].PkgCmds = t2204SetupDojo
		case templateReleases[k].Release == "21.04":
			templateReleases[k].PkgCmds = t2104SetupDojo
		}
	}
}

func getTemplateSetupDojo(bc *c.CmdPkg, t string) error {
	// Set setup DefectDojo as the commands to use
	setTemplateSetupDojo()

	// Cycle through Template install targets
	for k, v := range templateReleases {
		// Find a match for the target ID and the existing list of commands in templateReleases
		if strings.Compare(
			strings.ToLower(v.ID),
			strings.ToLower(t)) == 0 {
			bc.Targets = append(bc.Targets, templateReleases[k])
			return nil
		}
	}

	// No match for the target provided
	return fmt.Errorf("Unable to find commands for target %s\n", t)
}

// Template 22.04 setup DefectDojo Commands
var t2204SetupDojo = []c.SingleCmd{
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

// No command changes needed for Template 21.04
var t2104SetupDojo = append([]c.SingleCmd{}, t2204SetupDojo...)

//// TEMPLATE

///////////////////////////////////////////////////////////////////////////////
//                           SOMETHING commands                         //
///////////////////////////////////////////////////////////////////////////////

// Template 22.04 [SOMETHING] Commands
//var t2204uskeleton = []c.SingleCmd{
//	c.SingleCmd{
//		Cmd:        "",
//		Errmsg:     "",
//		Hard:       true,
//		Timeout:    0,
//		BeforeText: "",
//		AfterText:  "",
//	},
//	c.SingleCmd{
//		Cmd:        "",
//		Errmsg:     "",
//		Hard:       true,
//		Timeout:    0,
//		BeforeText: "",
//		AfterText:  "",
//	},
//}
//
//// No command changes needed for Template 21.04
//var t2104uskeleton = append([]c.SingleCmd{}, t2204uskeleton...)
