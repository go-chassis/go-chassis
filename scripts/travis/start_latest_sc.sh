#!/usr/bin/env bash

kill -9 $(lsof -i:30100 |awk '{print $2}' | tail -n 2)

echo "delete old contain  server-center ...."
docker rm -f service-center

echo "docker pull servicecomb/service-center images"
docker pull servicecomb/service-center

docker run -d -p 30100:30100 --name=service-center  servicecomb/service-center