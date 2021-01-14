package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	"gitlab.com/jtaimisto/bluewalker/hci"
	"gitlab.com/jtaimisto/bluewalker/host"
)

type thermometer struct {
	// Advertisement  []byte
	// Address        hci.BtAddress
	Mac            string
	Temperature    int16
	Humidity       uint8
	Battery        uint8
	BatteryVoltage uint16
	PacketCounter  uint8
}

func makeThermometer(scan *host.ScanReport) (*thermometer, error) {
	for _, data := range scan.Data {
		if data.Typ == hci.AdServiceData {

			if len(data.Data) != 15 {
				return nil, errors.New("Advertisement hat unorthodox length (expected 15)")
			}

			t := new(thermometer)

			d := data.Data

			// t.Advertisement = d
			// t.Address = scan.Address
			t.Mac = net.HardwareAddr.String(d[2:8])
			buf := bytes.NewBuffer(d[8:])
			binary.Read(buf, binary.BigEndian, &t.Temperature)
			binary.Read(buf, binary.BigEndian, &t.Humidity)
			binary.Read(buf, binary.BigEndian, &t.Battery)
			binary.Read(buf, binary.BigEndian, &t.BatteryVoltage)
			binary.Read(buf, binary.BigEndian, &t.PacketCounter)

			return t, nil
		}
	}

	return nil, errors.New("Found no AdServiceData")
}
