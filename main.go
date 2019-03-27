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
	// hostname := hostess.Hostname{"localhost", net.IPv4(127, 0, 0, 1), true, false}

	servers, _ := ReadJSONConfig()

	Serve(hostsfile, servers)
}

// Serve 本地运行的gin服务器
func Serve(hostfile *hostess.Hostfile, servers *[]Server) {
	r := gin.Default()
	r.Use(Transfer(hostfile, servers))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// listen and serve on 0.0.0.0:8080
	r.Run(":80")
}

// Transfer 将不同域名的请求转发至配置文件的指定端口
func Transfer(hostfile *hostess.Hostfile, servers *[]Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		target := "https://baidu.com"
		if c.Request.Host == "gui.xiaoma.cn" {
			target = "https://sohu.com"
		}
		res, _ := http.Get(target)
		bytes, _ := ioutil.ReadAll(res.Body)
		c.Writer.Write(bytes)
		c.Abort()
		return
		// c.Next()
	}
}

// ReadJSONConfig 读写本地转发列表servers.json
func ReadJSONConfig() (*[]Server, error) {
	file, err := os.Open("servers.json")
	var servers map[string]interface{}
	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		//
	}

	json.Unmarshal(bytes, &servers)
	structuredServers := make([]Server, len(servers))
	for name, s := range servers {
		server := s.(map[string]interface{})
		_ = append(structuredServers, Server{name, server["host"].(string), server["port"].(float64)})
	}

	fmt.Print(structuredServers)

	defer file.Close()

	return &structuredServers, nil
}
