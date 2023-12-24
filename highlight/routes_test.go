package highlight

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	t "github.com/sikozonpc/notebase/types"
)

func TestHandleUserHighlights(t *testing.T) {
	store := &mockHighlightStore{}
	userStore := &mockUserStore{}
	handler := NewHandler(store, userStore)

	t.Run("should handle get user highlights", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/user/1/highlight", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleGetUserHighlights(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should handle get highlight by id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/user/1/highlight/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleGetHighlight(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should handle create highlight", func(t *testing.T) {
		payload := CreateHighlightRequest{
			Text:     "test",
			Location: "test",
			Note:     "test",
			BookId:   1,
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
		handler.handleCreateHighlight(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should handle delete highlight", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/user/1/highlight/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleDeleteHighlight(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockHighlightStore struct{}

func (m *mockHighlightStore) GetUserHighlights(userID int) ([]*t.Highlight, error) {
	return []*t.Highlight{}, nil
}

func (m *mockHighlightStore) CreateHighlight(h t.Highlight) error {
	return nil
}

func (m *mockHighlightStore) GetHighlightByID(id, userID int) (*t.Highlight, error) {
	return &t.Highlight{}, nil
}

func (m *mockHighlightStore) DeleteHighlight(id int) error {
	return nil
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*t.User, error) {
	return &t.User{}, nil
}

func (m *mockUserStore) CreateUser(u t.User) error {
	return nil
}

func (m *mockUserStore) GetUserByID(id int) (*t.User, error) {
	return &t.User{}, nil
}
