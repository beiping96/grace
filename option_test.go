package grace

import (
	"context"
	"syscall"
	"testing"
	"time"
)

func TestOptionExpire(t *testing.T) {
	Signal(syscall.SIGTERM)
	Go(func(ctx context.Context) {
		t.Log("backend goroutine running")
		<-ctx.Done()
		t.Log("backend goroutine stopped")
	}, OptionExpire(time.Second))
	Go(func(ctx context.Context) {
		Go(func(ctx context.Context) {
			t.Log("dynamic goroutine running")
			<-ctx.Done()
			t.Log("dynamic goroutine stopped")
		}, OptionExpire(time.Millisecond))
	})
	Run(time.Second)
}

func TestOptionRestart(t *testing.T) {
	Signal(syscall.SIGTERM)
	Go(func(ctx context.Context) {
		t.Log("backend goroutine running")
	}, OptionRestart(1))
	Go(func(ctx context.Context) {
		Go(func(ctx context.Context) {
			t.Log("dynamic goroutine running")
		}, OptionRestart(2))
	})
	Run(time.Second)
}
