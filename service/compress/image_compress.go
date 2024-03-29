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
	if file.Size > 10*1024*1024 {
		c.String(http.StatusBadRequest, fmt.Sprintf("文件大小不允许超过10M: %s", err.Error()))
		return
	}
	size := c.PostForm("size")
	fmt.Println("期望压缩后的文本大小:", size)
	targetSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("获取图片大小失败: %s", err.Error()))
		return
	}
	targetSize = targetSize * 1000
	if targetSize >= file.Size {
		c.String(http.StatusBadRequest, fmt.Sprintf("期望压缩后的文本大小大于文件大小: %s", err.Error()))
		return
	}
	fileAbs, err := filepath.Abs(filepath.Join("./temp", file.Filename))
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("获取文件绝对路径失败: %s", err.Error()))
		return
	}
	fmt.Println("绝对路径=" + fileAbs)
	if err := c.SaveUploadedFile(file, fileAbs); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("临时文件上传错误: %s", err.Error()))
		return
	}
	defer func() {
		if err := os.Remove(fileAbs); err != nil {
			log.Println("警告：删除临时文件失败", fileAbs)
		}
	}()

	switch strings.ToLower(filepath.Ext(file.Filename)) {
	case ".png":
		data, err := compressPNG(fileAbs, targetSize)
		handleCompressionResult(c, err, data, "image/png")
	case ".jpg", ".jpeg":
		data, err := compressJPG(fileAbs, targetSize)
		handleCompressionResult(c, err, data, "image/jpeg")
	case ".gif":
		data, err := compressGIF(fileAbs, targetSize)
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
