package types

import "time"

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
	Password  string    `json:"password"`
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
