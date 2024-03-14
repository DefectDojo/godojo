package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/defectdojo/godojo/distros"
	c "github.com/mtesauro/commandeer"
)

// Setup a struct to use for DB commands
type sqlStr struct {
	os     string
	sql    string
	errMsg string
	creds  map[string]string
	kind   string
}

// saneDBConfig checks if the options configured in dojoConfig.yml are
// possible aka sane and will exist the installer with a message if they are not
func saneDBConfig(d *DDConfig) {
	// Remote database that doesn't exist - godojo can't help you here
	if !d.conf.Install.DB.Local && !d.conf.Install.DB.Exists {
		d.errorMsg("Remote database which doens't exist was confgiured in dojoConfig.yml.")
		d.errorMsg("This is an unsupported configuration.")
		d.statusMsg("Correct configuration and/or install a remote DB before running installer again.")
		fmt.Printf("Exiting...\n\n")
		os.Exit(1)
	}
}

// prepDBForDojo
func installDBForDojo(d *DDConfig, t *targetOS) {
	// Handle the case that the DB is local and doesn't exist
	if !d.conf.Install.DB.Exists {
		// Note that godojo won't try to install remote databases
		dbNotExist(d, t)
	}

	// Install DB clients for remote DBs
	if !d.conf.Install.DB.Local {
		dbClient(d, t)
	}

	// Start the database if local and didn't already exist
	if d.conf.Install.DB.Local && !d.conf.Install.DB.Exists {
		localDBStart(d, t)
	}

}

// dbNotExist takes a pointer to a DDConfig struct and a pointer to targetOS
// struct and runs the commands necesary to install a local database of the
// supported type (PostgreSQL, MySQL, etc)
func dbNotExist(d *DDConfig, t *targetOS) {
	// Handle the case that the DB is local and doesn't exist
	d.sectionMsg("Installing database needed for DefectDojo")

	// Create a new install DB command package
	cInstallDB := c.NewPkg("installdb")

	// Get commands for the right distro & DB
	switch {
	case t.distro == "ubuntu":
		d.traceMsg("DB needs to be installed on Ubuntu")
		err := distros.GetUbuntuDB(cInstallDB, t.id, d.conf.Install.DB.Engine)
		if err != nil {
			fmt.Printf("Error searching for commands to install DB on target OS %s was\n", t.id)
			fmt.Printf("\t%+v\n", err)
			os.Exit(1)
		}
		if strings.ToLower(d.conf.Install.DB.Engine) == "mysql" {
			d.warnMsg("WARNING: While supported, there is significantly more testing with PostreSQL than MySQL. YMMV.")
		}
	case t.distro == "rhel":
		d.traceMsg("DB needs to be installed on RHEL")
		err := distros.GetRHELDB(cInstallDB, t.id, d.conf.Install.DB.Engine)
		if err != nil {
			fmt.Printf("Error searching for commands to install DB on target OS %s was\n", t.id)
			fmt.Printf("\t%+v\n", err)
			os.Exit(1)
		}
		if strings.ToLower(d.conf.Install.DB.Engine) == "mysql" {
			d.warnMsg("WARNING: While supported, there is significantly more testing with PostreSQL than MySQL. YMMV.")
		}
	default:
		d.traceMsg(fmt.Sprintf("Distro identified (%s) is not supported", t.id))
		fmt.Printf("Distro identified by godojo (%s) is not supported, exiting...\n", t.id)
		os.Exit(1)
	}

	// Run the commands to install the chosen DB
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Installing " + d.conf.Install.DB.Engine + " database for DefectDojo..."
	d.spin.Start()
	// Run the install DB for the target OS
	tCmds, err := distros.CmdsForTarget(cInstallDB, t.id)
	if err != nil {
		fmt.Printf("Error getting commands to install DB target OS %s\n", t.id)
		os.Exit(1)
	}

	for i := range tCmds {
		sendCmd(d,
			d.cmdLogger,
			tCmds[i].Cmd,
			tCmds[i].Errmsg,
			tCmds[i].Hard)
	}
	d.spin.Stop()
	d.statusMsg("Installing Database complete")
}

