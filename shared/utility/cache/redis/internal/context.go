//go:generate gobox tools/easymock

package internal

import (
	"context"
	"net/http"
)

type Context interface {
	Title() string
	SetTitle(title string)
	Data(key interface{}) interface{}
	SetData(key interface{}, value interface{})
	PersistData(key interface{}, fn func() interface{}) interface{}
	Dep(key interface{}) interface{}
	SetDep(key interface{}, value interface{})
	PersistDep(key interface{}, fn func() interface{}) interface{}
	Ctx() context.Context
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}
