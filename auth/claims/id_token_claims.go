package claims

import "github.com/golang-jwt/jwt/v5"

type IDTokenClaims struct {
	BaseClaims
}

func NewIDTokenClaims(claims *jwt.MapClaims) *IDTokenClaims {
	return &IDTokenClaims{
		BaseClaims: *NewBaseClaims(claims),
	}
}
