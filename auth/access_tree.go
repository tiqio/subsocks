package main

import "sort"

type AccessTree struct {
	ServiceTrees []ServiceTree
	ServiceIds   map[string]struct{}
}

type ServiceTree struct {
	Id       string `json:"id"`
	Accesses Accesses
}

type Accesses []Access
type Access struct {
	Id    string  `json:"id"`
	Delay float64 `json:"delay"`
}

func (accessTree *AccessTree) AddService(serviceId string) {
	accessTree.ServiceIds[serviceId] = struct{}{}
}

func (accessTree *AccessTree) HasService(serviceId string) bool {
	_, exists := accessTree.ServiceIds[serviceId]
	return exists
}

func (jwtInfo *JWTInfo) AccessTree() AccessTree {
	accessTree := AccessTree{
		ServiceIds:   make(map[string]struct{}),
		ServiceTrees: make([]ServiceTree, 0),
	}

	for _, accessInfo := range jwtInfo.AccessTokenClaims.GetAccessInfos() {
		for _, service := range accessInfo.Services {
			if accessTree.HasService(service.Id) {
				for idx, serviceTree := range accessTree.ServiceTrees {
					if serviceTree.Id != service.Id {
						continue
					}

					accesses := serviceTree.Accesses

					var hasAccess bool
					for _, access := range accesses {
						if access.Id == accessInfo.Id {
							hasAccess = true
							break
						}
					}

					if hasAccess {
						continue
					} else {
						newAccess := Access{
							Id:    accessInfo.Id,
							Delay: service.Delay,
						}
						accesses = append(accesses, newAccess)
						accessTree.ServiceTrees[idx].Accesses = accesses
						sort.Slice(accessTree.ServiceTrees[idx].Accesses, func(i, j int) bool {
							return accessTree.ServiceTrees[idx].Accesses[i].Delay <
								accessTree.ServiceTrees[idx].Accesses[j].Delay
						})
					}
				}
			} else {
				var serviceTree ServiceTree
				serviceTree.Id = service.Id
				var access Access
				access.Id = accessInfo.Id
				access.Delay = service.Delay
				serviceTree.Accesses = append(serviceTree.Accesses, access)
				accessTree.ServiceTrees = append(accessTree.ServiceTrees, serviceTree)
				accessTree.AddService(service.Id)
			}
		}
	}

	return accessTree
}
