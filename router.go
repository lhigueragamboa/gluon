package gluon

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type router struct {
	basePath        string
	middlewareStack alice.Chain
	router          *httprouter.Router
}

func New(path string) *router {
	return &router{
		basePath:        path,
		middlewareStack: alice.New(),
		router:          httprouter.New(),
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *router) LogHandlerErrors() {
	Logger = log.New(os.Stderr, "gluon: ", log.Lshortfile)
}
