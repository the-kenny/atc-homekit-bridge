 # ATC HomeKit Bridge

 ## Setup

1. Set capabilities to run `atc-bridge` without root:

    sudo setcap 'cap_net_admin+ep' ./atc-bridge

2. Disable bluetooth service so `atc-bridge` can access the HCI socket

    sudo systemctl stop bluetooth
    sudo systemctl disable bluetooth