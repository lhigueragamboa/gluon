package gluon

import (
	"fmt"
	"net/http"
)

var ErrSetingRoutes error

func (r *router) Get(path string, handler interface{}) {
	ErrSetingRoutes = r.do("GET", path, handler)
}

func (r *router) Post(path string, handler interface{}) {
	ErrSetingRoutes = r.do("POST", path, handler)
}

func (r *router) Put(path string, handler interface{}) {
	ErrSetingRoutes = r.do("PUT", path, handler)
}

func (r *router) Patch(path string, handler interface{}) {
	ErrSetingRoutes = r.do("PATCH", path, handler)
}

func (r *router) Head(path string, handler interface{}) {
	ErrSetingRoutes = r.do("HEAD", path, handler)
}

func (r *router) Delete(path string, handler interface{}) {
	ErrSetingRoutes = r.do("DELETE", path, handler)
}

func (r *router) Options(path string, handler interface{}) {
	ErrSetingRoutes = r.do("OPTIONS", path, handler)
}

func (r *router) Trace(path string, handler interface{}) {
	ErrSetingRoutes = r.do("TRACE", path, handler)
}

func (r *router) do(method, path string, handler interface{}) error {
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
		return fmt.Errorf("Wrong type for handler: %T", h)
	}

	return nil
}
