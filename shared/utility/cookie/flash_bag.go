//go:generate gobox tools/gmock

package cookie

import (
	"encoding/json"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/constant"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

const flashBagSession = constant.FlashBagSession

// FlashBagValues maps a string key to a list of values.
type FlashBagValues map[string][]string

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v FlashBagValues) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value. It replaces any existing
// values.
func (v FlashBagValues) Set(key, value string) {
	v[key] = []string{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v FlashBagValues) Add(key, value string) {
	v[key] = append(v[key], value)
}

// Del deletes the values associated with key.
func (v FlashBagValues) Del(key string) {
	delete(v, key)
}

type FlashBag interface {
	GetFlashBag(context ctx.Context) FlashBagValues
	SaveFlashBagToSession(context ctx.Context)
}

func GetFlashBag(context ctx.BackgroundContext) FlashBag {
	type flashBagContext struct{}
	return context.Persist(flashBagContext{}, func() (interface{}, error) {
		return FlashBag(flashBag{
			session:      GetSession(context),
			errorService: loggers.GetErrorService(context),
		}), nil
	}).(FlashBag)
}

type flashBag struct {
	session      Session
	errorService loggers.ErrorService
}

func (f flashBag) GetFlashBag(context ctx.Context) FlashBagValues {
	type flashBagContext struct{}
	return context.PersistData(flashBagContext{}, func() interface{} {
		fB := FlashBagValues{}

		b := f.session.GetDel(context, flashBagSession)
		_ = json.Unmarshal(b, &fB)

		return fB
	}).(FlashBagValues)
}

func (f flashBag) SaveFlashBagToSession(context ctx.Context) {
	fB := f.GetFlashBag(context)
	b, err := json.Marshal(fB)
	f.errorService.CheckErrorAndPanic(err)
	f.session.Set(context, flashBagSession, b)
}
