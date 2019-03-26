package main

// 读写本地配置文件json

import (
	"encoding/json"
	"fmt"
	"github.com/cbednarski/hostess"
	"io/ioutil"
	"net"
	"os"
)

type Server struct {
	host []byte
	port int
}

func main() {
	hostsfile, errs := hostess.LoadHostfile()
	if len(errs) > 0 {

	}

	hostname := hostess.Hostname{"localhost", net.IPv4(127, 0, 0, 1), true, false}

	servers, err := ReadJSONConfig()
	if err != nil {

	}

	for name, s := range servers {
		server := s.(map[string]interface{})
		fmt.Println(server["host"], server["port"], name)

		if hostsfile.Hosts.Contains(&hostname) {
			fmt.Println("exists a server refer to local")
		}
	}
}

func ReadJSONConfig() (map[string]interface{}, error) {
	file, err := os.Open("servers.json")
	var servers map[string]interface{}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		//
	}

	json.Unmarshal(bytes, &servers)
	defer file.Close()

	return servers, nil

}
