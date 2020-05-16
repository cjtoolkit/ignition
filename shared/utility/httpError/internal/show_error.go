//go:generate gobox tools/easymock

package internal

import "net/http"

type ShowError interface {
	ShowError(req *http.Request, code int, message string)
}
