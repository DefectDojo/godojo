package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// prepInstaller takes a pointer to a gdjDefault struct and prepares for the
// installation by reading command-line args, making changes to the default
// based on command-line arguments and reading the environmental variables to
// override values from the config file
func prepInstaller(d *gdjDefault) {
	// Read the command-line arguments
	readArgs(d)

	// Setup logging

	// Handle default and dev installs
	if d.defInstall {
		// Set config options based on embedded default config
		defaultConfig(d)
	} else {
		// Read dojoConfig.yml file
		readConfigFile(d)
	}

	// Read in any environmental variables
	readEnvVars(&d.conf)

	// Write final install configuration to a file
	writeFinalConfig(d)

	// Initialize Redactatron
	d.initRedact()

	// Ensure installer has sufficient privileges
	checkUserPrivs(d)

	// Check that configured DB configuration is sane
	saneDBConfig(d)

	// Logging is setup, start using statusMsg and errorMsg functions for output
	d.traceMsg("Logging established, trace log begins here")
	d.sectionMsg("Starting the dojo install at " + time.Now().Format("Mon Jan 2, 2006 15:04:05 MST"))

}

// defaultConfig takes no arguements and setups a godojo installation to uses
// all the defaults in the config file
func defaultConfig(d *gdjDefault) {
	d.traceMsg("Inside of defaultConfig")
	// Temporarily write out the config file into current working directory
	writeDefaultConfig(d.cf, false)

	// Read the config file
	readConfigFile(d)

	// Clean-up the temporary config file
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine current working directory, exiting...")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	err = os.Remove(path + "/" + d.cf)
	if err != nil {
		// TODO Change me to error log
		fmt.Println("Unable to delete temporary config file")
		fmt.Printf("Error: %v\n", err)
		fmt.Println("File will remain for user to manually remove")
	}
}

// setDevDefaults has not been implemented
func setDevDefaults() {
	// TODO: Complete this option
	fmt.Println("")
	fmt.Println("Currently, this is not a supported option.")
	fmt.Println("Instead, please run ./godojo without any command-line options")
	fmt.Println("to create a default config file in the current working directory")
	fmt.Println("and edit that file as needed")
	fmt.Println("")
	fmt.Println("Alternatively, godojo can be run with \"-default\" to do an install")
	fmt.Println("using the default config options.")
	fmt.Println("")
	fmt.Println("Ask Matt nicely and he may knock this out for you. ;-)")
	fmt.Println("")
	os.Exit(1)
}

