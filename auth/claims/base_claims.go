package claims

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luyuhuang/subsocks/log"
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
		log.Debug("get access Debugs")
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
					log.Debug("decode accessinfo", "metadataMap['accessinfo']", metadataMap["accessinfo"].(string))
					log.Error("decode accessinfo error", "err", err)
				} else {
					if err := json.Unmarshal(decodedAccessInfo, &accessinfos); err != nil {
						log.Debug("accessinfo parse error", "err", err)
					} else {
						//log.Debug("unmarshal accessinfo", "accessinfo", accessinfo)
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
					log.Debug("decode metadata", "metadataMap['metadata']", metadataMap["metadata"].(string))
					log.Error("decode metadata error", "err", err)
				} else {
					if err := json.Unmarshal(decodedMetadata, &metadata); err != nil {
						log.Debug("metadata parse error", "err", err)
					} else {
						//log.Debug("unmarshal metadata", "metadata", metadata)
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
					log.Debug("decode endpoints", "metadataMap['endpoints']", metadataMap["endpoints"].(string))
					log.Error("decode endpoints error", "err", err)
				} else {
					if err := json.Unmarshal(decodedEndpoints, &endpointIds); err != nil {
						log.Debug("endpoints parse error", "err", err)
					} else {
						for _, endpointId := range endpointIds {
							endpoints = append(endpoints, Endpoint{Id: endpointId})
						}
						//log.Debug("unmarshal endpoints", "endpoints", endpoints)
						b.SetEndpoints(endpoints)
					}
				}
			}

			return
		} else {
			log.Error("access Debugs parse error")
			return
		}
	})
}

func (b *BaseClaims) PrintJWTClaims() {
	expirationTime, err := b.claims.GetExpirationTime()
	if err != nil {
		log.Error("expirationTime parse error", "err", err)
		return
	} else {
		log.Debug(fmt.Sprintf("expirationTime is %+v", expirationTime))
	}

	audience, err := b.claims.GetAudience()
	if err != nil {
		log.Error("audience parse error", "err", err)
		return
	} else {
		log.Debug(fmt.Sprintf("audience is %+v", audience))
	}

	issuedAt, err := b.claims.GetIssuedAt()
	if err != nil {
		log.Error("issuedAt parse error", "err", err)
		return
	} else {
		log.Debug(fmt.Sprintf("issuedAt is %+v", issuedAt))
	}

	notBefore, err := b.claims.GetNotBefore()
	if err != nil {
		log.Error("notBefore parse error", "err", err)
		return
	} else {
		log.Debug(fmt.Sprintf("notBefore is %+v", notBefore))
	}

	issuer, err := b.claims.GetIssuer()
	if err != nil {
		log.Error("issuer parse error", "err", err)
		return
	} else {
		log.Debug(fmt.Sprintf("issuer is %+v", issuer))
	}

	subject, err := b.claims.GetSubject()
	if err != nil {
		log.Error("subject parse error", "err", err)
		return
	} else {
		log.Debug(fmt.Sprintf("subject is %+v", subject))
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
