#!/usr/bin/env bash

docker pull servicecomb/service-center

docker run -d -p 30100:30100 --name=service-center  servicecomb/service-center:latest