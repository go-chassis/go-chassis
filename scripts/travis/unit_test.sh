#!/bin/sh
set -e
# Make the Coverage File
echo "mode: atomic" > coverage.txt
# Make Necessary directories needed by Test (Ideally it should get created automatically but Travis is not allowing to create it using os.MkdriAll)
# I know this is insane but nothing can be done
mkdir -p root/conf
mkdir -p $GOPATH/src/github.com/ServiceComb/go-chassis/core/transport/tls
mkdir -p $GOPATH/src/github.com/ServiceComb/go-chassis/examples/discovery/server/log
mkdir -p $GOPATH/src/github.com/ServiceComb/go-chassis/examples/discovery/server/conf

# For auth test
mkdir -p $GOPATH/test/auth/conf
mkdir -p $GOPATH/test/auth/cipher
mkdir -p $GOPATH/test/auth/lib
mkdir -p $GOPATH/test/auth/log

# For transport test
mkdir -p $GOPATH/test/transport/TestCreateTransport/tls

mkdir -p /tmp/conf
mkdir -p $GOPATH/conf/microservice1/schema
mkdir -p $GOPATH/conf/microservice2/schema
mkdir -p $GOPATH/conf/microservice3/schema
#Start the Test
for d in $(go list ./... | grep -v vendor |  grep -v third_party | grep -v examples | grep -v metrics); do
    echo $d
    echo $GOPATH
    cd $GOPATH/src/$d
    if [ $(ls | grep _test.go | wc -l) -gt 0 ]; then
        go test -cover -covermode atomic -coverprofile coverage.out
        if [ -f coverage.out ]; then
            sed '1d;$d' coverage.out >> $GOPATH/src/github.com/ServiceComb/go-chassis/coverage.txt
        fi
    fi
done
