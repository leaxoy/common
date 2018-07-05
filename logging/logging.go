package logging

import (
	"sync"
)

type (
	KV           map[string]interface{}
	LoggerConfig struct {
		LogDir  string `json:"log_dir" yaml:"log_dir" toml:"log_dir"`
		LogFile string `json:"log_file" yaml:"log_file" toml:"log_file"`
	}
	Logger interface {
		Debugln(kv KV, msg string)
		Infoln(kv KV, msg string)
		Warnln(kv KV, msg string)
		Errorln(err error, kv KV, msg string)
		Panicln(kv KV, msg string)
		Fatalln(kv KV, msg string)
	}
)

type nopLogger struct{}

func (*nopLogger) Debugln(kv KV, msg string) {}

func (*nopLogger) Infoln(kv KV, msg string) {}

func (*nopLogger) Warnln(kv KV, msg string) {}

func (*nopLogger) Errorln(err error, kv KV, msg string) {}

func (*nopLogger) Panicln(kv KV, msg string) {}

func (*nopLogger) Fatalln(kv KV, msg string) {}

var (
	defaultLogger Logger = &nopLogger{}
	once          sync.Once
)

func SetLogger(logger Logger) {
	once.Do(func() {
		defaultLogger = logger
	})
}

func Debugln(kv KV, args ...string) {
	if len(args) == 0 {
		defaultLogger.Debugln(kv, "")
	}
	defaultLogger.Debugln(kv, args[0])
}

func Infoln(kv KV, args ...string) {
	if len(args) == 0 {
		defaultLogger.Infoln(kv, "")
	}
	defaultLogger.Infoln(kv, args[0])
}

func Warnln(kv KV, args ...string) {
	if len(args) == 0 {
		defaultLogger.Warnln(kv, "")
	}
	defaultLogger.Warnln(kv, args[0])
}

func Errorln(err error, kv KV, args ...string) {
	if len(args) == 0 {
		defaultLogger.Errorln(err, kv, "")
	}
	defaultLogger.Errorln(err, kv, args[0])
}

func Panicln(kv KV, args ...string) {
	if len(args) == 0 {
		defaultLogger.Panicln(kv, "")
	}
	defaultLogger.Panicln(kv, args[0])
}
