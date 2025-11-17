package claims

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	BaseClaims
}

func NewAccessTokenClaims(claims *jwt.MapClaims) *AccessTokenClaims {
	return &AccessTokenClaims{
		BaseClaims: *NewBaseClaims(claims),
	}
}
