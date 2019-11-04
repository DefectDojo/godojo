package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Commands to bootstrap Ubuntu for the installer
func ubuntuInitOSInst(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
		b.id = "ubuntu:18.04"
		b.cmds = []string{
			"DEBIAN_FRONTEND=noninteractive apt-get update",
			"DEBIAN_FRONTEND=noninteractive apt-get -y upgrade",
			fmt.Sprintf("curl -sS %s | apt-key add - >/dev/null 2>&1", YarnGPG),
			fmt.Sprintf("echo -n %s > /etc/apt/sources.list.d/yarn.list", YarnRepo),
			fmt.Sprintf("curl -sL %s | sudo -E bash", NodeURL),
			"apt install -y apt-transport-https libjpeg-dev gcc libssl-dev python3-dev python3-pip nodejs yarn build-essential",
		}
		b.errmsg = []string{
			"Unable to update apt database",
			"Unable to upgrade OS packages with apt",
			"Unable to obtain the gpg key for Yarn",
			"Unable to add yard repo as an apt source",
			"Unable to install node",
			"Installing OS packages with apt failed",
		}
		b.hard = []bool{
			true,
			true,
			true,
			true,
			true,
			true,
		}
		// Currently, only Ubuntu 18.04 is supported
	}
	return
}

// Commands to install SQLite on Ubuntu
func ubuntuInstSQLite(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
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
	return
}

// Commands to install MariaDB on Ubuntu
func ubuntuInstMariaDB(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
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
	return
}

// Commands to install MySQL on Ubuntu
func ubuntuInstMySQL(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
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
	return
}

// Commands to install PostgreSQL on Ubuntu
func ubuntuInstPostgreSQL(id string, b *osCmds) {
	switch id {
	case "ubuntu:18.04":
		b.id = id
		b.cmds = []string{
			"DEBIAN_FRONTEND=noninteractive apt-get install -y libpq-dev postgresql postgresql-contrib",
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

// Determine the default creds for a database freshly installed in Ubuntu
func ubuntuDefaultDBCreds(db string, creds map[string]string) {
	// Installer currently assumes the default DB passwrod handling won't change by release
	// Switch on the DB type
	switch db {
	case "MySQL":
		ubuntuDefaultMySQL(creds)
	}

	return
}

func ubuntuDefaultMySQL(c map[string]string) {
	// Sent some intial values that ensure the connection will fail if the file read fails
	c["user"] = "debian-sys-maint"
	c["pass"] = "FAIL"

	// Pull the debian-sys-maint creds from /etc/mysql/debian.cnf
	f, err := os.Open("/etc/mysql/debian.cnf")
	if err != nil {
		// Exit with error code if we can't read the default creds file
		errorMsg("Unable to read file with defautl credentials, cannot continue")
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
		errorMsg("Unable to scan file with defautl credentials, cannot continue")
		os.Exit(1)
	}

}
