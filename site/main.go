package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cjtoolkit/ignition/site/errorPage"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/command/param"
	"github.com/cjtoolkit/ignition/shared/utility/router"
)

func boot() (http.Handler, param.Param) {
	context := ctx.NewBackgroundContext()
	defer ctx.ClearBackgroundContext(context)

	_param := param.GetParam(context)

	errorPage.Boot(context)

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
