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
)

var (
	defaultStopSignal = []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	}

	defaultLogger = func(format string, a ...interface{}) { fmt.Printf(format, a...) }

	defaultPidDir = "pid/"
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
// Default is ./pid/
func PID(path string) {
	if isRunning {
		panic("GRACE is running, PANIC set pid dir after running.")
	}
	defaultPidDir = path
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
		close(allStopped)
	}()
	defer cancel()

	pid := os.Getpid()
	pidFile := filepath.Join(defaultPidDir, fmt.Sprintf("%d.pid", pid))
	if !filepath.IsAbs(defaultPidDir) {
		dir, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("GRACE pid os.Getwd %v", err))
		}
		pidFile = filepath.Join(dir, defaultPidDir, fmt.Sprintf("%d.pid", pid))
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

	times := 0
	for {
		select {
		case s := <-signalChan:
			times++
			if times > 5 {
				defaultLogger("%s GRACE stopped by killed.\n",
					time.Now())
				return
			}
			defaultLogger("%s GRACE receive stop signal %s \n",
				time.Now(), s)
			go func() {
				cancel()
				defaultLogger("%s GRACE waitting all goroutines exit...\n",
					time.Now())
				select {
				case <-time.After(time.Duration(1) * time.Minute):
				case <-allStopped:
				}
				defaultLogger("%s GRACE stopped.\n",
					time.Now())
				os.Exit(0)
			}()
		case <-allStopped:
			defaultLogger("%s GRACE all goroutines exit...\n",
				time.Now())
			defaultLogger("%s GRACE stopped.\n",
				time.Now())
			return
		}
	}
}