// readEnvVars reads the DefectDojo supported environmental variables and
// overrides any options set in the configuration file. These variables
// are used to supply either install-time configurations or provide values
// that are used in DefectDojo's settings.py configuration file
func readEnvVars(gdConf *dojoConfig) {
	// Env variables pulled from repo. Add newly supported env vars below and
	// to the switch statement below after the for that ranges over overrides
	// TODO: Add non-setting.py ENV vars like DD_SourcCommit
	dojoEnvs := map[string]bool{
		"DD_ADMIN_FIRST_NAME":                            true,
		"DD_ADMIN_LAST_NAME":                             true,
		"DD_ADMIN_MAIL":                                  true,
		"DD_ADMIN_PASSWORD":                              true,
		"DD_ADMINS":                                      true,
		"DD_ADMIN_USER":                                  true,
		"DD_ALLOWED_HOSTS":                               true,
		"DD_CELERY_BEAT_SCHEDULE_FILENAME":               true,
		"DD_CELERY_BROKER_HOST":                          true,
		"DD_CELERY_BROKER_PASSWORD":                      true,
		"DD_CELERY_BROKER_PATH":                          true,
		"DD_CELERY_BROKER_PORT":                          true,
		"DD_CELERY_BROKER_SCHEME":                        true,
		"DD_CELERY_BROKER_URL":                           true,
		"DD_CELERY_BROKER_USER":                          true,
		"DD_CELERY_LOG_LEVEL":                            true,
		"DD_CELERY_RESULT_BACKEND":                       true,
		"DD_CELERY_RESULT_EXPIRES":                       true,
		"DD_CELERY_TASK_IGNORE_RESULT":                   true,
		"DD_CELERY_TASK_SERIALIZER":                      true,
		"DD_CREDENTIAL_AES_256_KEY":                      true,
		"DD_CSRF_COOKIE_HTTPONLY":                        true,
		"DD_CSRF_COOKIE_SECURE":                          true,
		"DD_DATABASE_ENGINE":                             true,
		"DD_DATABASE_HOST":                               true,
		"DD_DATABASE_NAME":                               true,
		"DD_DATABASE_PASSWORD":                           true,
		"DD_DATABASE_PORT":                               true,
		"DD_DATABASE_TYPE":                               true,
		"DD_DATABASE_URL":                                true,
		"DD_DATABASE_USER":                               true,
		"DD_DATA_UPLOAD_MAX_MEMORY_SIZE":                 true,
		"DD_DEBUG":                                       true,
		"DD_DJANGO_ADMIN_ENABLED":                        true,
		"DD_EMAIL_URL":                                   true,
		"DD_ENV":                                         true,
		"DD_ENV_PATH":                                    true,
		"DD_FORCE_LOWERCASE_TAGS":                        true,
		"DD_HOST":                                        true,
		"DD_INITIALIZE":                                  true,
		"DD_LANG":                                        true,
		"DD_LANGUAGE_CODE":                               true,
		"DD_LOGIN_REDIRECT_URL":                          true,
		"DD_MAX_TAG_LENGTH":                              true,
		"DD_MEDIA_ROOT":                                  true,
		"DD_MEDIA_URL":                                   true,
		"DD_PORT":                                        true,
		"DD_PORT_SCAN_CONTACT_EMAIL":                     true,
		"DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST":          true,
		"DD_PORT_SCAN_RESULT_EMAIL_FROM":                 true,
		"DD_PORT_SCAN_SOURCE_IP":                         true,
		"DD_ROOT":                                        true,
		"DD_SECRET_KEY":                                  true,
		"DD_SECURE_BROWSER_XSS_FILTER":                   true,
		"DD_SECURE_CONTENT_TYPE_NOSNIFF":                 true,
		"DD_SECURE_HSTS_INCLUDE_SUBDOMAINS":              true,
		"DD_SECURE_HSTS_SECONDS":                         true,
		"DD_SECURE_PROXY_SSL_HEADER":                     true,
		"DD_SECURE_SSL_REDIRECT":                         true,
		"DD_SESSION_COOKIE_HTTPONLY":                     true,
		"DD_SESSION_COOKIE_SECURE":                       true,
		"DD_SITE_ID":                                     true,
		"DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_ENABLED":   true,
		"DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_KEY":       true,
		"DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_RESOURCE":  true,
		"DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_SECRET":    true,
		"DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_TENANT_ID": true,
		"DD_SOCIAL_AUTH_GOOGLE_OAUTH2_ENABLE":            true,
		"DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY":               true,
		"DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET":            true,
		"DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL":             true,
		"DD_SOCIAL_AUTH_OKTA_OAUTH2_ENABLED":             true,
		"DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY":                 true,
		"DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET":              true,
		"DD_STATIC_ROOT":                                 true,
		"DD_STATIC_URL":                                  true,
		"DD_TEAM_NAME":                                   true,
		"DD_TEST_DATABASE_NAME":                          true,
		"DD_TEST_RUNNER":                                 true,
		"DD_TIME_ZONE":                                   true,
		"DD_TRACK_MIGRATIONS":                            true,
		"DD_URL_PREFIX":                                  true,
		"DD_USE_I18N":                                    true,
		"DD_USE_L10N":                                    true,
		"DD_USE_TZ":                                      true,
		"DD_UUID":                                        true,
		"DD_UWSGI_ENDPOINT":                              true,
		"DD_UWSGI_HOST":                                  true,
		"DD_UWSGI_MODE":                                  true,
		"DD_UWSGI_PASS":                                  true,
		"DD_UWSGI_PORT":                                  true,
		"DD_WHITENOISE":                                  true,
		"DD_WKHTMLTOPDF":                                 true,
		"DOJO_ADMIN_USER":                                true,
	} // End of dojoEnvs declaration

	match := false
	overrides := make(map[string]string)
	for _, e := range os.Environ() {
		// Pull out each env var into a slice
		env := strings.SplitN(e, "=", 2)

		// Check if the name of the env var matches the supported map
		if _, ok := dojoEnvs[env[0]]; ok {
			// If matched, add to overrides map
			overrides[env[0]] = env[1]
			match = true
		}
	}

	// Return early if no env variables are matched
	if !match {
		return
	}

	// Override config values if we found matching Env vars
	for k, v := range overrides {
		// Set DojoConfig struct values from Env variables to override config values
		// Have to do this as a switch statement as there's no sanity to DefectDojo env var naming
		switch k {
		case "DD_ADMIN_FIRST_NAME":
			gdConf.Settings.AdminFirstName = v
		case "DD_ADMIN_LAST_NAME":
			gdConf.Settings.AdminLastName = v
		case "DD_ADMIN_MAIL":
			gdConf.Settings.AdminMail = v
		case "DD_ADMIN_PASSWORD":
			gdConf.Settings.AdminPassword = v
		case "DD_ADMINS":
			gdConf.Settings.Admins = v
		case "DD_ADMIN_USER":
			gdConf.Settings.AdminUser = v
		case "DD_ALLOWED_HOSTS":
			gdConf.Settings.AllowedHosts = v
		case "DD_CELERY_BEAT_SCHEDULE_FILENAME":
			gdConf.Settings.CeleryBeatScheduleFilename = v
		case "DD_CELERY_BROKER_HOST":
			gdConf.Settings.CeleryBrokerHost = v
		case "DD_CELERY_BROKER_PASSWORD":
			gdConf.Settings.CeleryBrokerPassword = v
		case "DD_CELERY_BROKER_PATH":
			gdConf.Settings.CeleryBrokerPath = v
		case "DD_CELERY_BROKER_PORT":
			port := convInt(v, "DD_CELERY_BROKER_PORT provided via environmental variable isn't a valid port number")
			intLessThan(port, 65535, "DD_CELERY_BROKER_PORT provided via environmental variable is too large")
			gdConf.Settings.CeleryBrokerPort = port
		case "DD_CELERY_BROKER_SCHEME":
			gdConf.Settings.CeleryBrokerScheme = v
		case "DD_CELERY_BROKER_URL":
			gdConf.Settings.CeleryBrokerURL = v
		case "DD_CELERY_BROKER_USER":
			gdConf.Settings.CeleryBrokerUser = v
		case "DD_CELERY_LOG_LEVEL":
			gdConf.Settings.CeleryLogLevel = v
		case "DD_CELERY_RESULT_BACKEND":
			gdConf.Settings.CeleryResultBackend = v
		case "DD_CELERY_RESULT_EXPIRES":
			gdConf.Settings.CeleryResultExpires = convInt(v, "DD_CELERY_RESULT_EXPIRES provided via environmental variable isn't a valid number")
		case "DD_CELERY_TASK_IGNORE_RESULT":
			gdConf.Settings.CeleryTaskIgnoreResult = convBool(v, "DD_CELERY_TASK_IGNORE_RESULT environmental variable was not a boolean.")
		case "DD_CELERY_TASK_SERIALIZER":
			gdConf.Settings.CeleryTaskSerializer = v
		case "DD_CREDENTIAL_AES_256_KEY":
			gdConf.Settings.CredentialAES256Key = v
		case "DD_CSRF_COOKIE_HTTPONLY":
			gdConf.Settings.CSRFCookieHTTPOnly = convBool(v, "DD_CSRF_COOKIE_HTTPONLY environmental variable was not a boolean.")
		case "DD_CSRF_COOKIE_SECURE":
			gdConf.Settings.CSRFCookieSecure = convBool(v, "DD_CSRF_COOKIE_SECURE environmental variable was not a boolean.")
		case "DD_DATABASE_ENGINE":
			gdConf.Settings.DatabaseEngine = v
		case "DD_DATABASE_HOST":
			gdConf.Settings.DatabaseHost = v
		case "DD_DATABASE_NAME":
			gdConf.Settings.DatabaseName = v
		case "DD_DATABASE_PASSWORD":
			gdConf.Settings.DatabasePassword = v
		case "DD_DATABASE_PORT":
			gdConf.Settings.DatabasePort = v
		case "DD_DATABASE_TYPE":
			gdConf.Settings.DatabaseType = v
		case "DD_DATABASE_URL":
			gdConf.Settings.DatabaseURL = v
		case "DD_DATABASE_USER":
			gdConf.Settings.DatabaseUser = v
		case "DD_DATA_UPLOAD_MAX_MEMORY_SIZE":
			gdConf.Settings.DataUploadMaxMemorySize = convInt(v, "DD_DATA_UPLOAD_MAX_MEMORY_SIZE provided via environmental variable isn't a valid number")
		case "DD_DEBUG":
			gdConf.Settings.Debug = convBool(v, "DD_DEBUG environmental variable was not a boolean.")
		case "DD_DJANGO_ADMIN_ENABLED":
			gdConf.Settings.DjangoAdminEnabled = convBool(v, "DD_DJANGO_ADMIN_ENABLED environmental variable was not a boolean.")
		case "DD_EMAIL_URL":
			gdConf.Settings.EmailURL = v
		case "DD_ENV":
			gdConf.Settings.Env = v
		case "DD_ENV_PATH":
			gdConf.Settings.EnvPath = v
		case "DD_FORCE_LOWERCASE_TAGS":
			gdConf.Settings.ForceLowercaseTags = convBool(v, "DD_FORCE_LOWERCASE_TAGS environmental variable was not a boolean.")
		case "DD_HOST":
			gdConf.Settings.Host = v
		case "DD_INITIALIZE":
			gdConf.Settings.Initialize = v
		case "DD_LANG":
			gdConf.Settings.Lang = v
		case "DD_LANGUAGE_CODE":
			gdConf.Settings.LanguageCode = v
		case "DD_LOGIN_REDIRECT_URL":
			gdConf.Settings.LoginRedirectURL = v
		case "DD_MAX_TAG_LENGTH":
			// TODO: Look up maximum tag length in data model and check for that too
			gdConf.Settings.MaxTagLength = convInt(v, "DD_MAX_TAG_LENGTH provided via environmental variable isn't a valid number")
		case "DD_MEDIA_ROOT":
			gdConf.Settings.MediaRoot = v
		case "DD_MEDIA_URL":
			gdConf.Settings.MediaURL = v
		case "DD_PORT":
			gdConf.Settings.Port = v
		case "DD_PORT_SCAN_CONTACT_EMAIL":
			gdConf.Settings.PortScanContactEmail = v
		case "DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST":
			gdConf.Settings.PortScanExternalUnitEmailList = v
		case "DD_PORT_SCAN_RESULT_EMAIL_FROM":
			gdConf.Settings.PortScanResultEmailFrom = v
		case "DD_PORT_SCAN_SOURCE_IP":
			gdConf.Settings.PortScanSourceIP = v
		case "DD_ROOT":
			gdConf.Settings.Root = v
		case "DD_SECRET_KEY":
			gdConf.Settings.SecretKey = v
		case "DD_SECURE_BROWSER_XSS_FILTER":
			gdConf.Settings.SecureBrowserXSSFilter = convBool(v, "DD_SECURE_BROWSER_XSS_FILTER environmental variable was not a boolean.")
		case "DD_SECURE_CONTENT_TYPE_NOSNIFF":
			gdConf.Settings.SecureContentTypeNosniff = v
		case "DD_SECURE_HSTS_INCLUDE_SUBDOMAINS":
			gdConf.Settings.SecureHSTSIncludeSubdomains = convBool(v, "DD_SECURE_HSTS_INCLUDE_SUBDOMAINS environmental variable was not a boolean.")
		case "DD_SECURE_HSTS_SECONDS":
			gdConf.Settings.SecureHSTSSeconds = convInt(v, "DD_SECURE_HSTS_SECONDS provided via environmental variable isn't a valid number")
		case "DD_SECURE_PROXY_SSL_HEADER":
			gdConf.Settings.SecureProxySSLHeader = convBool(v, "DD_SECURE_PROXY_SSL_HEADER environmental variable was not a boolean.")
		case "DD_SECURE_SSL_REDIRECT":
			gdConf.Settings.SecureSSLRedirect = convBool(v, "DD_SECURE_SSL_REDIRECT environmental variable was not a boolean.")
		case "DD_SESSION_COOKIE_HTTPONLY":
			gdConf.Settings.SessionCookieHTTPOnly = convBool(v, "DD_SESSION_COOKIE_HTTPONLY environmental variable was not a boolean.")
		case "DD_SESSION_COOKIE_SECURE":
			gdConf.Settings.SessionCookieSecure = convBool(v, "DD_SESSION_COOKIE_SECURE environmental variable was not a boolean.")
		case "DD_SITE_ID":
			gdConf.Settings.SiteID = convInt(v, "DD_SITE_ID provided via environmental variable isn't a valid number")
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_ENABLED":
			gdConf.Settings.SocialAuthAzureadTenantOauth2Enabled = v
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_KEY":
			gdConf.Settings.SocialAuthAzureadTenantOauth2Key = v
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_RESOURCE":
			gdConf.Settings.SocialAuthAzureadTenantOauth2Resource = v
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_SECRET":
			gdConf.Settings.SocialAuthAzureadTenantOauth2Secret = v
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_TENANT_ID":
			gdConf.Settings.SocialAuthAzureadTenantOauth2TenantID = v
		case "DD_SOCIAL_AUTH_GOOGLE_OAUTH2_ENABLE":
			gdConf.Settings.SocialAuthGoogleOauth2Enable = v
		case "DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY":
			gdConf.Settings.SocialAuthGoogleOauth2Key = v
		case "DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET":
			gdConf.Settings.SocialAuthGoogleOauth2Secret = v
		case "DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL":
			gdConf.Settings.SocialAuthOktaOauth2APIURL = v
		case "DD_SOCIAL_AUTH_OKTA_OAUTH2_ENABLED":
			gdConf.Settings.SocialAuthOktaOauth2Enabled = v
		case "DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY":
			gdConf.Settings.SocialAuthOktaOauth2Key = v
		case "DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET":
			gdConf.Settings.SocialAuthOktaOauth2Secret = v
		case "DD_STATIC_ROOT":
			gdConf.Settings.StaticRoot = v
		case "DD_STATIC_URL":
			gdConf.Settings.StaticURL = v
		case "DD_TEAM_NAME":
			gdConf.Settings.TeamName = v
		case "DD_TEST_DATABASE_NAME":
			gdConf.Settings.TestDatabaseName = v
		case "DD_TEST_RUNNER":
			gdConf.Settings.TestRunner = v
		case "DD_TIME_ZONE":
			gdConf.Settings.TimeZone = v
		case "DD_TRACK_MIGRATIONS":
			gdConf.Settings.TrackMigrations = convBool(v, "DD_TRACK_MIGRATIONS environmental variable was not a boolean.")
		case "DD_URL_PREFIX":
			gdConf.Settings.URLPrefix = v
		case "DD_USE_I18N":
			gdConf.Settings.UseI18N = convBool(v, "DD_USE_I18N environmental variable was not a boolean.")
		case "DD_USE_L10N":
			gdConf.Settings.UseL10N = convBool(v, "DD_USE_L10N environmental variable was not a boolean.")
		case "DD_USE_TZ":
			gdConf.Settings.UseTZ = convBool(v, "DD_USE_TZ environmental variable was not a boolean.")
		case "DD_UUID":
			gdConf.Settings.UUID = v
		case "DD_UWSGI_ENDPOINT":
			gdConf.Settings.UwsgiEndpoint = v
		case "DD_UWSGI_HOST":
			gdConf.Settings.UwsgiHost = v
		case "DD_UWSGI_MODE":
			gdConf.Settings.UwsgiMode = v
		case "DD_UWSGI_PASS":
			gdConf.Settings.UwsgiPass = v
		case "DD_UWSGI_PORT":
			gdConf.Settings.UwsgiPort = v
		case "DD_WHITENOISE":
			gdConf.Settings.Whitenoise = convBool(v, "DD_WHITENOISE environmental variable was not a boolean.")
		case "DD_WKHTMLTOPDF":
			gdConf.Settings.Wkhtmltopdf = v
		case "DOJO_ADMIN_USER":
			gdConf.Settings.DojoAdminUser = v
			// TODO: Deprecate me
		}
	}

}

