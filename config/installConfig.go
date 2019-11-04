package config

// DojoConfig - "mother" struct to hold all the config options
type DojoConfig struct {
	Install  InstallConfig
	Settings SettingsConfig
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
	Set           string         // The install set or type: Single Server, Dev, Stand-alone
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
	// Configs for settings.py
	Debug       bool // Run Dojo in debug mode, default false
	Login       LoginSt
	Django      DjangoSt
	Session     SessionSt
	CSRF        CSRFSt
	Secure      SecureSt
	Time        TimeSt
	Lang        string
	Wkhtmltopdf string
	Team        TeamSt
	Admins      string // May be redundant
	Port        PortSt
	Whitenoise  bool
	Track       TrackSt
	Test        TestSt
	URL         URLSt
	Root        string
	Language    LanguageSt
	Site        SiteSt
	Use         UseSt
	Media       MediaStatic
	Static      MediaStatic
	Celery      CelerySt
	Force       ForceSt
	Max         MaxSt
	Database    DatabaseSt
	Secret      SecretSt
	Credential  CredentialSt
	Data        DataSt
	Social      SocialSt
	Allowed     AllowedSt
	Email       EmailUSt // WARNING - "U" Added to make definition unique
}

// LoginSt - struct for DD_LOGIN_REDIRECT_URL
type LoginSt struct {
	Redirect RedirectSt
}

// RedirectSt - struct for DD_LOGIN_REDIRECT_URL
type RedirectSt struct {
	URL string
}

// DjangoSt - struct for DD_DJANGO_ADMIN_ENABLED
type DjangoSt struct {
	Admin AdminSt
}

// AdminSt - struct for DD_DJANGO_ADMIN_ENABLED
type AdminSt struct {
	Enabled bool
}

// SessionSt - struct for DD_SESSION_COOKIE_HTTPONLY
type SessionSt struct {
	Cookie CookieSt
	Secure bool
}

// CookieSt - struct for DD_SESSION_COOKIE_HTTPONLY
type CookieSt struct {
	HTTPOnly bool
}

// CSRFSt - struct for DD_CSFR_COOKIE_HTTPONLY and DD_CSRF_COOKIE_SECURE
type CSRFSt struct {
	Cookie CookieSt // Reuse struct above for DD_SESSION_COOKIE_HTTPONLY
	Secure bool
}

// SecureSt - struct for DD_SECURE_SSL_REDIRECT, DD_SECURE_HSTS_INCLUDE_SUBDOMAINS, DD_SECURE_HSTS_SECONDS
// DD_SECURE_BROWSER_XSS_FILTER, and DD_SECURE_PROXY_SSL_HEADER
type SecureSt struct {
	SSL     SSLSt
	HSTS    HSTSSt
	Browser BrowserSt
	Proxy   ProxySt
}

// SSLSt - struct for DD_SECURE_SSL_REDIRECT, DD_SECURE_HSTS_INCLUDE_SUBDOMAINS, DD_SECURE_HSTS_SECONDS
// DD_SECURE_BROWSER_XSS_FILTER, and DD_SECURE_PROXY_SSL_HEADER
type SSLSt struct {
	Redirect bool
}

// HSTSSt - struct for DD_SECURE_SSL_REDIRECT, DD_SECURE_HSTS_INCLUDE_SUBDOMAINS, DD_SECURE_HSTS_SECONDS
// DD_SECURE_BROWSER_XSS_FILTER, and DD_SECURE_PROXY_SSL_HEADER
type HSTSSt struct {
	Include IncludeSt
	Seconds uint32
}

// IncludeSt - struct for DD_SECURE_SSL_REDIRECT, DD_SECURE_HSTS_INCLUDE_SUBDOMAINS, DD_SECURE_HSTS_SECONDS
// DD_SECURE_BROWSER_XSS_FILTER, and DD_SECURE_PROXY_SSL_HEADER
type IncludeSt struct {
	Subdomains bool
}

// BrowserSt - struct for DD_SECURE_SSL_REDIRECT, DD_SECURE_HSTS_INCLUDE_SUBDOMAINS, DD_SECURE_HSTS_SECONDS
// DD_SECURE_BROWSER_XSS_FILTER, and DD_SECURE_PROXY_SSL_HEADER
type BrowserSt struct {
	XSS XSSSt
}

