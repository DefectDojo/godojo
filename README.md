# godojo
Golang installer for DefectDojo

### Dependencies

None, just download godojo for your platform + architechture and either:

1. Accept the default configuration
2. Edit dojoConfig.yml to meet your needs
3. Set environmental variable(s) to override the default configuration in dojoConfig.yml

### Assumptions

* Installer is run as root or with sudo like:

```
$ sudo godojo
```

* Installer can create a 'logs' directory where it run to write a log of the install
* Installer can create a file in the directory where it is run to save the run config
