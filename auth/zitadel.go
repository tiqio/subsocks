package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/luyuhuang/subsocks/log"
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

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Error(fmt.Sprintf("Environment variable %s is not set. Exiting...", key))
	} else {
		log.Info("Environment variable found", key, value)
	}
	return value
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
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

	log.Init("", getEnv("LOG_LEVEL"))
}

func GetJWTInfo(zitadelTokenUrl string, clientId string, clientSecret string) (*JWTInfo, error) {
	// Encode the client ID and client Secret in Base64
	clientCredentials := fmt.Sprintf("%s:%s", clientId, clientSecret)
	base64ClientCredentials := base64.StdEncoding.EncodeToString([]byte(clientCredentials))

	// Prepare the request headers
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "Basic " + base64ClientCredentials,
	}

	// Prepare the request data
	data := fmt.Sprintf("grant_type=client_credentials&scope=openid profile email urn:zitadel:iam:org:project:id:%s:aud urn:zitadel:iam:org:projects:roles urn:zitadel:iam:user:metadata", PROJECT_ID)

	// Create a new request
	req, err := http.NewRequest("POST", zitadelTokenUrl, bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Error("Error creating request")
	}

	// Set the headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode == http.StatusOK {
		var responseData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %w", err)
		}

		return NewJWTInfo(&JWTResponse{
			AccessToken: responseData["access_token"].(string),
			TokenType:   responseData["token_type"].(string),
			ExpiresIn:   responseData["expires_in"].(float64),
			IDToken:     responseData["id_token"].(string),
		}), nil
	}

	return nil, fmt.Errorf("received status code %d", resp.StatusCode)
}
