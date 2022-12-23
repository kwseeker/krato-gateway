package middleware

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestKratosLog(t *testing.T) {
	//1 通过log函数输出日志，默认是输出到标准错误 os.Stderr
	log.Info("log.Info()")
	//log.Fatalf("log.Fatalf()") //和其他级别不同的是日志输出完后进程就退出了

	//2 定制一个Logger打印日志
	//如果指定使用DefaultLogger, 会改变log.Xxx()输出
	logger := log.With(log.DefaultLogger, "preKey", "preValue", "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)
	l := log.NewHelper(log.NewFilter(log.NewFilter(logger, log.FilterLevel(log.LevelWarn), log.FilterKey("passwd")))) //第一层套了个空的过滤器
	l.Log(log.LevelInfo, "foo", "bar")
	l.Debug("debug msg ...")
	l.Info("info msg ...")
	l.Warn("warn msg ...")
	l.Warnw("key3", "val3", "passwd", "123456") //w结尾的方法，用于打印键值对日志
	//l.Fatal()

	//如果不想改变原logger,但是又想用原logger配置好的功能
	w := log.NewWriter(logger, log.WithWriterLevel(log.LevelError), log.WithWriteMessageKey("nlMsg"))
	_, _ = w.Write([]byte("Hello new writer"))
}
