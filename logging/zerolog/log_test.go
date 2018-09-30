package zerolog

import (
	"github.com/leaxoy/common/logging"
	"os"
	"testing"
)

func TestHandler_Infoln(t *testing.T) {
	l := NewLogger(os.Stdout)
	l.Infoln(logging.KV{"message": 123}, "12345")
}
