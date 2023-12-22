package highlight

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sikozonpc/notebase/auth"
	t "github.com/sikozonpc/notebase/types"
	u "github.com/sikozonpc/notebase/utils"
)

type Handler struct {
	store     t.HighlightStore
	userStore t.UserStore
}

func NewHandler(store t.HighlightStore, userStore t.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user/{id}/highlight", auth.WithJWTAuth(u.MakeHTTPHandler(h.handleUserHighlights), h.userStore))
	router.HandleFunc("/highlight/{id}", auth.WithJWTAuth(u.MakeHTTPHandler(h.handleHighlightsById), h.userStore))
}

func (s *Handler) handleUserHighlights(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetUserHighlights(w, r)
	}

	if r.Method == http.MethodPost {
		return s.handleCreateHighlight(w, r)
	}

	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *Handler) handleGetUserHighlights(w http.ResponseWriter, r *http.Request) error {
	id, err := u.GetIdFromRequest(r)
	if err != nil {
		return err
	}

	hs, err := s.store.GetUserHighlights(id)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, hs)
}

func (s *Handler) handleHighlightsById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetHighlight(w, r)
	}

	if r.Method == http.MethodDelete {
		return s.handleDeleteHighlight(w, r)
	}

	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *Handler) handleDeleteHighlight(w http.ResponseWriter, r *http.Request) error {
	id, err := u.GetIdFromRequest(r)
	if err != nil {
		return err
	}

	err = s.store.DeleteHighlight(id)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, nil)
}

func (s *Handler) handleGetHighlight(w http.ResponseWriter, r *http.Request) error {
	id, err := u.GetIdFromRequest(r)
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

func (s *Handler) handleCreateHighlight(w http.ResponseWriter, r *http.Request) error {
	payload := new(t.CreateHighlightRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return err
	}

	highlight := New(payload.Text, payload.Location, payload.Note, payload.UserId, payload.BookId)

	if err := s.store.CreateHighlight(*highlight); err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, highlight)

}
