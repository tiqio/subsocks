package main

import (
	"testing"

	"github.com/luyuhuang/subsocks/log"
	"github.com/stretchr/testify/assert"
)

func TestGetUserInfo(t *testing.T) {
	jwtInfo, err := GetJWTInfo(ZITADEL_TOKEN_URL, CLIENT_ID, CLIENT_SECRET)
	assert.NotNil(t, jwtInfo)
	assert.Nil(t, err)

	err = jwtInfo.GetAccessTokenClaims()
	assert.Nil(t, err)

	err = jwtInfo.GetIDTokenClaims()
	assert.Nil(t, err)

	log.Info("----- jwtInfo.GetAccessTokenClaims -----")
	jwtInfo.AccessTokenClaims.PrintJWTClaims()
	log.Info("jwtInfo.AccessTokenClaims ----- ", "AccessInfos", jwtInfo.AccessTokenClaims.GetAccessInfos())
	log.Info("jwtInfo.AccessTokenClaims ----- ", "Metadata", jwtInfo.AccessTokenClaims.GetMetadata())
	log.Info("jwtInfo.AccessTokenClaims ----- ", "Endpoints", jwtInfo.AccessTokenClaims.GetEndpoints())

	log.Info("----- jwtInfo.GetIDTokenClaims -----")
	jwtInfo.IDTokenClaims.PrintJWTClaims()
	log.Info("jwtInfo.IDTokenClaims ----- ", "AccessInfos", jwtInfo.IDTokenClaims.GetAccessInfos())
	log.Info("jwtInfo.IDTokenClaims ----- ", "Metadata", jwtInfo.IDTokenClaims.GetMetadata())
	log.Info("jwtInfo.IDTokenClaims ----- ", "Endpoints", jwtInfo.IDTokenClaims.GetEndpoints())
}
