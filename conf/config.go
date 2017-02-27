package conf

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
)

const (
	adDefaultDuration = 6
)

var gConfig = &Config{}
var (
	configFile string
	listenPort int
)

type Config struct {
	Spiders      []string
	AllowedRefer []string
	CachePath    string
	ImgPath      string
	ExtImgSize   string
	AdImgPath    string
	AdDuration   map[string]int
	LogFile      string
}

func init() {
	flag.StringVar(&configFile, "c", "conf/config.yaml", "conf file path")
	flag.IntVar(&listenPort, "p", 6789, "listen port")
	flag.Parse()

	reloadConf(configFile)
}

func IsAllowedRefer(refer string) bool {
	for _, allowedRefer := range gConfig.AllowedRefer {
		if strings.Contains(refer, allowedRefer) {
			return true
		}
	}
	return false
}

func IsSpider(ua string) bool {
	for _, spiderKey := range gConfig.AllowedRefer {
		if strings.Contains(ua, spiderKey) {
			return true
		}
	}
	return false
}

func GetImgFullPath(imgPath string) string {
	return path.Join(gConfig.ImgPath, imgPath)
}

func GetCacheDir() string {
	return gConfig.CachePath
}

func GetAdFullPath() string {
	return path.Join(gConfig.ImgPath, gConfig.AdImgPath)
}

func GetAdDuration(host string) int {
	if duration, ok := gConfig.AdDuration[host]; ok {
		return duration
	}
	return adDefaultDuration
}

func GetRandomAdPath() string {
	adImgPath := path.Join(gConfig.ImgPath, gConfig.AdImgPath)
	imgs, err := ioutil.ReadDir(adImgPath)
	if err != nil || len(imgs) == 0 {
		return ""
	}
	return path.Join(gConfig.AdImgPath, imgs[rand.Intn(len(imgs))].Name())
}

func GetLogFile() string {
	return gConfig.LogFile
}

func GetExtImgSize() string {
	return gConfig.ExtImgSize
}

func GetListenPort() int {
	return listenPort
}

func reloadConf(configFile string) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicf("read config file err:%v\n", err)
	}

	err = yaml.Unmarshal(configData, gConfig)
	if err != nil {
		log.Panicf("parse config file err:%v\n", err)
	}

	log.Printf("config:%+v", gConfig)
	go reloadYamlFile(configFile, time.Minute, gConfig)
}

func reloadYamlFile(configFile string, duration time.Duration, serverConf *Config) {
	var lastMtime = getFileMtime(configFile)
	for {
		time.Sleep(duration)
		if curMtime := getFileMtime(configFile); curMtime > lastMtime {
			lastMtime = curMtime
			configData, err := ioutil.ReadFile(configFile)
			if err != nil {
				log.Panicf("read config file err:%v\n", err)
			}
			err = yaml.Unmarshal(configData, &serverConf)
			if err != nil {
				log.Panicf("parse config file err:%v\n", err)
			}
			log.Printf("config:%+v", serverConf)
		}
	}
}

func getFileMtime(file string) int64 {
	fileInfo, err := os.Stat(file)
	if err != nil {
		log.Fatalf("file stat err:%v\n", err)
		return 0
	}
	return fileInfo.ModTime().Unix()
}
