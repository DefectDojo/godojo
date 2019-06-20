package main

// TODO:
// Add Cobra for command-line args - https://github.com/spf13/cobra
// Add go-git to pull Dojo source - https://github.com/src-d/go-git
import (
	"fmt"

	"github.com/spf13/viper"
)

// TODO: Consider moving install-config.go into a subfolder and import it abovj
//type installConfig struct {
//	dojoVer string
//}

// getDojo retrives the supplied version of DefectDojo from the Git repo
// and places it in the specified dojoSource directory (default is /opt/dojo)
func getDojo(v string) string {
	return v
}

func main() {
	//inst := installConfig{dojoVer: "1.5.3.1"}
	fmt.Println("Getting dojo from repository")
	ver := getDojo("7.7.7")
	fmt.Printf("ver = %+v\n", ver)

	// Read the config and set default values
	err := readInstallConfig()
	if err != nil {
		panic("Can't read config nor setup default values")
	}

	// Setup viper config
	viper.AddConfigPath(".")
	viper.SetConfigName("dojoConfig")
	var conf dojoConfig

	err = viper.ReadInConfig()
	if err != nil {
		panic("TODO: Better error handling 01")
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic("TODO: Better error handling 02")
	}
	fmt.Printf("conf.install.ddDojoVer = %+v\n", conf.install.ddDojoVer)
	fmt.Printf("conf = %+v\n", conf)
	fmt.Printf("conf.install = %+v\n", conf.install)
	// Stopped here
	fmt.Println(testConfig())

	// Start stub'ing out stuff
	// Look at setup.bash's high-level workflow
}
