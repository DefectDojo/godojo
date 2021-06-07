# godojo - an installer for DefectDojo

godojo is an installer for DefectDojo created as a much more powerful replacement for setup.bash.  It provides a way to complete a 'server' install of DefectDojo. This is a traditional installation where DefectDojo is installed on the disk of a server/VM as part of the running OS.

godojo simplifies installing DefectDojo since the only thing needed to complete the install is the godojo binary. The installer handles pulling the requested version of the source code and any needed dependencies.

godojo supports the following types of installations:

* Installing a [numbered release of DefectDojo](https://github.com/DefectDojo/django-DefectDojo/releases) e.g. version 1.14.1
* Installing head of a specific branch e.g. [dev branch](https://github.com/DefectDojo/django-DefectDojo/tree/dev)
* Installing a specific commit

The currently supported Linux distros and database configurations are listed [here](https://docs.google.com/spreadsheets/d/1HuXh3Zr4mrmb6_YmKkDgzl-ZINYZCvVZn31UCqIGpUA/edit?usp=sharing)

godojo is developed targeting .deb (Debian) based distributions especially Ubuntu but should work on any Debian-based distro.

DefectDojo also supports [other methods of installation](https://github.com/DefectDojo/django-DefectDojo#supported-installation-options) that are not covered by godojo.

### Dependencies

None, just download the most recent [godojo release](https://github.com/DefectDojo/godojo/releases) and either:

1. Accept the default configuration (one will be created for you the first time you run godojo)
2. Edit dojoConfig.yml to meet your needs then run godojo
3. Set environmental variable(s) to override the default configuration in dojoConfig.yml when you run godojo

The defaults in dojoConfig.yml are pretty sane. All you really need to do is:

* decide what version of DefectDojo you want to install (a release, branch or commit)
* set a password for the initial Admin user (Install > Admin > Pass).

You can see all the configuration options with descriptions in the [example config file](https://github.com/DefectDojo/godojo/blob/master/example_dojoConfig.yml).

Note: godojo is built with go version 1.16.3 (or newer)

### Assumptions / requirements

* Bash is available and in $PATH
* Installer is run as root or with sudo like:

```
$ sudo ./godojo
```
or
```
# ./godojo
```

* Installer can create a 'logs' directory where the installer is run to write a log of the install
* Installer can create a file in the directory where it is run to save the runtime config
* Installer can create a base directory for the DefectDojo install (default is /opt/dojo).
* Installer can download the source code for DefectDojo and it's dependencies (Internet access)

### Other benefits of godojo

* The same installer can install multiple versions of DefectDojo
* Supports both MySQL and PostgreSQL databases
* Supports creating a new database or using an existing database. Database can be local (same host) or remote.
* godojo doesn't care where it is run from - the only important location is where DefectDojo will be installed which defaults to /opt/dojo
* godojo creates logs in a 'logs' subdirectory in the directory where it is run.
  * Logs are configurable from none ("Quiet: true" in dojoConfig.yml) to trace ("Trace: true" in dojoConfig.yml)
* Any passwords, keys or other sensitive data is redacted in the logs by default ("Redact: true" in dojoConfig.yml)
* All dojoConfig.yml configuration items can be overridden with environmental variables at run time

### Example installation

If you don't have a dojoConfig.yml in the same directory as godojo (or this is your first install), one will be created for you:

```
$ sudo ./godojo

NOTE: A dojoConfig.yml file was not found in the current directory:
	/home/example
A default configuration file was written there.

Please review the configuration settings, adjusting as needed and
re-run the godojo installer to begin the install you configured.
```

Once you have a dojoConfig.yml you're happy with, just run godojo:

```
 sudo ./godojo
        ____       ____          __     ____          _
       / __ \___  / __/__  _____/ /_   / __ \____    (_)___
      / / / / _ \/ /_/ _ \/ ___/ __/  / / / / __ \  / / __ \ 
     / /_/ /  __/ __/  __/ /__/ /_   / /_/ / /_/ / / / /_/ /
    /_____/\___/_/  \___/\___/\__/  /_____/\____/_/ /\____/
                                               /___/
    version  1.1.1

  Welcome to godojo, the official way to install DefectDojo.
  For more information on how goDojo does an install, see:
  https://github.com/DefectDojo/godojo

==============================================================================
  Starting the dojo install at Sun Apr 25, 2021 05:43:44 UTC
==============================================================================


==============================================================================
  Determining OS for installation
==============================================================================

OS was determined to be Linux, Ubuntu:20.10
DefectDojo installation on this OS is supported, continuing

==============================================================================
  Bootstrapping the godojo installer
==============================================================================

Bootstrapping...(-*--------)
```

Note: The above is a snippet of godojo in action. With Quiet set to 'false' you will see output for the various stages godojo completes with a progress bar for each part.


