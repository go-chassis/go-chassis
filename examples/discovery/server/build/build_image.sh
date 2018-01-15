#!/bin/bash
set -e 
set -x 

workspace=$(cd $(dirname $0);pwd)
cd $workspace

IMAGE=gosdk-discovery-server
TAG=latest

docker build -t $IMAGE:$TAG .

YUNLONG_IMAGE=csegosdk:latest

docker tag $IMAGE:$TAG $YUNLONG_IMAGE
