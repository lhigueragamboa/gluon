package main

import (
	"errors"
	"log"
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
	hola := gluon.New("/hola")
	hola.Get("/name/:name", name)
	hola.Get("/error", err)
	if gluon.ErrSetingRoutes != nil {
		log.Fatalln(gluon.ErrSetingRoutes)
	}

	hola.LogHandlerErrors()

	log.Fatalln(http.ListenAndServe(":8080", hola))
}
