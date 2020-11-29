package main

import (
	"context"
	"fmt"
	"os"

	"gitlab.com/jtaimisto/bluewalker/filter"
	"gitlab.com/jtaimisto/bluewalker/hci"
	"gitlab.com/jtaimisto/bluewalker/host"
)

func runBluetooth(ctx context.Context, device string, thermometerChan chan *thermometer) error {
	raw, err := hci.Raw(device)
	if err != nil {
		errorCritical(nil, fmt.Sprintf("Error while opening RAW HCI socket: %v\nAre you running as root and have you run sudo hciconfig %s down?", err, device))
	}

	host := host.New(raw)
	if err = host.Init(); err != nil {
		errorCritical(host, fmt.Sprintf("Unable to initialize host: %v", err))
	}

	filters := []filter.AdFilter{
		filter.ByPartialAddress([]byte{0xa4, 0xc1, 0x38}),
	}
	reportChan, err := host.StartScanning(false, filters)
	if err != nil {
		errorCritical(host, fmt.Sprintf("Unable to start scanning: %v", err))
	}

	for {
		select {
		case <-ctx.Done():
			break
		case sr := <-reportChan:
			// found := &FoundDevice{Structures: sr.Data,
			// 	Rssi:     sr.Rssi,
			// 	LastSeen: time.Now(),
			// 	Device:   sr.Address}

			// found.Types = []hci.AdvType{sr.Type}

			// fmt.Printf("device=%s\n", found)

			t, err := makeThermometer(sr)
			if err == nil {
				thermometerChan <- t
			} else {
				return err
			}
		}
	}

	host.Deinit()

	return nil
}

func errorMessage(message string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", message)
}

//error_critical will print given error message and terminate the program
// if host is non-nil, it will be deinitialized befor stoppping
func errorCritical(host *host.Host, message string) {
	errorMessage(message)
	if host != nil {
		host.Deinit()
	}

	os.Exit(255)
}
