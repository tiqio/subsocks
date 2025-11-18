package rule

import (
	"log"

	"github.com/luyuhuang/subsocks/control/access"
	"github.com/luyuhuang/subsocks/control/service"
)

const (
	ruleNone = iota
	ruleProxy
	ruleDirect
	ruleAuto
)

var ruleString2Rule = map[string]int{
	"proxy":  ruleProxy,
	"direct": ruleDirect,
	"auto":   ruleAuto,
	"P":      ruleProxy,
	"D":      ruleDirect,
	"A":      ruleAuto,
}

type Info struct {
	AccessInfo  access.Info  `json:"access_info"`
	ServiceInfo service.Info `json:"service_info"`
	Rule        string       `json:"rule"`
	Level       int          `json:"level"`
}

func NewInfo(accessInfo access.Info, serviceInfo service.Info, rule string) *Info {

	ruleLevel, ok := ruleString2Rule[rule]
	if !ok {
		log.Fatalf("NewInfo got %s, want proxy|direct|auto|P|D|A", rule)
	}

	return &Info{
		AccessInfo:  accessInfo,
		ServiceInfo: serviceInfo,
		Rule:        rule,
		Level:       ruleLevel,
	}
}
