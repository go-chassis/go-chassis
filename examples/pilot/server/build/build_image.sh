#!/bin/bash
set -e 


workspace=$(cd $(dirname $0);pwd)
cd $workspace

IMAGE=gosdk-istio-server
TAG=latest

docker build -t $IMAGE:$TAG .
