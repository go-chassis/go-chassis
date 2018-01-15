#!/bin/sh
set -e 
set -x 

# //////////////////////////////////////////////////////////////////////////// #
#                               go sdk                                         #
# //////////////////////////////////////////////////////////////////////////// #

if [ "$CSE_SERVICE_CENTER" ]; then
    export CSE_REGISTRY_ADDR=$CSE_SERVICE_CENTER
fi

if [ "$CSE_CONFIG_CENTER" ]; then
    export CSE_CONFIG_CENTER_ADDR=$CSE_CONFIG_CENTER
fi

name=gosdk-discovery-client

listen_addr="0.0.0.0"
advertise_addr=$(ifconfig eth0 | grep -E 'inet\W' | grep -o -E [0-9]+.[0-9]+.[0-9]+.[0-9]+ | head -n 1)

cd /home/$name
# replace ip addr
sed -i s/"listenAddress:\s\{1,\}[0-9]\{1,3\}.[0-9]\{1,3\}.[0-9]\{1,3\}.[0-9]\{1,3\}"/"listenAddress: $listen_addr"/g conf/chassis.yaml
sed -i s/"advertiseAddress:\s\{1,\}[0-9]\{1,3\}.[0-9]\{1,3\}.[0-9]\{1,3\}.[0-9]\{1,3\}"/"advertiseAddress: $advertise_addr"/g conf/chassis.yaml

./app

while true; do
    sleep 60
done

