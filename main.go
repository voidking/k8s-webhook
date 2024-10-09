package main

import (
	"fmt"
	"io"
	"k8s-webhook/internal/config"
	"k8s-webhook/internal/server"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	conf := config.GetConfig()
	// 设置全局时区
	loc, err := time.LoadLocation(conf.TimeZone) // 载入指定时区
	if err != nil {
		panic(err)
	}
	time.Local = loc
	// 设置全局日志配置
	if conf.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	logPath := path.Join(conf.LogPath, fmt.Sprintf("main-%s.log",
		time.Now().Format("20060102150401")))
	dir := filepath.Dir(logPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, file)
		logrus.SetOutput(mw)
	} else {
		logrus.Info("can not create log file, use default stderr")
	}
}

func main() {
	logrus.Infoln("Starting k8s-webhook")
	server.RunServer()
}
