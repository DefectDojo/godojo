package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	//if conf.Install.Redact {
	if on {
		// Config says to remove sentivite info from output
		r := "REDACTED"
		fmt.Printf("FIXME = %+v\n", r)
		//clean = strings.Replace(l, conf.Install.DB.Root, r, -1)
		//clean = strings.Replace(l, conf.Install.DB.Pass, r, -1)
		//clean = strings.Replace(l, conf.Install.OS.Pass, r, -1)
		//clean = strings.Replace(l, conf.Install.Admin.Pass, r, -1)
		//clean = strings.Replace(l, conf.Settings.Celery.Broker.Password, r, -1)
		//clean = strings.Replace(l, conf.Settings.Database.Password, r, -1)
		//clean = strings.Replace(l, conf.Settings.Secret.Key, r, -1)
		//clean = strings.Replace(l, conf.Settings.Credential.AES.B256.Key, r, -1)
		//clean = strings.Replace(l, conf.Settings.Social.Auth.Google.OAUTH2.Key, r, -1)
		//clean = strings.Replace(l, conf.Settings.Social.Auth.Google.OAUTH2.Secret, r, -1)
		//clean = strings.Replace(l, conf.Settings.Social.Auth.Okta.OAUTH2.Key, r, -1)
		//clean = strings.Replace(l, conf.Settings.Social.Auth.Okta.OAUTH2.Secret, r, -1)
		// Add more lines here if new sensitive data is added to the DojoConfig struct
	}
	return clean
}

// InitRedactatron - sets up the data to be redacted by Redactatron
func InitRedactatron(conf *config.DojoConfig) string {
	// Redact sensitive info from the files in ./logs/
	clean := "FIXME" //TODO
	if true {
		// Config says to remove sentivite info from output
		r := "REDACTED"
		fmt.Printf("FIXME = %+v\n", r)
		//clean = strings.Replace(l, conf.Install.DB.Root, r, -1)
		//clean = strings.Replace(l, conf.Install.DB.Pass, r, -1)
		//clean = strings.Replace(l, conf.Install.OS.Pass, r, -1)
		//clean = strings.Replace(l, conf.Install.Admin.Pass, r, -1)
		//clean = strings.Replace(l, conf.Settings.Celery.Broker.Password, r, -1)
		//clean = strings.Replace(l, conf.Settings.Database.Password, r, -1)
		//clean = strings.Replace(l, conf.Settings.Secret.Key, r, -1)
		//clean = strings.Replace(l, conf.Settings.Credential.AES.B256.Key, r, -1)
		//clean = strings.Replace(l, conf.Settings.Social.Auth.Google.OAUTH2.Key, r, -1)
		//clean = strings.Replace(l, conf.Settings.Social.Auth.Google.OAUTH2.Secret, r, -1)
		//clean = strings.Replace(l, conf.Settings.Social.Auth.Okta.OAUTH2.Key, r, -1)
		//clean = strings.Replace(l, conf.Settings.Social.Auth.Okta.OAUTH2.Secret, r, -1)
		// Add more lines here if new sensitive data is added to the DojoConfig struct
	}
	return clean
}
