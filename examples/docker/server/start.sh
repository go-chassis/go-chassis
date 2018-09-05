#!/bin/sh
set -e

CURRENT_DIR=$(cd $(dirname $0);pwd)
APP_ROOT=$CURRENT_DIR

TMP_DIR="/tmp"
CONF_DIR="$APP_ROOT/conf"

copy_tmp2mesher() {
    if [ -f $TMP_DIR/$1 ]; then
        echo "$1 is customed"
        cp -f $TMP_DIR/$1 $CONF_DIR/$1
    fi
}

copy_tmp2mesher chassis.yaml
copy_tmp2mesher microservice.yaml
copy_tmp2mesher circuit_breaker.yaml
copy_tmp2mesher load_balancing.yaml
copy_tmp2mesher monitoring.yaml
copy_tmp2mesher lager.yaml
copy_tmp2mesher rate_limiting.yaml
copy_tmp2mesher tls.yaml
copy_tmp2mesher auth.yaml
copy_tmp2mesher tracing.yaml
copy_tmp2mesher router.yaml

$APP_ROOT/provider-mesher
