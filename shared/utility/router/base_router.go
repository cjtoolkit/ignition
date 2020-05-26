package router

import (
	"github.com/cjtoolkit/ctx"
	"github.com/julienschmidt/httprouter"
)

type (
	Param  = httprouter.Param
	Params = httprouter.Params
)

func getBaseRouter(context ctx.Context) *httprouter.Router {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		router := httprouter.New()

		return router, nil
	}).(*httprouter.Router)
}
