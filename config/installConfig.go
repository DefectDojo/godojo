package config

// DojoConfig - "mother" struct to hold all the config options
type DojoConfig struct {
	Install  InstallConfig
	Settings SettingsConfig
	Options  OptionalConfig
}

// InstallConfig - struct to hold the install time options
type InstallConfig struct {
	// Installer settings
	Version       string         // Holds the version of Dojo to check out from the repo
	SourceInstall bool           // If true, do a source install instead of a versioned release
	SourceBranch  string         // Branch to checkout for a source install, if SourceCommit isn't "", SourceBranch will be ignored
	SourceCommit  string         // head or full commit hash to install a specific commit, SourceBranch will be ignored if this isn't ""
	Quiet         bool           // If true, suppress all output except for very early errors - logs will still be written in the log directory
	Trace         bool           // If true, log at the trace level
	Redact        bool           // If true, redact sensitive information from being logged.  Defaults to true
	Prompt        bool           // Prompt at run time for install config.  If true, user will be prompted
	Mac           bool           // The install set or type: Single Server, Dev, Stand-alone
	Root          string         // Install root defaults to /opt/dojo
	Source        string         // Directory to put the Dojo souce, child directory of Root
	Files         string         // Directory for locally generated files like uploads, static, media, etc
	App           string         // Directory where the Dojo Django app lives inside of Source above
	Sampledata    bool           // Install the sample data if true, defaults to false
	DB            DBTarget       // struct for DB configuration values
	OS            OSTarget       // struct for DB configuration values
	Settings      SettingsTarget // struct for DB configuration values
	Admin         AdminTarget    // struct for DB configuration values
	PullSource    bool           // If false, installer won't download source code - primarily for debugging
}

// DBTarget - struct to hold Install.DB options
type DBTarget struct {
	Engine string
	Local  bool
	Exists bool
	Ruser  string
	Rpass  string
	Name   string
	User   string
	Pass   string
	Host   string
	Port   int
	Drop   bool
}

// OSTarget - struct to hold Install.OS options
type OSTarget struct {
	User  string
	Pass  string
	Group string
}

// SettingsTarget - struct to hold Install.Settings options
type SettingsTarget struct {
	Dist string
	File string
	Env  string
}

// AdminTarget - struct to hold Install.Admin options
type AdminTarget struct {
	User  string
	Pass  string
	Email string
}

