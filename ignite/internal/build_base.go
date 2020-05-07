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
	baseDir              = "ignition-master/shared/"
	baseDirPattern       = "^" + baseDir
	baseGoPattern        = ".go$"
	baseGoModPattern     = "go.mod$"
	baseGitIgnorePattern = "gitignore.txt"
	baseReplace          = "\"github.com/cjtoolkit/ignition/shared"
	baseReplaceModule    = "module github.com/cjtoolkit/ignition/shared"
)

func BuildBase(dir, moduleName string) {
	moduleNamePrefix := "module " + moduleName
	dir = filepath.FromSlash(dir)
	err := os.Mkdir(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	tr := getFileFromGitHub()

	var (
		dirPattern       = regexp.MustCompile(baseDirPattern)
		goPattern        = regexp.MustCompile(baseGoPattern)
		goModPattern     = regexp.MustCompile(baseGoModPattern)
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

		if !dirPattern.MatchString(hdr.Name) || hdr.Name == baseDir {
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			err = os.Mkdir(dir+filepath.FromSlash("/"+strings.TrimPrefix(hdr.Name, baseDir)), hdr.FileInfo().Mode())
			if err != nil {
				log.Fatal(err)
			}
		case tar.TypeReg:
			b, err := ioutil.ReadAll(tr)
			if err != nil {
				log.Fatal(err)
			}

			fileName := dir + filepath.FromSlash("/"+strings.TrimPrefix(hdr.Name, baseDir))

			if goPattern.MatchString(hdr.Name) {
				b = bytes.ReplaceAll(b, []byte(baseReplace), []byte("\""+moduleName))
			} else if goModPattern.MatchString(hdr.Name) {
				b = bytes.ReplaceAll(b, []byte(baseReplaceModule), []byte(moduleNamePrefix))
			} else if gitIgnorePattern.MatchString(hdr.Name) {
				fileName = filepath.Dir(fileName) + filepath.FromSlash("/.gitignore")
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
