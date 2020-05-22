# Grace

[![Build Status](https://travis-ci.com/beiping96/grace.svg?branch=master)](https://travis-ci.com/beiping96/grace)
[![GoDoc](https://godoc.org/github.com/beiping96/grace?status.svg)](https://pkg.go.dev/github.com/beiping96/grace)
[![Go Report Card](https://goreportcard.com/badge/github.com/beiping96/grace)](https://goreportcard.com/report/github.com/beiping96/grace)
[![CI On Push](https://github.com/beiping96/grace/workflows/CI-On-Push/badge.svg)](https://github.com/beiping96/grace/actions)

<!-- [![codecov](https://codecov.io/gh/beiping96/grace/branch/master/graph/badge.svg)](https://codecov.io/gh/beiping96/grace) -->

A graceful way to manage node and goroutines. When running, it will listen system signals. After receiving stop signal, the context's cancel function will be called. If all goroutines are exit, node will stop immediately.

## Usage

``` go
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
    grace.Init(syscall.SIGQUIT, syscall.SIGTERM)

    // Log declare logger method
    // Default is fmt.Printf
    grace.Log(log.Printf)

    // Register goroutine
    grace.Go(manager0)
    grace.Go(manager1)

    // Never return
    // Stopped when receiving stop signal
    // or all goroutines are exit
    grace.Run(time.Second)
}

func do(work interface{}) {}

func manager0(ctx context.Context) {
    works := []interface{}{}
    for _, work := range works {
        select {
        case <-ctx.Done():
            // receive stop signal
            return
        default:
        }
        do(work)
    }
}

func manager1(ctx context.Context) {
    works := []interface{}{}
    for _, work := range works {
        select {
        case <-ctx.Done():
            // receive stop signal
            return
        default:
        }
        // start dynamic goroutine
        func(workLocal interface{}) {
            grace.Go(func(ctx context.Context) { do(workLocal) })
        }(work)
    }
}
```
