package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mtesauro/godojo/config"
)

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
// Based on https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func Untar(dst string, r io.Reader) error {

	// Setup new gzip Reader to extract tarball contents
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	// Loop through the file reading each header to determine if its a file or directory
	// then either create the directory (if needed) or create the file
	for {
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil
		// return any other error
		case err != nil:
			return err
		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// check the file type
		switch header.Typeflag {
		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			// TODO: Reformat me
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			// TODO: Reformat me
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

// Redactatron - redacts sensitive information from being written to the logs
// Redaction is configurable with Install's Redact boolean config.
// If true (the default), sensitive info will be redacted
func Redactatron(l string, on bool) string {
	// Redact sensitive info from the files in ./logs/
	clean := l
	r := "=[REDACTED]="
	// Redact sensitive data if it's turned on
	if on {
		for i := 0; i < len(sensStr); i++ {
			if strings.Contains(clean, sensStr[i]) {
				clean = strings.Replace(clean, sensStr[0], r, -1)
			}
		}
	}
	return clean
}

// InitRedactatron - sets up the data to be redacted by Redactatron
func InitRedact(conf *config.DojoConfig) {
	// Add the strings from DojoConfig to be redacted
	sensStr[0] = conf.Install.DB.Root
	sensStr[1] = conf.Install.DB.Pass
	sensStr[2] = conf.Install.OS.Pass
	sensStr[3] = conf.Install.Admin.Pass
	sensStr[4] = conf.Settings.Celery.Broker.Password
	sensStr[5] = conf.Settings.Database.Password
	sensStr[6] = conf.Settings.Secret.Key
	sensStr[7] = conf.Settings.Credential.AES.B256.Key
	sensStr[8] = conf.Settings.Social.Auth.Google.OAUTH2.Key
	sensStr[9] = conf.Settings.Social.Auth.Google.OAUTH2.Secret
	sensStr[10] = conf.Settings.Social.Auth.Okta.OAUTH2.Key
	sensStr[11] = conf.Settings.Social.Auth.Okta.OAUTH2.Secret
}
