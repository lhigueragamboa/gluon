package gluon

import "net/http"

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

func (r *router) do(method, path string, handler interface{}) {
	if path[0] != '/' {
		path = "/" + path
	}

	switch h := handler.(type) {
	case http.Handler:
		defaultRouter.Handle(
			method,
			r.basePath+path,
			httprouterHandler(
				r.middlewareStack.Then(h),
			),
		)
	case func(http.ResponseWriter, *http.Request):
		defaultRouter.Handle(
			method,
			r.basePath+path,
			httprouterHandler(
				r.middlewareStack.Then(
					http.HandlerFunc(h),
				),
			),
		)
	case func(http.ResponseWriter, *http.Request) *HandlerError:
		defaultRouter.Handle(
			method,
			r.basePath+path,
			httprouterHandler(
				r.middlewareStack.Then(
					handlerFuncError(h),
				),
			),
		)
	default:
		Logger.Fatalf("Wrong type for handler: %T\n", h)
	}
}
