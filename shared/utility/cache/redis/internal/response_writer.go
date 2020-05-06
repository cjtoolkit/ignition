//go:generate gobox tools/gmock

package internal

import "net/http"

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
