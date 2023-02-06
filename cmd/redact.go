package cmd

import (
	"strings"
)

// Redactatron - redacts sensitive information from being written to the logs
// Redaction is configurable with Install's Redact boolean config.
// If true (the default), sensitive info will be redacted
func (d *gdjDefault) redactatron(l string, on bool) string {
	// Redact sensitive data if it's turned on
	if on {
		// Redact sensitive info from the files in ./logs/
		clean := l
		r := "[~REDACTED~]"
		for i := range d.sensStr {
			if strings.Contains(clean, d.sensStr[i]) {
				clean = strings.Replace(clean, d.sensStr[i], r, -1)
			}
		}
		return clean
	}
	return l
}

// initRedact - sets up the data to be redacted by Redactatron
func (d *gdjDefault) initRedact() {
	// Setup Default strings to redact
	l := []string{
		d.conf.Install.DB.Rpass,
		d.conf.Install.DB.Pass,
		d.conf.Install.OS.Pass,
		d.conf.Install.Admin.Pass,
		d.conf.Settings.CeleryBrokerPassword,
		d.conf.Settings.DatabasePassword,
		d.conf.Settings.SecretKey,
		d.conf.Settings.CredentialAES256Key,
		d.conf.Settings.SocialAuthGoogleOauth2Key,
		d.conf.Settings.SocialAuthGoogleOauth2Secret,
		d.conf.Settings.SocialAuthOktaOauth2Key,
		d.conf.Settings.SocialAuthOktaOauth2Secret,
	}

	// Add the strings from DojoConfig to be redacted if they have content
	for i := range l {
		if len(l[i]) > 0 {
			d.sensStr = append(d.sensStr, l[i])
		}
	}
}

func (d *gdjDefault) addRedact(s string) {
	// Add an additional string to redact from the logs
	d.sensStr = append(d.sensStr, s)
}
