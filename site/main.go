package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/command/param"
	"github.com/cjtoolkit/ignition/shared/utility/router"
	"github.com/cjtoolkit/ignition/site/errorPage"
	"github.com/cjtoolkit/ignition/site/fileServer"
	"github.com/cjtoolkit/ignition/site/homePage"
)

var build = "Undefined"

func boot() (http.Handler, param.Param) {
	context := ctx.NewBackgroundContext()
	defer ctx.ClearBackgroundContext(context)

	_param := param.GetParam(context)

	errorPage.Boot(context)
	fileServer.Boot(context)

	homePage.Boot(context)

	fmt.Printf("Build: %q", build)
	fmt.Println()
	fmt.Println("Bootup up successfully.")
	fmt.Println("")

	return router.GetRouter(context), _param
}

func main() {
	handler, _param := boot()

	param.CheckIfTestRun(_param)

	fmt.Println("Now listening on", _param.Address)
	log.Print(http.ListenAndServe(_param.Address, handler))
}
