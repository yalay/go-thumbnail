package util

import (
	"flag"
)

var (
	CacheRoot = "./cache/"
	ImgRoot   = "./public/"
)

func init() {
	flag.StringVar(&CacheRoot, "cPath", "./cache/", "cache path")
	flag.StringVar(&ImgRoot, "sPath", "./public/", "source image path")
}
