package main

type dojoConfig struct {
	install  installConfig
	settings settingsConfig
}

type installConfig struct {
	// Installer settings
	ddDojoVer     string // Holds the version of Dojo to check out from the repo
	ddPrompt      bool
	ddInstallType string
	ddDBType      string
}

type settingsConfig struct {
	// Configs for settings.py
	ddDebug                         bool   // Run Dojo in debug mode, default false
	ddLoginRedirectURL              string // Where to redirect to after login, default "/"
	ddDjangoAdminEnabled            bool   // Is Django's admin enabled, default false
	ddSessionCookieHttponly         bool   // Add httponly cookie option to session, default true
	ddCsrfCookieHttponly            bool   // Add httponly cookie option to csrf, default true
	ddSecureSslRedirect             bool   // Redirect to https, default false
	ddSecureHstsInclueSubdomains    bool   // Add HSTS header to include subdomains
	ddSecureHstsSeconds             int    // Number of seconds for HSTS, default 31536000 aka 1 year
	ddCsrfCookieSecure              bool   // CSRF cookie has secure flag, default false
	ddSecureBrowserXSSFilter        bool   // TODO, default false
	ddTimeZone                      string // Time zone for Defect Dojo, default UTC
	ddLang                          string // Default language for Dojo, default en-us
	ddWkhtmltopdf                   string // path to wkhtmltopdf binary, default /usr/local/bin/wkhtmltopdf
	ddTeamName                      string // Name of the team in dojo, default "Security Team"
	ddAdmins                        string // admin email addresses, default "DefectDojo:dojo@localhost,Admin:admin@localhost"
	ddPortScanContactEmail          string // email address for port scans, default "email@localhost"
	ddPortScanResultEmailFrom       string // email address port scan results are sent from, default "email@localhost"
	ddPortScanExternalUnitEmailList string // TODO, default "email@localhost"
	ddPortScanSourceIP              string // What souce IP to use for port scanning, default "127.0.0.1"
	ddWhitenoise                    bool   // Should Whitenoise be used, default false
	ddTrackMigrations               bool   // TODO, default false
	ddSecureProxySslHeader          bool   // TODO, default false
	ddTestRunner                    string // TODO, default "django.test.runner.DiscoverRunner"
	ddURLPrefix                     string // Prefix to the Dojo URL, default ""
	ddRoot                          string // TODO, default "dojo"
	ddLanguageCode                  string // TODO, default "en-us"
	ddSiteID                        int    // TODO, default 1
	ddUseI18N                       bool   // TODO, default true
	ddUseL10n                       bool   // TODO, default true
	ddUseTz                         bool   // TODO, default true
	ddMediaURL                      string // URL prefix for Dojo media files, default "/media/"
	ddMediaRoot                     string // File path to the directory to store media files, default ddRoot + "media"
	ddStaticURL                     string // URL prefix for Dojo static files, default "/static/"
	ddStaticRoot                    string // File path to the directory to store static files, default ddRoot + "static"
	ddCeleryBrokerURL               string // TODO, default ""
	ddCeleryBrokerScheme            string // TODO, default "sqla+sqlite"
	ddCeleryBrokerUser              string // TODO, default ""
	ddCeleryBrokerPassword          string // TODO, default ""
	ddCeleryBrokerHost              string // TODO, default ""
	ddCeleryBrokerPort              int    // TODO, default -1
	ddCeleryBrokerPath              string // TODO, default "/dojo.celerydb.sqlite"
	ddCeleryTaskIgnoreResult        bool   // TODO, default true
	ddCeleryResultBackend           bool   // TODO, default "django-db"
	ddCeleryResultExpires           int    // TODO, default 86400 aka 24 hours
	ddCeleryBeatScheduleFilename    string // TODO, default ddRoot + "dojo.celery.beat.db"
	ddCeleryTaskSerializer          string // TODO, default "pickle"
	ddForceLowercaseTags            bool   // Force all tags to be lowercase, default true
	ddMaxTagLength                  int    // Max length for a tag in Dojo, default 25
	ddDatabaseEngine                string // Database for Dojo to use, default "django.db.backends.mysql"
	ddDatabaseHost                  string // Database host name, default "mysql"
	ddDatabaseName                  string // Database name for Dojo to use, default "defectdojo"
	ddTestDatabaseName              string // Database name used for testing, default "test_defectdojo"
	ddDatabasePassword              string // Dojo database user password, default "defectdojo"
	ddDatabasePort                  int    // Port that database is listening on, default 3306
	ddDatabaseUser                  string // Database user for the Dojo database, default "defectdojo"
	ddSecretKey                     string // TODO, default "."
	ddCredentialAES256Key           string // TODO, default "."
	ddDataUploadMaxMemorySize       int    // Max post size, default 8388608 aka 8 MB
	ddSocialAuthGoogleOAuth2Key     string // TODO, default ""
	ddSocialAuthGoogleOAuth2Secret  string // TODO, default ""
	ddSocialAuthOktaOAuth2Key       string // TODO, default ""
	ddSocialAuthOktaOAuthSecret     string // TODO, default ""
	ddSocialAuthOktaOAuthAPIURL     string // TODO, default "https://{your-org-url}/oauth2/default"
}

func testConfig() string {
	return "Worked"
}

func readInstallConfig() error {
	// Read the config

	// Set the default values

	return nil
}
