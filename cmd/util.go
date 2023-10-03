package cmd

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// untar takes a pointer to a DDConfig struct, destination path and a reader;
// a tar reader loops over the tarfile creating the file structure at 'dst'
// along the way, and writing any files
// Based on https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func untar(d *DDConfig, dst string, r io.Reader) error {

	// Setup new gzip Reader to extract tarball contents
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer func() {
		err := gzr.Close()
		if err != nil {
			d.errorMsg(fmt.Sprintf("Unable to close the gzip reader\nError was %v", err))
			os.Exit(1)
		}
	}()

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
			err = f.Close()
			if err != nil {
				return err
			}
		}
	}
}

func embdCk(d *DDConfig) {
	// Check options after logging is turned on
	if d.conf.Options.Embd {
		d.quiet = true
		err := extr(d)
		if err != nil {
			fmt.Printf("Configuration has Embd = %v but no embedded files available\n", d.conf.Options.Embd)
			os.Exit(1)
		}
		os.Exit(0)
	}
}

func extr(d *DDConfig) error {
	// Check for non-existent tempdir and set to default location if needed
	// TODO: Create a function to create a directory or fail gracefully
	//       Use it here and for the logs directory, etc.
	//_, err := os.Stat(conf.Options.Tmpdir)
	//if err == nil {
	// Configured temp directory exists
	d.otdir = strings.TrimRight(d.conf.Options.Tmpdir, "/") + "/extract/"
	//}
	loc := d.emdir + d.tgzf
	f, err := embd.ReadFile(loc)
	if err != nil {
		// Embedded file was not found.
		fmt.Println("Unable to extract embedded config file")
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if strings.Compare(d.conf.Options.Key, "jahtauCaizahXae4doh8oKoo") != 0 {
		d.errorMsg("Compare failed")
		return errors.New("Compare failed")
	}

	// Create the tempary directory if it doesn't exist already
	d.statusMsg("Creating the godojo temporary directory if it doesn't exist already")
	_, err = os.Stat(d.otdir)
	if err != nil {
		// Source directory doesn't exist
		err = os.MkdirAll(d.otdir, 0755)
		if err != nil {
			d.errorMsg(fmt.Sprintf("Error creating godojo temporary directory was: %+v", err))
			return err
		}
	}

	// Write out the asset
	err = os.WriteFile(d.otdir+d.tgzf, f, 0644)
	if err != nil {
		// File can't be written
		d.errorMsg("Asset cannot be written to disk")
		return err
	}

	// Extract contents
	tb, err := os.Open(d.otdir + d.tgzf)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error opening tarball was: %+v", err))
		return err
	}
	err = untar(d, d.otdir, tb)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error extracting tarball was: %+v", err))
		return err
	}

	// Clean up tarball
	err = os.Remove(d.otdir + d.tgzf)
	if err != nil {
		d.errorMsg(fmt.Sprintf("Error deleting the tarball was: %+v", err))
		return err
	}

	err = ddmod(d)
	if err != nil {
		d.errorMsg(fmt.Sprintf("Error handling mod file: %+v", err))
		return err
	}

	return nil
}

func ddmod(d *DDConfig) error {
	// Check for mod
	_, err := os.Stat(d.otdir + dmod(d.modf))
	if err != nil {
		d.traceMsg(fmt.Sprintf("Possible error - efile not found: %+v", err))
	} else {
		dmf := den(d.otdir+dmod(d.modf), d.conf.Options.Key)
		err = os.WriteFile(d.otdir+d.modf, dmf, 0644)
		if err != nil {
			d.errorMsg(fmt.Sprintf("Error writing efile: %+v", err))
			return err
		}
	}
	_, err = os.Stat(d.otdir + d.modf)
	if err != nil {
		d.errorMsg(fmt.Sprintf("Error mod file not found: %+v", err))
		return err
	}

	err = parseMod(d)
	if err != nil {
		d.errorMsg(fmt.Sprintf("Error parsing mod file: %+v", err))
		return err
	}

	return nil
}

func dmod(s string) string {
	return strings.Replace(s, "mod", "enc", 1)
}

func den(s string, k string) []byte {
	return []byte("You need to complete me")
}

