package service

import "fmt"

type Info struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

// map: id ——> service_info
var infos = map[string]Info{
	"zzz1": {
		Name: "IPSB",
		Host: "ip.sb",
		Port: 443,
	},
	"zzz2": {
		Name: "YouTube",
		Host: "www.youtube.com",
		Port: 443,
	},
}

func GetInfoById(id string) (Info, error) {
	if _, ok := infos[id]; !ok {
		return Info{}, fmt.Errorf("invalid access info id: %s", id)
	}

	return infos[id], nil
}
