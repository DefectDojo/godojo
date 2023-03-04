package cmd

import (
	"flag"
	"fmt"
	"os"
)

// readArgs() takes no arguements and returns filled launchArgs struct unless
// there are errors in the arguments or the argument provided calls for an
// early exist such as --version or --help
func readArgs(d *DDConfig) {
	d.traceMsg("Called readArgs")
	// Read in the supported command-line options
	var version, help, v, h bool
	flag.BoolVar(&d.defInstall, "default", false, "Do an install based on default config values")
	flag.BoolVar(&version, "version", false, "Print the version and exit")
	flag.BoolVar(&v, "v", false, "Print the version and exit")
	flag.BoolVar(&help, "help", false, "Print the help message and exit")
	flag.BoolVar(&h, "h", false, "Print the help message and exit")
	flag.Parse()

	// Print help
	if help || h {
		printHelp()
		os.Exit(0)
	}
	// Print version
	if version || v {
		fmt.Printf("godojo version %s\n", d.ver)
		os.Exit(0)
	}

	// Handle special install case of default installs
	if d.defInstall {
		return
	}

	// See if the dojoConfig.yml is in the local directory
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine current working directory, exiting...")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	_, err = os.Stat(path + "/" + d.cf)
	if err != nil {
		// No config file found, so create one and exit
		writeDefaultConfig(d.cf, true)
	}

	d.traceMsg("Reached the end of readArgs")
}

// printHelp takes no arguements and prints godojo's help content to stdout
func printHelp() {
	// Output the help info
	fmt.Println("")
	fmt.Println("Usage of godojo")
	fmt.Println("")
	fmt.Println("./godojo [optional arguments]")
	fmt.Println("")
	fmt.Println("  [No arguments]")
	fmt.Println("        Check for a dojoConfig.yml file in the current working directory")
	fmt.Println("        If found, use those values to configure the installation")
	fmt.Println("        If NOT found, create a default dojoConfig.yml in the current working directory and exit")
	fmt.Println("  -default")
	fmt.Println("        OPTIONAL - Do an install based on the default dojoConfig.yml values")
	fmt.Println("                   Must be used alone and without other arguments")
	fmt.Println("  -help, -h")
	fmt.Println("        Print this help message and exit, ignoring all other arguments")
	fmt.Println("  -version, -v")
	fmt.Println("        Print the version and exit, ignoring all other arguments")
	fmt.Println("")
	fmt.Println("  Note #1: GNU-style arguments like --name are also supported")
	fmt.Println("")
	fmt.Println("  Note #2: Any of the configuration values can be overridden with an environmental variable")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("$ ./godojo")
	fmt.Println("     (Either creates a default config file or installs based on the config file in the same directory)")
	fmt.Println("$ ./godojo -dev")
	fmt.Println("     (Does a dev aka development/test install using known and fixed values for the installation")
	// TODO Consider an example of overriding with an env variable
	fmt.Println("")
}
