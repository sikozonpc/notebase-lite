package types

import (
	"net/http"
	"time"
)

type EndpointHandler func(w http.ResponseWriter, r *http.Request) error

type Config struct {
	Env        string
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
	JWTSecret  string
}

type APIError struct {
	Error string `json:"error"`
}

type Highlight struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Location  string    `json:"location"`
	Note      string    `json:"note"`
	UserId    int       `json:"userId"`
	BookId    int       `json:"bookId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateHighlightRequest struct {
	Text     string `json:"text"`
	Location string `json:"location"`
	Note     string `json:"note"`
	UserId   int    `json:"userId"`
	BookId   int    `json:"bookId"`
}

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserStore is an interface that represents operations for storing and retrieving users.
type UserStore interface {
	CreateUser(User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
}

// HighlightStore is an interface that represents operations for storing and retrieving highlights.
type HighlightStore interface {
	GetUserHighlights(id int) ([]*Highlight, error)
	GetHighlightByID(id int) (*Highlight, error)
	CreateHighlight(Highlight) error
	DeleteHighlight(id int) error
}