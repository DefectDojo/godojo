package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// readConfigFile reads the yaml configuration file for godojo
// to determine runtime configuration.  The file is dojoConfig.yml
// and is expected to be in the same directory as the godojo binary
// It returns nohing but will exit early with a exit code of 1
// if there are errors reading the file or unmarshialling into a struct
func readConfigFile() {
	// Setup viper config
	viper.AddConfigPath(".")
	viper.SetConfigName("dojoConfig")

	// Setup ENV variables
	// TODO: Do these manually in readEnvVars() since they have odd names for Viper auto-magic
	//viper.SetEnvPrefix("DD")
	//replace := strings.NewReplacer(".", "_")
	//viper.SetEnvKeyReplacer(replace)
	//viper.AutomaticEnv()

	// Read the default config file dojoConfig.yml
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("")
		fmt.Println("Unable to read the godojo config file (dojoConfig.yml), exiting install")
		os.Exit(1)
	}
	// Marshall the config values into the DojoConfig struct
	err = viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println("")
		fmt.Println("Unable to set the config values based on config file and ENV variables, exiting install")
		os.Exit(1)
	}
}

// readEnvVars reads the DefectDojo supported environmental variables and
// overrides any options set in the configuration file. These variables
// are used to supply either install-time configurations or provide values
// that are used in DefectDojo's settings.py configuration file
func readEnvVars() {
}
