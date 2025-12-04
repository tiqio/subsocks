package auth

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserInfo(t *testing.T) {
	jwtInfo, err := GetJWTInfo(ZITADEL_TOKEN_URL, CLIENT_ID, CLIENT_SECRET, PROJECT_ID)
	assert.NotNil(t, jwtInfo)
	assert.Nil(t, err)

	err = jwtInfo.GetAccessTokenClaims()
	assert.Nil(t, err)

	err = jwtInfo.GetIDTokenClaims()
	assert.Nil(t, err)

	log.Println("----- jwtInfo.GetAccessTokenClaims -----")
	jwtInfo.AccessTokenClaims.PrintJWTClaims()
	log.Printf("jwtInfo.AccessTokenClaims ----- AccessInfos: %+v\n", jwtInfo.AccessTokenClaims.GetAccessInfos())
	log.Printf("jwtInfo.AccessTokenClaims ----- Metadata: %+v\n", jwtInfo.AccessTokenClaims.GetMetadata())
	log.Printf("jwtInfo.AccessTokenClaims ----- Endpoints: %+v\n", jwtInfo.AccessTokenClaims.GetEndpoints())

	log.Println("----- jwtInfo.GetIDTokenClaims -----")
	jwtInfo.IDTokenClaims.PrintJWTClaims()
	log.Printf("jwtInfo.IDTokenClaims ----- AccessInfos: %+v\n", jwtInfo.IDTokenClaims.GetAccessInfos())
	log.Printf("jwtInfo.IDTokenClaims ----- Metadata: %+v\n", jwtInfo.IDTokenClaims.GetMetadata())
	log.Printf("jwtInfo.IDTokenClaims ----- Endpoints: %+v\n", jwtInfo.IDTokenClaims.GetEndpoints())
}
