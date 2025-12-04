package claims

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"sync"

	"log"

	"github.com/golang-jwt/jwt/v5"
)

type BaseClaims struct {
	parseOnly   sync.Once
	claims      jwt.Claims
	accessInfos AccessInfos
	metadata    Metadata
	endpoints   Endpoints
}

func NewBaseClaims(claims *jwt.MapClaims) *BaseClaims {
	return &BaseClaims{
		parseOnly: sync.Once{},
		claims:    claims,
	}
}

func (b *BaseClaims) ParseClaims() {
	b.parseOnly.Do(func() {
		if mapClaims, ok := b.claims.(*jwt.MapClaims); ok {
			jwtPayload := (*mapClaims)["urn:zitadel:iam:user:metadata"]
			if metadataMap, ok := jwtPayload.(map[string]interface{}); ok {
				var accessinfos AccessInfos
				accessInfoStr := strings.TrimSpace(metadataMap["accessinfo"].(string))
				for len(accessInfoStr)%4 != 0 {
					accessInfoStr += "="
				}
				decodedAccessInfo, err := base64.StdEncoding.DecodeString(accessInfoStr)
				if err != nil {
					log.Printf("decode accessinfo: metadataMap['accessinfo'] = %s\n", metadataMap["accessinfo"].(string))
					log.Printf("decode accessinfo error: err = %v\n", err)
				} else {
					if err := json.Unmarshal(decodedAccessInfo, &accessinfos); err != nil {
						log.Printf("accessinfo parse error: err = %v\n", err)
					} else {
						b.SetAccessInfos(accessinfos)
					}
				}

				var metadata Metadata
				metadataStr := strings.TrimSpace(metadataMap["metadata"].(string))
				for len(metadataStr)%4 != 0 {
					metadataStr += "="
				}
				decodedMetadata, err := base64.StdEncoding.DecodeString(metadataStr)
				if err != nil {
					log.Printf("decode metadata: metadataMap['metadata'] = %s\n", metadataMap["metadata"].(string))
					log.Printf("decode metadata error: err = %v\n", err)
				} else {
					if err := json.Unmarshal(decodedMetadata, &metadata); err != nil {
						log.Printf("metadata parse error: err = %v\n", err)
					} else {
						b.SetMetadata(metadata)
					}
				}

				var endpoints Endpoints
				var endpointIds []string
				endpointsStr := strings.TrimSpace(metadataMap["endpoints"].(string))
				for len(endpointsStr)%4 != 0 {
					endpointsStr += "="
				}
				decodedEndpoints, err := base64.StdEncoding.DecodeString(endpointsStr)
				if err != nil {
					log.Printf("decode endpoints: metadataMap['endpoints'] = %s\n", metadataMap["endpoints"].(string))
					log.Printf("decode endpoints error: err = %v\n", err)
				} else {
					if err := json.Unmarshal(decodedEndpoints, &endpointIds); err != nil {
						log.Printf("endpoints parse error: err = %v\n", err)
					} else {
						for _, endpointId := range endpointIds {
							endpoints = append(endpoints, Endpoint{Id: endpointId})
						}
						b.SetEndpoints(endpoints)
					}
				}
			}

			return
		} else {
			log.Println("access Debugs parse error")
			return
		}
	})
}

func (b *BaseClaims) PrintJWTClaims() {
	expirationTime, err := b.claims.GetExpirationTime()
	if err != nil {
		log.Printf("expirationTime parse error: err = %v\n", err)
		return
	} else {
		log.Printf("expirationTime is %+v\n", expirationTime)
	}

	audience, err := b.claims.GetAudience()
	if err != nil {
		log.Printf("audience parse error: err = %v\n", err)
		return
	} else {
		log.Printf("audience is %+v\n", audience)
	}

	issuedAt, err := b.claims.GetIssuedAt()
	if err != nil {
		log.Printf("issuedAt parse error: err = %v\n", err)
		return
	} else {
		log.Printf("issuedAt is %+v\n", issuedAt)
	}

	notBefore, err := b.claims.GetNotBefore()
	if err != nil {
		log.Printf("notBefore parse error: err = %v\n", err)
		return
	} else {
		log.Printf("notBefore is %+v\n", notBefore)
	}

	issuer, err := b.claims.GetIssuer()
	if err != nil {
		log.Printf("issuer parse error: err = %v\n", err)
		return
	} else {
		log.Printf("issuer is %+v\n", issuer)
	}

	subject, err := b.claims.GetSubject()
	if err != nil {
		log.Printf("subject parse error: err = %v\n", err)
		return
	} else {
		log.Printf("subject is %+v\n", subject)
	}
}

func (b *BaseClaims) GetJWTClaims() jwt.Claims {
	return b.claims
}

func (b *BaseClaims) GetAccessInfos() AccessInfos {
	b.ParseClaims()
	return b.accessInfos
}

func (b *BaseClaims) SetAccessInfos(accessInfos AccessInfos) {
	b.accessInfos = accessInfos
}

func (b *BaseClaims) GetMetadata() Metadata {
	b.ParseClaims()
	return b.metadata
}

func (b *BaseClaims) SetMetadata(metadata Metadata) {
	b.metadata = metadata
}

func (b *BaseClaims) GetEndpoints() Endpoints {
	b.ParseClaims()
	return b.endpoints
}

func (b *BaseClaims) SetEndpoints(endpoints Endpoints) {
	b.endpoints = endpoints
}
