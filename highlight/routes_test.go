package highlight

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sikozonpc/notebase/storage"
	types "github.com/sikozonpc/notebase/types"
	u "github.com/sikozonpc/notebase/utils"
)

var fakeHighlight *types.Highlight

func TestHandleUserHighlights(t *testing.T) {
	memStore := storage.NewMemoryStorage()
	bookStore := &mockBookStore{}

	store := &mockHighlightStore{}
	userStore := &mockUserStore{}
	handler := NewHandler(store, userStore, memStore, bookStore)

	t.Run("should handle get user highlights", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/user/1/highlight", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}/highlight", u.MakeHTTPHandler(handler.handleGetUserHighlights))

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail to handle get highlight by ID if highlight does not exist", func(t *testing.T) {
		fakeHighlight = nil

		req, err := http.NewRequest(http.MethodGet, "/user/1/highlight/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}/highlight/{id}", u.MakeHTTPHandler(handler.handleGetHighlightByID)).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should handle get highlight by ID", func(t *testing.T) {
		fakeHighlight = &types.Highlight{
			ID:       1,
			Text:     "test",
			Location: "test",
			Note:     "test",
			BookID:   "B004XCFJ3E",
			UserID:   1,
		}

		req, err := http.NewRequest(http.MethodGet, "/user/1/highlight/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}/highlight/{id}", u.MakeHTTPHandler(handler.handleGetHighlightByID)).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should handle create highlight", func(t *testing.T) {
		payload := CreateHighlightRequest{
			Text:     "test",
			Location: "test",
			Note:     "test",
			BookId:   "B004XCFJ3E",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/user/1/highlight", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}/highlight", u.MakeHTTPHandler(handler.handleCreateHighlight))

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if rr.Body.String() == "" {
			t.Errorf("expected response body to be non-empty")
		}

		var response types.Highlight
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Text != payload.Text {
			t.Errorf("expected text to be %s, got %s", payload.Text, response.Text)
		}

		if response.Note != payload.Note {
			t.Errorf("expected location to be %s, got %s", payload.Location, response.Location)
		}
	})

	t.Run("should handle delete highlight", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/user/1/highlight/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}/highlight/{id}", u.MakeHTTPHandler(handler.handleDeleteHighlight)).Methods(http.MethodDelete)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail handle parse kindle extract if filename is not sent", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/user/1/parse-kindle-extract", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}/parse-kindle-extract", u.MakeHTTPHandler(handler.handleParseKindleFile))

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should handle parse kindle extract", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/user/1/parse-kindle-extract?filename=file.json", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}/parse-kindle-extract", u.MakeHTTPHandler(handler.handleParseKindleFile))

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockHighlightStore struct{}

func (m *mockHighlightStore) GetUserHighlights(userID int) ([]*types.Highlight, error) {
	return []*types.Highlight{}, nil
}

func (m *mockHighlightStore) CreateHighlight(h types.Highlight) error {
	return nil
}

func (m *mockHighlightStore) GetHighlightByID(id, userID int) (*types.Highlight, error) {
	return fakeHighlight, nil
}

func (m *mockHighlightStore) DeleteHighlight(id int) error {
	return nil
}

func (m *mockHighlightStore) CreateHighlights(highlights []types.Highlight) error {
	return nil
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) CreateUser(u types.User) error {
	return nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return &types.User{}, nil
}

type mockBookStore struct{}

func (m *mockBookStore) GetBookByISBN(ISBN string) (*types.Book, error) {
	return &types.Book{}, nil
}

func (m *mockBookStore) CreateBook(book types.Book) error {
	return nil
}
