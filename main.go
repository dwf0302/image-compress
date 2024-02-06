package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	compressService "utils/service/compress"
)

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()

	// 加载 HTML 模板文件
	r.LoadHTMLGlob("templates/*")
	// 指定静态文件目录
	r.Static("/static", "./static")
	// 定义路由处理函数
	r.GET("/index", func(c *gin.Context) {
		// 渲染 HTML 模板并发送给客户端
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "index",
			// 可以在 HTML 模板中使用这些数据
			//"imageURL": "/static/2178859284299776000.heic", // 静态文件的路径
		})
	})
	// 设置路由处理器，处理文件上传请求
	r.POST("/upload", compressService.ImageCompress)
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run(":8199")
}
