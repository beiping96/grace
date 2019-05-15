package grace

import (
	"context"
	"os"
	"sync"
	"syscall"
)

var (
	defaultStopSignal = []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
	}
)

// Init Declare stop signals
// Default is syscall.SIGINT or syscall.SIGQUIT
func Init(stopSignals ...os.Signal) {
	if len(stopSignals) == 0 {
		panic("beiping96/grace Init PANIC nil stopSignals")
	}
	defaultStopSignal = stopSignals
}

var (
	sysGoroutines = []Goroutine{}
	isRunning     = false
	cancel        func()
	ctx           context.Context
	mu            = new(sync.Mutex)
	wg            = new(sync.WaitGroup)
)

// Goroutine Goroutine function
type Goroutine func(ctx context.Context)

// Go Start a goroutine
func Go(g Goroutine) {
	if isRunning {
		go func() {
			wg.Add(1)
			defer wg.Done()
			g(ctx)
		}()
		return
	}
	mu.Lock()
	sysGoroutines = append(sysGoroutines, g)
	mu.Unlock()
}

// Run Start node
func Run() {
	mu.Lock()
	defer mu.Unlock()
	isRunning = true
	ctx, cancel = context.WithCancel(context.Background())
	go func() {

	}()
	wg.Wait()
}
