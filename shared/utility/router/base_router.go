package router

import (
	"github.com/cjtoolkit/ctx"
	"github.com/julienschmidt/httprouter"
)

type (
	Param  = httprouter.Param
	Params = httprouter.Params
)

func getBaseRouter(context ctx.BackgroundContext) *httprouter.Router {
	type routerContext struct{}
	return context.Persist(routerContext{}, func() (interface{}, error) {
		router := httprouter.New()

		return router, nil
	}).(*httprouter.Router)
}
