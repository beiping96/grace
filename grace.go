package grace

import (
	"context"
	"os"
)

type Config struct {
	StopSignal     []os.Signal
	RestartSignal  []os.Signal
	ShutdownSignal []os.Signal
}

func Init(cfg *Config) {}

func Go(name string, fn func(ctx context.Context)) {}

func Run() {}
