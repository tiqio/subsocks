package main

import (
	"testing"

	"github.com/luyuhuang/subsocks/auth/claims"
	"github.com/stretchr/testify/assert"
)

func TestAccessTree(t *testing.T) {
	jwtInfo := NewJWTInfo(&JWTResponse{})

	var accessInfos = []claims.AccessInfo{
		{
			Id: "yyy1",
			Services: []claims.Service{
				{
					Id:    "zzz1",
					Delay: 100,
				}, {
					Id:    "zzz2",
					Delay: 120,
				},
			},
		},
		{
			Id: "yyy2",
			Services: []claims.Service{
				{
					Id:    "zzz1",
					Delay: 130,
				},
			},
		},
	}

	jwtInfo.AccessTokenClaims.SetAccessInfos(accessInfos)

	accessTree := jwtInfo.AccessTree()

	expectedAccessTree := AccessTree{
		ServiceIds: map[string]struct{}{
			"zzz1": {},
			"zzz2": {},
		},
		ServiceTrees: []ServiceTree{
			{
				Id: "zzz1",
				Accesses: []Access{
					{
						Id:    "yyy1",
						Delay: 100,
					},
					{
						Id:    "yyy2",
						Delay: 130,
					},
				},
			},
			{
				Id: "zzz2",
				Accesses: []Access{
					{
						Id:    "yyy1",
						Delay: 120,
					},
				},
			},
		},
	}

	assert.Equal(t, expectedAccessTree, accessTree)
}
