package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	t "github.com/sikozonpc/notebase/types"
	u "github.com/sikozonpc/notebase/utils"
)

func permissionDenied(w http.ResponseWriter) {
	u.WriteJSON(w, http.StatusUnauthorized, t.APIError{
		Error: fmt.Errorf("permission denied").Error(),
	})
}

func GetUserFromToken(t string) (int, error) {
	token, err := validateJWT(t)
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	claimsUserID := claims["userID"].(string)

	userID, err := strconv.Atoi(claimsUserID)
	if err != nil {
		log.Println("failed to convert userID to int")
		return 0, err
	}

	return userID, nil
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

		userID, err := strconv.Atoi(claimsUserID)
		if err != nil {
			log.Println("failed to convert userID to int")
			permissionDenied(w)
			return
		}

		_, err = store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

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
