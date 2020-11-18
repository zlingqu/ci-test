package main

import (
	"net"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	hostname, _ := os.Hostname()
	addrs, _ := net.InterfaceAddrs()
	var ip string
	for _, address := range addrs {
		// 检查ip地址格式
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ip + ipnet.IP.String()
			}

		}
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, 我是默认页")
	})

	r.GET("/hostname", func(c *gin.Context) {
		c.String(200, hostname)
	})

	r.GET("/ip", func(c *gin.Context) {
		c.String(200, ip)
	})

	r.Run(":" + "80") // listen and serve on 0.0.0.0:80
}
