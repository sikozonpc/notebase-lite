package config

import (
	"fmt"
	"os"

	t "github.com/sikozonpc/notebase/types"
)

var Envs = initConfig()

func initConfig() t.Config {
	return t.Config{
		Env:                getEnv("ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnv("DB_PASSWORD", "password"),
		DBAddress:          fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:             getEnv("DB_NAME", "notebase"),
		JWTSecret:          getEnvOrPanic("JWT_SECRET", "JWT secret is required"),
		GCPID:              getEnvOrPanic("GOOGLE_CLOUD_PROJECT_ID", "Google Cloud Project ID is required"),
		GCPBooksBucketName: getEnvOrPanic("GOOGLE_CLOUD_BOOKS_BUCKET_NAME", "Google Cloud Books Bucket Name is required"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvOrPanic(key, err string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(err)
}
