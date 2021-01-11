package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mtesauro/godojo/config"
)

// saneDBConfig checks if the options configured in dojoConfig.yml are
// possible aka sane and will exist the installer with a message if they are not
func saneDBConfig(local bool, exists bool) {
	// Remote database that doesn't exist - godojo can't help you here
	if !local && !exists {
		errorMsg("Remote database which doens't exist was confgiured in dojoConfig.yml.")
		errorMsg("This is an unsupported configuration.")
		statusMsg("Correct configuration and/or install a remote DB before running installer again.")
		fmt.Printf("Exiting...\n\n")
		os.Exit(1)
	}
}

func installDB(osTar string, dbTar *config.DBTarget, dCmd *osCmds) {
	// Look at the dbTar and call function to install that DB target
	switch dbTar.Engine {
	case "SQLite":
		// Generate commands to install SQLite
		instSQLite(osTar, dCmd)
	case "MariaDB":
		// Generate commands to install MariaDB
		instMariaDB(osTar, dCmd)
	case "MySQL":
		// Generate commands to install MySQL
		instMySQL(osTar, dCmd)
	case "PostgreSQL":
		// Generate commands to install PostgreSQL
		instPostgreSQL(osTar, dCmd)
	}
	return
}

func startDB(osTar string, dbTar *config.DBTarget, dbCmd *osCmds) {
	// Look at the dbTar and call function to install that DB target
	switch dbTar.Engine {
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
		dbCmd.id = osTar
		// TODO: Propably time to convert this to systemctl calls
		//       also consider enabling the service just in case
		dbCmd.cmds = []string{
			"service postgresql start",
		}
		dbCmd.errmsg = []string{
			"Unable to start PostgreSQL",
		}
		dbCmd.hard = []bool{
			true,
		}
	}
}

func dbPrep(osTar string, dbTar *config.DBTarget) error {
	// Call the necessary function for the supported DB engines
	switch dbTar.Engine {
	case "SQLite":
		// Generate commands to install SQLite
		return prepSQLite(dbTar, osTar)
	case "MariaDB":
		// Generate commands to install MariaDB
		return prepMariaDB(dbTar, osTar)
	case "MySQL":
		// Generate commands to install MySQL
		return prepMySQL(dbTar, osTar)
	case "PostgreSQL":
		// Generate commands to install PostgreSQL
		return prepPostgreSQL(dbTar, osTar)
	}
	// Shouldn't get here but if we do, it's definitely an error
	return errors.New("Unknown database engine configured, cannot check connectivity")
}

func prepSQLite(dbTar *config.DBTarget, os string) error {
	// Open a connection the the configured SQLite DB
	// https://github.com/mattn/go-sqlite3#dsn-examples
	// TODO - write this code and test it
	//return nil
	return errors.New("Not implemented yet")
}

func prepMariaDB(dbTar *config.DBTarget, os string) error {
	// TODO - test that this works MariaDB sd
	// Open a connection the the configured MySQL DB
	// https://github.com/go-sql-driver/mysql/#dsn-data-source-name
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//conn := dbTar.User + ":" + dbTar.Pass + "@" + dbTar.Host + ":" + strconv.Itoa(dbTar.Port)
	//fmt.Println("DB conn is ", conn)
	//dbMySQL, err := sql.Open("mysql", conn)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(dbMySQL)

	//return nil
	return errors.New("Not implemented yet")
}

// Setup a struct to use for DB commands
type sqlStr struct {
	os     string
	sql    string
	errMsg string
	creds  map[string]string
	kind   string
}

