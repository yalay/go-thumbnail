package util

import (
	"crypto/md5"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	_ "image/png"
)

func LoadImage(imgPath string) (img image.Image, err error) {
	file, err := os.Open(ImgRoot + imgPath)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

func LoadFile(imgPath string) ([]byte, error) {
	return ioutil.ReadFile(ImgRoot + imgPath)
}

func ParseImgArg(imgArg string) (uint, uint) {
	args := strings.Split(imgArg, "x")
	if len(args) != 2 {
		return 0, 0
	}

	width, _ := strconv.Atoi(args[0])
	height, _ := strconv.Atoi(args[1])
	return uint(width), uint(height)
}

func Md5Sum(key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}
