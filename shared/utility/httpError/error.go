package httpError

import (
	"fmt"
	"net/http"
)

type NoError struct{}

func Halt() {
	panic(NoError{})
}

type HttpError struct {
	Code    int
	Message string
}

func HaltNotFound(message string) {
	panic(HttpError{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

func HaltForbidden(message string) {
	panic(HttpError{
		Code:    http.StatusForbidden,
		Message: message,
	})
}

func HaltInternalServerError(message string) {
	panic(HttpError{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

func HaltCustomError(code int, message string) {
	panic(HttpError{
		Code:    code,
		Message: message,
	})
}

type HttpRedirectError struct {
	Code     int
	Location string
}

func HaltMovedPermanently(location string) {
	panic(HttpRedirectError{
		Code:     http.StatusMovedPermanently,
		Location: location,
	})
}

func HaltSeeOther(location string) {
	panic(HttpRedirectError{
		Code:     http.StatusSeeOther,
		Location: location,
	})
}

type HttpErrorNoContent struct {
	Code int
}

func HaltNotModified() {
	panic(HttpErrorNoContent{
		Code: http.StatusNotModified,
	})
}

func CheckParamErr(err error) {
	if err != nil {
		HaltNotFound(fmt.Sprintf("Url Param failed Validation, %v", err))
	}
}
