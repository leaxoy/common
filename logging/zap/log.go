package zap

import (
	"github.com/leaxoy/common/logging"
	"github.com/leaxoy/logprovider"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"time"
)

type handler struct {
	syncTicker time.Duration
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

func InitLogger(config logging.LoggerConfig) logging.Logger {
	lp := logprovider.NewAsyncFrame(1, logprovider.NewFileProvider(filepath.Join(config.LogDir, config.LogFile), logprovider.DayDur))
	writer := io.MultiWriter(os.Stderr, lp)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		EncodeDuration: func(duration time.Duration, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(duration.String())
		},
		EncodeTime: zapcore.ISO8601TimeEncoder,
	}), zapcore.AddSync(writer), zap.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	h := &handler{syncTicker: time.Duration(10) * time.Second}
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
