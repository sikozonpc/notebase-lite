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

func GetParamFromRequest(r *http.Request, param string) (int, error) {
	vars := mux.Vars(r)
	id, ok := vars[param]
	if !ok {
		return 0, fmt.Errorf("no param: %v in request", param)
	}

	return strconv.Atoi(id)
}
