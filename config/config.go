package config

import (
	"os"

	t "github.com/sikozonpc/notebase/types"
)

var Envs = initConfig()

func initConfig() t.Config {
	return t.Config{
		Env:       getEnv("ENV", "development"),
		Port:      getEnv("PORT", "8080"),
		MongoURI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		PublicURL: getEnv("PUBLIC_URL", "http://localhost:3000"),
		JWTSecret: getEnv("JWT_SECRET", "JWT secret is required"),
		SendGridAPIKey:    getEnv("SENDGRID_API_KEY", "SendGrid API KEY is required"),
		SendGridFromEmail: getEnv("SENDGRID_FROM_EMAIL", "SendGrid From email is required"),
		APIKey:            getEnv("API_KEY", "API Key is required"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
