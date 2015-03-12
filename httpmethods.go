package gluon

import (
	"fmt"
	"net/http"
)

func (r *router) Get(path string, handler interface{}) {
	r.do("GET", path, handler)
}

func (r *router) Post(path string, handler interface{}) {
	r.do("POST", path, handler)
}

func (r *router) Put(path string, handler interface{}) {
	r.do("PUT", path, handler)
}

func (r *router) Patch(path string, handler interface{}) {
	r.do("PATCH", path, handler)
}

func (r *router) Head(path string, handler interface{}) {
	r.do("HEAD", path, handler)
}

func (r *router) Delete(path string, handler interface{}) {
	r.do("DELETE", path, handler)
}

func (r *router) Options(path string, handler interface{}) {
	r.do("OPTIONS", path, handler)
}

func (r *router) Trace(path string, handler interface{}) {
	r.do("TRACE", path, handler)
}

var ErrSetingRoutes error

func (r *router) do(method, path string, handler interface{}) {
	switch h := handler.(type) {
	case http.Handler:
		r.router.Handle(
			method,
			r.basePath+path,
			httprouterHandler(
				r.middlewareStack.Then(h),
			),
		)
	case func(http.ResponseWriter, *http.Request):
		r.router.Handle(
			method,
			r.basePath+path,
			httprouterHandler(
				r.middlewareStack.Then(
					http.HandlerFunc(h),
				),
			),
		)
	case func(http.ResponseWriter, *http.Request) *HandlerError:
		r.router.Handle(
			method,
			r.basePath+path,
			httprouterHandler(
				r.middlewareStack.Then(
					handlerFuncError(h),
				),
			),
		)
	default:
		ErrSetingRoutes = fmt.Errorf("Wrong type for handler: %T", h)
	}
}
