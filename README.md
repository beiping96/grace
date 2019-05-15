# Grace

[![GoDoc](https://godoc.org/github.com/beiping96/grace?status.svg)](https://godoc.org/github.com/beiping96/grace)
[![Build Status](https://travis-ci.com/beiping96/grace.svg?branch=master)](https://travis-ci.com/beiping96/grace)
[![Go Report Card](https://goreportcard.com/badge/github.com/beiping96/grace)](https://goreportcard.com/report/github.com/beiping96/grace)

A graceful way to manage node and goroutines.

## Usage

``` go
import (
    "context"
    "github.com/beiping96/grace"
)

func init() {
    grace.Init()
}

func main {
    grace.Go(manager0)
    grace.Go(manager1)
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
        grace.Go(func(workLocal interface{}) func(ctx context.Context) {
            return func(ctx context.Context) { do(workLocal) }
        })
    }
}

func do(work interface{}) {}

```
