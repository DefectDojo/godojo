package main

// TODO:
// Add Cobra for command-line args - https://github.com/spf13/cobra
// Add go-git to pull Dojo source - https://github.com/src-d/go-git
import (
	"fmt"
	"strings"

	"github.com/mtesauro/godojo/config"
	"github.com/spf13/viper"
)

// TODO: Consider moving install-config.go into a subfolder and import it abovj

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
	err := config.ReadInstallConfig()
	if err != nil {
		panic("Can't read config nor setup default values")
	}

	// Setup viper config
	viper.AddConfigPath(".")
	viper.SetConfigName("dojoConfig")
	var conf config.DojoConfig

	// Setup ENV variables
	viper.SetEnvPrefix("DD")
	replace := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replace)
	viper.AutomaticEnv()

	// DEBUG
	//os.Setenv("DD_DOJO_VER", "6.6.6")

	err = viper.ReadInConfig()
	if err != nil {
		panic("TODO: Better error handling 01")
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic("TODO: Better error handling 02")
	}
	fmt.Printf("conf.Install.Version = %+v\n\n", conf.Install.Version)
	fmt.Printf("conf = %+v\n\n", conf)
	fmt.Printf("conf.Install = %+v\n\n", conf.Install)
	fmt.Printf("viper.Get(\"DojoVer\") = %+v\n", viper.Get("install.version"))
	fmt.Printf("viper.AllKeys = %+v\n\n", viper.AllKeys())
	fmt.Printf("viper.AllSettings = %+v\n\n", viper.AllSettings())
	fmt.Printf("FROM VIPER viper.Get(\"install.db.engine\") = %+v\n\n", viper.Get("install.db.engine"))
	fmt.Printf("FROM STRUCT config.Install.db.engine= %+v\n", conf.Install.DB.Engine)
	// Stopped here
	fmt.Println(config.TestConfig())

	// DEBUG - Testing writing config
	err = viper.WriteConfigAs("runtime-install-config.yml")
	if err != nil {
		fmt.Printf("err from write was: %+v\n", err)
		//panic("TODO: Better error handling 03")
	}

	// Start stub'ing out stuff
	// Look at setup.bash's high-level workflow
}
