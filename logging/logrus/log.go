package logrus

import (
	"fmt"
	"io"
	"runtime"

	"github.com/leaxoy/common/logging"
	"github.com/sirupsen/logrus"
)

type (
	handler struct {
		skip int
		kv   map[string]interface{}
	}
	Option func(*handler)
)

func Skip(skip int) Option {
	return func(h *handler) {
		h.skip = skip
	}
}

func (h *handler) Debugln(kv logging.KV, msg string) {
	entry(h.skip, kv).Debugln(msg)
}

func (h *handler) Infoln(kv logging.KV, msg string) {
	entry(h.skip, kv).Infoln(msg)
}

func (h *handler) Warnln(kv logging.KV, msg string) {
	entry(h.skip, kv).Warnln(msg)
}

func (h *handler) Errorln(err error, kv logging.KV, msg string) {
	entry(h.skip, kv).WithError(err).Warnln(msg)
}

func (h *handler) Panicln(kv logging.KV, msg string) {
	entry(h.skip, kv).Panicln(msg)
}

func (h *handler) Fatalln(kv logging.KV, msg string) {
	entry(h.skip, kv).Fatalln(msg)
}

func NewLogger(w io.Writer, opts ...Option) logging.Logger {
	h := &handler{}
	for _, opt := range opts {
		opt(h)
	}
	logrus.SetOutput(w)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	return h
}

func entry(skip int, kv logging.KV) *logrus.Entry {
	if kv == nil {
		kv = make(map[string]interface{})
	}
	caller := "???"
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}
	kv["caller"] = caller
	return logrus.WithFields(logrus.Fields(kv))
}
