package environment

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/constant"
)

const (
	configDirectoryVar     = constant.ConfigDirectoryEnvVar
	defaultConfigDirectory = constant.DefaultConfigDirectory
)

type Environment struct {
	ConfigDirectoryVar     string
	DefaultConfigDirectory string
}

func GetEnvironment(context ctx.BackgroundContext) *Environment {
	type c struct{}
	return context.Persist(c{}, func() (i interface{}, e error) {
		return initEnvironment(), nil
	}).(*Environment)
}

func initEnvironment() *Environment {
	return &Environment{
		ConfigDirectoryVar:     configDirectoryVar,
		DefaultConfigDirectory: defaultConfigDirectory,
	}
}

func (e *Environment) ParseConfigDirectory() string {
	configDirectory := os.Getenv(e.ConfigDirectoryVar)
	if configDirectory != "" {
		return filepath.FromSlash(configDirectory)
	}

	if runtime.GOOS == "windows" {
		log.Fatalf("'%s' is mandatory on Windows", e.ConfigDirectoryVar)
	}

	return filepath.FromSlash(e.DefaultConfigDirectory)
}
