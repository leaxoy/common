package limiter

type Limiter interface {
	Allow(req string) bool
}

type LimiterFunc func(req string) bool

func (f LimiterFunc) Allow(req string) bool {
	return f(req)
}

type nopLimiter struct{}

func (*nopLimiter) Allow(req string) bool {
	return true
}

var (
	limitMap map[string]Limiter
)

func RegisterLimiter(name string, limiter Limiter) {
	if limiter == nil {
		panic("err: nil limiter")
	}
	limitMap[name] = limiter
}

func GetLimiter(name string) Limiter {
	l, ok := limitMap[name]
	if !ok {
		return &nopLimiter{}
	}
	return l
}
