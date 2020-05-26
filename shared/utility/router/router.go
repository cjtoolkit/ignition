package router

import (
	"log"
	"net/http"

	"github.com/cjtoolkit/ctx/ctxHttp"

	"github.com/cjtoolkit/ignition/shared/utility/command/param"

	"github.com/cjtoolkit/ctx"
	"github.com/julienschmidt/httprouter"
)

type Handle func(context ctx.Context, params Params)

type Router interface {
	DELETE(path string, handle Handle)
	GET(path string, handle Handle)
	HEAD(path string, handle Handle)
	Handle(method, path string, handle Handle)
	Handler(method, path string, handler http.Handler)
	HandlerFunc(method, path string, handler http.HandlerFunc)
	Lookup(method, path string) (httprouter.Handle, Params, bool)
	OPTIONS(path string, handle Handle)
	PATCH(path string, handle Handle)
	POST(path string, handle Handle)
	PUT(path string, handle Handle)
	ServeFiles(path string, root http.FileSystem)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	SetNotFound(h http.Handler)
	SetMethodNotAllowed(h http.Handler)
	SetPanicHandler(f func(http.ResponseWriter, *http.Request, interface{}))
}

func GetRouter(context ctx.Context) Router {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return initRouter(context), nil
	}).(Router)
}

type router struct {
	router     *httprouter.Router
	production bool
}

func initRouter(context ctx.Context) *router {
	return &router{
		router:     getBaseRouter(context),
		production: param.GetParam(context).Production,
	}
}

func (r *router) DELETE(path string, handle Handle)  { r.Handle(http.MethodDelete, path, handle) }
func (r *router) GET(path string, handle Handle)     { r.Handle(http.MethodGet, path, handle) }
func (r *router) HEAD(path string, handle Handle)    { r.Handle(http.MethodHead, path, handle) }
func (r *router) OPTIONS(path string, handle Handle) { r.Handle(http.MethodOptions, path, handle) }
func (r *router) PATCH(path string, handle Handle)   { r.Handle(http.MethodPatch, path, handle) }
func (r *router) POST(path string, handle Handle)    { r.Handle(http.MethodPost, path, handle) }
func (r *router) PUT(path string, handle Handle)     { r.Handle(http.MethodPut, path, handle) }
func (r *router) SetNotFound(h http.Handler)         { r.router.NotFound = h }
func (r *router) SetMethodNotAllowed(h http.Handler) { r.router.MethodNotAllowed = h }

func (r *router) Handle(method, path string, handle Handle) {
	r.router.Handle(method, path, func(_ http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handle(ctxHttp.Context(request), params)
	})
}

func (r *router) Handler(method, path string, handler http.Handler) {
	r.router.Handler(method, path, handler)
}

func (r *router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.router.HandlerFunc(method, path, handler)
}

func (r *router) ServeFiles(path string, root http.FileSystem) { r.router.ServeFiles(path, root) }

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.production {
		log.Printf("HTTP: %q: %q", req.Method, req.URL.String())
	}
	req = ctxHttp.NewContext(req, w)
	r.router.ServeHTTP(w, req)
}

func (r *router) SetPanicHandler(f func(http.ResponseWriter, *http.Request, interface{})) {
	r.router.PanicHandler = f
}

func (r *router) Lookup(method, path string) (httprouter.Handle, Params, bool) {
	return r.router.Lookup(method, path)
}