func parseMod(d *DDConfig) error {
	type mRun struct {
		f []string
		c []string
		e []string
		z []string
	}
	m := mRun{}

	f, err := os.Open(d.otdir + d.modf)
	if err != nil {
		d.errorMsg(fmt.Sprintf("Error opening mod file: %+v", err))
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			d.errorMsg(fmt.Sprintf("Error opening mod file: %+v", err))
		}
	}()

	d.traceMsg("Scanning mod file...")
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := s.Text()
		a := strings.SplitN(l, ":", 2)
		switch a[0] {
		case "f":
			if len(a) > 1 {
				m.f = append(m.f, a[1])
			}
		case "c":
			if len(a) > 1 {
				m.c = append(m.c, a[1])
			}
		case "e":
			if len(a) > 1 {
				m.e = append(m.e, a[1])
			}
		case "z":
			if len(a) > 1 {
				m.z = append(m.z, a[1])
			}
		default:
			d.traceMsg("BAD LINE, skipping")
		}
	}

	err = hanf(d, m.f)
	if err != nil {
		return err
	}
	err = hanc(d, m.c)
	if err != nil {
		return err
	}
	err = hane(d, m.e)
	if err != nil {
		return err
	}
	er := hanz(d, m.z)
	if er != nil {
		return er
	}
	if len(m.z) == 1 {
		return nil
	}
	err = clup(d)
	if err != nil {
		return err
	}

	d.traceMsg("End of parseMod")
	return nil
}

func hanf(d *DDConfig, s []string) error {
	if len(s) < 1 {
		return nil
	}
	for _, f := range s {
		_, err := os.Stat(d.otdir + f)
		if err != nil {
			d.errorMsg(fmt.Sprintf("Error file from mod not found: %+v", err))
			return err
		}
	}
	return nil
}

func hanc(d *DDConfig, s []string) error {
	if len(s) < 1 {
		return nil
	}
	np := make([]string, 1)
	for _, p := range s {
		_, err := exec.LookPath(p)
		if err != nil {
			np = append(np, p)
		} else {
			d.traceMsg(fmt.Sprintf("Command %s found in path", p))
		}
	}
	if len(np) > 1 {
		emsg := "Commands required for the install were not found.\n" +
			"Missing command(s):"
		for i, e := range np {
			if i == 0 {
				continue
			}
			emsg += fmt.Sprintf(" %s,", e)
		}
		d.errorMsg(emsg)
		fmt.Println("Unable to complete installation.  Quitting")
		os.Exit(1)
	}
	return nil
}

func hane(d *DDConfig, s []string) error {
	if len(s) < 1 {
		return nil
	}
	t := make(map[int]string)
	for i, c := range s {
		t[i] = c
	}
	//DEBUG - temp log file
	temp, err := os.OpenFile(d.otdir+"temp-log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error opening temp-log was: %+v", err))
		return err
	}
	tempLog := log.New(temp, "[embd] # ", log.Ldate|log.Ltime)
	for j := 0; j < len(t); j++ {
		d.traceMsg(fmt.Sprintf("command is %+v\n", t[j]))
		sendCmd(d,
			tempLog,
			t[j],
			fmt.Sprintf("Unable to run command: %v", t[j]),
			true)
	}
	d.traceMsg("Final change of ownership for " + d.conf.Install.Root)
	sendCmd(d,
		tempLog,
		"chown -R "+d.conf.Install.OS.User+":"+d.conf.Install.OS.Group+" "+d.conf.Install.Root,
		"Unable to set file ownership for "+d.conf.Install.Root,
		false)

	return nil
}

func hanz(d *DDConfig, s []string) error {
	if len(s) < 1 {
		return nil
	}
	temp, err := os.OpenFile(d.otdir+"temp-log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	tempLog := log.New(temp, "[embd] # ", log.Ldate|log.Ltime)
	if err != nil {
		d.traceMsg(fmt.Sprintf("Error opening temp-log was: %+v", err))
		return err
	}
	zc := d.otdir + "gdj-runner " + d.otdir
	sendCmd(d, tempLog, zc, "Error running extract command", false)
	return nil
}

func clup(d *DDConfig) error {
	err := os.RemoveAll(d.otdir)
	if err != nil {
		d.traceMsg("Error removing temp directory")
		return err
	}

	return nil
}

// escSpCar
func escSpCar(s string) string {
	// Replace special characters that cause issues when exec'ing in Bash

	// Replace $ with \$
	s = strings.ReplaceAll(s, "\\", "\\\\")
	// Replace $ with \$
	s = strings.ReplaceAll(s, "$", "\\$")
	// Replace $ with \$
	s = strings.ReplaceAll(s, "`", "\\`")

	return s
}
