package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")
	userID := 1

	token, err := CreateJWT(secret, userID)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	assert.NotEmpty(t, token)
}