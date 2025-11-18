package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luyuhuang/subsocks/auth/claims"
)

type JWTResponse struct {
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	ExpiresIn   float64 `json:"expires_in"`
	IDToken     string  `json:"id_token"`
}

type JWTInfo struct {
	JWTResponse       *JWTResponse
	AccessTokenClaims claims.Claims
	IDTokenClaims     claims.Claims
}

func NewJWTInfo(jwtResponse *JWTResponse) *JWTInfo {
	return &JWTInfo{
		JWTResponse:       jwtResponse,
		AccessTokenClaims: claims.NewAccessTokenClaims(&jwt.MapClaims{}),
		IDTokenClaims:     claims.NewIDTokenClaims(&jwt.MapClaims{}),
	}
}

func (jwtInfo *JWTInfo) GetAccessTokenClaims() error {
	_, _, err := jwt.NewParser().ParseUnverified(jwtInfo.JWTResponse.AccessToken, jwtInfo.AccessTokenClaims.GetJWTClaims())
	if err != nil {
		return fmt.Errorf("error parsing access token claims: %w", err)
	}

	return nil
}

func (jwtInfo *JWTInfo) GetIDTokenClaims() error {
	_, _, err := jwt.NewParser().ParseUnverified(jwtInfo.JWTResponse.IDToken, jwtInfo.IDTokenClaims.GetJWTClaims())
	if err != nil {
		return fmt.Errorf("error parsing ID token claims: %w", err)
	}

	return nil
}
