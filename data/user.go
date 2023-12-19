package data

import (
	t "github.com/sikozonpc/notebase/types"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(firstName, lastName, email, password string) *t.User {
	return &t.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswords(hashed string, plain []byte) bool {
	byteHash := []byte(hashed)

	err := bcrypt.CompareHashAndPassword(byteHash, plain)
	return err != nil
}
