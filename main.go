package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, 我是默认页")
	})
	r.Run(":" + "80") // listen and serve on 0.0.0.0:80
}
