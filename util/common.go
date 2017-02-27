package util

import (
	"conf"
	"crypto/md5"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	_ "github.com/golang/image/webp"
	_ "image/gif"
	_ "image/png"
)

func LoadImage(imgPath string) (img image.Image, err error) {
	file, err := os.Open(conf.GetImgFullPath(imgPath))
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

func LoadFile(imgPath string) ([]byte, error) {
	return ioutil.ReadFile(conf.GetImgFullPath(imgPath))
}

func ParseImgArg(imgArg string) (int, int) {
	args := strings.Split(imgArg, "x")
	if len(args) != 2 {
		return 0, 0
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	return width, height
}

func Md5Sum(key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}
