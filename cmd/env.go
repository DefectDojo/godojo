package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"
)

// Handles the template-based generation of env.prod for DefectDojo's settings.py

// Define the template
const envProd = `
# Django Debug, don't enable on production! - default is off
#DD_DEBUG={{.DD_DEBUG}}

# Enables Django Admin - default is on
#DD_DJANGO_ADMIN_ENABLED={{.DD_DJANGO_ADMIN_ENABLED}}

# A secret key for a particular Django installation.
DD_SECRET_KEY={{.DD_SECRET_KEY}}

# Key for encrypting credentials in the manager
DD_CREDENTIAL_AES_256_KEY={{.DD_CREDENTIAL_AES_256_KEY}}

# Database URL, options: postgres://, mysql://, sqlite://, to use unsafe characters encode with urllib.parse.encode
DD_DATABASE_URL={{.DD_DATABASE_URL}}

# Hosts/domain names that are valid for this site;
DD_ALLOWED_HOSTS={{.DD_ALLOWED_HOSTS}}

# CSRF Trusted Origins
DD_CSRF_TRUSTED_ORIGINS="{{.DD_APP_HOSTNAME}}"

# WhiteNoise allows your web app to serve its own static files,
# making it a self-contained unit that can be deployed anywhere without relying on nginx,
# if using nginx then disable Whitenoise
DD_WHITENOISE={{.DD_WHITENOISE}}

`

type envVals struct {
	DD_DEBUG                              bool
	DD_DJANGO_ADMIN_ENABLED               bool
	DD_SECRET_KEY                         string
	DD_CREDENTIAL_AES_256_KEY             string
	DD_DATABASE_URL                       string
	DD_ALLOWED_HOSTS                      string
	DD_APP_HOSTNAME                       string
	DD_WHITENOISE                         bool
	DD_TIME_ZONE                          string
	DD_TRACK_MIGRATIONS                   bool
	DD_SESSION_COOKIE_HTTPONLY            bool
	DD_CSRF_COOKIE_HTTPONLY               bool
	DD_SECURE_SSL_REDIRECT                bool
	DD_CSRF_COOKIE_SECURE                 bool
	DD_SECURE_BROWSER_XSS_FILTER          bool
	DD_LANG                               string
	DD_WKHTMLTOPDF                        string
	DD_TEAM_NAME                          string
	DD_ADMINS                             string
	DD_PORT_SCAN_CONTACT_EMAIL            string
	DD_PORT_SCAN_RESULT_EMAIL_FROM        string
	DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST string
	DD_PORT_SCAN_SOURCE_IP                string
}

func genAndWriteEnv(d *DDConfig, dbURL string) {
	// Generate randon values for the two keys below
	secretKey := d.conf.Settings.SecretKey
	if len(secretKey) < 28 {
		// Handle the case that the key wasn't configured
		s1 := make([]byte, 42)
		_, err := rand.Read(s1)
		if err != nil {
			d.errorMsg("Error generating random data for encryption keys")
			os.Exit(1)
		}
		secretKey = base64.StdEncoding.EncodeToString(s1)
	}
	credentialKey := d.conf.Settings.CredentialAES256Key
	if len(credentialKey) < 28 {
		// Handle the case that the key wasn't configured
		s2 := make([]byte, 42)
		_, err := rand.Read(s2)
		if err != nil {
			d.errorMsg("Error generating random data for encryption keys")
			os.Exit(1)
		}
		credentialKey = base64.StdEncoding.EncodeToString(s2)
	}

	// Set the values from the configuration file
	env := envVals{
		DD_DEBUG:                              d.conf.Settings.Debug,
		DD_DJANGO_ADMIN_ENABLED:               d.conf.Settings.DjangoAdminEnabled,
		DD_SECRET_KEY:                         secretKey,
		DD_CREDENTIAL_AES_256_KEY:             credentialKey,
		DD_DATABASE_URL:                       dbURL,
		DD_ALLOWED_HOSTS:                      d.conf.Settings.AllowedHosts,
		DD_APP_HOSTNAME:                       d.conf.Settings.AppHostname,
		DD_WHITENOISE:                         d.conf.Settings.Whitenoise,
		DD_TIME_ZONE:                          d.conf.Settings.TimeZone,
		DD_TRACK_MIGRATIONS:                   d.conf.Settings.TrackMigrations,
		DD_SESSION_COOKIE_HTTPONLY:            d.conf.Settings.SessionCookieHTTPOnly,
		DD_CSRF_COOKIE_HTTPONLY:               d.conf.Settings.CSRFCookieHTTPOnly,
		DD_SECURE_SSL_REDIRECT:                d.conf.Settings.SecureSSLRedirect,
		DD_CSRF_COOKIE_SECURE:                 d.conf.Settings.CSRFCookieSecure,
		DD_SECURE_BROWSER_XSS_FILTER:          d.conf.Settings.SecureBrowserXSSFilter,
		DD_LANG:                               d.conf.Settings.Lang,
		DD_WKHTMLTOPDF:                        d.conf.Settings.Wkhtmltopdf,
		DD_TEAM_NAME:                          d.conf.Settings.TeamName,
		DD_ADMINS:                             d.conf.Settings.Admins,
		DD_PORT_SCAN_CONTACT_EMAIL:            d.conf.Settings.PortScanContactEmail,
		DD_PORT_SCAN_RESULT_EMAIL_FROM:        d.conf.Settings.PortScanResultEmailFrom,
		DD_PORT_SCAN_EXTERNAL_UNIT_EMAIL_LIST: d.conf.Settings.PortScanExternalUnitEmailList,
		DD_PORT_SCAN_SOURCE_IP:                d.conf.Settings.PortScanSourceIP,
	}

	// Create a template based on the text above
	t := template.Must(template.New("envProd").Parse(envProd))

	// Open a file to write the contents of the parsed template
	d.traceMsg(fmt.Sprintf("Location of env file is %+v/django-DefectDojo/dojo/settings/.env.prod\n", d.conf.Install.Root))
	f, err := os.Create(d.conf.Install.Root + "/django-DefectDojo/dojo/settings/.env.prod")
	if err != nil {
		d.errorMsg("Unable to create .env.prod file for settings.py configuration")
		os.Exit(1)
	}
	defer f.Close()

	// Make substitutions in the template
	err = t.Execute(f, env)
	if err != nil {
		d.errorMsg("Failed to create .env.prod from template")
		os.Exit(1)
	}
}
