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

//func compressPNG(sourcePath string, targetSize int64) ([]byte, error) {
//	file, err := os.Open(sourcePath)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	img, err := png.Decode(file)
//	if err != nil {
//		return nil, err
//	}
//
//	quality := 100 // 初始 JPEG 压缩质量
//	for {
//		// 创建缓冲区以保存压缩后的图像
//		var buf bytes.Buffer
//		// 将图像以指定的质量编码为 JPEG 格式
//		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
//		if err != nil {
//			return nil, err
//		}
//		// 检查压缩后图像的大小是否满足目标大小
//		if int64(buf.Len()) <= targetSize {
//			return buf.Bytes(), nil
//		}
//		// 如果超过目标大小，则增加质量，重试压缩
//		quality -= 5
//		if quality < 0 {
//			return nil, err
//		}
//	}
//}

//func compressPNG(sourcePath string, targetSize int64) ([]byte, error) {
//	file, err := os.Open(sourcePath)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	img, _, err := image.Decode(file)
//	if err != nil {
//		return nil, err
//	}
//
//	targetWidth := findTargetWidth(img, targetSize)
//	targetHeight := uint(float64(img.Bounds().Dy()) * (float64(targetWidth) / float64(img.Bounds().Dx())))
//	resizedImg := resize.Resize(targetWidth, targetHeight, img, resize.Lanczos3)
//
//	var buf bytes.Buffer
//	err = png.Encode(&buf, resizedImg)
//	if err != nil {
//		return nil, err
//	}
//	return buf.Bytes(), nil
//}
//
//func findTargetWidth(img image.Image, targetSize int64) uint {
//	lower := uint(1)
//	upper := uint(img.Bounds().Dx())
//
//	for lower <= upper {
//		mid := (lower + upper) / 2
//		resizedImg := resize.Resize(mid, 0, img, resize.Lanczos3)
//
//		var buf bytes.Buffer
//		_ = png.Encode(&buf, resizedImg)
//
//		currentSize := int64(buf.Len())
//
//		if currentSize == targetSize {
//			return mid
//		} else if currentSize < targetSize {
//			lower = mid + 1
//		} else {
//			upper = mid - 1
//		}
//	}
//	return upper
//}
