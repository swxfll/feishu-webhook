package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func main() {
	r := gin.Default()

	r.POST("/feishu-webhook", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			// 错误处理
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		fmt.Println(string(body)) // 打印原始 JSON 数据到控制台

		// 其他处理逻辑...

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
