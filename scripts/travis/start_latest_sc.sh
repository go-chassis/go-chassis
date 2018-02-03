#!/usr/bin/env bash
mkdir latestsc
cd latestsc
git clone https://github.com/apache/incubator-servicecomb-service-center.git
cd incubator-servicecomb-service-center
gvt restore
bash -x scripts/release/make_release.sh linux latest latest
bash -x servicecomb-service-center-latest-linux-amd64/start.sh


