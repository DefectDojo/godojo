package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func defaultConfig() {
	traceMsg("Inside of defaultConfig")
	// Temporarily write out the config file into current working directory
	createDefaultConfig(cf, false)

	// Read the config file
	readConfigFile()

	// Clean-up the temporary config file
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine current working directory, exiting...")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	err = os.Remove(path + "/" + cf)
	if err != nil {
		// TODO Change me to error log
		fmt.Println("Unable to delete temporary config file")
		fmt.Printf("Error: %v\n", err)
		fmt.Println("File will remain for user to manually remove")
	}
}

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

// readConfigFile reads the yaml configuration file for godojo
// to determine runtime configuration.  The file is dojoConfig.yml
// and is expected to be in the same directory as the godojo binary
// It returns nohing but will exit early with a exit code of 1
// if there are errors reading the file or unmarshialling into a struct
func readConfigFile() {
	// Setup viper config
	viper.AddConfigPath(".")
	viper.SetConfigName("dojoConfig")
	viper.SetConfigType("yml")

	// Read the default config file dojoConfig.yml
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("")
		fmt.Println("Unable to read the godojo config file (dojoConfig.yml), exiting install")
		fmt.Printf("Error was: %v\n", err)
		os.Exit(1)
	}

	// Marshall the config values into the DojoConfig struct
	err = viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println("")
		fmt.Println("Unable to set the config values based on config file and ENV variables, exiting install")
		fmt.Printf("Error was: %v\n", err)
		os.Exit(1)
	}
}

