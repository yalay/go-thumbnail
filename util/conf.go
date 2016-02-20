package util

import (
	"flag"
	"strings"
)

var (
	CacheRoot = "./cache/"
	ImgRoot   = "./public/"
	ServePort = "6789"
)

func init() {
	flag.StringVar(&CacheRoot, "cPath", "./cache/", "cache path")
	flag.StringVar(&ImgRoot, "sPath", "./public/", "source image path")
	flag.StringVar(&ServePort, "port", "6789", "server port")
	flag.Parse()

	if !strings.HasSuffix(ImgRoot, "/") {
		ImgRoot += "/"
	}

	if !strings.HasSuffix(CacheRoot, "/") {
		CacheRoot += "/"
	}
}
