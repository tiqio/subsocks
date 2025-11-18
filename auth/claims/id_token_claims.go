package claims

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/luyuhuang/subsocks/log"
)

type IDTokenClaims struct {
	BaseClaims
}

func NewIDTokenClaims(claims *jwt.MapClaims) *IDTokenClaims {
	return &IDTokenClaims{
		BaseClaims: *NewBaseClaims(claims),
	}
}

func (i *IDTokenClaims) PrintJWTClaims() {
	log.Info("ID Token Claims", "metadata", i.GetMetadata())
	i.BaseClaims.PrintJWTClaims()
}
