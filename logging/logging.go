package logging

import (
	"fmt"
	"sync"
)

type (
	KV map[string]interface{}
	LoggerConfig struct {
		LogDir  string `json:"log_dir" yaml:"log_dir" toml:"log_dir"`
		LogFile string `json:"log_file" yaml:"log_file" toml:"log_file"`
	}
	Logger interface {
		Debugln(kv KV, msg string)
		Infoln(kv KV, msg string)
		Warnln(kv KV, msg string)
		Errorln(err error, kv KV, msg string)
		Fatalln(kv KV, msg string)

		Debugf(kv KV, msg string, args ...interface{})
		Infof(kv KV, msg string, args ...interface{})
		Warnf(kv KV, msg string, args ...interface{})
		Errorf(err error, kv KV, msg string, args ...interface{})
		Fatalf(kv KV, msg string, args ...interface{})

		Debug(kv KV, msg string)
		Info(kv KV, msg string)
		Warn(kv KV, msg string)
		Error(err error, kv KV, msg string)
		Fatal(kv KV, msg string)
	}
)

type nopLogger struct{}

func (*nopLogger) Debug(kv KV, msg string)            {}
func (*nopLogger) Info(kv KV, msg string)             {}
func (*nopLogger) Warn(kv KV, msg string)             {}
func (*nopLogger) Error(err error, kv KV, msg string) {}
func (*nopLogger) Fatal(kv KV, msg string)            {}

func (*nopLogger) Debugf(kv KV, msg string, args ...interface{})            {}
func (*nopLogger) Infof(kv KV, msg string, args ...interface{})             {}
func (*nopLogger) Warnf(kv KV, msg string, args ...interface{})             {}
func (*nopLogger) Errorf(err error, kv KV, msg string, args ...interface{}) {}
func (*nopLogger) Fatalf(kv KV, msg string, args ...interface{})            {}

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

func Debug(kv KV, msg string)            { defaultLogger.Debug(kv, msg) }
func Info(kv KV, msg string)             { defaultLogger.Info(kv, msg) }
func Warn(kv KV, msg string)             { defaultLogger.Warn(kv, msg) }
func Error(err error, kv KV, msg string) { defaultLogger.Error(err, kv, msg) }
func Fatal(kv KV, msg string)            { defaultLogger.Fatal(kv, msg) }

func Debugf(kv KV, msg string, args ...interface{})            { Debug(kv, fmt.Sprintf(msg, args...)) }
func Infof(kv KV, msg string, args ...interface{})             { Info(kv, fmt.Sprintf(msg, args...)) }
func Warnf(kv KV, msg string, args ...interface{})             { Warn(kv, fmt.Sprintf(msg, args...)) }
func Errorf(err error, kv KV, msg string, args ...interface{}) { Error(err, kv, fmt.Sprintf(msg, args...)) }

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

func Fatalln(kv KV, args ...string) {
	if len(args) == 0 {
		defaultLogger.Fatalln(kv, "")
	}
	defaultLogger.Fatalln(kv, args[0])
}
