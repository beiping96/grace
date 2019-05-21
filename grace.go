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

	defaultLogger = func(format string, a ...interface{}) { fmt.Printf(format, a...) }
)

// Init declare stop signals
// Default is syscall.SIGINT, syscall.SIGQUIT or syscall.SIGTERM
func Init(stopSignals ...os.Signal) {
	if len(stopSignals) == 0 {
		panic("GRACE Init PANIC nil stopSignals")
	}
	defaultStopSignal = stopSignals
}

// Log declare logger method
// Default is fmt.Printf
func Log(logger func(format string, a ...interface{})) {
	defaultLogger = logger
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
	defaultLogger("%s GRACE is running...\n",
		time.Now())
	defaultLogger("%s GRACE stop signal %s \n",
		time.Now(), defaultStopSignal)
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
		defaultLogger("%s GRACE receive stop signal %s \n",
			time.Now(), s)
		cancel()
		defaultLogger("%s GRACE waitting all goroutines exit...\n",
			time.Now())
		select {
		case <-time.After(time.Duration(1) * time.Minute):
		case <-allStopped:
		}
		defaultLogger("%s GRACE stopped.\n",
			time.Now())
	case <-allStopped:
		defaultLogger("%s GRACE all goroutines exit...\n",
			time.Now())
		defaultLogger("%s GRACE stopped.\n",
			time.Now())
	}
	os.Exit(0)
}
