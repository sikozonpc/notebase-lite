package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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


func getIdFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("no id in request")
	}

	return strconv.Atoi(id)
}