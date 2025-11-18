package access

import "fmt"

type Info struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

// map: id ——> access_info
var infos = map[string]Info{
	"yyy1": {
		Name: "Ubuntu25",
		Host: "192.168.235.128",
		Port: 1080,
	},
	"yyy2": {
		Name: "Ubuntu22",
		Host: "192.168.235.133",
		Port: 1080,
	},
}

func GetInfoById(id string) (Info, error) {
	if _, ok := infos[id]; !ok {
		return Info{}, fmt.Errorf("invalid access info id: %s", id)
	}

	return infos[id], nil
}
