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
	BookID    string    `json:"bookId"`
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

type Book struct {
	ISBN    string `json:"isbn"`
	Title   string `json:"title"`
	Authors string `json:"authors"`
}

// This is the format of the file that is downloaded from web tool
type RawExtractBook struct {
	ASIN       string                `json:"asin"`
	Title      string                `json:"title"`
	Authors    string                `json:"authors"`
	Highlights []RawExtractHighlight `json:"highlights"`
}

// This is the format of the file that is downloaded from web tool
type RawExtractHighlight struct {
	Text     string `json:"text"`
	Location struct {
		Value int    `json:"value"`
		URL   string `json:"url"`
	} `json:"location"`
	Note       string `json:"note"`
	IsNoteOnly bool   `json:"isNoteOnly"`
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
	CreateHighlights([]Highlight) error
	DeleteHighlight(id int) error
}

type BookStore interface {
	GetBookByISBN(isbn string) (*Book, error)
	CreateBook(Book) error
}
