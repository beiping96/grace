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
	wg            = new(sync.WaitGroup)
)

// Goroutine function
type Goroutine func(ctx context.Context)

// Go start a goroutine
func Go(g Goroutine) {
	if !isRunning {
		sysGoroutines = append(sysGoroutines, g)
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g(ctx)
	}()
}

// Run start node
func Run() {
	if isRunning {
		panic("GRACE is running, PANIC run twice.")
	}
	fmt.Printf("GRACE is running...\n")
	fmt.Printf("GRACE stop signal %s \n", defaultStopSignal)
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
		case <-time.After(time.Duration(1) * time.Minute):
		case <-allStopped:
		}
		fmt.Printf("GRACE stopped.\n")
	case <-allStopped:
		fmt.Printf("GRACE all goroutines exit...\n")
		fmt.Printf("GRACE stopped.\n")
	}
	os.Exit(0)
}
