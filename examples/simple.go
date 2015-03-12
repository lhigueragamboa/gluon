package main

import (
	"errors"
	"net/http"
	"runtime"

	"github.com/lhigueragamboa/gluon"
)

func name(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hola\n"))
}

func err(w http.ResponseWriter, r *http.Request) *gluon.HandlerError {
	return &gluon.HandlerError{
		Error:      errors.New("nuevo error"),
		Message:    "Esto no funciono",
		StatusCode: http.StatusInternalServerError,
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	mux := gluon.New("/")
	mux.Get("/", name)

	hola := gluon.New("/hola")
	hola.Get("/name/:name", name)
	hola.Get("/error", err)

	gluon.LogHandlerErrors = true

	if err := gluon.Start(":8080"); err != nil {
		gluon.Logger.Fatalln(err)
	}
}
