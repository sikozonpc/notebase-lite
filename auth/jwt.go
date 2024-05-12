package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sikozonpc/notebase/config"
	t "github.com/sikozonpc/notebase/types"
	u "github.com/sikozonpc/notebase/utils"
)

func permissionDenied(w http.ResponseWriter) {
	u.WriteJSON(w, http.StatusUnauthorized, t.APIError{
		Error: fmt.Errorf("permission denied").Error(),
	})
}

func GetUserFromToken(t string) (string, error) {
	token, err := validateJWT(t)
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	claimsUserID := claims["userID"].(string)

	return claimsUserID, nil
}

func WithAPIKey(handlerFunc http.HandlerFunc) http.HandlerFunc {
	apiKey := config.Envs.APIKey

	return func(w http.ResponseWriter, r *http.Request) {
		apiKeyFromRequest := r.Header.Get("X-API-KEY")

		if apiKeyFromRequest != apiKey {
			log.Println("invalid api key")
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store t.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := u.GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		claimsUserID := claims["userID"].(string)

		_, err = store.GetUserByID(context.Background(), claimsUserID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
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
