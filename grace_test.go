package grace

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		stopSignals []os.Signal
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.stopSignals...)
		})
	}
}
