package config

type DojoConfig struct {
	Install  InstallConfig
	Settings SettingsConfig
}

type InstallConfig struct {
	// Installer settings
	Version       string // Holds the version of Dojo to check out from the repo
	SourceInstall bool   // If true, do a source install instead of a versioned release
	SourceBranch  string // Branch to checkout for a source install, if SourceCommit isn't "", SourceBranch will be ignored
	SourceCommit  string // head or full commit hash to install a specific commit, SourceBranch will be ignored if this isn't ""
	Quiet         bool   // If true, suppress all output except for very early errors - logs will still be written in the log directory
	Trace         bool   // If true, log at the trace level
	Prompt        bool   // Prompt at run time for install config.  If true, user will be prompted
	Set           string // The install set or type: Single Server, Dev, Stand-alone
	Root          string // Install root defaults to /opt/dojo
	Source        string // Directory to put the Dojo souce, child directory of Root
	Files         string // Directory for locally generated files like uploads, static, media, etc
	App           string // Directory where the Dojo Django app lives inside of Source above
	Sampledata    bool   // Install the sample data if true, defaults to false
	DB            DBTarget
}

type DBTarget struct {
	Engine string
	Local  bool
}

type SettingsConfig struct {
	// Configs for settings.py
	Debug                         bool   // Run Dojo in debug mode, default false
	LoginRedirectURL              string // Where to redirect to after login, default "/"
	DjangoAdminEnabled            bool   // Is Django's admin enabled, default false
	SessionCookieHttponly         bool   // Add httponly cookie option to session, default true
	CsrfCookieHttponly            bool   // Add httponly cookie option to csrf, default true
	SecureSslRedirect             bool   // Redirect to https, default false
	SecureHstsInclueSubdomains    bool   // Add HSTS header to include subdomains
	SecureHstsSeconds             int    // Number of seconds for HSTS, default 31536000 aka 1 year
	CsrfCookieSecure              bool   // CSRF cookie has secure flag, default false
	SecureBrowserXSSFilter        bool   // TODO, default false
	TimeZone                      string // Time zone for DefectDojo, default UTC
	Lang                          string // Default language for Dojo, default en-us
	Wkhtmltopdf                   string // path to wkhtmltopdf binary, default /usr/local/bin/wkhtmltopdf
	TeamName                      string // Name of the team in dojo, default "Security Team"
	Admins                        string // admin email addresses, default "DefectDojo:dojo@localhost,Admin:admin@localhost"
	PortScanContactEmail          string // email address for port scans, default "email@localhost"
	PortScanResultEmailFrom       string // email address port scan results are sent from, default "email@localhost"
	PortScanExternalUnitEmailList string // TODO, default "email@localhost"
	PortScanSourceIP              string // What souce IP to use for port scanning, default "127.0.0.1"
	Whitenoise                    bool   // Should Whitenoise be used, default false
	TrackMigrations               bool   // TODO, default false
	SecureProxySslHeader          bool   // TODO, default false
	TestRunner                    string // TODO, default "django.test.runner.DiscoverRunner"
	URLPrefix                     string // Prefix to the Dojo URL, default ""
	Root                          string // TODO, default "dojo"
	LanguageCode                  string // TODO, default "en-us"
	SiteID                        int    // TODO, default 1
	UseI18N                       bool   // TODO, default true
	UseL10n                       bool   // TODO, default true
	UseTz                         bool   // TODO, default true
	MediaURL                      string // URL prefix for Dojo media files, default "/media/"
	MediaRoot                     string // File path to the directory to store media files, default ddRoot + "media"
	StaticURL                     string // URL prefix for Dojo static files, default "/static/"
	StaticRoot                    string // File path to the directory to store static files, default ddRoot + "static"
	CeleryBrokerURL               string // TODO, default ""
	CeleryBrokerScheme            string // TODO, default "sqla+sqlite"
	CeleryBrokerUser              string // TODO, default ""
	CeleryBrokerPassword          string // TODO, default ""
	CeleryBrokerHost              string // TODO, default ""
	CeleryBrokerPort              int    // TODO, default -1
	CeleryBrokerPath              string // TODO, default "/dojo.celerydb.sqlite"
	CeleryTaskIgnoreResult        bool   // TODO, default true
	CeleryResultBackend           bool   // TODO, default "django-db"
	CeleryResultExpires           int    // TODO, default 86400 aka 24 hours
	CeleryBeatScheduleFilename    string // TODO, default ddRoot + "dojo.celery.beat.db"
	CeleryTaskSerializer          string // TODO, default "pickle"
	ForceLowercaseTags            bool   // Force all tags to be lowercase, default true
	MaxTagLength                  int    // Max length for a tag in Dojo, default 25
	DatabaseEngine                string // Database for Dojo to use, default "django.db.backends.mysql"
	DatabaseHost                  string // Database host name, default "mysql"
	DatabaseName                  string // Database name for Dojo to use, default "defectdojo"
	TestDatabaseName              string // Database name used for testing, default "test_defectdojo"
	DatabasePassword              string // Dojo database user password, default "defectdojo"
	DatabasePort                  int    // Port that database is listening on, default 3306
	DatabaseUser                  string // Database user for the Dojo database, default "defectdojo"
	SecretKey                     string // TODO, default "."
	CredentialAES256Key           string // TODO, default "."
	DataUploadMaxMemorySize       int    // Max post size, default 8388608 aka 8 MB
	SocialAuthGoogleOAuth2Key     string // TODO, default ""
	SocialAuthGoogleOAuth2Secret  string // TODO, default ""
	SocialAuthOktaOAuth2Key       string // TODO, default ""
	SocialAuthOktaOAuthSecret     string // TODO, default ""
	SocialAuthOktaOAuthAPIURL     string // TODO, default "https://{your-org-url}/oauth2/default"
}
