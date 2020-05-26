package homePage

import (
	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/router"
	"github.com/cjtoolkit/ignition/site/homePage/controller"
)

func Boot(context ctx.Context) {
	homeBoot{
		homeController: controller.GetHomeController(context),
		router:         router.GetRouter(context),
	}.boot()
}
