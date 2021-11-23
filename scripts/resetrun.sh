#!/bin/bash

ID=$(virsh list | grep running | awk '{ print $1 }')
UUID=$(virsh list | grep running | awk '{ print $2 }')

virsh destroy $ID
virsh undefine $UUID

ip link delete netnsv0-1

./dist/wfctl sj -l 1 2
