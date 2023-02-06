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
func saneDBConfig(d *gdjDefault) {
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
func installDBForDojo(d *gdjDefault, o *targetOS) {
	// Handle the case that the DB is local and doesn't exist
	if !d.conf.Install.DB.Exists {
		// Note that godojo won't try to install remote databases
		dbNotExist(d, o)
	}

	// Install DB clients for remote DBs
	if !d.conf.Install.DB.Local {
		clientInst := osCmds{}
		dbClient(d, o.id, &clientInst)
	}

	// Start the database if local and didn't already exist
	if d.conf.Install.DB.Local && !d.conf.Install.DB.Exists {
		dbStart := osCmds{}
		localDBStart(d, o.id, &dbStart)
	}

}

// dbNotExist takes a pointer to a gdjDefault struct and a pointer to targetOS
// struct and runs the commands necesary to install a local database of the
// supported type (PostgreSQL, MySQL, etc)
func dbNotExist(d *gdjDefault, o *targetOS) {
	// Handle the case that the DB is local and doesn't exist
	d.sectionMsg("Installing database needed for DefectDojo")

	// Gather OS commands to install the DB
	dbInst := osCmds{}
	installDB(d, o.id, &dbInst)

	// Run the commands to install the chosen DB
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Installing " + d.conf.Install.DB.Engine + " database for DefectDojo..."
	d.spin.Start()
	for i := range dbInst.cmds {
		sendCmd(d,
			d.cmdLogger,
			dbInst.cmds[i],
			dbInst.errmsg[i],
			dbInst.hard[i])
	}
	d.spin.Stop()
	d.statusMsg("Installing Database complete")
}

// installDB
func installDB(d *gdjDefault, osTar string, dCmd *osCmds) {
	// Look at the dbTar and call function to install that DB target
	switch d.conf.Install.DB.Engine {
	case "SQLite":
		// Generate commands to install SQLite
		instSQLite(d, osTar, dCmd)
	case "MariaDB":
		// Generate commands to install MariaDB
		instMariaDB(d, osTar, dCmd)
	case "MySQL":
		// Generate commands to install MySQL
		instMySQL(d, osTar, dCmd)
	case "PostgreSQL":
		// Generate commands to install PostgreSQL
		instPostgreSQL(d, osTar, dCmd)
	}
	return
}

// instSQLite
func instSQLite(d *gdjDefault, id string, b *osCmds) {
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
		ubuntuInstSQLite(d, id, b)
	}
	return
}

// Commands to install SQLite on Ubuntu
func ubuntuInstSQLite(d *gdjDefault, id string, b *osCmds) {
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
			"DEBIAN_FRONTEND=noninteractive apt-get install -y sqlite3",
		}
		b.errmsg = []string{
			"Unable to install SQLite",
		}
		b.hard = []bool{
			true,
		}
	}
	d.warnMsg("sqlite is a deprecated database for DefectDojo")
	return
}

// instMariaDB
func instMariaDB(d *gdjDefault, id string, b *osCmds) {
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
		ubuntuInstMariaDB(d, id, b)
	}
	return
}

// Commands to install MariaDB on Ubuntu
func ubuntuInstMariaDB(d *gdjDefault, id string, b *osCmds) {
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
			"DEBIAN_FRONTEND=noninteractive apt-get install -y mariadb-server libmariadbclient-dev",
		}
		b.errmsg = []string{
			"Unable to install MariaDB",
		}
		b.hard = []bool{
			true,
		}
	}
	d.warnMsg("MariaDB is an unsupported database for DefectDojo")
	return
}

// instMySQL
func instMySQL(d *gdjDefault, id string, b *osCmds) {
	d.traceMsg(fmt.Sprintf("Installing MySQL for %s\n", id))
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
		ubuntuInstMySQL(d, id, b)
	}
	return
}

// Commands to install MySQL on Ubuntu
func ubuntuInstMySQL(d *gdjDefault, id string, b *osCmds) {
	d.traceMsg(fmt.Sprintf("Installing Ubuntu MySQL for %s\n", id))
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
			"DEBIAN_FRONTEND=noninteractive apt-get install -y mysql-server libmysqlclient-dev",
		}
		b.errmsg = []string{
			"Unable to install MySQL",
		}
		b.hard = []bool{
			true,
		}
	}
	d.warnMsg("Posgres is the preferred database for DefectDojo")
	return
}