// readEnvVars reads the DefectDojo supported environmental variables and
// overrides any options set in the configuration file. These variables
// are used to supply either install-time configurations or provide values
// that are used in DefectDojo's settings.py configuration file
func readEnvVars() {
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
			conf.Settings.AdminFirstName = v
		case "DD_ADMIN_LAST_NAME":
			conf.Settings.AdminLastName = v
		case "DD_ADMIN_MAIL":
			conf.Settings.AdminMail = v
		case "DD_ADMIN_PASSWORD":
			conf.Settings.AdminPassword = v
		case "DD_ADMINS":
			conf.Settings.Admins = v
		case "DD_ADMIN_USER":
			conf.Settings.AdminUser = v
		case "DD_ALLOWED_HOSTS":
			conf.Settings.AllowedHosts = v
		case "DD_CELERY_BEAT_SCHEDULE_FILENAME":
			conf.Settings.CeleryBeatScheduleFilename = v
		case "DD_CELERY_BROKER_HOST":
			conf.Settings.CeleryBrokerHost = v
		case "DD_CELERY_BROKER_PASSWORD":
			conf.Settings.CeleryBrokerPassword = v
		case "DD_CELERY_BROKER_PATH":
			conf.Settings.CeleryBrokerPath = v
		case "DD_CELERY_BROKER_PORT":
			port := convInt(v, "DD_CELERY_BROKER_PORT provided via environmental variable isn't a valid port number")
			intLessThan(port, 65535, "DD_CELERY_BROKER_PORT provided via environmental variable is too large")
			conf.Settings.CeleryBrokerPort = port
		case "DD_CELERY_BROKER_SCHEME":
			conf.Settings.CeleryBrokerScheme = v
		case "DD_CELERY_BROKER_URL":
			conf.Settings.CeleryBrokerURL = v
		case "DD_CELERY_BROKER_USER":
			conf.Settings.CeleryBrokerUser = v
		case "DD_CELERY_LOG_LEVEL":
			conf.Settings.CeleryLogLevel = v
		case "DD_CELERY_RESULT_BACKEND":
			conf.Settings.CeleryResultBackend = v
		case "DD_CELERY_RESULT_EXPIRES":
			conf.Settings.CeleryResultExpires = convInt(v, "DD_CELERY_RESULT_EXPIRES provided via environmental variable isn't a valid number")
		case "DD_CELERY_TASK_IGNORE_RESULT":
			conf.Settings.CeleryTaskIgnoreResult = convBool(v, "DD_CELERY_TASK_IGNORE_RESULT environmental variable was not a boolean.")
		case "DD_CELERY_TASK_SERIALIZER":
			conf.Settings.CeleryTaskSerializer = v
		case "DD_CREDENTIAL_AES_256_KEY":
			conf.Settings.CredentialAES256Key = v
		case "DD_CSRF_COOKIE_HTTPONLY":
			conf.Settings.CSRFCookieHTTPOnly = convBool(v, "DD_CSRF_COOKIE_HTTPONLY environmental variable was not a boolean.")
		case "DD_CSRF_COOKIE_SECURE":
			conf.Settings.CSRFCookieSecure = convBool(v, "DD_CSRF_COOKIE_SECURE environmental variable was not a boolean.")
		case "DD_DATABASE_ENGINE":
			conf.Settings.DatabaseEngine = v
		case "DD_DATABASE_HOST":
			conf.Settings.DatabaseHost = v
		case "DD_DATABASE_NAME":
			conf.Settings.DatabaseName = v
		case "DD_DATABASE_PASSWORD":
			conf.Settings.DatabasePassword = v
		case "DD_DATABASE_PORT":
			conf.Settings.DatabasePort = v
		case "DD_DATABASE_TYPE":
			conf.Settings.DatabaseType = v
		case "DD_DATABASE_URL":
			conf.Settings.DatabaseURL = v
		case "DD_DATABASE_USER":
			conf.Settings.DatabaseUser = v
		case "DD_DATA_UPLOAD_MAX_MEMORY_SIZE":
			conf.Settings.DataUploadMaxMemorySize = convInt(v, "DD_DATA_UPLOAD_MAX_MEMORY_SIZE provided via environmental variable isn't a valid number")
		case "DD_DEBUG":
			conf.Settings.Debug = convBool(v, "DD_DEBUG environmental variable was not a boolean.")
		case "DD_DJANGO_ADMIN_ENABLED":
			conf.Settings.DjangoAdminEnabled = convBool(v, "DD_DJANGO_ADMIN_ENABLED environmental variable was not a boolean.")
		case "DD_EMAIL_URL":
			conf.Settings.EmailURL = v
		case "DD_ENV":
			conf.Settings.Env = v
		case "DD_ENV_PATH":
			conf.Settings.EnvPath = v
		case "DD_FORCE_LOWERCASE_TAGS":
			conf.Settings.ForceLowercaseTags = convBool(v, "DD_FORCE_LOWERCASE_TAGS environmental variable was not a boolean.")
		case "DD_HOST":
			conf.Settings.Host = v
		case "DD_INITIALIZE":
			conf.Settings.Initialize = v
		case "DD_LANG":
			conf.Settings.Lang = v
		case "DD_LANGUAGE_CODE":
			conf.Settings.LanguageCode = v
		case "DD_LOGIN_REDIRECT_URL":
			conf.Settings.LoginRedirectURL = v
		case "DD_MAX_TAG_LENGTH":
			// TODO: Look up maximum tag length in data model and check for that too
			conf.Settings.MaxTagLength = convInt(v, "DD_MAX_TAG_LENGTH provided via environmental variable isn't a valid number")
		case "DD_MEDIA_ROOT":
			conf.Settings.MediaRoot = v
		case "DD_MEDIA_URL":
			conf.Settings.MediaURL = v
		case "DD_PORT":
			conf.Settings.Port = v
		case "DD_PORT_SCAN_CONTACT_EMAIL":
			conf.Settings.PortScanContactEmail = v
		case "DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST":
			conf.Settings.PortScanExternalUnitEmailList = v
		case "DD_PORT_SCAN_RESULT_EMAIL_FROM":
			conf.Settings.PortScanResultEmailFrom = v
		case "DD_PORT_SCAN_SOURCE_IP":
			conf.Settings.PortScanSourceIP = v
		case "DD_ROOT":
			conf.Settings.Root = v
		case "DD_SECRET_KEY":
			conf.Settings.SecretKey = v
		case "DD_SECURE_BROWSER_XSS_FILTER":
			conf.Settings.SecureBrowserXSSFilter = convBool(v, "DD_SECURE_BROWSER_XSS_FILTER environmental variable was not a boolean.")
		case "DD_SECURE_CONTENT_TYPE_NOSNIFF":
			conf.Settings.SecureContentTypeNosniff = v
		case "DD_SECURE_HSTS_INCLUDE_SUBDOMAINS":
			conf.Settings.SecureHSTSIncludeSubdomains = convBool(v, "DD_SECURE_HSTS_INCLUDE_SUBDOMAINS environmental variable was not a boolean.")
		case "DD_SECURE_HSTS_SECONDS":
			conf.Settings.SecureHSTSSeconds = convInt(v, "DD_SECURE_HSTS_SECONDS provided via environmental variable isn't a valid number")
		case "DD_SECURE_PROXY_SSL_HEADER":
			conf.Settings.SecureProxySSLHeader = convBool(v, "DD_SECURE_PROXY_SSL_HEADER environmental variable was not a boolean.")
		case "DD_SECURE_SSL_REDIRECT":
			conf.Settings.SecureSSLRedirect = convBool(v, "DD_SECURE_SSL_REDIRECT environmental variable was not a boolean.")
		case "DD_SESSION_COOKIE_HTTPONLY":
			conf.Settings.SessionCookieHTTPOnly = convBool(v, "DD_SESSION_COOKIE_HTTPONLY environmental variable was not a boolean.")
		case "DD_SESSION_COOKIE_SECURE":
			conf.Settings.SessionCookieSecure = convBool(v, "DD_SESSION_COOKIE_SECURE environmental variable was not a boolean.")
		case "DD_SITE_ID":
			conf.Settings.SiteID = convInt(v, "DD_SITE_ID provided via environmental variable isn't a valid number")
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_ENABLED":
			conf.Settings.SocialAuthAzureadTenantOauth2Enabled = v
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_KEY":
			conf.Settings.SocialAuthAzureadTenantOauth2Key = v
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_RESOURCE":
			conf.Settings.SocialAuthAzureadTenantOauth2Resource = v
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_SECRET":
			conf.Settings.SocialAuthAzureadTenantOauth2Secret = v
		case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_TENANT_ID":
			conf.Settings.SocialAuthAzureadTenantOauth2TenantID = v
		case "DD_SOCIAL_AUTH_GOOGLE_OAUTH2_ENABLE":
			conf.Settings.SocialAuthGoogleOauth2Enable = v
		case "DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY":
			conf.Settings.SocialAuthGoogleOauth2Key = v
		case "DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET":
			conf.Settings.SocialAuthGoogleOauth2Secret = v
		case "DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL":
			conf.Settings.SocialAuthOktaOauth2APIURL = v
		case "DD_SOCIAL_AUTH_OKTA_OAUTH2_ENABLED":
			conf.Settings.SocialAuthOktaOauth2Enabled = v
		case "DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY":
			conf.Settings.SocialAuthOktaOauth2Key = v
		case "DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET":
			conf.Settings.SocialAuthOktaOauth2Secret = v
		case "DD_STATIC_ROOT":
			conf.Settings.StaticRoot = v
		case "DD_STATIC_URL":
			conf.Settings.StaticURL = v
		case "DD_TEAM_NAME":
			conf.Settings.TeamName = v
		case "DD_TEST_DATABASE_NAME":
			conf.Settings.TestDatabaseName = v
		case "DD_TEST_RUNNER":
			conf.Settings.TestRunner = v
		case "DD_TIME_ZONE":
			conf.Settings.TimeZone = v
		case "DD_TRACK_MIGRATIONS":
			conf.Settings.TrackMigrations = convBool(v, "DD_TRACK_MIGRATIONS environmental variable was not a boolean.")
		case "DD_URL_PREFIX":
			conf.Settings.URLPrefix = v
		case "DD_USE_I18N":
			conf.Settings.UseI18N = convBool(v, "DD_USE_I18N environmental variable was not a boolean.")
		case "DD_USE_L10N":
			conf.Settings.UseL10N = convBool(v, "DD_USE_L10N environmental variable was not a boolean.")
		case "DD_USE_TZ":
			conf.Settings.UseTZ = convBool(v, "DD_USE_TZ environmental variable was not a boolean.")
		case "DD_UUID":
			conf.Settings.UUID = v
		case "DD_UWSGI_ENDPOINT":
			conf.Settings.UwsgiEndpoint = v
		case "DD_UWSGI_HOST":
			conf.Settings.UwsgiHost = v
		case "DD_UWSGI_MODE":
			conf.Settings.UwsgiMode = v
		case "DD_UWSGI_PASS":
			conf.Settings.UwsgiPass = v
		case "DD_UWSGI_PORT":
			conf.Settings.UwsgiPort = v
		case "DD_WHITENOISE":
			conf.Settings.Whitenoise = convBool(v, "DD_WHITENOISE environmental variable was not a boolean.")
		case "DD_WKHTMLTOPDF":
			conf.Settings.Wkhtmltopdf = v
		case "DOJO_ADMIN_USER":
			conf.Settings.DojoAdminUser = v
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
