package util

import (
	"crypto/md5"
	"encoding/hex"
	"image"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	if strings.HasPrefix(strings.TrimLeft(imgPath, "/"), AdImgPath) {
		imgs, err := ioutil.ReadDir(ImgRoot + AdImgPath)
		if err != nil || len(imgs) == 0 {
			return nil, err
		}
		randomImg := imgs[rand.Intn(len(imgs))]
		return ioutil.ReadFile(ImgRoot + AdImgPath + randomImg.Name())
	}
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

func Md5Sum(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
