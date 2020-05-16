package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type data struct {
	ImportPath string
	Name       string
}

func main() {
	cmd := exec.Command(runtime.Version(), "list", "-json")
	cmd.Stderr = os.Stderr
	b, err := cmd.Output()
	if err != nil {
		cmd = exec.Command("go", "list", "-json")
		cmd.Stderr = os.Stderr
		b, err = cmd.Output()
	}
	checkErr(err)

	d := data{}
	checkErr(json.Unmarshal(b, &d))

	goSrcFile := os.Getenv("GOFILE")
	var goDestFile string
	{
		split := strings.Split(goSrcFile, ".")
		split = split[:len(split)-1]
		goDestFile = strings.Join(split, ".") + ".mock.go"
	}

	runCommand("mockgen",
		"-write_package_comment=false",
		fmt.Sprintf("-package=%s", d.Name),
		fmt.Sprintf("-self_package=%s", d.ImportPath),
		fmt.Sprintf("-source=%s", goSrcFile),
		fmt.Sprintf("-destination=%s", goDestFile))

	runCommand("debugflag", goDestFile)
}

func runCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	checkErr(cmd.Run())
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
