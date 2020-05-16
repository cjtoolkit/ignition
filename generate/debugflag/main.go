//  Prepend '// +build debug' to top of file
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

const flag = "// +build debug"

func main() {
	if len(os.Args) < 1 {
		return
	}

	wd, err := os.Getwd()
	if nil != err {
		return
	}

	sourceFile := wd + "/" + os.Args[1]
	file, err := os.OpenFile(sourceFile, os.O_RDWR, 0666)
	if nil != err {
		return
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	defer buf.Reset()
	io.Copy(buf, file)

	file.Seek(0, io.SeekStart)
	file.Truncate(0)

	fmt.Fprintln(file, flag)
	fmt.Fprintln(file)

	io.Copy(file, buf)
}
