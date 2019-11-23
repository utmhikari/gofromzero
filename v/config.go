package v

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

var (
	InvalidPortError = errors.New("端口号需要在5000~25000之间！")
	NoServerError = errors.New("没有合法配置服务器")
)

const (
	defaultConfigPath = "./config.json"
)

type config struct {
	Port uint `json:"port"`
	Servers []*Server `json:"servers"`
}

func initLoadBalancer(configPath string) (*loadBalancer, error) {
	lb := &loadBalancer{
		ServerMap: make(map[uint]*Server),
	}
	// load json from file
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return lb, err
	}
	var config config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return lb, err
	}
	// check listen port
	if config.Port < 5000 || config.Port > 25000 {
		return lb, InvalidPortError
	}
	lb.Port = config.Port
	// filter servers, yet without connection
	if len(config.Servers) == 0 {
		return lb, NoServerError
	}
	var id uint = 1
	for _, server := range config.Servers {
		if server.Weight <= 0 {
			server.Weight = weightDefault
		}
		lb.ServerMap[id] = server
		lb.ServerMap[id].Active = true
		id++
	}
	lb.PrintServers()
	return lb, nil
}
