package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/sikozonpc/notebase/auth"
	"github.com/sikozonpc/notebase/data"
	t "github.com/sikozonpc/notebase/types"
	u "github.com/sikozonpc/notebase/utils"
)

type APIServer struct {
	addr  string
	store Storage
}

type EndpointHandler func(w http.ResponseWriter, r *http.Request) error

func (s *APIServer) Run() error {
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	router.HandleFunc("/register", makeHTTPHandler(s.handleRegister))
	router.HandleFunc("/login", makeHTTPHandler(s.handleLogin))

	router.HandleFunc("/highlight", auth.WithJWTAuth(makeHTTPHandler(s.handleHighlights)))
	router.HandleFunc("/highlight/{id}", auth.WithJWTAuth(makeHTTPHandler(s.handleHighlightsById)))

	log.Println("Listening on", s.addr)

	log.Println("Process PID", os.Getpid())

	env := Configs.Env
	if env == "development" {
		v := reflect.ValueOf(Configs)

		for i := 0; i < v.NumField(); i++ {
			log.Println(v.Type().Field(i).Name, "=", v.Field(i).Interface())
		}
	}

	return http.ListenAndServe(s.addr, router)
}

func (s *APIServer) handleHighlights(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetHighlights(w, r)
	}

	if r.Method == http.MethodPost {
		return s.handleCreateHighlight(w, r)
	}

	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *APIServer) handleHighlightsById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetHighlight(w, r)
	}

	if r.Method == http.MethodDelete {
		return s.handleDeleteHighlight(w, r)
	}

	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *APIServer) handleGetHighlights(w http.ResponseWriter, r *http.Request) error {
	highlights, err := s.store.GetHighlights()
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, highlights)
}

func (s *APIServer) handleDeleteHighlight(w http.ResponseWriter, r *http.Request) error {
	id, err := getIdFromRequest(r)
	if err != nil {
		return err
	}

	err = s.store.DeleteHighlight(id)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, nil)
}

func (s *APIServer) handleGetHighlight(w http.ResponseWriter, r *http.Request) error {
	id, err := getIdFromRequest(r)
	if err != nil {
		return err
	}

	h, err := s.store.GetHighlightByID(id)
	if err != nil {
		return err
	}

	if h.ID == 0 {
		return u.WriteJSON(w, http.StatusNotFound, t.APIError{Error: fmt.Errorf("highlight with id %d not found", id).Error()})
	}

	return u.WriteJSON(w, http.StatusOK, h)

}

func (s *APIServer) handleCreateHighlight(w http.ResponseWriter, r *http.Request) error {
	payload := new(t.CreateHighlightRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return err
	}

	highlight := data.NewHighlight(payload.Text, payload.Location, payload.Note, payload.UserId, payload.BookId)

	if err := s.store.CreateHighlight(*highlight); err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, highlight)

}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("method %s not allowed", r.Method)
	}

	payload := new(t.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return err
	}

	user, err := s.store.GetUserByEmail(payload.Email)
	if err != nil {
		return err
	}

	if data.ComparePasswords(user.Password, []byte(payload.Password)) {
		return fmt.Errorf("invalid password or user does not exist")
	}

	token, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, token)
}

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("method %s not allowed", r.Method)
	}

	payload := new(t.RegisterRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return err
	}

	hashedPassword, err := data.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	user := data.NewUser(payload.FirstName, payload.LastName, payload.Email, hashedPassword)

	if err := s.store.CreateUser(*user); err != nil {
		return err
	}

	token, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, token)
}

func createAndSetAuthCookie(userID int, w http.ResponseWriter) (string, error) {
	secret := []byte(Configs.JWTSecret)
	token, err := auth.CreateJWT(secret, userID)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
