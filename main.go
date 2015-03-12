package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type handlerError struct {
	Err        error
	Message    string
	StatusCode int
}

type handlerFuncError func(http.ResponseWriter, *http.Request) *handlerError

type errResponse struct {
	Message string `json:"error"`
}

func writeErrorJSON(w http.ResponseWriter, hErr *handlerError) error {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(hErr.StatusCode)

	return json.NewEncoder(w).Encode(&errResponse{hErr.Message})
}

func (h handlerFuncError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if hErr := h(w, r); hErr != nil {
		writeErrorJSON(w, hErr)
	}
}

func httprouterHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		q := r.URL.Query()
		for i := range ps {
			q.Set(":"+ps[i].Key, ps[i].Value)
		}
		r.URL.RawQuery = q.Encode()

		h.ServeHTTP(w, r)
	}
}

type router struct {
	basePath        string
	middlewareStack alice.Chain
	router          *httprouter.Router
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
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

func (r *router) Use(middleware ...alice.Constructor) {
	r.middlewareStack = r.middlewareStack.Append(middleware...)
}

var errWrongHandlerType = errors.New("Wrong type of handler parameter")
var ErrSetingRoutes error = nil

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
	case func(http.ResponseWriter, *http.Request) *handlerError:
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
		ErrSetingRoutes = errWrongHandlerType
	}
}

func (r *router) Get(path string, handler interface{}) {
	r.do("GET", path, handler)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello\n"))
}

func hello(w http.ResponseWriter, r *http.Request) *handlerError {
	return &handlerError{
		Err:        errors.New("XXX"),
		Message:    "fake error",
		StatusCode: http.StatusInternalServerError,
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	r := New("/")

	r.Get("/", index)
	r.Get("/hola/:name", hello)
	if ErrSetingRoutes != nil {
		log.Fatalln(ErrSetingRoutes)
	}

	log.Fatalln(http.ListenAndServe(":8080", r))
}
