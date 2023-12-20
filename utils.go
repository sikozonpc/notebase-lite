package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	t "github.com/sikozonpc/notebase/types"
	"github.com/sikozonpc/notebase/utils"
)

func NewAPIServer(addr string, store Storage) *APIServer {
	return &APIServer{
		addr:  addr,
		store: store,
	}
}

// Wraps a handler function by handling any errors that occur
func makeHTTPHandler(fn EndpointHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, t.APIError{
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
