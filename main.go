package main

import (
	"encoding/json"
	"fmt"
	"github.com/cbednarski/hostess"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

// Server 配置文件的结构表示
type Server struct {
	name string
	host string
	port float64
}

func main() {
	hostsfile, _ := hostess.LoadHostfile()
	servers, _ := ReadJSONConfig()
	Serve(hostsfile, servers)
}

// Serve 本地运行的gin服务器
func Serve(hostfile *hostess.Hostfile, servers *[]Server) {
	r := gin.Default()
	r.Use(Transfer(hostfile, servers))
	r.Run(":80")
}

// Transfer 将不同域名的请求转发至配置文件的指定端口
func Transfer(hostfile *hostess.Hostfile, servers *[]Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var url string
		var port int
		var res *http.Response

		for _, server := range *servers {
			if c.Request.Host == server.host {
				port = int(server.port)
				url = fmt.Sprintf("http://localhost:%d%s", port, c.Request.RequestURI)
				break
			}
		}

		res, _ = http.Get(url)
		defer res.Body.Close()

		if res.StatusCode == 404 {
			res, _ = http.Get(fmt.Sprintf("http://localhost:%d", port))
			defer res.Body.Close()
		}

		bytes, _ := ioutil.ReadAll(res.Body)
		c.Writer.WriteHeader(200)
		c.Writer.Write(bytes)

		c.Abort()
		return
	}
}

// ReadJSONConfig 读写本地转发列表servers.json
func ReadJSONConfig() (*[]Server, error) {
	var servers map[string]interface{}
	var structuredServers []Server

	file, _ := os.Open("servers.json")
	bytes, _ := ioutil.ReadAll(file)

	json.Unmarshal(bytes, &servers)
	for name, s := range servers {
		server := s.(map[string]interface{})
		structuredServers = append(structuredServers, Server{name, server["host"].(string), server["port"].(float64)})
	}

	defer file.Close()

	return &structuredServers, nil
}
