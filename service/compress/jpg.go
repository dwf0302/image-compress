package compressService

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func compressJPG(sourcePath string, targetSizeKB int64) ([]byte, error) {
	quality := 99
	for {
		// 使用 djpegOut 进行解码
		inputFile, err := os.ReadFile(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file", err)
		}
		// 调用 mozjpeg 命令行工具解码图片
		djpegCmd := exec.Command("djpeg")
		djpegCmd.Stdin = bytes.NewReader(inputFile)
		var djpegOut bytes.Buffer
		djpegCmd.Stdout = &djpegOut
		err = djpegCmd.Run()
		if err != nil {
			return nil, fmt.Errorf("error decoding image", err)
		}
		// 调用 mozjpeg 命令行工具压缩图片
		cjpegCmd := exec.Command("cjpeg", "-quality", fmt.Sprintf("%d", quality), "-optimize")
		cjpegCmd.Stdin = &djpegOut
		var out, errBuffer bytes.Buffer
		cjpegCmd.Stdout = &out
		cjpegCmd.Stderr = &errBuffer

		// 执行命令
		if err := cjpegCmd.Run(); err != nil {
			return nil, fmt.Errorf("command execution failed: %v, stderr: %s", err, errBuffer.String())
		}
		// 获取压缩后图像的大小
		compressedSize := int64(len(out.Bytes())) // 转换为 KB

		// 如果压缩后的图像大小接近目标大小或者已经小于目标大小，则返回压缩后的图像数据
		if compressedSize <= targetSizeKB || quality <= 0 {
			return out.Bytes(), nil
		}

		// 减少压缩质量，继续尝试
		quality -= 1
		if quality <= 0 {
			return nil, fmt.Errorf("failed to compress JPG image to %d KB", targetSizeKB)
		}
	}
}
