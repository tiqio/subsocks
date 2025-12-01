package access

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

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

// AccessInfo 定义结构体来映射 JSON 数据
type AccessInfo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Addr      string `json:"addr"`
	CaInfoId  *int   `json:"ca_info_id"`
	Timestamp string `json:"timestamp"`
}

// GetAccessInfo 获取数据的函数
func GetAccessInfo(name string) (*AccessInfo, error) {
	url := fmt.Sprintf("http://127.0.0.1:8000/myapp/info/access/%s", name)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var accessInfo AccessInfo
	if err := json.Unmarshal(body, &accessInfo); err != nil {
		return nil, err
	}

	return &accessInfo, nil
}

func GetInfoById(id string) (Info, error) {
	if info, ok := infos[id]; ok {
		return info, nil
	}

	accessInfo, err := GetAccessInfo(id)
	if err != nil {
		return Info{}, err
	}

	// 拆分 addr
	addrParts := strings.Split(accessInfo.Addr, ":")
	if len(addrParts) != 2 {
		return Info{}, fmt.Errorf("invalid address format: %s", accessInfo.Addr)
	}

	// 创建并返回Info结构体
	port, _ := strconv.Atoi(addrParts[1]) // 假设 addr 的格式是 "host:port"
	return Info{
		Name: accessInfo.Name,
		Host: addrParts[0],
		Port: port,
	}, nil
}
