 # ATC HomeKit Bridge

This project bridges Xiaomi Thermometers (LYWSD03MMC) running a [custom firmware](https://github.com/atc1441/ATC_MiThermometer) to Apple HomeKit. 

It uses [BlueWalker](https://gitlab.com/jtaimisto/bluewalker) code to open a raw HCI socket on Linux to read Bluetooth LE advertisements and parses data from above thermometers. This data is then published as HomeKit devices via [brutella/hc](https://github.com/brutella/hc).

Please note that the usage of the HCI socket needs exclusive access to the bluetooth stack. You can't use the normal Bluetooth tools while this application is running. It might be possible to use a second bluetooth adapter for normal bluetooth operations, but I haven't tried it.

 ## Building & Testing

These setup instructions use a Raspberry Pi running Raspbian. It should be similar for other devices running Linux.

1. Build the binary

        GOOS=linux GOARCH=arm GOARM=5 go build

    or, if you're compiling directly on the target platform

        go build

2. Set capabilities to run `atc-homekit-bridge` without root:

        sudo setcap 'cap_net_admin+ep' ./atc-homekit-bridge

3. Disable bluetooth service so `atc-homekit-bridge` can access the HCI socket

        sudo systemctl stop bluetooth
        sudo systemctl disable bluetooth

4. Run the binary

        ./atc-homekit-bridge

## Deployment

A `systemd` unit is provided. It expects `atc-homekit-bridge` in `/usr/local/bin/`. No other setup should be necessary. Note that `atc-homekit-bridge.service` conflicts with `bluetooth.service`. You need to disable `bluetooth.service` before `atc-homekit-bridge.service` can start. 