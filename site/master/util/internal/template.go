//go:generate gobox tools/gmock

package internal

import "io"

type Template interface {
	Execute(wr io.Writer, data interface{}) error
}
