#!/bin/sh
      
SRC="main.go testsuite.go helper.go"
if [ $# -le 0 ]; then
    echo $$>/tmp/testhub_testsuite.pid
    go build $SRC
    if [ $? == 0 ]; then
	    exec ./main 
    else
            echo "Failed to build"
            exit -1
    fi
fi
