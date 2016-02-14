package util

import (
	"image"
	"os"
	"strconv"
	"strings"
)

func LoadImage(imgName, imgDate, category string) (img image.Image, err error) {
	file, err := os.Open(ImgRoot + "/" + category + "/" + imgDate + "/" + imgName)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
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
