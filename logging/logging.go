package logging

import (
	"sync"
)

type (
	KV map[string]interface{}
	Logger interface {
		Debugln(kv KV, msg string)
		Infoln(kv KV, msg string)
		Warnln(kv KV, msg string)
		Errorln(err error, kv KV, msg string)
		Fatalln(kv KV, msg string)
	}
)

type nopLogger struct{}

func (*nopLogger) Debugln(kv KV, msg string)            {}
func (*nopLogger) Infoln(kv KV, msg string)             {}
func (*nopLogger) Warnln(kv KV, msg string)             {}
func (*nopLogger) Errorln(err error, kv KV, msg string) {}
func (*nopLogger) Fatalln(kv KV, msg string)            {}

var (
	defaultLogger Logger = &nopLogger{}
	once          sync.Once
)

func SetLogger(logger Logger) {
	once.Do(func() {
		defaultLogger = logger
	})
}

func Debugln(kv KV, msg string)            { defaultLogger.Debugln(kv, msg) }
func Infoln(kv KV, msg string)             { defaultLogger.Infoln(kv, msg) }
func Warnln(kv KV, msg string)             { defaultLogger.Warnln(kv, msg) }
func Errorln(err error, kv KV, msg string) { defaultLogger.Errorln(err, kv, msg) }
func Fatalln(kv KV, msg string)            { defaultLogger.Fatalln(kv, msg) }
