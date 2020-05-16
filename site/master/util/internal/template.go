//go:generate gobox tools/easymock

package internal

import "io"

type Template interface {
	Execute(wr io.Writer, data interface{}) error
}
