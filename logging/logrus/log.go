package logrus

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/leaxoy/common/logging"
	"github.com/leaxoy/logprovider"
	"github.com/sirupsen/logrus"
)

type handler struct{}

func (h *handler) Debugln(kv logging.KV, msg string) {
	With(2, kv).Debugln(msg)
}

func (h *handler) Infoln(kv logging.KV, msg string) {
	With(2, kv).Infoln(msg)
}

func (h *handler) Warnln(kv logging.KV, msg string) {
	With(2, kv).Warnln(msg)
}

func (h *handler) Errorln(err error, kv logging.KV, msg string) {
	With(2, kv).WithError(err).Warnln(msg)
}

func (h *handler) Panicln(kv logging.KV, msg string) {
	With(2, kv).Panicln(msg)
}

func (h *handler) Fatalln(kv logging.KV, msg string) {
	With(2, kv).Fatalln(msg)
}

func NewLogger(config logging.LoggerConfig) logging.Logger {
	lp := logprovider.NewAsyncFrame(1, logprovider.NewFileProvider(filepath.Join(config.LogDir, config.LogFile), logprovider.DayDur))
	writer := io.MultiWriter(os.Stderr, lp)
	logrus.SetOutput(writer)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	return &handler{}
}

func With(skip int, kv logging.KV) *logrus.Entry {
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
