package controller

import (
	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/site/homePage/model"
	"github.com/cjtoolkit/ignition/site/homePage/view"
)

type HomeController struct {
	view view.HomeView
}

func GetHomeController(context ctx.Context) *HomeController {
	return &HomeController{view: view.NewHomeView(context)}
}

func (h *HomeController) Index(context ctx.Context) {
	h.view.ExecIndexView(context, model.Index{Message: "Hello World"})
}
