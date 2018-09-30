package logging_test

import (
	"github.com/leaxoy/common/logging"
	"github.com/leaxoy/common/logging/logrus"
	"github.com/leaxoy/common/logging/zap"
	"github.com/leaxoy/common/logging/zerolog"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestZapInfoln(t *testing.T) {
	l := zap.NewLogger(os.Stdout, zap.Skip(2))
	l.Infoln(logging.KV{"message": 2, "latency": time.Second + time.Millisecond*123}, "zap")
}

func TestLogrusInfoln(t *testing.T) {
	l := logrus.NewLogger(os.Stdout, logrus.Skip(3))
	l.Infoln(logging.KV{"message": 2, "latency": time.Second + time.Millisecond*123}, "logrus")
}

func TestZeroInfoln(t *testing.T) {
	l := zerolog.NewLogger(os.Stdout, zerolog.Skip(3))
	l.Infoln(logging.KV{"message": 2, "latency": time.Second + time.Millisecond*123}, "logrus")
}

func BenchmarkZapInfoln(b *testing.B) {
	l := zap.NewLogger(ioutil.Discard, zap.Skip(2))
	n := time.Now()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Infoln(logging.KV{
			"string":  "Hello",
			"latency": time.Second + time.Millisecond*123,
			"ints":    2,
			"float":   1.2,
			"time":    n,
			"bool":    false,
			//"map":     logging.KV{},
			"byte": 'b',
		}, "zap")
	}
}

func BenchmarkLogrusInfoln(b *testing.B) {
	l := logrus.NewLogger(ioutil.Discard, logrus.Skip(3))
	n := time.Now()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Infoln(logging.KV{
			"string":  "Hello",
			"latency": time.Second + time.Millisecond*123,
			"ints":    2,
			"float":   1.2,
			"time":    n,
			"bool":    false,
			//"map":     logging.KV{},
			"byte": 'b',
		}, "logrus")
	}
}

func BenchmarkZeroInfoln(b *testing.B) {
	l := zerolog.NewLogger(ioutil.Discard, zerolog.Skip(4))
	n := time.Now()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Infoln(logging.KV{
			"string":  "Hello",
			"latency": time.Second + time.Millisecond*123,
			"ints":    2,
			"float":   1.2,
			"time":    n,
			"bool":    false,
			//"map":     logging.KV{},
			"byte": 'b',
		}, "logrus")
	}
}
