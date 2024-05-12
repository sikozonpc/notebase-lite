package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func GetStringParamFromRequest(r *http.Request, param string) (string, error) {
	vars := mux.Vars(r)
	str, ok := vars[param]
	if !ok {
		log.Printf("no param: %v in request", param)
		return "", fmt.Errorf("no param: %v in request", param)
	}

	return str, nil
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}