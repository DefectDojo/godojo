# godojo
Golang installer for DefectDojo

### Dependencies

None, just download the godojo release (TBD) for your platform + architechture and either:

1. Accept the default configuration
2. Edit dojoConfig.yml to meet your needs
3. Set environmental variable(s) to override the default configuration in dojoConfig.yml

### Assumptions

* Installer is run as root or with sudo like:

```
$ sudo godojo
```
or
```
# godojo
```

* Installer can create a 'logs' directory where the installer is run to write a log of the install
* Installer can create a file in the directory where it is run to save the runtime config
* Installer can create a base directory for the DefectDojo install (default is /opt/dojo).
