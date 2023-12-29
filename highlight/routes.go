package highlight

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sikozonpc/notebase/auth"
	"github.com/sikozonpc/notebase/medium"
	"github.com/sikozonpc/notebase/storage"
	t "github.com/sikozonpc/notebase/types"
	u "github.com/sikozonpc/notebase/utils"
)

type Handler struct {
	store     t.HighlightStore
	userStore t.UserStore
	storage   storage.Storage
	bookStore t.BookStore
	mailer    medium.Medium
}

func NewHandler(
	store t.HighlightStore,
	userStore t.UserStore,
	storage storage.Storage,
	bookStore t.BookStore,
	mailer medium.Medium,
) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
		storage:   storage,
		bookStore: bookStore,
		mailer:    mailer,
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

	router.HandleFunc(
		"/daily-insights",
		auth.WithJWTAuth(u.MakeHTTPHandler(h.handleSendDailyInsights), h.userStore),
	).
		Methods("GET")
}

func (s *Handler) handleSendDailyInsights(w http.ResponseWriter, r *http.Request) error {
	users, err := s.userStore.GetUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		user, err := s.userStore.GetUserByID(u.ID)
		if err != nil {
			return fmt.Errorf("user with id %d not found", u.ID)
		}

		hs, err := s.store.GetRandomHighlights(u.ID, 3)
		if err != nil {
			return err
		}

		// Don't send daily insights if there are none
		if len(hs) == 0 {
			continue
		}

		insights, err := buildInsights(hs, s.bookStore)
		if err != nil {
			return err
		}

		if err = s.mailer.SendInsights(user, insights); err != nil {
			return err
		}
	}

	return u.WriteJSON(w, http.StatusOK, nil)
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

	raw, err := parseKindleExtractFile(file, userID)
	if err != nil {
		return err
	}

	// Create book
	_, err = s.bookStore.GetBookByISBN(raw.ASIN)
	if err != nil {
		s.bookStore.CreateBook(t.Book{
			ISBN:    raw.ASIN,
			Title:   raw.Title,
			Authors: raw.Authors,
		})
	}

	// Create highlights
	hs := make([]t.Highlight, len(raw.Highlights))
	for i, h := range raw.Highlights {
		hs[i] = t.Highlight{
			Text:     h.Text,
			Location: h.Location.URL,
			Note:     h.Note,
			UserID:   userID,
			BookID:   raw.ASIN,
		}
	}

	err = s.store.CreateHighlights(hs)
	if err != nil {
		log.Println("Error creating highlights: ", err)
		return err
	}

	return u.WriteJSON(w, http.StatusOK, raw)
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

	highlight := New(payload.Text, payload.Location, payload.Note, payload.BookId, payload.UserId)

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
	BookId   string `json:"bookId"`
}

func buildInsights(hs []*t.Highlight, bookStore t.BookStore) ([]*t.DailyInsight, error) {
	var insights []*t.DailyInsight

	for _, h := range hs {
		book, err := bookStore.GetBookByISBN(h.BookID)
		if err != nil {
			log.Println("Error getting book: ", err)
			return nil, err
		}

		insights = append(insights, &t.DailyInsight{
			Text:        h.Text,
			Note:        h.Note,
			BookAuthors: book.Authors,
			BookTitle:   book.Title,
		})
	}

	return insights, nil
}