// XSSSt - struct for DD_SECURE_SSL_REDIRECT, DD_SECURE_HSTS_INCLUDE_SUBDOMAINS, DD_SECURE_HSTS_SECONDS
// DD_SECURE_BROWSER_XSS_FILTER, and DD_SECURE_PROXY_SSL_HEADER
type XSSSt struct {
	Filter bool
}

// ProxySt - struct for DD_SECURE_SSL_REDIRECT, DD_SECURE_HSTS_INCLUDE_SUBDOMAINS, DD_SECURE_HSTS_SECONDS
// DD_SECURE_BROWSER_XSS_FILTER, and DD_SECURE_PROXY_SSL_HEADER
type ProxySt struct {
	PSSL PSSLSt // WARNING - "P" Added to make definition unique
}

// PSSLSt - struct for DD_SECURE_SSL_REDIRECT, DD_SECURE_HSTS_INCLUDE_SUBDOMAINS, DD_SECURE_HSTS_SECONDS
// DD_SECURE_BROWSER_XSS_FILTER, and DD_SECURE_PROXY_SSL_HEADER
type PSSLSt struct {
	Header bool
}

// TimeSt - struct for DD_TIME_ZONE
type TimeSt struct {
	Zone string
}

// TeamSt - struct for DD_TEAM_NAME
type TeamSt struct {
	Name string
}

// PortSt - struct for DD_PORT_SCAN_CONTACT_EMAIL, DD_PORT_SCAN_RESULT_EMAIL_FROM, DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST,
// and DD_PORT_SCAN_SOURCE_IP
type PortSt struct {
	Scan ScanSt
}

// ScanSt - struct for DD_PORT_SCAN_CONTACT_EMAIL, DD_PORT_SCAN_RESULT_EMAIL_FROM, DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST,
// and DD_PORT_SCAN_SOURCE_IP
type ScanSt struct {
	Contact  ContactSt
	Result   ResultSt
	External ExternalSt
	Source   SourceSt
}

// ContactSt - struct for DD_PORT_SCAN_CONTACT_EMAIL, DD_PORT_SCAN_RESULT_EMAIL_FROM, DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST,
// and DD_PORT_SCAN_SOURCE_IP
type ContactSt struct {
	Email string
}

// ResultSt - struct for DD_PORT_SCAN_CONTACT_EMAIL, DD_PORT_SCAN_RESULT_EMAIL_FROM, DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST,
// and DD_PORT_SCAN_SOURCE_IP
type ResultSt struct {
	Email EmailSt
}

// EmailSt - struct for DD_PORT_SCAN_CONTACT_EMAIL, DD_PORT_SCAN_RESULT_EMAIL_FROM, DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST,
// and DD_PORT_SCAN_SOURCE_IP
type EmailSt struct {
	From string
}

// ExternalSt - struct for DD_PORT_SCAN_CONTACT_EMAIL, DD_PORT_SCAN_RESULT_EMAIL_FROM, DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST,
// and DD_PORT_SCAN_SOURCE_IP
type ExternalSt struct {
	Unit UnitSt
}

// UnitSt - struct for DD_PORT_SCAN_CONTACT_EMAIL, DD_PORT_SCAN_RESULT_EMAIL_FROM, DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST,
// and DD_PORT_SCAN_SOURCE_IP
type UnitSt struct {
	Email UEmailSt // WARNING - "U" Added to make definition unique
}

// UEmailSt - struct for DD_PORT_SCAN_CONTACT_EMAIL, DD_PORT_SCAN_RESULT_EMAIL_FROM, DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST,
// and DD_PORT_SCAN_SOURCE_IP
type UEmailSt struct {
	List string
}

// SourceSt - struct for DD_PORT_SCAN_CONTACT_EMAIL, DD_PORT_SCAN_RESULT_EMAIL_FROM, DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST,
// and DD_PORT_SCAN_SOURCE_IP
type SourceSt struct {
	IP string
}

// TrackSt - struct for DD_TRACK_MIGRATIONS
type TrackSt struct {
	Migrations bool
}

// TestSt - struct for DD_TEST_RUNNER and DD_TEST_DATABASE_NAME
type TestSt struct {
	Runner   string
	Database TDatabaseSt
}

