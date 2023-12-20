package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	t "github.com/sikozonpc/notebase/types"
	u "github.com/sikozonpc/notebase/utils"
)

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		token, err := validateJWT(tokenString)
		if err != nil {
			u.WriteJSON(w, http.StatusUnauthorized, t.APIError{
				Error: fmt.Errorf("invalid token").Error(),
			})
			return
		}

		if !token.Valid {
			u.WriteJSON(w, http.StatusUnauthorized, t.APIError{
				Error: fmt.Errorf("invalid token").Error(),
			})
			return
		}

		fmt.Println(token)

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
