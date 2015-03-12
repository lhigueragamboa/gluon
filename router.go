package gluon

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var defaultRouter = httprouter.New()

type router struct {
	basePath        string
	middlewareStack alice.Chain
}

func New(path string) *router {
	path = strings.TrimRight(path, "/")

	return &router{
		basePath:        path,
		middlewareStack: alice.New(),
	}
}

func Start(addr string) error {
	return http.ListenAndServe(addr, defaultRouter)
}
