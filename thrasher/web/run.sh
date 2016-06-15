#!/bin/sh
      
if [ $# -le 0 ]; then
    echo $$>/tmp/testhub_web.pid
    exec node app.js
fi
