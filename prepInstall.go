package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// readConfigFile reads the yaml configuration file for godojo
// to determine runtime configuration.  The file is dojoConfig.yml
// and is expected to be in the same directory as the godojo binary
// It returns nohing but will exit early with a exit code of 1
// if there are errors reading the file or unmarshialling into a struct
func readConfigFile() {
	// Setup viper config
	viper.AddConfigPath(".")
	viper.SetConfigName("dojoConfig")

	// Setup ENV variables
	// TODO: Do these manually in readEnvVars() since they have odd names for Viper auto-magic
	//viper.SetEnvPrefix("DD")
	//replace := strings.NewReplacer(".", "_")
	//viper.SetEnvKeyReplacer(replace)
	//viper.AutomaticEnv()

	// Read the default config file dojoConfig.yml
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("")
		fmt.Println("Unable to read the godojo config file (dojoConfig.yml), exiting install")
		os.Exit(1)
	}
	// Marshall the config values into the DojoConfig struct
	err = viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println("")
		fmt.Println("Unable to set the config values based on config file and ENV variables, exiting install")
		os.Exit(1)
	}
}

// readEnvVars reads the DefectDojo supported environmental variables and
// overrides any options set in the configuration file. These variables
// are used to supply either install-time configurations or provide values
// that are used in DefectDojo's settings.py configuration file
func readEnvVars() {
	// Env variables pulled from repo
	// Add newly supported env vars below
	// and to the case statement below after "if match {"
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
		// TYPO in the source - "DD_DATABASE_PASsWORD": true,
		// TODO: Fix https://github.com/mtesauro/godojo/issues/1
	} // End of dojoEnvs declaration

	match := false
	overrides := make(map[string]string)
	for _, e := range os.Environ() {
		// Pull out each env var into a slice
		env := strings.SplitN(e, "=", 2)

		// Check if the name of the env var matches the supported map
		if _, ok := dojoEnvs[env[0]]; ok {
			fmt.Printf("Found %+v\n", env[0])
			overrides[env[0]] = env[1]
			match = true
		}
	}

	// Override config values if we found matching Env vars
	if match {
		fmt.Println("Yup, we found a match")
		for k, v := range overrides {
			// Set DojoConfig struct values from Env variables to override config values
			// Have to do this as a switch statement as there's no sanity to DefectDojo env var naming
			switch k {
			case "DD_ADMIN_FIRST_NAME":
			case "DD_ADMIN_LAST_NAME":
			case "DD_ADMIN_MAIL":
			case "DD_ADMIN_PASSWORD":
			case "DD_ADMINS":
			case "DD_ADMIN_USER":
			case "DD_ALLOWED_HOSTS":
			case "DD_CELERY_BEAT_SCHEDULE_FILENAME":
			case "DD_CELERY_BROKER_HOST":
			case "DD_CELERY_BROKER_PASSWORD":
			case "DD_CELERY_BROKER_PATH":
			case "DD_CELERY_BROKER_PORT":
			case "DD_CELERY_BROKER_SCHEME":
			case "DD_CELERY_BROKER_URL":
			case "DD_CELERY_BROKER_USER":
			case "DD_CELERY_LOG_LEVEL":
			case "DD_CELERY_RESULT_BACKEND":
			case "DD_CELERY_RESULT_EXPIRES":
			case "DD_CELERY_TASK_IGNORE_RESULT":
			case "DD_CELERY_TASK_SERIALIZER":
			case "DD_CREDENTIAL_AES_256_KEY":
			case "DD_CSRF_COOKIE_HTTPONLY":
			case "DD_CSRF_COOKIE_SECURE":
			case "DD_DATABASE_ENGINE":
			case "DD_DATABASE_HOST":
			case "DD_DATABASE_NAME":
			case "DD_DATABASE_PASSWORD":
			case "DD_DATABASE_PORT":
			case "DD_DATABASE_TYPE":
			case "DD_DATABASE_URL":
			case "DD_DATABASE_USER":
			case "DD_DATA_UPLOAD_MAX_MEMORY_SIZE":
			case "DD_DEBUG":
			case "DD_DJANGO_ADMIN_ENABLED":
			case "DD_EMAIL_URL":
			case "DD_ENV":
			case "DD_ENV_PATH":
			case "DD_FORCE_LOWERCASE_TAGS":
			case "DD_HOST":
			case "DD_INITIALIZE":
			case "DD_LANG":
			case "DD_LANGUAGE_CODE":
			case "DD_LOGIN_REDIRECT_URL":
			case "DD_MAX_TAG_LENGTH":
			case "DD_MEDIA_ROOT":
			case "DD_MEDIA_URL":
			case "DD_PORT":
			case "DD_PORT_SCAN_CONTACT_EMAIL":
			case "DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST":
			case "DD_PORT_SCAN_RESULT_EMAIL_FROM":
			case "DD_PORT_SCAN_SOURCE_IP":
			case "DD_ROOT":
			case "DD_SECRET_KEY":
			case "DD_SECURE_BROWSER_XSS_FILTER":
			case "DD_SECURE_CONTENT_TYPE_NOSNIFF":
			case "DD_SECURE_HSTS_INCLUDE_SUBDOMAINS":
			case "DD_SECURE_HSTS_SECONDS":
			case "DD_SECURE_PROXY_SSL_HEADER":
			case "DD_SECURE_SSL_REDIRECT":
			case "DD_SESSION_COOKIE_HTTPONLY":
			case "DD_SESSION_COOKIE_SECURE":
			case "DD_SITE_ID":
			case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_ENABLED":
			case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_KEY":
			case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_RESOURCE":
			case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_SECRET":
			case "DD_SOCIAL_AUTH_AZUREAD_TENANT_OAUTH2_TENANT_ID":
			case "DD_SOCIAL_AUTH_GOOGLE_OAUTH2_ENABLE":
			case "DD_SOCIAL_AUTH_GOOGLE_OAUTH2_KEY":
			case "DD_SOCIAL_AUTH_GOOGLE_OAUTH2_SECRET":
			case "DD_SOCIAL_AUTH_OKTA_OAUTH2_API_URL":
			case "DD_SOCIAL_AUTH_OKTA_OAUTH2_ENABLED":
			case "DD_SOCIAL_AUTH_OKTA_OAUTH2_KEY":
			case "DD_SOCIAL_AUTH_OKTA_OAUTH2_SECRET":
			case "DD_STATIC_ROOT":
			case "DD_STATIC_URL":
			case "DD_TEAM_NAME":
			case "DD_TEST_DATABASE_NAME":
			case "DD_TEST_RUNNER":
			case "DD_TIME_ZONE":
			case "DD_TRACK_MIGRATIONS":
			case "DD_URL_PREFIX":
			case "DD_USE_I18N":
			case "DD_USE_L10N":
			case "DD_USE_TZ":
			case "DD_UUID":
				// TODO: Delete me - I for debugging only
				fmt.Println("DD_UUID found :-)")
				fmt.Printf("Will set config to %+v\n", v)
				// Set config var here
			case "DD_UWSGI_ENDPOINT":
			case "DD_UWSGI_HOST":
			case "DD_UWSGI_MODE":
			case "DD_UWSGI_PASS":
			case "DD_UWSGI_PORT":
			case "DD_WHITENOISE":
			case "DD_WKHTMLTOPDF":
			case "DOJO_ADMIN_USER":
			}
		}
	}

	//fmt.Printf("Supported environmental variables are %+v", dojoEnvs)
}
