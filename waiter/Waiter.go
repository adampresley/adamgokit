package waiter

import (
	"os"
	"os/signal"
	"syscall"
)

/*
Wait will return a channel with an interrupt signal attached. Callers
must use "<- waiter.Wait()" to block.
*/
func Wait() chan os.Signal {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return quit
}