// instPostgreSQL
func instPostgreSQL(d *gdjDefault, id string, b *osCmds) {
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
		ubuntuInstPostgreSQL(d, id, b)
	}
	return
}

// Commands to install PostgreSQL on Ubuntu
func ubuntuInstPostgreSQL(d *gdjDefault, id string, b *osCmds) {
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
			"DEBIAN_FRONTEND=noninteractive apt-get install -y libpq-dev postgresql postgresql-contrib postgresql-client-common",
		}
		b.errmsg = []string{
			"Unable to install PostgreSQL",
		}
		b.hard = []bool{
			true,
		}
	}
	return
}

// dbClientInstall
func dbClient(d *gdjDefault, osTar string, dCmd *osCmds) {
	// Setup commands for DB clients
	installDBClient(d, osTar, dCmd)

	// Run the commands to install the chosen DB
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Installing " + d.conf.Install.DB.Engine + " database client for DefectDojo..."
	d.spin.Start()
	for i := range dCmd.cmds {
		sendCmd(d,
			d.cmdLogger,
			dCmd.cmds[i],
			dCmd.errmsg[i],
			dCmd.hard[i])
	}
	d.spin.Stop()
	d.statusMsg("Installing Database client complete")

}

// installDBClient
func installDBClient(d *gdjDefault, osTar string, dCmd *osCmds) {
	// Look at the dbTar and call function to install that DB target
	switch d.conf.Install.DB.Engine {
	case "SQLite":
		// Generate commands to install SQLite
		// A remote SQLite DB makes no sense
		d.warnMsg("A remote sqlite database makes no sense.")
		return
	case "MariaDB":
		// Generate commands to install MariaDB
		// TODO: Write install for MariaDB client
		//instMariaDBClient(osTar, dCmd)
		d.warnMsg("MariaDB is not a supported DefectDojo databse")
		return
	case "MySQL":
		// Generate commands to install MySQL
		// TODO: Write install for MySQL client
		//instMySQLClient(osTar, dCmd)
		d.warnMsg("MySQL client install has not been implemented")
		return
	case "PostgreSQL":
		// Generate commands to install PostgreSQL
		instPostgreSQLClient(osTar, dCmd)
	}
	return
}

// instPostgreSQLClient
func instPostgreSQLClient(id string, b *osCmds) {
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
		ubuntuInstPostgreSQLClient(id, b)
	}
	return
}

// ubuntuInstPostgreSQLClient
func ubuntuInstPostgreSQLClient(id string, b *osCmds) {
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
			"DEBIAN_FRONTEND=noninteractive apt-get install -y postgresql-client-12",
			"/usr/sbin/groupadd -f postgres",                         // TODO: consider using os.Group.Lookup before calling this
			"/usr/sbin/useradd -s /bin/bash -m -g postgres postgres", // TODO: consider using os.User.Lookup before calling this
		}
		b.errmsg = []string{
			"Unable to install PostgreSQL client",
			"Unable to add postgres group",
			"Unable to add postgres user",
		}
		b.hard = []bool{
			true,
			true,
			false, // incase there is an existing postgres user, useradd returns a 9 exit code
		}
	}
	return
}

// localDBStart
func localDBStart(d *gdjDefault, id string, c *osCmds) {
	// Handle the case that the DB is local and doesn't exist
	d.sectionMsg("Starting the database needed for DefectDojo")

	// Gather OS commands to install the DB
	startDB(d, id, c)

	// Run the commands to install the chosen DB
	d.spin = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	d.spin.Prefix = "Starting " + d.conf.Install.DB.Engine + " database for DefectDojo..."
	d.spin.Start()
	for i := range c.cmds {
		sendCmd(d,
			d.cmdLogger,
			c.cmds[i],
			c.errmsg[i],
			c.hard[i])
	}
	d.spin.Stop()
	d.statusMsg("Starting Database complete")
}

