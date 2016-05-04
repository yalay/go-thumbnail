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
	flag.StringVar(&CacheRoot, "cPath", "./cache/", "Cache path")
	flag.StringVar(&ImgRoot, "sPath", "./public/", "Source image path")
	flag.StringVar(&ServePort, "port", "6789", "Server port")
	flag.Parse()

	if !strings.HasSuffix(ImgRoot, "/") {
		ImgRoot += "/"
	}

	if !strings.HasSuffix(CacheRoot, "/") {
		CacheRoot += "/"
	}
}