// SettingsConfig - struct to hold the config values for settings.py
type SettingsConfig struct {
	AdminFirstName                        string `yaml:"AdminFirstName"`
	AdminLastName                         string `yaml:"AdminLastName"`
	AdminMail                             string `yaml:"AdminMail"`
	AdminPassword                         string `yaml:"AdminPassword"`
	Admins                                string `yaml:"Admins"`
	AdminUser                             string `yaml:"AdminUser"`
	AllowedHosts                          string `yaml:"AllowedHosts"`
	CeleryBeatScheduleFilename            string `yaml:"CeleryBeatScheduleFilename"`
	CeleryBrokerHost                      string `yaml:"CeleryBrokerHost"`
	CeleryBrokerPassword                  string `yaml:"CeleryBrokerPassword"`
	CeleryBrokerPath                      string `yaml:"CeleryBrokerPath"`
	CeleryBrokerPort                      int    `yaml:"CeleryBrokerPort"`
	CeleryBrokerScheme                    string `yaml:"CeleryBrokerScheme"`
	CeleryBrokerURL                       string `yaml:"CeleryBrokerURL"`
	CeleryBrokerUser                      string `yaml:"CeleryBrokerUser"`
	CeleryLogLevel                        string `yaml:"CeleryLogLevel"`
	CeleryResultBackend                   string `yaml:"CeleryResultBackend"`
	CeleryResultExpires                   int    `yaml:"CeleryResultExpires"`
	CeleryTaskIgnoreResult                bool   `yaml:"CeleryTaskIgnoreResult"`
	CeleryTaskSerializer                  string `yaml:"CeleryTaskSerializer"`
	CredentialAES256Key                   string `yaml:"CredentialAES256Key"`
	CSRFCookieHTTPOnly                    bool   `yaml:"CSRFCookieHTTPOnly"`
	CSRFCookieSecure                      bool   `yaml:"CSRFCookieSecure"`
	DatabaseEngine                        string `yaml:"DatabaseEngine"`
	DatabaseHost                          string `yaml:"DatabaseHost"`
	DatabaseName                          string `yaml:"DatabaseName"`
	DatabasePassword                      string `yaml:"DatabasePassword"`
	DatabasePort                          string `yaml:"DatabasePort"`
	DatabaseType                          string `yaml:"DatabaseType"`
	DatabaseURL                           string `yaml:"DatabaseURL"`
	DatabaseUser                          string `yaml:"DatabaseUser"`
	DataUploadMaxMemorySize               int    `yaml:"DataUploadMaxMemorySize"`
	Debug                                 bool   `yaml:"Debug"`
	DjangoAdminEnabled                    bool   `yaml:"DjangoAdminEnabled"`
	EmailURL                              string `yaml:"EmailURL"`
	Env                                   string `yaml:"Env"`
	EnvPath                               string `yaml:"EnvPath"`
	ForceLowercaseTags                    bool   `yaml:"ForceLowercaseTags"`
	Host                                  string `yaml:"Host"`
	Initialize                            string `yaml:"Initialize"`
	Lang                                  string `yaml:"Lang"`
	LanguageCode                          string `yaml:"LanguageCode"`
	LoginRedirectURL                      string `yaml:"LoginRedirectURL"`
	MaxTagLength                          int    `yaml:"MaxTagLength"`
	MediaRoot                             string `yaml:"MediaRoot"`
	MediaURL                              string `yaml:"MediaURL"`
	Port                                  string `yaml:"Port"`
	PortScanContactEmail                  string `yaml:"PortScanContactEmail"`
	PortScanExternalUnitEmailList         string `yaml:"PortScanExternalUnitEmailList"`
	PortScanResultEmailFrom               string `yaml:"PortScanResultEmailFrom"`
	PortScanSourceIP                      string `yaml:"PortScanSourceIP"`
	Root                                  string `yaml:"Root"`
	SecretKey                             string `yaml:"SecretKey"`
	SecureBrowserXSSFilter                bool   `yaml:"SecureBrowserXSSFilter"`
	SecureContentTypeNosniff              string `yaml:"SecureContentTypeNosniff"`
	SecureHSTSIncludeSubdomains           bool   `yaml:"SecureHSTSIncludeSubdomains"`
	SecureHSTSSeconds                     int    `yaml:"SecureHSTSSeconds"`
	SecureProxySSLHeader                  bool   `yaml:"SecureProxySSLHeader"`
	SecureSSLRedirect                     bool   `yaml:"SecureSSLRedirect"`
	SessionCookieHTTPOnly                 bool   `yaml:"SessionCookieHTTPOnly"`
	SessionCookieSecure                   bool   `yaml:"SessionCookieSecure"`
	SiteID                                int    `yaml:"SiteID"`
	SocialAuthAzureadTenantOauth2Enabled  string `yaml:"SocialAuthAzureadTenantOauth2Enabled"`
	SocialAuthAzureadTenantOauth2Key      string `yaml:"SocialAuthAzureadTenantOauth2Key"`
	SocialAuthAzureadTenantOauth2Resource string `yaml:"SocialAuthAzureadTenantOauth2Resource"`
	SocialAuthAzureadTenantOauth2Secret   string `yaml:"SocialAuthAzureadTenantOauth2Secret"`
	SocialAuthAzureadTenantOauth2TenantID string `yaml:"SocialAuthAzureadTenantOauth2TenantID"`
	SocialAuthGoogleOauth2Enable          string `yaml:"SocialAuthGoogleOauth2Enable"`
	SocialAuthGoogleOauth2Key             string `yaml:"SocialAuthGoogleOauth2Key"`
	SocialAuthGoogleOauth2Secret          string `yaml:"SocialAuthGoogleOauth2Secret"`
	SocialAuthOktaOauth2APIURL            string `yaml:"SocialAuthOktaOauth2APIURL"`
	SocialAuthOktaOauth2Enabled           string `yaml:"SocialAuthOktaOauth2Enabled"`
	SocialAuthOktaOauth2Key               string `yaml:"SocialAuthOktaOauth2Key"`
	SocialAuthOktaOauth2Secret            string `yaml:"SocialAuthOktaOauth2Secret"`
	StaticRoot                            string `yaml:"StaticRoot"`
	StaticURL                             string `yaml:"StaticURL"`
	TeamName                              string `yaml:"TeamName"`
	TestDatabaseName                      string `yaml:"TestDatabaseName"`
	TestRunner                            string `yaml:"TestRunner"`
	TimeZone                              string `yaml:"TimeZone"`
	TrackMigrations                       bool   `yaml:"TrackMigrations"`
	URLPrefix                             string `yaml:"UrlPrefix"`
	UseI18N                               bool   `yaml:"UseI18n"`
	UseL10N                               bool   `yaml:"UseL10n"`
	UseTZ                                 bool   `yaml:"UseTZ"`
	UUID                                  string `yaml:"UUID"`
	UwsgiEndpoint                         string `yaml:"UwsgiEndpoint"`
	UwsgiHost                             string `yaml:"UwsgiHost"`
	UwsgiMode                             string `yaml:"UwsgiMode"`
	UwsgiPass                             string `yaml:"UwsgiPass"`
	UwsgiPort                             string `yaml:"UwsgiPort"`
	Whitenoise                            bool   `yaml:"Whitenoise"`
	Wkhtmltopdf                           string `yaml:"Wkhtmltopdf"`
	DojoAdminUser                         string `yaml:"DojoAdminUser"`
} // yaml:"Settings"

// OptionalConfig values added to make developing and testing godojo easier
// AKA you should never really need to change these.
type OptionalConfig struct {
	HelpURL    string `yaml:"HelpURL"`
	ReleaseURL string `yaml:"ReleaseURL"`
	CloneURL   string `yaml:"CloneURL"`
	YarnGPG    string `yaml:"YarnGPG"`
	NodeURL    string `yaml:"NodeURL"`
	Embd       bool   `yaml:"Embd"`
	Key        string `yaml:"Key"`
	Tmpdir     string `yaml:Tmpdir`
}
