# Grace

[![Build Status](https://travis-ci.com/beiping96/grace.svg?branch=master)](https://travis-ci.com/beiping96/grace)
[![GoDoc](https://godoc.org/github.com/beiping96/grace?status.svg)](https://godoc.org/github.com/beiping96/grace)
[![codecov](https://codecov.io/gh/beiping96/grace/branch/master/graph/badge.svg)](https://codecov.io/gh/beiping96/grace)
[![Go Report Card](https://goreportcard.com/badge/github.com/beiping96/grace)](https://goreportcard.com/report/github.com/beiping96/grace)

A graceful way to manage node and goroutines.

## Usage

``` go
import (
    "context"
    "syscall"
    "github.com/beiping96/grace"
)

func main {
    // Declare stop signals
    // Default is syscall.SIGINT, syscall.SIGQUIT or syscall.SIGTERM
    grace.Init(syscall.SIGTERM)
    // Register goroutine
    grace.Go(manager0)
    grace.Go(manager1)
    // Never return
    // Stopped when receiving stop signal
    // or all goroutines are exit
    grace.Run()
}

func manager0(ctx context.Context) {
    works := []interface{}{}
    for _, work := range works {
        select {
            case <- ctx.Done():
                // node receive stop signal
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
            case <- ctx.Done():
                // node receive stop signal
                return
            default:
        }
        // start dynamic goroutine
        grace.Go(func(workLocal interface{}) func(ctx context.Context) {
            return func(ctx context.Context) { do(workLocal) }
        })
    }
}

func do(work interface{}) {}

```
