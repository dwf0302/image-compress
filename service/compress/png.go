package compressService

import (
	"bytes"
	"fmt"
	"os/exec"
)

const pngquant = "/Users/devin/Documents/Weifeng/code/go/src/utils/bin/pngquant-linux"

func compressPNG(sourcePath string, targetSizeKB int64) ([]byte, error) {
	// 初始压缩质量
	quality := 95

	for {
		// 使用 pngquant-linux 进行 PNG 图像压缩
		cmd := exec.Command(pngquant, fmt.Sprintf("--quality=%d", quality), "--force", "--output", "-", sourcePath)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return nil, err
		}

		// 获取压缩后图像的大小
		compressedSize := int64(len(out.Bytes()) / 1024) // 转换为 KB

		// 如果压缩后的图像大小接近目标大小或者已经小于目标大小，则返回压缩后的图像数据
		if compressedSize <= targetSizeKB || quality <= 0 {
			return out.Bytes(), nil
		}
		// 减少压缩质量，继续尝试
		quality -= 5
		if quality <= 0 {
			return nil, fmt.Errorf("failed to compress PNG image to %d KB", targetSizeKB)
		}
	}
}
