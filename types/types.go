package types

import (
	"net/http"
	"time"
)

type EndpointHandler func(w http.ResponseWriter, r *http.Request) error

type Config struct {
	Env                string
	Port               string
	DBUser             string
	DBPassword         string
	DBAddress          string
	DBName             string
	JWTSecret          string
	GCPID              string
	GCPBooksBucketName string
}

type APIError struct {
	Error string `json:"error"`
}

type Highlight struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Location  string    `json:"location"`
	Note      string    `json:"note"`
	UserID    int       `json:"userId"`
	BookID    int       `json:"bookId"`
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

type UserStore interface {
	CreateUser(User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
}

type HighlightStore interface {
	GetUserHighlights(userID int) ([]*Highlight, error)
	GetHighlightByID(id, userID int) (*Highlight, error)
	CreateHighlight(Highlight) error
	DeleteHighlight(id int) error
}