func prepMySQL(dbTar *config.DBTarget, osTar string) error {
	// TODO: Path check any binaries called
	//       * mysqladmin
	// TODO: Check MySQL version and handle MySQL 8 and password format change

	// Set Creds based on dojoConfig.yml
	creds := map[string]string{"user": dbTar.Ruser, "pass": dbTar.Rpass}
	traceMsg(fmt.Sprintf("DB Creds from config are %s / %s", creds["user"], creds["pass"]))

	// Creds are unknown if DB is local and newly installed by godojo
	if dbTar.Local && !dbTar.Exists {
		// Determine default access for fresh install of that OS
		// AKA databse is local and didn't exist before the install
		creds = defaultDBCreds(dbTar.Engine, osTar)
		addRedact(creds["pass"])
	}
	traceMsg(fmt.Sprintf("DB Creds are now %s / %s", creds["user"], creds["pass"]))

	statusMsg("Validating DB connection")
	// Check connectivity to DB
	conCk := sqlStr{
		os:     osTar,
		sql:    "SHOW PROCESSLIST;",
		errMsg: "Unable to connect to the configured MySQL database",
		creds:  creds,
		kind:   "try",
	}
	_, err := runMySQLCmd(dbTar, conCk)
	if err != nil {
		traceMsg("validation of connection to MySQL failed")
		return err
	}

	// Drop existing DefectDojo database if it exists and configuration says to
	if dbTar.Drop {
		traceMsg("Dropping any existing database per Install.DB.Drop=True in dojoConfig.yml")
		// TODO: Convert this and the above call to a function
		// Query MySQL to see if the configured database name exists already
		// Another option is "show databases like '" + dbTar.Name + "';"
		dbCk := sqlStr{
			os:     osTar,
			sql:    "SELECT count(SCHEMA_NAME) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '" + dbTar.Name + "';",
			errMsg: "Unable to check for existing DefectDojo MySQL database",
			creds:  creds,
			kind:   "inspect",
		}
		out, err := runMySQLCmd(dbTar, dbCk)
		if err != nil {
			traceMsg("Check for existing DefectDojo MySQL database failed")
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
			traceMsg("Unable to convert existing DB check string to int")
			return err
		}
		if ck == 1 {
			traceMsg("DB EXISTS so droping that sucker")
			dropDB := sqlStr{
				os:     osTar,
				sql:    "DROP DATABASE " + dbTar.Name + ";",
				errMsg: "Unable to drop the existing MySQL database",
				creds:  creds,
				kind:   "try",
			}
			_, err := runMySQLCmd(dbTar, dropDB)
			if err != nil {
				traceMsg("Failed to drop existing database per configured option to drop existing")
				return err
			}
			fmt.Printf("Existing database %+v dropped since Database Drop was set to %+v\n", dbTar.Name, dbTar.Drop)
		}

	}

	// Create the DefectDojo database if it doesn't already exist
	traceMsg("Creating database for DefectDojo on MySQL")
	createDB := sqlStr{
		os:     osTar,
		sql:    "CREATE DATABASE IF NOT EXISTS " + dbTar.Name + "  CHARACTER SET UTF8;",
		errMsg: "Unable to create a new MySQL database for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runMySQLCmd(dbTar, createDB)
	if err != nil {
		traceMsg("Failed to create new database for DefectDojo to use")
		return err
	}

	// Drop user DefectDojo uses to connect to the database
	traceMsg("Dropping existing DefectDojo MySQL DB user, if any")
	dropUsr := sqlStr{
		os:     osTar,
		sql:    "DROP USER '" + dbTar.User + "'@'localhost';DROP USER '" + dbTar.User + "'@'%';",
		errMsg: "Unable to delete existing database user for DefectDojo or one didn't exist",
		creds:  creds,
		kind:   "inspect",
	}
	out, err := runMySQLCmd(dbTar, dropUsr)
	if err != nil {
		// No reason to return the error as this is expected for most cases
		// and create user will error out for edge cases
		traceMsg("Unable to delete existing database user for DefectDojo or one didn't exist")
		traceMsg(fmt.Sprintf("SQL DROP command output was %+v (in any)", squishSlice(out)))
		traceMsg("Continuing after error deleting user, non-fatal error")
	}

	// First set the appropriate host for the DefectDojo user to connect from
	usrHost := "localhost"
	if !dbTar.Local && dbTar.Exists {
		// DB is remote and exists so localhost won't work
		usrHost = "%"
	}
	// Create user for DefectDojo to use to connect to the database
	traceMsg("Creating MySQL DB user for DefectDojo")
	createUsr := sqlStr{
		os:     osTar,
		sql:    "CREATE USER '" + dbTar.User + "'@'" + usrHost + "' IDENTIFIED BY '" + dbTar.Pass + "';",
		errMsg: "Unable to create a MySQL database user for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runMySQLCmd(dbTar, createUsr)
	if err != nil {
		traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	// Grant the DefectDojo db user the necessary privileges
	traceMsg("Granting privileges to DefectDojo MySQL DB user")
	grantPrivs := sqlStr{
		os:     osTar,
		sql:    "GRANT ALL PRIVILEGES ON " + dbTar.Name + ".* TO '" + dbTar.User + "'@'" + dbTar.Host + "';",
		errMsg: "Unable to grant needed privileges to database user for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runMySQLCmd(dbTar, grantPrivs)
	if err != nil {
		traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	// Flush privileges to finalize changes to db
	traceMsg("Flushing privileges for DefectDojo MySQL DB user")
	flushPrivs := sqlStr{
		os:     osTar,
		sql:    "FLUSH PRIVILEGES;",
		errMsg: "Unable to flush database privileges",
		creds:  creds,
		kind:   "try",
	}
	_, err = runMySQLCmd(dbTar, flushPrivs)
	if err != nil {
		traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	return nil
}

func runMySQLCmd(dbTar *config.DBTarget, c sqlStr) ([]string, error) {
	out := make([]string, 1)
	traceMsg(fmt.Sprintf("MySQL query: %s", c.sql))
	DBCmds := osCmds{
		id: c.os,
		cmds: []string{"mysql --host=" + dbTar.Host +
			" --user=" + c.creds["user"] +
			" --port=" + strconv.Itoa(dbTar.Port) +
			" --password=" + c.creds["pass"] +
			" --execute=\"" + c.sql + "\""},
		errmsg: []string{c.errMsg},
		hard:   []bool{false},
	}

	// Swicht on how to run the command aka runCmds, tryCmds, inspectCmds
	err := errors.New("")
	switch c.kind {
	case "try":
		err = tryCmds(cmdLogger, DBCmds)
	case "inspect":
		out, err = inspectCmds(cmdLogger, DBCmds)
	default:
		traceMsg("Invalid 'kind' sent to runMySQLCmd, bug in godojo")
		fmt.Println("Bug discovered in godojo, see trace message. Quitting.")
		os.Exit(1)
	}

	// Handle errors from running the MySQL command
	if err != nil {
		traceMsg(fmt.Sprintf("Error running MySQL command - %s", c.sql))
		return out, err
	}

	return out, nil
}

func squishSlice(sl []string) string {
	str := ""
	for i := 0; i < len(sl); i++ {
		str += sl[i]
	}
	return str
}

func prepPostgreSQL(dbTar *config.DBTarget, os string) error {
	// Open a connection to the configured PostgreSQL database
	// https://godoc.org/github.com/lib/pq
	//conn := "user=" + dbTar.User + " password=" + dbTar.Pass + " host=" + dbTar.Host + " port=" + strconv.Itoa(dbTar.Port)
	//fmt.Println("DB conn is ", conn)

	//dbPostgreSQL, err := sql.Open("postgres", conn)
	//if err != nil {
	//	return err
	//}

	//fmt.Println(dbPostgreSQL)
	return nil
}
