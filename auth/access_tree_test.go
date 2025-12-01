package auth

import (
	"testing"

	"github.com/luyuhuang/subsocks/auth/claims"
	acc "github.com/luyuhuang/subsocks/control/access"
	srv "github.com/luyuhuang/subsocks/control/service"
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

func TestAccessTree_WithSwitch(t *testing.T) {
	jwtInfo := NewJWTInfo(&JWTResponse{})

	var accessInfos = []claims.AccessInfo{
		{
			Id: "yyy1",
			Services: []claims.Service{
				{
					Id:    "zzz2",
					Delay: 120,
				},
				{
					Id:    "zzz1",
					Delay: 100,
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
				Id: "zzz2",
				Accesses: []Access{
					{
						Id:    "yyy1",
						Delay: 120,
					},
				},
			},
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
		},
	}

	assert.Equal(t, expectedAccessTree, accessTree)
}

func TestAccessTree_WithSort(t *testing.T) {
	jwtInfo := NewJWTInfo(&JWTResponse{})

	var accessInfos = []claims.AccessInfo{
		{
			Id: "yyy1",
			Services: []claims.Service{
				{
					Id:    "zzz1",
					Delay: 150,
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
						Id:    "yyy2",
						Delay: 130,
					},
					{
						Id:    "yyy1",
						Delay: 150,
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

func TestAccessTree_FillInfo(t *testing.T) {
	accessTree := AccessTree{
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

	accessTree.FillInfo()

	expectedAccessTree := AccessTree{
		ServiceIds: map[string]struct{}{
			"zzz1": {},
			"zzz2": {},
		},
		ServiceTrees: []ServiceTree{
			{
				Id: "zzz1",
				ServiceInfo: srv.Info{
					Name: "IPSB",
					Host: "ip.sb",
					Port: 443,
				},
				Accesses: []Access{
					{
						Id: "yyy1",
						AccessInfo: acc.Info{
							Name: "Ubuntu25",
							Host: "192.168.235.128",
							Port: 1080,
						},
						Delay: 100,
					},
					{
						Id: "yyy2",
						AccessInfo: acc.Info{
							Name: "Ubuntu22",
							Host: "192.168.235.133",
							Port: 1080,
						},
						Delay: 130,
					},
				},
			},
			{
				Id: "zzz2",
				ServiceInfo: srv.Info{
					Name: "YouTube",
					Host: "www.youtube.com",
					Port: 443,
				},
				Accesses: []Access{
					{
						Id: "yyy1",
						AccessInfo: acc.Info{
							Name: "Ubuntu25",
							Host: "192.168.235.128",
							Port: 1080,
						},
						Delay: 120,
					},
				},
			},
		},
	}

	assert.Equal(t, expectedAccessTree, accessTree)
}

func TestAccessTree_ListService(t *testing.T) {
	accessTree := AccessTree{
		ServiceIds: map[string]struct{}{
			"zzz1": {},
			"zzz2": {},
		},
		ServiceTrees: []ServiceTree{
			{
				Id: "zzz1",
				ServiceInfo: srv.Info{
					Name: "IPSB",
					Host: "ip.sb",
					Port: 443,
				},
				Accesses: []Access{
					{
						Id: "yyy1",
						AccessInfo: acc.Info{
							Name: "Ubuntu25",
							Host: "192.168.235.128",
							Port: 1080,
						},
						Delay: 100,
					},
					{
						Id: "yyy2",
						AccessInfo: acc.Info{
							Name: "Ubuntu22",
							Host: "192.168.235.133",
							Port: 1080,
						},
						Delay: 130,
					},
				},
			},
			{
				Id: "zzz2",
				ServiceInfo: srv.Info{
					Name: "YouTube",
					Host: "www.youtube.com",
					Port: 443,
				},
				Accesses: []Access{
					{
						Id: "yyy1",
						AccessInfo: acc.Info{
							Name: "Ubuntu25",
							Host: "192.168.235.128",
							Port: 1080,
						},
						Delay: 120,
					},
				},
			},
		},
	}

	services := accessTree.ListServices()

	expectedService := []srv.Info{
		{
			Name: "IPSB",
			Host: "ip.sb",
			Port: 443,
		},
		{
			Name: "YouTube",
			Host: "www.youtube.com",
			Port: 443,
		},
	}

	assert.Equal(t, expectedService, services)
}
