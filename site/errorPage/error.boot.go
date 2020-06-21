package errorPage

import (
	"net/http"

	"github.com/cjtoolkit/ctx/v2/ctxHttp"

	"github.com/cjtoolkit/ignition/shared/utility/httpError"
	"github.com/cjtoolkit/ignition/shared/utility/router"
	"github.com/cjtoolkit/ignition/site/errorPage/controller"
)

func bootError(controller *controller.ErrorController, router router.Router) {
	showError := func(req *http.Request, code int, message string) {
		controller.ShowError(ctxHttp.Context(req), code, http.StatusText(code), message)
	}

	router.SetNotFound(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		showError(request, http.StatusNotFound, "Router could not find path")
	}))
	router.SetMethodNotAllowed(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		showError(request, http.StatusMethodNotAllowed, "Router found path, but however method is not allowed")
	}))

	router.SetPanicHandler(httpError.GetPanicHandler(showError))
}
