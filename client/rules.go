package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/luyuhuang/subsocks/control/access"
	"github.com/luyuhuang/subsocks/control/rule"
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

type domainNode struct {
	ruleInfo *rule.Info
	wild     bool
	children map[string]*domainNode
}

func newDomainNode() *domainNode {
	return &domainNode{children: make(map[string]*domainNode)}
}

type ipNode struct {
	ruleInfo *rule.Info
	bits     []byte
	children [2]*ipNode
}

// Rules represents proxy rules
type Rules struct {
	domainTree *domainNode
	ipv4Tree   *ipNode
	ipv6Tree   *ipNode
	otherInfo  *rule.Info

	mu        sync.RWMutex
	isProxy   map[string]bool
	cacheFile *os.File

	watcher   *fsnotify.Watcher
	rulesPath string
	ruleMu    sync.RWMutex
}

func newRules() *Rules {
	return &Rules{
		isProxy: make(map[string]bool),
	}
}

// NewRulesFromStructMap create a Rules object from struct map
func NewRulesFromStructMap(rules map[string]*rule.Info) (*Rules, error) {
	r := newRules()
	if err := r.loadCache(); err != nil {
		return nil, err
	}

	r.ipv4Tree = new(ipNode)
	r.ipv6Tree = new(ipNode)
	r.domainTree = newDomainNode()
	r.otherInfo = rule.NewInfo(access.Info{}, service.Info{}, "D")

	for addr, ruleInfo := range rules {
		ruleLevel, ok := ruleString2Rule[ruleInfo.Rule]
		if !ok {
			return nil, fmt.Errorf("ruleInfo of %q got %s, want proxy|direct|auto|P|D|A", addr, ruleInfo.Rule)
		}

		ruleInfo.Level = ruleLevel

		if err := setRule(r.ipv4Tree, r.ipv6Tree, r.domainTree, r.otherInfo, addr, ruleInfo); err != nil {
			return nil, fmt.Errorf("set ruleInfo failed: %s", err)
		}
	}

	return r, nil
}

func (r *Rules) loadCache() error {
	f, err := os.OpenFile(".proxy-cache", os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		r.isProxy[strings.TrimSpace(s.Text())] = true
	}
	r.cacheFile = f

	return nil
}

func setRule(ipv4Tree, ipv6Tree *ipNode, domainTree *domainNode, other *rule.Info, addr string, rule *rule.Info) error {
	if addr == "*" {
		// * default is O, and can be set by ruleInfo(D)
		other = rule
	} else if ip := net.ParseIP(addr); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			setIPRule(ipv4Tree, ipv4, 32, rule)
		} else {
			setIPRule(ipv6Tree, ip.To16(), 128, rule)
		}
	} else if _, cidr, err := net.ParseCIDR(addr); err == nil {
		ones, _ := cidr.Mask.Size()
		if ipv4 := cidr.IP.To4(); ipv4 != nil {
			setIPRule(ipv4Tree, ipv4, ones, rule)
		} else {
			setIPRule(ipv6Tree, cidr.IP.To16(), ones, rule)
		}
	} else {
		if err := setDomainRule(domainTree, addr, rule); err != nil {
			return err
		}
	}

	return nil
}

func setIPRule(root *ipNode, ip []byte, length int, rule *rule.Info) {
	var p, pp *ipNode
	p = root
	j := 0
	for i := 0; i < length; i++ {
		b := (ip[i/8] >> (8 - i%8 - 1)) & 1

		if j >= len(p.bits) {
			j = 0
			pp = p
			p = p.children[b]
		}

		var pnode, node *ipNode
		if p == nil {
			pnode = pp
		} else if p.bits[j] != b {
			// p: |------+---------|
			//           ^ j
			//    |--np--|----p----|

			np := new(ipNode)
			np.bits = make([]byte, j)
			copy(np.bits, p.bits[:j])

			pp.children[np.bits[0]] = np

			copy(p.bits, p.bits[j:])
			p.bits = p.bits[:len(p.bits)-j]

			np.children[p.bits[0]] = p
			pnode = np
		} else if i == length-1 {
			if j == len(p.bits)-1 {
				node = p
			} else {
				np := new(ipNode)
				np.bits = make([]byte, j+1)
				copy(np.bits, p.bits[:j+1])

				pp.children[np.bits[0]] = np
				node = np

				copy(p.bits, p.bits[j+1:])
				p.bits = p.bits[:len(p.bits)-j-1]

				np.children[p.bits[0]] = p
			}
		}

		if pnode != nil {
			node = new(ipNode)
			node.bits = make([]byte, length-i)
			for k := i; k < length; k++ {
				node.bits[k-i] = (ip[k/8] >> (8 - k%8 - 1)) & 1
			}
			pnode.children[node.bits[0]] = node
		}

		if node != nil {
			node.ruleInfo = rule
			break
		}

		j++
	}
}

func setDomainRule(p *domainNode, domain string, rule *rule.Info) error {
	if i := strings.IndexByte(domain, '*'); i != 0 && i != -1 ||
		strings.Count(domain, "*") > 1 {
		return fmt.Errorf("Domain %q contains illegal wildcards", domain)
	}

	parts := strings.Split(domain, ".")

	for i := len(parts) - 1; i > 0; i-- {
		part := parts[i]
		if p.children[part] == nil {
			p.children[part] = newDomainNode()
		}
		p = p.children[part]
	}

	if part := parts[0]; part == "*" {
		p.ruleInfo = rule
		p.wild = true
	} else {
		if p.children[part] == nil {
			p.children[part] = newDomainNode()
		}
		ruleCopy := *rule
		p.children[part].ruleInfo = &ruleCopy
	}

	return nil
}

func searchIPRule(root *ipNode, ip []byte) (rule *rule.Info) {
	p := root
	j := 0
	for i := 0; i < len(ip)*8; i++ {
		b := (ip[i/8] >> (8 - i%8 - 1)) & 1

		if j >= len(p.bits) {
			j = 0
			p = p.children[b]
		}

		if p == nil || p.bits[j] != b {
			break
		}

		if j == len(p.bits)-1 && p.ruleInfo.Level != ruleNone {
			rule = p.ruleInfo
		}

		j++
	}
	return
}

func (r *Rules) getRule(addr string) (ruleInfo *rule.Info) {
	r.ruleMu.RLock()
	if ip := net.ParseIP(addr); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil { // IPv4
			ruleInfo = searchIPRule(r.ipv4Tree, ipv4)
		} else { // IPv6
			ruleInfo = searchIPRule(r.ipv6Tree, ip.To16())
		}
	} else {
		parts := strings.Split(addr, ".")
		p := r.domainTree

		for i := len(parts) - 1; i >= 0; i-- {
			part := parts[i]
			p = p.children[part]
			if p == nil {
				break
			}

			if p.wild || i == 0 {
				if p.ruleInfo.Level != ruleNone {
					ruleInfo = p.ruleInfo
				}
			}
		}
	}
	r.ruleMu.RUnlock()

	if ruleInfo == nil {
		log.Printf("failed to find ruleInfo for %s\n", addr)
		return rule.NewInfo(access.Info{}, service.Info{}, "D")
	}
	return
}

func (r *Rules) setAsProxy(addr string) {
	if r == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isProxy[addr] {
		r.isProxy[addr] = true
		r.cacheFile.WriteString(addr + "\n")
	}
}
