#!/bin/sh

mkdir -p ./logs
touch ./logs/error.log
os=${OSTYPE//[0-9.]/}
if [ "$os" == "darwin" ]; then
    sudo nginx -p `pwd` -c conf/nginx_mac.conf
elif [ "$os" == "linux" ]; then
    sudo nginx -p `pwd` -c conf/nginx.conf
fi
