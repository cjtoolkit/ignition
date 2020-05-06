//go:generate gobox tools/gmock

package cookie

import (
	"encoding/json"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/constant"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

const flashBagSession = constant.FlashBagSession

type FlashBag interface {
	GetFlashBag(context ctx.Context) map[string]string
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

func (f flashBag) GetFlashBag(context ctx.Context) map[string]string {
	type flashBagContext struct{}
	return context.PersistData(flashBagContext{}, func() interface{} {
		fB := map[string]string{}

		b := f.session.GetDel(context, flashBagSession)
		_ = json.Unmarshal(b, &fB)

		return fB
	}).(map[string]string)
}

func (f flashBag) SaveFlashBagToSession(context ctx.Context) {
	fB := f.GetFlashBag(context)
	b, err := json.Marshal(fB)
	f.errorService.CheckErrorAndPanic(err)
	f.session.Set(context, flashBagSession, b)
}