// TDatabaseSt - struct for DD_TEST_RUNNER and DD_TEST_DATABASE_NAME
type TDatabaseSt struct {
	Name string
}

// URLSt - struct for DD_URL_PREFIX
type URLSt struct {
	Prefix string
}

// LanguageSt - struct for DD_LANGUAGE_CODE
type LanguageSt struct {
	Code string
}

// SiteSt - struct for DD_SITE_ID
type SiteSt struct {
	ID int
}

// UseSt - struct for DD_USE_I18N, DD_USE_L10N, and DD_USE_TZ
type UseSt struct {
	I18N bool
	L10N bool
	TZ   bool
}

// MediaStatic - stuct for DD_MEDIA_ROOT and DD_MEDIA_URL plus DD_STATIC_ROOT and DD_STATIC_URL
type MediaStatic struct { // Used for Media and Static config items since they are the same
	Root string
	URL  string
}

// CelerySt - struct for DD_CELERY_BROKER_URL, DD_CELERY_BROKER_SCHEME, DD_CELERY_BROKER_USER, DD_CELERY_BROKER_PASSWORD,
// DD_CELERY_BROKER_HOST, DD_CELERY_BROKER_PORT, DD_CELERY_BROKER_PATH, DD_CELERY_TASK_IGNORE_RESULT,
// DD_CELERY_RESULT_BACKEND, DD_CELERY_RESULT_EXPIRES, DD_CELERY_BEAT_SCHEDULE_FILENAME
type CelerySt struct {
	Broker BrokerSt
	Task   TaskSt
	Result CResultSt // WARNING - "C" Added to make definition unique
	Beat   BeatSt
}

// BrokerSt - struct for DD_CELERY_BROKER_URL, DD_CELERY_BROKER_SCHEME, DD_CELERY_BROKER_USER, DD_CELERY_BROKER_PASSWORD,
// DD_CELERY_BROKER_HOST, DD_CELERY_BROKER_PORT, DD_CELERY_BROKER_PATH, DD_CELERY_TASK_IGNORE_RESULT,
// DD_CELERY_RESULT_BACKEND, DD_CELERY_RESULT_EXPIRES, DD_CELERY_BEAT_SCHEDULE_FILENAME
type BrokerSt struct {
	URL      string
	Scheme   string
	User     string
	Password string
	Host     string
	Port     int
	Path     string
}

// TaskSt - struct for DD_CELERY_BROKER_URL, DD_CELERY_BROKER_SCHEME, DD_CELERY_BROKER_USER, DD_CELERY_BROKER_PASSWORD,
// DD_CELERY_BROKER_HOST, DD_CELERY_BROKER_PORT, DD_CELERY_BROKER_PATH, DD_CELERY_TASK_IGNORE_RESULT,
// DD_CELERY_RESULT_BACKEND, DD_CELERY_RESULT_EXPIRES, DD_CELERY_BEAT_SCHEDULE_FILENAME
type TaskSt struct {
	Ignore     IgnoreSt
	Serializer string
}

// IgnoreSt - struct for DD_CELERY_BROKER_URL, DD_CELERY_BROKER_SCHEME, DD_CELERY_BROKER_USER, DD_CELERY_BROKER_PASSWORD,
// DD_CELERY_BROKER_HOST, DD_CELERY_BROKER_PORT, DD_CELERY_BROKER_PATH, DD_CELERY_TASK_IGNORE_RESULT,
// DD_CELERY_RESULT_BACKEND, DD_CELERY_RESULT_EXPIRES, DD_CELERY_BEAT_SCHEDULE_FILENAME
type IgnoreSt struct {
	Result bool
}

// CResultSt - struct for DD_CELERY_BROKER_URL, DD_CELERY_BROKER_SCHEME, DD_CELERY_BROKER_USER, DD_CELERY_BROKER_PASSWORD,
// DD_CELERY_BROKER_HOST, DD_CELERY_BROKER_PORT, DD_CELERY_BROKER_PATH, DD_CELERY_TASK_IGNORE_RESULT,
// DD_CELERY_RESULT_BACKEND, DD_CELERY_RESULT_EXPIRES, DD_CELERY_BEAT_SCHEDULE_FILENAME
type CResultSt struct {
	Backend string
	Expires int
}

