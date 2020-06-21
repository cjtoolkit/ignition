package homePage

import (
	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/router"
	"github.com/cjtoolkit/ignition/site/homePage/controller"
	"github.com/cjtoolkit/ignition/site/urls"
)

type homeBoot struct {
	homeController *controller.HomeController
	router         router.Router
}

func (b homeBoot) boot() {
	b.router.GET(urls.HomeIndex, func(context ctx.Context) {
		b.homeController.Index(context)
	})
}
