package param

import (
	"flag"
	"os"

	"github.com/cjtoolkit/ctx/v2"
)

type Param struct {
	Address    string
	Production bool
	TestRun    bool
}

func GetParam(context ctx.Context) Param {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return initParam(), nil
	}).(Param)
}

func initParam() Param {
	flagSet := flag.NewFlagSet("main", flag.ExitOnError)

	showHelp := flagSet.Bool("help", false, "Show Help")

	address := flagSet.String("address", ":8080", "Set server listening address.")
	production := flagSet.Bool("prod", false, "Set to production mode")
	testRun := flagSet.Bool("test-run", false, "Set to test run, does not start server")

	_ = flagSet.Parse(os.Args[1:])

	if *showHelp {
		flagSet.PrintDefaults()
	}

	param := &Param{
		Address:    *address,
		Production: *production,
		TestRun:    *testRun,
	}

	return *param
}

func CheckIfTestRun(param Param) {
	if param.TestRun {
		os.Exit(0)
	}
}
