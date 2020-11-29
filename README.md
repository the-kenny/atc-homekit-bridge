 # ATC HomeKit Bridge

 ## Setup

1. Set capabilities to run `atc-homekit-bridge` without root:

        sudo setcap 'cap_net_admin+ep' ./atc-homekit-bridge

2. Disable bluetooth service so `atc-homekit-bridge` can access the HCI socket

        sudo systemctl stop bluetooth
        sudo systemctl disable bluetooth
