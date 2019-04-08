package signal

import (
	"os"
	"os/signal"
)

func WatchSignal(cb func(os.Signal), sigs ...os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, sigs...)
	for s := range c {
		cb(s)
	}
}

