package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetIdFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("no id in request")
	}

	return strconv.Atoi(id)
}


func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
