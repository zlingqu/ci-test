package router

import (
	h "ci-test/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, 我是默认页")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/getHostname", h.GetHostname) //返回hostname

	r.GET("/getIp", h.GetIp) //返回ip

	r.GET("/getQrcode", h.GetQrcode) //返回二维码

	r.GET("/getReqInfo", h.GetRequestInfo) //返回请求的相关信息，比如header信息等

	r.GET("/getByte", h.GetReadfile) //返回[]byte类型的数据，可返回html内容的格式

	// 	router.LoadHTMLGlob("templates/") //全局加载templates/下一级模板文件
	// router.LoadHTMLGlob("templates/**/") //全局加载templates//下二级模板文件
	// router.LoadHTMLFiles("templates/a.html") //加载指定路径模板文件

	r.LoadHTMLGlob("files/*")
	r.GET("/getHtml", h.GetHtml) //Content-Type: text/html; charset=utf-8

	r.GET("/getXml", func(c *gin.Context) { //Content-Type: application/xml; charset=utf-8
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/getYaml", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/getJson", h.GetJson) //Content-Type: application/json; charset=utf-8
	r.GET("/getFile", h.GetFile) //Content-Type: application/json; charset=utf-8

	return r
}
