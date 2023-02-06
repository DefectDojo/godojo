package cmd

func Main() {
	// Set godojo defaults
	defaults := gdjDefault{}
	defaults.setGodojoDefaults()

	// Prepeare the installer
	prepInstaller(&defaults)

	// Start the installation
	run(&defaults)
}
