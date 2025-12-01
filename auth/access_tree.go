package auth

import (
	"sort"

	acc "github.com/luyuhuang/subsocks/control/access"
	"github.com/luyuhuang/subsocks/control/rule"
	srv "github.com/luyuhuang/subsocks/control/service"
	"github.com/luyuhuang/subsocks/log"
)

type AccessTree struct {
	ServiceTrees []ServiceTree
	ServiceIds   map[string]struct{}
}

type ServiceTree struct {
	Id          string   `json:"id"`
	ServiceInfo srv.Info `json:"service_info"`
	Accesses    Accesses
}

type Accesses []Access
type Access struct {
	Id         string   `json:"id"`
	AccessInfo acc.Info `json:"access_info"`
	Delay      float64  `json:"delay"`
}

func (accessTree *AccessTree) AddService(serviceId string) {
	if serviceId == "" {
		return
	}
	accessTree.ServiceIds[serviceId] = struct{}{}
}

func (accessTree *AccessTree) HasService(serviceId string) bool {
	_, exists := accessTree.ServiceIds[serviceId]
	return exists
}

func (accessTree *AccessTree) ListServices() []srv.Info {
	var serviceInfos []srv.Info

	for _, serviceTree := range accessTree.ServiceTrees {
		serviceInfos = append(serviceInfos, serviceTree.ServiceInfo)
	}

	return serviceInfos
}

func (accessTree *AccessTree) FillInfo() {
	for svcIdx, service := range accessTree.ServiceTrees {
		if service.Id == "" {
			return
		}
		svcInfo, err := srv.GetInfoById(service.Id)
		if err != nil {
			log.Error("GetInfoById failed", service.Id, err)
			return
		}
		service.ServiceInfo = svcInfo
		for accIdx, access := range service.Accesses {
			accInfo, err := acc.GetInfoById(access.Id)
			if err != nil {
				log.Error("GetInfoById failed", access.Id, err)
				return
			}
			access.AccessInfo = accInfo
			service.Accesses[accIdx] = access
		}
		accessTree.ServiceTrees[svcIdx] = service
	}
}

func (accessTree *AccessTree) ListRule() []rule.Info {
	var ruleInfos []rule.Info

	for _, service := range accessTree.ServiceTrees {
		accessInfo := service.Accesses[0].AccessInfo
		ruleInfo := rule.NewInfo(accessInfo, service.ServiceInfo, "P")
		ruleInfos = append(ruleInfos, *ruleInfo)
	}

	log.Info("ListRule", "ruleInfos", ruleInfos)

	return ruleInfos
}

func (jwtInfo *JWTInfo) AccessTree() AccessTree {
	accessTree := AccessTree{
		ServiceIds:   make(map[string]struct{}),
		ServiceTrees: make([]ServiceTree, 0),
	}

	for _, accessInfo := range jwtInfo.IDTokenClaims.GetAccessInfos() {
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
