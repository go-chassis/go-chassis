#!/usr/bin/env bash
cd ../..
mkdir -p apache
cd apache
git clone https://github.com/apache/incubator-servicecomb-service-center.git
cd incubator-servicecomb-service-center
gvt restore
bash -x scripts/release/make_release.sh linux latest latest
cd servicecomb-service-center-latest-linux-amd64
bash -x start.sh

sleep 5

echo "Testing if Service-Center started succesfully"

curl -v http://127.0.0.1:30100/version

curl -v http://127.0.0.1:30100/v4/default/registry/version

curl -v http://127.0.0.1:30100/health