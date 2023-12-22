package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	t "github.com/sikozonpc/notebase/types"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Wraps a handler function by handling any errors that occur
func MakeHTTPHandler(fn t.EndpointHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, t.APIError{
				Error: err.Error(),
			})
		}
	}
}

func GetIdFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("no id in request")
	}

	return strconv.Atoi(id)
}
