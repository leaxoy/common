package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/leaxoy/common/task"
	"golang.org/x/net/context"
	"golang.org/x/sync/semaphore"
)

type (
	HandlerConfig struct {
		ProducerConfig *sarama.Config
		ConsumerConfig *cluster.Config

		Addrs   []string
		GroupID string
		Topics  []string

		MaxQueueSize int

		MaxConcurrent int64
	}
	Option func(*Handler)

	Handler struct {
		onProduce              func(producer sarama.AsyncProducer, msg *sarama.ProducerMessage) error
		onConsume              func(consumer *cluster.Consumer, msg *sarama.ConsumerMessage) error
		onProducerError        func(producer sarama.AsyncProducer, err *sarama.ProducerError) error
		onProducerSuccess      func(producer sarama.AsyncProducer, msg *sarama.ProducerMessage) error
		onConsumerError        func(consumer *cluster.Consumer, err error) error
		onConsumerNotification func(consumer *cluster.Consumer, n *cluster.Notification) error

		handleError func(err error)

		producerC   chan *sarama.ProducerMessage
		concurrency *semaphore.Weighted

		ctx context.Context

		producer sarama.AsyncProducer
		consumer *cluster.Consumer
	}
)

const (
	Name = "kafka"
)

func WithProduce(fn func(producer sarama.AsyncProducer, msg *sarama.ProducerMessage) error) Option {
	return func(handler *Handler) {
		handler.onProduce = fn
	}
}

func WithConsume(fn func(consumer *cluster.Consumer, msg *sarama.ConsumerMessage) error) Option {
	return func(handler *Handler) {
		handler.onConsume = fn
	}
}

func WithProducerError(fn func(producer sarama.AsyncProducer, err *sarama.ProducerError) error) Option {
	return func(handler *Handler) {
		handler.onProducerError = fn
	}
}

func WithProducerSuccess(fn func(producer sarama.AsyncProducer, msg *sarama.ProducerMessage) error) Option {
	return func(handler *Handler) {
		handler.onProducerSuccess = fn
	}
}

func WithConsumerError(fn func(consumer *cluster.Consumer, err error) error) Option {
	return func(handler *Handler) {
		handler.onConsumerError = fn
	}
}

func WithConsumerNotification(fn func(consumer *cluster.Consumer, n *cluster.Notification) error) Option {
	return func(handler *Handler) {
		handler.onConsumerNotification = fn
	}
}

func WithHandleError(fn func(err error)) Option {
	return func(handler *Handler) {
		handler.handleError = fn
	}
}

func InitTaskHandler(ctx context.Context, config *HandlerConfig, opts ...Option) *Handler {
	if err := config.check(); err != nil {
		panic(err)
	}
	h := &Handler{
		producerC:   make(chan *sarama.ProducerMessage, config.MaxQueueSize),
		ctx:         ctx,
		concurrency: semaphore.NewWeighted(config.MaxConcurrent),
	}
	if config.ProducerConfig != nil {
		producer, err := sarama.NewAsyncProducer(config.Addrs, config.ProducerConfig)
		if err != nil {
			panic(err)
		}
		h.producer = producer
	}
	if config.ConsumerConfig != nil {
		consumer, err := cluster.NewConsumer(config.Addrs, config.GroupID, config.Topics, config.ConsumerConfig)
		if err != nil {
			panic(err)
		}
		h.consumer = consumer
	}
	for _, opt := range opts {
		opt(h)
	}
	if h.handleError == nil {
		h.handleError = func(err error) {}
	}
	h.setDefault()
	if err := h.check(); err != nil {
		panic(err)
	}
	task.RegisterTasker(h)
	go h.runProducer()
	go h.runConsumer()
	return h
}

func (h *Handler) Name() string {
	return Name
}

func (h *Handler) Produce(ctx context.Context, msg interface{}) error {
	if h.producer == nil {
		return errors.New("err: handler not support produce message")
	}
	m, ok := msg.(*sarama.ProducerMessage)
	if !ok {
		return fmt.Errorf("err: not a valid sarama.ProducerMessage(%+v)", msg)
	}
	select {
	case h.producerC <- m:
	default:
		return errors.New("err: produce message timeout")
	}
	return nil
}

func (h *Handler) SetConsumer(fn func(ctx context.Context, msg interface{})) {}

func (h *Handler) check() error {
	if h.consumer == nil && h.producer == nil {
		return errors.New("err: both consumer and producer are nil")
	}
	if h.consumer != nil && h.onConsume == nil {
		return errors.New("err: consumer is not nil, but onConsume is nil")
	}
	if h.producer != nil && h.onProduce == nil {
		return errors.New("err: producer is not nil, but onProduce is nil")
	}
	return nil
}

func (h *Handler) setDefault() {
	if h.producer != nil {
		if h.onProducerSuccess == nil {
			h.onProducerSuccess = func(producer sarama.AsyncProducer, msg *sarama.ProducerMessage) error {
				return nil
			}
		}
		if h.onProducerError == nil {
			h.onProducerError = func(producer sarama.AsyncProducer, err *sarama.ProducerError) error {
				return nil
			}
		}
	}
	if h.consumer != nil {
		if h.onConsumerError == nil {
			h.onConsumerError = func(consumer *cluster.Consumer, err error) error {
				return nil
			}
		}
		if h.onConsumerNotification == nil {
			h.onConsumerNotification = func(consumer *cluster.Consumer, n *cluster.Notification) error {
				return nil
			}
		}
	}
}

func (h *Handler) runProducer() {
	if h.producer == nil || h.producerC == nil {
		return
	}
	for {
		select {
		case <-h.ctx.Done():
			return
		case err := <-h.producer.Errors():
			h.handleError(h.onProducerError(h.producer, err))
		case msg := <-h.producerC:
			h.handleError(h.onProduce(h.producer, msg))
		case success := <-h.producer.Successes():
			h.handleError(h.onProducerSuccess(h.producer, success))
		}
	}
}

func (h *Handler) runConsumer() {
	if h.consumer == nil {
		return
	}
	select {
	case <-h.ctx.Done():
		return
	case err := <-h.consumer.Errors():
		h.handleError(h.onConsumerError(h.consumer, err))
	case notification := <-h.consumer.Notifications():
		h.handleError(h.onConsumerNotification(h.consumer, notification))
	case msg := <-h.consumer.Messages():
		if err := h.concurrency.Acquire(context.Background(), 1); err == nil {
			go func() {
				h.handleError(h.onConsume(h.consumer, msg))
				h.concurrency.Release(1)
			}()
		}
	}
}

func (c *HandlerConfig) check() error {
	if c.MaxConcurrent <= 0 {
		return errors.New("[task] err: max concurrent must large than 0")
	}
	if c.MaxQueueSize < (1 << 8) {
		c.MaxQueueSize = 1 << 8
	}
	if c.MaxQueueSize > (1 << 12) {
		c.MaxQueueSize = 1 << 12
	}
	if len(c.Addrs) == 0 {
		return errors.New("[task] err: addrs should has at least one addr")
	}
	if c.ConsumerConfig == nil && c.ProducerConfig == nil {
		return errors.New("[task] err: producer and consumer config both nil")
	}
	if c.ConsumerConfig != nil {
		if len(c.Topics) == 0 {
			return errors.New("[task] err: consumer config is not nil, but no topics given")
		}
		if c.GroupID == "" {
			return errors.New("[task] err: consumer config it not nil, but groupID is nil")
		}
	}
	return nil
}
