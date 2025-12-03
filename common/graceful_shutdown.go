package common

import (
	"os"
	"syscall"
	"time"
)

var SignalStopDefault = []os.Signal{
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGTERM,
	syscall.SIGKILL,
}

type GracefulShutdown struct {
	Timeout  time.Duration
	Delay    time.Duration
	HardStop time.Duration
	Signal   []os.Signal
}