// dbClientInstall
func dbClient(d *DDConfig, t *targetOS) {
	// Handle the case that the DB is local and doesn't exist
	d.sectionMsg("Installing database client needed for DefectDojo")

	// Create a new install DB client command package
	cInstallDBClient := c.NewPkg("installdbclient")

	// Get the commands for the right distro & DB
	switch {
	case t.distro == "ubuntu":
		d.traceMsg("DB client needs to be installed on Ubuntu")
		err := distros.GetUbuntuDB(cInstallDBClient, t.id, d.conf.Install.DB.Engine)
		if err != nil {
			fmt.Printf("Error searching for commands to install DB client on target OS %s was\n", t.id)
			fmt.Printf("\t%+v\n", err)
			os.Exit(1)
		}
	case t.distro == "rhel":
		d.traceMsg("DB client needs to be installed on RHEL")
		err := distros.GetRHELDB(cInstallDBClient, t.id, d.conf.Install.DB.Engine)
		if err != nil {
			fmt.Printf("Error searching for commands to install DB client on target OS %s was\n", t.id)
			fmt.Printf("\t%+v\n", err)
			os.Exit(1)
		}
	default:
		d.traceMsg(fmt.Sprintf("Distro identified (%s) is not supported", t.id))
		fmt.Printf("Distro identified by godojo (%s) is not supported, exiting...\n", t.id)
		os.Exit(1)
	}

	// Run the commands to install the chosen DB
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Installing " + d.conf.Install.DB.Engine + " database client for DefectDojo..."
	d.spin.Start()
	// Run the install DB client for the target OS
	tCmds, err := distros.CmdsForTarget(cInstallDBClient, t.id)
	if err != nil {
		fmt.Printf("Error getting commands to install DB target OS %s\n", t.id)
		os.Exit(1)
	}

	for i := range tCmds {
		sendCmd(d,
			d.cmdLogger,
			tCmds[i].Cmd,
			tCmds[i].Errmsg,
			tCmds[i].Hard)
	}
	d.spin.Stop()
	d.statusMsg("Installing Database client complete")

}

// localDBStart
func localDBStart(d *DDConfig, t *targetOS) {
	// Handle the case that the DB is local and doesn't exist
	d.sectionMsg("Starting the database needed for DefectDojo")

	// Create new boostrap command package
	cStartDB := c.NewPkg("startdb")

	// Get commands for the right distro
	switch {
	case t.distro == "ubuntu":
		d.traceMsg("Searching for commands to start MySQL under Ubuntu")
		err := distros.GetUbuntuDB(cStartDB, t.id, d.conf.Install.DB.Engine)
		if err != nil {
			fmt.Printf("Error searching for commands to start database under target OS %s\n", t.id)
			os.Exit(1)
		}
	case t.distro == "rhel":
		d.traceMsg("Searching for commands to start MySQL under RHEL")
		err := distros.GetRHELDB(cStartDB, t.id, d.conf.Install.DB.Engine)
		if err != nil {
			fmt.Printf("Error searching for commands to start database under target OS %s\n", t.id)
			os.Exit(1)
		}
	default:
		d.traceMsg(fmt.Sprintf("Distro identified (%s) is not supported", t.id))
		fmt.Printf("Distro identified by godojo (%s) is not supported, exiting...\n", t.id)
		os.Exit(1)
	}

	// Run the commands to install the chosen DB
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Starting " + d.conf.Install.DB.Engine + " database for DefectDojo..."
	d.spin.Start()
	// Run the start DB command(s) for the target OS
	tCmds, err := distros.CmdsForTarget(cStartDB, t.id)
	if err != nil {
		fmt.Printf("Error getting commands to start DB on target OS %s\n", t.id)
		os.Exit(1)
	}

	for i := range tCmds {
		sendCmd(d,
			d.cmdLogger,
			tCmds[i].Cmd,
			tCmds[i].Errmsg,
			tCmds[i].Hard)
	}
	d.spin.Stop()
	d.statusMsg("Starting Database complete")
}

