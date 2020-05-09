package internal

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	appDir           = "ignition-master/site/"
	appDirPattern    = "^" + appDir
	appReplace       = "\"github.com/cjtoolkit/ignition/site"
	appReplaceModule = "module github.com/cjtoolkit/ignition/site"
)

func BuildApp(dir, moduleName, baseModuleName string) {
	moduleNamePrefix := "module " + moduleName
	dir = filepath.FromSlash(dir)
	err := os.Mkdir(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	tr := getFile()

	var (
		appDirPattern    = regexp.MustCompile(appDirPattern)
		goPattern        = regexp.MustCompile(baseGoPattern)
		gitIgnorePattern = regexp.MustCompile(baseGitIgnorePattern)
	)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if !appDirPattern.MatchString(hdr.Name) || hdr.Name == appDir {
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			err = os.Mkdir(dir+filepath.FromSlash("/"+strings.TrimPrefix(hdr.Name, appDir)), hdr.FileInfo().Mode())
			if err != nil {
				log.Fatal(err)
			}
		case tar.TypeReg:
			b, err := ioutil.ReadAll(tr)
			if err != nil {
				log.Fatal(err)
			}

			fileName := dir + filepath.FromSlash("/"+strings.TrimPrefix(hdr.Name, appDir))

			if goPattern.MatchString(hdr.Name) {
				b = bytes.ReplaceAll(b, []byte(baseReplace), []byte("\""+baseModuleName))
				b = bytes.ReplaceAll(b, []byte(appReplace), []byte("\""+moduleName))
			} else if hdr.Name == appDir+"go.mod" {
				b = bytes.Replace(b, []byte(appReplaceModule), []byte(moduleNamePrefix), 1)
			} else if gitIgnorePattern.MatchString(hdr.Name) {
				fileName = filepath.Dir(fileName) + filepath.FromSlash("/.gitignore")
			} else if hdr.Name == appDir+"doTest" {
				b = bytes.Replace(b, []byte(baseReplace[1:]), []byte(baseModuleName), 1)
			}

			// write a file
			w, err := os.OpenFile(fileName,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC, hdr.FileInfo().Mode())
			if err != nil {
				log.Fatal(err)
			}
			_, err = w.Write(b)
			if err != nil {
				log.Fatal(err)
			}
			err = w.Close()
			if err != nil {
				panic(err)
			}
		}
	}
}
