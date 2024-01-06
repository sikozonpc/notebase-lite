package highlight

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
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
		"/cloud/parse-kindle-extract/{fileName}",
		auth.WithAPIKey(u.MakeHTTPHandler(h.handleCloudKindleParse)),
	).
		Methods("POST")

	router.HandleFunc(
		"/cloud/daily-insights",
		auth.WithAPIKey(u.MakeHTTPHandler(h.handleSendDailyInsights)),
	).
		Methods("GET")

	router.HandleFunc(
		"/unsubscribe",
		auth.WithJWTAuth(u.MakeHTTPHandler(h.handleUnsubscribe), h.userStore),
	).
		Methods("GET")
}

func (s *Handler) handleUnsubscribe(w http.ResponseWriter, r *http.Request) error {
	token := u.GetTokenFromRequest(r)

	userID, err := auth.GetUserFromToken(token)
	if err != nil {
		return err
	}

	user, err := s.userStore.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.IsActive = false
	if err := s.userStore.UpdateUser(*user); err != nil {
		return err
	}

	log.Printf("User %s unsubscribed", user.Email)

	return u.WriteJSON(w, http.StatusOK, nil)
}

func (s *Handler) handleSendDailyInsights(w http.ResponseWriter, r *http.Request) error {
	authToken := u.GetTokenFromRequest(r)

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

		if err = s.mailer.SendInsights(user, insights, authToken); err != nil {
			return err
		}
	}

	return u.WriteJSON(w, http.StatusOK, nil)
}

func (s *Handler) handleCloudKindleParse(w http.ResponseWriter, r *http.Request) error {
	userID, err := u.GetParamFromRequest(r, "userID")
	if err != nil {
		return err
	}

	filename, err := u.GetStringParamFromRequest(r, "fileName")
	if err != nil {
		return u.WriteJSON(w, http.StatusBadRequest, fmt.Errorf("filename is required"))
	}

	file, err := s.storage.Read(filename)
	if err != nil {
		return u.WriteJSON(w, http.StatusInternalServerError, err)
	}

	raw, err := parseKindleExtractFromString(file)
	if err != nil {
		return err
	}

	err = s.createDataFromRawBook(raw, userID)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusOK, raw)
}

func (s *Handler) handleParseKindleFile(w http.ResponseWriter, r *http.Request) error {
	userID, err := u.GetParamFromRequest(r, "userID")
	if err != nil {
		return err
	}

	// Parse the multipart form in the request
	// Maximum memory 20MB
	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		log.Println("Error parsing multipart form: ", err, " (file might be too large)")
		return u.WriteJSON(w, http.StatusBadRequest, err)
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return u.WriteJSON(w, http.StatusBadRequest, err)
	}
	defer file.Close()

	raw, err := parseKindleExtractFromFile(file)
	if err != nil {
		return u.WriteJSON(w, http.StatusBadRequest, err)
	}

	err = s.createDataFromRawBook(raw, userID)
	if err != nil {
		return err
	}

	return u.WriteJSON(w, http.StatusNoContent, "")
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

type ParseKindleFileRequest struct {
	File multipart.File `json:"file"`
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

func (s *Handler) createDataFromRawBook(raw *t.RawExtractBook, userID int) error {
	// Create book
	_, err := s.bookStore.GetBookByISBN(raw.ASIN)
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

	return nil
}
