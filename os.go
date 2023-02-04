package main

import (
	"fmt"
	"strconv"

	"github.com/mtesauro/godojo/config"
)

// Location for all non-OS specific calls where case statements handle dispacting calls to OS specifc calls

func initOSInst(id string, b *osCmds) {
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
		ubuntuInitOSInst(id, b)

	}
	return
}

func instSQLite(id string, b *osCmds) {
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
		ubuntuInstSQLite(id, b)
	}
	return
}

func instMariaDB(id string, b *osCmds) {
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
		ubuntuInstMariaDB(id, b)
	}
	return
}

func instMySQL(id string, b *osCmds) {
	traceMsg(fmt.Sprintf("Installing MySQL for %s\n", id))
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
		ubuntuInstMySQL(id, b)
	}
	return
}

func instPostgreSQL(id string, b *osCmds) {
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
		ubuntuInstPostgreSQL(id, b)
	}
	return
}

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

func defaultDBCreds(db *config.DBTarget, os string) map[string]string {
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
		ubuntuDefaultDBCreds(db, creds)
	}

	return creds
}

func osPrep(id string, inst *config.InstallConfig, cmds *osCmds) {
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
		ubuntuOSPrep(id, inst, cmds)
	}
	return
}

func createSettingsPy(id string, inst *config.DojoConfig, cmds *osCmds) {
	// Setup the env.prod file used by settings.py

	// Create the database URL for the env file - https://github.com/kennethreitz/dj-database-url
	dbURL := ""
	switch inst.Install.DB.Engine {
	case "SQLite":
		// sqlite:///PATH
		dbURL = "sqlite:///defectdojo.db"
	case "MySQL":
		// mysql://USER:PASSWORD@HOST:PORT/NAME
		dbURL = "mysql://" + inst.Install.DB.User + ":" + inst.Install.DB.Pass + "@" + inst.Install.DB.Host + ":" +
			strconv.Itoa(inst.Install.DB.Port) + "/" + inst.Install.DB.Name
	case "PostgreSQL":
		// postgres://USER:PASSWORD@HOST:PORT/NAME
		dbURL = "postgres://" + inst.Install.DB.User + ":" + inst.Install.DB.Pass + "@" + inst.Install.DB.Host + ":" +
			strconv.Itoa(inst.Install.DB.Port) + "/" + inst.Install.DB.Name
	}

	// Setup env file for production
	genAndWriteEnv(inst, dbURL)

	// Create a settings.py for Dojo to use
	cmds.id = id
	cmds.cmds = []string{
		"ln -s " + inst.Install.Root + "/django-DefectDojo/dojo/settings " +
			inst.Install.Root + "/customizations",
		"chown " + inst.Install.OS.User + "." + inst.Install.OS.Group + " " + inst.Install.Root +
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

func setupDjango(id string, inst *config.DojoConfig, cmds *osCmds) {
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
		ubuntuSetupDDjango(id, &inst.Install, cmds)
	}
	return
}
