package grace

import (
	"context"
	"log"
	"syscall"
	"testing"
	"time"
)

func TestGrace(t *testing.T) {
	Signal(syscall.SIGTERM)
	PID(".")
	Log(log.Printf)
	Go(func(ctx context.Context) {
		t.Log("backend goroutine running")
	})
	Go(func(ctx context.Context) {
		Go(func(ctx context.Context) {
			t.Log("dynamic goroutine running")
		})
	})
	Run(time.Second)
}
