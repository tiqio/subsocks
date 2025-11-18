package claims

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/luyuhuang/subsocks/log"
)

type AccessTokenClaims struct {
	BaseClaims
}

func NewAccessTokenClaims(claims *jwt.MapClaims) *AccessTokenClaims {
	return &AccessTokenClaims{
		BaseClaims: *NewBaseClaims(claims),
	}
}

func (a *AccessTokenClaims) PrintJWTClaims() {
	log.Info("Access Token Claims", "metadata", a.GetMetadata())
	a.BaseClaims.PrintJWTClaims()
}