// prepDBForDojo
func prepDBForDojo(d *DDConfig, t *targetOS) {
	// Preapare the database for DefectDojo by:
	// (1) Checking connectivity to the DB,
	// (2) checking that the configured Dojo database name doesn't exit already
	// (3) Droping the existing database if Drop = true is configured (4) Create the DefectDojo database
	// (5) Add the DB user for DefectDojo to use
	// TODO: Validate this against @owasp - https://docs.google.com/spreadsheets/d/1HuXh3Zr4mrmb6_YmKkDgzl-ZINYZCvVZn31UCqIGpUA/edit#gid=0
	d.sectionMsg("Preparing the database needed for DefectDojo")
	err := dbPrep(d, t)
	if err != nil {
		d.errorMsg(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	// Start the installed DB
	if d.conf.Install.DB.Local {
		d.traceMsg("Starting the local DB")
		localDBStart(d, t)
	}

}

// dbPrep
func dbPrep(d *DDConfig, t *targetOS) error {
	// Call the necessary function for the supported DB engines
	switch d.conf.Install.DB.Engine {
	case "MySQL":
		// Generate commands to install MySQL
		return prepMySQL(d, t.id)
	case "PostgreSQL":
		// Generate commands to install PostgreSQL
		return prepPostgreSQL(d, t)
	}
	// Shouldn't get here but if we do, it's definitely an error
	return errors.New("Unknown database engine configured, cannot check connectivity")
}

func prepMySQL(d *DDConfig, osTar string) error {
	// TODO: Path check any binaries called
	//       * mysql
	// TODO: Check MySQL version and handle MySQL 8 and password format change

	// Set Creds based on dojoConfig.yml
	creds := map[string]string{"user": d.conf.Install.DB.Ruser, "pass": d.conf.Install.DB.Rpass}
	d.traceMsg(fmt.Sprintf("DB Creds from config are %s / %s", creds["user"], creds["pass"]))

	// Creds are unknown if DB is local and newly installed by godojo
	if d.conf.Install.DB.Local && !d.conf.Install.DB.Exists {
		// Determine default access for fresh install of that OS
		// AKA databse is local and didn't exist before the install
		creds = defaultDBCreds(d, osTar)
		d.addRedact(creds["pass"])
	}
	d.traceMsg(fmt.Sprintf("DB Creds are now %s / %s", creds["user"], creds["pass"]))

	d.statusMsg("Validating DB connection")
	// Check connectivity to DB
	conCk := sqlStr{
		os:     osTar,
		sql:    "SHOW PROCESSLIST;",
		errMsg: "Unable to connect to the configured MySQL database",
		creds:  creds,
		kind:   "try",
	}
	_, err := runMySQLCmd(d, conCk)
	if err != nil {
		d.traceMsg("validation of connection to MySQL failed")
		return err
	}

	// Drop existing DefectDojo database if it exists and configuration says to
	if d.conf.Install.DB.Drop {
		d.traceMsg("Dropping any existing database per Install.DB.Drop=True in dojoConfig.yml")
		// Query MySQL to see if the configured database name exists already
		// Another option is "show databases like '" + d.conf.Install.DB.Name + "';"
		dbCk := sqlStr{
			os: osTar,
			sql: "SELECT count(SCHEMA_NAME) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '" +
				d.conf.Install.DB.Name + "';",
			errMsg: "Unable to check for existing DefectDojo MySQL database",
			creds:  creds,
			kind:   "inspect",
		}
		out, err := runMySQLCmd(d, dbCk)
		if err != nil {
			d.traceMsg("Check for existing DefectDojo MySQL database failed")
			fmt.Println("Drop database set to true but no database found, continuing")
			//return err
		}

		// Clean up stdout from inspectCmd output
		strOut := squishSlice(out)
		resp := strings.Trim(
			strings.ReplaceAll(
				strings.ReplaceAll(strOut, "count(SCHEMA_NAME)", ""), "\n", ""), " ")

		// Check if there's an existing DB
		// if resp = 0 then DB doesn't exist
		// if resp = 1 then the DB exists already and needs to be dropped first
		ck, err := strconv.Atoi(resp)
		if err != nil {
			d.traceMsg("Unable to convert existing DB check string to int")
			return err
		}
		if ck == 1 {
			d.traceMsg("DB EXISTS so droping that sucker")
			dropDB := sqlStr{
				os:     osTar,
				sql:    "DROP DATABASE " + d.conf.Install.DB.Name + ";",
				errMsg: "Unable to drop the existing MySQL database",
				creds:  creds,
				kind:   "try",
			}
			_, err := runMySQLCmd(d, dropDB)
			if err != nil {
				d.traceMsg("Failed to drop existing database per configured option to drop existing")
				return err
			}
			fmt.Printf("Existing database %+v dropped since Database Drop was set to %+v\n",
				d.conf.Install.DB.Name, d.conf.Install.DB.Drop)
		}

	}

	// Create the DefectDojo database if it doesn't already exist
	d.traceMsg("Creating database for DefectDojo on MySQL")
	createDB := sqlStr{
		os:     osTar,
		sql:    "CREATE DATABASE IF NOT EXISTS " + d.conf.Install.DB.Name + "  CHARACTER SET UTF8;",
		errMsg: "Unable to create a new MySQL database for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runMySQLCmd(d, createDB)
	if err != nil {
		d.traceMsg("Failed to create new database for DefectDojo to use")
		return err
	}

	// Drop user DefectDojo uses to connect to the database
	d.traceMsg("Dropping existing DefectDojo MySQL DB user, if any")
	dropUsr := sqlStr{
		os: osTar,
		sql: "DROP USER '" + d.conf.Install.DB.User + "'@'localhost';DROP USER '" +
			d.conf.Install.DB.User + "'@'%';",
		errMsg: "Unable to delete existing database user for DefectDojo or one didn't exist",
		creds:  creds,
		kind:   "inspect",
	}
	out, err := runMySQLCmd(d, dropUsr)
	if err != nil {
		// No reason to return the error as this is expected for most cases
		// and create user will error out for edge cases
		d.traceMsg("Unable to delete existing database user for DefectDojo or one didn't exist")
		d.traceMsg(fmt.Sprintf("SQL DROP command output was %+v (in any)", squishSlice(out)))
		d.traceMsg("Continuing after error deleting user, non-fatal error")
	}

	// First set the appropriate host for the DefectDojo user to connect from
	usrHost := "localhost"
	if !d.conf.Install.DB.Local && d.conf.Install.DB.Exists {
		// DB is remote and exists so localhost won't work
		usrHost = "%"
	}
	// Create user for DefectDojo to use to connect to the database
	d.traceMsg("Creating MySQL DB user for DefectDojo")
	createUsr := sqlStr{
		os: osTar,
		sql: "CREATE USER '" + d.conf.Install.DB.User + "'@'" + usrHost +
			"' IDENTIFIED BY '" + d.conf.Install.DB.Pass + "';",
		errMsg: "Unable to create a MySQL database user for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runMySQLCmd(d, createUsr)
	if err != nil {
		d.traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	// Grant the DefectDojo db user the necessary privileges
	d.traceMsg("Granting privileges to DefectDojo MySQL DB user")
	grantPrivs := sqlStr{
		os: osTar,
		sql: "GRANT ALL PRIVILEGES ON " + d.conf.Install.DB.Name + ".* TO '" +
			d.conf.Install.DB.User + "'@'" + d.conf.Install.DB.Host + "';",
		errMsg: "Unable to grant needed privileges to database user for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runMySQLCmd(d, grantPrivs)
	if err != nil {
		d.traceMsg("Failed to grant privileges to database user for DefectDojo")
		return err
	}

	// Flush privileges to finalize changes to db
	d.traceMsg("Flushing privileges for DefectDojo MySQL DB user")
	flushPrivs := sqlStr{
		os:     osTar,
		sql:    "FLUSH PRIVILEGES;",
		errMsg: "Unable to flush database privileges",
		creds:  creds,
		kind:   "try",
	}
	_, err = runMySQLCmd(d, flushPrivs)
	if err != nil {
		d.traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	// Adjust requirements.txt to only use MySQL Python modules
	// TODO: Deprecate this as it's just been a problem generally
	err = trimRequirementsTxt(d, "MySQL")
	if err != nil {
		d.traceMsg("Unable to adjust requirements.txt for MySQL usage, exiting")
		return err
	}

	return nil
}

func runMySQLCmd(d *DDConfig, c sqlStr) ([]string, error) {
	out := make([]string, 1)
	d.traceMsg(fmt.Sprintf("MySQL query: %s", c.sql))
	DBCmds := osCmds{
		id: c.os,
		cmds: []string{"mysql --host=" + d.conf.Install.DB.Host +
			" --user=" + c.creds["user"] +
			" --port=" + strconv.Itoa(d.conf.Install.DB.Port) +
			" --password=" + c.creds["pass"] +
			" --execute=\"" + c.sql + "\""},
		errmsg: []string{c.errMsg},
		hard:   []bool{false},
	}

	// Switch on how to run the command aka runCmds, tryCmds, inspectCmds
	switch c.kind {
	case "try":
		err := tryCmds(d, DBCmds)
		if err != nil {
			d.traceMsg(fmt.Sprintf("Error running MySQL command - %s", c.sql))
			return out, err
		}
	case "inspect":
		out, err := inspectCmds(d, DBCmds)
		if err != nil {
			d.traceMsg(fmt.Sprintf("Error running MySQL command - %s", c.sql))
			return out, err
		}
	default:
		d.traceMsg("Invalid 'kind' sent to runMySQLCmd, bug in godojo")
		fmt.Println("Bug discovered in godojo, see trace message or re-run with trace logging. Quitting.")
		os.Exit(1)
	}

	return out, nil
}

// squishSlice
func squishSlice(sl []string) string {
	str := ""
	for i := 0; i < len(sl); i++ {
		str += sl[i]
	}
	return str
}

func prepPostgreSQL(d *DDConfig, t *targetOS) error {
	// TODO: Path check any binaries called
	//       * postgres

	// Set Creds based on dojoConfig.yml
	creds := map[string]string{"user": d.conf.Install.DB.Ruser, "pass": d.conf.Install.DB.Rpass}
	d.traceMsg(fmt.Sprintf("DB Creds from config are %s / %s", creds["user"], creds["pass"]))

	// Creds are unknown if DB is local and newly installed by godojo
	if d.conf.Install.DB.Local && !d.conf.Install.DB.Exists {
		// Determine default access for fresh install of that OS
		// AKA databse is local and didn't exist before the install
		creds = defaultDBCreds(d, t.id)
		d.addRedact(creds["pass"])
	}
	d.traceMsg(fmt.Sprintf("DB Creds are now %s / %s", creds["user"], creds["pass"]))

	// Update pg_hba.conf for RHEL only (shakes fist at RHEL)
	if !updatePgHba(d, t) {
		d.traceMsg("Failed to update pg_hba.conf, cannot connect to the DB. Quiting install")
		return errors.New("Unable to update pg_hba.conf so SQL to the DB will fail.  Exiting")
	}

	// Use pg_isready to check connectivity to PostgreSQL DB
	d.statusMsg("Checking connectivity to PostgreSQL")

	readyOut, err := isPgReady(d, creds)
	if err != nil {
		d.traceMsg(fmt.Sprintf("PostgreSQL is not available, error was %+v", err))
		return err
	}
	d.traceMsg(fmt.Sprintf("Output from pgReady: %+v", readyOut))

	d.statusMsg("Validating DB connection settings")
	// Check connectivity to DB
	conCk := sqlStr{
		os:     t.id,
		sql:    "\\l",
		errMsg: "Unable to connect to the configured Postgres database",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(d, conCk)
	if err != nil {
		d.traceMsg("validation of connection to Postgres failed")
		// TODO Fix this validation bypass
		//return err
	}

	// Drop existing DefectDojo database if it exists and configuration says to
	if d.conf.Install.DB.Drop {
		d.traceMsg("Dropping any existing database per Install.DB.Drop=True in dojoConfig.yml")
		// Query PostgreSQL to see if the configured database name exists already
		dbCk := sqlStr{
			os:     t.id,
			sql:    "\\l",
			errMsg: "Unable to check for existing DefectDojo PostgreSQL database",
			creds:  creds,
			kind:   "inspect",
		}
		out, err := runPgSQLCmd(d, dbCk)
		if err != nil {
			d.traceMsg("Check for existing DefectDojo PostgreSQL database failed")
			fmt.Println("Drop database set to true but no database found, continuing")
		}

		// Clean up stdout from inspectCmd output
		strOut := squishSlice(out)
		ck := pgParseDBList(d, strOut)

		// Check if there's an existing DB
		// if ck = 0 then DB doesn't exist
		// if ck = 1 then the DB exists already and needs to be dropped first
		if ck == 1 {
			d.traceMsg("DB EXISTS so droping that sucker")
			dropDB := sqlStr{
				os:     t.id,
				sql:    "DROP DATABASE IF EXISTS " + d.conf.Install.DB.Name + ";",
				errMsg: "Unable to drop the existing PostgreSQL database",
				creds:  creds,
				kind:   "try",
			}
			_, err := runPgSQLCmd(d, dropDB)
			if err != nil {
				d.traceMsg("Failed to drop existing database per configured option to drop existing")
				return err
			}
			fmt.Printf("Existing database %+v dropped since Database Drop was set to %+v\n",
				d.conf.Install.DB.Name, d.conf.Install.DB.Drop)
		}

	}

	// Create the DefectDojo database if it doesn't already exist
	d.traceMsg("Creating database for DefectDojo on PostgreSQL")
	createDB := sqlStr{
		os:     t.id,
		sql:    "CREATE DATABASE " + d.conf.Install.DB.Name + ";",
		errMsg: "Unable to create a new PostgreSQL database for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(d, createDB)
	if err != nil {
		d.traceMsg("Failed to create new database for DefectDojo to use")
		// TODO: DEGUGGING
		//return err
	}

	// Drop user DefectDojo uses to connect to the database
	d.traceMsg("Dropping existing DefectDojo PostgreSQL DB user, if any")
	dropUsr := sqlStr{
		os:     t.id,
		sql:    "DROP owned by " + d.conf.Install.DB.User + "; DROP USER IF EXISTS " + d.conf.Install.DB.User + ";",
		errMsg: "Unable to delete existing database user for DefectDojo or one didn't exist",
		creds:  creds,
		kind:   "inspect",
	}
	out, err := runPgSQLCmd(d, dropUsr)
	if err != nil {
		// No reason to return the error as this is expected for most cases
		// and create user will error out for edge cases
		d.traceMsg("Unable to delete existing database user for DefectDojo or one didn't exist")
		d.traceMsg(fmt.Sprintf("SQL DROP command output was %+v (in any)", squishSlice(out)))
		d.traceMsg("Continuing after error deleting user, non-fatal error")
	}

	// Create user for DefectDojo to use to connect to the database
	d.traceMsg("Creating PostgreSQL DB user for DefectDojo")
	createUsr := sqlStr{
		os: t.id,
		sql: "CREATE USER " + d.conf.Install.DB.User + " WITH ENCRYPTED PASSWORD '" +
			d.conf.Install.DB.Pass + "';",
		errMsg: "Unable to create a PostgreSQL database user for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(d, createUsr)
	if err != nil {
		d.traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	// Remote DBs cannot have their pg_hba.conf modified (duh)
	if !d.conf.Install.DB.Local {
		d.statusMsg("Note: pg_hba.conf has not been altered by godojo.")
		d.statusMsg("      It may need to be updated to allow DefectDojo to connect to the DB.")
		d.statusMsg("      Please consult the PostgreSQL documentation for further information.")
	}
	// Grant the DefectDojo db user the necessary privileges
	d.traceMsg("Granting privileges to DefectDojo PostgreSQL DB user")
	grantPrivs := sqlStr{
		os:     t.id,
		sql:    "GRANT ALL PRIVILEGES ON DATABASE " + d.conf.Install.DB.Name + " TO " + d.conf.Install.DB.User + ";",
		errMsg: "Unable to grant needed privileges to database user for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(d, grantPrivs)
	if err != nil {
		d.traceMsg("Failed to grant privileges to database user for DefectDojo")
		return err
	}

	// Set the DefectDojo db user as the owner of dojodb
	d.traceMsg("Seting DefectDojo db user as the owner of the DB")
	setPrivs := sqlStr{
		os:     t.id,
		sql:    "ALTER DATABASE " + d.conf.Install.DB.Name + " OWNER TO " + d.conf.Install.DB.User + ";",
		errMsg: "Unable to set database user as owner of DefectDojo DB",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(d, setPrivs)
	if err != nil {
		d.traceMsg("Failed to set ownership to database user for DefectDojo DB")
		return err
	}

	// Adjust requirements.txt to only use MySQL Python modules
	err = trimRequirementsTxt(d, "PostgreSQL")
	if err != nil {
		d.traceMsg("Unable to adjust requirements.txt for PostgreSQL usage, exiting")
		return err
	}

	return nil
}

func runPgSQLCmd(d *DDConfig, c sqlStr) ([]string, error) {
	out := make([]string, 1)
	d.traceMsg(fmt.Sprintf("Postgres query: %s", c.sql))
	DBCmds := osCmds{
		id: c.os,
		cmds: []string{"sudo -i -u postgres PGPASSWORD=\"" + c.creds["pass"] + "\"" +
			" psql --host=" + d.conf.Install.DB.Host +
			" --username=" + c.creds["user"] +
			" --port=" + strconv.Itoa(d.conf.Install.DB.Port) +
			" --command=\"" + c.sql + "\""},
		errmsg: []string{c.errMsg},
		hard:   []bool{false},
	}

	// Swicht on how to run the command aka runCmds, tryCmds, inspectCmds
	switch c.kind {
	case "try":
		err := tryCmds(d, DBCmds)
		// Handle errors from running the PostgreSQL command
		if err != nil {
			d.traceMsg(fmt.Sprintf("Error running Posgres command - %s", c.sql))
			return out, err
		}
	case "inspect":
		out, err := inspectCmds(d, DBCmds)
		// Handle errors from running the PostgreSQL command
		if err != nil {
			d.traceMsg(fmt.Sprintf("Error running Posgres command - %s", c.sql))
			return out, err
		}
	default:
		d.traceMsg("Invalid 'kind' sent to runPgSQLCmd, bug in godojo")
		fmt.Println("Bug discovered in godojo, see trace message. Quitting.")
		os.Exit(1)
	}

	return out, nil
}

func updatePgHba(d *DDConfig, t *targetOS) bool {
	// Only RHEL and binary compatible distros (e.g. Rocky Linux) need to have pg_hba.conf modified)
	if !strings.Contains(t.distro, "rhel") {
		// return early
		return true
	}

	// For remote DBs, it's not possible to edit pg_hba.conf
	if !d.conf.Install.DB.Local {
		// return early
		return true
	}

	d.traceMsg("RHEL or variant - pg_hba.conf needs to be updated.")
	f, err := os.OpenFile("/var/lib/pgsql/data/pg_hba.conf", os.O_RDWR, 0600)
	if err != nil {
		// Exit with error code if we can't read the default creds file
		d.errorMsg("Unable to read pg_hba.conf file, cannot continue")
		os.Exit(1)
	}
	defer f.Close()

	//// Create a new bufferend reader
	fr := bufio.NewReader(f)

	// Use a scanner to run through the config file to update access
	scanner := bufio.NewScanner(fr)
	content := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "127.0.0.1/32") {
			line = strings.Replace(line, "ident", "md5", 1)
			d.traceMsg("Replaced IPv4 localhost")
		}
		if strings.Contains(line, "::1/128") {
			line = strings.Replace(line, "ident", "md5", 1)
			d.traceMsg("Replaced IPv6 localhost")
		}

		content += line + "\n"
	}

	if err = scanner.Err(); err != nil {
		// Exit with error code if we can't scan the default creds file
		d.errorMsg("Unable to scan the pg_hba.conf file, exiting")
		os.Exit(1)
	}

	// Truncate the file to make sure its empty before writing
	_ = f.Truncate(0)

	// Write new config file by starting at the begining of the file
	_, err = f.WriteAt([]byte(content), 0)
	if err != nil {
		// Exit with error code if we can't scan the default creds file
		d.errorMsg("Unable to write the pg_hba.conf file, exiting")
		os.Exit(1)
	}
	d.traceMsg("Wrote the updated config file")

	// Reload pg_hba.conf using a SQL statement
	d.traceMsg("Re-reading the pg_hba.conf file")
	creds := map[string]string{"user": d.conf.Install.DB.Ruser, "pass": d.conf.Install.DB.Rpass}
	DBCmds := osCmds{
		id: t.id,
		cmds: []string{"sudo -i -u postgres PGPASSWORD=\"" + creds["pass"] + "\"" +
			" psql " + " --username=" + creds["user"] +
			" --port=" + strconv.Itoa(d.conf.Install.DB.Port) +
			" --command=\"SELECT pg_reload_conf();\""},
		errmsg: []string{"Unable to reload the pg_hba.conf file"},
		hard:   []bool{false},
	}
	err = tryCmds(d, DBCmds)
	if err != nil {
		d.traceMsg("Unable to reload the pg_hba.conf file")
		d.errorMsg("Unable to reload the pg_hba.conf file, exiting")
		os.Exit(1)
	}
	d.traceMsg("Restarted PostgreSQL")

	return true
}

// TODO: REPLACE THIS WITH CLIENT CALLS
func isPgReady(d *DDConfig, creds map[string]string) (string, error) {
	d.traceMsg("isPgReady called")

	// Call ps_isready and check exit code
	pgReady := osCmds{
		id: "Linux", cmds: []string{"PGPASSWORD=\"" + creds["pass"] + "\" pg_isready" +
			" --host=" + d.conf.Install.DB.Host +
			" --username=" + creds["user"] +
			" --port=" + strconv.Itoa(d.conf.Install.DB.Port) + " "},
		errmsg: []string{"Unable to run pg_isready to validate PostgreSQL DB status."},
		hard:   []bool{false},
	}

	d.traceMsg(fmt.Sprintf("Running command: %+v", pgReady.cmds))

	out, err := inspectCmds(d, pgReady)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error running pg_isready was: %+v", err))
		// TODO Fix this error bypass
		return squishSlice(out), nil
		//return "", err
	}

	return squishSlice(out), nil
}

// Parse a list of existng PostgreSQL DBs for a specific DB name
// if the DB name is found, return 1 else return 0
func pgParseDBList(d *DDConfig, tbl string) int {
	d.traceMsg(fmt.Sprintf("Parsing DB list for existing DefectDojo DB named  %+v", d.conf.Install.DB.Name))

	// Create a slice for matches
	matches := make([]string, 1)

	// Split up the string by newlines
	lines := strings.Split(tbl, "\n")
	for _, l := range lines {
		trim := strings.TrimLeft(l, " ")
		if len(trim) > 1 {
			// Look at the first character
			// Possibly a lossy match but doubtful
			switch trim[0:1] {
			case "L":
				continue
			case "N":
				continue
			case "-":
				continue
			case "|":
				continue
			case "(":
				continue
			default:
				cells := strings.Split(trim, "|")
				matches = append(matches, strings.TrimSpace(cells[0]))
			}
		}
	}

	// Look for matches and return 1 if a match is found
	for _, v := range matches {
		if strings.Contains(v, d.conf.Install.DB.Name) {
			d.traceMsg("Match found, there is an existing DefectDojo DB")
			return 1
		}
	}

	d.traceMsg("No match found, no existing DefectDojo DB")
	return 0
}

func trimRequirementsTxt(d *DDConfig, dbUsed string) error {
	d.traceMsg("Called trimRequirementsTxt")

	d.traceMsg("WARNING trimRequirementsTxt is deprecated")
	return nil

}

func defaultDBCreds(d *DDConfig, os string) map[string]string {
	// Setup a map to return
	creds := map[string]string{"user": "foo", "pass": "bar"}

	getDefaultDBCreds(d, creds)

	return creds
}

// Determine the default creds for a database freshly installed in Ubuntu
func getDefaultDBCreds(d *DDConfig, creds map[string]string) {
	// Installer currently assumes the default DB passwrod handling won't change by release
	// Switch on the DB type
	switch d.conf.Install.DB.Engine {
	case "MySQL":
		ubuntuDefaultMySQL(d, creds)
		d.warnMsg("MySQL default credentials are not implemented for RHEL Linux")
	case "PostgreSQL":
		// Set creds as the Ruser & Rpass for Postgres
		creds["user"] = d.conf.Install.DB.Ruser
		creds["pass"] = d.conf.Install.DB.Rpass
		setDefaultPgSQL(d, creds)
	}
}

func ubuntuDefaultMySQL(d *DDConfig, c map[string]string) {
	// Sent some initial values that ensure the connection will fail if the file read fails
	c["user"] = "debian-sys-maint"
	c["pass"] = "FAIL"

	// Pull the debian-sys-maint creds from /etc/mysql/debian.cnf
	f, err := os.Open("/etc/mysql/debian.cnf")
	if err != nil {
		// Exit with error code if we can't read the default creds file
		d.errorMsg("Unable to read file with defautl credentials, cannot continue")
		os.Exit(1)
	}

	// Create a new buffered reader
	fr := bufio.NewReader(f)

	// Create a scanner to run through the lines of the file
	scanner := bufio.NewScanner(fr)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "password") {
			l := strings.Split(line, "=")
			// Make sure there was a "=" in l
			if len(l) > 1 {
				c["pass"] = strings.Trim(l[1], " ")
				break
			}
		}
	}
	if err = scanner.Err(); err != nil {
		// Exit with error code if we can't scan the default creds file
		d.errorMsg("Unable to scan file with defautl credentials, cannot continue")
		os.Exit(1)
	}

}

func setDefaultPgSQL(d *DDConfig, creds map[string]string) {
	d.traceMsg("Called setDefaultPgSQL")

	// Set user to postgres as that's the default DB user for any new install
	creds["user"] = "postgres"

	// Use the default local OS user to set the postgres DB user
	pgAlter := osCmds{
		id:     "linux",
		cmds:   []string{"sudo -i -u postgres psql -c \"ALTER USER postgres with encrypted password '" + creds["pass"] + "';\""},
		errmsg: []string{"Unable to set initial password for PostgreSQL database user postgres"},
		hard:   []bool{false},
	}

	// Try command
	err := tryCmds(d, pgAlter)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error updating PostgreSQL DB user with %+v", squishSlice(pgAlter.cmds)))
		d.errorMsg("Unable to update default PostgreSQL DB user, quitting")
		os.Exit(1)
	}

	d.traceMsg("No error return from setDefaultPgSQL")
}
