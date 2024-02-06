package compressService

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ImageCompress(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("表单解析错误: %s", err.Error()))
		return
	}
	size := c.PostForm("size")
	fmt.Println("文本数据:", size)

	dst := filepath.Join("./temp", file.Filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("文件上传错误: %s", err.Error()))
		return
	}
	defer func() {
		if err := os.Remove(dst); err != nil {
			log.Println("警告：删除临时文件失败", dst)
		}
	}()

	targetSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("获取图片大小失败: %s", err.Error()))
		return
	}
	targetSize = targetSize * 1000
	switch strings.ToLower(filepath.Ext(file.Filename)) {
	case ".png":
		data, err := compressPNG(dst, targetSize)
		handleCompressionResult(c, err, data, "image/png")
	case ".jpg", ".jpeg":
		data, err := compressJPG(dst, targetSize)
		handleCompressionResult(c, err, data, "image/jpeg")
	case ".gif":
		data, err := compressGIF(dst, targetSize)
		handleCompressionResult(c, err, data, "image/gif")
	default:
		log.Println("不支持的文件类型:", filepath.Ext(file.Filename))
		c.String(http.StatusNotImplemented, "不支持的文件类型!")
	}
}

func handleCompressionResult(c *gin.Context, err error, data []byte, contentType string) {
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("图片压缩失败: %s", err.Error()))
		return
	}
	c.Data(http.StatusOK, contentType, data)
}
