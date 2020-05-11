package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type copyData struct {
	dest string
	src  string
	info os.FileInfo
}

func CopyFolder(dst, src string) error {
	err, data := walkDirectory(dst, src)
	if err != nil {
		return err
	}

	dst = filepath.FromSlash(dst)
	exist, err := exists(dst)
	if err != nil {
		return err
	}

	if !exist {
		err = os.Mkdir(filepath.FromSlash(dst), 0755)
		if err != nil {
			return err
		}
	}

	for _, datum := range data {
		if datum.info.IsDir() {
			fmt.Printf("Creating: %q -> %q", datum.src, datum.dest)
			fmt.Println()
			exist, err := exists(datum.dest)
			if err != nil {
				return err
			}

			if !exist {
				err := os.Mkdir(datum.dest, datum.info.Mode())
				if err != nil {
					return err
				}
			}
		} else {
			err := CopyFile(datum.dest, datum.src)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func walkDirectory(dst string, src string) (error, []copyData) {
	dst = filepath.FromSlash(dst)
	src = filepath.FromSlash(src)
	data := []copyData{}

	prefixToTrim := src + filepath.FromSlash("/")
	if src == "." {
		prefixToTrim = filepath.FromSlash("./")
	}

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if path == src {
			return nil
		}
		if err != nil {
			return err
		}
		_dst := dst + filepath.FromSlash("/") + strings.TrimPrefix(path, prefixToTrim)
		if dst == "" {
			_dst = _dst[1:]
		}
		data = append(data, copyData{
			dest: _dst,
			src:  path,
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
	fmt.Printf("Copying: %q -> %q", src, dst)
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
