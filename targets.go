package main

var InstallTargets = map[string][]string{
	"ubuntu": {"18.04", "18.10", "19.04"},
	"debian": {"stretch", "buster"},
}

// Supported OSes
type targetOS struct {
	id      string
	distro  string
	release string
}

type bootstrap struct {
	id     string   // Holds distro + release e.g. ubuntu:18.04
	cmds   []string // Holds the bootstrap commands
	errmsg []string // Holds the error messages if the matching command fails
	hard   []bool   // Flag to know if an error on the matching command is fatal
}

func initBootstrap(id string, b *bootstrap) {
	switch id {
	case "ubuntu:18.04":
		b.id = "ubuntu:18.04"
		b.cmds = []string{
			"DEBIAN_FRONTEND=noninteractive apt-get update",
			"DEBIAN_FRONTEND=noninteractive apt-get -y upgrade",
			"DEBIAN_FRONTEND=noninteractive apt-get -y install python3 python3-virtualenv ca-certificates",
			//"DEBIAN_FRONTEND=noninteractive apt -y install curl python3 python3-virtualenv expect wget gnupg2",
		}
		b.errmsg = []string{
			"Unable to update apt database",
			"Unable to upgrade OS packages with apt",
			"Unable to install prerequisites for installer via apt",
		}
		b.hard = []bool{
			true,
			true,
			true,
		}

		return

	}
}

//var (
//	u1804 = Bootstrap{
//		"ubuntu:18.04",
//		{
//			"DEBIAN_FRONTEND=noninteractive apt update",
//			"DEBIAN_FRONTEND=noninteractive apt -y upgrade",
//			"DEBIAN_FRONTEND=noninteractive apt -y install curl python3 python3-virtualenv expect wget gnupg2",
//		},
//		{
//			"Unable to update apt database",
//			"Unable to upgrade OS packages with apt",
//			"Unable to install prerequisites for installer via apt",
//		},
//		{
//			true,
//			true,
//			true,
//		},
//	}
//)

//var (
//	u1804 = Bootstrap{
//		id: "ubuntu:18.04",
//		cmds: {
//			"DEBIAN_FRONTEND=noninteractive apt update",
//			"DEBIAN_FRONTEND=noninteractive apt -y upgrade",
//			"DEBIAN_FRONTEND=noninteractive apt -y install curl python3 python3-virtualenv expect wget gnupg2",
//		},
//		errmsg: {
//			"Unable to update apt database",
//			"Unable to upgrade OS packages with apt",
//			"Unable to install prerequisites for installer via apt",
//		},
//		hard: {
//			true,
//			true,
//			true,
//		},
//	}
//)
