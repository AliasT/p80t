package main

import (
	"encoding/json"
	"fmt"
	"github.com/cbednarski/hostess"
	"github.com/gin-gonic/gin"
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

	Serve()
}

func Serve() {
	r := gin.Default()
	r.Use(Transfer())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

// middleware
func Transfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Write([]byte("hello"))
		c.Abort()
		return
		// c.Next()
	}
}

// 读写本地转发列表servers.json

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
