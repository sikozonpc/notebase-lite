package user

import (
	t "github.com/sikozonpc/notebase/types"
)

func New(firstName, lastName, email, password string) *t.User {
	return &t.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
}
