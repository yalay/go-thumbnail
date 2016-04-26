package util

import (
	"flag"
	"strings"
)

var (
	CacheRoot    = "./cache/"
	ImgRoot      = "./public/"
	ServePort    = "6789"
	ExtImgSize   = "200x300" // 外链图片尺寸大小
	WaterMarkImg = "water.png"
	WaterSize    = "1x1" // 1x1用来标记添加水印
	RedirectUrl  = "http://localhost"
)

var (
	Spiders      = []string{"Baiduspider", "Googlebot", "360Spider"}
	AllowedRefer = "127.0.0.1"
)

func init() {
	flag.StringVar(&CacheRoot, "cPath", "./cache/", "Cache path")
	flag.StringVar(&ImgRoot, "sPath", "./public/", "Source image path")
	flag.StringVar(&WaterMarkImg, "wImg", "water.png", "Water mark image")
	flag.StringVar(&ServePort, "port", "6789", "Server port")
	flag.StringVar(&RedirectUrl, "rurl", "http://localhost", "Redirect url")
	flag.StringVar(&AllowedRefer, "aRefer", "127.0.0.1", "Allowed refer")
	flag.Parse()

	if !strings.HasSuffix(ImgRoot, "/") {
		ImgRoot += "/"
	}

	if !strings.HasSuffix(CacheRoot, "/") {
		CacheRoot += "/"
	}
}
