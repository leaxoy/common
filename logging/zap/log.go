package zap

import (
	"github.com/leaxoy/common/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

type handler struct {
	syncTicker time.Duration
	logger     *zap.Logger
	skip       int
}

type Option func(*handler)

func Skip(skip int) Option {
	return func(i *handler) {
		i.skip = skip
	}
}

func SyncTicker(t time.Duration) Option {
	return func(i *handler) {
		i.syncTicker = t
	}
}

func (h *handler) sync(core zapcore.Core) {
	ticker := time.NewTicker(h.syncTicker)
	for range ticker.C {
		core.Sync()
	}
}

func (*handler) Debugln(kv logging.KV, msg string) {
	zap.L().Debug(msg, wrapFields(nil, kv)...)
}

func (*handler) Infoln(kv logging.KV, msg string) {
	zap.L().Info(msg, wrapFields(nil, kv)...)
}

func (*handler) Warnln(kv logging.KV, msg string) {
	zap.L().Warn(msg, wrapFields(nil, kv)...)
}

func (*handler) Errorln(err error, kv logging.KV, msg string) {
	zap.L().Error(msg, wrapFields(err, kv)...)
}

func (*handler) Panicln(kv logging.KV, msg string) {
	zap.L().Panic(msg, wrapFields(nil, kv)...)
}

func (*handler) Fatalln(kv logging.KV, msg string) {
	zap.L().Fatal(msg, wrapFields(nil, kv)...)
}

func NewLogger(w io.Writer, opts ...Option) logging.Logger {
	h := &handler{syncTicker: time.Duration(10) * time.Second, skip: 2}
	for _, opt := range opts {
		opt(h)
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		CallerKey:      "caller",
		EncodeCaller:   zapcore.FullCallerEncoder,
		TimeKey:        "time",
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LevelKey:       "level",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		MessageKey:     "msg",
	}), zapcore.AddSync(w), zap.InfoLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(h.skip))
	zap.ReplaceGlobals(logger)
	go h.sync(core)
	return h
}

func wrapFields(err error, kv logging.KV) []zap.Field {
	var fields []zap.Field
	if err != nil {
		fields = make([]zap.Field, 0, len(kv)+1)
		fields[0] = zap.Error(err)
	} else {
		fields = make([]zap.Field, 0, len(kv))
	}
	for key, val := range kv {
		fields = append(fields, zap.Any(key, val))
	}
	return fields
}
