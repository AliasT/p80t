package main

import (
	"encoding/json"
	"fmt"
	"github.com/cbednarski/hostess"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// Server 配置文件的结构表示
type Server struct {
	name string
	host string
	port float64
}

// BaseHandle proxy base
type BaseHandle struct {
	servers *[]Server
}

// ServeHTTP
// https://stackoverflow.com/questions/21055182/golang-reverse-proxy-with-multiple-apps
func (h *BaseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Bridge", "Golang Http Proxy")
	for _, server := range *h.servers {
		// 如果请求的主机地址符合条件, 转发至本地端口
		if r.Host == server.host {
			port := int(server.port)
			target := fmt.Sprintf("http://localhost:%d", port)
			remote, _ := url.Parse(target)
			proxy := httputil.NewSingleHostReverseProxy(remote)
			proxy.ServeHTTP(w, r)
			break
		}
	}

	w.Write([]byte("没有符合条件的转发 ！"))
}

// ReadJSONConfig 读写本地转发列表servers.json
// 将列表中的域名动态写入hosts文件
func ReadJSONConfig(hostfile *hostess.Hostfile) (*[]Server, error) {
	var servers map[string]interface{}
	var structuredServers []Server

	file, _ := os.Open("servers.json")
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &servers)

	for name, s := range servers {
		server := s.(map[string]interface{})
		host := server["host"].(string)
		port := server["port"].(float64)
		ip := net.IPv4(127, 0, 0, 1)

		hostname := hostess.Hostname{
			Domain:  host,
			IP:      ip,
			Enabled: true,
			IPv6:    false,
		}

		hostfile.Hosts.Add(&hostname)
		structuredServers = append(structuredServers, Server{name, host, port})
	}

	hostfile.Save()

	return &structuredServers, nil
}

func main() {
	hostsfile, _ := hostess.LoadHostfile()
	servers, _ := ReadJSONConfig(hostsfile)

	h := &BaseHandle{servers}
	http.Handle("/", h)

	server := &http.Server{
		Addr:    ":80",
		Handler: h,
	}

	log.Fatal(server.ListenAndServe())
}
