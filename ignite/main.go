package main

import (
	"log"
	"os"

	"github.com/cjtoolkit/ignition/ignite/internal"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You must at least enter a option")
	}

	switch os.Args[1] {
	case "base":
		if len(os.Args) < 4 {
			log.Fatal("You must at least enter three parameters (dir, moduleName)")
		}
		internal.BuildBase(os.Args[2], os.Args[3])
		return
	case "app":
		if len(os.Args) < 5 {
			log.Fatal("You must at least enter four parameters (dir, moduleName, baseModuleName)")
		}
		internal.BuildApp(os.Args[2], os.Args[3], os.Args[4])
		return
	}

	log.Fatal("Unknown option selected")
}
