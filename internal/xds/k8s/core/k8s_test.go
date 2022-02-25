package core

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient("", "test", time.Minute*5)
	_ = c

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	<-signalChan
}
