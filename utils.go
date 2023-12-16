package main

import (
	"encoding/json"
	"net/http"
)

func NewAPIServer(addr string, store Storage) *APIServer {
	return &APIServer{
		addr:  addr,
		store: store,
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Wraps a handler function by handling any errors that occur
func makeHTTPHandler(fn EndpointHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, APIError{
				Error: err.Error(),
			})
		}
	}
}
