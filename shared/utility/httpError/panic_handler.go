package httpError

import (
	"fmt"
	"net/http"
)

type (
	ShowError    func(req *http.Request, code int, message string)
	PanicHandler func(res http.ResponseWriter, req *http.Request, i interface{})
)

func GetPanicHandler(showError ShowError) PanicHandler {
	return func(res http.ResponseWriter, req *http.Request, i interface{}) {
		switch i := i.(type) {
		case NoError:
			// Do nothing
		case HttpError:
			showError(req, i.Code, i.Message)
		case HttpErrorNoContent:
			res.WriteHeader(i.Code)
		case HttpRedirectError:
			res.Header().Set("Location", i.Location)
			res.WriteHeader(i.Code)
		default:
			showError(req, http.StatusInternalServerError, fmt.Sprint(i))
		}
	}
}
