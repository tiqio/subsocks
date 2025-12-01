package service

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

// ServiceInfo 定义结构体来映射 JSON 数据
type ServiceInfo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Addr      string `json:"addr"`
	Protocol  string `json:"protocol"`
	Timestamp string `json:"timestamp"`
}

// GetServiceInfo 获取数据的函数
func GetServiceInfo(name string) (*ServiceInfo, error) {
	url := fmt.Sprintf("http://127.0.0.1:8000/myapp/info/service/%s", name)

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

	var serviceInfo ServiceInfo
	if err := json.Unmarshal(body, &serviceInfo); err != nil {
		return nil, err
	}

	return &serviceInfo, nil
}

func GetInfoById(id string) (Info, error) {
	if info, ok := infos[id]; ok {
		return info, nil
	}

	serviceInfo, err := GetServiceInfo(id)
	if err != nil {
		return Info{}, err
	}

	// 拆分 addr
	addrParts := strings.Split(serviceInfo.Addr, ":")
	if len(addrParts) != 2 {
		return Info{}, fmt.Errorf("error: invalid address format: %s", serviceInfo.Addr)
	}

	// 创建并返回 ServiceInfo结构体
	port, _ := strconv.Atoi(addrParts[1]) // 假设 addr 的格式是 "host:port"
	return Info{
		Name: serviceInfo.Name,
		Host: addrParts[0],
		Port: port,
	}, nil
}
