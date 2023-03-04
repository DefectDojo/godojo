package cmd

func Main() {
	// Set godojo defaults
	defaults := DDConfig{}
	defaults.setGodojoDefaults()

	// Prepeare the installer
	prepInstaller(&defaults)

	// Start the installation
	run(&defaults)
}
