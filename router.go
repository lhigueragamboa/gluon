package gluon

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type router struct {
	basePath        string
	middlewareStack alice.Chain
	router          *httprouter.Router
}

func New(path string) *router {
	if path == "/" {
		path = ""
	}

	if last := len(path); path != "" && path[last] == '/' {
		path = path[:last-1]
	}

	return &router{
		basePath:        path,
		middlewareStack: alice.New(),
		router:          httprouter.New(),
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
