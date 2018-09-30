package logrus

import (
	"github.com/leaxoy/common/logging"
	"os"
	"testing"
)

func TestHandler_Infoln(t *testing.T) {
	l := NewLogger(os.Stdout, Skip(2))
	l.Infoln(logging.KV{"msg": "1"}, "xx")
}
