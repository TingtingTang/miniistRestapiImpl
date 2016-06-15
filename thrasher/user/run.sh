#!/bin/sh
      
SRC="main.go join.go login.go respond.go user.go helper.go"
if [ $# -le 0 ]; then
    echo $$>/tmp/testhub_user.pid
    go build $SRC
    if [ $? == 0 ]; then
	    exec ./main 
    else
            echo "Failed to build"
            exit -1
    fi
fi
