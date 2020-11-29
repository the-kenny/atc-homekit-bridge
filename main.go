package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hc/log"
)

func main() {
	device := "hci0"

	// log.Debug.Enable()

	stateDirectory := os.Getenv("STATE_DIRECTORY")
	if stateDirectory != "" {
		os.Chdir(stateDirectory)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	termChan := make(chan int)

	ctx, shutdown := context.WithCancel(context.Background())

	thermometers := make(chan *thermometer)
	go runBluetooth(ctx, device, thermometers)

	homekit := makeHomekit()
	go homekit.Start(ctx)

	for {
		select {
		case thermometer := <-thermometers:
			// log.Info.Printf("Found thermometer %s", *thermometer)
			homekit.thermometerUpdates <- thermometer
		case <-termChan:
		case s := <-sig:
			log.Info.Printf("Received signal %s, stopping ", s)
			shutdown()
			return
		}
	}
}
