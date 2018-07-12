package limiter

import (
	"testing"
)

func TestNopLimiter_Allow(t *testing.T) {
	l := &nopLimiter{}
	for i := 0; i < 1000; i++ {
		if !l.Allow("") {
			t.Error("limiter should always allow any request")
		}
	}
}

func TestLimiterFunc_Allow(t *testing.T) {
	f := LimiterFunc(func(req string) bool {
		switch req {
		case "allow":
			return true
		}
		return false
	})
	if !f.Allow("allow") {
		t.Error("should allow request `allow`")
	}
	if f.Allow("xxx") {
		t.Error("should not allow request `xxx`")
	}
}

func TestGetLimiter(t *testing.T) {
	l := GetLimiter("xxx")
	if l == nil {
		t.Error("should return default limiter")
	}
}

func TestRegisterLimiter(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			if e.(string) != "err: nil limiter" {
				t.Error("err: panic message")
			}
		}
	}()
	RegisterLimiter("x", nil)
}
