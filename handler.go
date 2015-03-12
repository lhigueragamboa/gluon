package gluon

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type handlerFuncError func(http.ResponseWriter, *http.Request) *HandlerError

type HandlerError struct {
	Error      error
	Message    string
	StatusCode int
}

type errResponse struct {
	Message string `json:"error"`
}

func (h handlerFuncError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if hErr := h(w, r); hErr != nil {
		if LogHandlerErrors {
			go func() { Logger.Println(hErr.Error) }()
		}
		writeErrorJSON(w, hErr)
	}
}

func writeErrorJSON(w http.ResponseWriter, hErr *HandlerError) error {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(hErr.StatusCode)

	return json.NewEncoder(w).Encode(&errResponse{hErr.Message})
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
