package internal

import (
	"archive/zip"
	"io"
	"os"
	"time"
)

func CreateZip(fileName, src string) error {
	err, data := walkDirectory("", src)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer out.Close()

	w := zip.NewWriter(out)

	now := time.Now()

	for _, datum := range data {
		if datum.info.IsDir() {
			continue
		}
		f, err := w.CreateHeader(&zip.FileHeader{
			Name:     datum.dest,
			Method:   zip.Store,
			Modified: now,
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
