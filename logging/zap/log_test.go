package zap

import (
	"github.com/leaxoy/common/logging"
	"os"
	"testing"
)

func TestHandler_Infoln(t *testing.T) {
	h := NewLogger(os.Stdout, Skip(1))
	h.Infoln(logging.KV{"x": 1}, "x")
}
