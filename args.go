package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type LaunchArgs struct {
	Dev      bool
	Default  bool
	Componse bool // TODO: Implement this
	K8s      bool // TODO: Implement this
}

func readArgs() LaunchArgs {
	//TODO TUrn me to a log writer fmt.Println("Called readArgs")
	// Read in the supported command-line options
	var version, help, v, h bool
	opts := LaunchArgs{}
	flag.BoolVar(&opts.Default, "default", false, "Do an install based on default config values")
	flag.BoolVar(&opts.Dev, "dev", false, "Do a development install with known config values")
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
		fmt.Printf("godojo version %s\n", ver)
		os.Exit(0)
	}

	// Handle double options (dev and default currently)
	if opts.Dev && opts.Default {
		// Only 1 option should be provided, not both
		fmt.Println("Error: godojo only supports a single install option")
		fmt.Println("Both -dev and --default have been provided so exiting.")
		os.Exit(1)
	}

	// Handle special install cases of default and dev
	if opts.Default || opts.Dev {
		return opts
	}

	// See if the dojoConfig.yml is in the local directory
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine current working directory, exiting...")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	_, err = os.Stat(path + "/" + cf)
	if err != nil {
		// No config file found, so create one and exit
		createDefaultConfig(cf, true)
	}
	traceMsg("Reached the end of readArgs")
	return opts
}

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
	fmt.Println("  -dev")
	fmt.Println("        OPTIONAL - Do an dev install with fixed values especially for testing")
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
	fmt.Println("")
}

func createDefaultConfig(c string, ex bool) {
	// Get the current working directory for future operations
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine current working directory, exiting...")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	data, err := Asset(emdir + c)
	if err != nil {
		// Asset was not found.
		fmt.Println("Default dojoConfig.yml was not found")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Write out the embedded default dojoConfig.yml
	err = ioutil.WriteFile(path+"/"+c, data, 0644)
	if err != nil {
		// Cannot write config file
		fmt.Printf("Unable to write configuration file in %s, exiting...\n", path)
		fmt.Printf("Error: %v\n", err)
	}

	if ex {
		fmt.Println("\nNOTE: A dojoConfig.yml file was not found in the current directory:")
		fmt.Printf("\t%s\nA default configuration file was written there.\n\n", path)
		fmt.Println("Please review the configuration settings, adjusting as needed and")
		fmt.Println("re-run the godojo installer to begin the install you configured.")
		os.Exit(0)
	}
}
