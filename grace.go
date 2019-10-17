package grace

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	_ "go.uber.org/automaxprocs"
)

var (
	defaultStopSignal = []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	}

	defaultLogger = func(format string, a ...interface{}) { fmt.Printf(format, a...) }

	defaultPidDir *string = nil
)

// Init declare stop signals
// Default is syscall.SIGINT, syscall.SIGQUIT or syscall.SIGTERM
func Init(stopSignals ...os.Signal) {
	if len(stopSignals) == 0 {
		panic("GRACE Init PANIC nil stopSignals")
	}
	if isRunning {
		panic("GRACE is running, PANIC set stop signals after running.")
	}
	defaultStopSignal = stopSignals
}

// Log declare logger method
// Default is fmt.Printf
func Log(logger func(format string, a ...interface{})) {
	if isRunning {
		panic("GRACE is running, PANIC set log after running.")
	}
	defaultLogger = logger
}

// PID configure pid file path
// Default is unable
func PID(path string) {
	if isRunning {
		panic("GRACE is running, PANIC set pid dir after running.")
	}
	defaultPidDir = &path
}

var (
	sysGoroutines = []Goroutine{}
	isRunning     = false
	cancel        func()
	wg            = new(sync.WaitGroup)
)

var (
	CTX context.Context
)

// Goroutine function
type Goroutine func(ctx context.Context)

// Go start a goroutine
func Go(g Goroutine, options ...Option) {
	wrapG := g
	for _, option := range options {
		wrapG = option.wrap(wrapG)
	}

	if !isRunning {
		sysGoroutines = append(sysGoroutines, wrapG)
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g(CTX)
	}()
}

// Run start node
func Run(exitExpire time.Duration) {
	if exitExpire <= 0 {
		exitExpire = time.Minute
	}
	if isRunning {
		panic("GRACE is running, PANIC run twice.")
	}
	defaultLogger("%s GRACE is running...\n",
		time.Now())
	defaultLogger("%s GRACE stop signal %s \n",
		time.Now(), defaultStopSignal)
	isRunning = true
	CTX, cancel = context.WithCancel(context.Background())
	for _, g := range sysGoroutines {
		wg.Add(1)
		go func(goroutine Goroutine) {
			defer wg.Done()
			goroutine(CTX)
		}(g)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, defaultStopSignal...)
	allStopped := make(chan struct{})
	go func() {
		wg.Wait()
		allStopped <- struct{}{}
		close(allStopped)
	}()
	defer cancel()

	if defaultPidDir != nil && len(*defaultPidDir) != 0 {
		pidDir := *defaultPidDir
		pid := os.Getpid()
		pidFile := filepath.Join(pidDir, fmt.Sprintf("%d.pid", pid))
		if !filepath.IsAbs(pidDir) {
			dir, err := os.Getwd()
			if err != nil {
				panic(fmt.Errorf("GRACE pid os.Getwd %v", err))
			}
			pidFile = filepath.Join(dir, pidDir, fmt.Sprintf("%d.pid", pid))
		}
		os.MkdirAll(filepath.Dir(pidFile), 0777)
		err := ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0777)
		if err != nil {
			panic(fmt.Errorf("GRACE pid monitor %s %v", pidFile, err))
		}
		defer func() {
			err := os.Remove(pidFile)
			if err != nil {
				panic(fmt.Errorf("GRACE pid exit %s %v", pidFile, err))
			}
		}()
	}

	select {
	case s := <-signalChan:
		defaultLogger("%s GRACE receive stop signal %s \n",
			time.Now(), s)
		cancel()
		defaultLogger("%s GRACE waitting all goroutines exit...\n",
			time.Now())
		select {
		case <-time.After(exitExpire):
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
}
