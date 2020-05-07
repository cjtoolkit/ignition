package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type copyData struct {
	dest string
	src  string
	info os.FileInfo
}

func CopyFolder(dst, src string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(filepath.FromSlash(src))
	if err != nil {
		return err
	}

	err, data := walkDirectory(dst, src, err)
	if err != nil {
		return err
	}

	os.Chdir(wd)
	err = os.Mkdir(dst, 0755)
	if err != nil {
		return err
	}

	for _, datum := range data {
		if datum.info.IsDir() {
			fmt.Printf("Creating: %s -> %s", datum.src, datum.dest)
			fmt.Println()
			err = os.Mkdir(datum.dest, datum.info.Mode())
			if err != nil {
				return err
			}
		}
		err = CopyFile(datum.dest, datum.src)
		if err != nil {
			return err
		}
	}
	return nil
}

func walkDirectory(dst string, src string, err error) (error, []copyData) {
	data := []copyData{}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if path == "." {
			return nil
		}
		if err != nil {
			return err
		}
		data = append(data, copyData{
			dest: dst + filepath.FromSlash("/"+path),
			src:  src + filepath.FromSlash("/"+path),
			info: info,
		})
		return nil
	})
	return err, data
}

func CopyFile(dst, src string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	fmt.Printf("Copying: %s -> %s", src, dst)
	fmt.Println()
	err = copyFileContents(dst, src, sfi)
	return
}

func copyFileContents(dst, src string, stat os.FileInfo) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, stat.Mode())
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
