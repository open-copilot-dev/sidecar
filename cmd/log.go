package main

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
	"time"
)

func initLog(debug *bool) {
	// 日志目录
	logDir := "./logs/"
	if err := os.MkdirAll(logDir, 0o777); err != nil {
		log.Println(err.Error())
		return
	}

	// 日志文件
	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logDir, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return
		}
	}

	// 日志对象
	logger := hertzlogrus.NewLogger(hertzlogrus.WithHook(&LogIdHook{}))
	logger.SetOutput(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    20, // 一个文件最大可达 20M。
		MaxAge:     10, // 一个文件最多可以保存 10 天。
		MaxBackups: 5,  // 最多同时保存 5 个文件。
		LocalTime:  true,
		Compress:   true, // 用 gzip 压缩。
	})
	hlog.SetLogger(logger)

	// 日志级别
	if debug != nil && *debug {
		hlog.SetLevel(hlog.LevelDebug)
	} else {
		hlog.SetLevel(hlog.LevelInfo)
	}
}

type LogIdHook struct{}

func (h *LogIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LogIdHook) Fire(e *logrus.Entry) error {
	ctx := e.Context
	if ctx == nil {
		return nil
	}
	value := ctx.Value("X-Request-ID")
	if value != nil {
		e.Data["log_id"] = value
	}
	return nil
}
