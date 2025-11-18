package rule

import (
	"github.com/luyuhuang/subsocks/control/access"
	"github.com/luyuhuang/subsocks/control/service"
)

type Info struct {
	AccessInfo  access.Info  `json:"access_info"`
	ServiceInfo service.Info `json:"service_info"`
	Rule        string       `json:"rule"`
	Level       int          `json:"level"`
}

func NewInfo(accessInfo access.Info, serviceInfo service.Info, rule string) *Info {
	return &Info{
		AccessInfo:  accessInfo,
		ServiceInfo: serviceInfo,
		Rule:        rule,
	}
}
