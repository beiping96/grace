package grace

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	defaultStopSignal = []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	}
)

// Init declare stop signals
// Default is syscall.SIGINT, syscall.SIGQUIT or syscall.SIGTERM
func Init(stopSignals ...os.Signal) {
	if len(stopSignals) == 0 {
		panic("GRACE Init PANIC nil stopSignals")
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

// Goroutine function
type Goroutine func(ctx context.Context)

// Go start a goroutine
func Go(g Goroutine) {
	if isRunning {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g(ctx)
		}()
		return
	}
	mu.Lock()
	sysGoroutines = append(sysGoroutines, g)
	mu.Unlock()
}

// Run start node
func Run() {
	mu.Lock()
	defer mu.Unlock()
	isRunning = true
	ctx, cancel = context.WithCancel(context.Background())
	for _, g := range sysGoroutines {
		wg.Add(1)
		go func(goroutine Goroutine) {
			defer wg.Done()
			goroutine(ctx)
		}(g)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, defaultStopSignal...)
	allStopped := make(chan struct{})
	go func() {
		wg.Wait()
		allStopped <- struct{}{}
	}()
	defer cancel()
	select {
	case s := <-signalChan:
		fmt.Printf("GRACE receive stop signal %s \n", s)
		cancel()
		fmt.Printf("GRACE waitting all goroutines exit...\n")
		select {
		case <-time.After(time.Duration(5) * time.Second):
		case <-allStopped:
		}
		fmt.Printf("GRACE stopped.\n")
	case <-allStopped:
		fmt.Printf("GRACE all goroutines exit...\n")
		fmt.Printf("GRACE stopped.\n")
	}
	os.Exit(0)
}
