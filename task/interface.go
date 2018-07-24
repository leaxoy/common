package task

import (
	"golang.org/x/net/context"
)

type Tasker interface {
	Name() string
	Produce(ctx context.Context, msg interface{}) error
	SetConsumer(fn func(ctx context.Context, msg interface{}))
}

type nopTasker struct{}

func (*nopTasker) Name() string                                              { return "nop" }
func (*nopTasker) Produce(ctx context.Context, msg interface{}) error        { return nil }
func (*nopTasker) SetConsumer(fn func(ctx context.Context, msg interface{})) {}

var taskMap = make(map[string]Tasker)

func RegisterTasker(tasker Tasker) {
	taskMap[tasker.Name()] = tasker
}

func GetTasker(name string) Tasker {
	if t, ok := taskMap[name]; ok {
		return t
	}
	return &nopTasker{}
}
