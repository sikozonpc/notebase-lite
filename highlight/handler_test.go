package highlight

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sikozonpc/notebase/storage"
	types "github.com/sikozonpc/notebase/types"
	u "github.com/sikozonpc/notebase/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var fakeHighlight *types.Highlight

func TestHandleUserHighlights(t *testing.T) {
	memStore := storage.NewMemoryStorage()
	bookStore := &mockBookStore{}
	mockMailer := &mockMailer{}

	store := &mockHighlightStore{}
	userStore := &mockUserStore{}
	handler := NewHandler(store, userStore, memStore, bookStore, mockMailer)

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
			ID:       primitive.NewObjectID(),
			Text:     "test",
			Location: "test",
			Note:     "test",
			BookID:   "B004XCFJ3E",
			UserID:   primitive.NewObjectID(),
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

	t.Run("should handle parse kindle extract", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/user/1/cloud/parse-kindle-extract/file.json", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}/cloud/parse-kindle-extract/{fileName}", u.MakeHTTPHandler(handler.handleCloudKindleParse))

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockHighlightStore struct{}

func (m *mockHighlightStore) CreateHighlight(context.Context, *types.CreateHighlightRequest) (primitive.ObjectID, error) {
	return primitive.NilObjectID, nil
}

func (m *mockHighlightStore) GetHighlightByID(context.Context, primitive.ObjectID, primitive.ObjectID) (*types.Highlight, error) {
	return fakeHighlight, nil
}

func (m *mockHighlightStore) GetUserHighlights(context.Context, primitive.ObjectID) ([]*types.Highlight, error) {
	return []*types.Highlight{}, nil
}

func (m *mockHighlightStore) DeleteHighlight(context.Context, primitive.ObjectID) error {
	return nil
}

func (m *mockHighlightStore) GetRandomHighlights(context.Context, primitive.ObjectID, int) ([]*types.Highlight, error) {
	return []*types.Highlight{}, nil
}

type mockUserStore struct{}

func (m *mockUserStore) Create(context.Context, types.RegisterRequest) (primitive.ObjectID, error) {
	return primitive.NilObjectID, nil
}

func (m *mockUserStore) GetUserByID(context.Context, string) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) GetUsers(context.Context) ([]*types.User, error) {
	return []*types.User{}, nil
}

func (m *mockUserStore) GetUserByEmail(context.Context, string) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) UpdateUser(context.Context, types.User) error {
	return nil
}

type mockBookStore struct{}

func (m *mockBookStore) GetByISBN(context.Context, string) (*types.Book, error) {
	return &types.Book{}, nil
}

func (m *mockBookStore) Create(context.Context, *types.CreateBookRequest) (primitive.ObjectID, error) {
	return primitive.NilObjectID, nil
}

type mockMailer struct{}

func (m *mockMailer) SendMail(string, string, string) error {
	return nil
}

func (m *mockMailer) SendInsights(*types.User, []*types.DailyInsight, string) error {
	return nil
}
