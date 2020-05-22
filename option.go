package grace

import (
	"context"
	"time"
)

// Option define how goroutine works
type Option interface{ wrap(g Goroutine) Goroutine }

// OptionExpire declare goroutine expire
func OptionExpire(expire time.Duration) Option {
	return optionExpire{expire: expire}
}

type optionExpire struct{ expire time.Duration }

func (o optionExpire) wrap(g Goroutine) Goroutine {
	return func(parentCTX context.Context) {
		ctx, cancel := context.WithTimeout(parentCTX, o.expire)
		defer cancel()
		g(ctx)
	}
}

// OptionRestart declare restart option
// -1 means always restart
func OptionRestart(times int) Option { return optionRestart{times} }

type optionRestart struct{ times int }

func (o optionRestart) wrap(g Goroutine) Goroutine {
	times := o.times
	return func(ctx context.Context) {
		for {
			if times == 0 {
				break
			}
			if times > 0 {
				times--
			}
			g(ctx)
		}
	}
}
