package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CreateZip(fileName, src string) error {
	src = filepath.FromSlash(src)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(src)
	if err != nil {
		return err
	}

	err, data := walkDirectory("", src, err)
	if err != nil {
		return err
	}

	os.Chdir(wd)

	out, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer out.Close()

	w := zip.NewWriter(out)

	for _, datum := range data {
		if datum.info.IsDir() {
			continue
		}
		f, err := w.CreateHeader(&zip.FileHeader{
			Name:   strings.Trim(strings.Trim(datum.src, "."), fmt.Sprintf("%c", filepath.Separator)),
			Method: zip.Store,
		})
		if err != nil {
			return err
		}
		in, err := os.Open(datum.src)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, in)
		if err != nil {
			return err
		}
		in.Close()
	}

	return w.Close()
}