func convInt(i string, s string) int {
	convI, err := strconv.Atoi(i)
	if err != nil {
		fmt.Println("ERROR:")
		fmt.Printf("  %s\n", s)
		fmt.Printf("  Error was: %v\n", err)
		os.Exit(1)
	}
	return convI
}

func intLessThan(i int, max int, s string) {
	if i > max {
		fmt.Println("ERROR:")
		fmt.Printf("  %s\n", s)
		os.Exit(1)
	}
}

func convBool(b string, s string) bool {
	res, err := strconv.ParseBool(b)
	if err != nil {
		fmt.Println("ERROR:")
		fmt.Printf("  %s\n", s)
		fmt.Println("  Valid values are 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.")
		fmt.Printf("  Error was: %v\n", err)
		os.Exit(1)
	}
	return res
}

// checkUserPrivs takes a pointer to gdjDefault struct and verifies that the
// user running godojo has sufficient privileges to complete the install and
// exits with a 1 if privileges are lacking
func checkUserPrivs(d *gdjDefault) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if usr.Uid != "0" && !d.conf.Options.UsrInst {
		fmt.Println("")
		fmt.Println("##############################################################################")
		fmt.Println("  ERROR: This program must be run as root or with sudo\n  Please correct and run installer again")
		fmt.Println("##############################################################################")
		fmt.Println("")
		os.Exit(1)
	}
}
