#!/bin/bash
set -e 
set -x 

BINARY_NAME="rest-service"
PKG_NAME="service.tar.gz"

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

cd $ROOT_PATH
if [ -f $BINARY_NAME ]; then
    rm $BINARY_NAME
fi

if [ -f $PKG_NAME ]; then
    rm $PKG_NAME
fi

go build -a -o "$BINARY_NAME"
tar -zcvf $PKG_NAME conf $BINARY_NAME start.sh

echo "Build success!"
