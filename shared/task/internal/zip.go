package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func CreateZip(fileName, src string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(filepath.FromSlash(src))
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
			Name:   datum.src,
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
	}

	return w.Close()
}
