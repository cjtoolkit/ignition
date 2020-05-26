//go:generate gobox tools/easymock

package internal

type Context interface {
	Set(key, value interface{})
	Get(key interface{}) (interface{}, bool)

	// The fn function only gets called if there is a cache miss. Return error as nil to bypass health check.
	Persist(key interface{}, fn func() (interface{}, error)) interface{}
}
