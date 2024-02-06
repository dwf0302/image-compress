package compressService

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const Out = "./temp/outGif"

type CompressService struct {
	GifsiclePath string
}

func compressGIF(sourceGifPath string, targetSize int64) ([]byte, error) {
	colorNum := 256
	metaColor := 256

	sourceGifFile, err := os.Stat(sourceGifPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("错误：读取gif文件出错 %s %v", sourceGifPath, err))
	}

	sourceSize := sourceGifFile.Size()

	if sourceSize < targetSize {
		fileData, err := ioutil.ReadFile(sourceGifPath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("警告：源gif文件小于目标大小，不做处理,复制文件失败 %s %v", sourceGifPath, err))
		}
		return fileData, nil
	}

	for {
		gifCmd := exec.Command("gifsicle", "--colors", strconv.Itoa(colorNum), "-O3", sourceGifPath, "-o", Out)
		resByte, err := gifCmd.CombinedOutput()
		if err != nil {
			resStr := string(resByte)
			return nil, errors.New(fmt.Sprintf("错误：压缩出错 %s %v %s", sourceGifPath, err, resStr))
		}

		outGifFile, err := os.Stat(Out)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("错误：读取gif文件出错 %s %v", Out, err))
		}

		outsize := outGifFile.Size()
		diffCurrent := targetSize - outsize
		if diffCurrent >= 0 && colorNum >= 256 {
			break
		}

		metaColor = metaColor / 2

		if metaColor <= 0 {
			if diffCurrent < 0 {
				colorNum = colorNum - 1
				continue
			}
			outGifName := filepath.Base(Out)
			outGifName = strings.Replace(outGifName, ".gif", "_"+strconv.Itoa(colorNum)+".gif", -1)
			if err := os.Rename(Out, filepath.Join(filepath.Dir(Out), outGifName)); err != nil {
				log.Println("重命名文件失败", Out)
			}
			break
		}
		if diffCurrent < 0 {
			colorNum = colorNum - metaColor
		} else {
			if colorNum >= 256 {
				colorNum = 0
			}
			colorNum = colorNum + metaColor
		}
	}
	fileData, err := ioutil.ReadFile(Out)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}
