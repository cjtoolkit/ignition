package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFolder(dst, src string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(wd + filepath.FromSlash("/"+src))
	defer os.Chdir(wd)
	if err != nil {
		return err
	}
	dstPath := wd + filepath.FromSlash("/"+dst)

	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if path == "." {
			return nil
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			fmt.Printf("Creating: %s -> %s", path, dst+filepath.FromSlash("/"+path))
			fmt.Println()
			return os.Mkdir(dstPath+filepath.FromSlash("/"+path), info.Mode())
		}
		fmt.Printf("Copying: %s -> %s", path, dst+filepath.FromSlash("/"+path))
		fmt.Println()
		return copyFileContents(dstPath+filepath.FromSlash("/"+path), path)
	})
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
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
	if err = os.Link(src, dst); err == nil {
		return
	}
	fmt.Printf("Copying: %s -> %s", src, dst)
	err = copyFileContents(dst, src)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(dst, src string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	stat, err := in.Stat()
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
