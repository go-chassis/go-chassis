#!/bin/bash
set -e 
set -x 

appname="gosdk-kube-server"

BUILD_PATH=$(cd $(dirname $0);pwd)
ROOT_PATH=$(cd $BUILD_PATH/..;pwd)
RELEASE_PATH=$ROOT_PATH/release

cd $ROOT_PATH

mkdir -p $RELEASE_PATH/$appname
rm -rf $RELEASE_PATH/*

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build --ldflags " -extldflags '-static'" -a -o "$RELEASE_PATH/$appname/app"

cp -rf conf $RELEASE_PATH/$appname
if [ -d "lib" ]; then
	cp -rf lib $RELEASE_PATH/$appname
fi

cd $RELEASE_PATH

package=$appname.tar.gz
tar -zcvf $package $appname

cp $RELEASE_PATH/$package $BUILD_PATH
bash $BUILD_PATH/build_image.sh

echo "build success!"
