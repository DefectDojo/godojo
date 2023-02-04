package main

import (
	"bufio"
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

func installDBClient(osTar string, dbTar *config.DBTarget, dCmd *osCmds) {
	// Look at the dbTar and call function to install that DB target
	switch dbTar.Engine {
	case "SQLite":
		// Generate commands to install SQLite
		// A remote SQLite DB makes no sense
		// TODO: Log this error
		return
	case "MariaDB":
		// Generate commands to install MariaDB
		// TODO: Write install for MariaDB client
		//instMariaDBClient(osTar, dCmd)
		return
	case "MySQL":
		// Generate commands to install MySQL
		// TODO: Write install for MySQL client
		//instMySQLClient(osTar, dCmd)
		return
	case "PostgreSQL":
		// Generate commands to install PostgreSQL
		instPostgreSQLClient(osTar, dCmd)
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
	//       * mysql
	// TODO: Check MySQL version and handle MySQL 8 and password format change

	// Set Creds based on dojoConfig.yml
	creds := map[string]string{"user": dbTar.Ruser, "pass": dbTar.Rpass}
	traceMsg(fmt.Sprintf("DB Creds from config are %s / %s", creds["user"], creds["pass"]))

	// Creds are unknown if DB is local and newly installed by godojo
	if dbTar.Local && !dbTar.Exists {
		// Determine default access for fresh install of that OS
		// AKA databse is local and didn't exist before the install
		creds = defaultDBCreds(dbTar, osTar)
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
		traceMsg("Failed to grant privileges to database user for DefectDojo")
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

	// Adjust requirements.txt to only use MySQL Python modules
	err = trimRequirementsTxt("MySQL")
	if err != nil {
		traceMsg("Unable to adjust requirements.txt for PostgreSQL usage, exiting")
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

func prepPostgreSQL(dbTar *config.DBTarget, osTar string) error {
	// TODO: Path check any binaries called
	//       * postgres

	// Set Creds based on dojoConfig.yml
	creds := map[string]string{"user": dbTar.Ruser, "pass": dbTar.Rpass}
	traceMsg(fmt.Sprintf("DB Creds from config are %s / %s", creds["user"], creds["pass"]))

	// Creds are unknown if DB is local and newly installed by godojo
	if dbTar.Local && !dbTar.Exists {
		// Determine default access for fresh install of that OS
		// AKA databse is local and didn't exist before the install
		creds = defaultDBCreds(dbTar, osTar)
		addRedact(creds["pass"])
	}
	traceMsg(fmt.Sprintf("DB Creds are now %s / %s", creds["user"], creds["pass"]))

	// Use pg_isready to check connectivity to PostgreSQL DB
	statusMsg("Checking connectivity to PostgreSQL")

	t, err := isPgReady(dbTar, creds)
	if err != nil {
		traceMsg(fmt.Sprintf("PostgreSQL is not available, error was %+v", err))
		return err
	}
	traceMsg(fmt.Sprintf("Output from pgReady: %+v", t))

	statusMsg("Validating DB connection settings")
	// Check connectivity to DB
	conCk := sqlStr{
		os:     osTar,
		sql:    "\\l",
		errMsg: "Unable to connect to the configured Postgres database",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(dbTar, conCk)
	if err != nil {
		traceMsg("validation of connection to Postgres failed")
		// TODO Fix this validation bypass
		//return err
	}

	// Drop existing DefectDojo database if it exists and configuration says to
	if dbTar.Drop {
		traceMsg("Dropping any existing database per Install.DB.Drop=True in dojoConfig.yml")
		// Query PostgreSQL to see if the configured database name exists already
		dbCk := sqlStr{
			os:     osTar,
			sql:    "\\l",
			errMsg: "Unable to check for existing DefectDojo PostgreSQL database",
			creds:  creds,
			kind:   "inspect",
		}
		out, err := runPgSQLCmd(dbTar, dbCk)
		if err != nil {
			traceMsg("Check for existing DefectDojo PostgreSQL database failed")
			fmt.Println("Drop database set to true but no database found, continuing")
		}

		// Clean up stdout from inspectCmd output
		strOut := squishSlice(out)
		ck := pgParseDBList(strOut, dbTar.Name)

		// Check if there's an existing DB
		// if ck = 0 then DB doesn't exist
		// if ck = 1 then the DB exists already and needs to be dropped first
		if ck == 1 {
			traceMsg("DB EXISTS so droping that sucker")
			dropDB := sqlStr{
				os:     osTar,
				sql:    "DROP DATABASE IF EXISTS " + dbTar.Name + ";",
				errMsg: "Unable to drop the existing PostgreSQL database",
				creds:  creds,
				kind:   "try",
			}
			_, err := runPgSQLCmd(dbTar, dropDB)
			if err != nil {
				traceMsg("Failed to drop existing database per configured option to drop existing")
				return err
			}
			fmt.Printf("Existing database %+v dropped since Database Drop was set to %+v\n", dbTar.Name, dbTar.Drop)
		}

	}

	// Create the DefectDojo database if it doesn't already exist
	traceMsg("Creating database for DefectDojo on PostgreSQL")
	createDB := sqlStr{
		os:     osTar,
		sql:    "CREATE DATABASE " + dbTar.Name + ";",
		errMsg: "Unable to create a new PostgreSQL database for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(dbTar, createDB)
	if err != nil {
		traceMsg("Failed to create new database for DefectDojo to use")
		return err
	}

	// Drop user DefectDojo uses to connect to the database
	traceMsg("Dropping existing DefectDojo PostgreSQL DB user, if any")
	dropUsr := sqlStr{
		os:     osTar,
		sql:    "DROP USER IF EXISTS " + dbTar.User + ";",
		errMsg: "Unable to delete existing database user for DefectDojo or one didn't exist",
		creds:  creds,
		kind:   "inspect",
	}
	out, err := runPgSQLCmd(dbTar, dropUsr)
	if err != nil {
		// No reason to return the error as this is expected for most cases
		// and create user will error out for edge cases
		traceMsg("Unable to delete existing database user for DefectDojo or one didn't exist")
		traceMsg(fmt.Sprintf("SQL DROP command output was %+v (in any)", squishSlice(out)))
		traceMsg("Continuing after error deleting user, non-fatal error")
	}

	// Create user for DefectDojo to use to connect to the database
	traceMsg("Creating PostgreSQL DB user for DefectDojo")
	createUsr := sqlStr{
		os:     osTar,
		sql:    "CREATE USER " + dbTar.User + " WITH ENCRYPTED PASSWORD '" + dbTar.Pass + "';",
		errMsg: "Unable to create a PostgreSQL database user for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(dbTar, createUsr)
	if err != nil {
		traceMsg("Failed to create database user for DefectDojo")
		return err
	}

	statusMsg("Note: pg_hba.conf has not been altered by godojo.")
	statusMsg("      It may need to be updated to allow DefectDojo to connect to the DB.")
	statusMsg("      Please consult the PostgreSQL documentation for further information.")

	// Grant the DefectDojo db user the necessary privileges
	traceMsg("Granting privileges to DefectDojo PostgreSQL DB user")
	grantPrivs := sqlStr{
		os:     osTar,
		sql:    "GRANT ALL PRIVILEGES ON DATABASE " + dbTar.Name + " TO " + dbTar.User + ";",
		errMsg: "Unable to grant needed privileges to database user for DefectDojo",
		creds:  creds,
		kind:   "try",
	}
	_, err = runPgSQLCmd(dbTar, grantPrivs)
	if err != nil {
		traceMsg("Failed to grant privileges to database user for DefectDojo")
		return err
	}

	// Adjust requirements.txt to only use MySQL Python modules
	err = trimRequirementsTxt("PostgreSQL")
	if err != nil {
		traceMsg("Unable to adjust requirements.txt for PostgreSQL usage, exiting")
		return err
	}

	return nil
}

func runPgSQLCmd(dbTar *config.DBTarget, c sqlStr) ([]string, error) {
	out := make([]string, 1)
	traceMsg(fmt.Sprintf("Postgres query: %s", c.sql))
	DBCmds := osCmds{
		id: c.os,
		cmds: []string{"sudo -u postgres PGPASSWORD=\"" + c.creds["pass"] + "\"" +
			" psql --host=" + dbTar.Host +
			" --username=" + c.creds["user"] +
			" --port=" + strconv.Itoa(dbTar.Port) +
			" --command=\"" + c.sql + "\""},
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
		traceMsg("Invalid 'kind' sent to runPgSQLCmd, bug in godojo")
		fmt.Println("Bug discovered in godojo, see trace message. Quitting.")
		os.Exit(1)
	}

	// Handle errors from running the PostgreSQL command
	if err != nil {
		traceMsg(fmt.Sprintf("Error running Posgres command - %s", c.sql))
		return out, err
	}

	return out, nil
}

func isPgReady(dbTar *config.DBTarget, creds map[string]string) (string, error) {
	traceMsg("isPgReady called")

	// Call ps_isready and check exit code
	pgReady := osCmds{
		id: "Linux", cmds: []string{"PGPASSWORD=\"" + creds["pass"] + "\" pg_isready" +
			" --host=" + dbTar.Host +
			" --username=" + creds["user"] +
			" --port=" + strconv.Itoa(dbTar.Port) + " "},
		errmsg: []string{"Unable to run pg_isready to validate PostgreSQL DB status."},
		hard:   []bool{false},
	}

	traceMsg(fmt.Sprintf("Running command: %+v", pgReady.cmds))

	out, err := inspectCmds(cmdLogger, pgReady)
	if err != nil {
		traceMsg(fmt.Sprintf("Error running pg_isready was: %+v", err))
		// TODO Fix this error bypass
		return squishSlice(out), nil
		//return "", err
	}

	return squishSlice(out), nil
}

// Parse a list of existng PostgreSQL DBs for a specific DB name
// if the DB name is found, return 1 else return 0
func pgParseDBList(tbl string, name string) int {
	traceMsg(fmt.Sprintf("Parsing DB list for existing DefectDojo DB named  %+v", name))

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
		if strings.Contains(v, name) {
			traceMsg("Match found, there is an existing DefectDojo DB")
			return 1
		}
	}

	traceMsg("No match found, no existing DefectDojo DB")
	return 0
}

func trimRequirementsTxt(dbUsed string) error {
	traceMsg("Called trimRequirementsTxt")

	req := make([]string, 1)
	switch dbUsed {
	case "MySQL":
		traceMsg("MySQL requirements fix")
		req[0] = "psycopg2-binary"
		// append any additional requirements
	case "PostgreSQL":
		traceMsg("PostgreSQL requirements fix")
		req[0] = "mysqlclient"
		// append any additional requirements
	default:
		return errors.New("Unknown database provided to trimRequirementsTxt()")
	}

	traceMsg("No error return from trimRequirementsTxt")
	return removeRequirement(req)
}

func removeRequirement(r []string) error {
	traceMsg("Called removeRequirement")

	// Set th path for requirements.txt
	p := conf.Install.Root + "/" + conf.Install.Source + "/"

	cyaCopy := osCmds{
		id:     "Linux",
		cmds:   []string{"cp " + p + "requirements.txt " + p + "requirements.cya"},
		errmsg: []string{"Unable to make a cya copy of requirements.txt"},
		hard:   []bool{false},
	}

	err := tryCmds(cmdLogger, cyaCopy)
	if err != nil {
		traceMsg("Failed creating a backup copy of requirements.txt")
		return err
	}

	// Open requirements.txt
	f, err := os.OpenFile(p+"requirements.txt", os.O_RDWR, 0744)
	if err != nil {
		traceMsg("Unable to open requirements.txt")
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
				traceMsg(fmt.Sprintf("Keeping requirements line: %+v", line))
				newf += line + "\n"
			} else {
				traceMsg(fmt.Sprintf("Skiping requirements line: %+v", line))
				skip += "#\n"
			}
		}

	}

	// Clearn the file contents
	err = f.Truncate(0)
	if err != nil {
		traceMsg("Unable to truncate requirements.txt file")
		return err
	}

	// Add a couple of extra lines just in case
	skip += "# File re-written by godojo\n"
	newf += skip

	// Write resulting file
	_, err = f.WriteAt([]byte(newf), 0)
	if err != nil {
		traceMsg("Unable to write new requirements.txt file")
		return err
	}

	traceMsg("No error return from removeRequirement")
	return nil
}
