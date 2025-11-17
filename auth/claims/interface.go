package claims

import "github.com/golang-jwt/jwt/v5"

type Metadata struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type AccessInfos []AccessInfo
type AccessInfo struct {
	Id       string `json:"access_id"`
	Services Services
}

type Services []Service
type Service struct {
	Id    string  `json:"service_id"`
	Delay float64 `json:"delay"`
}

type Endpoints []Endpoint
type Endpoint struct {
	Id string `json:"id"`
}

type Claims interface {
	PrintJWTClaims()
	GetJWTClaims() jwt.Claims
	GetAccessInfos() AccessInfos
	SetAccessInfos(AccessInfos)
	GetMetadata() Metadata
	SetMetadata(Metadata)
	GetEndpoints() Endpoints
	SetEndpoints(Endpoints)
}
