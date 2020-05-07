package homePage

import (
	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/router"
	"github.com/cjtoolkit/ignition/site/homePage/controller"
	"github.com/cjtoolkit/ignition/site/urls"
)

type homeBoot struct {
	homeController controller.HomeController
	router         router.Router
}

func (b homeBoot) boot() {
	b.router.GET(urls.Index, func(context ctx.Context, _ router.Params) {
		b.homeController.Index(context)
	})
}
