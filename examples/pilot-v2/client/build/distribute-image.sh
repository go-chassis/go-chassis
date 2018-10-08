#!/bin/bash
imagename=$1
if [ "$imagename" == "" ]; then
	echo "need an image name"
	exit 1
fi

imagefilename=$(echo $imagename | sed 's/\//_/g')

# nodeIPs=$(kubectl get nodes -o wide | grep -v 'NAME.*STATUS' | awk '{print $1}')
nodeIPs=$(kubectl get nodes -o jsonpath='{$.items[*].status.addresses[?(@.type=="ExternalIP")].address}')
if [ "$nodeIPs" == "" ]; then
	nodeIPs=$(kubectl get nodes -o jsonpath='{$.items[*].status.addresses[?(@.type=="InternalIP")].address}')
fi
if [ "$nodeIPs" == "" ]; then
	nodeIPs=$(kubectl get nodes -o jsonpath='{$.items[*].status.addresses[?(@.type=="Hostname")].address}')
fi

if [ "$nodeIPs" == "" ]; then
	echo "Failed to get nodes' IP"
	exit 1
fi

echo saving to $imagefilename.tar
docker save $imagename > /tmp/$imagefilename.tar

for nodeip in ${nodeIPs[@]}
do
	echo distributing $imagename "to" $nodeip
	scp /tmp/$imagefilename.tar root@$nodeip:/root/
	ssh root@$nodeip "docker load < /root/$imagefilename.tar"
	ssh root@$nodeip "rm /root/$imagefilename.tar"
done

rm /tmp/$imagefilename.tar
