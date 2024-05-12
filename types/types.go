package types

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EndpointHandler func(w http.ResponseWriter, r *http.Request) error

type Config struct {
	Env                string
	Port               string
	MongoURI           string
	JWTSecret          string // Used for signing JWT tokens
	GCPID              string // Google Cloud Project ID
	GCPBooksBucketName string // Google CLoud Storage Bucket Name from where upload books are parsed
	SendGridAPIKey     string
	SendGridFromEmail  string
	PublicURL          string // Used for generating links in emails
	APIKey             string // Used for authentication with external clients like GCP pub/sub
}

type APIError struct {
	Error string `json:"error"`
}

type Highlight struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Text      string             `json:"text" bson:"text"`
	Location  string             `json:"location" bson:"location"`
	Note      string             `json:"note" bson:"note"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	BookID    string             `json:"bookId" bson:"bookId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"`
	IsActive  bool               `json:"isActive" bson:"isActive"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

type Book struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	ISBN      string             `json:"isbn" bson:"isbn"`
	Title     string             `json:"title" bson:"title"`
	Authors   string             `json:"authors" bson:"authors"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
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
	Create(context.Context, RegisterRequest) (primitive.ObjectID, error)
	GetUserByEmail(context.Context, string) (*User, error)
	GetUserByID(context.Context, string) (*User, error)
	GetUsers(context.Context) ([]*User, error)
	UpdateUser(context.Context, User) error
}

type HighlightStore interface {
	CreateHighlight(context.Context, *CreateHighlightRequest) (primitive.ObjectID, error)
	GetHighlightByID(context.Context, primitive.ObjectID, primitive.ObjectID) (*Highlight, error)
	GetUserHighlights(context.Context, primitive.ObjectID) ([]*Highlight, error)
	DeleteHighlight(context.Context, primitive.ObjectID) error
	GetRandomHighlights(context.Context, primitive.ObjectID, int) ([]*Highlight, error)
}

type BookStore interface {
	GetByISBN(context.Context, string) (*Book, error)
	Create(context.Context, *CreateBookRequest) (primitive.ObjectID, error)
}

type CreateBookRequest struct {
	ISBN    string `json:"isbn" bson:"isbn"`
	Title   string `json:"title" bson:"title"`
	Authors string `json:"authors" bson:"authors"`
}

type CreateHighlightRequest struct {
	Text     string             `json:"text" bson:"text"`
	Location string             `json:"location" bson:"location"`
	Note     string             `json:"note" bson:"note"`
	UserID   primitive.ObjectID `json:"userId" bson:"userId"`
	BookID   string             `json:"bookId" bson:"bookId"`
}

type DailyInsight struct {
	Text        string
	Note        string
	BookAuthors string
	BookTitle   string
}

type RegisterRequest struct {
	FirstName string `json:"firstName" bson:"firstName" validate:"required"`
	LastName  string `json:"lastName" bson:"lastName" validate:"required"`
	Email     string `json:"email" bson:"email" validate:"required"`
	Password  string `json:"password" bson:"password" validate:"required"`
}
