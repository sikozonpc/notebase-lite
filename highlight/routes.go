package highlight

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sikozonpc/notebase/auth"
	"github.com/sikozonpc/notebase/storage"
	t "github.com/sikozonpc/notebase/types"
	u "github.com/sikozonpc/notebase/utils"
)

type Handler struct {
	store     t.HighlightStore
	userStore t.UserStore
	storage   storage.Storage
}

func NewHandler(store t.HighlightStore, userStore t.UserStore, storage storage.Storage) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
		storage:   storage,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(
		"/user/{userID}/highlight",
		auth.WithJWTAuth(u.MakeHTTPHandler(h.handleGetUserHighlights), h.userStore),
	).Methods("GET")

	router.HandleFunc(
		"/user/{userID}/highlight",
		auth.WithJWTAuth(u.MakeHTTPHandler(h.handleCreateHighlight), h.userStore),
	).Methods("POST")

	router.HandleFunc(
		"/user/{userID}/highlight/{id}",
		auth.WithJWTAuth(u.MakeHTTPHandler(h.handleGetHighlightByID), h.userStore),
	).Methods("GET")

	router.HandleFunc(
		"/user/{userID}/highlight/{id}",
		auth.WithJWTAuth(u.MakeHTTPHandler(h.handleDeleteHighlight), h.userStore),
	).Methods("DELETE")

	router.HandleFunc(
		"/user/{userID}/parse-kindle-extract",
		auth.WithJWTAuth(u.MakeHTTPHandler(h.handleParseKindleFile), h.userStore),
	).
		Methods("POST")
}

func (s *Handler) handleParseKindleFile(w http.ResponseWriter, r *http.Request) error {
	userID, err := u.GetParamFromRequest(r, "userID")
	if err != nil {
		return err
	}

	query := r.URL.Query()
	filename := query.Get("filename")

	if filename == "" {
		return u.WriteJSON(w, http.StatusBadRequest, fmt.Errorf("filename is required"))
	}

	file, err := s.storage.Read(filename)
	if err != nil {
		return u.WriteJSON(w, http.StatusInternalServerError, err)
	}

	hs, err := parseKindleExtractFile(file, userID)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, hs)
}

func (s *Handler) handleGetUserHighlights(w http.ResponseWriter, r *http.Request) error {
	userID, err := u.GetParamFromRequest(r, "userID")
	if err != nil {
		return err
	}

	hs, err := s.store.GetUserHighlights(userID)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, hs)
}

func (s *Handler) handleDeleteHighlight(w http.ResponseWriter, r *http.Request) error {
	id, err := u.GetParamFromRequest(r, "id")
	if err != nil {
		return err
	}

	err = s.store.DeleteHighlight(id)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, nil)
}

func (s *Handler) handleGetHighlightByID(w http.ResponseWriter, r *http.Request) error {
	userID, err := u.GetParamFromRequest(r, "userID")
	if err != nil {
		return err
	}

	id, err := u.GetParamFromRequest(r, "id")
	if err != nil {
		return err
	}

	h, err := s.store.GetHighlightByID(id, userID)
	if err != nil {
		return err
	}

	if h == nil {
		return u.WriteJSON(w, http.StatusNotFound, t.APIError{Error: fmt.Errorf("highlight with id %d not found", id).Error()})
	}

	return u.WriteJSON(w, http.StatusOK, h)

}

func (s *Handler) handleCreateHighlight(w http.ResponseWriter, r *http.Request) error {
	payload := new(CreateHighlightRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return err
	}

	highlight := New(payload.Text, payload.Location, payload.Note, payload.UserId, payload.BookId)

	if err := s.store.CreateHighlight(*highlight); err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, highlight)

}

type CreateHighlightRequest struct {
	Text     string `json:"text"`
	Location string `json:"location"`
	Note     string `json:"note"`
	UserId   int    `json:"userId"`
	BookId   int    `json:"bookId"`
}
