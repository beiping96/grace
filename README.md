# Grace

[![Build Status](https://travis-ci.com/beiping96/grace.svg?branch=master)](https://travis-ci.com/beiping96/grace)
[![GoDoc](https://godoc.org/github.com/beiping96/grace?status.svg)](https://pkg.go.dev/github.com/beiping96/grace)
[![Go Report Card](https://goreportcard.com/badge/github.com/beiping96/grace)](https://goreportcard.com/report/github.com/beiping96/grace)
[![CI On Push](https://github.com/beiping96/grace/workflows/CI-On-Push/badge.svg)](https://github.com/beiping96/grace/actions)

<!-- [![codecov](https://codecov.io/gh/beiping96/grace/branch/master/graph/badge.svg)](https://codecov.io/gh/beiping96/grace) -->

Grace manages long-running goroutines gracefully by trapping system signals and canceling `context.Context`.

`import "github.com/beiping96/grace"`

The core of Grace is `type Goroutine func(ctx context.Context)`.

## Usage

### Simple
```go
package main

import(
    "context"
    "time"
    "github.com/beiping96/grace"
)

func main() {
    // Register worker
    grace.Go(worker)

    // Start worker and hold main goroutine.
    // After process received system signal, 
    // worker should exit in two seconds,
    // otherwise, be killed.
    grace.Run(time.Second * time.Duration(2)) 
}

func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            // receive stop signal
            return
        default:
        }
        // do some job costs second
        time.Sleep(time.Second)      
    }
}
```

### Pub/Sub
```go
package main

import(
    "context"
    "time"
    "github.com/beiping96/grace"
)

func main() {
    // Register puber
    grace.Go(puber)
    // Register 2 suber
    for i := 0; i < 2; i++ { 
        grace.Go(suber) 
    }

    grace.Run(time.Second) // start
}

var (
    jobChan = make(chan struct{}, 1)
)

func puber(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            // receive stop signal
            close(jobChan)
            return
        default:
        }
        // load job
        var (
            job struct{}
        )
        // pub job
        jobChan <- job 
    }
}

func suber(ctx context.Context) {
    for job := range jobChan {
        // handle job
    	println(job)
    }
}
```

### HTTPD
```go
package main

import (
    "context"
    "net/http"
    "time"
    "github.com/beiping96/grace"
)

func httpd(ctx context.Context) {
    server := &http.Server{}
    
    grace.Go(func(ctx context.Context) {
        <-ctx.Done()
        server.Shutdown(context.Background())
    })
    
    server.ListenAndServe()
}

func main() {
    // Register httpd
    grace.Go(httpd)
    
    // Start
    // After process received system signal, 
    // wait all http connection closed is ten seconds
    grace.Run(time.Second * time.Duration(10)) 
}
```

### Option
```go
package main

import (
    "context"
    "time"
    "github.com/beiping96/grace"
)

func main() {
    // After one minute, the worker's ctx will be canceled
    grace.Go(worker, grace.OptionExpire(time.Minute))    
    
    // If worker stopped, it will be restart
    // The max number of restart is 2
    grace.Go(worker, grace.OptionRestart(2))
    
    // If worker stopped, it will be restart
    // The max number of restart is 2
    // Each worker can executed one minute
    grace.Go(worker, grace.OptionRestart(2),
                     grace.OptionExpire(time.Minute))
    
    // If worker stopped, it will be restart
    // The max number of restart is 2
    // All workers only have one minute to executed
    grace.Go(worker, grace.OptionExpire(time.Minute),
                     grace.OptionRestart(2))

    grace.Run(time.Second)
}


func worker(ctx context.Context) {}
```

### Custom
``` go
package main

import (
    "log"
    "context"
    "syscall"
    "time"
    "github.com/beiping96/grace"
)

func main() {
    // Declare stop signals
    // Default is syscall.SIGINT, syscall.SIGQUIT or syscall.SIGTERM
    grace.Signal(syscall.SIGQUIT, syscall.SIGTERM)

    // Declare logger method
    // Default is os.Stdout
    grace.Log(log.Printf)

    // Declare process id folder
    // Default is unable
    grace.PID("./")
    
    // Register goroutine
    grace.Go(worker)
    
    // Register dynamic goroutine
    grace.Go(func(ctx context.Context) { grace.Go(worker) })

    // Never return
    // Stopped when receiving stop signal
    // or all goroutines are exit
    grace.Run(time.Second)
}

func worker(ctx context.Context) {}
```
