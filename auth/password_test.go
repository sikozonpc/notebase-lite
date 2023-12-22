package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	assert.NotEqual(t, hash, "password")
	assert.NotEmpty(t, hash)
}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	assert.True(t, ComparePasswords(hash, []byte("password")))
	assert.False(t, ComparePasswords(hash, []byte("notpassword")))
}