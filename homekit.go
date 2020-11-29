package main

import (
	"context"
	"strings"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/log"
	"github.com/brutella/hc/service"
)

type homekitThermometer struct {
	*accessory.Accessory

	temperatureSensor *service.TemperatureSensor
	humiditySensor    *service.HumiditySensor
	battery           *service.BatteryService
}

type homekit struct {
	accessories        map[string]*homekitThermometer
	transports         map[string]hc.Transport
	thermometerUpdates chan *thermometer
}

func (h *homekit) updateThermometerInternal(t *thermometer) *error {
	if h.accessories[t.Mac] == nil {
		log.Info.Printf("No Homekit Thermometer for %s found, creating...", t.Mac)
		config := hc.Config{Pin: "11112222"}

		shortMac := t.Mac
		shortMac = strings.TrimPrefix(shortMac, "a4:c1:38:")
		shortMac = strings.ReplaceAll(shortMac, ":", "")

		info := accessory.Info{
			ID:           1,
			Name:         "ATC " + shortMac,
			SerialNumber: t.Mac,
		}
		a := &homekitThermometer{}
		a.Accessory = accessory.New(info, accessory.TypeSensor)

		a.temperatureSensor = service.NewTemperatureSensor()
		a.temperatureSensor.CurrentTemperature.SetMinValue(-40)
		a.temperatureSensor.CurrentTemperature.SetMaxValue(100)
		a.temperatureSensor.CurrentTemperature.SetStepValue(0.1)
		a.Accessory.AddService(a.temperatureSensor.Service)

		a.humiditySensor = service.NewHumiditySensor()
		a.Accessory.AddService(a.humiditySensor.Service)

		a.battery = service.NewBatteryService()
		a.Accessory.AddService(a.battery.Service)

		transport, err := hc.NewIPTransport(config, a.Accessory)
		if err != nil {
			return &err
		}

		h.accessories[t.Mac] = a
		h.transports[t.Mac] = transport

		go transport.Start()
	}

	a := h.accessories[t.Mac]

	newBattery := int(t.Battery)
	if a.battery.BatteryLevel.GetValue() != newBattery {
		log.Info.Printf("Updating battery of thermometer %s to %d", t.Mac, newBattery)
		a.battery.BatteryLevel.SetValue(newBattery)
	}

	newTemperature := float64(t.Temperature) / 10.0
	if a.temperatureSensor.CurrentTemperature.GetValue() != newTemperature {
		log.Info.Printf("Updating temperature of thermometer %s to %f", t.Mac, newTemperature)
		a.temperatureSensor.CurrentTemperature.SetValue(newTemperature)
	}

	newHumidity := float64(t.Humidity)
	if a.humiditySensor.CurrentRelativeHumidity.GetValue() != newHumidity {
		log.Info.Printf("Updating humidity of thermometer %s to %f", t.Mac, newHumidity)
		a.humiditySensor.CurrentRelativeHumidity.SetValue(newHumidity)
	}

	return nil
}

func (h *homekit) Start(ctx context.Context) {
	for {
		select {
		case t := <-h.thermometerUpdates:
			err := h.updateThermometerInternal(t)
			if err != nil {
				log.Info.Fatal(*err)
			}

		case <-ctx.Done():
			log.Info.Println("Shutting down Homekit")
			for _, t := range h.transports {
				<-t.Stop()
			}
			break
		}
	}
}

func makeHomekit() *homekit {
	return &homekit{
		accessories:        make(map[string]*homekitThermometer),
		transports:         make(map[string]hc.Transport),
		thermometerUpdates: make(chan *thermometer),
	}
}
