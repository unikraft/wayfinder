#!/usr/bin/env bash

set -xe

apt-get install -y curl

BRIDGE=wayfinder$WAYFINDER_CORE_ID0 # create a unique bridge
BRIDGE_IP="172.${WAYFINDER_CORE_ID0}.${WAYFINDER_CORE_ID1}.1"
UNIKERNEL_IMAGE=${UNIKERNEL_IMAGE:-"nginx"}
UNIKERNEL_IP="172.${WAYFINDER_CORE_ID0}.${WAYFINDER_CORE_ID1}.2"
NUM_PARALLEL_CONNS=${NUM_PARALLEL_CONNS:-30}
DURATION=${DURATION:-30}

function cleanup {
  ifconfig $BRIDGE down || true
  brctl delbr $BRIDGE || true
  pkill qemu-system-x86_64 || true
}

trap "cleanup" EXIT

echo "Creating bridge..."
brctl addbr $BRIDGE || true
ifconfig $BRIDGE down
ifconfig $BRIDGE $BRIDGE_IP
ifconfig $BRIDGE up

echo "Starting unikernel..."
qemu-system-x86_64 \
  -machine q35 -display none -serial stdio -m 256M -machine accel=kvm:tcg  \
   -vga none -cpu host \
  -drive file=/root/.ops/images/nginx,format=raw,if=none,id=hd0 \
  -device pcie-root-port,port=0x10,chassis=1,id=pci.1,bus=pcie.0,multifunction=on \
  -device pcie-root-port,port=0x11,chassis=2,id=pci.2,bus=pcie.0 \
  -device pcie-root-port,port=0x12,chassis=3,id=pci.3,bus=pcie.0 \
  -device virtio-scsi-pci,bus=pci.2,id=scsi0 \
  -device scsi-hd,bus=scsi0.0,drive=hd0 \
  -device virtio-net,bus=pci.3,netdev=n0 \
  -netdev user,id=n0,hostfwd=tcp::8084-:8084

# make sure that the server has properly started
sleep 5

curl -Lvk http://$UNIKERNEL_IP:80

echo "Starting experiment..."
taskset -c $WAYFINDER_CORE_ID2 \
  wrk \
    -d $DURATION --latency \
    -t $NUM_PARALLEL_CONNS \
    -c $NUM_PARALLEL_CONNS http://$UNIKERNEL_IP:80/payload.txt &> /results.txt

cat /results.txt

echo "Done!"