// startDB
func startDB(d *gdjDefault, osTar string, dbCmd *osCmds) {
	// Look at the dbTar and call function to install that DB target
	switch d.conf.Install.DB.Engine {
	case "SQLite":
		startSQLite(osTar, dbCmd)
	case "MariaDB":
		startMariaDB(osTar, dbCmd)
	case "MySQL":
		startMySQL(osTar, dbCmd)
	case "PostgreSQL":
		startPostgres(osTar, dbCmd)
	}
	return
}

func startSQLite(osTar string, dbCmd *osCmds) {
	// Generate commands to start SQLite
	switch osTar {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		dbCmd.id = osTar
		dbCmd.cmds = []string{
			"echo 'Nothing to start for SQLite'",
		}
		dbCmd.errmsg = []string{
			"Starting SQLite should never error since there's nothing to start",
		}
		dbCmd.hard = []bool{
			true,
		}
	}
}

func startMariaDB(osTar string, dbCmd *osCmds) {
	// Generate commands to start MariaDB
	switch osTar {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		dbCmd.id = osTar
		// TODO: Propably time to convert this to systemctl calls
		//       also consider enabling the service just in case
		dbCmd.cmds = []string{
			"service mysql start",
		}
		dbCmd.errmsg = []string{
			"Unable to start MariaDB",
		}
		dbCmd.hard = []bool{
			true,
		}
	}
}

func startMySQL(osTar string, dbCmd *osCmds) {
	// Generate commands to start MySQL
	switch osTar {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		dbCmd.id = osTar
		// TODO: Propably time to convert this to systemctl calls
		//       also consider enabling the service just in case
		dbCmd.cmds = []string{
			"service mysql start",
		}
		dbCmd.errmsg = []string{
			"Unable to start MySQL",
		}
		dbCmd.hard = []bool{
			true,
		}
	}
}

func startPostgres(osTar string, dbCmd *osCmds) {
	// Generate commands to start PostgreSQL
	switch osTar {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		dbCmd.id = osTar
		// TODO: Propably time to convert this to systemctl calls
		//       also consider enabling the service just in case
		dbCmd.cmds = []string{
			"/usr/sbin/service postgresql start",
		}
		dbCmd.errmsg = []string{
			"Unable to start PostgreSQL",
		}
		dbCmd.hard = []bool{
			true,
		}
	}
}

