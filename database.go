package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
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

func startDB(osTar string, dbTar *config.DBTarget, dCmd *osCmds) {
	// Look at the dbTar and call function to install that DB target
	switch dbTar.Engine {
	case "SQLite":
		// Generate commands to start SQLite
		switch osTar {
		case "ubuntu:18.04":
			fallthrough
		case "ubuntu:20.04":
			fallthrough
		case "ubuntu:20.10":
			dCmd.id = osTar
			dCmd.cmds = []string{
				"echo 'Nothing to start for SQLite'",
			}
			dCmd.errmsg = []string{
				"Starting SQLite should never error since there's nothing to start",
			}
			dCmd.hard = []bool{
				true,
			}
		}
	case "MariaDB":
		// Generate commands to start MariaDB
		switch osTar {
		case "ubuntu:18.04":
			fallthrough
		case "ubuntu:20.04":
			fallthrough
		case "ubuntu:20.10":
			dCmd.id = osTar
			// TODO: Propably time to convert this to systemctl calls
			//       also consider enabling the service just in case
			dCmd.cmds = []string{
				"service mysql start",
			}
			dCmd.errmsg = []string{
				"Unable to start MariaDB",
			}
			dCmd.hard = []bool{
				true,
			}
		}
	case "MySQL":
		// Generate commands to start MySQL
		switch osTar {
		case "ubuntu:18.04":
			fallthrough
		case "ubuntu:20.04":
			fallthrough
		case "ubuntu:20.10":
			dCmd.id = osTar
			// TODO: Propably time to convert this to systemctl calls
			//       also consider enabling the service just in case
			dCmd.cmds = []string{
				"service mysql start",
			}
			dCmd.errmsg = []string{
				"Unable to start MySQL",
			}
			dCmd.hard = []bool{
				true,
			}
		}
	case "PostgreSQL":
		// Generate commands to start PostgreSQL
		switch osTar {
		case "ubuntu:18.04":
			fallthrough
		case "ubuntu:20.04":
			fallthrough
		case "ubuntu:20.10":
			dCmd.id = osTar
			// TODO: Propably time to convert this to systemctl calls
			//       also consider enabling the service just in case
			dCmd.cmds = []string{
				"service postgresql start",
			}
			dCmd.errmsg = []string{
				"Unable to start PostgreSQL",
			}
			dCmd.hard = []bool{
				true,
			}
		}
	}
	return
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
	// TODO: Convert this and the below calls to a function
	// Check connectivity to DB
	DBCmds := osCmds{
		id: osTar,
		cmds: []string{"mysqladmin --host=" + dbTar.Host +
			" --user=" + creds["user"] +
			" --port=" + strconv.Itoa(dbTar.Port) +
			" --password=" + creds["pass"] +
			" processlist"},
		errmsg: []string{"Unable to connect to the configured MySQL database"},
		hard:   []bool{false},
	}
	err := tryCmds(cmdLogger, DBCmds)
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
		sql := "SELECT count(SCHEMA_NAME) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '" + dbTar.Name + "';"
		DBCk := osCmds{
			id: osTar,
			cmds: []string{"mysql --host=" + dbTar.Host +
				" --user=" + creds["user"] +
				" --port=" + strconv.Itoa(dbTar.Port) +
				" --password=" + creds["pass"] +
				" --execute=\"" + sql + "\""},
			errmsg: []string{"Unable to connect to the configured MySQL database"},
			hard:   []bool{false},
		}
		out, err := inspectCmd(cmdLogger, DBCk.cmds[0], DBCk.errmsg[0], DBCk.hard[0])
		if err != nil {
			traceMsg("validation of connection to MySQL failed")
			return err
		}
		// Clean up stdout from inspectCmd output
		resp := strings.Trim(
			strings.ReplaceAll(
				strings.ReplaceAll(out, "count(SCHEMA_NAME)", ""), "\n", ""), " ")

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
			sql := "DROP DATABASE " + dbTar.Name + ";"
			traceMsg(fmt.Sprintf("%+v\n", sql))
			// TODO: Convert this and the above call to a function
			DropDB := osCmds{
				id: osTar,
				cmds: []string{"mysql --host=" + dbTar.Host +
					" --user=" + creds["user"] +
					" --port=" + strconv.Itoa(dbTar.Port) +
					" --password=" + creds["pass"] +
					" --execute=\"" + sql + "\""},
				errmsg: []string{"Unable to drop the existing MySQL database"},
				hard:   []bool{false},
			}
			err := tryCmd(cmdLogger, DropDB.cmds[0], DropDB.errmsg[0], DropDB.hard[0])
			if err != nil {
				traceMsg("Failed to drop existing database per configured option to drop existing")
				return err
			}
			fmt.Printf("Existing database %+v dropped since Database Drop was set to %+v\n", dbTar.Name, dbTar.Drop)
		}

	}

	// Create the DefectDojo database if it doesn't already exist
	sql := "CREATE DATABASE IF NOT EXISTS " + dbTar.Name + "  CHARACTER SET UTF8;"
	traceMsg(fmt.Sprintf("%+v\n", sql))
	// TODO: Convert this and the above call to a function
	CreateDB := osCmds{
		id: osTar,
		cmds: []string{"mysql --host=" + dbTar.Host +
			" --user=" + creds["user"] +
			" --port=" + strconv.Itoa(dbTar.Port) +
			" --password=" + creds["pass"] +
			" --execute=\"" + sql + "\""},
		errmsg: []string{"Unable to create a new MySQL database for DefectDojo"},
		hard:   []bool{false},
	}
	err = tryCmd(cmdLogger, CreateDB.cmds[0], CreateDB.errmsg[0], CreateDB.hard[0])
	if err != nil {
		traceMsg("Failed to create new database for DefectDojo to use")
		return err
	}

	// Drop user DefectDojo uses to connect to the database
	sql = "DROP USER '" + dbTar.User + "'@'localhost';DROP USER '" + dbTar.User + "'@'%';"
	traceMsg(fmt.Sprintf("%+v\n", sql))
	// TODO: Convert this and the above call to a function
	dropUsr := osCmds{
		id: osTar,
		cmds: []string{"mysql --host=" + dbTar.Host +
			" --user=" + creds["user"] +
			" --port=" + strconv.Itoa(dbTar.Port) +
			" --password=" + creds["pass"] +
			" --execute=\"" + sql + "\""},
		errmsg: []string{"Unable to delete existing database user for DefectDojo or one didn't exist"},
		hard:   []bool{false},
	}
	s, err := inspectCmd(cmdLogger, dropUsr.cmds[0], dropUsr.errmsg[0], dropUsr.hard[0])
	if err != nil {
		// No reason to return the error as this is expected for most cases
		// and create user will error out for edge cases
		traceMsg("Unable to delete existing database user for DefectDojo or one didn't exist")
		traceMsg(fmt.Sprintf("SQL DROP command output was %+v (in any)", s))
		traceMsg("Continuing after error deleting user, non-fatal error")
	}

	// If Drop DB, first delete any existing DD user
	// TODO: ^

	// First set the appropriate host for the DefectDojo user to connect from
	usrHost := "localhost"
	if !dbTar.Local && dbTar.Exists {
		// DB is remote and exists so localhost won't work
		usrHost = "%"
	}
	// Create user for DefectDojo to use to connect to the database
	sql = "CREATE USER '" + dbTar.User + "'@'" + usrHost + "' IDENTIFIED BY '" + dbTar.Pass + "';"
	traceMsg(fmt.Sprintf("%+v\n", sql))
	// TODO: Convert this and the above call to a function
	CreateUsr := osCmds{
		id: osTar,
		cmds: []string{"mysql --host=" + dbTar.Host +
			" --user=" + creds["user"] +
			" --port=" + strconv.Itoa(dbTar.Port) +
			" --password=" + creds["pass"] +
			" --execute=\"" + sql + "\""},
		errmsg: []string{"Unable to create a MySQL database user for DefectDojo"},
		hard:   []bool{false},
	}
	err = tryCmd(cmdLogger, CreateUsr.cmds[0], CreateUsr.errmsg[0], CreateUsr.hard[0])
	if err != nil {
		traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	// Grant the DefectDojo db user the necessary privileges
	sql = "GRANT ALL PRIVILEGES ON " + dbTar.Name + ".* TO '" + dbTar.User + "'@'" + dbTar.Host + "';"
	traceMsg(fmt.Sprintf("%+v\n", sql))
	// TODO: Convert this and the above call to a function
	grantPrivs := osCmds{
		id: osTar,
		cmds: []string{"mysql --host=" + dbTar.Host +
			" --user=" + creds["user"] +
			" --port=" + strconv.Itoa(dbTar.Port) +
			" --password=" + creds["pass"] +
			" --execute=\"" + sql + "\""},
		errmsg: []string{"Unable to grant needed privileges to database user for DefectDojo"},
		hard:   []bool{false},
	}
	err = tryCmd(cmdLogger, grantPrivs.cmds[0], grantPrivs.errmsg[0], grantPrivs.hard[0])
	if err != nil {
		traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	// Flush privileges to finalize changes to db
	sql = "FLUSH PRIVILEGES;"
	traceMsg(fmt.Sprintf("%+v\n", sql))
	// TODO: Convert this and the above call to a function
	flushPrivs := osCmds{
		id: osTar,
		cmds: []string{"mysql --host=" + dbTar.Host +
			" --user=" + creds["user"] +
			" --port=" + strconv.Itoa(dbTar.Port) +
			" --password=" + creds["pass"] +
			" --execute=\"" + sql + "\""},
		errmsg: []string{"Unable to flush database privileges"},
		hard:   []bool{false},
	}
	err = tryCmd(cmdLogger, flushPrivs.cmds[0], flushPrivs.errmsg[0], flushPrivs.hard[0])
	if err != nil {
		traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	return nil
}

func prepPostgreSQL(dbTar *config.DBTarget, os string) error {
	// Open a connection to the configured PostgreSQL database
	// https://godoc.org/github.com/lib/pq
	conn := "user=" + dbTar.User + " password=" + dbTar.Pass + " host=" + dbTar.Host + " port=" + strconv.Itoa(dbTar.Port)
	fmt.Println("DB conn is ", conn)

	dbPostgreSQL, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}

	fmt.Println(dbPostgreSQL)
	return nil
}
