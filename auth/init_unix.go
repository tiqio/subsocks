//go:build !windows
// +build !windows

package auth

import (
	"log"

	"github.com/joho/godotenv"
)

var (
	PROJECT_ID                string
	ZITADEL_DOMAIN            string
	ZITADEL_TOKEN_URL         string
	CLIENT_ID                 string
	CLIENT_SECRET             string
	ZITADEL_INTROSPECTION_URL string
	API_CLIENT_ID             string
	API_CLIENT_SECRET         string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return
	}

	PROJECT_ID = getEnv("PROJECT_ID")
	ZITADEL_DOMAIN = getEnv("ZITADEL_DOMAIN")
	ZITADEL_TOKEN_URL = getEnv("ZITADEL_TOKEN_URL")
	CLIENT_ID = getEnv("CLIENT_ID")
	CLIENT_SECRET = getEnv("CLIENT_SECRET")
	ZITADEL_INTROSPECTION_URL = getEnv("ZITADEL_INTROSPECTION_URL")
	API_CLIENT_ID = getEnv("API_CLIENT_ID")
	API_CLIENT_SECRET = getEnv("API_CLIENT_SECRET")
}
