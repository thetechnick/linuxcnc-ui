package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/thetechnick/linuxcnc-ui/internal/linuxcnc"
)

func main() {
	c, err := linuxcnc.NewClient()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	p := linuxcnc.NewStatusPoller(c, 500*time.Millisecond)

	serverStopCh := make(chan struct{})
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		<-signalCh
		close(serverStopCh)
	}()

	go func() {
		for range p.OnChange() {
			status := p.Status()
			fmt.Printf("%#v\n", status)
		}
	}()

	if err := p.Run(serverStopCh); err != nil {
		panic(err)
	}
}
