#!/bin/bash
set -e 


workspace=$(cd $(dirname $0);pwd)
cd $workspace

IMAGE=gosdk-discovery-client
TAG=latest

docker build -t $IMAGE:$TAG .
