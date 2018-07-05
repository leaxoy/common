package logrus

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/leaxoy/logprovider"
)

type (
	KV map[string]interface{}
	LogConf struct {
		LogDir  string
		LogFile string
	}
)

func InitLogger(config LogConf) {
	lp := logprovider.NewAsyncFrame(1, logprovider.NewFileProvider(filepath.Join(config.LogDir, config.LogFile), logprovider.DayDur))
	writer := io.MultiWriter(os.Stderr, lp)
	logrus.SetOutput(writer)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func With(skip int, kv KV) *logrus.Entry {
	if kv == nil {
		kv = make(map[string]interface{})
	}
	caller := "???"
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}
	kv["caller"] = caller
	kv["func"] = runtime.FuncForPC(pc).Name()
	return logrus.WithFields(logrus.Fields(kv))
}

// Infoln wrap `logrus.Infoln` With caller info.
func Infoln(kv KV, args ...interface{}) {
	With(2, kv).Infoln(args...)
}

// Debugln wrap `logrus.Debugln` With caller info.
func Debugln(kv KV, args ...interface{}) {
	With(2, kv).Debugln(args...)
}

// Println wrap `logrus.Println` With caller info.
func Println(kv KV, args ...interface{}) {
	With(2, kv).Println(args...)
}

// Warnln wrap `logrus.Warnln` With caller info.
func Warnln(kv KV, args ...interface{}) {
	With(2, kv).Warnln(args...)
}

// Errorln wrap `logrus.Errorln` With caller info.
func Errorln(err error, kv KV, args ...interface{}) {
	With(2, kv).WithError(err).Errorln(args...)
}

// Fatalln wrap `logrus.Fatalln` With caller info.
func Fatalln(kv KV, args ...interface{}) {
	With(2, kv).Fatalln(args...)
}

// Panicln wrap `logrus.Panicln` With caller info.
func Panicln(kv KV, args ...interface{}) {
	With(2, kv).Panicln(args...)
}
