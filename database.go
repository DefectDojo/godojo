package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/mtesauro/godojo/config"
)

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
			dCmd.id = osTar
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
			dCmd.id = osTar
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
			dCmd.id = osTar
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
	return nil
}

func prepMariaDB(dbTar *config.DBTarget, os string) error {
	// TODO - test that this works MariaDB sd
	// Open a connection the the configured MySQL DB
	// https://github.com/go-sql-driver/mysql/#dsn-data-source-name
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	conn := dbTar.User + ":" + dbTar.Pass + "@" + dbTar.Host + ":" + strconv.Itoa(dbTar.Port)
	fmt.Println("DB conn is ", conn)
	dbMySQL, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	fmt.Println(dbMySQL)

	return nil
}

func prepMySQL(dbTar *config.DBTarget, os string) error {
	// Open a connection the the configured MySQL DB
	// https://github.com/go-sql-driver/mysql/#dsn-data-source-name

	// Generate a connection string like:
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	conn := ""
	if dbTar.Local && !dbTar.Exists {
		// Determine default access for fresh install of that OS
		// AKA databse is local and didn't exist before the install
		creds := defaultDBCreds(dbTar.Engine, os)
		conn = creds["user"] + ":" + creds["pass"] + "@tcp(" + dbTar.Host + ":" + strconv.Itoa(dbTar.Port) + ")/mysql"
	} else {
		// If the database is remote or pre-existing and local, there's
		// no way for the installer to reliably determine the correct creds
		// so it must rely on the provided DB root user login creds
		conn = dbTar.Ruser + ":" + dbTar.Rpass + "@tcp(" + dbTar.Host + ":" + strconv.Itoa(dbTar.Port) + ")/mysql"
	}

	// User the connction string above to open a DB connection
	fmt.Println(conn)
	dbMySQL, err := sql.Open("mysql", conn)
	if err != nil {
		traceMsg("Unable to run sql.Open against MySQL, connection string was:")
		traceMsg(conn)
		return err
	}

	// Create a context to use with following queries that has a 3 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Ping the database to extablish a connection to it - give the DB 3 seconds to respond
	err = dbMySQL.PingContext(ctx)
	if err != nil {
		traceMsg("Attempt to ping MySQL database failed")
		return err
	}

	// Drop existing DefectDojo database if it exists and configuration says to
	if dbTar.Drop {
		fmt.Println("INSIDE DROP IF")
		// Query MySQL to see if the configured database name exists already
		sql := "SELECT count(SCHEMA_NAME) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '" + dbTar.Name + "';"
		rows, err := dbMySQL.QueryContext(ctx, sql)
		if err != nil {
			traceMsg("Attempt to query MySQL database for the configured database name failed")
			return err
		}
		defer rows.Close()

		// Get the count from the query above - should be a single row returned
		_ = rows.Next()
		var r int
		if err := rows.Scan(&r); err != nil {
			traceMsg("Attempt to scan rows from MySQL for database name count database failed")
			return err
		}

		// If count is 1, we need to drop the configured databas in MySQL before moving on
		if r == 1 {
			// Drop the configured database
			sql := "DROP DATABASE " + dbTar.Name + ";"
			_, err := dbMySQL.ExecContext(ctx, sql)
			if err != nil {
				traceMsg("Attempt to drop existing database failed")
				fmt.Println("DOH!")
				return err
			}
		}
	}

	// Create the DefectDojo database if it doesn't already exist
	sql := "CREATE DATABASE IF NOT EXISTS " + dbTar.Name + "  CHARACTER SET UTF8;"
	result, err := dbMySQL.ExecContext(ctx, sql)
	if err != nil {
		traceMsg("Unable to create database for DefectDojo")
		return err
	}

	// Get a count of rows affected to ensure database was created correctly
	rows, err := result.RowsAffected()
	if err != nil {
		traceMsg("Failed to get the rows affected after creating the DefectDojo database")
		return err
	}
	if rows != 1 {
		return errors.New("Error occured when creating DefectDojo database")
	}

	// Create user for DefectDojo to use to connect to the database
	// Note: setup.bash would drop the DefectDojo DB user here - I'm not going to because:
	// (1) If db is remote or existing, we're already using the root/superuser creds anyway and
	// (2) If db is local and new (aka existing=false), then there won't be a DefectDojo user
	sql = "CREATE USER '" + dbTar.User + "'@'" + dbTar.Host + "' IDENTIFIED BY '" + dbTar.Pass + "';"
	result, err = dbMySQL.ExecContext(ctx, sql)
	if err != nil {
		traceMsg("Unable to create database user for DefectDojo")
		return err
	}

	// Grant the DefectDojo db user the necessary privileges
	sql = "GRANT ALL PRIVILEGES ON " + dbTar.Name + ".* TO '" + dbTar.User + "'@'" + dbTar.Host + "';"
	result, err = dbMySQL.ExecContext(ctx, sql)
	if err != nil {
		traceMsg("Unable to grant database user privileges")
		return err
	}

	// Flush privileges to finalize changes to db
	sql = "FLUSH PRIVILEGES;"
	result, err = dbMySQL.ExecContext(ctx, sql)
	if err != nil {
		traceMsg("Unable to flush privileges")
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
