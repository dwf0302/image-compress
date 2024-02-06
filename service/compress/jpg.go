package compressService

import (
	"bytes"
	"fmt"
	"os/exec"
)

//func compressJPG(sourcePath string, targetSize int64) ([]byte, error) {
//	file, err := os.Open(sourcePath)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	img, err := jpeg.Decode(file)
//	if err != nil {
//		return nil, err
//	}
//
//	quality := 100
//	for {
//		compressedImage := make(chan []byte)
//		go func() {
//			compressedImage <- compressImage(img, quality)
//		}()
//
//		select {
//		case imgData := <-compressedImage:
//			if int64(len(imgData)) <= targetSize {
//				return imgData, nil
//			}
//			quality -= 3
//			if quality <= 0 {
//				return nil, errors.New("压缩质量过低，无法满足目标大小")
//			}
//		}
//	}
//}
//
//func compressImage(img image.Image, quality int) []byte {
//	var buf bytes.Buffer
//	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
//	if err != nil {
//		log.Println(err)
//		return nil
//	}
//	return buf.Bytes()
//}

//func compressPNG(sourcePath string, targetSizeKB int64) ([]byte, error) {
//	// 初始压缩质量
//	quality := 95
//
//	for {
//		// 使用 pngquant-linux 进行 PNG 图像压缩
//		cmd := exec.Command(pngquant, fmt.Sprintf("--quality=%d", quality), "--force", "--output", "-", sourcePath)
//		var out bytes.Buffer
//		cmd.Stdout = &out
//		err := cmd.Run()
//		if err != nil {
//			return nil, err
//		}
//
//		// 获取压缩后图像的大小
//		compressedSize := int64(len(out.Bytes()) / 1024) // 转换为 KB
//
//		// 如果压缩后的图像大小接近目标大小或者已经小于目标大小，则返回压缩后的图像数据
//		if compressedSize <= targetSizeKB || quality <= 0 {
//			return out.Bytes(), nil
//		}
//		// 减少压缩质量，继续尝试
//		quality -= 5
//		if quality <= 0 {
//			return nil, fmt.Errorf("failed to compress PNG image to %d KB", targetSizeKB)
//		}
//	}
//}

func compressJPG(sourcePath string, targetSizeKB int64) ([]byte, error) {
	quality := 95
	for {
		//使用 pngquant-linux 进行 PNG 图像压缩
		cmd := exec.Command("cjpeg", fmt.Sprintf("--quality=%d", quality), "--force", "--output", "-", sourcePath)
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
