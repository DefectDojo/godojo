package main

// Location for all non-OS specific calls where case statements handle dispacting calls to OS specifc calls

func initOSInst(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
		ubuntuInitOSInst(id, b)

	}
	return
}

func instSQLite(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
		ubuntuInstSQLite(id, b)
	}
	return
}

func instMariaDB(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
		ubuntuInstMariaDB(id, b)
	}
	return
}

func instMySQL(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
		ubuntuInstMySQL(id, b)
	}
	return
}

func instPostgreSQL(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
		ubuntuInstPostgreSQL(id, b)
	}
	return
}

func defaultDBCreds(db string, os string) map[string]string {
	// Setup a map to return
	creds := map[string]string{"user": "foo", "pass": "bar"}

	// Get the default creds based on OS
	switch os {
	case "ubuntu:18.04":
		ubuntuDefaultDBCreds(db, creds)
	}

	return creds
}
