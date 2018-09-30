package zerolog

import (
	"github.com/leaxoy/common/logging"
	"github.com/rs/zerolog"
	"io"
)

type handler struct {
	logger zerolog.Logger
	skip   int
}

type Option func(*handler)

func Skip(skip int) Option {
	return func(i *handler) {
		i.skip = skip
	}
}

func (h *handler) Debugln(kv logging.KV, msg string) {
	h.logger.Debug().Fields(kv).Msg(msg)
}

func (h *handler) Infoln(kv logging.KV, msg string) {
	h.logger.Info().Caller().Timestamp().Fields(kv).Msg(msg)
}

func (h *handler) Warnln(kv logging.KV, msg string) {
	h.logger.Warn().Fields(kv).Msg(msg)
}

func (h *handler) Errorln(err error, kv logging.KV, msg string) {
	h.logger.Error().Err(err).Fields(kv).Msg(msg)
}

func (h *handler) Fatalln(kv logging.KV, msg string) {
	h.logger.Fatal().Caller().Timestamp().Fields(kv).Msg(msg)
}

func NewLogger(w io.Writer, opts ...Option) logging.Logger {
	h := &handler{skip: 3}
	for _, opt := range opts {
		opt(h)
	}
	zerolog.CallerSkipFrameCount = h.skip
	zerolog.MessageFieldName = "msg"
	l := zerolog.New(w)
	h.logger = l
	return h
}
