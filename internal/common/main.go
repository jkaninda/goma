package common

import (
	"github.com/getsentry/sentry-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Init() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, syscall.SIGINT)
	signal.Notify(interrupt, syscall.SIGTERM)
	return interrupt
}

func Defer() {
	sentry.Flush(time.Second * 5)
	sentry.Recover()
}
