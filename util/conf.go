package util

import (
	"flag"
	"strings"
)

var (
	CacheRoot    = "./cache/"
	ImgRoot      = "./public/"
	ServePort    = "6789"
	AllowedRefer = "127.0.0.1"
	LogFile      = "log"
)

func init() {
	flag.StringVar(&CacheRoot, "cPath", "./cache/", "Cache path")
	flag.StringVar(&ImgRoot, "sPath", "./public/", "Source image path")
	flag.StringVar(&ServePort, "port", "6789", "Server port")
	flag.StringVar(&AllowedRefer, "aRefer", "127.0.0.1", "Allowed refer")
	flag.StringVar(&LogFile, "log", "log", "log file pre name")
	flag.Parse()

	if !strings.HasSuffix(ImgRoot, "/") {
		ImgRoot += "/"
	}

	if !strings.HasSuffix(CacheRoot, "/") {
		CacheRoot += "/"
	}
}