// prepDBForDojo
func prepDBForDojo(d *gdjDefault, o *targetOS) {
	// Preapare the database for DefectDojo by:
	// (1) Checking connectivity to the DB,
	// (2) checking that the configured Dojo database name doesn't exit already
	// (3) Droping the existing database if Drop = true is configured (4) Create the DefectDojo database
	// (5) Add the DB user for DefectDojo to use
	// TODO: Validate this against @owasp - https://docs.google.com/spreadsheets/d/1HuXh3Zr4mrmb6_YmKkDgzl-ZINYZCvVZn31UCqIGpUA/edit#gid=0
	d.sectionMsg("Preparing the database needed for DefectDojo")
	err := dbPrep(d, o.id)
	if err != nil {
		d.errorMsg(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

}

// dbPrep
func dbPrep(d *gdjDefault, osTar string) error {
	// Call the necessary function for the supported DB engines
	switch d.conf.Install.DB.Engine {
	case "SQLite":
		// Generate commands to install SQLite
		return prepSQLite()
	case "MariaDB":
		// Generate commands to install MariaDB
		return prepMariaDB()
	case "MySQL":
		// Generate commands to install MySQL
		return prepMySQL(d, osTar)
	case "PostgreSQL":
		// Generate commands to install PostgreSQL
		return prepPostgreSQL(d, osTar)
	}
	// Shouldn't get here but if we do, it's definitely an error
	return errors.New("Unknown database engine configured, cannot check connectivity")
}

func prepSQLite() error {
	// Open a connection the the configured SQLite DB
	// https://github.com/mattn/go-sqlite3#dsn-examples
	// TODO - write this code and test it
	return errors.New("Not implemented yet")
}

func prepMariaDB() error {
	// TODO - Decide if this should even be supported
	return errors.New("Not implemented yet")
}

func prepMySQL(d *gdjDefault, osTar string) error {
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

func runMySQLCmd(d *gdjDefault, c sqlStr) ([]string, error) {
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

	// Swicht on how to run the command aka runCmds, tryCmds, inspectCmds
	err := errors.New("")
	switch c.kind {
	case "try":
		err = tryCmds(d, DBCmds)
	case "inspect":
		out, err = inspectCmds(d, DBCmds)
	default:
		d.traceMsg("Invalid 'kind' sent to runMySQLCmd, bug in godojo")
		fmt.Println("Bug discovered in godojo, see trace message. Quitting.")
		os.Exit(1)
	}

	// Handle errors from running the MySQL command
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error running MySQL command - %s", c.sql))
		return out, err
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

func prepPostgreSQL(d *gdjDefault, osTar string) error {
	// TODO: Path check any binaries called
	//       * postgres

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

	// Use pg_isready to check connectivity to PostgreSQL DB
	d.statusMsg("Checking connectivity to PostgreSQL")

	t, err := isPgReady(d, creds)
	if err != nil {
		d.traceMsg(fmt.Sprintf("PostgreSQL is not available, error was %+v", err))
		return err
	}
	d.traceMsg(fmt.Sprintf("Output from pgReady: %+v", t))

	d.statusMsg("Validating DB connection settings")
	// Check connectivity to DB
	conCk := sqlStr{
		os:     osTar,
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
			os:     osTar,
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
				os:     osTar,
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
		os:     osTar,
		sql:    "CREATE DATABASE " + d.conf.Install.DB.Name + ";",
		errMsg: "Unable to create a new PostgreSQL database for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(d, createDB)
	if err != nil {
		d.traceMsg("Failed to create new database for DefectDojo to use")
		return err
	}

	// Drop user DefectDojo uses to connect to the database
	d.traceMsg("Dropping existing DefectDojo PostgreSQL DB user, if any")
	dropUsr := sqlStr{
		os:     osTar,
		sql:    "DROP USER IF EXISTS " + d.conf.Install.DB.User + ";",
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
		os: osTar,
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

	d.statusMsg("Note: pg_hba.conf has not been altered by godojo.")
	d.statusMsg("      It may need to be updated to allow DefectDojo to connect to the DB.")
	d.statusMsg("      Please consult the PostgreSQL documentation for further information.")

	// Grant the DefectDojo db user the necessary privileges
	d.traceMsg("Granting privileges to DefectDojo PostgreSQL DB user")
	grantPrivs := sqlStr{
		os:     osTar,
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

	// Adjust requirements.txt to only use MySQL Python modules
	err = trimRequirementsTxt(d, "PostgreSQL")
	if err != nil {
		d.traceMsg("Unable to adjust requirements.txt for PostgreSQL usage, exiting")
		return err
	}

	return nil
}

func runPgSQLCmd(d *gdjDefault, c sqlStr) ([]string, error) {
	out := make([]string, 1)
	d.traceMsg(fmt.Sprintf("Postgres query: %s", c.sql))
	DBCmds := osCmds{
		id: c.os,
		cmds: []string{"sudo -u postgres PGPASSWORD=\"" + c.creds["pass"] + "\"" +
			" psql --host=" + d.conf.Install.DB.Host +
			" --username=" + c.creds["user"] +
			" --port=" + strconv.Itoa(d.conf.Install.DB.Port) +
			" --command=\"" + c.sql + "\""},
		errmsg: []string{c.errMsg},
		hard:   []bool{false},
	}

	// Swicht on how to run the command aka runCmds, tryCmds, inspectCmds
	err := errors.New("")
	switch c.kind {
	case "try":
		err = tryCmds(d, DBCmds)
	case "inspect":
		out, err = inspectCmds(d, DBCmds)
	default:
		d.traceMsg("Invalid 'kind' sent to runPgSQLCmd, bug in godojo")
		fmt.Println("Bug discovered in godojo, see trace message. Quitting.")
		os.Exit(1)
	}

	// Handle errors from running the PostgreSQL command
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error running Posgres command - %s", c.sql))
		return out, err
	}

	return out, nil
}

func isPgReady(d *gdjDefault, creds map[string]string) (string, error) {
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
func pgParseDBList(d *gdjDefault, tbl string) int {
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

func trimRequirementsTxt(d *gdjDefault, dbUsed string) error {
	d.traceMsg("Called trimRequirementsTxt")

	req := make([]string, 1)
	switch dbUsed {
	case "MySQL":
		d.traceMsg("MySQL requirements fix")
		req[0] = "psycopg2-binary"
		// append any additional requirements
	case "PostgreSQL":
		d.traceMsg("PostgreSQL requirements fix")
		req[0] = "mysqlclient"
		// append any additional requirements
	default:
		return errors.New("Unknown database provided to trimRequirementsTxt()")
	}

	d.traceMsg("No error return from trimRequirementsTxt")
	return removeRequirement(d, req)
}

func removeRequirement(d *gdjDefault, r []string) error {
	d.traceMsg("Called removeRequirement")

	// Set th path for requirements.txt
	p := d.conf.Install.Root + "/" + d.conf.Install.Source + "/"

	cyaCopy := osCmds{
		id:     "Linux",
		cmds:   []string{"cp " + p + "requirements.txt " + p + "requirements.cya"},
		errmsg: []string{"Unable to make a cya copy of requirements.txt"},
		hard:   []bool{false},
	}

	err := tryCmds(d, cyaCopy)
	if err != nil {
		d.traceMsg("Failed creating a backup copy of requirements.txt")
		return err
	}

	// Open requirements.txt
	f, err := os.OpenFile(p+"requirements.txt", os.O_RDWR, 0744)
	if err != nil {
		d.traceMsg("Unable to open requirements.txt")
		return err
	}
	defer f.Close()

	// Scan the file for a match
	newf := ""
	skip := ""
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		// Run through Python modules to remove from requirements.txt
		for _, v := range r {
			if !strings.Contains(line, v) {
				d.traceMsg(fmt.Sprintf("Keeping requirements line: %+v", line))
				newf += line + "\n"
			} else {
				d.traceMsg(fmt.Sprintf("Skiping requirements line: %+v", line))
				skip += "#\n"
			}
		}

	}

	// Clearn the file contents
	err = f.Truncate(0)
	if err != nil {
		d.traceMsg("Unable to truncate requirements.txt file")
		return err
	}

	// Add a couple of extra lines just in case
	skip += "# File re-written by godojo\n"
	newf += skip

	// Write resulting file
	_, err = f.WriteAt([]byte(newf), 0)
	if err != nil {
		d.traceMsg("Unable to write new requirements.txt file")
		return err
	}

	d.traceMsg("No error return from removeRequirement")
	return nil
}

func defaultDBCreds(d *gdjDefault, os string) map[string]string {
	// Setup a map to return
	creds := map[string]string{"user": "foo", "pass": "bar"}

	// Get the default creds based on OS
	switch os {
	case "ubuntu:18.04":
		fallthrough
	case "ubuntu:20.04":
		fallthrough
	case "ubuntu:20.10":
		fallthrough
	case "ubuntu:21.04":
		fallthrough
	case "ubuntu:22.04":
		ubuntuDefaultDBCreds(d, creds)
	}

	return creds
}

// Determine the default creds for a database freshly installed in Ubuntu
func ubuntuDefaultDBCreds(d *gdjDefault, creds map[string]string) {
	// Installer currently assumes the default DB passwrod handling won't change by release
	// Switch on the DB type
	switch d.conf.Install.DB.Engine {
	case "MySQL":
		ubuntuDefaultMySQL(d, creds)
	case "PostgreSQL":
		// Set creds as the Ruser & Rpass for Postgres
		creds["user"] = d.conf.Install.DB.Ruser
		creds["pass"] = d.conf.Install.DB.Rpass
		ubuntuDefaultPgSQL(d, creds)
	}

	return
}

func ubuntuDefaultMySQL(d *gdjDefault, c map[string]string) {
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

func ubuntuDefaultPgSQL(d *gdjDefault, creds map[string]string) {
	d.traceMsg("Called ubuntuDefaultPgSQL")

	// Set user to postgres as that's the default DB user for any new install
	creds["user"] = "postgres"

	// Use the default local OS user to set the postgres DB user
	pgAlter := osCmds{
		id:     "linux",
		cmds:   []string{"sudo -u postgres psql -c \"ALTER USER postgres with encrypted password '" + creds["pass"] + "';\""},
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

	d.traceMsg("No error return from ubuntuDefaultPgSQL")
	return
}
