package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

type ReqQrcode struct {
	Url string `form:"url"`
}

func GetQrcode(c *gin.Context) { //获取二维码
	var reqQrcode ReqQrcode
	var png []byte
	c.ShouldBind(&reqQrcode)

	png, _ = qrcode.Encode(reqQrcode.Url, qrcode.Medium, 256)
	c.String(http.StatusOK, string(png))
}

func GetRequestInfo(c *gin.Context) { //返回请求的一些信息

	var head string
	// var r c.Request
	for k, v := range c.Request.Header { //将head转换为string
		var value string
		for _, j := range v {
			value = value + j
		}
		head = head + k + ":" + value + "\n"
	}
	body, _ := ioutil.ReadAll(c.Request.Body) //获取body

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	respon :=
		"本次请求客户端的IP和端口是:" + c.Request.RemoteAddr + "\n" +
			"请求完整的url是:" + strings.Join([]string{scheme, "://", c.Request.Host, c.Request.RequestURI}, "") + "\n" +
			"请求协议是:" + scheme + "\n" +
			"请求方式是:" + c.Request.Method + "\n" +
			"请求path是：" + c.Request.URL.Path + "\n" +
			"请求的http版本是:" + c.Request.Proto + "\n" +
			"请求host是：" + c.Request.Host + "\n" +
			"请求RequestURI是：" + c.Request.RequestURI + "\n" +
			"请求Referer是：" + c.Request.Referer() + "\n" + //显示上一跳的信息，可用于防盗链、网站流量来源分析等领域
			"请求header如下：\n" + head +
			"请求RawQuery是：\n" + c.Request.URL.RawQuery + "\n" +
			"请求body是:" + string(body) + "\n"

	c.Writer.Header().Add("name", "zlingqu") //返回的头部中，添加一个特殊的
	c.String(http.StatusOK, respon)

}

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

	r.GET("/qrcode", GetQrcode)

	r.GET("/req-info", GetRequestInfo)

	r.Run(":" + "80") // listen and serve on 0.0.0.0:80
}
