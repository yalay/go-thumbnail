package util

import (
	"bytes"
	"os"
	"time"
)

var (
	logBuffer     = &bytes.Buffer{}
	writeDutation = 5 * time.Minute
)

func init() {
	go writeLog()
}

func GetLogBuf() *bytes.Buffer {
	return logBuffer
}

func Log(msg string) {
	logBuffer.WriteString(msg)
}

func writeLog() {
	for {
		todady := time.Now().Format("2006-01-02")
		file, err := os.OpenFile(LogFile+"."+todady, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			break
		}
		if logBuffer.Len() > 0 {
			logBuffer.WriteTo(file)
		}
		file.Close()
		time.Sleep(writeDutation)
	}
}
