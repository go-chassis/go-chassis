SHELL := /bin/bash

all: binary docker
	# ./build./distribute-image.sh go-chassis/pilotv2client:latest
	-kubectl delete -f ./build/pilotv2client.dep.yaml
	kubectl apply -f ./build/pilotv2client.dep.yaml

binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w' -o ./build/pilotv2client

docker: binary
	cp -r ./conf ./build/
	docker build -t go-chassis/pilotv2client:latest ./build
