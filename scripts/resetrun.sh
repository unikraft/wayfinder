#!/bin/bash

ID=$(virsh list | grep running | awk '{ print $1 }')
UUID=$(virsh list | grep running | awk '{ print $2 }')

virsh destroy $ID
virsh undefine $UUID

INT_TO_UNIQUE=$(sudo ip link show | grep -ho "ethc[0-9]*")
INT_TO_CLEAR=$(echo -e "${INT_TO_UNIQUE// /\\n}" | sort -u)

for interface in $INT_TO_CLEAR; do
    sudo ip link set $interface down
    sudo ip link delete $interface
done