// BeatSt - struct for DD_CELERY_BROKER_URL, DD_CELERY_BROKER_SCHEME, DD_CELERY_BROKER_USER, DD_CELERY_BROKER_PASSWORD,
// DD_CELERY_BROKER_HOST, DD_CELERY_BROKER_PORT, DD_CELERY_BROKER_PATH, DD_CELERY_TASK_IGNORE_RESULT,
// DD_CELERY_RESULT_BACKEND, DD_CELERY_RESULT_EXPIRES, DD_CELERY_BEAT_SCHEDULE_FILENAME
type BeatSt struct {
	Schedule ScheduleSt
}

// ScheduleSt - struct for DD_CELERY_BROKER_URL, DD_CELERY_BROKER_SCHEME, DD_CELERY_BROKER_USER, DD_CELERY_BROKER_PASSWORD,
// DD_CELERY_BROKER_HOST, DD_CELERY_BROKER_PORT, DD_CELERY_BROKER_PATH, DD_CELERY_TASK_IGNORE_RESULT,
// DD_CELERY_RESULT_BACKEND, DD_CELERY_RESULT_EXPIRES, DD_CELERY_BEAT_SCHEDULE_FILENAME
type ScheduleSt struct {
	Filename string
}

// ForceSt - struct for DD_FORCE_LOWERCASE_TAGS
type ForceSt struct {
	Lowercase LowercaseSt
}

// LowercaseSt - struct for DD_FORCE_LOWERCASE_TAGS
type LowercaseSt struct {
	tags bool
}

// MaxSt - struct for DD_MAX_TAG_LENGTH
type MaxSt struct {
	Tag TagSt
}

// TagSt - struct for DD_MAX_TAG_LENGTH
type TagSt struct {
	Length int
}

// DatabaseSt - struct for DD_DATABSE_ENGINE, DD_DATABSE_HOST, DD_DATABSE_NAME, DD_DATABSE_PASSWORD, DD_DATABSE_PORT,
// DD_DATABSE_USER
type DatabaseSt struct {
	Engine   string
	Host     string
	Name     string
	Password string
	Port     string
	User     string
}

// SecretSt - struct for DD_SECRET_KEY
type SecretSt struct {
	Key string
}

// CredentialSt - struct for DD_CREDENTIAL_AES_256_KEY
type CredentialSt struct {
	AES AESSt
}

// AESSt - struct for DD_CREDENTIAL_AES_256_KEY
type AESSt struct {
	B256 B256St
}

// B256St - struct for DD_CREDENTIAL_AES_256_KEY
type B256St struct {
	Key string
}

// DataSt - struct for DD_DATA_UPLOAD_MAX_MEMORY_SIZE
type DataSt struct {
	Upload UploadSt
}

// UploadSt - struct for DD_DATA_UPLOAD_MAX_MEMORY_SIZE
type UploadSt struct {
	Max UMaxSt // WARNING - "U" Added to make definition unique
}

// UMaxSt - struct for DD_DATA_UPLOAD_MAX_MEMORY_SIZE
type UMaxSt struct {
	Memory MemorySt
}

// MemorySt - struct for DD_DATA_UPLOAD_MAX_MEMORY_SIZE
type MemorySt struct {
	Size uint64
}

// SocialSt - struct for DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY, DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY,
// DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL
type SocialSt struct {
	Auth AuthSt
}

// AuthSt - struct for DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY, DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY,
// DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL
type AuthSt struct {
	Google GoogleSt
	Okta   OktaSt
}

// GoogleSt - struct for DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY, DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY,
// DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL
type GoogleSt struct {
	OAUTH2 OAUTH2St
}

// OAUTH2St - struct for DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY, DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY,
// DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL
type OAUTH2St struct { // Overloaded to hold data for both Google and Okta - Google only uses Key and Secret
	Key    string
	Secret string
	API    APISt
}

// APISt - struct for DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY, DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY,
// DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL
type APISt struct {
	URL string
}

// OktaSt - struct for DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY, DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY,
// DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET, DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL
type OktaSt struct {
	OAUTH2 OAUTH2St // Struct shared with Google above ^
}

// AllowedSt - struct for DD_ALLOWED_HOSTS
type AllowedSt struct {
	Hosts string
}

// EmailUSt - struct for DD_EMAIL_URL
type EmailUSt struct {
	URL string
}
