package errorPage

import (
	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/router"
	"github.com/cjtoolkit/ignition/site/errorPage/controller"
)

func Boot(context ctx.Context) {
	bootError(controller.NewErrorController(context), router.GetRouter(context))
}
